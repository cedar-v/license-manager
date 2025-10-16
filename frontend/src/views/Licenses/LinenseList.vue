<template>
  <div class="license-list-page">
    <!-- 顶部操作区 -->
    <div class="top-actions">
      <el-button type="primary" @click="handleCreate">
        创建授权
      </el-button>

      <div class="filters">
        <span class="filter-label">筛选：</span>
        <el-select v-model="filterStatus" placeholder="状态" clearable class="filter-select">
          <el-option label="正常" value="normal" />
          <el-option label="锁定" value="locked" />
          <el-option label="过期" value="expired" />
        </el-select>

        <el-date-picker
          v-model="filterCreateDate"
          type="date"
          placeholder="创建时间"
          clearable
          class="filter-date"
        />

        <el-date-picker
          v-model="filterExpiryDate"
          type="date"
          placeholder="到期时间"
          clearable
          class="filter-date"
        />

        <el-button type="primary" @click="handleQuery">
          查询
        </el-button>
      </div>
    </div>

    <!-- 表格 -->
    <div class="table-container">
      <el-table
        :data="tableData"
        stripe
        style="width: 100%"
        v-loading="loading"
        empty-text="暂无数据"
      >
        <el-table-column prop="code" label="授权码" width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="current_activations" label="激活状态" width="120" />
        <el-table-column prop="end_date" label="到期时间" width="180" />
        <el-table-column prop="description" label="备注" min-width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="handleDetail(row)">
              详情
            </el-button>
            <el-button 
              v-if="row.is_locked" 
              type="success" 
              size="small" 
              @click="handleUnlock(row)"
            >
              解锁
            </el-button>
            <el-button 
              v-else 
              type="warning" 
              size="small" 
              @click="handleLock(row)"
            >
              锁定
            </el-button>
            <el-button type="danger" size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 分页 -->
    <div class="pagination-container">
      <el-pagination
        :current-page="currentPage"
        :page-size="pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
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
import { getLicenses, lockAuthorizationCode, type AuthorizationCode, type LicenseQueryRequest } from '@/api/license'

const router = useRouter()
const route = useRoute()

// 从路由参数获取客户信息
const customerInfo = computed(() => {
  return {
    id: route.query.customerId as string || '',
    name: route.query.customerName as string || ''
  }
})

// 响应式数据
const loading = ref(false)
const filterStatus = ref('')
const filterCreateDate = ref('')
const filterExpiryDate = ref('')
const tableData = ref<AuthorizationCode[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

// 方法
const getStatusTagType = (status: string) => {
  switch (status) {
    case 'normal':
      return 'success'
    case 'locked':
      return 'warning'
    case 'expired':
      return 'danger'
    default:
      return 'info'
  }
}

const getStatusText = (status: string) => {
  switch (status) {
    case 'normal':
      return '正常'
    case 'locked':
      return '锁定'
    case 'expired':
      return '过期'
    default:
      return status
  }
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

const handleDetail = (_row: any) => {
  ElMessage.info('查看详情功能待实现')
}

const handleLock = async (row: any) => {
  ElMessageBox.confirm('确定要锁定此授权吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await lockAuthorizationCode(row.id, { 
        is_locked: true,
        lock_reason: '手动锁定',
        reason: '管理员手动锁定'
      })
      ElMessage.success('锁定成功')
      loadData()
    } catch (error) {
      console.error('锁定失败:', error)
      ElMessage.error('锁定失败')
    }
  }).catch(() => {
    // 取消操作
  })
}

const handleUnlock = async (row: any) => {
  ElMessageBox.confirm('确定要解锁此授权吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      await lockAuthorizationCode(row.id, { 
        is_locked: false,
        lock_reason: '手动解锁',
        reason: '管理员手动解锁'
      })
      ElMessage.success('解锁成功')
      loadData()
    } catch (error) {
      console.error('解锁失败:', error)
      ElMessage.error('解锁失败')
    }
  }).catch(() => {
    // 取消操作
  })
}

const handleDelete = (_row: any) => {
  ElMessageBox.confirm('确定要删除此授权吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(() => {
    ElMessage.success('删除成功')
    loadData()
  }).catch(() => {
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

    console.log('查询参数:', queryParams)
    const response = await getLicenses(queryParams)
    console.log('API响应:', response)
    tableData.value = response.data.list
    total.value = response.data.total
    console.log('表格数据:', tableData.value)
  } catch (error) {
    console.error('加载数据失败:', error)
    ElMessage.error('加载数据失败')
  } finally {
    loading.value = false
  }
}

// 生命周期
onMounted(() => {
  console.log('LinenseList组件已挂载')
  console.log('客户信息:', customerInfo.value)
  loadData()
})
</script>

<style scoped lang="scss">
.license-list-page {
  padding: 20px;
  background: #fff;
  min-height: calc(100vh - 80px);
}

.customer-info {
  margin-bottom: 20px;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  border-left: 4px solid #409eff;

  h2 {
    margin: 0;
    font-size: 18px;
    font-weight: 600;
    color: #303133;
  }
}

.debug-info {
  margin-bottom: 20px;
  padding: 16px;
  background: #fff3cd;
  border-radius: 8px;
  border-left: 4px solid #ffc107;

  p {
    margin: 4px 0;
    font-size: 14px;
    color: #856404;
  }
}

.top-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  gap: 16px;
  flex-wrap: wrap;
}

.filters {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
}

.filter-label {
  font-size: 14px;
  color: #606266;
}

.filter-select {
  width: 120px;
}

.filter-date {
  width: 160px;
}

.table-container {
  margin-bottom: 20px;
}

.pagination-container {
  display: flex;
  justify-content: center;
}

// 响应式设计
@media (max-width: 768px) {
  .license-list-page {
    padding: 12px;
  }

  .top-actions {
    flex-direction: column;
    align-items: stretch;
  }

  .filters {
    width: 100%;

    .filter-select,
    .filter-date {
      flex: 1;
      min-width: 100px;
    }
  }
}
</style>
