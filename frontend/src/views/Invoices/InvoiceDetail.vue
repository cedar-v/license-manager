<template>
  <Layout>
    <div class="invoice-detail-container" v-loading="loading">
      <div class="detail-header">
        <div class="header-left">
          <span class="back-link" @click="handleBack">发票管理</span>
          <span class="separator">/</span>
          <span class="current-page">发票申请详情</span>
        </div>
        <div class="header-actions">
          <el-button v-if="detailData.status === 'pending'" class="btn-upload" type="primary" @click="handleUpload">上传发票</el-button>
          <el-button v-if="detailData.status === 'pending'" class="btn-reject" type="danger" @click="handleReject">驳回申请</el-button>
          <el-button @click="handleBack">返回</el-button>
        </div>
      </div>

      <div class="detail-content">
        <!-- 基本信息 -->
        <div class="detail-card">
          <div class="card-title">
            <span class="title-line"></span>
            <h3>基本信息</h3>
          </div>
          <el-descriptions :column="4" class="info-descriptions">
            <el-descriptions-item label="发票申请号">{{ detailData.invoiceNo }}</el-descriptions-item>
            <el-descriptions-item label="关联订单号">{{ detailData.orderNo }}</el-descriptions-item>
            <el-descriptions-item label="申请时间">{{ detailData.applyTime }}</el-descriptions-item>
            <el-descriptions-item label="开票状态">
              <span :class="['status-tag', detailData.status]">
                {{ detailData.statusLabel }}
              </span>
            </el-descriptions-item>
            <el-descriptions-item label="用户信息">
              {{ detailData.user }} <span class="user-type">({{ detailData.userTypeLabel }})</span>
            </el-descriptions-item>
            <el-descriptions-item label="联系电话">{{ detailData.phone }}</el-descriptions-item>
            <el-descriptions-item label="开票金额">
              <span class="amount-value">¥{{ (detailData.amount || 0).toLocaleString(undefined, { minimumFractionDigits: 2 }) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="收票邮箱">{{ detailData.email }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 开票信息 -->
        <div class="detail-card">
          <div class="card-title">
            <span class="title-line"></span>
            <h3>开票信息</h3>
          </div>
          <el-descriptions :column="4" class="info-descriptions">
            <el-descriptions-item label="发票类型">{{ detailData.invoiceType }}</el-descriptions-item>
            <el-descriptions-item label="发票抬头">{{ detailData.invoiceTitle }}</el-descriptions-item>
            <el-descriptions-item label="纳税人识别号">{{ detailData.taxId }}</el-descriptions-item>
            <el-descriptions-item label="开票内容">{{ detailData.content }}</el-descriptions-item>
            <el-descriptions-item label="备注信息" :span="4">{{ detailData.remark }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <!-- 驳回信息 (仅在状态为驳回时显示) -->
        <div class="detail-card reject-info-card" v-if="detailData.status === 'rejected'">
          <div class="card-title">
            <span class="title-line"></span>
            <h3>驳回信息</h3>
          </div>
          <div class="info-content">
            <div class="info-row">
              <span class="label">驳回原因：</span>
              <span class="value">{{ detailData.rejectReason }}</span>
            </div>
            <div class="info-row">
              <span class="label">驳回时间：</span>
              <span class="value">{{ detailData.rejectTime }}</span>
            </div>
            <div class="info-row">
              <span class="label">驳回人员：</span>
              <span class="value">{{ detailData.rejectUser }}</span>
            </div>
          </div>
        </div>

        <!-- 通过信息 (仅在状态为已开票时显示) -->
        <div class="detail-card approve-info-card" v-if="detailData.status === 'completed'">
          <div class="card-title">
            <span class="title-line"></span>
            <h3>通过信息</h3>
          </div>
          <div class="info-content">
            <div class="info-row">
              <span class="label">开票完成时间：</span>
              <span class="value">{{ detailData.finishTime }}</span>
            </div>
            <div class="info-row">
              <span class="label">上传人员：</span>
              <span class="value">{{ detailData.approveUser }}</span>
            </div>
            <div class="info-row">
              <span class="label">发票文件：</span>
              <span class="value">
                <el-link type="primary" :underline="false">{{ detailData.fileName }} ({{ detailData.fileSize }})</el-link>
              </span>
            </div>
            <div class="info-row">
              <span class="label">上传时间：</span>
              <span class="value">{{ detailData.approveTime }}</span>
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
import { ElMessage } from 'element-plus'
import Layout from '@/components/common/layout/Layout.vue'
import UploadInvoiceDialog from './components/UploadInvoiceDialog.vue'
import RejectInvoiceDialog from './components/RejectInvoiceDialog.vue'
import { getInvoiceDetail, uploadInvoice, rejectInvoice, type Invoice } from '@/api/invoice'

const router = useRouter()
const route = useRoute()

const uploadVisible = ref(false)
const rejectVisible = ref(false)
const loading = ref(false)

const detailData = ref<Partial<Invoice>>({
  invoiceNo: '',
  orderNo: '',
  applyTime: '',
  status: 'pending',
  statusLabel: '待处理',
  user: '',
  userType: '',
  userTypeLabel: '',
  phone: '',
  amount: 0,
  email: '',
  invoiceType: '',
  invoiceTitle: '',
  taxId: '',
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
      detailData.value = res.data
    }
  } catch (error: any) {
    console.error('Fetch invoice detail error:', error)
    ElMessage.error(error.backendMessage || '获取详情失败')
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
  try {
    const res = await uploadInvoice(detailData.value.id!, data)
    if (res.code === '000000') {
      ElMessage.success('上传成功')
      fetchDetail()
    }
  } catch (error: any) {
    ElMessage.error(error.backendMessage || '上传失败')
  }
}

const handleReject = () => {
  rejectVisible.value = true
}

const handleRejectSubmit = async (data: any) => {
  try {
    const res = await rejectInvoice(detailData.value.id!, data)
    if (res.code === '000000') {
      ElMessage.success('已驳回')
      fetchDetail()
    }
  } catch (error: any) {
    ElMessage.error(error.backendMessage || '驳回失败')
  }
}
</script>

<style lang="scss" scoped>
.invoice-detail-container {
  padding: 24px;
  background-color: #F0F2F5;
  min-height: calc(100vh - 80px);
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  background: #fff;
  padding: 16px 24px;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);

  .header-left {
    display: flex;
    align-items: center;
    font-size: 14px;

    .back-link {
      color: #8C8C8C;
      cursor: pointer;
      &:hover { color: #019C7C; }
    }

    .separator {
      margin: 0 8px;
      color: #BFBFBF;
    }

    .current-page {
      color: #262626;
      font-weight: 500;
    }
  }
}

.header-actions {
  display: flex;
  gap: 12px;

  .btn-upload {
    background-color: #019C7C;
    border-color: #019C7C;
    &:hover { background-color: #017c63; border-color: #017c63; }
  }

  .btn-reject {
    background-color: #F5222D;
    border-color: #F5222D;
    &:hover { background-color: #d32029; border-color: #d32029; }
  }
}

.detail-content {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.detail-card {
  background: #fff;
  padding: 24px;
  border-radius: 4px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);

  .card-title {
    display: flex;
    align-items: center;
    margin-bottom: 24px;

    .title-line {
      width: 4px;
      height: 16px;
      background: #019C7C;
      margin-right: 12px;
    }

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
      color: #262626;
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
        color: #595959;
        width: 120px;
        flex-shrink: 0;
      }

      .value {
        color: #262626;
      }
    }
  }
}

:deep(.info-descriptions) {
  .el-descriptions__label {
    width: 120px;
    color: #595959;
    font-weight: normal;
  }
  .el-descriptions__content {
    color: #262626;
  }
}

.status-tag {
  padding: 2px 8px;
  border-radius: 2px;
  font-size: 12px;
  &.pending { background: #FFF7E6; color: #FAAD14; border: 1px solid #FFE7BA; }
  &.completed { background: #F6FFED; color: #52C41A; border: 1px solid #B7EB8F; }
  &.rejected { background: #FFF1F0; color: #F5222D; border: 1px solid #FFA39E; }
}

.amount-value {
  font-weight: 600;
}

.user-type {
  color: #8C8C8C;
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
