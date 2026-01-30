<template>
  <el-dialog
    v-model="visible"
    title="驳回发票申请"
    width="640px"
    class="reject-invoice-dialog"
    :show-close="true"
    @close="handleClose"
  >
    <div class="reject-dialog-content">
      <div class="reject-alert">
        <el-icon class="alert-icon"><warning-filled /></el-icon>
        <span class="alert-text">驳回后，用户将收到通知并需要重新提交申请。</span>
      </div>

      <div class="info-summary" v-if="invoiceData">
        <div class="summary-item">
          <span class="label">发票申请号：</span>
          <span class="value">{{ invoiceData.invoiceNo }}</span>
        </div>
        <div class="summary-item">
          <span class="label">用户信息：</span>
          <span class="value">{{ invoiceData.user }}</span>
        </div>
        <div class="summary-item full-width">
          <span class="label">发票抬头：</span>
          <span class="value">{{ invoiceData.invoiceTitle }}</span>
        </div>
      </div>

      <el-form :model="form" ref="formRef" label-position="top">
        <div class="label-row">
          <span class="section-label required">驳回原因</span>
          <span class="label-tip">请选择或输入驳回的具体原因</span>
        </div>
        <el-form-item prop="reason" :rules="[{ required: true, message: '请选择驳回原因', trigger: 'change' }]">
          <el-select v-model="form.reason" placeholder="请选择驳回原因" style="width: 100%">
            <el-option label="发票抬头信息有误" value="title_error" />
            <el-option label="纳税人识别号不正确" value="tax_id_error" />
            <el-option label="附件不清晰或有误" value="attachment_error" />
            <el-option label="其他原因" value="other" />
          </el-select>
        </el-form-item>

        <div class="label-row">
          <span class="section-label">详细说明</span>
        </div>
        <el-form-item prop="remark">
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="4"
            placeholder="请输入详细的驳回说明，帮助用户更准确地修改信息"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">取消</el-button>
        <el-button type="danger" @click="handleSubmit">确认驳回</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { WarningFilled } from '@element-plus/icons-vue'

const props = defineProps<{
  modelValue: boolean
  invoiceData: any
}>()

const emit = defineEmits(['update:modelValue', 'submit'])

const visible = ref(props.modelValue)
const formRef = ref()
const form = reactive({
  reason: '',
  remark: ''
})

watch(() => props.modelValue, (val) => {
  visible.value = val
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleClose = () => {
  if (formRef.value) {
    formRef.value.resetFields()
  }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate((valid: boolean) => {
    if (valid) {
      ElMessage.success('发票申请已驳回')
      emit('submit', { ...form })
      visible.value = false
    }
  })
}
</script>

<style lang="scss" scoped>
@mixin dialog-style-base {
  border-radius: 4px;
  overflow: hidden;
  padding: 0;

  :deep(.el-dialog__header) {
    background-color: #019C7C;
    margin: 0;
    padding: 16px 24px;
    
    .el-dialog__title {
      color: #fff;
      font-size: 16px;
      font-weight: 500;
    }

    .el-dialog__close {
      color: #fff;
      &:hover { color: rgba(255, 255, 255, 0.8); }
    }
  }

  :deep(.el-dialog__body) {
    padding: 24px;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px 24px;
    border-top: 1px solid #f0f0f0;
  }
}

.reject-invoice-dialog {
  @include dialog-style-base;
  
  :deep(.el-form-item__label) {
    width: 100%;
    padding: 0;
  }
}

.reject-dialog-content {
  .reject-alert {
    background-color: #FFFBE6;
    border: 1px solid #FFE58F;
    border-radius: 4px;
    padding: 8px 16px;
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 24px;

    .alert-icon {
      color: #FAAD14;
      font-size: 16px;
    }

    .alert-text {
      color: #262626;
      font-size: 14px;
    }
  }

  .info-summary {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    padding: 20px;
    background: #FAFAFA;
    border-radius: 4px;
    border: 1px solid #F0F0F0;
    margin-bottom: 24px;

    .summary-item {
      font-size: 14px;
      .label { color: #8C8C8C; }
      .value { color: #262626; }
      &.full-width {
        grid-column: span 2;
      }
    }
  }

  .label-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    margin-bottom: 8px;

    .section-label {
      font-size: 14px;
      color: #262626;
      font-weight: 500;
      &.required::before {
        content: '*';
        color: #F5222D;
        margin-right: 4px;
      }
    }

    .label-tip {
      font-size: 12px;
      color: #BFBFBF;
      font-weight: normal;
    }
  }
}

.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 16px;

  .el-button {
    min-width: 80px;
  }
}
</style>