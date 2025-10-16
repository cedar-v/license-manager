<template>
  <div class="license-form-container">
    <!-- 顶部横向区域 - 与表单平级 -->
    <div class="top-section">
      <!-- 面包屑导航 -->
      <div class="breadcrumb">
        <span class="breadcrumb-item">授权管理</span>
        <span class="breadcrumb-separator">></span>
        <span class="breadcrumb-current">{{ isEdit ? '编辑授权' : '创建授权' }}</span>
      </div>
      
      <!-- 操作按钮 -->
      <div class="form-actions">
        <el-button @click="handleCancel">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{ isEdit ? '更新' : '创建' }}</el-button>
      </div>
    </div>

    <!-- 主表单区域 -->
    <div>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-position="top">
        <!-- 基本信息 -->
        <div class="license-form">
          <h3 class="section-title">基本信息</h3>
          
          <!-- 第一行：关联客户名称，客户ID，创建人 -->
          <div class="fields-row">
            <el-form-item label="关联客户名称" prop="customer_id" required class="field-item">
              <el-select
                v-model="formData.customer_id"
                placeholder="请选择客户"
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
            
            <el-form-item label="客户ID" prop="customer_code" required class="field-item">
              <el-input
                v-model="customerCode"
                placeholder="自动生成或手动输入"
                readonly
              />
            </el-form-item>
            
            <el-form-item label="创建人" prop="created_by" required class="field-item">
              <el-input
                v-model="createdBy"
                placeholder="当前登录用户"
                readonly
              />
            </el-form-item>
          </div>
          
          <!-- 第二行：备注单独一行 -->
          <div class="fields-row">
            <el-form-item label="备注" prop="description" required class="field-item-full">
              <el-input
                v-model="formData.description"
                placeholder="请输入备注信息"
              />
            </el-form-item>
          </div>
        </div>

        <!-- 授权配置 -->
        <div class="license-form">
          <h3 class="section-title">授权配置</h3>
          
          <!-- 第三行：授权期限，起止时间 -->
          <div class="fields-row">
            <el-form-item label="授权期限" prop="validity_days" required class="field-item">
              <el-input-number
                v-model="formData.validity_days"
                :min="1"
                :max="36500"
                placeholder="请输入授权期限（天）"
                style="width: 100%"
              />
            </el-form-item>
            
            <el-form-item label="起止时间" prop="date_range" required class="field-item">
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                range-separator="至"
                start-placeholder="开始日期"
                end-placeholder="结束日期"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                @change="handleDateRangeChange"
              />
            </el-form-item>
          </div>
          
          <!-- 第四行：最大激活设备数，部署类型，加密类型 -->
          <div class="fields-row">
            <el-form-item label="最大激活设备数" prop="max_activations" required class="field-item">
              <el-input-number
                v-model="formData.max_activations"
                :min="1"
                placeholder="请输入最大激活设备数"
                style="width: 100%"
              />
            </el-form-item>
            
            <el-form-item label="部署类型" prop="deployment_type" required class="field-item">
              <el-select
                v-model="formData.deployment_type"
                placeholder="请选择部署类型"
                style="width: 100%"
              >
                <el-option
                  v-for="option in deploymentTypeOptions"
                  :key="option.key"
                  :label="option.display"
                  :value="option.key"
                />
              </el-select>
            </el-form-item>
            
            <el-form-item label="加密类型" prop="encryption_type" required class="field-item">
              <el-select
                v-model="formData.encryption_type"
                placeholder="请选择加密类型"
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
import { createLicense, updateLicense, getLicenseDetail, type AuthorizationCodeCreateRequest, type LicenseUpdateRequest } from '@/api/license'
import { getCustomers, type Customer } from '@/api/customer'
import { getEnumOptions, type RawEnumItem } from '@/api/enum'

