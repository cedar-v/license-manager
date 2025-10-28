<template>
  <div class="license-search-container">
    <div class="background-section">
      <div class="center-title">
        <h1 class="platform-title">{{ t('pages.licenses.platform') }}</h1>
      </div>

      <div class="action-section">
        <div class="action-row">
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
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getCustomers, type Customer } from '@/api/customer'
import { getLicenses } from '@/api/license'

const { t } = useI18n()
const router = useRouter()

const selectedCustomer = ref<string>('')
const customers = ref<Customer[]>([])

const selectedCustomerInfo = computed(() => {
  if (!selectedCustomer.value) return null
  return customers.value.find(c => c.id === selectedCustomer.value) || null
})

const loadCustomers = async () => {
  try {
    const response = await getCustomers({ status: 'active', page_size: 100 })
    const allCustomers = response.data.list

    const customersWithLicenses = []
    for (const customer of allCustomers) {
      try {
        const licenseResponse = await getLicenses({
          customer_id: customer.id,
          page_size: 1
        })
        if (licenseResponse.data.list.length > 0) {
          customersWithLicenses.push(customer)
        }
      } catch (error) {
        console.warn(`��7 ${customer.customer_name} ��C1%:`, error)
        customersWithLicenses.push(customer)
      }
    }

    customers.value = customersWithLicenses
  } catch (error) {
    console.error('Failed to load customers:', error)
    ElMessage.error('�}�7h1%')
  }
}

const handleCustomerChange = () => {
  // �7	����;�
}

const handleQuery = async () => {
  if (!selectedCustomer.value) {
    ElMessage.warning(t('pages.licenses.message.selectCustomerFirst'))
    return
  }

  try {
    const response = await getLicenses({
      customer_id: selectedCustomer.value,
      page_size: 1
    })

    if (response.data.list.length === 0) {
      ElMessage.warning(t('pages.licenses.message.noLicenseWarning'))
      return
    }

    router.push({
      name: 'licenses-list',
      query: {
        customerId: selectedCustomer.value,
        customerName: selectedCustomerInfo.value?.customer_name || ''
      }
    })
  } catch (error) {
    console.error('��C1%:', error)
    ElMessage.error(t('pages.licenses.message.queryLicenseError'))
  }
}

const handleCreateLicense = () => {
  router.push({
    name: 'licenses-create',
    query: {
      customerId: selectedCustomer.value || '',
      customerName: selectedCustomerInfo.value?.customer_name || ''
    }
  })
}

onMounted(() => {
  loadCustomers()
})
</script>

<style scoped lang="scss">
@use '@/assets/styles/variables.scss' as *;
@use 'sass:color';

.license-search-container {
  height: calc(100vh - 80px);
  width: 100%;
  overflow: hidden;
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

  :deep(.el-select) {
    margin: 0;
  }

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

  &:deep(.el-button) {
    margin: 0;
    height: 40px;
    line-height: 1;
  }

  &.query-btn {
    min-width: 91px;
    background-color: $primary-color;

    &:hover {
      background-color: color.adjust($primary-color, $lightness: -10%);
    }
  }

  &.create-btn {
    min-width: 137px;
    background-color: $primary-color;

    &:hover {
      background-color: color.adjust($primary-color, $lightness: -10%);
    }
  }
}

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
