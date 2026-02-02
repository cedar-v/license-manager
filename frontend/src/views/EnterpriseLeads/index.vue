<template>
  <Layout :page-title="t('enterpriseLeads.title')">
    <div class="enterprise-leads">
      <div class="stats-grid">
        <div v-for="stat in stats" :key="stat.key" class="stat-card">
          <div class="stat-info">
            <p class="stat-label">{{ t(`enterpriseLeads.stats.${stat.key}`) }}</p>
            <p class="stat-value">{{ stat.value.toLocaleString() }}</p>
          </div>
          <div class="stat-icon" :class="stat.key">
            <el-icon><component :is="stat.icon" /></el-icon>
          </div>
        </div>
      </div>

      <div class="filter-card">
        <div class="filter-row">
          <div class="search-group">
            <el-input
              v-model="filters.keyword"
              :placeholder="t('enterpriseLeads.filter.searchPlaceholder')"
              class="search-input"
              clearable
              @keyup.enter="handleFilter"
            >
              <template #append>
                <el-button class="search-btn" @click="handleFilter">
                  <el-icon><Search /></el-icon>
                </el-button>
              </template>
            </el-input>
          </div>
          <div class="status-group">
            <el-select v-model="filters.status" class="status-select" :placeholder="t('enterpriseLeads.filter.statusPlaceholder')" clearable @change="handleFilter">
              <el-option :label="t('enterpriseLeads.filter.allStatus')" value="" />
              <el-option :label="t('enterpriseLeads.status.pending')" value="pending" />
              <el-option :label="t('enterpriseLeads.status.contacting')" value="contacting" />
              <el-option :label="t('enterpriseLeads.status.completed')" value="completed" />
              <el-option :label="t('enterpriseLeads.status.rejected')" value="rejected" />
            </el-select>
          </div>
        </div>
      </div>

      <div class="list-card" v-loading="loading">
        <div class="list-header">
          <div class="list-title">
            <span class="title-bar"></span>
            <span>{{ t('enterpriseLeads.table.title') }}</span>
          </div>
          <el-button class="refresh-btn" :icon="Refresh" text @click="handleRefresh">
            {{ t('enterpriseLeads.actions.refresh') }}
          </el-button>
        </div>

        <el-table
          :data="tableData"
          stripe
          class="leads-table"
          :header-cell-style="{
            backgroundColor: '#E6F7F3',
            color: '#4F4F4F',
            fontWeight: '600'
          }"
        >
          <el-table-column prop="id" :label="t('enterpriseLeads.table.id')" />
          <el-table-column prop="company_name" :label="t('enterpriseLeads.table.company')" />
          <el-table-column prop="contact_name" :label="t('enterpriseLeads.table.contact')" width="100" />
          <el-table-column prop="contact_phone" :label="t('enterpriseLeads.table.phone')" />
          <el-table-column prop="created_at" :label="t('enterpriseLeads.table.submittedAt')" />
          <el-table-column :label="t('enterpriseLeads.table.status')" width="80">
            <template #default="{ row }">
              <span class="status-tag" :class="row.status">{{ t(`enterpriseLeads.status.${row.status}`) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('enterpriseLeads.table.actions')" width="200" fixed="right">
            <template #default="{ row }">
              <div class="action-buttons">
                <el-button size="small" class="btn-view" @click="handleView(row)">{{ t('enterpriseLeads.actions.view') }}</el-button>
                <el-button size="small" class="btn-edit" @click="handleEdit(row)">{{ t('enterpriseLeads.actions.edit') }}</el-button>
                <el-button size="small" class="btn-delete" @click="handleDelete(row)">{{ t('enterpriseLeads.actions.delete') }}</el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-container">
          <el-pagination
            v-model:current-page="page"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            layout="prev, pager, next, jumper, sizes, total"
            :total="total"
            background
          />
        </div>
      </div>
    </div>

    <LeadDetailDialog
      v-model="detailVisible"
      :id="selectedLead?.id"
    />

    <LeadEditDialog
      v-model="editVisible"
      :id="selectedLead?.id"
      @save="handleUpdate"
    />
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/common/layout/Layout.vue'
import { Phone, CircleCheck, Refresh, User, OfficeBuilding, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import LeadDetailDialog from './components/LeadDetailDialog.vue'
import LeadEditDialog from './components/LeadEditDialog.vue'
import { getLeads, getLeadSummary, updateLead, deleteLead, type Lead } from '@/api/lead'
import { formatDateTime } from '@/utils/date'

const { t } = useI18n()

const stats = ref([
  { key: 'total', value: 0, icon: OfficeBuilding, field: 'total_count' },
  { key: 'pending', value: 0, icon: Phone, field: 'pending_count' },
  { key: 'contacting', value: 0, icon: CircleCheck, field: 'contacted_count' },
  { key: 'completed', value: 0, icon: User, field: 'converted_count' }
])

const filters = ref({
  keyword: '',
  status: ''
})

const tableData = ref<Lead[]>([])
const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const detailVisible = ref(false)
const editVisible = ref(false)
const selectedLead = ref<any>(null)

const fetchSummary = async () => {
  try {
    const res = await getLeadSummary()
    if (res.code === '000000' && res.data) {
      stats.value.forEach(item => {
        if (res.data[item.field] !== undefined) {
          item.value = res.data[item.field]
        }
      })
    }
  } catch (error) {
    console.error('Fetch lead summary error:', error)
  }
}

const fetchData = async () => {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value,
      search: filters.value.keyword,
      status: filters.value.status
    }
    const res = await getLeads(params)
    if (res.code === '000000' && res.data) {
      tableData.value = res.data.leads.map((item: any) => ({
        ...item,
        created_at: formatDateTime(item.created_at)
      }))
      total.value = res.data.total_count
    }
  } catch (error: any) {
    console.error('Fetch leads error:', error)
    ElMessage.error(error.backendMessage || t('enterpriseLeads.messages.fetchError'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
  fetchSummary()
})

watch([page, pageSize], () => {
  fetchData()
})

const handleFilter = () => {
  page.value = 1
  fetchData()
}

const handleRefresh = () => {
  fetchData()
  fetchSummary()
}

const handleView = (row: any) => {
  selectedLead.value = row
  detailVisible.value = true
}

const handleEdit = (row: any) => {
  selectedLead.value = row
  editVisible.value = true
}

const handleDelete = (row: any) => {
  ElMessageBox.confirm(
    t('enterpriseLeads.messages.deleteConfirm'),
    t('common.confirm'),
    {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'warning',
    }
  ).then(async () => {
    try {
      const res = await deleteLead(row.id)
      if (res.code === '000000') {
        ElMessage.success(t('enterpriseLeads.messages.deleteSuccess'))
        fetchData()
        fetchSummary()
      }
    } catch (error: any) {
      console.error('Delete lead error:', error)
      ElMessage.error(error.backendMessage || t('enterpriseLeads.messages.deleteError'))
    }
  }).catch(() => {
    // Cancelled
  })
}

const handleUpdate = async (updatedData: any) => {
  try {
    const { id, ...data } = updatedData
    const res = await updateLead(id, data)
    if (res.code === '000000') {
      ElMessage.success(t('enterpriseLeads.messages.updateSuccess'))
      fetchData()
      fetchSummary()
    }
  } catch (error: any) {
    console.error('Update lead error:', error)
    ElMessage.error(error.backendMessage || t('enterpriseLeads.messages.updateError'))
  }
}
</script>

<style lang="scss" scoped>
.enterprise-leads {
  padding: 24px;
  background-color: #f0f2f5;
  min-height: calc(100vh - 80px);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24px;
  margin-bottom: 24px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px 24px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  border: none;
  height: 100px;
}

.stat-label {
  font-size: 14px;
  color: #999;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: bold;
  color: #00a870;
  line-height: 1.2;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
}

.stat-icon.total {
  background: #f0f4ff;
  color: #409eff;
}

.stat-icon.pending {
  background: #fff8e6;
  color: #e6a23c;
}

.stat-icon.contacting {
  background: #e6f7f3;
  color: #00a870;
}

.stat-icon.completed {
  background: #f5f0ff;
  color: #7c5cfc;
}

.filter-card {
  background: #fff;
  padding: 24px;
  border-radius: 8px;
  margin-bottom: 24px;
  display: flex;
  align-items: center;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  border: none;
}

.filter-row {
  width: 100%;
  display: flex;
  gap: 20px;
  align-items: center;
}

.search-group {
  width: 320px;
}

.search-input :deep(.el-input__wrapper) {
  border-radius: 4px 0 0 4px;
  box-shadow: 1px 0 0 0 var(--el-input-border-color) inset, 0 1px 0 0 var(--el-input-border-color) inset, 0 -1px 0 0 var(--el-input-border-color) inset !important;
}

.search-btn {
  background: #00a870 !important;
  border-color: #00a870 !important;
  color: #fff !important;
  border-radius: 0 4px 4px 0;
  width: 44px;
  height: 32px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.status-select {
  width: 160px;
}

.status-select :deep(.el-input__wrapper) {
  border-radius: 4px;
}

.list-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  padding: 24px;
  border: none;
}

.list-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
}

.list-title {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.title-bar {
  width: 4px;
  height: 18px;
  background: #00a870;
  border-radius: 2px;
}

.refresh-btn {
  color: #666;
  font-size: 14px;
}

.leads-table :deep(.el-table__header th) {
  font-weight: 600;
  height: 50px;
}

.status-tag {
  font-size: 14px;
  font-weight: 500;
}

.status-tag.contacting {
  color: #409eff;
}

.status-tag.pending {
  color: #e6a23c;
}

.status-tag.completed {
  color: #333;
}

.status-tag.rejected {
  color: #999;
}

.action-buttons {
  display: flex;
  gap: 8px;
}

.btn-view, .btn-edit, .btn-delete {
  border: none;
  padding: 4px 12px;
  height: 28px;
  font-size: 12px;
  border-radius: 4px;
}

.btn-view, .btn-edit {
  background: #e6f7f3 !important;
  color: #00a870 !important;
}

.btn-delete {
  background: #fff1f0 !important;
  color: #f5222d !important;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 24px;
}

.pagination-container :deep(.el-pagination.is-background .el-pager li:not(.is-disabled).is-active) {
  background-color: #00a870;
}

@media (max-width: 1400px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .enterprise-leads {
    padding: 16px;
  }

  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }

  .status-select {
    width: 100%;
  }
}
</style>
