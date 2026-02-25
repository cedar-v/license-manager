<template>
  <Layout :page-title="t('navigation.menu.invoices')">
    <div class="invoice-container">
      <!-- 统计卡片 -->
      <div class="stats-row">
        <div v-for="stat in stats" :key="stat.key" class="stat-card" :class="stat.type">
          <div class="stat-info">
            <div class="stat-label">{{ t(`invoices.stats.${stat.type}`) }}</div>
            <div class="stat-value">{{ stat.value }}</div>
          </div>
          <div class="stat-icon-wrapper">
            <el-icon v-if="stat.icon === 'document'">
              <Document />
            </el-icon>
            <el-icon v-else-if="stat.icon === 'warning'">
              <Warning />
            </el-icon>
            <el-icon v-else-if="stat.icon === 'success'">
              <CircleCheck />
            </el-icon>
            <el-icon v-else-if="stat.icon === 'error'">
              <CircleClose />
            </el-icon>
          </div>
        </div>
      </div>

      <!-- 搜索和筛选 -->
      <div class="filter-section">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <div class="form_box">
            <el-form-item :label="t('invoices.filter.searchLabel')">
              <el-input
                v-model="filterForm.keyword"
                :placeholder="t('invoices.filter.searchPlaceholder')"
              />
            </el-form-item>
            <el-form-item :label="t('invoices.filter.statusLabel')">
              <el-select
                v-model="filterForm.status"
                :placeholder="t('invoices.filter.allStatus')"
                style="width: 150px"
              >
                <el-option :label="t('invoices.filter.allStatus')" value="" />
                <el-option :label="t('invoices.filter.pending')" value="pending" />
                <el-option :label="t('invoices.filter.completed')" value="issued" />
                <el-option :label="t('invoices.filter.rejected')" value="rejected" />
              </el-select>
            </el-form-item>
            <el-form-item :label="t('invoices.filter.timeLabel')">
              <el-date-picker
                v-model="filterForm.dateRange"
                type="daterange"
                range-separator="-"
                :start-placeholder="t('chart.licenseTrend.datePicker.startPlaceholder')"
                :end-placeholder="t('chart.licenseTrend.datePicker.endPlaceholder')"
                format="YYYY/MM/DD"
              />
            </el-form-item>
            <div class="filter-right">
              <el-form-item>
                <el-button @click="resetFilter">{{ t('invoices.filter.reset') }}</el-button>
                <el-button type="primary" class="btn-filter" @click="handleFilter">{{
                  t('invoices.filter.query')
                }}</el-button>
              </el-form-item>
            </div>
          </div>
        </el-form>
      </div>

      <!-- 数据表格 -->
      <div class="table-section" v-loading="loading">
        <div class="table-header-title">
          <span class="title-line"></span>
          <h3>{{ t('invoices.table.listTitle') }}</h3>
          <!-- <el-button
            type="primary"
            class="btn-create"
            @click="handleCreate"
            style="margin-left: auto"
            >{{ t('invoices.actions.create') }}</el-button
          > -->
        </div>
        <el-table
          :data="tableData"
          style="width: 100%"
          :header-cell-style="{
            backgroundColor: 'var(--app-bg-color)',
            color: 'var(--app-text-primary)',
            borderRight: '1px solid var(--app-border-light)'
          }"
          :row-style="getRowStyle"
        >
          <el-table-column
            prop="invoice_no"
            :label="t('invoices.table.invoiceNo')"
            min-width="150"
          />
          <el-table-column :label="t('invoices.table.user')" min-width="120">
            <template #default="{ row }">
              <div class="user-info">
                <span class="user-name">{{ row.cu_user_phone }}</span>
                <!-- <span class="user-type" :class="row.order_package_name">{{ row.order_package_name }}</span> -->
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="order_no" :label="t('invoices.table.orderNo')" min-width="150" />
          <el-table-column
            prop="created_at"
            :label="t('invoices.table.applyTime')"
            min-width="180"
          />
          <el-table-column :label="t('invoices.table.status')" width="100">
            <template #default="{ row }">
              <span class="status-text" :class="row.status">{{ row.status_display }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="amount" :label="t('invoices.table.amount')" width="120">
            <template #default="{ row }">
              ¥{{ (row.amount || 0).toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
            </template>
          </el-table-column>
          <el-table-column :label="t('invoices.table.operation')" width="200" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" style="color: #019c7c" @click="handleView(row)">{{
                t('invoices.actions.view')
              }}</el-button>
              <template v-if="row.status === 'pending'">
                <el-button link type="success" style="color: #52c41a" @click="handleUpload(row)">{{
                  t('invoices.actions.upload')
                }}</el-button>
                <el-button link type="danger" style="color: #f5222d" @click="handleReject(row)">{{
                  t('invoices.actions.reject')
                }}</el-button>
              </template>
            </template>
          </el-table-column>
        </el-table>

        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="prev, pager, next, jumper, sizes, total"
            :total="totalCount"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          >
          </el-pagination>
        </div>
      </div>

      <!-- 上传发票对话框 -->
      <UploadInvoiceDialog
        v-model="uploadVisible"
        :invoice-data="currentRow"
        @submit="handleUploadSubmit"
      />

      <!-- 驳回申请对话框 -->
      <RejectInvoiceDialog
        v-model="rejectVisible"
        :invoice-data="currentRow"
        @submit="handleRejectSubmit"
      />

      <!-- 新增申请对话框 -->
      <el-dialog v-model="createVisible" title="新增发票申请" width="500px">
        <el-form :model="createForm" label-width="100px">
          <el-form-item label="订单号">
            <el-input v-model="createForm.orderNo" placeholder="请输入订单号" />
          </el-form-item>
          <el-form-item label="开票金额">
            <el-input-number
              v-model="createForm.amount"
              :precision="2"
              :step="100"
              style="width: 100%"
            />
          </el-form-item>
          <el-form-item label="抬头类型">
            <el-radio-group v-model="createForm.userType">
              <el-radio label="personal">个人</el-radio>
              <el-radio label="enterprise">企业</el-radio>
            </el-radio-group>
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="createVisible = false">取消</el-button>
          <el-button type="primary" @click="submitCreate">提交</el-button>
        </template>
      </el-dialog>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/common/layout/Layout.vue'
import { Document, Warning, CircleCheck, CircleClose } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import UploadInvoiceDialog from './components/UploadInvoiceDialog.vue'
import RejectInvoiceDialog from './components/RejectInvoiceDialog.vue'
import {
  getInvoices,
  getInvoiceSummary,
  uploadInvoice,
  rejectInvoice,
  issueInvoice,
  getInvoiceDetail,
  type Invoice
} from '@/api/invoice'
import { formatDateTime } from '@/utils/date'

const { t } = useI18n()
const router = useRouter()

const stats = ref([
  { label: '全部发票申请', value: '0', type: 'all', icon: 'document', key: 'total_count' },
  { label: '待处理申请', value: '0', type: 'pending', icon: 'warning', key: 'pending_count' },
  { label: '已开票', value: '0', type: 'completed', icon: 'success', key: 'issued_count' },
  { label: '已驳回', value: '0', type: 'rejected', icon: 'error', key: 'rejected_count' }
])

const filterForm = reactive({
  keyword: '',
  status: '',
  dateRange: [] as any[]
})

const tableData = ref<Invoice[]>([])
const loading = ref(false)
const totalCount = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const createVisible = ref(false)
const uploadVisible = ref(false)
const rejectVisible = ref(false)
const currentRow = ref<any>(null)

const getRowStyle = ({ rowIndex }: { rowIndex: number }) => {
  return rowIndex % 2 === 1
    ? { backgroundColor: 'var(--app-bg-color)' }
    : { backgroundColor: 'var(--app-content-bg)' }
}

const fetchSummary = async () => {
  try {
    const res = await getInvoiceSummary()
    if (res.code === '000000' && res.data) {
      stats.value.forEach((item: any) => {
        const value = (res.data as any)[item.key]
        if (value !== undefined) {
          item.value = String(value)
        }
      })
    }
  } catch (error) {
    console.error('Fetch summary error:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
      keyword: filterForm.keyword,
      status: filterForm.status
    }
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.start_date = filterForm.dateRange[0]
      params.end_date = filterForm.dateRange[1]
    }
    const res = await getInvoices(params)
    if (res.code === '000000') {
      fetchSummary()
      res.data.invoices.forEach((item: any) => {
        item.created_at = formatDateTime(item.created_at)
      })
      tableData.value = res.data.invoices
      totalCount.value = res.data.total_count
    }
  } catch (error: any) {
    console.error('Fetch invoices error:', error)
    ElMessage.error(error.backendMessage || t('invoices.messages.fetchError'))
  } finally {
    loading.value = false
  }
}

