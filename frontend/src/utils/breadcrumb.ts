/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-12 11:28:54
 * @FilePath: /frontend/src/utils/breadcrumb.ts
 * @Description: 面包屑导航工具函数
 */

import { computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// 面包屑项接口
export interface BreadcrumbItem {
  title: string
  path?: string
  icon?: string
}

// 路由面包屑配置
const routeBreadcrumbConfig: Record<string, BreadcrumbItem[]> = {
  '/dashboard': [
    { title: '仪表盘', path: '/dashboard' }
  ],
  '/customers': [
    { title: '客户管理', path: '/customers' }
  ],
  '/licenses': [
    { title: '授权管理', path: '/licenses' }
  ],
  '/roles': [
    { title: '角色权限', path: '/roles' }
  ],
  '/users': [
    { title: '系统用户', path: '/users' }
  ],
  '/login': [
    { title: '登录' }
  ]
}

export function useBreadcrumb() {
  const route = useRoute()
  const router = useRouter()

  // 生成面包屑
  const breadcrumbs = computed<BreadcrumbItem[]>(() => {
    const currentPath = route.path
    
    // 从配置中获取面包屑
    const configBreadcrumbs = routeBreadcrumbConfig[currentPath]
    if (configBreadcrumbs) {
      return configBreadcrumbs.map((item, index) => ({
        ...item,
        // 最后一项不需要链接
        path: index === configBreadcrumbs.length - 1 ? undefined : item.path
      }))
    }

    // 如果没有配置，自动生成
    return generateAutoBreadcrumb(currentPath)
  })

  // 当前页面标题
  const pageTitle = computed(() => {
    return (route.meta?.title as string) || '未知页面'
  })

  // 自动生成面包屑
  function generateAutoBreadcrumb(path: string): BreadcrumbItem[] {
    const segments = path.split('/').filter(Boolean)
    const breadcrumbs: BreadcrumbItem[] = []

    // 总是以工作台开始（除了登录页）
    if (path !== '/login') {
      breadcrumbs.push({ title: '工作台', path: '/dashboard' })
    }

    // 根据路径段生成面包屑
    let currentPath = ''
    segments.forEach((segment, index) => {
      currentPath += `/${segment}`
      
      if (currentPath === '/dashboard') return // 已经添加过了
      
      const title = getSegmentTitle(segment)
      breadcrumbs.push({
        title,
        path: index === segments.length - 1 ? undefined : currentPath
      })
    })

    return breadcrumbs
  }

  // 根据路径段获取标题
  function getSegmentTitle(segment: string): string {
    const titleMap: Record<string, string> = {
      'dashboard': '工作台',
      'customers': '客户管理',
      'licenses': '授权管理',
      'roles': '角色权限',
      'users': '系统用户',
      'login': '登录',
      'settings': '设置',
      'profile': '个人资料'
    }
    
    return titleMap[segment] || segment
  }

  // 导航到面包屑项
  function navigateTo(item: BreadcrumbItem) {
    if (item.path && item.path !== route.path) {
      router.push(item.path)
    }
  }

  // 监听路由变化，可以在这里做一些额外处理
  watch(() => route.path, (newPath) => {
    // 这里可以添加路由变化时的逻辑
    console.log('Route changed to:', newPath)
  })

  return {
    breadcrumbs,
    pageTitle,
    navigateTo
  }
}