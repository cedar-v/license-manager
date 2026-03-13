<template>
  <div class="license-view-page">
    <!-- 面包屑导航 -->
    <div class="breadcrumb-section">
      <span class="breadcrumb-item clickable" @click="goBack">{{ $t('pages.licenses.detail.breadcrumb.licenseManagement') }}</span>/<span class="breadcrumb-item active">{{ $t('pages.licenses.detail.breadcrumb.licenseDetail') }}</span>
    </div>

    <!-- 授权详情内容 -->
    <div class="content-section" v-loading="loading">
      <!-- 左侧菜单栏 -->
      <div class="left-sidebar">
        <div
          class="sidebar-item"
          :class="{ active: activeTab === 'basic' }"
          @click="activeTab = 'basic'"
        >
          <LicenseTabIcon name="basic" :active="activeTab === 'basic'" />
          <span class="text">{{ $t('pages.licenses.detail.tabs.basic') }}</span>
        </div>
        <div
          class="sidebar-item"
          :class="{ active: activeTab === 'authorization' }"
          @click="activeTab = 'authorization'"
        >
          <LicenseTabIcon name="authorization" :active="activeTab === 'authorization'" />
          <span class="text">{{ $t('pages.licenses.detail.tabs.authorization') }}</span>
        </div>
        <div
          class="sidebar-item"
          :class="{ active: activeTab === 'license' }"
          @click="activeTab = 'license'"
        >
          <LicenseTabIcon name="license" :active="activeTab === 'license'" />
          <span class="text">{{ $t('pages.licenses.detail.tabs.license') }}</span>
        </div>
        <div
          class="sidebar-item"
          :class="{ active: activeTab === 'history' }"
          @click="activeTab = 'history'"
        >
          <LicenseTabIcon name="history" :active="activeTab === 'history'" />
          <span class="text">{{ $t('pages.licenses.detail.tabs.history') }}</span>
        </div>
      </div>

      <div class="right-content">
        <!-- 授权码和状态标签 -->
        <div class="license-header">
          <div class="license-code-group">
            <h2 class="license-code">{{ licenseData?.code || '-' }}</h2>
            <el-tag :type="getStatusType(licenseData?.status)" class="status-tag">
              {{ licenseData?.status_display || '-' }}
            </el-tag>
          </div>
          <div class="action-buttons-inline">
            <el-button class="action-btn-inline copy-btn" @click="handleCopyCode">
              {{ $t('pages.licenses.detail.actions.copyCode') }}
            </el-button>
            <el-button class="action-btn-inline update-btn" @click="openUpdateDialog">
              {{ $t('pages.licenses.detail.actions.updateLicense') }}
            </el-button>
            <el-button class="action-btn-inline renew-btn" @click="openChangeValidityDialog">
              {{ $t('pages.licenses.detail.actions.changeValidity') }}
            </el-button>
            <!-- update button moved -->
            <el-button class="action-btn-inline download-btn" @click="handleDownloadCertificate">
              {{ $t('pages.licenses.detail.actions.downloadCertificate') }}
            </el-button>
          </div>
        </div>

        <!-- 变更授权期限弹窗 -->
        <el-dialog
          :title="t('pages.licenses.detail.dialogs.changeValidityTitle')"
          v-model="changeDialogVisible"
          width="520px"
        >
          <div>
            <el-form label-width="100px" :model="{ validity_type: changeValidityType, date_range: changeDateRange, reason: changeReason }">
              <el-form-item :label="t('pages.licenses.detail.dialogs.validityType')">
                <el-radio-group v-model="changeValidityType">
                  <el-radio label="limited">{{ t('pages.licenses.detail.dialogs.limited') }}</el-radio>
                  <el-radio label="permanent">{{ t('pages.licenses.detail.dialogs.permanent') }}</el-radio>
                </el-radio-group>
              </el-form-item>

              <el-form-item v-if="changeValidityType === 'limited'" :label="t('pages.licenses.detail.dialogs.dateRange')">
                <el-date-picker
                  v-model="changeDateRange"
                  type="daterange"
                  range-separator=" - "
                  start-placeholder="Start date"
                  end-placeholder="End date"
                  :disabled-date="disabledDate"
                  value-format="YYYY-MM-DD"
                  @change="handleChangeDateRange"
                />
              </el-form-item>

              <el-form-item :label="t('pages.licenses.detail.dialogs.reason')">
                <el-input
                  type="textarea"
                  v-model="changeReason"
                  :placeholder="t('pages.licenses.detail.dialogs.reasonPlaceholder')"
                  :rows="3"
                  maxlength="500"
                />
              </el-form-item>
            </el-form>
          </div>
          <template #footer>
            <el-button @click="changeDialogVisible = false">{{ t('common.cancel') }}</el-button>
            <el-button type="primary" :loading="changeSubmitting" @click="handleConfirmChangeValidity">{{ t('common.confirm') }}</el-button>
          </template>
        </el-dialog>

        <!-- 变更授权弹窗（编辑激活数/配置等） -->
        <el-dialog
          :title="t('pages.licenses.detail.updateDialog.title')"
          v-model="updateDialogVisible"
          width="720px"
        >
          <div>
            <el-form label-width="140px" class="update-dialog-form">
              <el-form-item :label="t('pages.licenses.detail.updateDialog.fields.changeType')">
                <el-select v-model="updateChangeType" placeholder="Select" style="width:100%">
                  <el-option
                    v-for="opt in changeTypeOptions"
                    :key="opt.key"
                    :label="opt.display"
                    :value="opt.key"
                  />
                </el-select>
              </el-form-item>

              <el-form-item :label="t('pages.licenses.detail.updateDialog.fields.maxActivations')">
                <el-input-number v-model="updateMaxActivations" :min="1" :controls="true" />
              </el-form-item>

              <el-form-item :label="t('pages.licenses.detail.updateDialog.fields.featureConfig')">
                <JsonEditor ref="featureConfigRef" v-model="featureConfigData" />
              </el-form-item>

              <el-form-item :label="t('pages.licenses.detail.updateDialog.fields.usageLimits')">
                <JsonEditor ref="usageLimitsRef" v-model="usageLimitsData" />
              </el-form-item>

              <el-form-item :label="t('pages.licenses.detail.updateDialog.fields.customParameters')">
                <JsonEditor ref="customParamsRef" v-model="customParamsData" />
              </el-form-item>

              <el-form-item :label="t('pages.licenses.detail.updateDialog.fields.reason')">
                <el-input type="textarea" v-model="updateReason" :rows="3" maxlength="500" />
              </el-form-item>
            </el-form>
          </div>
          <template #footer>
            <el-button @click="updateDialogVisible = false">{{ t('common.cancel') }}</el-button>
            <el-button type="primary" :loading="updateSubmitting" @click="handleConfirmUpdate">{{ t('common.confirm') }}</el-button>
          </template>
        </el-dialog>

        <!-- Tab 内容区域 -->
        <div class="tab-content-area">
          <BasicInfo v-if="activeTab === 'basic'" :license-data="licenseData" />
          <AuthorizationInfo v-if="activeTab === 'authorization'" :license-data="licenseData" />
          <LicenseInfo
            v-if="activeTab === 'license'"
            :license-data="licenseData"
            @license-revoked="loadLicenseData"
          />
          <ChangeHistory v-if="activeTab === 'history'" :license-data="licenseData" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getLicenseDetail, downloadAuthorizationFile, type AuthorizationCode } from '@/api/license'
