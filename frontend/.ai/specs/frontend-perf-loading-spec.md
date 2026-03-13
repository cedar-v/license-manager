# 技术契约与架构设计书：前端加载性能优化

**版本**：v1.0  
**日期**：2026-03-11  
**作者**：system-architect  
**状态**：待 task-coder 实施

---

## 一、问题诊断（根因分析）

### 1.1 每次页面跳转都重新加载 Element Plus CSS

**根因**：`vite.config.ts` 中 `ElementPlusResolver({ importStyle: 'css' })` 配置让 `unplugin-vue-components` 在每个使用到 Element Plus 组件的 `.vue` 文件中注入独立的 CSS import（`el-dropdown-menu.css`、`el-button-group.css` 等）。与此同时 `main.ts` 还保留了 `import 'element-plus/dist/index.css'`（全量 CSS），两者并存形成**双重加载**。

**表现**：网络面板中看到 `el-dropdown-menu.css`、`el-button-group.css`、`element_plus_es_components_loading_s...` 等碎片 CSS 在每次路由切换时被重新请求。

**正确做法**：二选一——
- **方案 A（推荐）**：继续使用 `element-plus/dist/index.css` 全量 CSS（已在 main.ts 中），关闭 resolver 的 `importStyle`。全量 CSS 首次加载后长期缓存，零碎的 per-component CSS 完全消失。
- 方案 B：彻底删除 `element-plus/dist/index.css`，完全依赖 resolver 按需注入。需要确保每个组件都被 resolver 覆盖，否则样式丢失。

### 1.2 字体每次导航重复加载（视觉上的误判）

**根因**：浏览器缓存策略问题。`@fontsource/noto-sans-sc` 输出的 woff2 文件名**不含内容 hash**（如 `noto-sans-sc-116-400-normal.woff2`），而 Vite dev 服务器对字体文件**不设置强缓存 header**（`Cache-Control: no-cache`）。结果是浏览器每次导航都发起带 If-None-Match 的条件请求，即使返回 304，devtools 也会显示为"加载"。

**生产环境额外问题**：vite.config.ts 的 `assetFileNames` 对字体匹配正确会加 hash，但 `@fontsource` 的 CSS 中 `src` 路径是相对路径，构建时 Vite 自动处理后会有 hash，缓存正常。

### 1.3 PangMenZhengDao.ttf：1004 KB，仅 LicenseSearch.vue 使用

这是最严重的单点问题。一个装饰性标题字体（仅 `LicenseSearch.vue` 的 `.platform-title` 用到）以 **1 MB 的 TTF 格式**全局加载。

- 全局加载：`main.ts` → `fonts.css` → 预加载所有字体
- 即使用户从未访问授权查询页，这 1 MB 也会被下载

### 1.4 Noto Sans SC：两个字重 = 每字重 ~15 个 woff2 分片

400 + 700 两个字重，每次页面含中文时浏览器匹配 unicode-range 触发下载，共约 20-30 个分片文件（每个 25-33 KB）。这是 CJK 字体的固有代价，无法消除，但可以：
- 替换为系统字体（零网络开销）
- 转为自托管 WOFF2 子集（仅保留常用汉字）

---

## 二、优化方案设计

### 方案一（必做）：移除 ElementPlusResolver CSS 双重加载

**影响层**：构建配置层  
**改动意图**：将 `importStyle: 'css'` 改为 `importStyle: false`，消除每个组件的碎片 CSS 请求，改由全局 `element-plus/dist/index.css` 统一提供样式。

**数据收益**：消除约 10-20 个 CSS 请求（每次路由跳转）。

---

### 方案二（必做）：PangMenZhengDao 懒加载隔离

**影响层**：CSS 层 + 字体加载策略层  
**改动意图**：
1. 从 `fonts.css` 中删除 `@font-face` 全局声明
2. 将 `@font-face` 声明移入 `LicenseSearch.vue` 的 `<style scoped>`，利用路由懒加载天然隔离，仅在用户访问该页时触发字体下载

**数据收益**：首屏及非授权查询页节省 1 MB 传输。

---

### 方案三（推荐）：Noto Sans SC → 系统字体降级

**影响层**：CSS 变量层（variables.scss）  
**改动意图**：将 `$font-family-primary` 中的 `'Noto Sans SC'` 降级为第二选择，保留兜底。同时评估是否保留 `@fontsource` 导入（影响跨平台一致性）。

如果要保留字体一致性，最小化方案是：
- **只导入 400**（删除 700 的导入），weight:700 的文字交给浏览器合成（faux-bold），视觉差异极小
- 用 `font-display: optional` 替代 `swap` 以避免 FOUT

**数据收益**：减少 ~50% 的字体请求数（约 15 个）。

---

### 方案四（加分）：开发服务器字体缓存优化

**影响层**：vite.config.ts server 配置  
**改动意图**：为 dev server 添加字体文件 `Cache-Control` 响应头，避免浏览器 devtools 中字体反复显示为"加载中"的视觉误导。

---

## 三、数据模型层（Schema）

无数据库变更。

---

## 四、后端服务层（API）

无需改动。

---

## 五、前端交互层（View）— 目标文件清单

| # | 文件路径 | 操作 | 改动说明 |
|---|---------|------|---------|
| 1 | `frontend/vite.config.ts` | **修改** | `ElementPlusResolver` 的 `importStyle: 'css'` → `false`；添加 dev server font cache header |
| 2 | `frontend/src/assets/styles/fonts.css` | **修改** | 删除 `PangMenZhengDao` 的 `@font-face` 块 |
| 3 | `frontend/src/views/Licenses/LicenseSearch.vue` | **修改** | 在 `<style scoped>` 顶部加入 `PangMenZhengDao` 的 `@font-face` 声明 |
| 4 | `frontend/src/main.ts` | **修改** | 删除 `@fontsource/noto-sans-sc/700.css` 导入（可选，评估后决定）|

**共 4 个文件，不超过 task-coder 单批限制（5个）。**

---

## 六、实施优先级

```
P0（必做，无副作用）:
  - 文件1：ElementPlusResolver importStyle fix
  - 文件2 + 文件3：PangMenZhengDao 懒加载隔离

P1（可选，有视觉权衡）:
  - 文件4：删除 700 字重导入
```

---

## 七、预期收益对比

| 指标 | 优化前 | 优化后（P0） |
|------|-------|------------|
| 路由跳转触发的 CSS 请求 | ~15 个碎片 CSS | 0 |
| 非授权查询页的字体传输 | 1 MB (PangMenZhengDao) | 0 |
| 首次访问总字体请求数 | ~35 个 | ~20 个 |

---

*→ 移交 `task-coder` 按上述文件清单实施。*
