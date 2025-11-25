<template>
  <div class="license-info-container" v-loading="loading">
    <el-dialog
      v-model="addLicenseVisible"
      class="license-add-dialog"
      :title="t('pages.licenses.detail.licenseModal.title')"
      width="520px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      align-center
      @closed="resetAddForm"
    >
      <p class="dialog-description">
        {{ t('pages.licenses.detail.licenseModal.description') }}
      </p>
      <el-form
        ref="addFormRef"
        :model="addForm"
        :rules="addFormRules"
        label-position="top"
        class="license-form-dialog"
      >
        <el-form-item :label="t('pages.licenses.detail.licenseModal.fields.authorizationCode')">
          <div class="auth-code-row">
            <el-input :model-value="authorizationCodeValue" readonly />
            <el-button text type="primary" class="copy-btn" @click="handleCopyAuthorizationCode">
              {{ t('pages.licenses.detail.licenseModal.buttons.copy') }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item
          :label="t('pages.licenses.detail.licenseModal.fields.hardwareFingerprint')"
          prop="hardware_fingerprint"
        >
          <el-input
            v-model="addForm.hardware_fingerprint"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 5 }"
            :placeholder="t('pages.licenses.detail.licenseModal.placeholders.hardwareFingerprint')"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-form-item :label="t('pages.licenses.detail.licenseModal.fields.remark')" prop="remark">
          <el-input
            v-model="addForm.remark"
            type="textarea"
            :autosize="{ minRows: 3, maxRows: 5 }"
            :placeholder="t('pages.licenses.detail.licenseModal.placeholders.remark')"
            maxlength="3000"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="handleDialogCancel" :disabled="submitting">
            {{ t('pages.licenses.detail.licenseModal.buttons.cancel') }}
          </el-button>
          <el-button type="primary" @click="submitAddLicense" :loading="submitting">
            {{ t('pages.licenses.detail.licenseModal.buttons.submit') }}
          </el-button>
        </div>
      </template>
    </el-dialog>

    <!-- 新增许可证按钮 -->
    <div class="action-bar">
      <el-button type="primary" @click="handleAddLicense">{{
        t('pages.licenses.detail.licenseInfo.addButton')
      }}</el-button>
    </div>

    <!-- 设备列表 -->
    <div v-if="!loading && devices.length > 0" class="devices-list">
        <div v-for="(device, index) in devices" :key="device.id" class="device-card">
        <!-- 设备头部 -->
        <div class="device-header">
          <div class="device-title-row">
              <span class="device-title">
                {{ t('pages.licenses.detail.licenseInfo.deviceTitle', { index: index + 1 }) }}
              </span>
            <el-tag :type="device.is_online ? 'success' : 'danger'" size="small">
                {{ device.is_online_display }}
            </el-tag>
          </div>
          <div class="device-info-row">
              <span class="device-info-item">
                {{ t('pages.licenses.detail.licenseInfo.deviceAddress') }}：
                {{ device.last_online_ip || device.activation_ip || '-' }}
              </span>
              <span class="device-info-item">
                {{ t('pages.licenses.detail.licenseInfo.lastActiveTime') }}：
                {{ formatDateTime(device.last_heartbeat || '') }}
              </span>
          </div>
        </div>

        <!-- 硬件信息 -->
        <div class="info-section hardware-section">
          <div class="section-header">
            <span class="section-title">
              {{ t('pages.licenses.detail.licenseInfo.hardwareSectionTitle') }}
            </span>
            <span class="section-subtitle">
              {{ t('pages.licenses.detail.licenseInfo.hardwareSectionSubtitle') }}
            </span>
          </div>
          <div class="info-grid">
            <div class="info-row">
              <div class="info-label">
                {{ t('pages.licenses.detail.licenseInfo.hardwareFingerprint') }}
              </div>
              <div class="info-value">{{ device.hardware_fingerprint || '-' }}</div>
            </div>
          </div>
        </div>

        <!-- 软件许可证 -->
        <div class="info-section license-section">
          <div class="section-header">
            <span class="section-title">
              {{ t('pages.licenses.detail.licenseInfo.softwareSectionTitle') }}
            </span>
            <span class="section-subtitle">
              {{ t('pages.licenses.detail.licenseInfo.softwareSectionSubtitle') }}
            </span>
          </div>
          <div class="license-content">
            <div class="license-row">
              <div class="license-label">
                {{ t('pages.licenses.detail.licenseInfo.licenseKey') }}
              </div>
              <div class="license-value">{{ device.license_key || '-' }}</div>
            </div>
            <div class="license-row">
              <div class="license-label">
                {{ t('pages.licenses.detail.licenseInfo.licenseStatus') }}
              </div>
              <div class="license-value status-text">
                {{
                  device.status_display ||
                  t('pages.licenses.detail.licenseInfo.statusDescription')
                }}
              </div>
            </div>
          </div>
        </div>

        <!-- 下载按钮 -->
        <div class="download-section">
          <el-button
            class="download-license-btn"
            type="primary"
            :loading="isDownloading(device.id)"
            @click="handleDownloadLicense(device)"
          >
            {{ t('pages.licenses.detail.licenseInfo.downloadButton') }}
          </el-button>
          <el-button
            class="revoke-license-btn"
            type="danger"
            plain
            :loading="isRevoking(device.id)"
            :disabled="device.status === 'revoked'"
            @click="handleRevokeLicense(device)"
          >
            {{ t('pages.licenses.detail.actions.revokeLicense') }}
          </el-button>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && devices.length === 0" class="empty-state">
      <el-empty :description="t('pages.licenses.detail.licenseInfo.emptyDescription')">
        <el-button type="primary" @click="handleAddLicense">
          {{ t('pages.licenses.detail.licenseInfo.addButton') }}
        </el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, reactive, computed, nextTick } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'
