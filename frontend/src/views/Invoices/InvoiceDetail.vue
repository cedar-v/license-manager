<template>
  <Layout :page-title="t('invoices.detailTitle')">
    <div class="invoice-detail-container" v-loading="loading">
      <div class="detail-header">
        <div class="header-left">
          <el-button link @click="handleBack">
            <el-icon><ArrowLeft /></el-icon>
            {{ t('invoices.title') }}
          </el-button>
        </div>
        <div class="header-actions">
          <el-button v-if="detailData.status === 'pending'" class="btn-upload" type="primary" @click="handleUpload">{{ t('invoices.actions.upload') }}</el-button>
          <el-button v-if="detailData.status === 'pending'" class="btn-reject" type="danger" @click="handleReject">{{ t('invoices.actions.reject') }}</el-button>
          <el-button @click="handleBack">{{ t('invoices.actions.back') }}</el-button>
        </div>
      </div>

      <div class="detail-content">
        <!-- 基本信息 -->
        <div class="detail-card">
          <div class="card-title">
            <span class="title-line"></span>
            <h3>{{ t('invoices.detail.basicInfo') }}</h3>
          </div>
          <el-descriptions :column="4" class="info-descriptions">
            <el-descriptions-item :label="t('invoices.detail.invoiceNo')">{{ detailData.invoice_no }}</el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.orderNo')">{{ detailData.order_no }}</el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.applyTime')">{{ detailData.created_at }}</el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.status')">
              <span :class="['status-tag', detailData.status]">
                {{ detailData.status_display }}
              </span>
            </el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.userInfo')">
              {{ detailData.applicant_name }} <span class="user-type">({{ detailData.order_package_name }})</span>
            </el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.phone')">{{ detailData.applicant_phone }}</el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.amount')">
              <span class="amount-value">¥{{ (detailData.amount || 0).toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</span>
            </el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.email')">{{ detailData.receiver_email }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 开票信息 -->
        <div class="detail-card">
          <div class="card-title">
            <span class="title-line"></span>
            <h3>{{ t('invoices.detail.invoiceInfo') }}</h3>
          </div>
          <el-descriptions :column="4" class="info-descriptions">
            <el-descriptions-item :label="t('invoices.detail.invoiceType')">{{ detailData.invoice_type_display }}</el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.invoiceTitle')">{{ detailData.title }}</el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.taxId')">{{ detailData.taxpayer_id }}</el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.content')">{{ detailData.content }}</el-descriptions-item>
            <el-descriptions-item :label="t('invoices.detail.remark')" :span="4">{{ detailData.remark }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 驳回信息 (仅在状态为驳回时显示) -->
        <div class="detail-card reject-info-card" v-if="detailData.status === 'rejected'">
          <div class="card-title">
            <span class="title-line"></span>
            <h3>{{ t('invoices.detail.rejectInfo') }}</h3>
          </div>
          <div class="info-content">
            <div class="info-row">
              <span class="label">{{ t('invoices.detail.rejectReason') }}：</span>
              <span class="value">{{ detailData.reject_reason }}</span>
            </div>
            <div class="info-row">
              <span class="label">{{ t('invoices.detail.rejectTime') }}：</span>
              <span class="value">{{ detailData.rejected_at }}</span>
            </div>
            <div class="info-row">
              <span class="label">{{ t('invoices.detail.rejectUser') }}：</span>
              <span class="value">{{ detailData.rejected_by }}</span>
            </div>
          </div>
        </div>

        <!-- 通过信息 (仅在状态为已开票时显示) -->
        <div class="detail-card approve-info-card" v-if="detailData.status === 'completed'">
          <div class="card-title">
            <span class="title-line"></span>
            <h3>{{ t('invoices.detail.approveInfo') }}</h3>
          </div>
          <div class="info-content">
            <div class="info-row">
              <span class="label">{{ t('invoices.detail.finishTime') }}：</span>
              <span class="value">{{ detailData.uploaded_at }}</span>
            </div>
            <div class="info-row">
              <span class="label">{{ t('invoices.detail.approveUser') }}：</span>
              <span class="value">{{ detailData.uploader_name }}</span>
            </div>
            <div class="info-row">
              <span class="label">{{ t('invoices.detail.fileName') }}：</span>
              <span class="value">
                <el-link v-if="detailData.invoice_file_url" type="primary" :underline="false" :href="detailData.invoice_file_url" target="_blank">
                  {{ detailData.invoice_file_url.split('/').pop() }}
                </el-link>
              </span>
            </div>
            <div class="info-row">
              <span class="label">{{ t('invoices.detail.approveTime') }}：</span>
              <span class="value">{{ detailData.uploaded_at }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 上传发票对话框 -->
      <UploadInvoiceDialog
        v-model="uploadVisible"
        :invoice-data="detailData"
        @submit="handleUploadSubmit"
      />

      <!-- 驳回申请对话框 -->
      <RejectInvoiceDialog
        v-model="rejectVisible"
        :invoice-data="detailData"
        @submit="handleRejectSubmit"
      />
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { ArrowLeft } from '@element-plus/icons-vue'
import Layout from '@/components/common/layout/Layout.vue'
import UploadInvoiceDialog from './components/UploadInvoiceDialog.vue'
import RejectInvoiceDialog from './components/RejectInvoiceDialog.vue'
import { getInvoiceDetail, uploadInvoice, rejectInvoice, type Invoice } from '@/api/invoice'
import { formatDateTime } from '@/utils/date'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()

const uploadVisible = ref(false)
const rejectVisible = ref(false)
const loading = ref(false)

const detailData = ref<Partial<Invoice>>({
  invoice_no: '',
  order_no: '',
  created_at: '',
  status: 'pending',
  status_display: t('invoices.filter.pending'),
  applicant_name: '',
  order_package_name: '',
  applicant_phone: '',
  amount: 0,
  receiver_email: '',
  invoice_type_display: '',
  title: '',
  taxpayer_id: '',
  content: '',
  remark: ''
})

const fetchDetail = async () => {
  const id = route.params.id as string
  if (!id) return
  loading.value = true
  try {
    const res = await getInvoiceDetail(id)
    if (res.code === '000000' && res.data) {
      const data = res.data
      data.created_at = formatDateTime(data.created_at)
      if (data.rejected_at) data.rejected_at = formatDateTime(data.rejected_at)
      if (data.uploaded_at) data.uploaded_at = formatDateTime(data.uploaded_at)
      detailData.value = data
    }
  } catch (error: any) {
    console.error('Fetch invoice detail error:', error)
    ElMessage.error(error.backendMessage || t('invoices.messages.fetchDetailError'))
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchDetail()
})

const handleBack = () => {
  router.push('/invoices')
}

const handleUpload = () => {
  uploadVisible.value = true
}

const handleUploadSubmit = async (data: any) => {
  console.log('handleUploadSubmit detail triggered:', data)
  if (!detailData.value.invoice_no) {
    console.error('No detail data invoice_no')
    return
  }
  try {
    const formData = new FormData()
    formData.append('invoice_no', detailData.value.id!)
    if (data.fileList && data.fileList.length > 0) {
      const file = data.fileList[0].raw
      console.log('File to upload (detail):', file)
      formData.append('file', file)
    } else {
      ElMessage.warning(t('invoices.dialog.uploadTip'))
      return
    }

    const res = await uploadInvoice(formData)
    console.log('Upload response (detail):', res)
    if (res.code === '000000') {
      ElMessage.success(t('invoices.messages.uploadSuccess'))
      fetchDetail()
    }
  } catch (error: any) {
    console.error('Upload error detail:', error)
    ElMessage.error(error.backendMessage || t('invoices.messages.fetchDetailError'))
  }
}

const handleReject = () => {
  rejectVisible.value = true
}

const handleRejectSubmit = async (data: any) => {
  if (!detailData.value.id) return
  try {
    const reasonText = data.reason === 'other' ? data.suggestion : t(`invoices.dialog.reasons.${data.reason}`)
    const res = await rejectInvoice(detailData.value.id, {
      reject_reason: reasonText,
      suggestion: data.suggestion
    })
    if (res.code === '000000') {
      ElMessage.success(t('invoices.messages.rejectSuccess'))
      fetchDetail()
    }
  } catch (error: any) {
    console.error('Reject error detail:', error)
    ElMessage.error(error.backendMessage || t('invoices.messages.fetchDetailError'))
  }
}
</script>

<style lang="scss" scoped>
.invoice-detail-container {
  padding: 24px;
  background-color: var(--app-bg-color);
  min-height: calc(100vh - 80px);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: var(--app-content-bg);
  padding: 16px 24px;
  border-radius: 4px;
  border: 1px solid var(--app-border-light);
  box-shadow: var(--app-shadow);

  .header-left {
    display: flex;
    align-items: center;

    :deep(.el-button) {
      font-size: 14px;
      color: var(--app-text-secondary);
      &:hover {
        color: var(--el-color-primary);
      }
      .el-icon {
        margin-right: 4px;
      }
    }
  }
}

.header-actions {
  display: flex;
  gap: 12px;

  .btn-upload {
    background-color: var(--el-color-primary);
    border-color: var(--el-color-primary);
    &:hover { background-color: var(--el-color-primary-dark-2); border-color: var(--el-color-primary-dark-2); }
  }

  .btn-reject {
    background-color: var(--el-color-danger);
    border-color: var(--el-color-danger);
    &:hover { background-color: var(--el-color-danger-dark-2); border-color: var(--el-color-danger-dark-2); }
  }
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-card {
  background: var(--app-content-bg);
  padding: 24px;
  border-radius: 4px;
  border: 1px solid var(--app-border-light);
  box-shadow: var(--app-shadow);

  .card-title {
    display: flex;
    align-items: center;
    margin-bottom: 24px;

    .title-line {
      width: 4px;
      height: 16px;
      background: var(--el-color-primary);
      margin-right: 12px;
    }

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: var(--app-text-primary);
    }
  }

  .info-content {
    .info-row {
      margin-bottom: 12px;
      font-size: 14px;
      line-height: 22px;
      display: flex;
      
      &:last-child {
        margin-bottom: 0;
      }

      .label {
        color: var(--app-text-secondary);
        width: 120px;
        flex-shrink: 0;
      }

      .value {
        color: var(--app-text-primary);
      }
    }
  }
}

:deep(.info-descriptions) {
  .el-descriptions__label {
    width: 120px;
    color: var(--app-text-secondary);
    font-weight: normal;
  }
  .el-descriptions__content {
    color: var(--app-text-primary);
  }
}

.status-tag {
  padding: 2px 8px;
  border-radius: 2px;
  font-size: 12px;
  &.pending { background: var(--el-color-warning-light-9); color: var(--el-color-warning); border: 1px solid var(--el-color-warning-light-5); }
  &.completed { background: var(--el-color-success-light-9); color: var(--el-color-success); border: 1px solid var(--el-color-success-light-5); }
  &.rejected { background: var(--el-color-danger-light-9); color: var(--el-color-danger); border: 1px solid var(--el-color-danger-light-5); }
}

.amount-value {
  font-weight: 600;
}

.user-type {
  color: var(--app-text-secondary);
}

// 弹窗通用基础样式
@mixin dialog-style-base {
  border-radius: 4px;
  overflow: hidden;
  padding: 0;

  .el-dialog__header {
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

  .el-dialog__body {
    padding: 24px;
  }

  .el-dialog__footer {
    padding: 16px 24px 24px;
    border-top: 1px solid #f0f0f0;
  }
}

// 上传弹窗样式还原
:deep(.upload-invoice-dialog) {
  @include dialog-style-base;
}

// 驳回弹窗样式
:deep(.reject-invoice-dialog) {
  @include dialog-style-base;

  .el-form-item__label {
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
