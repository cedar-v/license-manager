<template>
  <div class="license-form-container">
    <!-- 顶部横向区域 - 与表单平级 -->
    <div class="top-section">
      <!-- 面包屑导航 -->
      <div class="breadcrumb">
        <span class="breadcrumb-item">{{
          t('pages.licenses.form.breadcrumb.licenseManagement')
        }}</span>
        <span class="breadcrumb-separator">></span>
        <span class="breadcrumb-current">{{
          isEdit
            ? t('pages.licenses.form.breadcrumb.editLicense')
            : t('pages.licenses.form.breadcrumb.createLicense')
        }}</span>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <el-button @click="handleCancel">{{ t('pages.licenses.form.actions.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{
          isEdit ? t('pages.licenses.form.actions.update') : t('pages.licenses.form.actions.create')
        }}</el-button>
      </div>
    </div>

    <!-- 主表单区域 -->
    <div>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-position="top">
        <!-- 基本信息 -->
        <div class="license-form">
          <h3 class="section-title">{{ t('pages.licenses.form.sections.basicInfo') }}</h3>

          <!-- 第一行：关联客户名称，客户ID -->
          <div class="fields-row">
            <el-form-item
              :label="t('pages.licenses.form.fields.customerName')"
              prop="customer_id"
              required
              class="field-item"
            >
              <el-select
                v-model="formData.customer_id"
                :placeholder="t('pages.licenses.form.placeholders.selectCustomer')"
                filterable
                remote
                :remote-method="searchCustomers"
                :loading="customerLoading"
                style="width: 100%"
              >
                <el-option
                  v-for="customer in customerOptions"
                  :key="customer.id"
                  :label="customer.customer_name"
                  :value="customer.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item
              :label="t('pages.licenses.form.fields.customerId')"
              prop="customer_code"
              required
              class="field-item"
            >
              <el-input
                v-model="formData.customer_code"
                :placeholder="t('pages.licenses.form.placeholders.autoGenerateOrManual')"
                disabled
              />
            </el-form-item>
          </div>

          <!-- 第二行：备注单独一行 -->
          <div class="fields-row">
            <el-form-item
              :label="t('pages.licenses.form.fields.description')"
              prop="description"
              required
              class="field-item-full"
            >
              <el-input
                v-model="formData.description"
                type="textarea"
                :rows="4"
                :placeholder="t('pages.licenses.form.placeholders.enterDescription')"
              />
            </el-form-item>
          </div>
        </div>

        <!-- 授权配置 -->
        <div class="license-form">
          <h3 class="section-title">{{ t('pages.licenses.form.sections.licenseConfig') }}</h3>

          <!-- 第三行：授权期限类型，起止时间 -->
          <div class="fields-row">
            <el-form-item
              :label="t('pages.licenses.form.fields.validityPeriod')"
              prop="validity_type"
              required
              class="field-item"
            >
              <el-select
                v-model="formData.validity_type"
                :placeholder="t('pages.licenses.form.placeholders.selectValidityType')"
                style="width: 100%"
                @change="handleValidityTypeChange"
              >
                <el-option
                  :label="t('pages.licenses.form.validityTypes.permanent')"
                  value="permanent"
                />
                <el-option
                  :label="t('pages.licenses.form.validityTypes.limited')"
                  value="limited"
                />
              </el-select>
            </el-form-item>

            <el-form-item
              v-if="formData.validity_type === 'limited'"
              :label="t('pages.licenses.form.fields.dateRange')"
              prop="date_range"
              required
              class="field-item"
            >
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                :range-separator="t('chart.licenseTrend.datePicker.rangeSeparator')"
                :start-placeholder="t('pages.licenses.form.placeholders.startDate')"
                :end-placeholder="t('pages.licenses.form.placeholders.endDate')"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled-date="disabledDate"
                @change="handleDateRangeChange"
              />
            </el-form-item>

            <el-form-item
              v-if="formData.validity_type === 'permanent'"
              :label="t('pages.licenses.form.fields.validityDays')"
              class="field-item"
            >
              <el-input
                :value="t('pages.licenses.form.validityTypes.permanent') + '（365000天）'"
                readonly
                style="width: 100%"
              />
            </el-form-item>
          </div>

          <!-- 第四行：最大激活设备数，部署类型，加密类型 -->
          <div class="fields-row">
            <el-form-item
              :label="t('pages.licenses.form.fields.maxActivations')"
              prop="max_activations"
              required
              class="field-item"
            >
              <el-input-number
                v-model="formData.max_activations"
                :min="1"
                :placeholder="t('pages.licenses.form.placeholders.enterMaxActivations')"
                style="width: 100%"
              />
            </el-form-item>

            <el-form-item
              :label="t('pages.licenses.form.fields.deploymentType')"
              prop="deployment_type"
              required
              class="field-item"
            >
              <el-select
                v-model="formData.deployment_type"
                :placeholder="t('pages.licenses.form.placeholders.selectDeploymentType')"
                style="width: 100%"
                @change="handleDeploymentTypeChange"
              >
                <el-option
                  v-for="option in deploymentTypeOptions"
                  :key="option.key"
                  :label="option.display"
                  :value="option.key"
                />
              </el-select>
            </el-form-item>

            <el-form-item
              :label="t('pages.licenses.form.fields.encryptionType')"
              prop="encryption_type"
              required
              class="field-item"
            >
              <el-select
                v-model="formData.encryption_type"
                :placeholder="t('pages.licenses.form.placeholders.selectEncryptionType')"
                style="width: 100%"
              >
                <el-option
                  v-for="option in encryptionTypeOptions"
                  :key="option.key"
                  :label="option.display"
                  :value="option.key"
                />
              </el-select>
            </el-form-item>
          </div>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  createLicense,
  updateLicense,
  getLicenseDetail,
  type AuthorizationCodeCreateRequest,
  type LicenseUpdateRequest
} from '@/api/license'
import { getCustomers, type Customer } from '@/api/customer'
import { getEnumOptions, type RawEnumItem } from '@/api/enum'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

