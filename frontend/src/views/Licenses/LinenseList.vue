<template>
  <div class="license-list-page">
    <!-- 顶部操作区 -->
    <div class="top-actions">
      <el-button type="primary" @click="handleCreate"> {{ t('pages.licenses.list.createLicense') }} </el-button>

      <div class="filters">
        <span class="filter-label">{{ t('pages.licenses.list.filter.label') }}</span>

        <el-input v-model="filterCode" :placeholder="t('pages.licenses.list.filter.codePlaceholder')" clearable class="filter-input" />

        <el-select v-model="filterStatus" :placeholder="t('pages.licenses.list.filter.statusPlaceholder')" clearable class="filter-select">
          <el-option
            v-for="option in statusOptions"
            :key="option.key"
            :label="option.display"
            :value="option.key"
          />
        </el-select>

        <el-button type="primary" @click="handleQuery"> {{ t('pages.licenses.list.filter.query') }} </el-button>
      </div>
    </div>

    <!-- 数据表格 -->
    <div class="table-container">
      <div class="table-wrapper">
        <el-table
          :data="tableData"
          v-loading="loading"
          :element-loading-text="t('pages.licenses.list.table.loading')"
          stripe
          border
          style="width: 100%"
          :header-cell-style="{
            backgroundColor: 'var(--app-bg-color)',
            color: 'var(--app-text-primary)'
          }"
          :row-style="getRowStyle"
          :max-height="'calc(100vh - 280px)'"
        >
          <el-table-column
            prop="code"
            :label="t('pages.licenses.list.table.code')"
            :width="220"
            show-overflow-tooltip
            align="center"
          >
            <template #default="scope">
              <span
                class="license-code"
                style="white-space: nowrap; overflow: visible; text-overflow: initial"
                >{{ scope.row.code }}</span
              >
            </template>
          </el-table-column>
          <el-table-column prop="status" :label="t('pages.licenses.list.table.status')" :width="100" align="center">
            <template #default="scope">
              <div class="status-tag" :class="getStatusClass(scope.row.status)">
                {{ scope.row.status_display }}
              </div>
            </template>
          </el-table-column>
          <el-table-column
            :label="t('pages.licenses.list.table.activationProgress')"
            :width="220"
            align="center"
          >
            <template #default="scope">
              <div
                class="progress-container"
                :class="{ 'zero-activated': getActivatedCount(scope.row) === 0 }"
              >
                <el-progress
                  :percentage="getActivationProgress(scope.row)"
                  :stroke-width="8"
                  :show-text="true"
                  :text-inside="true"
                  :format="() => `${getActivatedCount(scope.row)}/${scope.row.max_activations}`"
                />
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="end_date" :label="t('pages.licenses.list.table.endDate')" :width="180" align="center">
            <template #default="scope">
              {{ formatDate(scope.row.end_date) }}
            </template>
          </el-table-column>
          <el-table-column
            prop="description"
            :label="t('pages.licenses.list.table.description')"
            :min-width="200"
            show-overflow-tooltip
            align="center"
          >
            <template #default="scope">
              <span class="ellipsis-text">{{ scope.row.description }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('pages.licenses.list.table.operation')" fixed="right" class-name="action-column" align="center">
            <template #default="scope">
              <div class="action-buttons">
                <button class="action-btn primary" @click="handleDetail(scope.row)">{{ t('pages.licenses.list.actions.detail') }}</button>
                <button
                  class="action-btn"
                  :class="scope.row.is_locked ? 'success' : 'warning'"
                  @click="scope.row.is_locked ? handleUnlock(scope.row) : handleLock(scope.row)"
                >
                  {{ scope.row.is_locked ? t('pages.licenses.list.actions.unlock') : t('pages.licenses.list.actions.lock') }}
                </button>
                <button class="action-btn danger" @click="handleDelete(scope.row)">{{ t('pages.licenses.list.actions.delete') }}</button>
              </div>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <!-- 分页组件 -->
    <div class="pagination-container">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[16, 32, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        :pager-count="7"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  getLicenses,
  lockAuthorizationCode,
  deleteLicense,
  type AuthorizationCode,
  type LicenseQueryRequest
} from '@/api/license'
import { getAuthorizationStatusEnums, type RawEnumItem } from '@/api/enum'
import { formatDate } from '@/utils/date'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

// 从路由参数获取客户信息
const customerInfo = computed(() => {
  return {
    id: (route.query.customerId as string) || '',
    name: (route.query.customerName as string) || ''
  }
})

// 响应式数据
const loading = ref(false)
const filterStatus = ref('')
const filterCode = ref('')
const tableData = ref<AuthorizationCode[]>([])
const currentPage = ref(1)
const pageSize = ref(16)
const total = ref(0)
const statusOptions = ref<RawEnumItem[]>([])

// 方法

const getRowStyle = ({ rowIndex }: { rowIndex: number }) => {
  return rowIndex % 2 === 1
    ? { backgroundColor: 'var(--app-bg-color)' }
    : { backgroundColor: 'var(--app-content-bg)' }
}

const getStatusClass = (status: string) => {
  return {
    'status-normal': status === 'normal',
    'status-locked': status === 'locked',
    'status-expired': status === 'expired'
  }
}

const getActivatedCount = (row: AuthorizationCode) => {
  if (typeof row.current_activations === 'number') {
    return row.current_activations
  }
  if (typeof row.activated_licenses_count === 'number') {
    return row.activated_licenses_count
  }
  return 0
}

const getActivationProgress = (row: AuthorizationCode) => {
  if (!row.max_activations || row.max_activations === 0) {
    return 0
  }
  const activated = getActivatedCount(row)
  return Math.round((activated / row.max_activations) * 100)
}


const handleCreate = () => {
  router.push({
    name: 'licenses-create',
    query: {
      customerId: customerInfo.value.id,
      customerName: customerInfo.value.name
    }
  })
}

const handleQuery = () => {
  console.log('点击查询按钮')
  loadData()
}

const handleDetail = (row: any) => {
  // 跳转到授权详情页面
  router.push({
    name: 'licenses-view',
    params: {
      id: row.id
    },
    query: {
      customerId: customerInfo.value.id,
      customerName: customerInfo.value.name
    }
  })
}

const handleLock = async (row: any) => {
  ElMessageBox.confirm(t('pages.licenses.list.confirm.lockMessage'), t('pages.licenses.list.confirm.lockTitle'), {
    confirmButtonText: t('pages.licenses.list.confirm.confirm'),
    cancelButtonText: t('pages.licenses.list.confirm.cancel'),
    type: 'warning'
  })
    .then(async () => {
      try {
        await lockAuthorizationCode(row.id, {
          is_locked: true,
          lock_reason: '手动锁定',
          reason: '管理员手动锁定'
        })
        ElMessage.success(t('pages.licenses.list.message.lockSuccess'))
        loadData()
      } catch (error) {
        console.error('锁定失败:', error)
        ElMessage.error(t('pages.licenses.list.message.lockError'))
      }
    })
    .catch(() => {
      // 取消操作
    })
}

const handleUnlock = async (row: any) => {
  ElMessageBox.confirm(t('pages.licenses.list.confirm.unlockMessage'), t('pages.licenses.list.confirm.unlockTitle'), {
    confirmButtonText: t('pages.licenses.list.confirm.confirm'),
    cancelButtonText: t('pages.licenses.list.confirm.cancel'),
    type: 'warning'
  })
    .then(async () => {
      try {
        await lockAuthorizationCode(row.id, {
          is_locked: false,
          lock_reason: '手动解锁',
          reason: '管理员手动解锁'
        })
        ElMessage.success(t('pages.licenses.list.message.unlockSuccess'))
        loadData()
      } catch (error) {
        console.error('解锁失败:', error)
        ElMessage.error(t('pages.licenses.list.message.unlockError'))
      }
    })
    .catch(() => {
      // 取消操作
    })
}

const handleDelete = async (row: any) => {
  ElMessageBox.confirm(t('pages.licenses.list.confirm.deleteMessage'), t('pages.licenses.list.confirm.deleteTitle'), {
    confirmButtonText: t('pages.licenses.list.confirm.deleteConfirm'),
    cancelButtonText: t('pages.licenses.list.confirm.cancel'),
    type: 'warning'
  })
    .then(async () => {
      try {
        await deleteLicense(row.id)
        ElMessage.success(t('pages.licenses.list.message.deleteSuccess'))
        loadData()
      } catch (error) {
        console.error('删除失败:', error)
        ElMessage.error(t('pages.licenses.list.message.deleteError'))
      }
    })
    .catch(() => {
      // 取消操作
    })
}

const handleSizeChange = (size: number) => {
  pageSize.value = size
  loadData()
}

const handlePageChange = (page: number) => {
  currentPage.value = page
  loadData()
}

const loadStatusOptions = async () => {
  try {
    const response = await getAuthorizationStatusEnums()
    statusOptions.value = response.data.items
    console.log('状态选项:', statusOptions.value)
  } catch (error) {
    console.error('加载状态选项失败:', error)
    // 如果接口失败，使用默认选项
    statusOptions.value = [
      { key: 'normal', display: '正常' },
      { key: 'locked', display: '锁定' },
      { key: 'expired', display: '过期' }
    ]
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const queryParams: LicenseQueryRequest = {
      page: currentPage.value,
      page_size: pageSize.value,
      customer_id: customerInfo.value.id || '',
      sort: 'created_at',
      order: 'desc'
    }

    // 添加状态筛选
    if (filterStatus.value) {
      queryParams.status = filterStatus.value as 'normal' | 'locked' | 'expired'
    }

    // 添加授权码搜索
    if (filterCode.value) {
      queryParams.code = filterCode.value
    }

    console.log('查询参数:', queryParams)
    const response = await getLicenses(queryParams)
    console.log('API响应:', response)
    tableData.value = response.data.list
    total.value = response.data.total
    console.log('表格数据:', tableData.value)
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error(t('pages.licenses.list.message.loadError'))
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(async () => {
  console.log('LinenseList组件已挂载')
  console.log('客户信息:', customerInfo.value)
  await loadStatusOptions()
  loadData()
})
</script>

<style lang="scss" scoped>
// Variables and mixins are auto-injected via Vite配置
@use 'sass:math';

.license-list-page {
  max-width: none !important;
  height: 100%;
  padding: 0;
  background-color: var(--app-content-bg);
  display: flex;
  flex-direction: column;
  overflow: hidden;
  box-sizing: border-box;
  width: 100%;
  margin: 0 auto;
  @include mobile {
    padding: $spacing-base;
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

.filters {
  @include flex-center-vertical;
  gap: $spacing-small;
  min-width: 288px; // 改为288px (36*8)
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
  min-width: 48px; // 改为48px (6*8)

  @include mobile {
    width: auto;
    margin-right: $spacing-small;
  }

  @include mobile {
    align-self: flex-start;
  }
}

.filter-input {
  width: 150px;
}

.filter-select {
  width: 120px; // 符合8px栅格(15*8)
}

.table-container {
  width: calc(100% - #{$spacing-large * 2});
  margin: 0 $spacing-large;
  border-radius: $border-radius-base;
  overflow: hidden;
  margin-bottom: $spacing-medium;
  flex: 1;
  display: flex;
  flex-direction: column;
  @include mobile {
    @include table-responsive;
  }
}

.table-wrapper {
  width: 100%;
  height: 100%;
  position: relative;
  overflow: auto;
  @include smooth-scroll;
  @include mobile {
    -webkit-overflow-scrolling: touch;
  }
}

.pagination-container {
  @include flex-end;
  gap: $spacing-small;
  padding: 0 $spacing-large $spacing-large $spacing-large;
  flex-shrink: 0;
  @include mobile {
    @include flex-center;
    flex-wrap: wrap;
    gap: $spacing-extra-small;
  }
}

:deep(.el-table) {
  width: 100% !important;
  min-width: 1408px !important;
  border: 1px solid var(--app-border-light);
  table-layout: fixed !important;
  .el-table__header-wrapper,
  .el-table__body-wrapper {
    width: 100% !important;
    overflow-x: auto !important;
  }
  .el-table__body-wrapper {
    overflow-y: auto !important;
    flex: 1 !important;
    margin-right: 0 !important;
    padding-right: 0 !important;
  }
  .el-table__header,
  .el-table__body {
    width: 100% !important;
    table-layout: fixed !important;
  }
  th, td {
    box-sizing: border-box !important;
    overflow: hidden !important;
    text-overflow: ellipsis !important;
    white-space: nowrap !important;
    padding: 0 !important;
    margin: 0 !important;
    height: 48px !important;
    font-family: 'Source Han Sans CN', sans-serif !important;
    font-size: 14px !important;
    vertical-align: middle !important;
  }
  colgroup {
    width: 100% !important;
  }
  .el-table__header-wrapper th {
    background-color: var(--app-bg-color) !important;
    font-weight: 500 !important;
    font-size: 16px !important;
    line-height: 24px !important;
    color: var(--app-text-primary) !important;
    border-bottom: 1px solid var(--app-border-light) !important;
  }
  .el-table__body-wrapper td {
    font-weight: 400 !important;
    line-height: 22px !important;
    color: var(--app-text-primary) !important;
    border-bottom: 1px solid var(--app-border-light) !important;
    white-space: nowrap !important;
    overflow: hidden !important;
    text-overflow: ellipsis !important;
  }
  .cell {
    white-space: nowrap !important;
    overflow: hidden !important;
    text-overflow: ellipsis !important;
    padding: 0 !important;
  }
}

.status-tag {
  padding: $spacing-extra-small $spacing-small;
  margin: 0 10px;
  border-radius: $border-radius-small;
  font-size: $font-size-base;
  font-weight: $font-weight-secondary;
  height: 24px;
  min-width: fit-content;
  @include flex-center;
  @include text-ellipsis;
  &.status-normal {
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
  }
  &.status-locked {
    background: rgba(230, 162, 60, 0.08);
    color: #e6a23c;
  }
  &.status-expired {
    background: rgba(244, 108, 108, 0.08);
    color: #f46c6c;
  }
}

.license-code {
  color: var(--el-color-primary);
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 400;
  font-size: 14px;
  line-height: 24px; // 改为24px符合8px栅格
}

.progress-container {
  width: 100%;
  padding: 0 8px;
  display: flex;
  align-items: center;
  justify-content: center;

  :deep(.el-progress) {
    width: 100%;
    min-width: 160px;
  }
  
  :deep(.el-progress-bar__outer) {
    background: #4763FF1A !important;
    min-height: 24px;
    position: relative;
  }
  
  :deep(.el-progress-bar__inner) {
    background: #4763FF !important;
    min-height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
  }
  
  :deep(.el-progress-bar__innerText) {
    color: #FFFFFF !important;
    font-size: 12px;
    font-weight: 500;
    white-space: nowrap;
    z-index: 10;
    pointer-events: none;
  }
  
  :deep(.el-progress__text) {
    color: #FFFFFF !important;
    font-size: 12px;
    font-weight: 500;
    z-index: 10;
  }

  // 当激活数量为 0 时，文字显示为红色
  &.zero-activated {
    :deep(.el-progress-bar__innerText),
    :deep(.el-progress__text) {
      color: #F56C6C !important;
      margin-left: 50px !important;
    }
    
    // 当百分比为0时，确保外部文字可见且为红色
    :deep(.el-progress-bar__outer) {
      position: relative;
    }
  }
}

// 操作按钮组，严格复制 customers 页 action-btn/action-buttons 样式
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
  height: 32px; // 改为32px符合8px栅格
  border-radius: $border-radius-small;
  font-size: $font-size-small;
  font-weight: $font-weight-primary;
  line-height: 20px; // 改为20px符合8px栅格
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
    color: #f46c6c;
    @include non-touch-device {
      &:hover {
        background: rgba(244, 108, 108, 0.15);
      }
    }
  }
  &.warning {
    background: rgba(230, 162, 60, 0.08);
    color: #e6a23c;
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

.pagination-container {
  @include flex-end;
  gap: $spacing-small;
  padding: 0 $spacing-large $spacing-large $spacing-large;
  flex-shrink: 0;
  @include mobile {
    @include flex-center;
    flex-wrap: wrap;
    gap: $spacing-extra-small;
  }
}
:deep(.el-pagination) {
  justify-content: flex-end;
  gap: 8px;
  --el-pagination-font-size: 14px;
  --el-pagination-bg-color: var(--app-content-bg);
  --el-pagination-text-color: var(--app-text-primary);
  --el-pagination-border-color: var(--app-border-color);
  --el-pagination-hover-color: #019c7c;
  font-family: 'Source Han Sans CN', sans-serif;
}
:deep(.el-pagination .btn-prev), :deep(.el-pagination .btn-next), :deep(.el-pagination .el-pager li) {
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
:deep(.el-pagination .el-pager li.is-active) {
  background: var(--el-color-primary) !important;
  border-color: var(--el-color-primary) !important;
  color: white !important;
}
:deep(.el-pagination .btn-prev:hover), :deep(.el-pagination .btn-next:hover), :deep(.el-pagination .el-pager li:hover) {
  color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}
:deep(.el-pagination .el-pagination__total), :deep(.el-pagination .el-pagination__sizes) {
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 350;
  font-size: 14px;
  color: var(--app-text-primary);
}
:deep(.el-pagination .el-pagination__sizes .el-select) {
  width: 88px;
}
:deep(.el-pagination .el-pagination__sizes .el-select .el-select__wrapper) {
  height: 32px;
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  background: var(--app-content-bg);
}
:deep(.el-pagination .el-pagination__jump) {
  font-family: 'Source Han Sans CN', sans-serif;
  font-weight: 350;
  font-size: 14px;
  color: var(--app-text-primary);
  margin-left: 8px;
}
:deep(.el-pagination .el-pagination__jump .el-input) {
  width: 48px;
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
:deep(.el-table__body tr:hover > td) {
  background: #e3f7f1 !important;
  transition: background 0.2s;
}
[data-theme='dark'] :deep(.el-table__body tr:hover > td) {
  background: rgba(16, 185, 129, 0.15) !important;
}
</style>