// 获取详情
const fetchDetail = async (id: string) => {
  if (!id) return
  loading.value = true
  try {
    const res = await getInvoiceDetail(id)
    if (res.code === '000000' && res.data) {
      const data = res.data
      data.created_at = formatDateTime(data.created_at)
      if (data.rejected_at) data.rejected_at = formatDateTime(data.rejected_at)
      if (data.uploaded_at) data.uploaded_at = formatDateTime(data.uploaded_at)
      currentRow.value = data
    }
  } catch (error: any) {
    console.error('Fetch invoice detail error:', error)
    ElMessage.error(error.backendMessage || t('invoices.messages.fetchDetailError'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const createForm = reactive({
  orderNo: '',
  amount: 0,
  userType: 'personal'
})

const submitCreate = () => {
  ElMessage.success(t('invoices.messages.createSuccess'))
  createVisible.value = false
}

const handleFilter = () => {
  currentPage.value = 1
  fetchData()
}

const resetFilter = () => {
  filterForm.keyword = ''
  filterForm.status = ''
  filterForm.dateRange = []
  currentPage.value = 1
  fetchData()
}

const handleView = (row: any) => {
  router.push(`/invoices/detail/${row.id || row.invoice_no}`)
}

const handleUpload = (row: any) => {
  // currentRow.value = row
  fetchDetail(row.id)
  setTimeout(() => {
    uploadVisible.value = true
  }, 300)
}

const handleUploadSubmit = async (data: any) => {
  console.log('handleUploadSubmit triggered:', data)
  if (!currentRow.value || !currentRow.value.id) {
    console.error('No current row or id')
    return
  }
  loading.value = true
  try {
    const formData = new FormData()
    formData.append('invoice_no', currentRow.value.id)
    if (data.fileList && data.fileList.length > 0) {
      const file = data.fileList[0].raw
      console.log('File to upload:', file)
      formData.append('file', file)
    } else {
      ElMessage.warning(t('invoices.dialog.uploadTip'))
      loading.value = false
      return
    }

    // 1. 上传文件
    const uploadRes = await uploadInvoice(formData)
    console.log('Upload response:', uploadRes)

    if (uploadRes.code === '000000' && (uploadRes.data?.file_url || uploadRes.data?.url)) {
      const fileUrl = uploadRes.data.file_url || uploadRes.data.url
      // 2. 调用开票接口
      const issueRes = await issueInvoice(currentRow.value.id, {
        invoice_file_url: fileUrl,
        issued_at: data.issued_at
          ? new Date(data.issued_at).toISOString().replace(/\.\d{3}/, '')
          : ''
      })

      if (issueRes.code === '000000') {
        ElMessage.success(t('invoices.messages.uploadSuccess'))
        uploadVisible.value = false
        fetchData()
        fetchSummary()
      }
    } else {
      throw new Error(uploadRes.message || 'Upload failed')
    }
  } catch (error: any) {
    console.error('Upload invoice error:', error)
    ElMessage.error(error.backendMessage || error.message || t('invoices.messages.fetchError'))
  } finally {
    loading.value = false
  }
}

const handleReject = (row: any) => {
  fetchDetail(row.id)
  // currentRow.value = row
  setTimeout(() => {
    rejectVisible.value = true
  }, 300)
}

const handleRejectSubmit = async (data: any) => {
  if (!currentRow.value?.id) return
  try {
    const reasonText =
      data.reason === 'other' ? data.suggestion : t(`invoices.dialog.reasons.${data.reason}`)
    const res = await rejectInvoice(currentRow.value.id, {
      reject_reason: reasonText,
      suggestion: data.suggestion
    })
    if (res.code === '000000') {
      ElMessage.success(t('invoices.messages.rejectSuccess'))
      fetchData()
      fetchSummary()
    }
  } catch (error: any) {
    console.error('Reject error:', error)
    ElMessage.error(error.backendMessage || t('invoices.messages.fetchError'))
  }
}

const handleSizeChange = (val: number) => {
  pageSize.value = val
  fetchData()
}

const handleCurrentChange = (val: number) => {
  currentPage.value = val
  fetchData()
}
</script>

<style lang="scss" scoped>
.invoice-container {
  padding: 24px;
  background-color: var(--app-bg-color);
  min-height: calc(100vh - 80px);
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  margin-bottom: 24px;
}

.stat-card {
  background: var(--app-content-bg);
  border-radius: 4px;
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: var(--app-shadow);
  border: 1px solid var(--app-border-light);
  height: 96px;

  .stat-label {
    font-size: 14px;
    color: var(--app-text-secondary);
    margin-bottom: 8px;
  }

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    color: var(--el-color-primary);
    line-height: 1.2;
  }

  .stat-icon-wrapper {
    width: 48px;
    height: 48px;
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
  }

  &.all .stat-icon-wrapper {
    background: var(--el-color-primary-light-9);
    color: var(--el-color-primary);
  }

  &.pending .stat-icon-wrapper {
    background: var(--el-color-warning-light-9);
    color: var(--el-color-warning);
  }

  &.completed .stat-icon-wrapper {
    background: var(--el-color-success-light-9);
    color: var(--el-color-success);
  }

  &.rejected .stat-icon-wrapper {
    background: var(--el-color-danger-light-9);
    color: var(--el-color-danger);
  }
}

.filter-section {
  background: var(--app-content-bg);
  padding: 16px 24px;
  border-radius: 4px;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  gap: 10px;
  border: 1px solid var(--app-border-light);
  .form_box{
    width: 100%;
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    gap: 10px;
  }

  .filter-form {
    width: 100%;
    display: flex;
    align-items: center;
    flex-wrap: wrap;

    :deep(.el-form-item) {
      margin-bottom: 0;
      margin-right: 24px;

      .el-form-item__label {
        font-weight: 500;
        color: var(--app-text-primary);
      }
    }
  }

  .filter-right {
    display: flex;
    align-items: center;
    margin-left: auto;

    :deep(.el-form-item) {
      margin-right: 12px;

      &:last-child {
        margin-right: 0;
      }
    }
  }
}

.selected-count {
  margin-right: 24px;

  .count-badge {
    background: var(--el-color-danger);
    color: #fff;
    padding: 0 12px;
    border-radius: 4px;
    font-size: 20px;
    font-weight: bold;
    height: 32px;
    line-height: 32px;
    display: inline-block;
  }
}

.btn-filter {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
  padding: 8px 24px;

  &:hover {
    background-color: var(--el-color-primary-dark-2);
    border-color: var(--el-color-primary-dark-2);
  }
}

.table-section {
  background: var(--app-content-bg);
  padding: 24px;
  border-radius: 4px;
  border: 1px solid var(--app-border-light);
  box-shadow: var(--app-shadow);
}

.table-header-title {
  display: flex;
  align-items: center;
  margin-bottom: 20px;

  .title-line {
    width: 4px;
    height: 16px;
    background: #019c7c;
    margin-right: 12px;
  }

  h3 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #262626;
  }
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;

  .user-name {
    color: #262626;
  }

  .user-type {
    font-size: 12px;
    padding: 0 8px;
    border-radius: 2px;
    height: 20px;
    line-height: 20px;

    &.personal {
      background: #e6f7ff;
      color: #1890ff;
    }

    &.enterprise {
      background: #f6ffed;
      color: #52c41a;
    }
  }
}

.status-text {
  &.completed {
    color: #1890ff;
  }

  &.pending {
    color: #faad14;
  }

  &.rejected {
    color: #bfbfbf;
  }
}

.pagination-container {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;

  :deep(.el-pagination) {
    .el-pager li.is-active {
      background-color: #019c7c;
      color: #fff;
    }
  }
}

:deep(.el-table) {
  --el-table-header-bg-color: #e6f7ff;
  border: 1px solid #bae7ff;

  .el-table__header th {
    color: #262626;
    height: 54px;
    border-right: 1px solid #bae7ff;

    &:last-child {
      border-right: none;
    }
  }

  .el-table__body td {
    height: 64px;
    color: #595959;
  }

  .el-table__row--striped {
    background-color: #fafafa;
  }
}

@media (max-width: 1400px) {
  .stats-row {
    gap: 16px;
  }

  .stat-card {
    padding: 16px;

    .stat-value {
      font-size: 24px;
    }
  }
}

@media (max-width: 1200px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
