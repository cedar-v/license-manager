<template>
  <!-- 页面模式 -->
  <div v-if="isPageMode" class="license-form-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">{{ isEdit ? '编辑授权' : t('pages.licenses.actions.createLicense') }}</h1>
        <div class="header-actions">
          <el-button @click="handleBack">
            {{ t('customers.actions.cancel') }}
          </el-button>
          <el-button
            type="primary"
            :loading="submitting"
            @click="handleSubmit"
          >
            {{ t('customers.actions.save') }}
          </el-button>
        </div>
      </div>
    </div>
    
    <div class="page-content">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="120px"
        class="license-form"
      >
        <div class="form-section">
          <h3 class="section-title">{{ t('pages.licenses.form.basicInfo') }}</h3>

          <el-form-item
            :label="t('pages.licenses.form.selectCustomer')"
            prop="customer_id"
          >
            <el-select
              v-model="formData.customer_id"
              :placeholder="t('pages.licenses.form.placeholder.selectCustomer')"
              filterable
              style="width: 100%"
              :disabled="isEdit"
            >
              <el-option
                v-for="customer in customers"
                :key="customer.id"
                :label="customer.customer_name"
                :value="customer.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item
            :label="t('pages.licenses.form.description')"
            prop="description"
          >
            <el-input
              v-model="formData.description"
              :placeholder="t('pages.licenses.form.placeholder.enterDescription')"
              type="textarea"
              :rows="3"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>

          <el-form-item
            :label="t('pages.licenses.form.licenseType')"
            prop="license_type"
          >
            <el-select
              v-model="formData.license_type"
              :placeholder="t('pages.licenses.form.placeholder.selectLicenseType')"
              style="width: 100%"
            >
              <el-option label="标准版" value="standard" />
              <el-option label="专业版" value="professional" />
              <el-option label="企业版" value="enterprise" />
              <el-option label="定制版" value="custom" />
            </el-select>
          </el-form-item>

          <el-form-item
            :label="t('pages.licenses.form.expiryDate')"
            prop="expiry_date"
          >
            <el-date-picker
              v-model="formData.expiry_date"
              type="datetime"
              :placeholder="t('pages.licenses.form.placeholder.selectExpiryDate')"
              style="width: 100%"
              format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss"
              :disabled-date="disabledDate"
            />
          </el-form-item>

          <el-form-item
            :label="t('pages.licenses.form.maxUsers')"
            prop="max_users"
          >
            <el-input-number
              v-model="formData.max_users"
              :min="1"
              :max="10000"
              :placeholder="t('pages.licenses.form.placeholder.enterMaxUsers')"
              style="width: 100%"
            />
          </el-form-item>
        </div>

        <div class="form-section">
          <h3 class="section-title">高级配置</h3>

          <el-form-item
            :label="t('pages.licenses.form.features')"
            prop="features"
          >
            <el-checkbox-group v-model="formData.features">
              <el-checkbox label="basic" value="basic">基础功能</el-checkbox>
              <el-checkbox label="advanced" value="advanced">高级功能</el-checkbox>
              <el-checkbox label="api" value="api">API接口</el-checkbox>
              <el-checkbox label="export" value="export">数据导出</el-checkbox>
              <el-checkbox label="backup" value="backup">数据备份</el-checkbox>
              <el-checkbox label="analytics" value="analytics">数据分析</el-checkbox>
            </el-checkbox-group>
          </el-form-item>

          <el-form-item
            :label="t('pages.licenses.form.ipRestrictions')"
            prop="ip_restrictions"
          >
            <el-input
              v-model="ipRestrictionsText"
              :placeholder="t('pages.licenses.form.placeholder.enterIpRestrictions')"
              type="textarea"
              :rows="3"
              @blur="handleIpRestrictionsChange"
            />
            <div class="form-tip">多个IP地址请用换行分隔，支持CIDR格式（如：192.168.1.0/24）</div>
          </el-form-item>

          <el-form-item
            :label="t('pages.licenses.form.hardwareId')"
            prop="hardware_id"
          >
            <el-input
              v-model="formData.hardware_id"
              :placeholder="t('pages.licenses.form.placeholder.enterHardwareId')"
              maxlength="100"
            />
            <div class="form-tip">硬件ID用于绑定特定设备，留空表示不限制设备</div>
          </el-form-item>
        </div>
      </el-form>
    </div>
  </div>

  <!-- 对话框模式 -->
  <el-dialog
    v-else
    :model-value="visible"
    :title="isEdit ? '编辑授权' : t('pages.licenses.actions.createLicense')"
    width="800px"
    :before-close="handleClose"
    @update:model-value="$emit('update:visible', $event)"
  >
    <el-form
      ref="formRef"
      :model="formData"
      :rules="formRules"
      label-width="120px"
      class="license-form"
    >
      <div class="form-section">
        <h3 class="section-title">{{ t('pages.licenses.form.basicInfo') }}</h3>

        <el-form-item
          :label="t('pages.licenses.form.selectCustomer')"
          prop="customer_id"
        >
          <el-select
            v-model="formData.customer_id"
            :placeholder="t('pages.licenses.form.placeholder.selectCustomer')"
            filterable
            style="width: 100%"
            :disabled="isEdit"
          >
            <el-option
              v-for="customer in customers"
              :key="customer.id"
              :label="customer.customer_name"
              :value="customer.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item
          :label="t('pages.licenses.form.description')"
          prop="description"
        >
          <el-input
            v-model="formData.description"
            :placeholder="t('pages.licenses.form.placeholder.enterDescription')"
            type="textarea"
            :rows="3"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>

        <el-form-item
          :label="t('pages.licenses.form.licenseType')"
          prop="license_type"
        >
          <el-select
            v-model="formData.license_type"
            :placeholder="t('pages.licenses.form.placeholder.selectLicenseType')"
            style="width: 100%"
          >
            <el-option label="标准版" value="standard" />
            <el-option label="专业版" value="professional" />
            <el-option label="企业版" value="enterprise" />
            <el-option label="定制版" value="custom" />
          </el-select>
        </el-form-item>

        <el-form-item
          :label="t('pages.licenses.form.expiryDate')"
          prop="expiry_date"
        >
          <el-date-picker
            v-model="formData.expiry_date"
            type="datetime"
            :placeholder="t('pages.licenses.form.placeholder.selectExpiryDate')"
            style="width: 100%"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            :disabled-date="disabledDate"
          />
        </el-form-item>

        <el-form-item
          :label="t('pages.licenses.form.maxUsers')"
          prop="max_users"
        >
          <el-input-number
            v-model="formData.max_users"
            :min="1"
            :max="10000"
            :placeholder="t('pages.licenses.form.placeholder.enterMaxUsers')"
            style="width: 100%"
          />
        </el-form-item>
      </div>

      <div class="form-section">
        <h3 class="section-title">高级配置</h3>

        <el-form-item
          :label="t('pages.licenses.form.features')"
          prop="features"
        >
          <el-checkbox-group v-model="formData.features">
            <el-checkbox label="basic" value="basic">基础功能</el-checkbox>
            <el-checkbox label="advanced" value="advanced">高级功能</el-checkbox>
            <el-checkbox label="api" value="api">API接口</el-checkbox>
            <el-checkbox label="export" value="export">数据导出</el-checkbox>
            <el-checkbox label="backup" value="backup">数据备份</el-checkbox>
            <el-checkbox label="analytics" value="analytics">数据分析</el-checkbox>
          </el-checkbox-group>
        </el-form-item>

        <el-form-item
          :label="t('pages.licenses.form.ipRestrictions')"
          prop="ip_restrictions"
        >
          <el-input
            v-model="ipRestrictionsText"
            :placeholder="t('pages.licenses.form.placeholder.enterIpRestrictions')"
            type="textarea"
            :rows="3"
            @blur="handleIpRestrictionsChange"
          />
          <div class="form-tip">多个IP地址请用换行分隔，支持CIDR格式（如：192.168.1.0/24）</div>
        </el-form-item>

        <el-form-item
          :label="t('pages.licenses.form.hardwareId')"
          prop="hardware_id"
        >
          <el-input
            v-model="formData.hardware_id"
            :placeholder="t('pages.licenses.form.placeholder.enterHardwareId')"
            maxlength="100"
          />
          <div class="form-tip">硬件ID用于绑定特定设备，留空表示不限制设备</div>
        </el-form-item>
      </div>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">
          {{ t('customers.actions.cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="handleSubmit"
        >
          {{ t('customers.actions.save') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createLicense, updateLicense, type License, type AuthorizationCode, type AuthorizationCodeCreateRequest, type LicenseUpdateRequest } from '@/api/license'
import { getCustomers, type Customer } from '@/api/customer'

// Props & Emits
interface Props {
  visible?: boolean
  licenseData?: AuthorizationCode | null
  customers?: Customer[]
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  licenseData: null,
  customers: () => []
})

const emit = defineEmits<{
  'update:visible': [value: boolean]
  success: []
}>()

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

// 判断是否为页面模式
const isPageMode = computed(() => {
  return route.name === 'licenses-create' || route.name === 'licenses-edit'
})

// Refs
const formRef = ref<FormInstance>()
const submitting = ref(false)
const pageCustomers = ref<Customer[]>([])

// Computed
const isEdit = computed(() => !!props.licenseData)

// 获取客户列表（页面模式下使用）
const customers = computed(() => {
  return isPageMode.value ? pageCustomers.value : props.customers
})

// Form Data
const formData = reactive({
  customer_id: '',
  description: '',
  license_type: '',
  expiry_date: '',
  max_users: 1,
  features: [] as string[],
  ip_restrictions: [] as string[],
  hardware_id: ''
})

// IP限制文本（用于显示和编辑）
const ipRestrictionsText = ref('')

// Form Rules
const formRules: FormRules = {
  customer_id: [
    { required: true, message: t('pages.licenses.validation.customerRequired'), trigger: 'change' }
  ],
  description: [
    { required: true, message: t('pages.licenses.validation.descriptionRequired'), trigger: 'blur' }
  ],
  license_type: [
    { required: true, message: t('pages.licenses.validation.licenseTypeRequired'), trigger: 'change' }
  ],
  expiry_date: [
    { required: true, message: t('pages.licenses.validation.expiryDateRequired'), trigger: 'change' }
  ],
  max_users: [
    { required: true, message: t('pages.licenses.validation.maxUsersMin'), trigger: 'blur' },
    { type: 'number', min: 1, message: t('pages.licenses.validation.maxUsersMin'), trigger: 'blur' }
  ]
}

// Methods
const resetForm = () => {
  formData.customer_id = ''
  formData.description = ''
  formData.license_type = ''
  formData.expiry_date = ''
  formData.max_users = 1
  formData.features = []
  formData.ip_restrictions = []
  formData.hardware_id = ''
  ipRestrictionsText.value = ''
}

const loadFormData = () => {
  if (props.licenseData) {
    formData.customer_id = props.licenseData.customer_id || ''
    formData.description = props.licenseData.description || ''
    formData.license_type = ''
    formData.expiry_date = props.licenseData.end_date || ''
    formData.max_users = props.licenseData.max_activations || 1
    formData.features = []
    formData.ip_restrictions = []
    formData.hardware_id = ''
    ipRestrictionsText.value = ''
  } else {
    resetForm()
  }
}

const handleIpRestrictionsChange = () => {
  formData.ip_restrictions = ipRestrictionsText.value
    .split('\n')
    .map(ip => ip.trim())
    .filter(ip => ip.length > 0)
}

const disabledDate = (time: Date) => {
  return time.getTime() < Date.now() - 24 * 60 * 60 * 1000
}

const validateForm = async (): Promise<boolean> => {
  if (!formRef.value) return false

  try {
    await formRef.value.validate()

    // 验证IP地址格式
    if (formData.ip_restrictions.length > 0) {
      const ipRegex = /^((\d{1,3}\.){3}\d{1,3}(\/\d{1,2})?|(\d{1,3}\.){3}\d{1,3})$/
      for (const ip of formData.ip_restrictions) {
        if (!ipRegex.test(ip)) {
          ElMessage.error(`${t('pages.licenses.validation.ipFormatError')}: ${ip}`)
          return false
        }
      }
    }

    return true
  } catch {
    return false
  }
}

const handleSubmit = async () => {
  if (!(await validateForm())) return

  try {
    submitting.value = true

    if (isEdit.value && props.licenseData) {
      // 编辑授权
      const updateData: LicenseUpdateRequest = {
        description: formData.description,
        expiry_date: formData.expiry_date,
        max_users: formData.max_users,
        features: formData.features,
        ip_restrictions: formData.ip_restrictions.length > 0 ? formData.ip_restrictions : undefined,
        hardware_id: formData.hardware_id || undefined
      }

      await updateLicense(props.licenseData.id, updateData)
      ElMessage.success(t('pages.licenses.message.updateSuccess'))
    } else {
      // 创建授权
      const createData: AuthorizationCodeCreateRequest = {
        customer_id: formData.customer_id,
        description: formData.description,
        license_type: formData.license_type,
        expiry_date: formData.expiry_date,
        max_users: formData.max_users,
        features: formData.features,
        ip_restrictions: formData.ip_restrictions.length > 0 ? formData.ip_restrictions : undefined,
        hardware_id: formData.hardware_id || undefined
      }

      await createLicense(createData)
      ElMessage.success(t('pages.licenses.message.createSuccess'))
    }

    if (isPageMode.value) {
      // 页面模式下跳转回主页面
      router.push({ name: 'licenses' })
    } else {
      // 对话框模式下触发成功事件
      emit('success')
    }
  } catch (error) {
    console.error('Failed to save license:', error)
    ElMessage.error(isEdit.value
      ? t('pages.licenses.message.updateError')
      : t('pages.licenses.message.createError')
    )
  } finally {
    submitting.value = false
  }
}

const handleClose = () => {
  emit('update:visible', false)
  // 延迟重置表单，避免关闭动画期间表单闪烁
  setTimeout(() => {
    formRef.value?.clearValidate()
    resetForm()
  }, 300)
}

// 页面模式下的返回功能
const handleBack = () => {
  router.push({ name: 'licenses' })
}

// 加载客户数据（页面模式）
const loadCustomers = async () => {
  try {
    const response = await getCustomers({ status: 'active', page_size: 1000 })
    pageCustomers.value = response.data.list
  } catch (error) {
    console.error('Failed to load customers:', error)
    ElMessage.error('加载客户列表失败')
  }
}

// Watchers
watch(() => props.visible, (visible) => {
  if (visible) {
    loadFormData()
  }
})

// 页面模式下的初始化
watch(() => route.query, (query) => {
  if (isPageMode.value && query.customerId) {
    formData.customer_id = query.customerId as string
  }
}, { immediate: true })

// 生命周期
onMounted(() => {
  if (isPageMode.value) {
    loadCustomers()
    loadFormData()
  }
})

</script>

<style scoped lang="scss">
@import '@/assets/styles/variables.scss';

// 页面模式样式
.license-form-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: #fff;
}

.page-header {
  padding: 24px;
  border-bottom: 1px solid $border-color-light;
  background: #fff;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  max-width: 1200px;
  margin: 0 auto;
}

.page-title {
  font-size: 24px;
  font-weight: 600;
  color: $text-color-primary;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.page-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  box-sizing: border-box;
}

.license-form {
  max-height: none;
  overflow-y: visible;
  padding-right: 0;
}

.form-section {
  margin-bottom: 32px;

  &:last-child {
    margin-bottom: 0;
  }
}

.section-title {
  font-size: 16px;
  font-weight: 600;
  color: $text-color-primary;
  margin: 0 0 16px 0;
  padding-bottom: 8px;
  border-bottom: 1px solid $border-color-lighter;
}

.form-tip {
  font-size: 12px;
  color: $text-color-secondary;
  margin-top: 4px;
  line-height: 1.4;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

:deep(.el-checkbox-group) {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;

  .el-checkbox {
    margin-right: 0;
  }
}

:deep(.el-form-item__label) {
  font-weight: 500;
  color: $text-color-primary;
}

:deep(.el-input__wrapper) {
  box-shadow: 0 0 0 1px $border-color-base inset;

  &:hover {
    box-shadow: 0 0 0 1px $primary-color inset;
  }

  &.is-focus {
    box-shadow: 0 0 0 1px $primary-color inset;
  }
}

:deep(.el-textarea__inner) {
  box-shadow: 0 0 0 1px $border-color-base inset;

  &:hover {
    box-shadow: 0 0 0 1px $primary-color inset;
  }

  &:focus {
    box-shadow: 0 0 0 1px $primary-color inset;
  }
}

// 响应式设计
@media (max-width: 768px) {
  :deep(.el-dialog) {
    width: 95% !important;
    margin: 2.5vh auto;
  }

  .license-form {
    max-height: 60vh;
  }

  :deep(.el-form-item__label) {
    text-align: left !important;
  }
}
</style>