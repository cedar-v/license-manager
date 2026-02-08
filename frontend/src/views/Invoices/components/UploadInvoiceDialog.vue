<template>
  <div>
    <el-dialog v-model="visible" :title="t('invoices.dialog.uploadTitle')" width="680px" class="upload-invoice-dialog"
      :show-close="true" @close="handleClose">
      <div class="upload-dialog-content">
        <!-- 信息摘要 -->
        <div class="info-summary" v-if="invoiceData">
          <div class="summary-item">
            <span class="label">{{ t('invoices.detail.invoiceNo') }}：</span>
            <span class="value">{{ invoiceData.invoice_no }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ t('invoices.detail.orderNo') }}：</span>
            <span class="value">{{ invoiceData.order_no }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ t('invoices.detail.userInfo') }}：</span>
            <span class="value">{{ invoiceData.applicant_name }} ({{ invoiceData.order_package_name || '个人版' }})</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ t('invoices.detail.invoiceTitle') }}：</span>
            <span class="value">{{ invoiceData.title }}</span>
          </div>
          <div class="summary-item">
            <span class="label">{{ t('invoices.detail.amount') }}：</span>
            <span class="value">¥{{ (invoiceData.amount || 0).toLocaleString(undefined, { minimumFractionDigits: 2 })
            }}</span>
          </div>
        </div>

        <!-- 上传区域 -->
        <div class="upload-section">
          <p class="section-label">{{ t('invoices.dialog.uploadLabel') }}</p>
          <el-upload class="invoice-uploader" drag action="#" :auto-upload="false" v-model:file-list="form.fileList">
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              <span class="upload-link">{{ t('invoices.dialog.uploadDraggerText').split(' / ')[0] }}</span> / {{
                t('invoices.dialog.uploadDraggerText').split(' / ')[1] }}
            </div>
            <template #tip>
              <div class="el-upload__tip">
                {{ t('invoices.dialog.uploadTip') }}
              </div>
            </template>
          </el-upload>
        </div>

        <!-- 表单项 -->
        <div class="form-section">
          <div class="form-item">
            <p class="section-label">{{ t('invoices.dialog.finishTimeLabel') }}</p>
            <el-date-picker v-model="form.issued_at" type="datetime"
              :placeholder="t('invoices.dialog.finishTimePlaceholder')" style="width: 100%" format="YYYY-MM-DD HH:mm:ss"
              value-format="YYYY-MM-DD HH:mm:ss" />
          </div>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="visible = false">{{ t('invoices.actions.cancel') }}</el-button>
          <el-button type="primary" class="btn-submit" @click="handleSubmit">{{ t('invoices.actions.confirmUpload')
          }}</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'

const { t } = useI18n()
const props = defineProps<{
  modelValue: boolean
  invoiceData: any
}>()

const emit = defineEmits(['update:modelValue', 'submit'])

const visible = ref(props.modelValue)
const form = reactive({
  fileList: [] as any[],
  issued_at:""
})

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val) {
    // Reset form when opening
    form.fileList = []
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})
const handleClose = () => {
  // Reset form if needed
}

const handleSubmit = () => {
  console.log('handleSubmit clicked, fileList:', form.fileList)
  if (form.fileList.length === 0) {
    ElMessage.warning(t('invoices.dialog.uploadTip'))
    return
  }
  if (!form.issued_at) {
    ElMessage.warning(t('invoices.dialog.finishTimePlaceholder'))
    return
  }
  emit('submit', { ...form })
  visible.value = false
}
</script>

<style lang="scss" scoped>
:deep(.el-dialog__headerbtn) {
  top: 10px !important;
}

:deep(.el-dialog) {
  border-radius: 8px;
  overflow: hidden;
  padding: 0;
  padding-bottom: 20px;
  box-sizing: border-box;
}

:deep(.el-dialog__header) {
  background: linear-gradient(90deg, #00928A 0%, #00D19E 100%) !important;
  margin: 0;
  padding: 16px 24px;

  .el-dialog__title {
    color: #fff;
    font-size: 16px;
    font-weight: 500;
  }

  .el-dialog__close {
    color: #fff;

    &:hover {
      color: rgba(255, 255, 255, 0.8);
    }
  }
}

@mixin dialog-style-base {
  border-radius: 4px;
  overflow: hidden;
  padding: 0;



  :deep(.el-dialog__body) {
    padding: 24px;
  }

  :deep(.el-dialog__footer) {
    padding: 16px 24px 24px;
    border-top: 1px solid #f0f0f0;
  }
}

.upload-invoice-dialog {
  @include dialog-style-base;
}

.upload-dialog-content {
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

      .label {
        color: #8C8C8C;
      }

      .value {
        color: #262626;
      }
    }
  }

  .section-label {
    font-size: 14px;
    color: #262626;
    margin-bottom: 8px;
    font-weight: 500;
  }

  .upload-section {
    margin-bottom: 24px;

    .invoice-uploader {
      :deep(.el-upload-dragger) {
        border: 1px dashed #D9D9D9;
        padding: 30px;

        &:hover {
          border-color: #019C7C;
        }
      }

      .upload-link {
        color: #019C7C;
        text-decoration: underline;
      }

      .el-upload__tip {
        text-align: center;
        margin-top: 8px;
        color: #BFBFBF;
      }
    }
  }

  .form-section {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }
}

.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 16px;

  .btn-submit {
    background-color: #019C7C;
    border-color: #019C7C;
    padding: 8px 32px;

    &:hover {
      background-color: #017c63;
      border-color: #017c63;
    }
  }

  .el-button {
    min-width: 80px;
  }
}
</style>