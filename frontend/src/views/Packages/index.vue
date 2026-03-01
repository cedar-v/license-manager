<template>
  <Layout :page-title="t('packages.title')">
    <div class="packages-page">
      <div class="list-card" v-loading="loading">
        <div class="list-header">
          <div class="list-title">
            <span class="title-bar"></span>
            <span>{{ t('packages.table.title') }}</span>
          </div>
          <el-button 
            class="refresh-btn" 
            :icon="Refresh" 
            circle 
            @click="fetchData"
            :loading="loading"
          />
        </div>

        <el-table
          :data="tableData"
          stripe
          class="packages-table"
          :header-cell-style="{
            backgroundColor: '#E6F7F3',
            color: '#4F4F4F',
            fontWeight: '600',
            height: '50px'
          }"
        >
          <!-- <el-table-column prop="id" :label="t('packages.table.id')" width="100" /> -->
          <el-table-column prop="name" :label="t('packages.table.name')" min-width="150" />
          <el-table-column prop="price" :label="t('packages.table.price')" min-width="120">
            <template #default="{ row }">
              <span :class="{ 'price-free': (row.type === 'trial' && (row.price === 0 || row.price === '0')) || row.type === 'custom' }">
                <template v-if="row.type === 'custom'">
                  {{ t('packages.table.customPlan') }}
                </template>
                <template v-else-if="row.type === 'trial' && (row.price === 0 || row.price === '0')">
                  {{ t('packages.table.free') }}
                </template>
                <template v-else>
                  {{ typeof row.price === 'number' ? `${t('packages.table.priceSymbol')}${row.price}` : row.price }}
                </template>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="duration_description" :label="t('packages.table.cycle')" min-width="180" />
          <el-table-column prop="created_at" :label="t('packages.table.createdAt')" min-width="180" />
          <el-table-column prop="updated_at" :label="t('packages.table.updatedAt')" min-width="180" />
          <el-table-column :label="t('packages.table.status')" width="120">
            <template #default="{ row }">
              <span class="status-tag" :class="row.status === 1 ? 'enabled' : 'disabled'">
                {{ row.status === 1 ? t('packages.status.enabled') : t('packages.status.disabled') }}
              </span>
            </template>
          </el-table-column>
          <el-table-column :label="t('packages.table.actions')" width="100" fixed="right">
            <template #default="{ row }">
              <el-button size="small" class="btn-edit" @click="handleEdit(row)">
                {{ t('packages.actions.edit') }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <PackageEditDialog
      v-model="editVisible"
      :data="selectedPackage"
      @save="handleUpdate"
    />
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/common/layout/Layout.vue'
import PackageEditDialog from './components/PackageEditDialog.vue'
import { getPackages, updatePackage, getPackageDetail, type Package } from '@/api/package'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import { formatDateTime } from '@/utils/date'

const { t } = useI18n()

const tableData = ref<Package[]>([])
const loading = ref(false)

const fetchData = async () => {
  loading.value = true
  try {
    const res = await getPackages()
    if (res.code === '000000' && res.data) {
      tableData.value = res.data.packages.map((item: Package) => {
        let features = []
        try {
          if (typeof item.features === 'string' && item.features) {
            const parsed = JSON.parse(item.features)
            features = Array.isArray(parsed) ? parsed : [item.features]
          } else if (Array.isArray(item.features)) {
            features = item.features
          }
        } catch (e) {
          features = item.features ? [item.features] : []
        }

        return {
          ...item,
          features,
          created_at: formatDateTime(item.created_at),
          updated_at: formatDateTime(item.updated_at)
        }
      })
    }
  } catch (error: any) {
    console.error('Fetch packages error:', error)
    ElMessage.error(error.backendMessage || t('packages.messages.fetchError'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})

const editVisible = ref(false)
const selectedPackage = ref<any>(null)

const handleEdit = async (row: any) => {
  loading.value = true
  try {
    const res = await getPackageDetail(row.id)
    if (res.code === '000000' && res.data) {
      const detail = res.data
      let features = []
      try {
        if (typeof detail.features === 'string' && detail.features) {
          const parsed = JSON.parse(detail.features)
          features = Array.isArray(parsed) ? parsed : [detail.features]
        } else if (Array.isArray(detail.features)) {
          features = detail.features
        }
      } catch (e) {
        features = detail.features ? [detail.features] : []
      }

      selectedPackage.value = {
        ...detail,
        features
      }
      editVisible.value = true
    }
  } catch (error: any) {
    console.error('Fetch package detail error:', error)
    ElMessage.error(error.backendMessage || t('packages.messages.fetchDetailError'))
  } finally {
    loading.value = false
  }
}

const handleUpdate = async (updatedData: any) => {
  try {
    const { id, created_at, updated_at, ...data } = updatedData
    // 将 features 数组转回字符串供后端存储
    if (Array.isArray(data.features)) {
      data.features = JSON.stringify(data.features)
    }
    const res = await updatePackage(id, data)
    if (res.code === '000000') {
      ElMessage.success(t('packages.messages.updateSuccess'))
      fetchData()
    }
  } catch (error: any) {
    console.error('Update package error:', error)
    ElMessage.error(error.backendMessage || t('packages.messages.updateError'))
  }
}

</script>

<style lang="scss" scoped>
.packages-page {
  padding: 24px;
  background-color: #f0f2f5;
  min-height: calc(100vh - 80px);
}

.list-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
  padding: 24px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
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

.refresh-btn {
  border-color: #E6F7F3;
  background-color: #E6F7F3 !important;
  color: #00a870 !important;
  transition: all 0.3s;
  
  &:hover {
    border-color: #00a870;
    background-color: #00a870 !important;
    color: #fff !important;
    transform: rotate(180deg);
  }
}

.title-bar {
  width: 4px;
  height: 18px;
  background: #00a870;
  border-radius: 2px;
}

.status-tag {
  font-size: 14px;
  font-weight: 500;
  &.enabled {
    color: #409eff;
  }
  &.disabled {
    color: #999;
  }
}

.btn-edit {
  border: none;
  background: #e6f7f3 !important;
  color: #00a870 !important;
  padding: 4px 12px;
  height: 28px;
  font-size: 12px;
  border-radius: 4px;
}

.price-free {
  color: #666;
}
</style>
