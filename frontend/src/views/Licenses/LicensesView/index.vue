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
            <el-button class="action-btn-inline update-btn" @click="handleUpdateLicense">
              {{ $t('pages.licenses.detail.actions.updateLicense') }}
            </el-button>
            <el-button class="action-btn-inline renew-btn" @click="handleRenewLicense">
              {{ $t('pages.licenses.detail.actions.renewLicense') }}
            </el-button>
            <el-button class="action-btn-inline revoke-btn" @click="handleRevokeLicense">
              {{ $t('pages.licenses.detail.actions.revokeLicense') }}
            </el-button>
            <el-button class="action-btn-inline download-btn" @click="handleDownloadCertificate">
              {{ $t('pages.licenses.detail.actions.downloadCertificate') }}
            </el-button>
          </div>
        </div>

        <!-- Tab 内容区域 -->
        <div class="tab-content-area">
          <BasicInfo v-if="activeTab === 'basic'" :license-data="licenseData" />
          <AuthorizationInfo v-if="activeTab === 'authorization'" :license-data="licenseData" />
          <LicenseInfo v-if="activeTab === 'license'" :license-data="licenseData" />
          <ChangeHistory v-if="activeTab === 'history'" :history-data="historyData" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { getLicenseDetail, type AuthorizationCode } from '@/api/license'
import LicenseTabIcon from '@/components/common/icons/LicenseTabIcon.vue'
import BasicInfo from './components/BasicInfo.vue'
import AuthorizationInfo from './components/AuthorizationInfo.vue'
import LicenseInfo from './components/LicenseInfo.vue'
import ChangeHistory from './components/ChangeHistory.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()

const loading = ref(false)
const activeTab = ref('basic')
const licenseData = ref<AuthorizationCode | null>(null)
const historyData = ref<any[]>([])

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

const handleUpdateLicense = () => {
  ElMessage.info(t('pages.licenses.detail.messages.updateComingSoon'))
}

const handleRenewLicense = () => {
  ElMessage.info(t('pages.licenses.detail.messages.renewComingSoon'))
}

const handleRevokeLicense = () => {
  ElMessageBox.confirm(t('pages.licenses.detail.messages.revokeConfirm'), t('pages.licenses.detail.messages.revokeTitle'), {
    confirmButtonText: t('pages.licenses.detail.messages.revokeConfirmButton'),
    cancelButtonText: t('pages.licenses.detail.messages.revokeCancelButton'),
    type: 'warning'
  }).then(() => {
    ElMessage.info(t('pages.licenses.detail.messages.revokeComingSoon'))
  }).catch(() => {
    // 取消操作
  })
}

const handleDownloadCertificate = () => {
  ElMessage.info(t('pages.licenses.detail.messages.downloadComingSoon'))
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
    font-family: 'Source Han Sans CN', sans-serif;
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
    justify-content: center;
    gap: 12px;
    padding: 16px 0;
    margin: 8px 16px;
    cursor: pointer;
    transition: all 0.3s;
    color: #B2B8C2;
    background: transparent;
    border-radius: 28px;

    .text {
      font-family: 'Source Han Sans CN', sans-serif;
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
  font-family: 'Source Han Sans CN', sans-serif;
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
  font-family: 'Source Han Sans CN', sans-serif;
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

  &.revoke-btn {
    background: #F0142F;
    color: white;
    border: none;

    &:hover {
      background: #d01228;
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
