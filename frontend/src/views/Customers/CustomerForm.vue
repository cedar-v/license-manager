<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-20 12:00:00
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-20 12:00:00
 * @FilePath: /frontend/src/views/Customers/CustomerForm.vue
 * @Description: 客户新增/编辑表单组件
-->
<template>
  <!-- 顶部横向区域 - 与表单平级 -->
  <div class="top-section">
    <!-- 面包屑导航 -->
    <div class="breadcrumb">
      <span class="breadcrumb-item">{{ t('customers.breadcrumb.customerManagement') }}</span>
      <span class="breadcrumb-separator">></span>
      <span class="breadcrumb-current">{{ isEdit ? t('customers.breadcrumb.editCustomer') : t('customers.breadcrumb.addCustomer') }}</span>
    </div>
    
    <!-- 操作按钮 -->
    <div class="form-actions">
      <el-button @click="handleCancel">{{ t('customers.actions.cancel') }}</el-button>
      <el-button type="primary" @click="handleSave" :loading="loading">{{ t('customers.actions.save') }}</el-button>
    </div>
  </div>

  <!-- 主表单区域 -->
  <div >
    <el-form :model="formData" :rules="formRules" ref="formRef" label-position="top">
      <!-- 基本信息 -->
      <div  class="customer-form">
        <h3 class="section-title">{{ t('customers.form.basicInfo') }}</h3>
        
        <!-- Flex布局 - 7个字段横向排列 -->
        <div class="fields-row-flex">
          <el-form-item :label="t('customers.form.customerName')" prop="name" required class="field-item">
            <el-input v-model="formData.name" :placeholder="t('customers.form.placeholder.enter')" />
          </el-form-item>
          
          <el-form-item :label="t('customers.form.customerType')" prop="type" required class="field-item">
            <el-select v-model="formData.type" :placeholder="t('customers.form.placeholder.select')" style="width: 100%">
              <el-option 
                v-for="option in customerTypeOptions" 
                :key="option.key" 
                :label="option.display" 
                :value="option.key" 
              />
            </el-select>
          </el-form-item>
          
          <el-form-item :label="t('customers.form.customerLevel')" prop="level" required class="field-item">
            <el-select v-model="formData.level" :placeholder="t('customers.form.placeholder.select')" style="width: 100%">
              <el-option 
                v-for="option in customerLevelOptions" 
                :key="option.key" 
                :label="option.display" 
                :value="option.key" 
              />
            </el-select>
          </el-form-item>
          
          <el-form-item :label="t('customers.form.contactPerson')" prop="contact" required class="field-item">
            <el-input v-model="formData.contact" :placeholder="t('customers.form.placeholder.enter')" />
          </el-form-item>
          
          <el-form-item :label="t('customers.form.email')" prop="email" class="field-item">
            <el-input v-model="formData.email" :placeholder="t('customers.form.placeholder.enter')" />
          </el-form-item>
          
          <el-form-item :label="t('customers.form.status')" prop="status" required class="field-item status-field">
            <el-radio-group v-model="formData.status">
              <el-radio 
                v-for="option in statusOptions" 
                :key="option.key" 
                :value="option.key"
              >
                {{ option.display }}
              </el-radio>
            </el-radio-group>
          </el-form-item>
          
          <el-form-item :label="t('customers.form.phone')" prop="phone" class="field-item">
            <el-input v-model="formData.phone" :placeholder="t('customers.form.placeholder.enter')" />
          </el-form-item>
        </div>
        
        <!-- Grid布局 - 地址跨列 -->
        <div class="fields-row-grid">
          <el-form-item :label="t('customers.form.address')" prop="address" class="address-field">
            <el-input v-model="formData.address" :placeholder="t('customers.form.placeholder.enter')" />
          </el-form-item>
        </div>
      </div>

      <!-- 商业信息 -->
      <div  class="customer-form">
        <h3 class="section-title">{{ t('customers.form.businessInfo') }}</h3>
        
        <div class="business-fields">
          <el-form-item :label="t('customers.form.companySize')" prop="companySize" required class="company-size-field">
            <el-select v-model="formData.companySize" :placeholder="t('customers.form.placeholder.select')" style="width: 100%">
              <el-option 
                v-for="option in companySizeOptions" 
                :key="option.key" 
                :label="option.display" 
                :value="option.key" 
              />
            </el-select>
          </el-form-item>
          
          <el-form-item :label="t('customers.form.description')" prop="description" class="description-field">
            <el-input
              v-model="formData.description"
              type="textarea"
              :rows="4"
              :placeholder="t('customers.form.placeholder.enter')"
              :maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </div>
      </div>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { createCustomer, updateCustomer, getCustomerDetail } from '@/api/customer'
