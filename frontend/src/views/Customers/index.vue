<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-11-03 11:19:50
 * @FilePath: /frontend/src/views/Customers/index.vue
 * @Description: 客户管理页面  
-->
<template>
  <Layout app-name="Cedar-V" :page-title="getPageTitle()">
    <!-- 客户列表页面 -->
    <div v-if="!showCustomerForm && !showCustomerView" class="content-container">
      <!-- 顶部操作区域 -->
      <div class="top-actions">
        <!-- 左侧操作区域 -->
        <div class="left-actions">
          <el-button type="primary" class="add-customer-btn" @click="handleAddCustomer">
            {{ t('customers.actions.add') }}
          </el-button>
        </div>
        
        <!-- 右侧操作区域 -->
        <div class="right-actions">
          <!-- 筛选区域 -->
          <div class="filter-section">
            <span class="filter-label">{{ t('customers.filter.label') }}</span>
            <el-select
              v-model="filterCustomerType"
              :placeholder="t('customers.filter.customerType')"
              class="filter-select"
              clearable
              @change="handleFilterChange"
            >
              <el-option 
                v-for="option in customerTypeOptions" 
                :key="option.key" 
                :label="option.display" 
                :value="option.key" 
              />
            </el-select>
            <el-select
              v-model="filterCustomerLevel"
              :placeholder="t('customers.filter.customerLevel')"
              class="filter-select"
              clearable
              @change="handleFilterChange"
            >
              <el-option 
                v-for="option in customerLevelOptions" 
                :key="option.key" 
                :label="option.display" 
                :value="option.key" 
              />
            </el-select>
            <el-select
              v-model="filterStatus"
              :placeholder="t('customers.filter.status')"
              class="filter-select"
              clearable
              @change="handleFilterChange"
            >
              <el-option 
                v-for="option in statusOptions" 
                :key="option.key" 
                :label="option.display" 
                :value="option.key" 
              />
            </el-select>
          </div>
          
          <!-- 搜索区域 -->
          <div class="search-section">
            <el-input
              v-model="searchKeyword"
              :placeholder="t('customers.search.placeholder')"
              class="search-input"
            >
              <template #append>
                <el-button :icon="Search" @click="handleSearch" class="search-btn" />
              </template>
            </el-input>
          </div>
        </div>
      </div>

      <!-- 数据表格 -->
      <div class="table-container">
        <div class="table-wrapper">
          <el-table
            :data="tableData"
            v-loading="loading"
            :element-loading-text="t('customers.table.loading')"
            stripe
            border
            style="width: 100%;"
            :header-cell-style="{ backgroundColor: 'var(--app-bg-color)', color: 'var(--app-text-primary)' }"
            :row-style="getRowStyle"
            :max-height="'calc(100vh - 265px)'"
          >
          <el-table-column prop="customer_code" :label="t('customers.table.customerCode')" :width="170" show-overflow-tooltip  align="left">
            <template #default="scope">
              <span class="customer-code customer-code-link" @click="handleViewCustomerDetail(scope.row)" style="cursor:pointer; color:var(--el-color-primary)">{{ scope.row.customer_code }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="customer_name" :label="t('customers.table.customerName')" :width="180" show-overflow-tooltip  align="center">
            <template #default="scope">
              <span class="ellipsis-text">{{ scope.row.customer_name }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="customer_type_display" :label="t('customers.table.customerType')" :width="130"  align="center">
          </el-table-column>
          <el-table-column prop="contact_person" :label="t('customers.table.contactPerson')" :width="140"    align="center" />
          <el-table-column prop="email" :label="t('customers.table.email')" :width="160" show-overflow-tooltip  align="center">
             <template #default="scope">
              <span class="ellipsis-text">{{ scope.row.email }}</span>
            </template>
            </el-table-column>
          <el-table-column prop="customer_level_display" :label="t('customers.table.customerLevel')" :width="130" align="center">
           </el-table-column>
          <el-table-column prop="status_display" :label="t('customers.table.status')" :width="100"  align="center">
            <template #default="scope">
              <div class="status-tag" :class="getStatusClass(scope.row.status)">
                {{scope.row.status_display}}
              </div>
            </template>
          </el-table-column>
          <el-table-column :label="t('customers.table.createTime')" :width="140"  align="center">
            <template #default="scope">
              {{ formatDateShort(scope.row.created_at) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('customers.table.operation')" fixed="right" class-name="action-column" align="center" :width="320">
            <template #default="scope">
              <div class="action-buttons">
                <button class="action-btn primary" @click="handleViewLicense(scope.row)">{{ t('customers.actions.viewLicense') }}</button>
                <button class="action-btn primary" @click="handleEdit(scope.row)">{{ t('customers.actions.edit') }}</button>
                <button 
                  class="action-btn" 
                  :class="scope.row.status === 'active' ? 'warning' : 'success'" 
                  @click="handleDisable(scope.row)"
                >
                  {{ scope.row.status === 'active' ? t('customers.actions.disable') : t('customers.actions.enable') }}
                </button>
                <button class="action-btn danger" @click="handleDelete(scope.row)">{{ t('customers.actions.delete') }}</button>
              </div>
            </template>
          </el-table-column>
          </el-table>
        </div>
      </div>

      <!-- 分页组件 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[16, 32, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          :pager-count="7"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </div>
    
    <!-- 客户表单页面 -->
    <div v-if="showCustomerForm" class="form-page-container">
      <CustomerForm 
        :customer-id="isEditMode ? currentCustomerId : undefined"
        :is-edit="isEditMode"
        @save="handleFormSave"
        @cancel="handleFormCancel"
      />
    </div>

    <!-- 客户查看页面 -->
    <div v-if="showCustomerView" class="form-page-container">
      <CustomerView 
        :customer-id="currentCustomerId"
        @back="handleViewBack"
      />
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import Layout from '@/components/common/layout/Layout.vue'
import CustomerForm from './CustomerForm.vue'
import CustomerView from './CustomerView.vue'
import { Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox, ElLoading } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { 
  getCustomers, 
  deleteCustomer, 
  toggleCustomerStatus,
  type Customer,
  type CustomerQueryRequest 
} from '@/api/customer'
import { 
  getCustomerTypeEnums,
  getCustomerLevelEnums,
  getStatusEnums,
  type RawEnumItem
} from '@/api/enum'
import { formatDateShort } from '@/utils/date'
import { useRouter } from 'vue-router'
const router = useRouter()

// Customer类型已从API文件导入

const { t } = useI18n()
const searchKeyword = ref('')
const filterCustomerType = ref('')
const filterCustomerLevel = ref('')
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(16)
const total = ref(0)
const loading = ref(false)

const tableData = ref<Customer[]>([])
const showCustomerForm = ref(false)
const showCustomerView = ref(false)
const isEditMode = ref(false)
const currentCustomerId = ref<string>('')

// 获取枚举选项
const customerTypeOptions = ref<RawEnumItem[]>([])
const customerLevelOptions = ref<RawEnumItem[]>([])
const statusOptions = ref<RawEnumItem[]>([])

// 加载枚举数据
const loadEnums = async () => {
  try {
    const [typeRes, levelRes, statusRes] = await Promise.all([
      getCustomerTypeEnums(),
      getCustomerLevelEnums(),
      getStatusEnums()
    ])
    
    if (typeRes.code === '000000') {
      customerTypeOptions.value = typeRes.data.items
   
    }
    if (levelRes.code === '000000') {
      customerLevelOptions.value = levelRes.data.items

    }
    if (statusRes.code === '000000') {
      statusOptions.value = statusRes.data.items
    }
  } catch (error) {
    console.error(t('customers.message.loadEnumError'), error)
  }
}

// 页面标题
const getPageTitle = () => {
  if (showCustomerView.value) return t('customers.viewCustomer')
  if (showCustomerForm.value) return isEditMode.value ? t('customers.editCustomer') : t('customers.addCustomer')
  return t('customers.title')
}



const getRowStyle = ({ rowIndex }: { rowIndex: number }) => {
  return rowIndex % 2 === 1 ? { backgroundColor: 'var(--app-bg-color)' } : { backgroundColor: 'var(--app-content-bg)' }
}

// 加载数据
const loadData = async () => {
  try {
    loading.value = true
    const params: CustomerQueryRequest = {
      page: currentPage.value,
      page_size: pageSize.value,
      search: searchKeyword.value || undefined,
      customer_type: filterCustomerType.value as any || undefined,
      customer_level: filterCustomerLevel.value as any || undefined,
      status: filterStatus.value as any || undefined
    }
    
    const response = await getCustomers(params)
    
    if (response.code === '000000') {
      tableData.value = response.data.list
      total.value = response.data.total
    } else {
      ElMessage.error(response.message)
      tableData.value = []
      total.value = 0
    }
  } catch (error: any) {
    console.error('Load data error:', error)
    
    // 使用后端返回的错误信息
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(errorMessage)
    }
    tableData.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}


const handleSearch = () => {
  currentPage.value = 1
  loadData()
}

const handleFilterChange = () => {
  currentPage.value = 1
  loadData()
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  currentPage.value = 1
  loadData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  loadData()
}

const handleEdit = (row: Customer) => {
  showCustomerForm.value = true
  showCustomerView.value = false
  isEditMode.value = true
  currentCustomerId.value = row.id  // 设置当前客户ID
}

const handleDisable = async (row: Customer) => {
  const newStatus = row.status === 'active' ? 'disabled' : 'active'
  const actionText = newStatus === 'disabled' ? t('customers.actions.disable') : t('customers.actions.enable')
  
  try {
    await ElMessageBox.confirm(
      newStatus === 'disabled'
        ? t('customers.confirm.disableMessage', { name: row.customer_name })
        : t('customers.confirm.enableMessage', { name: row.customer_name }),
      newStatus === 'disabled'
        ? t('customers.confirm.disableTitle')
        : t('customers.confirm.enableTitle'),
      {
        confirmButtonText: t('customers.confirm.confirm'),
        cancelButtonText: t('customers.confirm.cancel'),
        type: 'warning'
      }
    )
    
    const loadingInstance = ElLoading.service({
      lock: true,
      text: newStatus === 'disabled'
        ? t('customers.message.disableLoading')
        : t('customers.message.enableLoading'),
      background: 'rgba(0, 0, 0, 0.7)'
    })
    
    try {
      const response = await toggleCustomerStatus(row.id, newStatus)
      
      if (response.code === '000000') {
        ElMessage.success(response.message)
        // 重新加载数据
        await loadData()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error: any) {
      console.error(`${actionText} error:`, error)
      
      // 使用后端返回的错误信息
      let errorMessage = error.backendMessage || error.response?.data?.message || error.message
      if (errorMessage) {
        ElMessage.error(errorMessage)
      }
    } finally {
      loadingInstance.close()
    }
  } catch {
    ElMessage.info(t('customers.message.statusChangeCancel', { action: actionText }))
  }
}

const getStatusClass = (status: string) => {
  return {
    'status-normal': status === 'active',
    'status-disabled': status === 'disabled'
  }
}

const handleAddCustomer = () => {
  showCustomerForm.value = true
  showCustomerView.value = false
  isEditMode.value = false
}

// --- 新增：客户编码点击跳转客户详情 ---
const handleViewCustomerDetail = (row: Customer) => {
  showCustomerView.value = true
  showCustomerForm.value = false
  currentCustomerId.value = row.id
}
// --- 修改：查看授权跳转到LinenseList页面并带参数 ---
const handleViewLicense = (row: Customer) => {
  router.push({
    path: '/licenses/list',
    query: { customerId: row.id, customerName: row.customer_name }
  })
}

const handleDelete = async (row: Customer) => {
  try {
    await ElMessageBox.confirm(
      t('customers.confirm.deleteMessage', { name: row.customer_name }),
      t('customers.confirm.deleteTitle'),
      {
        confirmButtonText: t('customers.confirm.confirm'),
        cancelButtonText: t('customers.confirm.cancel'),
        type: 'warning'
      }
    )
    
    const loadingInstance = ElLoading.service({
      lock: true,
      text: t('customers.message.deleteLoading'),
      background: 'rgba(0, 0, 0, 0.7)'
    })
    
    try {
      const response = await deleteCustomer(row.id)
      
      if (response.code === '000000') {
        ElMessage.success(response.message)
        // 重新加载数据
        await loadData()
      } else {
        ElMessage.error(response.message)
      }
    } catch (error: any) {
      console.error('Delete error:', error)
      
      // 使用后端返回的错误信息
      let errorMessage = error.backendMessage || error.response?.data?.message || error.message
      if (errorMessage) {
        ElMessage.error(errorMessage)
      }
    } finally {
      loadingInstance.close()
    }
  } catch {
    ElMessage.info(t('customers.message.deleteCancel'))
  }
}

// 处理客户表单操作
const handleFormSave = async (data: any) => {
  // 这里处理保存逻辑
  console.log('Save customer data:', data)
  showCustomerForm.value = false
  await loadData() // 重新加载数据
}

const handleFormCancel = () => {
  showCustomerForm.value = false
}

const handleViewBack = () => {
  showCustomerView.value = false
}

onMounted(async () => {
  await loadEnums()
  loadData()
})
</script>

<style lang="scss" scoped>
// Variables and mixins are auto-injected via Vite configuration
@use 'sass:math';

.content-container {
  // 所有CSS声明放在@include之前
  max-width: none !important; // 移除最大宽度限制，允许充满整个屏幕
  height: 100%; // 占满父容器高度
  padding: 0; // 移除padding，由父容器处理
  background-color: var(--app-content-bg);
  display: flex;
  flex-direction: column;
  overflow: hidden; // 防止内容超出
  box-sizing: border-box;
  width: 100%;
  margin: 0 auto;
  
  @include mobile {
    padding: $spacing-base;
  }
}

.form-page-container {
  max-width: none !important; 
  min-height: calc(100vh - 80px);
  padding: $spacing-large;
  background: var(--app-bg-color);
  display: flex;
  flex-direction: column;
  
  width: 100%;
  max-width: $breakpoint-desktop;
  margin: 0 auto;
  padding: 0 $spacing-medium;
  
  @include mobile {
    padding: 0 $spacing-base;
  }
}

.top-actions {
  @include flex-between;
  gap: $spacing-large;
  margin-bottom: $spacing-medium;
  padding: $spacing-large $spacing-large 0 $spacing-large; // 只在顶部和两侧添加padding
  flex-shrink: 0;
  
  @include mobile {
    @include flex-responsive;
    align-items: stretch;
  }
}

.left-actions {
  @include flex-center-vertical;
  gap: $spacing-large;
  
  @include mobile {
    justify-content: flex-start;
    width: 100%;
  }
}

.add-customer-btn {
  height: 32px;  // 符合8px栅格
  padding: $spacing-small $spacing-medium;
  font-size: $font-size-base;
  font-weight: $font-weight-primary;
  border-radius: $border-radius-small;

  @include button-primary;
}

.right-actions {
  @include flex-center-vertical;
  gap: $spacing-extra-large;
  min-width: 0;
  flex-shrink: 1;
  
  @include mobile {
    flex-direction: column;
    gap: $spacing-medium;
    width: 100%;
  }
}

.filter-section {
  @include flex-center-vertical;
  gap: $spacing-small;
  min-width: 288px;  // 改为288px (36*8)
  flex-shrink: 1;
  
  @include mobile {
    width: 100%;
    flex-wrap: wrap;
    gap: $spacing-small;
  }
  
  @include mobile {
    flex-direction: column;
    align-items: stretch;
  }
}

.filter-label {
  font-size: $font-size-base;
  font-weight: $font-weight-primary;
  color: var(--app-text-primary);
  margin-right: $spacing-small;
  min-width: 48px;   // 改为48px (6*8)
  
  @include mobile {
    width: auto;
    margin-right: $spacing-small;
  }
  
  @include mobile {
    align-self: flex-start;
  }
}

.search-section {
  min-width: 288px;  // 改为288px (36*8)
  flex-shrink: 1;
  
  @include mobile {
    width: 100%;
  }
}

.search-input {
  width: 100%;
}

.table-container {
  width: calc(100% - #{$spacing-large * 2}); // 计算宽度减去左右padding
  margin: 0 $spacing-large; // 水平居中并添加左右边距
  border-radius: $border-radius-base;
  overflow: hidden;
  margin-bottom: $spacing-medium;
  flex: 1; // 自动占满剩余空间
  display: flex;
  flex-direction: column;
  
  @include mobile {
    @include table-responsive;
  }
}

.pagination-container {
  @include flex-end;
  gap: $spacing-small;
  padding: 0 $spacing-large $spacing-large $spacing-large; // 添加下部和两侧padding
  flex-shrink: 0;
  
  @include mobile {
    @include flex-center;
    flex-wrap: wrap;
    gap: $spacing-extra-small;
  }
}

.table-wrapper {
  width: 100%;
  height: 100%; // 占满表格容器高度
  position: relative;
  overflow: auto; // 启用滚动以支持固定列
  @include smooth-scroll;

  @include mobile {
    -webkit-overflow-scrolling: touch;
  }
}

:deep(.el-table) {
  width: 100% !important;
  min-width: 1408px; /* 改为1408px (176*8) */
  border: 1px solid var(--app-border-light);
  table-layout: fixed; /* 固定表格布局防止错位 */
  
  /* 表头和表体宽度一致 */
  .el-table__header-wrapper,
  .el-table__body-wrapper {
    width: 100% !important;
    overflow-x: auto !important; /* 允许水平滚动以支持固定列 */
  }
  
  .el-table__body-wrapper {
    overflow-y: auto !important;
    /* 使用flex布局自动计算高度 */
    flex: 1 !important;
    /* 预留滚动条宽度空间 */
    margin-right: 0 !important;
    padding-right: 0 !important;
  }
  
  /* 表格主体和表头宽度一致 */
  .el-table__header,
  .el-table__body {
    width: 100% !important;
    table-layout: fixed !important;
  }
  
  /* 确保所有列宽度一致 */
  th, td {
    box-sizing: border-box;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin: 0 !important;
  }
  
  /* 列宽度分配 - 防止错位 */
  colgroup {
    width: 100% !important;
  }
}

/* 分页组件样式 - 匹配Figma设计 */
:deep(.el-pagination) {
  justify-content: flex-end;
  gap: 8px;
  --el-pagination-font-size: 14px;
  --el-pagination-bg-color: var(--app-content-bg);
  --el-pagination-text-color: var(--app-text-primary);
  --el-pagination-border-color: var(--app-border-color);
  --el-pagination-hover-color: #019C7C;
  font-family: 'Source Han Sans CN', sans-serif;
}

/* 页码按钮样式 */
:deep(.el-pagination .btn-prev),
:deep(.el-pagination .btn-next),
:deep(.el-pagination .el-pager li) {
  width: 32px;
  height: 32px;
  min-width: 32px;
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  margin: 0 4px;
  background: var(--app-content-bg);
  color: var(--app-text-primary);
  font-family: 'Segoe UI', sans-serif;
  font-weight: 400;
  font-size: 14px;
  line-height: 1.33;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 当前页样式 */
:deep(.el-pagination .el-pager li.is-active) {
  background: var(--el-color-primary) !important;
  border-color: var(--el-color-primary) !important;
  color: white !important;
}

/* 悬停效果 */
:deep(.el-pagination .btn-prev:hover),
:deep(.el-pagination .btn-next:hover),
:deep(.el-pagination .el-pager li:hover) {
  color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

/* 总数和每页条数样式 */
:deep(.el-pagination .el-pagination__total),
:deep(.el-pagination .el-pagination__sizes) {
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 350;
  font-size: 14px;
  color: var(--app-text-primary);
}

/* 每页条数选择器 */
:deep(.el-pagination .el-pagination__sizes .el-select) {
  width: 88px;  // 88px符合8px栅格(11*8)
}

:deep(.el-pagination .el-pagination__sizes .el-select .el-select__wrapper) {
  height: 32px;
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  background: var(--app-content-bg);
}

/* 跳转输入框 */
:deep(.el-pagination .el-pagination__jump) {
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 350;
  font-size: 14px;
  color: var(--app-text-primary);
  margin-left: 8px;
}

:deep(.el-pagination .el-pagination__jump .el-input) {
  width: 48px;  // 符合8px栅格
  margin: 0 8px;
}

:deep(.el-pagination .el-pagination__jump .el-input__wrapper) {
  height: 32px;
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  background: var(--app-content-bg);
}

:deep(.el-pagination .el-pagination__jump .el-input__inner) {
  text-align: center;
  font-family: 'Source Han Sans CN', sans-serif;
  font-size: 14px;
  color: var(--app-text-primary);
}

/* 桌面端响应式 - 中等屏幕 */
@media (max-width: 1200px) and (min-width: 769px) {
  .right-actions {
    gap: 24px;  // 改为24px符合8px栅格
  }
  
  .filter-section {
    min-width: 240px;
  }
  
  .search-section {
    min-width: 240px;
  }
  
  /* 平板端保持一行布局 */
  .action-buttons {
    gap: 8px;  // 改为8px符合8px栅格
    flex-wrap: nowrap;
    justify-content: center;
    width: auto;
  }
  
  .action-btn {
    padding: 3px 8px;
    font-size: 11px;
    flex-shrink: 1;
    min-width: 40px;
  }
}

/* 1920*1080分辨率专用修复 */
@media (min-width: 1900px) and (max-width: 1940px) and (min-height: 1060px) and (max-height: 1100px) {
  .table-wrapper {
    width: 100% !important;
  }
  
  :deep(.el-table) {
    width: 100% !important;
    min-width: 100% !important;
    
    .el-table__header-wrapper,
    .el-table__body-wrapper {
      width: 100% !important;
    }
    
    .el-table__header,
    .el-table__body {
      width: 100% !important;
    }
    
    /* 1920*1080下的列宽度优化 */
    th, td {
      min-width: auto !important;
    }
  }
}

/* 桌面端：使用vw单位实现2K/4K适配 */
@media (min-width: 1025px) {
  .content-container {
    height: 100% !important; /* 占满父容器高度 */
    padding: 0 !important; /* 移除padding */
    width: 100%; /* 充满整个屏幕 */
    margin: 0;
    box-sizing: border-box;
    display: flex; /* 桌面端使用flex布局传递高度给子组件 */
    flex-direction: column;
    background-color: var(--app-content-bg);
    overflow: hidden; /* 防止内容超出 */
  }
  
  .form-page-container {
    height: calc(100vh - 4.17vw);
    padding: 1.25vw;
    width: 100%;
    margin: 0;
    box-sizing: border-box;
    display: flex;
    flex-direction: column;
    background: var(--app-bg-color);
  }

  
  .table-container {
    border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
    margin-bottom: 0.83vw; /* 16px/1920 = 0.83vw */
    margin-left: 1.25vw; /* 24px/1920 = 1.25vw */
    margin-right: 1.25vw;
    width: calc(100% - 2.5vw); /* 减去左右边距 */
    flex: 1; /* 自动占满剩余空间 */
    display: flex;
    flex-direction: column;
  }
  
  .table-wrapper {
    width: 100%;
    flex: 1; /* 占满表格容器高度 */
  }
  
  .pagination-container {
    flex-shrink: 0;
    gap: 0.42vw; /* 8px/1920 = 0.42vw */
    padding: 0 1.25vw 1.25vw 1.25vw; /* 添加padding */
  }
  
  .top-actions {
    padding: 1.25vw 1.25vw 0 1.25vw; /* 添加padding */
  }
  
  .right-actions {
    gap: 2.08vw; /* 40px/1920 = 2.08vw */
  }
  
  .filter-section {
    min-width: 14.58vw; /* 280px/1920 = 14.58vw */
    gap: 0.21vw; /* 4px/1920 = 0.21vw */
  }
  
  .search-section {
    min-width: 14.58vw; /* 280px/1920 = 14.58vw */
  }
  
  .action-buttons {
    gap: 0.42vw; /* 8px/1920 = 0.42vw */
    justify-content: center;
    width: auto;
  }
  
  .action-btn {
    padding: 0.16vw 0.625vw; /* 3px 12px/1920 = 0.16vw 0.625vw */
    font-size: 0.625vw; /* 12px/1920 = 0.625vw */
    height: 1.46vw; /* 28px/1920 = 1.46vw */
    min-width: fit-content;
  }
  
  /* 表格样式使用vw - 确保表格充满容器且对齐 */
  .table-wrapper {
    width: 100% !important;
  }
  
  :deep(.el-table) {
    width: 100% !important;
    min-width: 73vw !important; /* 1400px/1920 = 73vw */
    
    .el-table__header-wrapper,
    .el-table__body-wrapper {
      width: 100% !important;
      overflow-x: hidden !important;
    }
    
    .el-table__body-wrapper {
      overflow-y: auto !important;
      flex: 1 !important; /* 使用flex自动计算高度 */
    }
    
    .el-table__header,
    .el-table__body {
      width: 100% !important;
      table-layout: fixed !important;
    }
    
    /* 桌面端列宽度优化和对齐修复 */
    th, td {
      overflow: hidden !important;
      text-overflow: ellipsis !important;
      white-space: nowrap !important;
      padding: 0 !important; /* 12px 20px/1920 = 0.625vw 1.04vw */
      box-sizing: border-box !important;
      margin: 0 !important;
    }
    
    /* 操作列自适应宽度 */
    .action-column {
      width: auto !important;
      min-width: max-content !important;
    }
  }
  
  :deep(.el-table .el-table__header-wrapper) {
    th,
    th.el-table__cell {
      height: 2.5vw !important; /* 48px/1920 = 2.5vw */
      padding: 0 !important; /* 12px 20px/1920 = 0.625vw 1.04vw */
      margin: 0 !important;
      box-sizing: border-box !important;
      font-size: 0.83vw !important; /* 16px/1920 = 0.83vw */
      line-height: 1.25vw !important; /* 24px/1920 = 1.25vw */
    }
  }

  :deep(.el-table .el-table__body-wrapper) {
    td.el-table__cell {
      height: 50px !important;
      padding: 0!important;
      margin: 0 !important;
      box-sizing: border-box !important;
      font-size: 14px !important;
      line-height: 24px !important;
    }
  }
  
  /* 筛选器和搜索框使用vw - 提高优先级 */
  :deep(.filter-select) {
    width: 6.25vw !important; /* 120px/1920 = 6.25vw */
    height: 1.875vw !important; /* 36px/1920 = 1.875vw */
  }
  
  :deep(.filter-select .el-select__wrapper) {
    min-height: 1.875vw !important; /* 36px/1920 = 1.875vw */
    height: 1.875vw !important;
    padding: 0 0.625vw !important; /* 12px/1920 = 0.625vw */
  }
  
  :deep(.search-input) {
    height: 1.875vw !important; /* 36px/1920 = 1.875vw */
  }
  
  :deep(.search-input .el-input__wrapper) {
    height: 1.875vw !important; /* 36px/1920 = 1.875vw */
  }
  
  :deep(.search-input .el-input-group__append) {
    width: 2.5vw !important; /* 48px/1920 = 2.5vw */
    height: 1.875vw !important; /* 36px/1920 = 1.875vw */
  }
}

/* 移动端响应式布局 - 修复表头错位 */
/* content-container和form-page-container响应式样式 */
@media (max-width: 768px) {
  .content-container, .form-page-container {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .content-container, .form-page-container {
    padding: 12px;
  }
}

@media (max-width: 768px) {
  /* 表格容器移动端优化 */
  .table-wrapper {
    overflow-x: scroll !important;
    -webkit-overflow-scrolling: touch;
    width: 100% !important;
  }
  
  /* 移动端表格样式修复 */
  :deep(.el-table) {
    min-width: 1200px !important; /* 移动端最小宽度 */
    
    .el-table__header-wrapper,
    .el-table__body-wrapper {
      width: 100% !important;
      overflow-x: visible !important;
    }
    
    .el-table__body-wrapper {
      max-height: calc(100vh - 300px) !important;
    }
    
    /* 移动端表头和表体对齐 */
    .el-table__header,
    .el-table__body {
      width: 100% !important;
      min-width: 1200px !important;
    }
    
    /* 移动端列宽度修复 */
    th, td {
      min-width: 80px !important;
      white-space: nowrap !important;
      padding: 8px 12px !important;
    }
    
    /* 移动端操作列自适应宽度 */
    .action-column {
      width: auto !important;
      min-width: max-content !important;
    }
  }
  /* 设置表格容器 */
  .table-container {
    --action-column-width: 160px;
    overflow-x: auto;
    width: 100%;
  }
  
  .table-wrapper {
    min-width: 100%;
    overflow-x: auto;
  }
  
  /* 确保表格在移动端正确显示 */
  :deep(.el-table .el-table__body-wrapper td) {
    vertical-align: middle !important;
  }
  
  /* 操作列保持顶部对齐 */
  :deep(.el-table .el-table__body-wrapper td:last-child) {
    vertical-align: top !important;
  }
  
  :deep(.el-table .el-table__body-wrapper .cell) {
    white-space: normal !important;
    overflow: visible !important;
  }
  
  /* 特别针对操作列的单元格 */
  :deep(.el-table .el-table__body-wrapper td:last-child .cell) {
    white-space: normal !important;
    overflow: visible !important;
    text-overflow: initial !important;
  }
  .top-actions {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .left-actions {
    justify-content: flex-start;
  }
  
  .right-actions {
    flex-direction: column;
    gap: 16px;
  }
  
  .filter-section {
    width: 100%;
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .filter-label {
    width: auto;
    margin-right: 8px;
  }
  
  .search-section {
    width: 100%;
  }
  
  /* 移动端操作列宽度调整 - 覆盖所有可能的Element Plus样式 */
  :deep(.el-table) {
    --el-table-action-column-width: 160px;
  }
  
  :deep(.el-table__header-wrapper colgroup col:last-child),
  :deep(.el-table__body-wrapper colgroup col:last-child),
  :deep(.el-table__fixed-right colgroup col:last-child) {
    width: 160px !important;
  }
  
  :deep(.el-table__header-wrapper .action-column),
  :deep(.el-table__body-wrapper .action-column),
  :deep(.el-table th:last-child),
  :deep(.el-table td:last-child) {
    width: 160px !important;
    min-width: 160px !important;
    max-width: 160px !important;
  }
  
  :deep(.el-table__fixed-right) {
    width: 160px !important;
    right: 0 !important;
  }
  
  :deep(.el-table__fixed-right .el-table__fixed-body-wrapper),
  :deep(.el-table__fixed-right .el-table__fixed-header-wrapper),
  :deep(.el-table__fixed-right .el-table__fixed-footer-wrapper) {
    width: 160px !important;
  }
  
  :deep(.el-table__fixed-right-patch) {
    width: 160px !important;
  }
  
  /* 强制覆盖内联样式 */
  :deep(.el-table [style*="width: 300px"]) {
    width: 160px !important;
  }
  
  /* 移动端操作按钮样式 - 2行2列布局 */
  .action-buttons {
    display: flex !important;
    flex-wrap: wrap !important;
    gap: 4px !important;
    width: 100% !important;
    max-width: 150px !important;
    justify-content: space-between !important;
    align-content: flex-start !important;
    padding: 6px 0 !important;
  }
  
  .action-btn {
    flex: 0 0 calc(50% - 2px) !important;
    width: calc(50% - 2px) !important;
    max-width: calc(50% - 2px) !important;
    padding: 2px 3px !important;
    font-size: 9px !important;
    height: 20px !important;
    min-width: 0 !important;
    white-space: nowrap !important;
    text-overflow: ellipsis !important;
    overflow: hidden !important;
    box-sizing: border-box !important;
    border-radius: 2px !important;
  }
  
  /* 移动端分页样式 */
  :deep(.el-pagination) {
    justify-content: center;
    flex-wrap: wrap;
    gap: 4px;
  }
  
  :deep(.el-pagination .el-pagination__total),
  :deep(.el-pagination .el-pagination__sizes) {
    order: -1;
    width: 100%;
    text-align: center;
    margin-bottom: 8px;
  }
}

@media (max-width: 480px) {
  .filter-section {
    flex-direction: column;
    align-items: stretch;
  }
  
  .filter-label {
    align-self: flex-start;
  }
  
  .filter-select {
    width: 100% !important;
  }
  
  /* 超小屏幕进一步压缩 */
  .action-buttons {
    max-width: 130px !important;
    gap: 3px !important;
  }
  
  .action-btn {
    padding: 1px 2px !important;
    font-size: 8px !important;
    height: 18px !important;
    flex: 0 0 calc(50% - 1.5px) !important;
    width: calc(50% - 1.5px) !important;
    max-width: calc(50% - 1.5px) !important;
  }
}

/* 筛选器样式 */
:deep(.filter-select) {
  width: 120px;  // 符合8px栅格(15*8)
  height: 40px;  // 改为40px符合8px栅格
}

:deep(.filter-select .el-select__wrapper) {
  align-items: center;
  background-color: var(--el-fill-color-blank);
  border-radius: var(--el-border-radius-base);
  box-shadow: 0 0 0 1px #DCDEE2 inset;
  box-sizing: border-box;
  cursor: pointer;
  display: flex;
  font-size: 14px;
  gap: 8px;  // 改为8px符合8px栅格
  line-height: 24px;
  min-height: 36px;
  padding: 0px 12px;
  position: relative;
  text-align: left;
  transform: translateZ(0);
  transition: var(--el-transition-duration);
}

:deep(.filter-select .el-select__selected-item) {
  font-size: 14px;
  color: var(--app-text-secondary);
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 400;
  line-height: 24px;
}

:deep(.filter-select .el-select__placeholder) {
  font-size: 14px;
  color: var(--app-text-secondary);
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 400;
  line-height: 24px;
}

:deep(.filter-select .el-select__suffix) {
  display: flex;
  align-items: center;
}

:deep(.filter-select .el-select__icon) {
  width: 14px;
  height: 14px;
  color: var(--app-text-secondary);
}

:deep(.filter-select.is-focus .el-select__wrapper) {
  box-shadow: 0 0 0 1px #019C7C inset;
}

:deep(.filter-select:hover .el-select__wrapper) {
  box-shadow: 0 0 0 1px #019C7C inset;
}

/* 搜索框样式 */
:deep(.search-input) {
  height: 40px;  // 改为40px符合8px栅格
}

:deep(.search-input .el-input__wrapper) {
  height: 40px;  // 改为40px符合8px栅格
  border-color: var(--app-border-color);
  border-radius: 4px 0 0 4px;
  background-color: var(--app-content-bg);
  display: flex;
  align-items: center;
}

:deep(.search-input .el-input__inner) {
  height: 26px;
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 400;
  font-size: 14px;
  line-height: 24px;  // 改为24px符合8px栅格
  color: var(--app-text-primary);
}

:deep(.search-input .el-input__inner::placeholder) {
  color: var(--app-text-secondary);
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 400;
  font-size: 14px;
}

:deep(.search-input.is-focus .el-input__wrapper) {
  border-color: var(--el-color-primary);
  box-shadow: 0 0 0 1px #019C7C;
}

:deep(.search-input .el-input-group__append) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
  border-radius: 0 4px 4px 0;
  padding: 10px;
  width: 48px;  // 符合8px栅格
  height: 40px;  // 改为40px符合8px栅格
  display: flex;
  justify-content: center;
  align-items: center;
}

:deep(.search-input .el-input-group__append .el-button) {
  background-color: transparent;
  border: none;
  color: white;
  padding: 0;
  width: 20px;
  height: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
}

:deep(.search-input .el-input-group__append .el-button:hover) {
  background-color: transparent;
  color: white;
}

/* 客户编码样式 */
.customer-code {
  color: var(--el-color-primary);
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 400;
  font-size: 14px;
  line-height: 24px;  // 改为24px符合8px栅格
}

/* 状态标签样式 */
.status-tag {
  padding: $spacing-extra-small $spacing-small;
  border-radius: $border-radius-small;
  font-size: $font-size-base;
  font-weight: $font-weight-secondary;
  height: 24px;  // 改为24px符合8px栅格
  min-width: fit-content;

  @include flex-center;
  @include text-ellipsis;
  
  &.status-normal {
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
  }

  &.status-disabled {
    background: rgba(244, 108, 108, 0.08);
    color: #F46C6C;
  }

  &.status-paused {
    background: rgba(230, 162, 60, 0.08);
    color: #E6A23C;
  }
}

/* 操作按钮组样式 */
.action-buttons {
  gap: $spacing-small;
  padding: 0;
  width: auto;
  min-width: max-content;

  @include flex-center-vertical;
  justify-content: center;
  flex-wrap: nowrap;
  overflow: visible;
  
  @include mobile {
    display: flex !important;
    flex-wrap: wrap !important;
    gap: $spacing-extra-small !important;
    width: 100% !important;
    max-width: 150px !important;
    justify-content: space-between !important;
    align-content: flex-start !important;
    padding: 0 !important;
  }
  
  @include mobile {
    max-width: 130px !important;
    gap: math.div($spacing-extra-small, 2) !important;
  }
}

.action-btn {
  padding: $spacing-extra-small $spacing-base;
  height: 32px;  // 改为32px符合8px栅格
  border-radius: $border-radius-small;
  font-size: $font-size-small;
  font-weight: $font-weight-primary;
  line-height: 20px;  // 改为20px符合8px栅格
  flex-shrink: 0;
  width: auto;
  min-width: fit-content;
  white-space: nowrap;

  @include button-base;
  
  @include mobile {
    flex: 0 0 calc(50% - #{math.div($spacing-extra-small, 2)}) !important;
    width: calc(50% - #{math.div($spacing-extra-small, 2)}) !important;
    max-width: calc(50% - #{math.div($spacing-extra-small, 2)}) !important;
    padding: math.div($spacing-extra-small, 2) $spacing-extra-small !important;
    font-size: 9px !important;
    height: 20px !important;
    min-width: 0 !important;
  }
  
  @include mobile {
    padding: 1px 2px !important;
    font-size: 8px !important;
    height: 18px !important;
    flex: 0 0 calc(50% - 1.5px) !important;
    width: calc(50% - 1.5px) !important;
    max-width: calc(50% - 1.5px) !important;
  }
}

.action-btn {
  &.primary {
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);

    @include non-touch-device {
      &:hover {
        background: var(--el-color-primary-light-8);
      }
    }
  }

  &.danger {
    background: rgba(244, 108, 108, 0.08);
    color: #F46C6C;

    @include non-touch-device {
      &:hover {
        background: rgba(244, 108, 108, 0.15);
      }
    }
  }

  &.warning {
    background: rgba(230, 162, 60, 0.08);
    color: #E6A23C;

    @include non-touch-device {
      &:hover {
        background: rgba(230, 162, 60, 0.15);
      }
    }
  }

  &.success {
    background: rgba(103, 194, 58, 0.08);
    color: #67C23A;

    @include non-touch-device {
      &:hover {
        background: rgba(103, 194, 58, 0.15);
      }
    }
  }
}

:deep(.el-table__body tr:hover > td) {
  background: #e3f7f1 !important; /* 明亮突出，适合主色系，可根据品牌色微调 */
  transition: background 0.2s;
}
[data-theme='dark'] :deep(.el-table__body tr:hover > td) {
  background: rgba(16, 185, 129, 0.15) !important; /* 暗色模式下主色加透明度 */
}
</style>