<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-11 09:07:12
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-09-03 10:24:43
 * @FilePath: /frontend/docs/Claude + Figma 开发 Vue 3 + TypeScript + Element Plus 应用指南.md
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
# Claude + Figma 开发 Vue 3 + TypeScript + Element Plus 应用指南

## 1. 准备工作

### 1.1 工具安装

**必需工具:**
- **Figma**: 设计工具，用于获取UI设计资源和导出图标、图片资源
- **Node.js**: 确保安装最新LTS版本 (推荐18.x或更高版本)
- **包管理器**: npm、yarn 或 pnpm (推荐使用 pnpm)
- **Vue CLI 或 Vite**: 项目脚手架工具 (推荐使用 Vite，构建速度更快)
- **Git**: 版本控制工具
- **Claude**: AI助手，用于代码生成和开发协助

**推荐工具:**
- **VS Code**: 代码编辑器
- **Vue DevTools**: 浏览器调试插件
- **Figma Desktop**: 桌面版 Figma 应用
- **Element Plus Helper**: VS Code 插件，提供 Element Plus 组件智能提示

### 1.2 环境配置

**Node.js 环境检查:**
```bash
node --version  # 应显示 v18.x.x 或更高
npm --version   # 检查 npm 版本
```

## 2. 项目结构说明

### 2.1 完整项目目录结构

```
license-manager-frontend/
├── .claude/                        # Claude Code 配置
│   └── settings.local.json        # 本地设置
├── .vscode/                        # VS Code 配置
│   └── extensions.json            # 推荐插件
├── docs/                           # 项目文档
│   ├── Claude + Figma 开发 Vue 3 + TypeScript + Element Plus 应用指南.md
│   └── 多语言框架设计方案.md
├── public/                         # 静态资源
│   └── favicon.ico                # 网站图标
├── src/                           # 源码目录
│   ├── api/                       # API 接口层
│   │   ├── https/                 # HTTP 请求工具
│   │   │   ├── errorCodeType.ts   # 错误码类型定义
│   │   │   └── index.ts          # HTTP 请求封装
│   │   └── user.ts               # 用户相关 API
│   ├── assets/                    # 静态资源
│   │   ├── icons/                 # SVG 图标
│   │   ├── images/                # 图片资源
│   │   │   ├── login-background.png      # 登录背景图
│   │   │   └── login-background-m.png    # 移动端登录背景
│   │   └── styles/                # 样式文件系统
│   │       ├── global.scss        # 全局样式重置和基础样式
│   │       ├── variables.scss     # SCSS 变量（颜色、尺寸等）
│   │       ├── element-theme.scss # Element Plus 主题定制
│   │       └── global.css         # 旧版全局样式（待清理）
│   ├── components/                # 组件目录
│   │   ├── common/                # 通用组件
│   │   │   └── layout/            # 布局组件
│   │   │       ├── Layout.vue     # 主布局组件
│   │   │       ├── Sidebar.vue    # 侧边栏组件
│   │   │       ├── NavContent.vue # 顶部导航组件
│   │   │       └── index.ts       # 类型定义和导出
│   │   └── business/              # 业务相关组件
│   │       └── cusCard/           # 自定义卡片组件
│   ├── i18n/                      # 国际化配置
│   │   ├── index.ts               # i18n 初始化配置
│   │   └── locales/               # 语言包
│   │       ├── zh.json            # 中文语言包
│   │       ├── en.json            # 英文语言包
│   │       └── ja.json            # 日文语言包
│   ├── router/                    # 路由配置
│   │   └── index.ts               # Vue Router 配置
│   ├── store/                     # 状态管理 (Pinia)
│   │   ├── index.ts               # Pinia 配置和导出
│   │   └── modules/               # 状态模块
│   │       ├── app.ts             # 应用全局状态
│   │       └── user.ts            # 用户状态管理
│   ├── utils/                     # 工具函数库
│   ├── views/                     # 页面组件
│   │   ├── Dashboard.vue          # 仪表盘页面
│   │   └── Login.vue              # 登录页面
│   ├── App.vue                    # 根组件
│   ├── main.ts                    # 应用入口文件
│   └── vite-env.d.ts              # Vite 环境类型声明
├── .env.development               # 开发环境变量
├── .env.production                # 生产环境变量
├── .env.test                      # 测试环境变量
├── .eslintrc.cjs                  # ESLint 配置
├── .prettierrc                    # Prettier 代码格式化配置
├── index.html                     # HTML 入口文件
├── package.json                   # 项目依赖和脚本配置
├── tsconfig.json                  # TypeScript 配置
├── tsconfig.node.json             # Node.js TypeScript 配置
└── vite.config.ts                 # Vite 构建配置
```