import LicenseTabIcon from '@/components/common/icons/LicenseTabIcon.vue'
import BasicInfo from './components/BasicInfo.vue'
import AuthorizationInfo from './components/AuthorizationInfo.vue'
import LicenseInfo from './components/LicenseInfo.vue'
import ChangeHistory from './components/ChangeHistory.vue'
import JsonEditor from '@/components/common/JsonEditor.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const loading = ref(false)
const activeTab = ref('basic')
const licenseData = ref<AuthorizationCode | null>(null)

const goBack = () => {
  router.back()
}

const getStatusType = (status?: string) => {
  switch (status) {
    case 'normal':
      return 'success'
    case 'locked':
      return 'warning'
    case 'expired':
      return 'danger'
    default:
      return 'info'
  }
}

const handleCopyCode = () => {
  if (!licenseData.value?.code) {
    ElMessage.warning(t('pages.licenses.detail.messages.codeEmpty'))
    return
  }

  navigator.clipboard.writeText(licenseData.value.code).then(() => {
    ElMessage.success(t('pages.licenses.detail.messages.copySuccess'))
  }).catch(() => {
    ElMessage.error(t('pages.licenses.detail.messages.copyError'))
  })
}

import { updateLicense } from '@/api/license'
import { getEnumOptions } from '@/api/enum'

