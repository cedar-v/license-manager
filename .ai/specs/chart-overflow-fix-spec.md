# 技术契约与架构设计书
## 功能：授权趋势图表自适应溢出修复

**版本**：v1.0  
**日期**：2026-03-13  
**设计者**：System Architect  
**状态**：待实现

---

## 一、问题定位（Root Cause Analysis）

通过审计截图与源码，确认存在 **三条独立的溢出根因**，相互叠加导致图表在不同视口宽度下随机超出卡片边界。

### 根因 A — CSS 高度链断裂（最高优先级）

**文件**：`frontend/src/views/Dashboard.vue`

```
.content-container { min-height: calc(100vh - 80px); }   ← ①
  └─ .dashboard     { height: 100%; }                     ← ② ❌ 100% of min-height = 0！
       └─ .content-section { flex: 1; min-height: 0; }    ← ③
            └─ .chart-section { flex: 0.5; min-height: 300px; }
                 └─ .license-trend-chart { height: 100%; } ← ④ 同样失效
```

CSS 规范：`height: 100%` 只解析相对于父元素的 **`height`** 属性，而非 `min-height`。步骤②处 `.dashboard` 获取的 `height: 100%` 实际解析为 `auto`，导致整条 flex 高度链全部退化为 `auto`/内容高度，图表实际渲染高度由内容撑开，触发溢出。

**现象**：桌面端图表卡片高度不稳定，ECharts canvas 随视口/数据量变化而溢出卡片底部。

### 根因 B — ECharts Grid 使用绝对像素 + `containLabel: false`

**文件**：`frontend/src/components/charts/LicenseTrendChart.vue`

```js
grid: {
  left: 60,         // 固定像素，容器小时 Y 轴标签溢出左侧
  right: 50,        // 固定像素，容器小时溢出右侧
  top: 30,
  bottom: 60,
  containLabel: false   // ❌ 不包含坐标轴标签在 grid 范围内
}
```

当容器宽度较小时，ECharts 按绝对像素保留边距并绘制轴标签于 grid 外部，标签实际渲染于 canvas 边缘外，出现截断/溢出。

**现象**：窄屏或侧边栏展开时，X 轴日期标签被裁切或与右侧内容重叠。

### 根因 C — 父容器 `overflow: hidden` 裁切 ECharts tooltip

**文件**：`frontend/src/components/charts/LicenseTrendChart.vue`

```scss
.license-trend-chart {
  overflow: hidden;   // ❌ 裁切了 ECharts tooltip 的绝对定位弹层
}
```

ECharts 的 tooltip 默认渲染为 `position: absolute` 的 `div` 元素，若其最近的 `overflow: hidden` 祖先正好是卡片容器，tooltip 悬浮层会被硬裁。移动端尤其明显（截图二：tooltip "授权数量: 8" 被截断于卡片底部）。

---

## 二、数据模型层（Schema）

本次修复为纯前端 CSS/布局层变更，**不涉及任何数据模型、API 接口或后端逻辑修改**。API 契约保持不变。

---

## 三、后端服务层（API）

**无变更**。后端 `GET /api/v1/dashboard/authorization-trend` 接口及响应格式 `AuthorizationTrendResponse` 不受影响。

---

## 四、前端交互层（View）变更契约

### 4.1 变更文件清单

| 操作 | 文件路径 | 变更范围 |
|------|----------|----------|
| 修改 | `frontend/src/views/Dashboard.vue` | `.content-container` 的高度策略 |
| 修改 | `frontend/src/components/charts/LicenseTrendChart.vue` | Grid 配置 + overflow + 高度策略 |

---

### 4.2 `frontend/src/views/Dashboard.vue` — 修改意图

**目标**：修复高度链断裂，使 flex 容器的高度能正确传递给子组件。

#### 改动 1：`.content-container` 改用 `height` 替代 `min-height`

```scss
// ❌ 修改前
.content-container {
  min-height: calc(100vh - 80px);
  // ...
}

// ✅ 修改后
.content-container {
  height: calc(100vh - 80px);   // 使用确定的 height，让 flex 子项 height: 100% 可以解析
  min-height: 600px;            // 保留最小高度防止内容压扁，但用 min-height 做限底
  // ...
}
```

> **风险评估**：`Dashboard.vue` 的 `height` 被 `Layout` 组件 `overflow: auto` 的滚动容器包裹，切换为固定 `height` 不会影响滚动行为，但需验证 Layout 组件是否已设置 `overflow-y: auto`。

---

### 4.3 `frontend/src/components/charts/LicenseTrendChart.vue` — 修改意图