import type { AuthorizationCode, LicenseDevice, LicenseDeviceCreateRequest } from '@/api/license'
import { getLicenseDevices, getLicenseDeviceDetail, createLicenseDevice, downloadLicenseFile, revokeLicense } from '@/api/license'
import { formatDate } from '@/utils/date'

interface Props {
  licenseData: AuthorizationCode | null
}

const emit = defineEmits<{
  (e: 'licenseRevoked'): void
}>()

const props = defineProps<Props>()
const { t } = useI18n()

const devices = ref<LicenseDevice[]>([])
const loading = ref(false)
const addLicenseVisible = ref(false)
const submitting = ref(false)
const addFormRef = ref<FormInstance>()
const addForm = reactive({
  hardware_fingerprint: '',
  remark: ''
})
const downloadingMap = reactive<Record<string, boolean>>({})
const revokingMap = reactive<Record<string, boolean>>({})

const authorizationCodeId = computed(() => props.licenseData?.id || '')
const authorizationCodeValue = computed(() => props.licenseData?.code || '--')

const validateHardware = (_rule: any, value: string, callback: (error?: Error) => void) => {
  const trimmed = value?.trim()
  if (!trimmed) {
    callback(new Error(t('pages.licenses.detail.licenseModal.validation.hardwareRequired')))
    return
  }
  if (trimmed.length > 200) {
    callback(new Error(t('pages.licenses.detail.licenseModal.validation.hardwareMax')))
    return
  }
  callback()
}

const validateRemark = (_rule: any, value: string, callback: (error?: Error) => void) => {
  const trimmed = value?.trim()
  if (trimmed && trimmed.length > 3000) {
    callback(new Error(t('pages.licenses.detail.licenseModal.validation.remarkMax')))
    return
  }
  callback()
}

const addFormRules: FormRules = {
  hardware_fingerprint: [
    { validator: validateHardware, trigger: ['blur', 'change'] }
  ],
  remark: [{ validator: validateRemark, trigger: ['blur', 'change'] }]
}

