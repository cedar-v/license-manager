<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-20 16:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-09-03 11:41:35
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
          <span class="value">{{ enumLabels.companySize || '--' }}</span>
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
        <div class="info-item inline-group full-width">
          <div class="inline-pair">
            <span class="label">创建人：</span>
            <span class="value">{{ customerData.statusRecords?.[0]?.user || '--' }}</span>
          </div>
          <div class="inline-pair">
            <span class="label">创建时间：</span>
            <span class="value">{{ formatDate(customerData.statusRecords?.[0]?.time) }}</span>
          </div>
        </div>
        <div class="info-item inline-group full-width">
          <div class="inline-pair">
            <span class="label">更新人：</span>
            <span class="value">{{ customerData.statusRecords?.[1]?.user || '--' }}</span>
          </div>
          <div class="inline-pair">
            <span class="label">更新时间：</span>
            <span class="value">{{ formatDate(customerData.statusRecords?.[1]?.time) }}</span>
          </div>
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
        <div class="info-item inline-group full-width">
          <div class="inline-pair">
            <span class="label">即将过期：</span>
            <span class="value">{{ customerData.expiringSoonLicenses || 0 }}</span>
          </div>
          <div class="inline-pair">
            <span class="label">即将过期时间：</span>
            <span class="value">{{ formatDate(customerData.expiringSoonTime) }}</span>
          </div>
        </div>
        <div class="info-item inline-group full-width">
          <div class="inline-pair">
            <span class="label">过期授权：</span>
            <span class="value">{{ customerData.expiredLicenses || 0 }}</span>
          </div>
          <div class="inline-pair">
            <span class="label">过期时间：</span>
            <span class="value">{{ formatDate(customerData.expiredTime) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getCustomerDetail } from '@/api/customer'

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

// 枚举标签缓存
const enumLabels = ref({
  customerType: '',
  customerLevel: '',
  status: '',
  companySize: ''
})

// 获取客户详情数据
const fetchCustomerDetail = async () => {
  if (!props.customerId) return
  
  loading.value = true
  try {
    // 调用真实API接口
    const response = await getCustomerDetail(props.customerId)
    const customer = response.data
    
    if (!customer) {
      ElMessage.error('客户数据不存在')
      return
    }
    
    // 使用原始枚举值作为显示值 (移除映射功能)
    enumLabels.value = {
      customerType: customer.customer_type_display || customer.customer_type,
      customerLevel: customer.customer_level_display || customer.customer_level,
      status: customer.status_display || customer.status,
      companySize: customer.company_size_display || customer.company_size || ''
    }
    
    // 转换API数据为组件需要的格式
    customerData.value = {
      name: customer.customer_name,
      contact: customer.contact_person,
      email: customer.email || '',
      phone: customer.phone || '',
      address: customer.address || '',
      type: customer.customer_type,
      level: customer.customer_level,
      status: customer.status,
      companySize: customer.company_size || '',
      preferredLicense: '在线授权', // 暂时写死，等后端添加字段
      description: customer.description || '',
      statusRecords: [
        {
          action: '创建人',
          user: customer.created_by,
          time: customer.created_at
        },
        {
          action: '更新人',
          user: customer.updated_by,
          time: customer.updated_at
        }
      ],
      // 授权统计数据等待后端接口
      totalLicenses: 0,
      activeLicenses: 0,
      expiredLicenses: 0,
      expiringSoonLicenses: 0,
      latestLicenseTime: undefined,
      expiringSoonTime: undefined,
      expiredTime: undefined
    }
  } catch (error: any) {
    console.error('获取客户详情失败:', error)
    // 使用后端返回的错误信息
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(errorMessage)
    }
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
  return enumLabels.value.customerType || '--'
}

// 获取客户等级标签
const getCustomerLevelLabel = (level: string | undefined) => {
  if (!level) return '--'
  return enumLabels.value.customerLevel || '--'
}

// 获取状态标签
const getStatusLabel = (status: string | undefined) => {
  if (!status) return '--'
  return enumLabels.value.status || '--'
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
  padding: 15px;
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

.info-item.inline-group {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
  gap: 40px;
}

.inline-pair {
  display: flex;
  align-items: center;
  min-width: 240px;
}

.inline-pair .label {
  min-width: 100px;
  margin-right: 8px;
  flex-shrink: 0;
}

.inline-pair .value {
  min-width: 120px;
  flex-shrink: 0;
}

.label {
  font-size: 14px;
  color: #666;
  font-weight: 500;
  min-width: 100px;
  flex-shrink: 0;
  margin-right: 8px;
}

.value {
  font-size: 14px;
  color: #1d1d1d;
  word-break: break-all;
  line-height: 1.5;
  min-width: 120px;
  flex-shrink: 0;
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