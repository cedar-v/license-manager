<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-09-29 17:29:42
 * @FilePath: /frontend/src/views/Licenses.vue
 * @Description: 授权管理页面
-->
<template>
  <Layout app-name="Cedar-V" :page-title="t('pages.licenses.title')">
    <div class="license-container">
      <!-- 主要内容区域 -->
      <div class="main-content">
        <!-- 背景图片 -->
        <div class="background-section">
          <!-- 中央标题 -->
          <div class="center-title">
            <h1 class="platform-title">{{ t('pages.licenses.platform') }}</h1>
          </div>

          <!-- 操作区域 -->
          <div class="action-section">
            <!-- 客户选择下拉框和操作按钮 -->
            <div class="action-row">
              <!-- 客户选择下拉框 -->
              <div class="customer-select-wrapper">
                <el-select
                  v-model="selectedCustomer"
                  :placeholder="t('pages.licenses.search.selectCustomer')"
                  clearable
                  filterable
                  size="large"
                  class="customer-select"
                  @change="handleCustomerChange"
                >
                  <el-option
                    :label="t('pages.licenses.search.allCustomers')"
                    value=""
                  />
                  <el-option
                    v-for="customer in customers"
                    :key="customer.id"
                    :label="customer.customer_name"
                    :value="customer.id"
                  />
                </el-select>
              </div>

              <!-- 操作按钮 -->
              <div class="action-buttons">
                <el-button
                  type="primary"
                  size="large"
                  class="action-btn query-btn"
                  @click="handleQuery"
                >
                  {{ t('pages.licenses.actions.query') }}
                </el-button>

                <el-button
                  type="primary"
                  size="large"
                  class="action-btn create-btn"
                  @click="handleCreateLicense"
                >
                  {{ t('pages.licenses.actions.createLicense') }}
                </el-button>
              </div>
            </div>
          </div>
        </div>

        <!-- 授权列表区域 -->
        <div v-if="showLicenseList" class="license-list-section">
          <div class="list-header">
            <div class="search-controls">
              <el-input
                v-model="searchText"
                :placeholder="t('pages.licenses.search.placeholder')"
                prefix-icon="Search"
                clearable
                size="large"
                class="search-input"
                @input="handleSearch"
              />

              <el-date-picker
                v-model="dateRange"
                type="daterange"
                :range-separator="t('pages.licenses.search.dateRange')"
                :start-placeholder="t('pages.licenses.search.startDate')"
                :end-placeholder="t('pages.licenses.search.endDate')"
                size="large"
                class="date-picker"
                @change="handleDateRangeChange"
              />

              <el-button
                type="primary"
                size="large"
                @click="refreshLicenseList"
              >
                {{ t('pages.licenses.actions.refresh') }}
              </el-button>
            </div>
          </div>

          <!-- 授权表格 -->
          <div class="table-container">
            <el-table
              :data="licenseList"
              v-loading="loading"
              :element-loading-text="t('pages.licenses.table.loading')"
              stripe
              style="width: 100%"
              class="license-table"
            >
              <el-table-column
                prop="license_code"
                :label="t('pages.licenses.table.licenseCode')"
                width="200"
                show-overflow-tooltip
              />

              <el-table-column
                prop="customer_name"
                :label="t('pages.licenses.table.customerName')"
                width="180"
                show-overflow-tooltip
              />

              <el-table-column
                prop="description"
                :label="t('pages.licenses.table.description')"
                min-width="200"
                show-overflow-tooltip
              />

              <el-table-column
                prop="status"
                :label="t('pages.licenses.table.status')"
                width="120"
              >
                <template #default="{ row }">
                  <el-tag
                    :type="getStatusType(row.status)"
                    effect="plain"
                  >
                    {{ getStatusText(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>

              <el-table-column
                prop="license_type_display"
                :label="t('pages.licenses.table.licenseType')"
                width="140"
                show-overflow-tooltip
              />

              <el-table-column
                prop="expiry_date"
                :label="t('pages.licenses.table.expiryDate')"
                width="180"
              >
                <template #default="{ row }">
                  {{ formatDate(row.expiry_date) }}
                </template>
              </el-table-column>

              <el-table-column
                prop="created_at"
                :label="t('pages.licenses.table.createTime')"
                width="180"
              >
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>

              <el-table-column
                :label="t('pages.licenses.table.operation')"
                width="200"
                fixed="right"
              >
                <template #default="{ row }">
                  <el-button
                    type="primary"
                    size="small"
                    @click="handleEdit(row)"
                  >
                    {{ t('pages.licenses.actions.edit') }}
                  </el-button>

                  <el-button
                    v-if="row.status === 'inactive'"
                    type="success"
                    size="small"
                    @click="handleActivate(row)"
                  >
                    {{ t('pages.licenses.actions.activate') }}
                  </el-button>

                  <el-button
                    v-if="row.status === 'active'"
                    type="warning"
                    size="small"
                    @click="handleDeactivate(row)"
                  >
                    {{ t('pages.licenses.actions.deactivate') }}
                  </el-button>

                  <el-button
                    type="danger"
                    size="small"
                    @click="handleDelete(row)"
                  >
                    {{ t('pages.licenses.actions.delete') }}
                  </el-button>
                </template>
              </el-table-column>
            </el-table>

            <!-- 分页 -->
            <div class="pagination-container">
              <el-pagination
                v-model:current-page="currentPage"
                v-model:page-size="pageSize"
                :page-sizes="[10, 20, 50, 100]"
                :total="total"
                layout="total, sizes, prev, pager, next, jumper"
                @size-change="handleSizeChange"
                @current-change="handleCurrentChange"
              />
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建/编辑授权对话框 -->
    <LicenseForm
      v-model:visible="formDialogVisible"
      :license-data="editingLicense"
      :customers="customers"
      @success="handleFormSuccess"
    />
  </Layout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import Layout from '@/components/common/layout/Layout.vue'
import LicenseForm from './LicenseForm.vue'
import { getLicenses, deleteLicense, activateLicense, deactivateLicense, type License, type LicenseQueryRequest } from '@/api/license'
import { getCustomers, type Customer } from '@/api/customer'

const { t } = useI18n()

// 响应式数据
const selectedCustomer = ref<string>('')
const customers = ref<Customer[]>([])
const showLicenseList = ref(false)
const searchText = ref('')
const dateRange = ref<[Date, Date] | null>(null)
const loading = ref(false)
const licenseList = ref<License[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const formDialogVisible = ref(false)
const editingLicense = ref<License | null>(null)

// 查询参数
const queryParams = reactive<LicenseQueryRequest>({
  page: 1,
  page_size: 20,
  search: '',
  customer_id: '',
  sort: 'created_at',
  order: 'desc'
})

// 计算属性
// const formattedCustomers = computed(() => {
//   return customers.value.filter(customer => customer.status === 'active')
// })

// 方法
const loadCustomers = async () => {
  try {
    const response = await getCustomers({ status: 'active', page_size: 1000 })
    customers.value = response.data.list
  } catch (error) {
    console.error('Failed to load customers:', error)
    ElMessage.error('加载客户列表失败')
  }
}

const loadLicenseList = async () => {
  try {
    loading.value = true
    const response = await getLicenses(queryParams)
    licenseList.value = response.data.list
    total.value = response.data.total
    currentPage.value = response.data.page
  } catch (error) {
    console.error('Failed to load licenses:', error)
    ElMessage.error(t('pages.licenses.message.loadError'))
  } finally {
    loading.value = false
  }
}

const handleCustomerChange = (value: string) => {
  queryParams.customer_id = value
}

const handleQuery = () => {
  showLicenseList.value = true
  queryParams.page = 1
  currentPage.value = 1
  loadLicenseList()
}

const handleCreateLicense = () => {
  editingLicense.value = null
  formDialogVisible.value = true
}

const handleSearch = () => {
  queryParams.search = searchText.value
  queryParams.page = 1
  currentPage.value = 1
  if (showLicenseList.value) {
    loadLicenseList()
  }
}

const handleDateRangeChange = (value: [Date, Date] | null) => {
  if (value) {
    queryParams.start_date = value[0].toISOString().split('T')[0]
    queryParams.end_date = value[1].toISOString().split('T')[0]
  } else {
    queryParams.start_date = undefined
    queryParams.end_date = undefined
  }

  if (showLicenseList.value) {
    queryParams.page = 1
    currentPage.value = 1
    loadLicenseList()
  }
}

const refreshLicenseList = () => {
  if (showLicenseList.value) {
    loadLicenseList()
  }
}

const handleSizeChange = (size: number) => {
  queryParams.page_size = size
  queryParams.page = 1
  currentPage.value = 1
  loadLicenseList()
}

const handleCurrentChange = (page: number) => {
  queryParams.page = page
  loadLicenseList()
}

const handleEdit = (license: License) => {
  editingLicense.value = license
  formDialogVisible.value = true
}

const handleActivate = async (license: License) => {
  try {
    await ElMessageBox.confirm(
      t('pages.licenses.confirm.activateMessage', { code: license.license_code }),
      t('pages.licenses.confirm.activateTitle'),
      {
        confirmButtonText: t('pages.licenses.confirm.confirm'),
        cancelButtonText: t('pages.licenses.confirm.cancel'),
        type: 'warning'
      }
    )

    await activateLicense(license.id)
    ElMessage.success(t('pages.licenses.message.activateSuccess'))
    loadLicenseList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to activate license:', error)
      ElMessage.error('激活授权失败')
    }
  }
}

const handleDeactivate = async (license: License) => {
  try {
    await ElMessageBox.confirm(
      t('pages.licenses.confirm.deactivateMessage', { code: license.license_code }),
      t('pages.licenses.confirm.deactivateTitle'),
      {
        confirmButtonText: t('pages.licenses.confirm.confirm'),
        cancelButtonText: t('pages.licenses.confirm.cancel'),
        type: 'warning'
      }
    )

    await deactivateLicense(license.id)
    ElMessage.success(t('pages.licenses.message.deactivateSuccess'))
    loadLicenseList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to deactivate license:', error)
      ElMessage.error('停用授权失败')
    }
  }
}

const handleDelete = async (license: License) => {
  try {
    await ElMessageBox.confirm(
      t('pages.licenses.confirm.deleteMessage', { code: license.license_code }),
      t('pages.licenses.confirm.deleteTitle'),
      {
        confirmButtonText: t('pages.licenses.confirm.confirm'),
        cancelButtonText: t('pages.licenses.confirm.cancel'),
        type: 'warning'
      }
    )

    await deleteLicense(license.id)
    ElMessage.success(t('pages.licenses.message.deleteSuccess'))
    loadLicenseList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete license:', error)
      ElMessage.error(t('pages.licenses.message.deleteError'))
    }
  }
}

const handleFormSuccess = () => {
  formDialogVisible.value = false
  if (showLicenseList.value) {
    loadLicenseList()
  }
}

const getStatusType = (status: string) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'info'
    case 'expired':
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: string) => {
  return t(`pages.licenses.status.${status}`)
}