**目标 A**：修复 ECharts Grid overflow（根因 B）。

#### 改动 1：`chartOption.grid` 启用 `containLabel: true` 并将边距改为百分比

```js
// ❌ 修改前
grid: {
  left: 60,
  right: 50,
  top: 30,
  bottom: 60,
  containLabel: false
}

// ✅ 修改后
grid: {
  left: '3%',        // 改为百分比，随容器缩放
  right: '3%',
  top: 30,
  bottom: 30,        // containLabel: true 后 ECharts 自动为标签留空间，无需预留 60px
  containLabel: true // 轴标签始终包含在 grid 区域内，不会超出 canvas
}
```

**目标 B**：解决 tooltip 被 `overflow: hidden` 裁切（根因 C）。

#### 改动 2：`v-chart` 的 tooltip 改为挂载到 body

在 `chartOption.tooltip` 中追加：

```js
tooltip: {
  // ...现有配置不变...
  appendToBody: true,   // tooltip DOM 挂载到 <body>，完全逃离 overflow: hidden 容器
}
```

> 此方案比改 `overflow: visible` 更稳健，不破坏卡片的 `border-radius` 裁切效果。`overflow: hidden` 保留，仅 tooltip 逃逸。

**目标 C**：稳定图表 canvas 在 flex 链中的高度（根因 A 的组件侧保障）。

#### 改动 3：`.chart-container` 添加 `position: relative` + `.trend-chart` 使用位置绝对铺满

为防止在高度链无法修复的回退场景下图表仍能铺满容器：

```scss
// ✅ 修改后
.chart-container {
  padding: 1.35vw 1.25vw 1.25vw;
  flex: 1;
  display: flex;
  min-height: 0;
  position: relative;            // 新增：建立定位上下文

  .trend-chart {
    width: 100%;
    height: 100%;
    min-height: 12.8vw;          // 保持现有移动端 fallback
  }
}

// 移动端覆盖（已有，维持原样）
@media (max-width: 1024px) {
  .chart-container .trend-chart {
    height: 200px;
    min-height: unset;
  }
}
```

---

## 五、全栈影响面分析（Full-Stack Impact Analysis）

| 层级 | 受影响组件 | 影响类型 | 处理方式 |
|------|-----------|----------|----------|
| 数据层 | 无 | — | 无需处理 |
| API 层 | 无 | — | 无需处理 |
| 状态层 | 无 | — | 无需处理 |
| 视图组件 | `Dashboard.vue` | 高度策略调整 | 必须修改 |
| 子组件 | `LicenseTrendChart.vue` | Grid 配置 + tooltip + overflow | 必须修改 |
| 其他页面 | 其他使用 `Layout` 组件的页面 | 可能受 `.content-container` 修改影响 | 修改前审查；若 `.content-container` 属全局 class 则需回归测试 |

> ⚠️ **注意**：如果 `.content-container` 是全局共用 class（定义在全局 CSS 而非 scoped），则修改会影响所有使用该 class 的页面。需先确认作用域：当前 Dashboard.vue 的 `<style lang="scss" scoped>` 限定了 `.content-container` 只作用于该页面，因此是安全的局部修改。

---

## 六、待验证事项（Acceptance Criteria）

- [ ] 1920×1080 桌面端：图表无溢出，X 轴日期标签完整显示
- [ ] 1440×900：图表无溢出，轴标签不被截断
- [ ] 1280×800：图表卡片高度稳定，边距正常
- [ ] 768px 平板：`.chart-section min-height: 300px` 生效，图表可读
- [ ] 375px 手机：固定 200px 高度正常，tooltip 不被卡片边界截断
- [ ] 暗模式：tooltip 背景色与轴标签颜色不变
- [ ] 快捷选择切换 7 天 / 30 天：图表 resize 正常，无残影溢出

---

## 七、文件清单汇总

### 必须修改
1. `frontend/src/views/Dashboard.vue`
   - SCSS `.content-container`：`min-height` → `height` + 添加 `min-height` 保底

2. `frontend/src/components/charts/LicenseTrendChart.vue`  
   - JS `chartOption.grid`：`containLabel: false` + 固定 px → `containLabel: true` + `%` 边距
   - JS `chartOption.tooltip`：新增 `appendToBody: true`
   - SCSS `.chart-container`：新增 `position: relative`

### 无需新建
本次修复不需要创建任何新文件。

---

*本文档由 System Architect 技能生成，作为 task-coder 的实现输入。实现时严格遵循上述契约，禁止超出本文档定义的修改范围。*
