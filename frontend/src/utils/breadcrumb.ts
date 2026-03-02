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
import { useI18n } from 'vue-i18n'

// 面包屑项接口
export interface BreadcrumbItem {
  title: string
  path?: string
  icon?: string
}

// 路由面包屑配置的key
const routeBreadcrumbKeys: Record<string, string[]> = {
  '/dashboard': ['dashboard'],
  '/customers': ['customers'],
  '/enterprise-leads': ['enterpriseLeads'],
  '/licenses': ['licenses'],
  '/invoices': ['invoices'],
  '/invoices/detail': ['invoices', 'invoiceDetail'],
  '/packages': ['packages'],
  '/roles': ['roles'],
  '/users': ['users'],
  '/login': ['login']
}

export function useBreadcrumb() {
  const route = useRoute()
  const router = useRouter()
  const { t } = useI18n()

  // 生成面包屑
  const breadcrumbs = computed<BreadcrumbItem[]>(() => {
    const currentPath = route.path

    // 特殊处理授权管理子路由
    if (currentPath.startsWith('/licenses/')) {
      const items: BreadcrumbItem[] = []

      // 添加授权管理作为父级
      items.push({
        title: t('navigation.breadcrumb.licenses'),
        path: '/licenses'
      })

      // 如果有客户名称查询参数，显示客户名称
      if (route.query.customerName && typeof route.query.customerName === 'string') {
        items.push({
          title: route.query.customerName
        })
      } else {
        // 否则显示当前路由的标题
        const pathSegments = currentPath.split('/').filter(Boolean)
        const lastSegment = pathSegments[pathSegments.length - 1]

        if (lastSegment === 'list') {
          items.push({
            title: t('navigation.breadcrumb.licenseList')
          })
        } else if (lastSegment === 'create') {
          items.push({
            title: t('navigation.breadcrumb.createLicense')
          })
        }
      }

      return items
    }

    // 特殊处理发票详情子路由（带ID）
    if (currentPath.startsWith('/invoices/detail/')) {
      const items: BreadcrumbItem[] = []

      // 添加发票管理作为父级
      items.push({
        title: t('navigation.breadcrumb.invoices'),
        path: '/invoices'
      })

      // 添加发票详情作为当前项
      items.push({
        title: t('navigation.breadcrumb.invoiceDetail')
      })

      return items
    }

    // 从配置中获取面包屑
    const configKeys = routeBreadcrumbKeys[currentPath]
    if (configKeys) {
      return configKeys.map((key, index) => ({
        title: t(`navigation.breadcrumb.${key}`),
        // 最后一项不需要链接
        path: index === configKeys.length - 1 ? undefined : currentPath
      }))
    }

    // 如果没有配置，自动生成
    return generateAutoBreadcrumb(currentPath)
  })

  // 当前页面标题
  const pageTitle = computed(() => {
    return (route.meta?.title as string) || t('navigation.breadcrumb.unknown')
  })

  // 自动生成面包屑
  function generateAutoBreadcrumb(path: string): BreadcrumbItem[] {
    const segments = path.split('/').filter(Boolean)
    const breadcrumbs: BreadcrumbItem[] = []

    // 总是以工作台开始（除了登录页）
    if (path !== '/login') {
      breadcrumbs.push({ title: t('navigation.breadcrumb.workspace'), path: '/dashboard' })
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
    const knownSegments = ['dashboard', 'customers', 'enterprise-leads', 'licenses', 'roles', 'users', 'login', 'settings', 'profile']
    
    if (knownSegments.includes(segment)) {
      return t(`navigation.breadcrumb.${segment}`)
    }
    
    return segment
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