// --- change validity dialog state & handlers ---
const changeDialogVisible = ref(false)
const changeValidityType = ref<'permanent' | 'limited'>('limited')
const changeDateRange = ref<[string, string] | null>(null)
const changeValidityDays = ref<number>(365)
const changeReason = ref<string>('')
const changeSubmitting = ref(false)

  const disabledDate = (time: Date) => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return time.getTime() < today.getTime()
}

const openChangeValidityDialog = () => {
  // initialize with current license data
  console.log('[Licenses] openChangeValidityDialog called')
  if (licenseData.value && licenseData.value.start_date && licenseData.value.end_date) {
    changeDateRange.value = [licenseData.value.start_date, licenseData.value.end_date]
    const startDate = new Date(changeDateRange.value[0])
    const endDate = new Date(changeDateRange.value[1])
    const diffTime = endDate.getTime() - startDate.getTime()
    changeValidityDays.value = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    changeValidityType.value = changeValidityDays.value >= 365000 ? 'permanent' : 'limited'
  } else {
    changeDateRange.value = null
    changeValidityDays.value = 365
    changeValidityType.value = 'limited'
  }
  changeReason.value = ''
  changeDialogVisible.value = true
}

// watch for dialog visibility changes and verify DOM mount
watch(changeDialogVisible, async (val) => {
  console.log('[Licenses] changeDialogVisible ->', val)
  if (val) {
    await nextTick()
    const wrapper = document.querySelector('.el-dialog__wrapper') || document.querySelector('.el-dialog')
    console.log('[Licenses] dialog DOM element:', wrapper)
    if (!wrapper) {
      ElMessage.error(t('pages.licenses.detail.messages.changeValidityError') + ' (dialog not mounted)')
    }
  }
})

const handleChangeDateRange = (dates: [string, string] | null) => {
  if (dates && dates.length === 2) {
    const startDate = new Date(dates[0])
    const endDate = new Date(dates[1])
    const diffTime = endDate.getTime() - startDate.getTime()
    changeValidityDays.value = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  }
}

const handleConfirmChangeValidity = async () => {
  if (!licenseData.value?.id) {
    ElMessage.error(t('pages.licenses.detail.messages.missingId'))
    return
  }

  if (!changeReason.value || !changeReason.value.trim()) {
    ElMessage.warning(t('pages.licenses.detail.messages.changeValidityReasonRequired'))
    return
  }

  changeSubmitting.value = true
  try {
    // build start_date / end_date payload (preferred). Format YYYY-MM-DD.
    const formatDateYMD = (d: Date) => {
      const y = d.getFullYear()
      const m = `${d.getMonth() + 1}`.padStart(2, '0')
      const day = `${d.getDate()}`.padStart(2, '0')
      return `${y}-${m}-${day}`
    }

    let startDateStr: string | null = null
    let endDateStr: string | null = null

    if (changeValidityType.value === 'permanent') {
      const today = new Date()
      startDateStr = formatDateYMD(today)
      const endDate = new Date(today)
      // permanent maps to 365000 days (back-end convention)
      endDate.setDate(endDate.getDate() + 365000 - 1)
      endDateStr = formatDateYMD(endDate)
    } else {
      if (!changeDateRange.value || changeDateRange.value.length !== 2 || !changeDateRange.value[0] || !changeDateRange.value[1]) {
        ElMessage.warning(t('pages.licenses.detail.form.validation.dateRangeRequired'))
        changeSubmitting.value = false
        return
      }
  // date picker returns strings when value-format is set; handle both and avoid timezone drift by parsing Y-M-D explicitly
  const rawStart = changeDateRange.value[0]
  const rawEnd = changeDateRange.value[1]
  const parseYMD = (s: string) => {
    const parts = s.split('-')
    return new Date(Number(parts[0]), Number(parts[1]) - 1, Number(parts[2]))
  }
  const start = typeof rawStart === 'string' ? parseYMD(rawStart) : (rawStart as any as Date)
  const end = typeof rawEnd === 'string' ? parseYMD(rawEnd) : (rawEnd as any as Date)
      startDateStr = formatDateYMD(start)
      endDateStr = formatDateYMD(end)
    }

    console.log('[Licenses] submitting changeValidity', { changeValidityType: changeValidityType.value, changeDateRange: changeDateRange.value, startDateStr, endDateStr })

    const payload: any = {
      change_type: 'renewal',
      reason: changeReason.value
    }
    if (startDateStr && endDateStr) {
      payload.start_date = startDateStr
      payload.end_date = endDateStr
    }
    // include update dialog fields if dialog visible
    if (updateDialogVisible.value) {
      if (updateMaxActivations.value !== null && updateMaxActivations.value !== undefined) {
        payload.max_activations = updateMaxActivations.value
      }
      // 使用 JsonEditor 数据
      payload.feature_config = featureConfigData.value || {}
      payload.usage_limits = usageLimitsData.value || {}
      payload.custom_parameters = customParamsData.value || {}
      payload.change_type = updateChangeType.value || payload.change_type
      // if reason in update dialog provided, override
      if (updateReason.value && updateReason.value.trim()) {
        payload.reason = updateReason.value
      }
    }

    const res = await updateLicense(licenseData.value.id, payload)
    if (res.code === '000000') {
      ElMessage.success(res.message || t('pages.licenses.detail.messages.changeValiditySuccess'))
      changeDialogVisible.value = false
      await loadLicenseData()
    } else {
      throw new Error(res.message || t('pages.licenses.detail.messages.changeValidityError'))
    }
  } catch (err: any) {
    console.error('change validity error', err)
    const backendMsg = err?.response?.data?.message
    ElMessage.error(backendMsg || err.message || t('pages.licenses.detail.messages.changeValidityError'))
  } finally {
    changeSubmitting.value = false
  }
}

