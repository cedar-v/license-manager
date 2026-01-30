<template>
  <Layout :page-title="t('navigation.menu.invoices')">
    <div class="invoice-container">
      <!-- 统计卡片 -->
      <div class="stats-row">
        <div v-for="stat in stats" :key="stat.label" class="stat-card" :class="stat.type">
          <div class="stat-info">
            <div class="stat-label">{{ stat.label }}</div>
            <div class="stat-value">{{ stat.value }}</div>
          </div>
          <div class="stat-icon-wrapper">
            <el-icon v-if="stat.icon === 'document'"><Document /></el-icon>
            <el-icon v-else-if="stat.icon === 'warning'"><Warning /></el-icon>
            <el-icon v-else-if="stat.icon === 'success'"><CircleCheck /></el-icon>
            <el-icon v-else-if="stat.icon === 'error'"><CircleClose /></el-icon>
          </div>
        </div>
      </div>

      <!-- 搜索和筛选 -->
      <div class="filter-section">
        <el-form :inline="true" :model="filterForm" class="filter-form">
          <el-form-item label="搜索发票号/订单号:">
            <el-input v-model="filterForm.keyword" placeholder="输入发票号或订单号" />
          </el-form-item>
          <el-form-item label="开票状态:">
            <el-select v-model="filterForm.status" placeholder="全部状态" style="width: 150px">
              <el-option label="全部状态" value="" />
              <el-option label="待处理" value="pending" />
              <el-option label="已开票" value="success" />
              <el-option label="已驳回" value="rejected" />
            </el-select>
          </el-form-item>
          <div class="filter-right">
            <div class="selected-count">
              <span class="count-badge">{{ totalCount }}</span>
            </div>
            <el-form-item label="申请时间:">

              <el-date-picker
                v-model="filterForm.dateRange"
                type="daterange"
                range-separator="-"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                format="YYYY年MM月DD日"
              />
            </el-form-item>
            <el-form-item>
              <el-button @click="resetFilter">重置</el-button>
              <el-button type="primary" class="btn-filter" @click="handleFilter">筛选</el-button>
            </el-form-item>
          </div>
        </el-form>
      </div>

      <!-- 数据表格 -->
      <div class="table-section" v-loading="loading">
        <div class="table-header-title">
          <span class="title-line"></span>
          <h3>发票申请列表</h3>
          <el-button type="primary" class="btn-create" @click="handleCreate" style="margin-left: auto">新增发票申请</el-button>
        </div>
        <el-table :data="tableData" stripe style="width: 100%" :header-cell-style="{ background: '#E6F7FF', color: '#000', borderRight: '1px solid #BAE7FF' }">
          <el-table-column prop="invoiceNo" label="发票申请号" min-width="150" />
          <el-table-column label="用户" min-width="120">
            <template #default="{ row }">
              <div class="user-info">
                <span class="user-name">{{ row.user }}</span>
                <span class="user-type" :class="row.userType">{{ row.userTypeLabel }}</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="orderNo" label="订单号" min-width="150" />
          <el-table-column prop="time" label="申请时间" min-width="180" />
          <el-table-column label="开票状态" width="100">
            <template #default="{ row }">
              <span class="status-text" :class="row.status">{{ row.statusLabel }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="发票金额" width="120">
            <template #default="{ row }">
              ¥{{ (row.amount || 0).toLocaleString(undefined, { minimumFractionDigits: 2 }) }}
            </template>
          </el-table-column>
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" style="color: #019C7C" @click="handleView(row)">查看</el-button>
              <template v-if="row.status === 'pending'">
                <el-button link type="success" style="color: #52C41A" @click="handleUpload(row)">上传</el-button>
                <el-button link type="danger" style="color: #F5222D" @click="handleReject(row)">驳回</el-button>
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
            <el-input-number v-model="createForm.amount" :precision="2" :step="100" style="width: 100%" />
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
import { getInvoices, getInvoiceSummary, type Invoice } from '@/api/invoice'

const { t } = useI18n()
const router = useRouter()

const stats = ref([
  { label: '全部发票申请', value: '0', type: 'all', icon: 'document', key: 'total' },
  { label: '待处理申请', value: '0', type: 'pending', icon: 'warning', key: 'pending' },
  { label: '已开票', value: '0', type: 'success', icon: 'success', key: 'completed' },
  { label: '已驳回', value: '0', type: 'rejected', icon: 'error', key: 'rejected' }
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

const fetchSummary = async () => {
  try {
    const res = await getInvoiceSummary()
    if (res.code === '000000' && res.data) {
      stats.value.forEach(stat => {
        if (stat.key && (res.data as any)[stat.key] !== undefined) {
          stat.value = (res.data as any)[stat.key].toString()
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
      status: filterForm.status,
    }
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.start_date = filterForm.dateRange[0]
      params.end_date = filterForm.dateRange[1]
    }
    const res = await getInvoices(params)
    if (res.code === '000000') {
      tableData.value = res.data.list
      totalCount.value = res.data.total
    }
  } catch (error: any) {
    console.error('Fetch invoices error:', error)
    ElMessage.error(error.backendMessage || '获取数据失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
  fetchSummary()
})


const createForm = reactive({
  orderNo: '',
  amount: 0,
  userType: 'personal'
})

const handleCreate = () => {
  createVisible.value = true
}

const submitCreate = () => {
  ElMessage.success('提交申请成功')
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
  router.push(`/invoices/detail/${row.id || row.invoiceNo}`)
}

const handleUpload = (row: any) => {
  currentRow.value = row
  uploadVisible.value = true
}

const handleUploadSubmit = (data: any) => {
  console.log('Upload Submit:', data)
  fetchData()
  fetchSummary()
}

const handleReject = (row: any) => {
  currentRow.value = row
  rejectVisible.value = true
}

const handleRejectSubmit = (data: any) => {
  console.log('Reject Submit:', data)
  fetchData()
  fetchSummary()
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
  background-color: #F0F2F5;
  min-height: calc(100vh - 80px);
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  margin-bottom: 24px;
}

.stat-card {
  background: #fff;
  border-radius: 4px;
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
  border: 1px solid transparent;
  height: 96px;

  &.all { border: 2px solid #1890FF; }

  .stat-label {
    font-size: 14px;
    color: #8C8C8C;
    margin-bottom: 8px;
  }

  .stat-value {
    font-size: 28px;
    font-weight: 600;
    color: #019C7C;
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

  &.all .stat-icon-wrapper { background: #E6F7FF; color: #1890FF; }
  &.pending .stat-icon-wrapper { background: #FFF7E6; color: #FAAD14; }
  &.success .stat-icon-wrapper { background: #F6FFED; color: #52C41A; }
  &.rejected .stat-icon-wrapper { background: #F9F0FF; color: #722ED1; }
}

.filter-section {
  background: #fff;
  padding: 16px 24px;
  border-radius: 4px;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  border: 1px solid #F0F0F0;

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
        color: #262626;
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
    background: #FF4D4F;
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
  background-color: #019C7C;
  border-color: #019C7C;
  padding: 8px 24px;
  &:hover {
    background-color: #017c63;
    border-color: #017c63;
  }
}

.table-section {
  background: #fff;
  padding: 24px;
  border-radius: 4px;
  border: 1px solid #F0F0F0;
}

.table-header-title {
  display: flex;
  align-items: center;
  margin-bottom: 20px;

  .title-line {
    width: 4px;
    height: 16px;
    background: #019C7C;
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
    &.personal { background: #E6F7FF; color: #1890FF; }
    &.enterprise { background: #F6FFED; color: #52C41A; }
  }
}

.status-text {
  &.success { color: #1890FF; }
  &.pending { color: #FAAD14; }
  &.rejected { color: #BFBFBF; }
}

.pagination-container {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
  
  :deep(.el-pagination) {
    .el-pager li.is-active {
      background-color: #019C7C;
      color: #fff;
    }
  }
}

:deep(.el-table) {
  --el-table-header-bg-color: #E6F7FF;
  border: 1px solid #BAE7FF;
  
  .el-table__header th {
    color: #262626;
    height: 54px;
    border-right: 1px solid #BAE7FF;
    &:last-child {
      border-right: none;
    }
  }

  .el-table__body td {
    height: 64px;
    color: #595959;
  }

  .el-table__row--striped {
    background-color: #FAFAFA;
  }
}

@media (max-width: 1400px) {
  .stats-row {
    gap: 16px;
  }
  .stat-card {
    padding: 16px;
    .stat-value { font-size: 24px; }
  }
}

@media (max-width: 1200px) {
  .stats-row {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
