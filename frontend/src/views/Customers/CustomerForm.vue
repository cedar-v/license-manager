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
      <span class="breadcrumb-item">客户管理</span>
      <span class="breadcrumb-separator">></span>
      <span class="breadcrumb-current">{{ isEdit ? '编辑客户' : '添加新客户' }}</span>
    </div>
    
    <!-- 操作按钮 -->
    <div class="form-actions">
      <el-button @click="handleCancel">取 消</el-button>
      <el-button type="primary" @click="handleSave" :loading="loading">保 存</el-button>
    </div>
  </div>

  <!-- 主表单区域 -->
  <div >
    <el-form :model="formData" :rules="formRules" ref="formRef" label-position="top">
      <!-- 基本信息 -->
      <div  class="customer-form">
        <h3 class="section-title">基本信息</h3>
        
        <!-- Flex布局 - 7个字段横向排列 -->
        <div class="fields-row-flex">
          <el-form-item label="客户名称" prop="name" required class="field-item">
            <el-input v-model="formData.name" placeholder="请输入" />
          </el-form-item>
          
          <el-form-item label="客户类型" prop="type" required class="field-item">
            <el-select v-model="formData.type" placeholder="请选择" style="width: 100%">
              <el-option label="企业客户" value="enterprise" />
              <el-option label="个人客户" value="individual" />
              <el-option label="政府机构" value="government" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="客户等级" prop="level" required class="field-item">
            <el-select v-model="formData.level" placeholder="请选择" style="width: 100%">
              <el-option label="VIP客户" value="vip" />
              <el-option label="重要客户" value="important" />
              <el-option label="普通客户" value="normal" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="联系人" prop="contact" required class="field-item">
            <el-input v-model="formData.contact" placeholder="请输入" />
          </el-form-item>
          
          <el-form-item label="邮箱" prop="email" class="field-item">
            <el-input v-model="formData.email" placeholder="请输入" />
          </el-form-item>
          
          <el-form-item label="状态" prop="status" required class="field-item status-field">
            <el-radio-group v-model="formData.status">
              <el-radio value="normal">正常</el-radio>
              <el-radio value="disabled">禁用</el-radio>
            </el-radio-group>
          </el-form-item>
          
          <el-form-item label="联系电话" prop="phone" class="field-item">
            <el-input v-model="formData.phone" placeholder="请输入" />
          </el-form-item>
        </div>
        
        <!-- Grid布局 - 地址跨列 -->
        <div class="fields-row-grid">
          <el-form-item label="地址" prop="address" class="address-field">
            <el-input v-model="formData.address" placeholder="请输入" />
          </el-form-item>
        </div>
      </div>

      <!-- 商业信息 -->
      <div  class="customer-form">
        <h3 class="section-title">商业信息</h3>
        
        <div class="business-fields">
          <el-form-item label="企业规模" prop="companySize" required class="company-size-field">
            <el-select v-model="formData.companySize" placeholder="请选择" style="width: 100%">
              <el-option label="大型企业(500人以上)" value="large" />
              <el-option label="中型企业(100-500人)" value="medium" />
              <el-option label="小型企业(50-100人)" value="small" />
              <el-option label="微型企业(50人以下)" value="micro" />
            </el-select>
          </el-form-item>
          
          <el-form-item label="客户描述" prop="description" class="description-field">
            <el-input
              v-model="formData.description"
              type="textarea"
              :rows="4"
              placeholder="请输入"
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
import { ref, reactive, computed } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'

interface CustomerFormData {
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
  customerData?: Partial<CustomerFormData>
  isEdit?: boolean
}>()

const emit = defineEmits<{
  save: [data: CustomerFormData]
  cancel: []
}>()

const formRef = ref<FormInstance>()
const loading = ref(false)

// 表单数据
const formData = reactive<CustomerFormData>({
  name: '',
  contact: '',
  email: '',
  phone: '',
  address: '',
  type: '',
  level: '',
  status: 'normal',
  companySize: '',
  description: ''
})

// 表单验证规则
const formRules: FormRules = {
  name: [
    { required: true, message: '请输入客户名称', trigger: 'blur' }
  ],
  contact: [
    { required: true, message: '请输入联系人', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择客户类型', trigger: 'change' }
  ],
  level: [
    { required: true, message: '请选择客户等级', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ],
  companySize: [
    { required: true, message: '请选择企业规模', trigger: 'change' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  phone: [
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ]
}

// 初始化表单数据
if (props.customerData) {
  Object.assign(formData, props.customerData)
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
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    emit('save', { ...formData })
    ElMessage.success(props.isEdit ? '编辑成功' : '添加成功')
  } catch (error) {
    console.error('表单验证失败:', error)
  } finally {
    loading.value = false
  }
}
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
  color: #888;
}

.breadcrumb-item {
  color: #888;
}

.breadcrumb-separator {
  margin: 0 8px;
  color: #b6bad2;
}

.breadcrumb-current {
  color: #1d1d1d;
}

.form-actions {
  display: flex;
  gap: 20px;
}

/* 主表单区域 */
.customer-form {
  padding: 20px;
  background: #fff;
  border-radius: 4px;
  width: 100%;
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  color: #1d1d1d;
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
  color: #1d1d1d;
  font-size: 14px;
  line-height: 1.5714285714285714em;
  margin-bottom: 8px;
}

:deep(.el-form-item.is-required .el-form-item__label::before) {
  content: '*';
  color: #eb3737;
  margin-right: 4px;
}

:deep(.el-input__wrapper) {
  border: 1px solid #c2c6ce;
  border-radius: 4px;
  height: 32px;
}

:deep(.el-input__inner) {
  color: #1d1d1d;
  font-size: 14px;
  line-height: 1.5em;
}

:deep(.el-input__inner::placeholder) {
  color: #b2b8c2;
  font-size: 14px;
}

:deep(.el-select .el-input__wrapper) {
  border: 1px solid #c2c6ce;
  height: 32px;
}

:deep(.el-textarea__inner) {
  border: 1px solid #c2c6ce;
  border-radius: 4px;
  font-size: 14px;
  color: #1d1d1d;
}

:deep(.el-textarea__inner::placeholder) {
  color: #b2b8c2;
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
  color: #1d1d1d;
  font-size: 14px;
}

:deep(.el-radio__input.is-checked .el-radio__inner) {
  background-color: #00c27c;
  border-color: #00c27c;
}

:deep(.el-radio__input.is-checked .el-radio__inner::after) {
  background-color: #ffffff;
}

:deep(.el-button--primary) {
  background-color: #019c7c;
  border-color: #019c7c;
  font-size: 14px;
  padding: 6px 16px;
}

:deep(.el-button--primary:hover) {
  background-color: #00c27c;
  border-color: #00c27c;
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