// --- update authorization dialog ---
import type { RawEnumItem } from '@/api/enum'

const updateDialogVisible = ref(false)
const changeTypeOptions = ref<RawEnumItem[]>([])
const updateChangeType = ref<string>('other')
const updateMaxActivations = ref<number | null>(null)
const featureConfigData = ref<Record<string, any>>({})
const usageLimitsData = ref<Record<string, any>>({})
const customParamsData = ref<Record<string, any>>({})
const updateReason = ref<string>('')
const updateSubmitting = ref(false)

// JsonEditor 组件引用
const featureConfigRef = ref<InstanceType<typeof JsonEditor>>()
const usageLimitsRef = ref<InstanceType<typeof JsonEditor>>()
const customParamsRef = ref<InstanceType<typeof JsonEditor>>()

const openUpdateDialog = async () => {
  updateChangeType.value = 'other'
  updateReason.value = ''
  updateMaxActivations.value = licenseData.value?.max_activations || null

  // 使用 JsonEditor 加载数据
  featureConfigData.value = typeof licenseData.value?.feature_config === 'string'
    ? JSON.parse(licenseData.value?.feature_config || '{}')
    : (licenseData.value?.feature_config || {})
  usageLimitsData.value = typeof licenseData.value?.usage_limits === 'string'
    ? JSON.parse(licenseData.value?.usage_limits || '{}')
    : (licenseData.value?.usage_limits || {})
  customParamsData.value = typeof licenseData.value?.custom_parameters === 'string'
    ? JSON.parse(licenseData.value?.custom_parameters || '{}')
    : (licenseData.value?.custom_parameters || {})

  // load change type options from backend enum
  try {
    const res = await getEnumOptions('authorization_change_type')
    if (res.code === '000000' && res.data?.items) {
      changeTypeOptions.value = res.data.items
      if (!changeTypeOptions.value.find(i => i.key === updateChangeType.value) && changeTypeOptions.value.length) {
        updateChangeType.value = changeTypeOptions.value[0].key
      }
    }
  } catch (err) {
    console.error('load change type enums failed', err)
  }

  updateDialogVisible.value = true
}