// 获取设备列表
const fetchDevices = async () => {
  if (!authorizationCodeId.value) {
    devices.value = []
    return
  }

  loading.value = true
  try {
    const response = await getLicenseDevices({
      authorization_code_id: authorizationCodeId.value,
      page: 1,
      page_size: 100
    })

    if (response.data?.list && response.data.list.length > 0) {
      // 检查第一个设备是否有 device_info，如果没有则需要获取详情
      const needDetails = !response.data.list[0].device_info

      if (needDetails) {
        // 并发获取所有设备的详细信息
        const detailPromises = response.data.list.map(device =>
          getLicenseDeviceDetail(device.id)
            .then(res => res.data)
            .catch(err => {
              console.error(`获取设备 ${device.id} 详情失败:`, err)
              return device // 如果获取失败，使用原数据
            })
        )

        const detailedDevices = await Promise.all(detailPromises)
        devices.value = detailedDevices as LicenseDevice[]
      } else {
        devices.value = response.data.list
      }
    } else {
      devices.value = []
    }
  } catch (error) {
    console.error('获取设备列表失败:', error)
    ElMessage.error(t('pages.licenses.detail.licenseInfo.messages.loadError'))
    devices.value = []
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchDevices()
})

// 监听licenseData变化，重新获取数据
watch(
  () => authorizationCodeId.value,
  () => {
    fetchDevices()
  }
)

// 格式化日期时间
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return '-'
  return formatDate(dateTime)
}

const resetAddForm = () => {
  addForm.hardware_fingerprint = ''
  addForm.remark = ''
  nextTick(() => {
    addFormRef.value?.clearValidate()
  })
}

const handleDialogCancel = () => {
  addLicenseVisible.value = false
  resetAddForm()
}

const handleAddLicense = () => {
  if (!authorizationCodeId.value) {
    ElMessage.warning(t('pages.licenses.detail.licenseModal.messages.missingAuthorization'))
    return
  }
  addLicenseVisible.value = true
  nextTick(() => {
    addFormRef.value?.clearValidate()
  })
}

const handleCopyAuthorizationCode = async () => {
  if (!authorizationCodeValue.value || authorizationCodeValue.value === '--') {
    ElMessage.warning(t('pages.licenses.detail.licenseModal.messages.missingAuthorization'))
    return
  }

  try {
    await navigator.clipboard.writeText(authorizationCodeValue.value)
    ElMessage.success(t('pages.licenses.detail.licenseModal.messages.copySuccess'))
  } catch (error) {
    console.error('复制授权码失败:', error)
    ElMessage.error(t('pages.licenses.detail.licenseModal.messages.copyError'))
  }
}

const isDownloading = (licenseId: string) => {
  return Boolean(licenseId && downloadingMap[licenseId])
}

const isRevoking = (licenseId: string) => {
  return Boolean(licenseId && revokingMap[licenseId])
}

const extractFilename = (contentDisposition?: string) => {
  if (!contentDisposition) return ''
  const utfMatch = /filename\*=utf-8''([^;]+)/i.exec(contentDisposition)
  if (utfMatch?.[1]) {
    try {
      return decodeURIComponent(utfMatch[1])
    } catch (error) {
      console.warn('Failed to decode filename from header:', error)
    }
  }

  const asciiMatch = /filename="?([^"]+)"?/i.exec(contentDisposition)
  return asciiMatch?.[1] || ''
}

const handleDownloadLicense = async (device: LicenseDevice) => {
  if (!device?.id) {
    ElMessage.warning(t('pages.licenses.detail.licenseInfo.messages.missingLicense'))
    return
  }

  downloadingMap[device.id] = true
  try {
    const response = await downloadLicenseFile(device.id)
    const disposition =
      response.headers?.['content-disposition'] || response.headers?.['Content-Disposition']
    const filename =
      extractFilename(disposition) ||
      `${device.customer_name || 'license'}-${device.id}.zip`

    const blobUrl = window.URL.createObjectURL(response.data)
    const link = document.createElement('a')
    link.href = blobUrl
    link.download = filename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(blobUrl)

    ElMessage.success(t('pages.licenses.detail.licenseInfo.messages.downloadSuccess'))
  } catch (error: any) {
    console.error('下载许可证失败:', error)
    const backendMessage = error?.backendMessage || error?.message
    ElMessage.error(
      backendMessage || t('pages.licenses.detail.licenseInfo.messages.downloadError')
    )
  } finally {
    downloadingMap[device.id] = false
  }
}