// 表单引用
const formRef = ref<FormInstance>()

// 响应式数据
const submitting = ref(false)
const customerLoading = ref(false)
const customerOptions = ref<Customer[]>([])

// 枚举选项
const deploymentTypeOptions = ref<RawEnumItem[]>([])
const encryptionTypeOptions = ref<RawEnumItem[]>([])

// 表单数据
const formData = reactive<
  AuthorizationCodeCreateRequest & {
    customer_code: string
    validity_type: 'permanent' | 'limited'
  }
>({
  customer_id: '',
  customer_code: '',
  validity_type: 'limited',
  description: '',
  validity_days: 365,
  deployment_type: 'standalone',
  encryption_type: 'standard',
  max_activations: 1,
  feature_config: {},
  usage_limits: {},
  custom_parameters: {}
})

// 独立的响应式变量
const dateRange = ref<[string, string] | null>(null)

// 计算属性
const isEdit = computed(() => {
  return route.name === 'licenses-edit' && route.params.id
})

// 禁止选择今天之前的日期
const disabledDate = (time: Date) => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return time.getTime() < today.getTime()
}

// 表单验证规则
const formRules: FormRules = {
  customer_id: [
    {
      required: true,
      message: t('pages.licenses.form.validation.customerRequired'),
      trigger: 'change'
    },
    {
      validator: (_rule: any, value: any, callback: any) => {
        if (value && customerOptions.value.length > 0) {
          const selectedCustomer = customerOptions.value.find(c => c.id === value)
          if (!selectedCustomer) {
            callback(new Error(t('pages.licenses.form.validation.customerInvalid')))
          } else {
            callback()
          }
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  description: [
    {
      required: true,
      message: t('pages.licenses.form.validation.descriptionRequired'),
      trigger: 'blur'
    },
    {
      min: 1,
      max: 500,
      message: t('pages.licenses.form.validation.descriptionLength'),
      trigger: 'blur'
    }
  ],
  validity_type: [
    {
      required: true,
      message: t('pages.licenses.form.validation.validityTypeRequired'),
      trigger: 'change'
    }
  ],
  date_range: [
    {
      required: true,
      message: t('pages.licenses.form.validation.dateRangeRequired'),
      trigger: 'change',
      validator: (_rule: any, _value: any, callback: any) => {
        if (
          formData.validity_type === 'limited' &&
          (!dateRange.value || dateRange.value.length !== 2)
        ) {
          callback(new Error(t('pages.licenses.form.validation.dateRangeRequired')))
        } else {
          callback()
        }
      }
    },
    {
      validator: (_rule: any, _value: any, callback: any) => {
        if (
          formData.validity_type === 'limited' &&
          dateRange.value &&
          dateRange.value.length === 2
        ) {
          const startDate = new Date(dateRange.value[0])
          const endDate = new Date(dateRange.value[1])
          const today = new Date()
          today.setHours(0, 0, 0, 0)

          if (startDate < today) {
            callback(new Error(t('pages.licenses.form.validation.startDateNotPast')))
          } else if (endDate <= startDate) {
            callback(new Error(t('pages.licenses.form.validation.endDateAfterStart')))
          } else {
            // 检查日期范围是否超过合理范围（比如10年）
            const maxDays = 3650 // 10年
            const diffTime = endDate.getTime() - startDate.getTime()
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
            if (diffDays > maxDays) {
              callback(new Error(t('pages.licenses.form.validation.validityPeriodTooLong')))
            } else {
              callback()
            }
          }
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  deployment_type: [
    {
      required: true,
      message: t('pages.licenses.form.validation.deploymentTypeRequired'),
      trigger: 'change'
    }
  ],
  encryption_type: [
    {
      required: true,
      message: t('pages.licenses.form.validation.encryptionTypeRequired'),
      trigger: 'change'
    }
  ],
  max_activations: [
    {
      required: true,
      message: t('pages.licenses.form.validation.maxActivationsRequired'),
      trigger: 'blur'
    },
    {
      type: 'number',
      min: 1,
      max: 999999,
      message: t('pages.licenses.form.validation.maxActivationsRange'),
      trigger: 'blur'
    },
    {
      validator: (_rule: any, value: any, callback: any) => {
        if (value && !Number.isInteger(Number(value))) {
          callback(new Error(t('pages.licenses.form.validation.maxActivationsInteger')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 方法
const loadEnumOptions = async () => {
  try {
    // 加载部署类型选项
    const deploymentResponse = await getEnumOptions('deployment_type')
    if (deploymentResponse.code === '000000') {
      deploymentTypeOptions.value = deploymentResponse.data.items
    } else {
      throw new Error(deploymentResponse.message || '加载部署类型选项失败')
    }

    // 加载加密类型选项
    const encryptionResponse = await getEnumOptions('encryption_type')
    if (encryptionResponse.code === '000000') {
      encryptionTypeOptions.value = encryptionResponse.data.items
    } else {
      throw new Error(encryptionResponse.message || '加载加密类型选项失败')
    }
  } catch (error: any) {
    console.error('加载枚举选项失败:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('pages.licenses.form.messages.loadEnumError') + ': ' + errorMessage)
    } else {
      ElMessage.error(t('pages.licenses.form.messages.loadEnumErrorRetry'))
    }
  }
}

// 搜索或加载所有客户
const searchCustomers = async (query: string) => {
  try {
    customerLoading.value = true
    const response = await getCustomers({
      page: 1,
      page_size: 100,
      customer_name: query || undefined,
      status: 'active'
    })
    customerOptions.value = response.data.list
  } catch (error: any) {
    console.error('搜索客户失败:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('pages.licenses.form.messages.searchCustomerError') + ': ' + errorMessage)
    } else {
      ElMessage.error(t('pages.licenses.form.messages.searchCustomerErrorRetry'))
    }
    customerOptions.value = []
  } finally {
    customerLoading.value = false
  }
}

// 表单初始化
const loadCustomerInfo = () => {
  // 从路由参数获取客户信息
  const customerId = route.query.customerId as string
  const customerName = route.query.customerName as string

  if (customerId && customerName) {
    formData.customer_id = customerId
    customerOptions.value = [{
      id: customerId,
      customer_name: customerName,
      customer_code: customerId,
      customer_type: '',
      customer_type_display: '',
      contact_person: '',
      email: '',
      customer_level: '',
      customer_level_display: '',
      status: 'active',
      status_display: '',
      created_at: '',
      updated_at: '',
      created_by: '',
      updated_by: ''
    } as Customer]
  }

}

// 处理授权期限类型变化
const handleValidityTypeChange = (type: 'permanent' | 'limited') => {
  if (type === 'permanent') {
    formData.validity_days = 365000
    dateRange.value = null
  } else {
    // 有限期限时，如果已有日期范围，重新计算天数
    if (dateRange.value && dateRange.value.length === 2) {
      const startDate = new Date(dateRange.value[0])
      const endDate = new Date(dateRange.value[1])
      const diffTime = endDate.getTime() - startDate.getTime()
      formData.validity_days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    }
  }
}

// 处理部署类型变化
const handleDeploymentTypeChange = () => {
  // 部署类型变化时的处理逻辑（如需要可以在这里添加）
}

// 处理日期范围变化
const handleDateRangeChange = (dates: [string, string] | null) => {
  if (dates && dates.length === 2) {
    const startDate = new Date(dates[0])
    const endDate = new Date(dates[1])
    const diffTime = endDate.getTime() - startDate.getTime()
    formData.validity_days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  }
}

// 根据id获取授权详情
const loadLicenseDetail = async () => {
  if (!isEdit.value) return

  try {
    const id = route.params.id as string
    const response = await getLicenseDetail(id)

    if (response.code === '000000' && response.data) {
      const data = response.data

      // 填充表单数据
      formData.customer_id = data.customer_id || ''
      formData.description = data.description || ''
      formData.deployment_type = data.deployment_type || 'standalone'
      formData.encryption_type = data.encryption_type || 'standard'
      formData.max_activations = data.max_activations || 1

      // 设置表单字段
      formData.customer_code = (data as any).customer_code || data.customer_id || ''

      // 计算有效期天数和设置日期范围
      if (data.start_date && data.end_date) {
        const startDate = new Date(data.start_date)
        const endDate = new Date(data.end_date)
        const diffTime = endDate.getTime() - startDate.getTime()
        const days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
        formData.validity_days = days

        // 根据天数判断是永久还是有限
        if (days >= 365000) {
          formData.validity_type = 'permanent'
          dateRange.value = null
        } else {
          formData.validity_type = 'limited'
          dateRange.value = [data.start_date, data.end_date]
        }
      } else {
        // 如果没有日期信息，根据validity_days判断
        const validityDays = (data as any).validity_days
        if (validityDays && validityDays >= 365000) {
          formData.validity_type = 'permanent'
          formData.validity_days = 365000
        } else {
          formData.validity_type = 'limited'
          formData.validity_days = validityDays || 365
        }
      }

      // 设置客户选项
      if (data.customer_name) {
        customerOptions.value = [
          {
            id: data.customer_id,
            customer_name: data.customer_name,
            customer_code: (data as any).customer_code || '',
            customer_type: '',
            customer_type_display: '',
            contact_person: '',
            email: '',
            customer_level: '',
            customer_level_display: '',
            status: '',
            status_display: '',
            created_at: '',
            updated_at: '',
            created_by: '',
            updated_by: ''
          } as Customer
        ]
      }
    } else {
      throw new Error(response.message || t('pages.licenses.form.messages.loadDetailError'))
    }
  } catch (error: any) {
    console.error('加载授权详情失败:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('pages.licenses.form.messages.loadDetailError') + ': ' + errorMessage)
    } else {
      ElMessage.error(t('pages.licenses.form.messages.loadDetailErrorRetry'))
    }
  }
}

// 提交（新增/更新）授权数据
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    // 执行表单验证
    await formRef.value.validate()

    submitting.value = true

    if (isEdit.value) {
      // 更新授权
      const id = route.params.id as string
      const updateData: LicenseUpdateRequest = {
        description: formData.description,
        validity_days: formData.validity_days,
        deployment_type: formData.deployment_type,
        encryption_type: formData.encryption_type,
        max_activations: formData.max_activations,
        feature_config: formData.feature_config,
        usage_limits: formData.usage_limits,
        custom_parameters: formData.custom_parameters,
        change_type: 'other',
        reason: t('pages.licenses.form.messages.updateSuccess')
      }

      const response = await updateLicense(id, updateData)
      if (response.code === '000000') {
        ElMessage.success(response.message || t('pages.licenses.form.messages.updateSuccess'))
        // 返回列表页面
        router.back()
      } else {
        throw new Error(response.message || t('pages.licenses.form.messages.submitError'))
      }
    } else {
      // 创建授权
      const response = await createLicense(formData)
      if (response.code === '000000') {
        ElMessage.success(response.message || t('pages.licenses.form.messages.createSuccess'))
        // 返回列表页面
        router.back()
      } else {
        throw new Error(response.message || t('pages.licenses.form.messages.submitError'))
      }
    }
  } catch (error: any) {
    console.error('提交失败:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('pages.licenses.form.messages.submitError') + ': ' + errorMessage)
    } else {
      ElMessage.error(t('pages.licenses.form.messages.submitErrorRetry'))
    }
  } finally {
    submitting.value = false
  }
}

// 取消函数
const handleCancel = () => {
  ElMessageBox.confirm(
    t('pages.licenses.form.messages.cancelConfirm'),
    t('pages.licenses.form.messages.cancelTitle'),
    {
      confirmButtonText: t('pages.licenses.form.messages.cancelConfirmButton'),
      cancelButtonText: t('pages.licenses.form.messages.cancelCancelButton'),
      type: 'warning'
    }
  )
    .then(() => {
      router.back()
    })
    .catch(() => {
      // 用户选择继续编辑
    })
}

// 监听客户选择变化，自动填充客户ID
watch(
  () => formData.customer_id,
  newCustomerId => {
    if (newCustomerId) {
      const selectedCustomer = customerOptions.value.find(c => c.id === newCustomerId)
      if (selectedCustomer) {
        formData.customer_code = selectedCustomer.customer_code || selectedCustomer.id
      } else {
        // 如果没有找到客户，清空客户代码
        formData.customer_code = ''
      }
    } else {
      formData.customer_code = ''
    }
  }
)

// 生命周期
onMounted(async () => {
  await loadEnumOptions()
  loadCustomerInfo()
  // 如果没有从路由获取到客户信息，则加载所有客户
  if (!route.query.customerId) {
    await searchCustomers('')
  }
  await loadLicenseDetail()
})
</script>

<style scoped>
/* 主容器 */
.license-form-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--app-bg-color);
  padding: 15px;
}

/* 顶部横向区域 */
.top-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: transparent;
  margin-bottom: 10px;
}

.breadcrumb {
  font-size: 14px;
  color: var(--app-text-secondary);
}

.breadcrumb-item {
  color: var(--app-text-secondary);
}

.breadcrumb-separator {
  margin: 0 8px;
  color: var(--app-text-secondary);
}

.breadcrumb-current {
  color: var(--app-text-primary);
}

.form-actions {
  display: flex;
  gap: 20px;
}

/* 主表单区域 */
.license-form {
  padding: 20px;
  background: var(--app-content-bg);
  border-radius: 4px;
  width: 100%;
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--app-text-primary);
  margin: 0 0 20px 0;
}

/* 字段行布局 */
.fields-row {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 20px;
}

/* 基础字段：按照Figma设计的325px宽度 */
.fields-row .field-item {
  flex: 0 0 325px; /* 固定宽度325px */
  min-width: 280px; /* 移动端最小宽度 */
  margin-bottom: 0;
}

/* 备注字段：占满整行 */
.fields-row .field-item-full {
  flex: 1 1 100%;
  min-width: 100%;
  margin-bottom: 0;
}

/* Element Plus 样式覆盖 */
:deep(.el-form-item__label) {
  font-weight: 400;
  color: var(--app-text-primary);
  font-size: 14px;
  line-height: 1.5714285714285714em;
  margin-bottom: 8px;
}

:deep(.el-form-item.is-required .el-form-item__label::before) {
  content: '*';
  color: #f56c6c;
  margin-right: 4px;
}

:deep(.el-input__wrapper) {
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  height: 32px;
}

:deep(.el-input__inner) {
  color: var(--app-text-primary);
  font-size: 14px;
  line-height: 1.5em;
}

:deep(.el-input__inner::placeholder) {
  color: var(--app-text-secondary);
  font-size: 14px;
}

:deep(.el-select .el-input__wrapper) {
  border: 1px solid var(--app-border-color);
  height: 32px;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number .el-input__wrapper) {
  border: 1px solid var(--app-border-color);
  height: 32px;
}

:deep(.el-date-editor) {
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  height: 32px;
}

:deep(.el-date-editor .el-input__wrapper) {
  border: none;
  height: 30px;
}

:deep(.el-button--primary) {
  background-color: var(--el-color-primary-dark-2);
  border-color: var(--el-color-primary-dark-2);
  font-size: 14px;
  padding: 6px 16px;
}

:deep(.el-button--primary:hover) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

:deep(.el-button) {
  font-size: 14px;
  padding: 6px 16px;
}

/* 响应式设计 */
@media (max-width: 1440px) {
  .fields-row {
    flex-wrap: wrap;
  }

  .fields-row .field-item {
    min-width: 200px;
  }
}

@media (max-width: 768px) {
  .top-section {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .fields-row {
    flex-direction: column;
    gap: 16px;
  }

  .fields-row .field-item,
  .fields-row .field-item-full {
    width: 100%;
    flex: none;
  }
}
</style>
