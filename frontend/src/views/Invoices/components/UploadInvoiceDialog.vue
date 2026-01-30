<template>
  <el-dialog
    v-model="visible"
    title="上传发票文件"
    width="680px"
    class="upload-invoice-dialog"
    :show-close="true"
    @close="handleClose"
  >
    <div class="upload-dialog-content">
      <!-- 信息摘要 -->
      <div class="info-summary" v-if="invoiceData">
        <div class="summary-item">
          <span class="label">发票申请号：</span>
          <span class="value">{{ invoiceData.invoiceNo }}</span>
        </div>
        <div class="summary-item">
          <span class="label">关联订单号：</span>
          <span class="value">{{ invoiceData.orderNo }}</span>
        </div>
        <div class="summary-item">
          <span class="label">用户信息：</span>
          <span class="value">{{ invoiceData.user }} ({{ invoiceData.userTypeLabel || '个人版' }})</span>
        </div>
        <div class="summary-item">
          <span class="label">发票抬头：</span>
          <span class="value">{{ invoiceData.invoiceTitle }}</span>
        </div>
        <div class="summary-item">
          <span class="label">开票金额：</span>
          <span class="value">¥{{ (invoiceData.amount || 0).toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</span>
        </div>
      </div>

      <!-- 上传区域 -->
      <div class="upload-section">
        <p class="section-label">上传发票文件</p>
        <el-upload
          class="invoice-uploader"
          drag
          action="#"
          :auto-upload="false"
          v-model:file-list="form.fileList"
        >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">
            <span class="upload-link">点击上传</span> / 拖拽到此区域
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持PDF格式。最大文件大小10MB
            </div>
          </template>
        </el-upload>
      </div>

      <!-- 表单项 -->
      <div class="form-section">
        <div class="form-item">
          <p class="section-label">开票完成时间</p>
          <el-date-picker
            v-model="form.finishTime"
            type="datetime"
            placeholder="选择日期时间"
            format="YYYY/MM/DD HH:mm"
            value-format="YYYY/MM/DD HH:mm"
            style="width: 100%"
          />
        </div>
        <div class="form-item">
          <p class="section-label">备注 (可选)</p>
          <el-input
            v-model="form.remark"
            type="textarea"
            :rows="3"
            placeholder="输入备注信息"
            maxlength="500"
            show-word-limit
          />
        </div>
      </div>
    </div>
    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" class="btn-submit" @click="handleSubmit">确认上传</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'

const props = defineProps<{
  modelValue: boolean
  invoiceData: any
}>()

const emit = defineEmits(['update:modelValue', 'submit'])

const visible = ref(props.modelValue)
const form = reactive({
  finishTime: '2026/01/14 06:20',
  remark: '',
  fileList: []
})

watch(() => props.modelValue, (val) => {
  visible.value = val
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

const handleClose = () => {
  // Reset form if needed
}

const handleSubmit = () => {
  ElMessage.success('发票上传成功')
  emit('submit', { ...form })
  visible.value = false
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
      .label { color: #8C8C8C; }
      .value { color: #262626; }
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
        &:hover { border-color: #019C7C; }
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
    &:hover { background-color: #017c63; border-color: #017c63; }
  }

  .el-button {
    min-width: 80px;
  }
}
</style>