### 2.2 核心目录详解

#### 📁 src/assets/styles/ - SCSS 样式系统
- **variables.scss**: 定义主题色彩、字体大小、间距等设计变量
- **mixins.scss**: 包含常用的样式混合器（清除浮动、文本省略、响应式等）
- **global.scss**: 全局样式重置和基础样式定义
- **element-theme.scss**: Element Plus 组件深度样式定制

#### 📁 src/components/ - 组件分层架构
- **common/**: 跨项目可复用的通用组件（如布局、表单、表格等）
- **business/**: 特定业务场景的组件（与许可证管理相关）

#### 📁 src/store/modules/ - Pinia 状态管理
- **app.ts**: 应用级状态（主题、语言、侧边栏状态等）
- **user.ts**: 用户状态（登录信息、权限、个人设置等）

#### 📁 配置文件说明
- **.eslintrc.cjs**: ESLint 代码规范检查配置
- **.prettierrc**: 代码格式化规则配置
- **vite.config.ts**: Vite 构建工具配置（包含性能优化设置）

### 2.3 样式文件组织最佳实践

#### SCSS 样式系统架构
```scss
// 1. 在 variables.scss 中定义设计变量
$primary-color: #019C7C;
$border-radius: 6px;

// 2. 在 mixins.scss 中创建可复用的样式片段
@mixin button-style($color: $primary-color) {
  background-color: $color;
  border-radius: $border-radius;
  // ...
}

// 3. 在组件中导入和使用
@import '@/assets/styles/variables.scss';
@import '@/assets/styles/mixins.scss';

.my-button {
  @include button-style($success-color);
}
```

#### 组件样式规范
- **优先使用** `<style lang="scss" scoped>` 避免样式污染
- **导入变量**: 在组件中导入 variables.scss 和 mixins.scss 
- **复用样式**: 将常用样式抽取到 mixins 中
- **主题支持**: 使用 SCSS 变量配合 Element Plus 主题定制


**一、整体响应式方案（不使用rem）:**

1. **JavaScript 辅助监听**: 在 `src/utils/resize.js` 等文件中监听窗口 resize 事件，动态计算屏幕宽度
2. **设备类型管理**: 将设备类型（desktop/mobile）等信息存入 Pinia，供组件使用
3. **CSS 媒体查询**: 定义不同屏幕宽度下的样式规则（核心方案）
4. **Flex 布局和百分比宽度**: 使用 Flex 布局实现内容弹性流动，使用 % 或 vw/vh 定义容器宽度

**二、移动端体验深度优化:**

1. **点击目标大小**: 确保按钮、链接等可点击元素尺寸不小于 44x44px（苹果HIG标准），适应手指触摸，避免误操作

2. **避免悬停（Hover）事件**: 移动端没有鼠标悬停，所有依赖 `:hover` 状态的功能（工具提示、二级菜单）都必须提供替代交互方式，如点击触发或长按触发

3. **原生控件和滚动**: 
   - 使用 `<input type="date">` 调用原生日期选择器，体验优于自定义组件
   - 处理 `-webkit-overflow-scrolling: touch` 启用iOS弹性滚动，让列表滚动更顺滑

4. **视口（Viewport）设置**: 确保 `<meta name="viewport">` 标签正确设置，根据需求决定是否禁止缩放以防止布局错乱

**三、复杂组件的响应策略:**

**1. 表格（Table）:**
   - **横向滚动**: 为表格容器设置 `overflow-x: auto`，允许用户横向滑动查看隐藏列
   - **列省略/优先级**: 隐藏次要列，只保留关键信息
   - **卡片化视图**: 小屏幕下（特别是移动端），将表格行转换为独立卡片，垂直展示数据

**2. 表单（Form）:**
   - **布局调整**: 将 label 和 input 的左右布局改为上下（block）布局

**3. 图表（Charts）:**
   - **动态重绘**: 使用 ECharts 等库的 `resize()` 方法，在容器尺寸变化时重绘图表
   - **简化显示**: 复杂图表在小屏幕下隐藏图例、简化配置，或提示用户"请在更大屏幕上查看此图表"

**4. 模态框（Modal/Dialog）:**
   - **全屏显示**: 移动端模态框应全屏显示，而非居中弹出浮层，更符合移动端交互习惯

**四、性能优化:**

1. **防抖和节流**: 对 resize、scroll 等频繁触发事件使用防抖或节流，避免性能浪费
2. **按需加载和代码分割**: 使用 Vue 异步组件和 Webpack 动态 `import()` 语法，将不同设备的组件代码分别加载
3. **图片和媒体优化**: 使用 `srcset` 和 `sizes` 属性为不同屏幕提供不同分辨率图片，避免移动端加载桌面大图

## 3. Figma设计转代码

### 3.1 设计标记获取流程

**步骤1：获取设计标记**
1. 在 Figma 中选择设计稿的根节点
2. 查看右侧面板的"设计"标签获取：
   - 颜色值 (Color)
   - 字体规格 (Typography)
   - 间距系统 (Spacing)
   - 圆角值 (Border radius)
   - 阴影效果 (Drop shadow)

**步骤2：导出设计资源**
```bash
# Figma 资源导出建议
1. 图标导出为 SVG 格式
2. 图片导出为 PNG/WebP 格式 (2x 倍图)
3. 获取精确的间距、字号、颜色值
```

### 3.3 Element Plus 主题定制

**主题配置 (`/src/assets/styles/element-theme.scss`):**
```scss
// 覆盖 Element Plus CSS 变量
:root {
  --el-color-primary: #{$primary-color};
  --el-border-radius-base: #{$border-radius-base};
  --el-font-size-base: #{$font-size-base};
  // 更多变量覆盖...
}

// 深度选择器示例
:deep(.el-button) {
  &--primary {
    background-color: $primary-color;
    border-color: $primary-color;
    
    &:hover {
      background-color: lighten($primary-color, 10%);
      border-color: lighten($primary-color, 10%);
    }
  }
}
```

### 3.4 组件样式实现示例

**按钮组件定制:**
```vue
<template>
  <el-button 
    :type="type" 
    :size="size" 
    :loading="loading"
    class="custom-button"
  >
    <slot></slot>
  </el-button>
</template>

<style lang="scss" scoped>
.custom-button {
  border-radius: $border-radius-base;
  font-weight: $font-weight-primary;
  
  // 根据 Figma 设计稿调整
  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.12);
  }
}
</style>
```
### 3.5 图标处理方案
**优先使用element-plus的icon图标库,可以在element-plus的官方文档中找到所有的图标,并将其引入到项目中**
```bash

**自定义 SVG 图标:在还原设计图的情况下如果element-plus的icon图标库不满足需求,可以将设计图中的icon下载为svg,**
```vue
<template>
  <svg-icon name="custom-icon" class="icon" />
</template>

<script setup lang="ts">
// 创建全局 SVG 图标组件
defineProps<{
  name: string
  size?: string
}>()
</script>
```

### 3.6 公共样式抽取

**工具类样式 (`/src/assets/styles/global.scss`):**


## 4. 开发流程和最佳实践

### 4.1 Claude 协作开发流程

**步骤1：需求分析**
```
1. 向 Claude 描述功能需求
2. 提供 Figma 设计稿链接或截图
3. 说明技术栈要求（Vue 3 + TypeScript + Element Plus）
```

**步骤2：代码生成**
```
提问示例：
"请根据这个 Figma 设计稿生成一个用户管理页面，包含：
- 用户列表表格
- 搜索和筛选功能
- 新增/编辑用户弹窗
- 使用 Element Plus 组件库"
```

**步骤3：代码审查与优化**
```
代码审查要点：
1. TypeScript 类型定义是否完整
2. 组件设计是否合理
3. 样式是否符合设计规范
4. 性能优化是否到位
```

### 4.2 组件开发最佳实践

**组件设计原则:**
```vue
<script setup lang="ts">
// 1. 定义清晰的 Props 接口
interface Props {
  title: string
  data: any[]
  loading?: boolean
}

// 2. 使用 withDefaults 设置默认值
const props = withDefaults(defineProps<Props>(), {
  loading: false
})

// 3. 定义明确的 Emit 事件
const emit = defineEmits<{
  update: [value: any]
  delete: [id: string]
}>()

// 4. 使用 Composition API 组织逻辑
const { data, loading, error } = useUserData()
</script>

<template>
  <!-- 5. 合理的模板结构 -->
  <div class="user-table">
    <el-table :data="data" :loading="loading">
      <!-- ... -->
    </el-table>
  </div>
</template>

<style lang="scss" scoped>
// 6. 作用域样式
.user-table {
  // 样式定义
}
</style>
```
### 4.3 样式开发最佳实践

**响应式设计:**
```scss
// 响应式断点
$breakpoint-mobile: 768px;
$breakpoint-tablet: 992px;
$breakpoint-desktop: 1200px;

// 响应式混入
@mixin mobile {
  @media (max-width: #{$breakpoint-mobile - 1px}) {
    @content;
  }
}

@mixin tablet {
  @media (min-width: #{$breakpoint-mobile}) and (max-width: #{$breakpoint-tablet - 1px}) {
    @content;
  }
}

@mixin desktop {
  @media (min-width: #{$breakpoint-desktop}) {
    @content;
  }
}

// 使用示例
.user-card {
  padding: 16px;
  
  @include mobile {
    padding: 8px;
  }
  
  @include desktop {
    padding: 24px;
  }
}
```

**主题切换支持:**
```scss
// 支持暗色模式
:root {
  --bg-color: #ffffff;
  --text-color: #333333;
}

[data-theme="dark"] {
  --bg-color: #1a1a1a;
  --text-color: #ffffff;
}

.app-container {
  background-color: var(--bg-color);
  color: var(--text-color);
}
```

### 4.5 性能优化建议

**组件懒加载:**
```typescript
// router/index.ts
const routes = [
  {
    path: '/users',
    name: 'UserManagement',
    component: () => import('@/views/UserManagement.vue')
  }
]
```

**图片懒加载和优化:**
```vue
<template>
  <el-image
    :src="imageUrl"
    fit="cover"
    lazy
    :preview-src-list="[imageUrl]"
    :placeholder="placeholderImage"
  />
</template>
```

## 5. 常见问题和解决方案







**问题3：SCSS 变量未定义**
```scss
// 解决方案：在 vite.config.ts 中配置全局 SCSS 变量
export default defineConfig({
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/assets/styles/variables.scss";`
      }
    }
  }
})
```
### 5.3 样式开发问题

**问题1：移动端适配问题**
```scss
// 解决方案：使用视口单位和媒体查询
.container {
  width: 100vw;
  margin: 0 auto;
  padding: 0 4vw;
  
  @media (max-width: 768px) {
    padding: 0 5vw;
  }
}

// 字体大小适配
.title {
  font-size: clamp(1.5rem, 4vw, 2.5rem);
}
```

**问题2：Figma导出的颜色值与实际不符**
```scss
// 解决方案：检查色彩空间和透明度
// Figma 中显示 rgba(255, 0, 0, 0.5)
// 可能需要转换为
.element {
  background-color: rgba(255, 0, 0, 0.5);
  /* 或者使用 HSL */
  background-color: hsla(0, 100%, 50%, 0.5);
}
```

**问题3：Element Plus 主题定制不生效**
```scss
// 错误写法：直接覆盖 CSS 变量可能不生效
:root {
  --el-color-primary: #ff0000;
}

// 正确写法：按照 Element Plus 规范
// 1. 使用 SCSS 变量
$--color-primary: #ff0000;

// 2. 或在根元素上设置
html {
  --el-color-primary: #ff0000;
}

// 3. 使用 Element Plus 主题配置
import { ElConfigProvider } from 'element-plus'

这个完善的指南应该能帮助你更好地使用 Claude 和 Figma 进行 Vue 3 + TypeScript + Element Plus 应用的开发。如果在开发过程中遇到其他问题，可以随时参考这个文档或向 Claude 寻求帮助。