import {
  getCustomerTypeEnums,
  getCustomerLevelEnums,
  getStatusEnums,
  getCompanySizeEnums,
  type RawEnumItem
} from '@/api/enum'

interface CustomerFormData {
  id?: string
  name: string
  contact: string
  email: string
  phone: string
  address: string
  type: string
  level: string
  status: string
  companySize: string
  description: string
}

const props = defineProps<{
  customerId?: string
  isEdit?: boolean
}>()

const emit = defineEmits<{
  save: [data: CustomerFormData]
  cancel: []
}>()

const { t } = useI18n()

const formRef = ref<FormInstance>()
const loading = ref(false)

// 枚举选项
const customerTypeOptions = ref<RawEnumItem[]>([])
const customerLevelOptions = ref<RawEnumItem[]>([])
const statusOptions = ref<RawEnumItem[]>([])
const companySizeOptions = ref<RawEnumItem[]>([])

// 加载枚举数据
const loadEnums = async () => {
  try {
    const [typeRes, levelRes, statusRes, sizeRes] = await Promise.all([
      getCustomerTypeEnums(),
      getCustomerLevelEnums(),
      getStatusEnums(),
      getCompanySizeEnums()
    ])
    
    if (typeRes.code === '000000') {
      customerTypeOptions.value = typeRes.data.items
    }
    if (levelRes.code === '000000') {
      customerLevelOptions.value = levelRes.data.items
    }
    if (statusRes.code === '000000') {
      statusOptions.value = statusRes.data.items
      // 设置默认状态（如果是新增模式）
      if (!props.isEdit && statusRes.data.items.length > 0) {
        formData.status = statusRes.data.items[0].key
      }
    }
    if (sizeRes.code === '000000') {
      companySizeOptions.value = sizeRes.data.items
    }
  } catch (error) {
    console.error(t('customers.message.loadEnumError'), error)
  }
}

// 表单数据
const formData = reactive<CustomerFormData>({
  name: '',
  contact: '',
  email: '',
  phone: '',
  address: '',
  type: '',
  level: '',
  status: '',
  companySize: '',
  description: ''
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: t('customers.validation.nameRequired'), trigger: 'blur' }
  ],
  contact: [
    { required: true, message: t('customers.validation.contactRequired'), trigger: 'blur' }
  ],
  type: [
    { required: true, message: t('customers.validation.typeRequired'), trigger: 'change' }
  ],
  level: [
    { required: true, message: t('customers.validation.levelRequired'), trigger: 'change' }
  ],
  status: [
    { required: true, message: t('customers.validation.statusRequired'), trigger: 'change' }
  ],
  companySize: [
    { required: true, message: t('customers.validation.companySizeRequired'), trigger: 'change' }
  ],
  email: [
    { type: 'email', message: t('customers.validation.emailFormat'), trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: t('customers.validation.phoneFormat'), trigger: 'blur' }
  ]
}

// 加载客户详情数据
const loadCustomerData = async () => {
  if (!props.customerId || !props.isEdit) return
  
  try {
    loading.value = true
    const response = await getCustomerDetail(props.customerId)
    
    if (response.code === '000000' && response.data) {
      const customerData = response.data
      // 将API返回的字段映射到表单字段
      Object.assign(formData, {
        id: customerData.id,
        name: customerData.customer_name,
        contact: customerData.contact_person,
        email: customerData.email || '',
        phone: customerData.phone || '',
        address: customerData.address || '',
        type: customerData.customer_type,
        level: customerData.customer_level,
        status: customerData.status,
        companySize: customerData.company_size || '',
        description: customerData.description || ''
      })
    } else {
      ElMessage.error(t('customers.message.getDetailError') + response.message)
    }
  } catch (error: any) {
    console.error('Get customer detail error:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('customers.message.getDetailError') + errorMessage)
    }
  } finally {
    loading.value = false
  }
}