const handleRevokeLicense = async (device: LicenseDevice) => {
  if (!device?.id) {
    ElMessage.error(t('pages.licenses.detail.messages.missingId'))
    return
  }

  if (device.status === 'revoked') {
    return
  }

  let reason = ''
  try {
    const { value } = await ElMessageBox.prompt(
      t('pages.licenses.detail.messages.revokeConfirm'),
      t('pages.licenses.detail.messages.revokeTitle'),
      {
        confirmButtonText: t('pages.licenses.detail.messages.revokeConfirmButton'),
        cancelButtonText: t('pages.licenses.detail.messages.revokeCancelButton'),
        type: 'warning',
        inputType: 'textarea',
        inputPlaceholder: t('pages.licenses.detail.messages.revokeReasonPlaceholder'),
        inputValidator: (val: string) => {
          if (!val || !val.trim()) {
            return t('pages.licenses.detail.messages.revokeReasonRequired')
          }
          if (val.trim().length > 500) {
            return t('pages.licenses.detail.messages.revokeReasonMax')
          }
          return true
        },
        inputErrorMessage: t('pages.licenses.detail.messages.revokeReasonRequired')
      }
    )
    reason = value.trim()
  } catch {
    return
  }

  revokingMap[device.id] = true
  try {
    await revokeLicense(device.id, reason)
    ElMessage.success(t('pages.licenses.detail.messages.revokeSuccess'))
    await fetchDevices()
    emit('licenseRevoked')
  } catch (error: any) {
    ElMessage.error(error?.message || t('pages.licenses.detail.messages.revokeError'))
  } finally {
    revokingMap[device.id] = false
  }
}

const submitAddLicense = async () => {
  if (!addFormRef.value) return

  try {
    await addFormRef.value.validate()
    if (!authorizationCodeId.value) {
      throw new Error(t('pages.licenses.detail.licenseModal.messages.missingAuthorization'))
    }

    submitting.value = true
    const payload: LicenseDeviceCreateRequest = {
      authorization_code_id: authorizationCodeId.value,
      hardware_fingerprint: addForm.hardware_fingerprint.trim()
    }

    if (addForm.remark.trim()) {
      payload.device_info = { remark: addForm.remark.trim() }
    }

    const response = await createLicenseDevice(payload)
    if (response.code === '000000') {
      ElMessage.success(
        response.message || t('pages.licenses.detail.licenseModal.messages.success')
      )
      addLicenseVisible.value = false
      resetAddForm()
      await fetchDevices()
    } else {
      throw new Error(response.message || t('pages.licenses.detail.licenseModal.messages.error'))
    }
  } catch (error: any) {
    console.error('创建许可证失败:', error)
    ElMessage.error(error.message || t('pages.licenses.detail.licenseModal.messages.error'))
  } finally {
    submitting.value = false
  }
}
</script>

<style lang="scss" scoped>
.license-info-container {
  width: 100%;
}

.license-add-dialog {
  :deep(.el-dialog) {
    border-radius: 12px;
    padding-bottom: 12px;
  }

  :deep(.el-dialog__header) {
    margin-right: 0;
    border-bottom: 1px solid #f0f2f5;
    padding: 16px 24px;
  }

  :deep(.el-dialog__title) {
    font-family: 'Source Han Sans CN', sans-serif;
    font-size: 18px;
    font-weight: 600;
    color: #1d1d1d;
  }

  :deep(.el-dialog__body) {
    padding: 16px 24px 0;
  }

  :deep(.el-dialog__footer) {
    padding: 0 24px 24px;
  }
}

.dialog-description {
  font-size: 14px;
  color: #666;
  line-height: 22px;
  margin-bottom: 16px;
}

.license-form-dialog {
  :deep(.el-form-item__label) {
    font-weight: 500;
    color: #1d1d1d;
  }
}