const handleConfirmUpdate = async () => {
  if (!licenseData.value?.id) {
    ElMessage.error(t('pages.licenses.detail.messages.missingId'))
    return
  }
  if (!updateChangeType.value || !updateChangeType.value.trim()) {
    ElMessage.warning(t('pages.licenses.detail.updateDialog.validation.changeTypeRequired'))
    return
  }
  if (!updateReason.value || !updateReason.value.trim()) {
    ElMessage.warning(t('pages.licenses.detail.updateDialog.validation.reasonRequired'))
    return
  }

  // 验证 JsonEditor 数据
  const featureValid = featureConfigRef.value?.validate() ?? true
  const limitValid = usageLimitsRef.value?.validate() ?? true
  const customValid = customParamsRef.value?.validate() ?? true
  if (!featureValid || !limitValid || !customValid) {
    ElMessage.error(t('pages.licenses.form.keyValue.validationFailed'))
    changeSubmitting.value = false
    return
  }

  // 构建 payload
  const payload: any = {
    change_type: updateChangeType.value,
    reason: updateReason.value
  }
  if (updateMaxActivations.value !== null && updateMaxActivations.value !== undefined) {
    payload.max_activations = updateMaxActivations.value
  }

  // 使用 JsonEditor 数据
  payload.feature_config = featureConfigData.value || {}
  payload.usage_limits = usageLimitsData.value || {}
  payload.custom_parameters = customParamsData.value || {}

  updateSubmitting.value = true
  try {
    const res = await updateLicense(licenseData.value.id, payload)
    if (res.code === '000000') {
      ElMessage.success(res.message || t('pages.licenses.detail.updateDialog.messages.updateSuccess'))
      updateDialogVisible.value = false
      await loadLicenseData()
    } else {
      throw new Error(res.message || t('pages.licenses.detail.updateDialog.messages.updateError'))
    }
  } catch (err: any) {
    console.error('update authorization error', err)
    const backendMsg = err?.response?.data?.message
    ElMessage.error(backendMsg || err.message || t('pages.licenses.detail.updateDialog.messages.updateError'))
  } finally {
    updateSubmitting.value = false
  }
}

const getFileNameFromDisposition = (disposition?: string | null) => {
  const defaultName = 'authorization_package.zip'
  if (!disposition) return defaultName

  const utf8Match = disposition.match(/filename\*=UTF-8''([^;]+)/i)
  if (utf8Match?.[1]) {
    try {
      return decodeURIComponent(utf8Match[1])
    } catch {
      return utf8Match[1]
    }
  }

  const asciiMatch = disposition.match(/filename="?([^";]+)"?/i)
  if (asciiMatch?.[1]) {
    return asciiMatch[1]
  }

  return defaultName
}

const handleDownloadCertificate = async () => {
  if (!licenseData.value?.id) {
    ElMessage.error(t('pages.licenses.detail.messages.missingId'))
    return
  }

  try {
    const response = await downloadAuthorizationFile(licenseData.value.id)
    const fileName = getFileNameFromDisposition(response.headers['content-disposition'])
    const blob = new Blob([response.data], { type: 'application/zip' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success(t('pages.licenses.detail.messages.downloadSuccess'))
  } catch (error: any) {
    console.error('下载授权文件失败:', error)
    let message = t('pages.licenses.detail.messages.downloadError')
    if (error?.response?.data instanceof Blob) {
      try {
        const text = await error.response.data.text()
        const parsed = JSON.parse(text)
        if (parsed?.message) {
          message = parsed.message
        }
      } catch {
        // ignore parse errors
      }
    }
    ElMessage.error(message)
  }
}

const loadLicenseData = async () => {
  const id = route.params.id as string
  if (!id) {
    ElMessage.error(t('pages.licenses.detail.messages.missingId'))
    router.back()
    return
  }

  loading.value = true
  try {
    const response = await getLicenseDetail(id)
    if (response.code === '000000' && response.data) {
      licenseData.value = response.data
      console.log(licenseData.value)
    } else {
      throw new Error(response.message || t('pages.licenses.detail.messages.loadError'))
    }
  } catch (error: any) {
    console.error('加载授权详情失败:', error)
    ElMessage.error(error.message || t('pages.licenses.detail.messages.loadError'))
    router.back()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadLicenseData()
})
// load change type enums on mount for dropdown
const loadChangeTypeEnums = async () => {
  try {
    const res = await getEnumOptions('authorization_change_type')
    if (res.code === '000000' && res.data?.items) {
      changeTypeOptions.value = res.data.items
    }
  } catch (err) {
    console.error('load change type enums failed', err)
  }
}
onMounted(() => {
  loadChangeTypeEnums()
})
</script>

<style lang="scss" scoped>
.license-view-page {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--app-bg-color);
  padding: 80px 60px 40px 60px;
  overflow: hidden;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 400px;
    background-image: url('@/assets/images/licenseView-bg.png');
    background-size: cover;
    background-position: center;
    background-repeat: no-repeat;
    z-index: 0;
    pointer-events: none;
  }

  > * {
    position: relative;
    z-index: 1;
  }
}