// 取消操作
const handleCancel = () => {
  emit('cancel')
}

// 保存操作
const handleSave = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    // 准备API数据
    const requestData = {
      customer_name: formData.name,
      customer_type: formData.type as any,
      contact_person: formData.contact,
      email: formData.email || undefined,
      phone: formData.phone || undefined,
      address: formData.address || undefined,
      customer_level: formData.level as any,
      status: formData.status as any,
      company_size: formData.companySize as any || undefined,
      description: formData.description || undefined
    }
    
    // 调用API
    let response
    if (props.isEdit && props.customerId) {
      response = await updateCustomer(props.customerId, requestData)
    } else {
      response = await createCustomer(requestData)
    }
    
    if (response.code === '000000') {
      emit('save', { ...formData })
      ElMessage.success(response.message)
    } else {
      ElMessage.error(response.message)
    }
  } catch (error: any) {
    console.error('Save failed:', error)
    // 使用后端返回的错误信息
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(errorMessage)
    }
  } finally {
    loading.value = false
  }
}

// 组件挂载时加载枚举数据和客户数据
onMounted(async () => {
  await loadEnums()
  await loadCustomerData()
})
</script>

<style scoped>
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
.customer-form {
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

/* Flex布局 - 7个字段横向排列（方案3：混合布局） */
.fields-row-flex {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 20px;
}

/* 基础字段：按照Figma设计的325px宽度 */
.fields-row-flex .field-item {
  flex: 0 0 325px; /* 固定宽度325px */
  min-width: 280px; /* 移动端最小宽度 */
  margin-bottom: 0;
}

/* 状态字段：特殊处理，较窄宽度 */
.fields-row-flex .field-item.status-field {
  flex: 0 0 160px; /* 状态字段较窄，容纳两个单选按钮 */
  min-width: 160px;
}

/* Grid布局 - 地址跨列 */
.fields-row-grid {
  display: grid;
  grid-template-columns: 1fr 1fr 1fr;
  gap: 20px;
}

.address-field {
  grid-column: span 2; /* 跨两列，730px宽 */
  margin-bottom: 0;
}

/* 商业信息区块 */
.business-fields {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.company-size-field {
  width: 325px;
  margin-bottom: 0;
}

.description-field {
  width: 100%;
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
  color: #F56C6C;
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

:deep(.el-textarea__inner) {
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  font-size: 14px;
  color: var(--app-text-primary);
}

:deep(.el-textarea__inner::placeholder) {
  color: var(--app-text-secondary);
}

:deep(.el-radio-group) {
  display: flex;
  gap: 16px;
}

/* 状态字段单选按钮组特殊样式 - 最高优先级 */
.customer-form .status-inline :deep(.el-radio-group) {
  display: flex !important;
  flex-direction: row !important;
  gap: 16px !important;
  align-items: center !important;
}


:deep(.el-radio__label) {
  color: var(--app-text-primary);
  font-size: 14px;
}

:deep(.el-radio__input.is-checked .el-radio__inner) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

:deep(.el-radio__input.is-checked .el-radio__inner::after) {
  background-color: var(--app-content-bg);
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
  .fields-row-flex {
    flex-wrap: wrap;
  }
  
  .fields-row-flex .field-item {
    min-width: 200px;
  }
}

@media (max-width: 768px) {
  .top-section {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .fields-row-flex {
    flex-direction: column;
    gap: 16px;
  }
  
  .fields-row-grid {
    grid-template-columns: 1fr;
  }
  
  .address-field {
    grid-column: span 1;
  }
  
  .company-size-field {
    width: 100%;
  }
}
</style>