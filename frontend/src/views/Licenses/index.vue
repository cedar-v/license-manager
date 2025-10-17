<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-12 00:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-10-17 10:51:55
 * @FilePath: /frontend/src/views/Licenses.vue
 * @Description: 授权管理页面
-->
<template>
  <Layout app-name="Cedar-V" :page-title="t('pages.licenses.title')">
    <div class="license-container">
      <!-- 主要内容区域 -->
      <div v-if="!showSubRoute" class="main-content">
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

      </div>

      <!-- 子路由视图区域 -->
      <div v-if="showSubRoute" class="sub-route-section">
        <router-view />
      </div>
    </div>

  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import Layout from '@/components/common/layout/Layout.vue'
import { getCustomers, type Customer } from '@/api/customer'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

// 响应式数据
const selectedCustomer = ref<string>('')
const customers = ref<Customer[]>([])
const showSubRoute = ref(false)

// 选中的客户信息
const selectedCustomerInfo = computed(() => {
  if (!selectedCustomer.value) return null
  return customers.value.find(c => c.id === selectedCustomer.value) || null
})


// 方法
const loadCustomers = async () => {
  try {
    const response = await getCustomers({ status: 'active', page_size: 100 })
    customers.value = response.data.list
  } catch (error) {
    console.error('Failed to load customers:', error)
    ElMessage.error('加载客户列表失败')
  }
}


const handleCustomerChange = () => {
  // 客户选择变化时的处理逻辑
}

const handleQuery = () => {
  // 跳转到列表页，并传递客户信息参数
  showSubRoute.value = true
  router.push({
    name: 'licenses-list',
    query: {
      customerId: selectedCustomer.value || '',
      customerName: selectedCustomerInfo.value?.customer_name || ''
    }
  })
}

const handleCreateLicense = () => {
  // 跳转到创建授权页，并传递客户信息参数
  showSubRoute.value = true
  router.push({
    name: 'licenses-create',
    query: {
      customerId: selectedCustomer.value || '',
      customerName: selectedCustomerInfo.value?.customer_name || ''
    }
  })
}





// 监听路由变化，控制子路由显示
watch(() => route.name, (newName) => {
  // 如果是子路由，显示子路由区域
  showSubRoute.value = newName === 'licenses-list' || newName === 'licenses-create'
}, { immediate: true })

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

.sub-route-section {
  flex: 1;
  background: #fff;
  display: flex;
  flex-direction: column;
  overflow: hidden;
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