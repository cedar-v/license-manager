/**
 * 日期格式化工具函数
 */

/**
 * 格式化日期字符串为中文格式
 * @param dateString 日期字符串
 * @returns 格式化后的日期字符串，如果为空则返回"永久"
 */
export const formatDate = (dateString: string | null | undefined): string => {
  if (!dateString) return '永久'
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return dateString
    
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    console.error('日期格式化失败:', error)
    return dateString
  }
}

/**
 * 格式化日期为简短格式（仅日期，不包含时间）
 * @param dateString 日期字符串
 * @returns 格式化后的日期字符串，如果为空则返回"永久"
 */
export const formatDateShort = (dateString: string | null | undefined): string => {
  if (!dateString) return '永久'
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return dateString
    
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit'
    })
  } catch (error) {
    console.error('日期格式化失败:', error)
    return dateString
  }
}

/**
 * 格式化日期为完整格式（包含秒）
 * @param dateString 日期字符串
 * @returns 格式化后的日期字符串，如果为空则返回"永久"
 */
export const formatDateTime = (dateString: string | null | undefined): string => {
  if (!dateString) return '永久'
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return dateString
    
    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit'
    })
  } catch (error) {
    console.error('日期格式化失败:', error)
    return dateString
  }
}

/**
 * 检查日期是否已过期
 * @param dateString 日期字符串
 * @returns 如果日期已过期返回 true，否则返回 false
 */
export const isDateExpired = (dateString: string | null | undefined): boolean => {
  if (!dateString) return false
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return false
    
    return date < new Date()
  } catch (error) {
    console.error('日期检查失败:', error)
    return false
  }
}

/**
 * 获取相对时间描述（如：3天前、1个月前）
 * @param dateString 日期字符串
 * @returns 相对时间描述
 */
export const getRelativeTime = (dateString: string | null | undefined): string => {
  if (!dateString) return '永久'
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return dateString
    
    const now = new Date()
    const diff = now.getTime() - date.getTime()
    const days = Math.floor(diff / (1000 * 60 * 60 * 24))
    
    if (days === 0) return '今天'
    if (days === 1) return '昨天'
    if (days < 7) return `${days}天前`
    if (days < 30) return `${Math.floor(days / 7)}周前`
    if (days < 365) return `${Math.floor(days / 30)}个月前`
    return `${Math.floor(days / 365)}年前`
  } catch (error) {
    console.error('相对时间计算失败:', error)
    return dateString
  }
}