const router = useRouter()
const route = useRoute()

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
const formData = reactive<AuthorizationCodeCreateRequest>({
  customer_id: '',
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
const customerCode = ref('')
const createdBy = ref('')
const dateRange = ref<[string, string] | null>(null)

// 计算属性
const isEdit = computed(() => {
  return route.name === 'licenses-edit' && route.params.id
})

// 表单验证规则
const formRules: FormRules = {
  customer_id: [
    { required: true, message: '请选择客户', trigger: 'change' }
  ],
  customer_code: [
    { required: true, message: '客户ID不能为空', trigger: 'blur' }
  ],
  created_by: [
    { required: true, message: '创建人不能为空', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入备注信息', trigger: 'blur' }
  ],
  validity_days: [
    { required: true, message: '请输入授权期限', trigger: 'blur' },
    { type: 'number', min: 1, max: 36500, message: '授权期限必须在1-36500天之间', trigger: 'blur' }
  ],
  date_range: [
    { required: true, message: '请选择起止时间', trigger: 'change' }
  ],
  deployment_type: [
    { required: true, message: '请选择部署类型', trigger: 'change' }
  ],
  encryption_type: [
    { required: true, message: '请选择加密类型', trigger: 'change' }
  ],
  max_activations: [
    { required: true, message: '请输入最大激活设备数', trigger: 'blur' },
    { type: 'number', min: 1, message: '最大激活设备数必须大于0', trigger: 'blur' }
  ]
}

// 方法
const loadEnumOptions = async () => {
  try {
    // 加载部署类型选项
    const deploymentResponse = await getEnumOptions('deployment_type')
    deploymentTypeOptions.value = deploymentResponse.data.items

    // 加载加密类型选项
    const encryptionResponse = await getEnumOptions('encryption_type')
    encryptionTypeOptions.value = encryptionResponse.data.items
  } catch (error) {
    console.error('加载枚举选项失败:', error)
    // 使用默认选项
    deploymentTypeOptions.value = [
      { key: 'standalone', display: '单机版' },
      { key: 'cloud', display: '云端版' },
      { key: 'hybrid', display: '混合版' }
    ]
    encryptionTypeOptions.value = [
      { key: 'standard', display: '标准加密' },
      { key: 'advanced', display: '高级加密' }
    ]
  }
}

const searchCustomers = async (query: string) => {
  if (!query) {
    customerOptions.value = []
    return
  }
  
  try {
    customerLoading.value = true
    const response = await getCustomers({
      page: 1,
      page_size: 20,
      customer_name: query
    })
    customerOptions.value = response.data.list
  } catch (error) {
    console.error('搜索客户失败:', error)
    ElMessage.error('搜索客户失败')
  } finally {
    customerLoading.value = false
  }
}

const loadCustomerInfo = () => {
  // 从路由参数获取客户信息
  const customerId = route.query.customerId as string
  const customerName = route.query.customerName as string
  
  if (customerId && customerName) {
    formData.customer_id = customerId
    customerOptions.value = [{
      id: customerId,
      customer_name: customerName,
      customer_code: '',
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
    } as Customer]
  }
  
  // 设置创建人（这里可以从用户store获取当前登录用户）
  createdBy.value = '当前用户' // TODO: 从用户store获取
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

const loadLicenseDetail = async () => {
  if (!isEdit.value) return
  
  try {
    const id = route.params.id as string
    const response = await getLicenseDetail(id)
    const data = response.data
    
    if (data) {
      // 填充表单数据
      formData.customer_id = data.customer_id || ''
      formData.description = data.description || ''
      formData.deployment_type = data.deployment_type || 'standalone'
      formData.encryption_type = data.encryption_type || 'standard'
      formData.max_activations = data.max_activations || 1
      
      // 设置独立字段
      customerCode.value = (data as any).customer_code || data.customer_id || ''
      createdBy.value = (data as any).created_by || '当前用户'
      
      // 计算有效期天数和设置日期范围
      if (data.start_date && data.end_date) {
        const startDate = new Date(data.start_date)
        const endDate = new Date(data.end_date)
        const diffTime = endDate.getTime() - startDate.getTime()
        formData.validity_days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
        dateRange.value = [data.start_date, data.end_date]
      }
      
      // 设置客户选项
      if (data.customer_name) {
        customerOptions.value = [{
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
        } as Customer]
      }
    }
  } catch (error) {
    console.error('加载授权详情失败:', error)
    ElMessage.error('加载授权详情失败')
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
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
        reason: '管理员更新授权配置'
      }
      
      await updateLicense(id, updateData)
      ElMessage.success('更新成功')
    } else {
      // 创建授权
      await createLicense(formData)
      ElMessage.success('创建成功')
    }
    
    // 返回列表页面
    router.back()
  } catch (error) {
    console.error('提交失败:', error)
    ElMessage.error('提交失败')
  } finally {
    submitting.value = false
  }
}

const handleCancel = () => {
  ElMessageBox.confirm('确定要取消吗？未保存的更改将丢失。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '继续编辑',
    type: 'warning'
  }).then(() => {
    router.back()
  }).catch(() => {
    // 用户选择继续编辑
  })
}

// 监听客户选择变化，自动填充客户ID
watch(() => formData.customer_id, (newCustomerId) => {
  if (newCustomerId) {
    const selectedCustomer = customerOptions.value.find(c => c.id === newCustomerId)
    if (selectedCustomer) {
      customerCode.value = selectedCustomer.customer_code || selectedCustomer.id
    }
  }
})

// 生命周期
onMounted(async () => {
  await loadEnumOptions()
  loadCustomerInfo()
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