.auth-code-row {
  display: flex;
  align-items: center;
  gap: 12px;

  :deep(.el-input__wrapper) {
    background: rgba(247, 248, 250, 0.8);
  }
}

.copy-btn {
  padding: 0 8px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.action-bar {
  margin-bottom: 16px;

  :deep(.el-button) {
    padding: 7px 16px;
    font-family: 'Source Han Sans CN', sans-serif;
    font-size: 14px;
    font-weight: 500;
    line-height: 21px;
    border-radius: 4px;
  }
}

.devices-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.device-card {
  background: #FFFFFF;
  border: 1px solid #E2E2E2;
  border-radius: 4px;
  padding: 16px 20px;
}

.device-header {
  margin-bottom: 20px;

  .device-title-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;

    .device-title {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 16px;
      font-weight: 700;
      line-height: 18px;
      color: #1D1D1D;
    }

    :deep(.el-tag) {
      padding: 7px 16px;
      height: 24px;
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 500;
      line-height: 21px;
      border-radius: 4px;

      &.el-tag--success {
        background: rgba(0, 194, 124, 0.08);
        border-color: transparent;
        color: #019C7C;
      }

      &.el-tag--danger {
        background: rgba(240, 20, 47, 0.08);
        border-color: transparent;
        color: #F0142F;
      }
    }
  }

  .device-info-row {
    display: flex;
    gap: 40px;

    .device-info-item {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      color: #1D1D1D;
    }
  }
}

.info-section {
  background: rgba(247, 248, 250, 0.7);
  border: 1px solid rgba(226, 226, 226, 0.6);
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 16px;

  &:last-child {
    margin-bottom: 0;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-bottom: 16px;

    .section-title {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 16px;
      font-weight: 700;
      line-height: 18px;
      color: #1D1D1D;
    }

    .section-subtitle {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 18px;
      color: #666666;
    }
  }
}

.hardware-section {
  .info-grid {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .info-row {
      display: flex;
      align-items: center;
      gap: 30px;

      .info-label {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        white-space: nowrap;
        min-width: 100px;
      }

      .info-value {
        flex: 1;
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        padding: 8px 12px;
        background: rgba(136, 165, 209, 0.2);
        border-radius: 8px;
        word-break: break-all;
      }
    }
  }
}

.license-section {
  .license-content {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .license-row {
      display: flex;
      align-items: center;
      gap: 30px;

      .license-label {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        white-space: nowrap;
        min-width: 100px;
      }

      .license-value {
        flex: 1;
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        padding: 8px 12px;
        background: rgba(136, 165, 209, 0.2);
        border-radius: 8px;
        word-break: break-all;

        &.status-text {
          color: #4763FF;
          background: transparent;
          padding: 0;
        }
      }
    }
  }

  .download-license-btn {
    min-width: 148px;
    background: #00C27C;
    border: none;
    color: #FFFFFF;
    font-family: 'Source Han Sans CN', sans-serif;
    font-size: 14px;
    font-weight: 500;
    line-height: 21px;
    border-radius: 4px;

    &:hover {
      background: #019C7C;
    }
  }
}

.download-section {
  display: flex;
  justify-content: flex-start;
  gap: 12px;
  flex-wrap: wrap;
}

.revoke-license-btn {
  min-width: 148px;
  border: 1px solid #F0142F;
  color: #F0142F;
  background: transparent;

  &:hover:not(.is-disabled) {
    background: rgba(240, 20, 47, 0.08);
    color: #d01228;
    border-color: #d01228;
  }
}

.empty-state {
  padding: 40px;
  text-align: center;
}

// 响应式布局
@media (max-width: 768px) {
  .device-header {
    .device-info-row {
      flex-direction: column;
      gap: 8px;
    }
  }

  .hardware-section {
    .info-grid {
      .info-row {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;

        .info-label {
          min-width: auto;
        }

        .info-value {
          width: 100%;
        }
      }
    }
  }

  .license-section {
    .license-content {
      .license-row {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;

        .license-label {
          min-width: auto;
        }

        .license-value {
          width: 100%;
        }
      }
    }
  }
}
</style>
