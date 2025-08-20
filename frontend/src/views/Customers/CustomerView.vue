<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-20 16:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-20 16:00:00
 * @FilePath: /frontend/src/views/Customers/CustomerView.vue
 * @Description: 客户查看组件
-->
<template>
  <!-- 顶部横向区域 - 面包屑和返回按钮 -->
  <div class="top-section">
    <!-- 面包屑导航 -->
    <div class="breadcrumb">
      <span class="breadcrumb-item">客户管理</span>
      <span class="breadcrumb-separator">></span>
      <span class="breadcrumb-current">查看客户</span>
    </div>
    
    <!-- 返回按钮 -->
    <div class="form-actions">
      <el-button @click="handleBack">返 回</el-button>
    </div>
  </div>

  <!-- 主内容区域 -->
  <div class="customer-view" v-loading="loading">
    <!-- 基本信息卡片 -->
    <div class="info-card basic-info">
      <h3 class="section-title">基本信息</h3>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">客户名称：</span>
          <span class="value">{{ customerData.name || '--' }}</span>
        </div>
        <div class="info-item">
          <span class="label">客户类型：</span>
          <span class="value">{{ getCustomerTypeLabel(customerData.type) }}</span>
        </div>
        <div class="info-item">
          <span class="label">客户等级：</span>
          <span class="value">{{ getCustomerLevelLabel(customerData.level) }}</span>
        </div>
        <div class="info-item">
          <span class="label">联系人：</span>
          <span class="value">{{ customerData.contact || '--' }}</span>
        </div>
        <div class="info-item">
          <span class="label">邮箱：</span>
          <span class="value">{{ customerData.email || '--' }}</span>
        </div>
        <div class="info-item">
          <span class="label">状态：</span>
          <span class="value">{{ getStatusLabel(customerData.status) }}</span>
        </div>
        <div class="info-item">
          <span class="label">联系电话：</span>
          <span class="value">{{ customerData.phone || '--' }}</span>
        </div>
        <div class="info-item full-width">
          <span class="label">地址：</span>
          <span class="value">{{ customerData.address || '--' }}</span>
        </div>
      </div>
    </div>

    <!-- 商业信息卡片 -->
    <div class="info-card business-info">
      <h3 class="section-title">商业信息</h3>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">企业规模：</span>
          <span class="value">{{ getCompanySizeLabel(customerData.companySize) }}</span>
        </div>
        <div class="info-item full-width">
          <span class="label">首选授权：</span>
          <span class="value">{{ customerData.preferredLicense || '--' }}</span>
        </div>
        <div class="info-item full-width">
          <span class="label">客户描述：</span>
          <span class="value description">{{ customerData.description || '--' }}</span>
        </div>
      </div>
    </div>

    <!-- 状态信息卡片 -->
    <div class="info-card status-info">
      <h3 class="section-title">状态信息</h3>
      <div class="info-grid">
        <div class="info-item full-width" v-for="(record, index) in customerData.statusRecords" :key="index">
          <span class="label">{{ record.action }}：</span>
          <span class="value">{{ record.user }}，{{ record.action === '创建人' ? '创建时间' : '更新时间' }}：{{ formatDate(record.time) }}</span>
        </div>
      </div>
    </div>

    <!-- 授权统计卡片 -->
    <div class="info-card license-stats">
      <h3 class="section-title">授权统计</h3>
      <div class="info-grid">
        <div class="info-item">
          <span class="label">总授权数：</span>
          <span class="value">{{ customerData.totalLicenses || 0 }}</span>
        </div>
        <div class="info-item">
          <span class="label">活跃授权：</span>
          <span class="value">{{ customerData.activeLicenses || 0 }}</span>
        </div>
        <div class="info-item">
          <span class="label">最近授权：</span>
          <span class="value">{{ formatDate(customerData.latestLicenseTime) }}</span>
        </div>
        <div class="info-item">
          <span class="label">即将过期：</span>
          <span class="value">{{ customerData.expiringSoonLicenses || 0 }}</span>
        </div>
        <div class="info-item">
          <span class="label">即将过期时间：</span>
          <span class="value">{{ formatDate(customerData.expiringSoonTime) }}</span>
        </div>
        <div class="info-item">
          <span class="label">过期授权：</span>
          <span class="value">{{ customerData.expiredLicenses || 0 }}</span>
        </div>
        <div class="info-item">
          <span class="label">过期时间：</span>
          <span class="value">{{ formatDate(customerData.expiredTime) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'

interface StatusRecord {
  action: string
  user: string
  time: string
}

interface CustomerViewData {
  name: string
  contact: string
  email: string
  phone: string
  address: string
  type: string
  level: string
  status: string
  companySize: string
  preferredLicense: string
  description: string
  statusRecords?: StatusRecord[]
  totalLicenses?: number
  activeLicenses?: number
  expiredLicenses?: number
  expiringSoonLicenses?: number
  latestLicenseTime?: string
  expiringSoonTime?: string
  expiredTime?: string
}

const props = defineProps<{
  customerId: string
}>()

const emit = defineEmits<{
  back: []
}>()

const loading = ref(false)
const customerData = ref<Partial<CustomerViewData>>({})

// 获取客户详情数据
const fetchCustomerDetail = async () => {
  if (!props.customerId) return
  
  loading.value = true
  try {
    // TODO: 调用真实API接口
    // const response = await getCustomerDetail(props.customerId)
    // customerData.value = response.data
    
    // 模拟数据
    await new Promise(resolve => setTimeout(resolve, 500))
    customerData.value = {
      name: '随州市海留有限公司',
      contact: '祁瑾',
      email: '13988887963@qq.com',
      phone: '13988887963',
      address: '湖北省随州市经济开发区创业大厦A座1001室',
      type: 'enterprise',
      level: 'vip',
      status: 'normal',
      companySize: 'medium',
      preferredLicense: '在线授权',
      description: '专业从事软件开发和技术服务的高新技术企业，主要业务包括企业管理系统开发、移动应用开发等。',
      statusRecords: [
        {
          action: '创建人',
          user: '张三',
          time: '2025-12-01T12:25:51Z'
        },
        {
          action: '更新人',
          user: '李四',
          time: '2025-12-15T16:30:25Z'
        }
      ],
      totalLicenses: 15,
      activeLicenses: 1200,
      expiredLicenses: 1,
      expiringSoonLicenses: 1,
      latestLicenseTime: '2025-12-01T12:25:51Z',
      expiringSoonTime: '2025-12-01T12:25:51Z',
      expiredTime: '2025-12-01T12:25:51Z'
    }
  } catch (error) {
    console.error('获取客户详情失败:', error)
    ElMessage.error('获取客户详情失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 返回操作
const handleBack = () => {
  emit('back')
}

// 获取客户类型标签
const getCustomerTypeLabel = (type: string | undefined) => {
  if (!type) return '--'
  const typeMap: { [key: string]: string } = {
    'enterprise': '企业客户',
    'individual': '个人客户',
    'government': '政府机构'
  }
  return typeMap[type] || '--'
}

// 获取客户等级标签
const getCustomerLevelLabel = (level: string | undefined) => {
  if (!level) return '--'
  const levelMap: { [key: string]: string } = {
    'vip': 'VIP客户',
    'important': '重要客户',
    'normal': '普通客户'
  }
  return levelMap[level] || '--'
}

// 获取状态标签
const getStatusLabel = (status: string | undefined) => {
  if (!status) return '--'
  const statusMap: { [key: string]: string } = {
    'normal': '正常',
    'disabled': '禁用'
  }
  return statusMap[status] || '--'
}

// 获取企业规模标签
const getCompanySizeLabel = (size: string | undefined) => {
  if (!size) return '--'
  const sizeMap: { [key: string]: string } = {
    'large': '大型企业(500人以上)',
    'medium': '中型企业(100-500人)',
    'small': '小型企业(50-100人)',
    'micro': '微型企业(50人以下)'
  }
  return sizeMap[size] || '--'
}

// 格式化日期
const formatDate = (date: string | undefined) => {
  if (!date) return '--'
  return new Date(date).toLocaleString('zh-CN')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchCustomerDetail()
})
</script>

<style scoped>
/* 顶部横向区域 */
.top-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: transparent;
  margin-bottom: 10px;
}

.breadcrumb {
  font-size: 14px;
  color: #888;
}

.breadcrumb-item {
  color: #888;
}

.breadcrumb-separator {
  margin: 0 8px;
  color: #b6bad2;
}

.breadcrumb-current {
  color: #1d1d1d;
}

.form-actions {
  display: flex;
  gap: 20px;
}

/* 主内容区域 */
.customer-view {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* 信息卡片 */
.info-card {
  background: #fff;
  border-radius: 4px;
  padding: 20px;
}

.basic-info {
  min-height: 172px;
}

.business-info {
  min-height: 172px;
}

.status-info {
  min-height: 134px;
}

.license-stats {
  min-height: 172px;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  color: #1d1d1d;
  margin: 0 0 20px 0;
}

/* 信息网格布局 */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 16px 40px;
  align-items: start;
}

.info-item {
  display: flex;
  align-items: flex-start;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
  min-width: 80px;
  flex-shrink: 0;
  margin-right: 8px;
}

.value {
  font-size: 14px;
  color: #1d1d1d;
  word-break: break-all;
  line-height: 1.5;
}

.value.description {
  line-height: 1.6;
  white-space: pre-wrap;
}


/* Element Plus 按钮样式 */
:deep(.el-button) {
  font-size: 14px;
  padding: 6px 16px;
  border: 1px solid #c2c6ce;
  color: #1d1d1d;
}

:deep(.el-button:hover) {
  border-color: #019c7c;
  color: #019c7c;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .top-section {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .info-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
  
  .info-item {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .label {
    margin-bottom: 4px;
    margin-right: 0;
  }
}

@media (max-width: 480px) {
  .info-card {
    padding: 16px;
  }
  
  .customer-view {
    gap: 16px;
  }
}
</style>