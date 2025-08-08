/**
 * Layout布局组件导出文件
 * 统一管理布局相关组件和类型定义的导出
 */

// 布局组件导出
export { default as Layout } from './Layout.vue'
export { default as Sidebar } from './Sidebar.vue' 
export { default as NavContent } from './NavContent.vue'

// 类型定义导出

// 导航菜单项类型定义
export interface NavItem {
  id: string // 导航项唯一标识
  label: string // 显示文本
  href: string // 链接地址
  icon?: string // 图标类名（可选）
  active?: boolean // 是否为当前活跃项（可选）
  children?: NavItem[] // 子菜单项（可选）
}

// 面包屑导航项类型定义
export interface BreadcrumbItem {
  title: string // 显示标题
  path?: string // 链接路径（可选，最后一项通常不需要链接）
}