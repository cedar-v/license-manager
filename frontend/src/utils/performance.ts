/**
 * 性能优化工具函数
 * 包含防抖、节流等常用性能优化函数
 */

/**
 * 防抖函数
 * 在事件停止触发n毫秒后才执行函数，适用于input输入、window resize等场景
 * @param func 要执行的函数
 * @param wait 延迟时间（毫秒）
 * @param immediate 是否立即执行（首次触发时）
 */
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number,
  immediate = false
): (...args: Parameters<T>) => void {
  let timeout: number | undefined
  return function executedFunction(...args: Parameters<T>) {
    const later = () => {
      timeout = undefined
      if (!immediate) func(...args)
    }
    
    const callNow = immediate && !timeout
    
    if (timeout !== undefined) {
      clearTimeout(timeout)
    }
    
    timeout = window.setTimeout(later, wait)
    
    if (callNow) func(...args)
  }
}

/**
 * 节流函数  
 * 限制函数在指定时间内最多执行一次，适用于scroll、mousemove等高频事件
 * @param func 要执行的函数
 * @param limit 时间间隔（毫秒）
 */
export function throttle<T extends (...args: any[]) => any>(
  func: T,
  limit: number
): (...args: Parameters<T>) => void {
  let inThrottle: boolean
  return function executedFunction(...args: Parameters<T>) {
    if (!inThrottle) {
      func(...args)
      inThrottle = true
      setTimeout(() => (inThrottle = false), limit)
    }
  }
}

/**
 * requestAnimationFrame 节流
 * 使用浏览器的 requestAnimationFrame 来节流，适合动画相关的操作
 * @param func 要执行的函数
 */
export function rafThrottle<T extends (...args: any[]) => any>(
  func: T
): (...args: Parameters<T>) => void {
  let rafId: number | undefined
  return function executedFunction(...args: Parameters<T>) {
    if (rafId === undefined) {
      rafId = requestAnimationFrame(() => {
        func(...args)
        rafId = undefined
      })
    }
  }
}

/**
 * 延迟执行函数
 * @param func 要执行的函数
 * @param delay 延迟时间（毫秒）
 */
export function delay(func: () => void, delay: number): number {
  return window.setTimeout(func, delay)
}

/**
 * 取消延迟执行
 * @param timerId setTimeout 返回的 ID
 */
export function cancelDelay(timerId: number): void {
  clearTimeout(timerId)
}