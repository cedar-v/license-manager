<template>
  <Layout :page-title="t('packages.title')">
    <div class="packages-page">
      <div class="list-card">
        <div class="list-header">
          <div class="list-title">
            <span class="title-bar"></span>
            <span>{{ t('packages.table.title') }}</span>
          </div>
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
          <el-table-column prop="id" :label="t('packages.table.id')" width="100" />
          <el-table-column prop="name" :label="t('packages.table.name')" min-width="150" />
          <el-table-column prop="price" :label="t('packages.table.price')" min-width="120">
            <template #default="{ row }">
              <span :class="{ 'price-free': row.price === '免费' }">{{ row.price }}</span>
            </template>
          </el-table-column>
          <el-table-column prop="cycle" :label="t('packages.table.cycle')" min-width="180" />
          <el-table-column prop="createdAt" :label="t('packages.table.createdAt')" min-width="180" />
          <el-table-column prop="updatedAt" :label="t('packages.table.updatedAt')" min-width="180" />
          <el-table-column :label="t('packages.table.status')" width="120">
            <template #default="{ row }">
              <span class="status-tag" :class="row.status">{{ t(`packages.status.${row.status}`) }}</span>
            </template>
          </el-table-column>
          <el-table-column :label="t('packages.table.actions')" width="100" fixed="right">
            <template #default="{ row }">
              <el-button size="small" class="btn-edit" @click="handleEdit(row)">
                {{ t('enterpriseLeads.actions.edit') }}
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
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/common/layout/Layout.vue'
import PackageEditDialog from './components/PackageEditDialog.vue'

const { t } = useI18n()

const tableData = ref([
  {
    id: '1',
    name: '试用版',
    price: '免费',
    cycle: '有效期至本月25日',
    createdAt: '2024-01-15 10:30',
    updatedAt: '2024-01-15 10:30',
    status: 'enabled',
    type: 'trial',
    sort: 1
  },
  {
    id: '2',
    name: '基础版',
    price: '¥300',
    cycle: '永久授权 (V1/V2/V3)',
    createdAt: '2024-01-15 10:30',
    updatedAt: '2024-01-15 10:30',
    status: 'enabled',
    type: 'basic',
    sort: 2
  },
  {
    id: '3',
    name: '专业版',
    price: '¥2000',
    cycle: '永久授权 (V1/V2/V3)',
    createdAt: '2024-01-15 10:30',
    updatedAt: '2024-01-15 10:30',
    status: 'enabled',
    type: 'pro',
    sort: 3
  },
  {
    id: '4',
    name: '定制版',
    price: '定制方案',
    cycle: '按需定价',
    createdAt: '2024-01-15 10:30',
    updatedAt: '2024-01-15 10:30',
    status: 'enabled',
    type: 'custom',
    sort: 4
  }
])

const editVisible = ref(false)
const selectedPackage = ref<any>(null)

const handleEdit = (row: any) => {
  selectedPackage.value = row
  editVisible.value = true
}

const handleUpdate = (updatedData: any) => {
  const index = tableData.value.findIndex(item => item.id === updatedData.id)
  if (index !== -1) {
    tableData.value[index] = { ...tableData.value[index], ...updatedData }
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