const formatDate = (dateString: string) => {
  if (!dateString) return '--'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 生命周期
onMounted(() => {
  loadCustomers()
})
</script>

<style scoped lang="scss">
@import '@/assets/styles/variables.scss';

.license-container {
  height: calc(100vh - 80px);
  width: 100%;
  overflow: hidden;
}

.main-content {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.background-section {
  position: relative;
  height: 100%;
  background-image: url('/src/assets/images/license-bg.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  transition: height 0.3s ease;

  &.collapsed {
    height: 300px;
  }
}

.background-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.center-title {
  position: relative;
  z-index: 2;
  margin-bottom: 60px;
}

.platform-title {
  font-family: 'PangMenZhengDao', 'Source Han Sans CN', sans-serif;
  font-size: 56px;
  font-weight: 400;
  color: #1C1C28;
  text-align: center;
  margin: 0;
  letter-spacing: 0.06em;
  line-height: 1.2;
}

.action-section {
  position: relative;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  max-width: 1200px;
  padding: 0 24px;
}

.action-row {
  display: flex;
  align-items: stretch;
  gap: 24px;
  width: 100%;
  max-width: 1200px;

  // 重置所有子元素的基线对齐
  > * {
    vertical-align: top;
    margin: 0;
  }
}

.customer-select-wrapper {
  position: relative;
  flex: 1;
  max-width: 900px;
  height: 60px;
  display: flex;
  align-items: center;
}

.customer-select {
  width: 100%;

  // 重置Element Plus的默认margin
  :deep(.el-select) {
    margin: 0;
  }

  // 重置el-input的默认样式
  :deep(.el-input) {
    margin: 0;
    height: 60px;
  }

  :deep(.el-input__wrapper) {
    height: 60px !important;
    font-size: 20px;
    background: rgba(255, 255, 255, 0.9);
    border: 1px solid #E3E3E3;
    border-radius: 8px;
    padding: 0 48px 0 16px;
    box-sizing: border-box;
    display: flex;
    align-items: center;
    margin: 0;
    box-shadow: none;

    .el-input__inner {
      color: #A9A9AF;
      font-size: 20px;
      line-height: 1.5;
      height: 100%;
      margin: 0;
      padding: 0;

      &::placeholder {
        color: #A9A9AF;
      }
    }
  }

  :deep(.el-input__suffix) {
    display: none;
  }

  // 确保选择器本身没有额外的margin或padding
  :deep(.el-select__wrapper) {
    margin: 0;
    padding: 0 20px;
  }
}

.action-buttons {
  display: flex;
  gap: 16px;
  flex-shrink: 0;
  align-items: center;
  height: 60px;
}

.action-btn {
  height: 40px !important;
  padding: 0 24px;
  font-size: 20px;
  font-weight: 500;
  letter-spacing: 0.15em;
  border-radius: 8px;
  border: none;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0;
  vertical-align: top;

  // 重置Element Plus按钮的默认样式
  &:deep(.el-button) {
    margin: 0;
    height: 40px;
    line-height: 1;
  }

  &.query-btn {
    min-width: 91px;
    background-color: $primary-color;

    &:hover {
      background-color: darken($primary-color, 10%);
    }
  }

  &.create-btn {
    min-width: 137px;
    background-color: $primary-color;

    &:hover {
      background-color: darken($primary-color, 10%);
    }
  }
}

.license-list-section {
  flex: 1;
  background: #fff;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.list-header {
  padding: 24px;
  border-bottom: 1px solid $border-color-light;
}

.search-controls {
  display: flex;
  gap: 16px;
  align-items: center;
  flex-wrap: wrap;
}

.search-input {
  min-width: 300px;
  flex: 1;
}

.date-picker {
  min-width: 300px;
}

.table-container {
  flex: 1;
  padding: 0 24px 24px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.license-table {
  flex: 1;
  overflow: auto;
}

.pagination-container {
  padding: 16px 0;
  display: flex;
  justify-content: center;
}

// 响应式设计
@media (max-width: 768px) {
  .platform-title {
    font-size: 36px;
  }

  .action-section {
    padding: 0 16px;
  }

  .action-row {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .customer-select-wrapper {
    max-width: none;
  }

  .action-buttons {
    justify-content: center;
    width: 100%;

    .action-btn {
      flex: 1;
      max-width: 200px;
    }
  }

  .search-controls {
    flex-direction: column;
    align-items: stretch;

    .search-input,
    .date-picker {
      min-width: auto;
      width: 100%;
    }
  }

  .table-container {
    padding: 0 16px 16px;
  }
}

@media (max-width: 480px) {
  .platform-title {
    font-size: 28px;
  }

  .action-btn {
    height: 32px;
    font-size: 16px;
  }

  .customer-select {
    :deep(.el-input__wrapper) {
      height: 48px;
      font-size: 16px;
    }
  }
}
</style>