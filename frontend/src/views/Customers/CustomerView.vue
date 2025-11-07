<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-20 16:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-11-03 15:44:00
 * @FilePath: /frontend/src/views/Customers/CustomerView.vue
 * @Description: 客户查看组件
-->
<template>
  <div class="customer-view-container">
    <!-- 顶部横向区域 - 面包屑和返回按钮 -->
    <div class="top-section">
      <!-- 面包屑导航 -->
      <div class="breadcrumb">
        <span class="breadcrumb-item">{{ t('customers.breadcrumb.customerManagement') }}</span>
        <span class="breadcrumb-separator">></span>
        <span class="breadcrumb-current">{{ t('customers.viewCustomer') }}</span>
      </div>

      <!-- 返回按钮 -->
      <div class="form-actions">
        <el-button @click="handleBack">{{ t('customers.view.back') }}</el-button>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="customer-view" v-loading="loading">
      <!-- 基本信息卡片 -->
      <div class="info-card basic-info">
        <h3 class="section-title">{{ t('customers.form.basicInfo') }}</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">{{ t('customers.form.customerName') }}：</span>
            <span class="value">{{ customerData.name || t('customers.view.noData') }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('customers.form.customerType') }}：</span>
            <span class="value">{{ getCustomerTypeLabel(customerData.type) }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('customers.form.customerLevel') }}：</span>
            <span class="value">{{ getCustomerLevelLabel(customerData.level) }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('customers.form.contactPerson') }}：</span>
            <span class="value">{{ customerData.contact || t('customers.view.noData') }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('customers.form.email') }}：</span>
            <span class="value">{{ customerData.email || t('customers.view.noData') }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('customers.form.status') }}：</span>
            <span class="value">{{ getStatusLabel(customerData.status) }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('customers.form.phone') }}：</span>
            <span class="value">{{ customerData.phone || t('customers.view.noData') }}</span>
          </div>
          <div class="info-item full-width">
            <span class="label">{{ t('customers.form.address') }}：</span>
            <span class="value">{{ customerData.address || t('customers.view.noData') }}</span>
          </div>
        </div>
      </div>

      <!-- 商业信息卡片 -->
      <div class="info-card business-info">
        <h3 class="section-title">{{ t('customers.form.businessInfo') }}</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">{{ t('customers.form.companySize') }}：</span>
            <span class="value">{{ enumLabels.companySize || t('customers.view.noData') }}</span>
          </div>
          <div class="info-item full-width">
            <span class="label">{{ t('customers.form.preferredLicense') }}：</span>
            <span class="value">{{
              customerData.preferredLicense || t('customers.view.noData')
            }}</span>
          </div>
          <div class="info-item full-width">
            <span class="label">{{ t('customers.form.description') }}：</span>
            <span class="value description">{{
              customerData.description || t('customers.view.noData')
            }}</span>
          </div>
        </div>
      </div>

      <!-- 状态信息卡片 -->
      <div class="info-card status-info">
        <h3 class="section-title">{{ t('customers.view.statusInfo') }}</h3>
        <div class="info-grid">
          <div class="info-item inline-group full-width">
            <div class="inline-pair">
              <span class="label">{{ t('customers.view.creator') }}：</span>
              <span class="value">{{
                customerData.statusRecords?.[0]?.user || t('customers.view.noData')
              }}</span>
            </div>
            <div class="inline-pair">
              <span class="label">{{ t('customers.view.createTime') }}：</span>
              <span class="value">{{ formatDate(customerData.statusRecords?.[0]?.time) }}</span>
            </div>
          </div>
          <div class="info-item inline-group full-width">
            <div class="inline-pair">
              <span class="label">{{ t('customers.view.updater') }}：</span>
              <span class="value">{{
                customerData.statusRecords?.[1]?.user || t('customers.view.noData')
              }}</span>
            </div>
            <div class="inline-pair">
              <span class="label">{{ t('customers.view.updateTime') }}：</span>
              <span class="value">{{ formatDate(customerData.statusRecords?.[1]?.time) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 授权统计卡片 -->
      <div class="info-card license-stats">
        <h3 class="section-title">{{ t('customers.view.licenseStats') }}</h3>
        <div class="info-grid">
          <div class="info-item">
            <span class="label">{{ t('customers.view.totalLicenses') }}：</span>
            <span class="value">{{ customerData.total_auth_codes || 0 }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('customers.view.activeLicenses') }}：</span>
            <span class="value">{{ customerData.expired_auth_codes || 0 }}</span>
          </div>
          <div class="info-item">
            <span class="label">{{ t('customers.view.recentLicense') }}：</span>
            <span class="value">{{ customerData.expiring_soon_auth_codes || 0 }}</span>
          </div>
          <div class="info-item inline-group full-width">
            <div class="inline-pair">
              <span class="label">{{ t('customers.view.expiringSoon') }}：</span>
              <span class="value">{{ customerData.total_licenses || 0 }}</span>
            </div>
            <div class="inline-pair">
              <span class="label">{{ t('customers.view.expiringSoonTime') }}：</span>
              <span class="value">{{ customerData.active_licenses || 0 }}</span>
            </div>
          </div>
          <div class="info-item inline-group full-width">
            <div class="inline-pair">
              <span class="label">{{ t('customers.view.expiredLicenses') }}：</span>
              <span class="value">{{ customerData.expired_licenses || 0 }}</span>
            </div>
            <div class="inline-pair">
              <span class="label">{{ t('customers.view.expiredTime') }}：</span>
              <span class="value">{{ customerData.inactive_licenses }}</span>
            </div>
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
import { useI18n } from 'vue-i18n'
import { formatDate } from '@/utils/date'

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
  active_licenses: number
  expired_auth_codes: number
  expired_licenses: number
  expiring_soon_auth_codes: number
  inactive_licenses: number
  total_auth_codes: number
  total_licenses: number
}