.breadcrumb-section {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding: 0 4px;
  position: absolute;
    left: 20px;
    top: 15px;

  .breadcrumb-item {
    
    font-size: 14px;
    font-weight: 400;
    line-height: 1.4285714285714286em;
    color: #888888;

    &.clickable {
      cursor: pointer;
      transition: color 0.3s;

      &:hover {
        color: #019C7C;
      }
    }

    &.active {
      color: #1D1D1D;
    }
  }

  .breadcrumb-separator {
    flex-shrink: 0;
  }
}

.content-section {
  display: flex;
  gap: 0;
  flex: 1;
  overflow: hidden;
  background: white;
  border-radius: 12px;
  box-shadow: 0px 0px 20px 0px rgba(0, 0, 0, 0.2);
}

.left-sidebar {
  width: 300px;
  background: var(--app-content-bg);
  border-radius: 12px 0 0 12px;
  border-right: 1px solid #DCDEE2;
  padding: 13px 0;
  flex-shrink: 0;

  .sidebar-item {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    gap: 12px;
    padding: 16px 24px;
    margin: 8px 16px;
    cursor: pointer;
    transition: all 0.3s;
    color: #B2B8C2;
    background: transparent;
    border-radius: 28px;

    .text {
      
      font-size: 16px;
      font-weight: 500;
      line-height: 1.3;
    }

    &:hover {
      background: rgba(0, 194, 124, 0.08);
    }

    &.active {
      background: rgba(0, 194, 124, 0.12);
      color: #019C7C;
      font-weight: 700;
      box-shadow: 0px 2px 32px 0px rgba(0, 0, 0, 0.02);

      .text {
        font-weight: 700;
      }
    }
  }
}

.right-content {
  flex: 1;
  background: white;
  border-radius: 0 12px 12px 0;
  padding: 0;
  overflow-y: auto;
  border-top: 1px solid #DCDEE2;
  display: flex;
  flex-direction: column;
}

.license-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 24px;
  border-bottom: 1px solid #DCDEE2;
  flex-shrink: 0;
}

.license-code-group {
  display: flex;
  align-items: center;
  gap: 16px;
}

.license-code {
  
  font-size: 18px;
  font-weight: 700;
  line-height: 1.5;
  color: #1D1D1D;
  margin: 0;
}

.status-tag {
  padding: 7px 16px;
  font-size: 14px;
  font-weight: 500;
  line-height: 1.5;
  border-radius: 4px;
}

.action-buttons-inline {
  display: flex;
  gap: 4px;
}

.action-btn-inline {
  padding: 7px 16px;
  font-size: 14px;
  font-weight: 500;
  line-height: 1.5;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s;
  
  height: 32px;

  &.copy-btn {
    background: transparent;
    color: #019C7C;
    border: 1px solid #019C7C;

    &:hover {
      background: rgba(1, 156, 124, 0.1);
    }
  }

  &.update-btn {
    background: transparent;
    color: #019C7C;
    border: 1px solid #019C7C;

    &:hover {
      background: rgba(1, 156, 124, 0.1);
    }
  }

  &.renew-btn {
    background: #00C27C;
    color: white;
    border: none;

    &:hover {
      background: #019C7C;
    }
  }

  &.download-btn {
    background: #00C27C;
    color: white;
    border: none;

    &:hover {
      background: #019C7C;
    }
  }
}

.tab-content-area {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

@media (max-width: 1024px) {
  .content-section {
    flex-direction: column;
  }

  .left-sidebar {
    width: 100%;
    display: flex;
    padding: 8px;
    overflow-x: auto;

    .sidebar-item {
      flex-shrink: 0;
      padding: 8px 16px;
      white-space: nowrap;

      &.active {
        border-left: none;
        border-bottom: 3px solid var(--el-color-primary);
      }
    }
  }
}

@media (max-width: 768px) {
  .license-view-page {
    padding: 16px;
  }

  .license-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .action-buttons-inline {
    width: 100%;
    flex-wrap: wrap;

    .action-btn-inline {
      flex: 1 1 calc(50% - 2px);
      min-width: 100px;
    }
  }
}
</style>

