<template>
  <el-dialog
    v-model="visible"
    :title="t('invoices.dialog.rejectTitle')"
    width="640px"
    class="reject-invoice-dialog"
    :show-close="true"
    @close="handleClose"
  >
    <div class="reject-dialog-content">
      <div class="reject-alert">
        <el-icon class="alert-icon"><warning-filled /></el-icon>
        <span class="alert-text">{{ t('invoices.dialog.rejectAlert') }}</span>
      </div>

      <div class="info-summary" v-if="invoiceData">
        <div class="summary-item">
          <span class="label">{{ t('invoices.detail.invoiceNo') }}：</span>
          <span class="value">{{ invoiceData.invoice_no }}</span>
        </div>
        <div class="summary-item">
          <span class="label">{{ t('invoices.detail.userInfo') }}：</span>
          <span class="value">{{ invoiceData.applicant_name }}</span>
        </div>
        <div class="summary-item full-width">
          <span class="label">{{ t('invoices.detail.invoiceTitle') }}：</span>
          <span class="value">{{ invoiceData.title }}</span>
        </div>
      </div>

      <el-form :model="form" ref="formRef" label-position="top">
        <div class="label-row">
          <span class="section-label required">{{ t('invoices.dialog.reasonLabel') }}</span>
          <span class="label-tip">{{ t('invoices.dialog.reasonTip') }}</span>
        </div>
        <el-form-item prop="reason" :rules="[{ required: true, message: t('invoices.dialog.reasonPlaceholder'), trigger: 'change' }]">
          <el-select v-model="form.reason" :placeholder="t('invoices.dialog.reasonPlaceholder')" style="width: 100%">
            <el-option :label="t('invoices.dialog.reasons.title_error')" value="title_error" />
            <el-option :label="t('invoices.dialog.reasons.tax_id_error')" value="tax_id_error" />
            <el-option :label="t('invoices.dialog.reasons.attachment_error')" value="attachment_error" />
            <el-option :label="t('invoices.dialog.reasons.other')" value="other" />
          </el-select>
        </el-form-item>

        <div class="label-row">
          <span class="section-label">{{ t('invoices.dialog.remarkLabel') }}</span>
        </div>
        <el-form-item prop="suggestion">
          <el-input
            v-model="form.suggestion"
            type="textarea"
            :rows="4"
            :placeholder="t('invoices.dialog.remarkPlaceholder')"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">{{ t('invoices.actions.cancel') }}</el-button>
        <el-button type="danger" @click="handleSubmit">{{ t('invoices.actions.confirmReject') }}</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { WarningFilled } from '@element-plus/icons-vue'

const { t } = useI18n()
const props = defineProps<{
  modelValue: boolean
  invoiceData: any
}>()

const emit = defineEmits(['update:modelValue', 'submit'])

const visible = ref(props.modelValue)
const formRef = ref()
const form = reactive({
  reason: '',
  suggestion: ''
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