const props = defineProps<{
  customerId: string
}>()

const emit = defineEmits<{
  back: []
}>()

const { t } = useI18n()
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
      ElMessage.error(t('customers.view.customerNotExist'))
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
      preferredLicense: t('customers.view.onlineLicense'), // 暂时写死，等后端添加字段
      description: customer.description || '',
      statusRecords: [
        {
          action: t('customers.view.creator'),
          user: customer.created_by,
          time: customer.created_at
        },
        {
          action: t('customers.view.updater'),
          user: customer.updated_by,
          time: customer.updated_at
        }
      ],
      // 授权统计数据
      total_auth_codes: customer.authorization_stats?.total_auth_codes || 0, // 总授权码数量
      expired_auth_codes: customer.authorization_stats?.expired_auth_codes || 0, // 已过期授权码数量
      expiring_soon_auth_codes: customer.authorization_stats?.expiring_soon_auth_codes || 0, // 30日内即将到期授权码数量
      total_licenses: customer.authorization_stats?.total_licenses || 0, // 总许可证数量
      active_licenses: customer.authorization_stats?.active_licenses || 0, // 已激活许可证数量
      expired_licenses: customer.authorization_stats?.expired_licenses || 0, // 已过期许可证数量
      inactive_licenses: customer.authorization_stats?.inactive_licenses || 0 // 未激活许可证数量
    }
  } catch (error: any) {
    console.error('Get customer detail failed:', error)
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
  if (!type) return t('customers.view.noData')
  return enumLabels.value.customerType || t('customers.view.noData')
}

// 获取客户等级标签
const getCustomerLevelLabel = (level: string | undefined) => {
  if (!level) return t('customers.view.noData')
  return enumLabels.value.customerLevel || t('customers.view.noData')
}

// 获取状态标签
const getStatusLabel = (status: string | undefined) => {
  if (!status) return t('customers.view.noData')
  return enumLabels.value.status || t('customers.view.noData')
}

// 组件挂载时获取数据
onMounted(() => {
  fetchCustomerDetail()
})
</script>

<style scoped>
/* 容器 */
.customer-view-container {
  display: flex;
  flex-direction: column;
}

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
  color: var(--app-text-secondary);
}

.breadcrumb-item {
  color: var(--app-text-secondary);
}

.breadcrumb-separator {
  margin: 0 8px;
  color: var(--app-text-secondary);
}

.breadcrumb-current {
  color: var(--app-text-primary);
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
  background: var(--app-content-bg);
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
  color: var(--app-text-primary);
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
  color: var(--app-text-regular);
  font-weight: 500;
  min-width: 100px;
  flex-shrink: 0;
  margin-right: 8px;
}

.value {
  font-size: 14px;
  color: var(--app-text-primary);
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
  border: 1px solid var(--app-border-color);
  color: var(--app-text-primary);
}

:deep(.el-button:hover) {
  border-color: var(--el-color-primary);
  color: var(--el-color-primary);
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
