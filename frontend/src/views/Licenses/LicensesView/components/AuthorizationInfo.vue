<template>
  <div class="authorization-info-container">
    <!-- 基本信息卡片 -->
    <div class="info-card">
      <div class="card-title">{{ t('pages.licenses.detail.authorizationInfo.basicInfo') }}</div>
      <div class="info-content">
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.customerName') }}：</div>
          <div class="info-value">{{ licenseData?.customer_info?.customer_name || '-' }}</div>
        </div>
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.customerCode') }}：</div>
          <div class="info-value">{{ licenseData?.customer_info?.customer_code || '-' }}</div>
        </div>
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.description') }}：</div>
          <div class="info-value">{{ licenseData?.description || '-' }}</div>
        </div>
      </div>
    </div>

    <!-- 授权信息卡片 -->
    <div class="info-card">
      <div class="card-title">{{ t('pages.licenses.detail.authorizationInfo.title') }}</div>
      <div class="info-content">
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.validityPeriodLabel') }}：</div>
          <div class="info-value">{{ getLicensePeriodType() }}</div>
        </div>
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.validityTime') }}：</div>
          <div class="info-value">{{ formatDateRange() }}</div>
        </div>
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.maxActivations') }}：</div>
          <div class="info-value">{{ licenseData?.max_activations || 0 }}{{ t('pages.licenses.detail.authorizationInfo.devices') }}</div>
        </div>
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.deploymentType') }}：</div>
          <div class="info-value">{{ licenseData?.deployment_type_display || '-' }}</div>
        </div>
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.encryptionType') }}：</div>
          <div class="info-value">{{ licenseData?.encryption_type_display || t('pages.licenses.detail.authorizationInfo.standardEncryption') }}</div>
        </div>
      </div>
    </div>

    <!-- 授权码卡片 -->
    <div class="info-card authorization-code-card">
      <div class="card-title-row">
        <div class="card-title">{{ t('pages.licenses.detail.authorizationInfo.authorizationCode') }}</div>
        <div class="card-subtitle">{{ t('pages.licenses.detail.authorizationInfo.forActivation') }}</div>
      </div>
      <div class="info-content">
        <div class="info-row">
          <div class="info-label">{{ t('pages.licenses.detail.authorizationInfo.code') }}：</div>
          <div class="info-value code-value">{{ licenseData?.code || '-' }}</div>
        </div>
        <div class="info-tip">
          <svg class="tip-icon" width="20" height="20" viewBox="0 0 20 20" fill="none">
            <path d="M10 0C4.48 0 0 4.48 0 10C0 15.52 4.48 20 10 20C15.52 20 20 15.52 20 10C20 4.48 15.52 0 10 0ZM11 15H9V9H11V15ZM11 7H9V5H11V7Z" fill="#4763FF"/>
          </svg>
          <span class="tip-text">{{ t('pages.licenses.detail.authorizationInfo.tipText') }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import type { AuthorizationCode } from '@/api/license'
import { formatDateShort as formatDateUtil } from '@/utils/date'

interface Props {
  licenseData: AuthorizationCode | null
}

const props = defineProps<Props>()
const { t } = useI18n()

// 获取授权期限类型
const getLicensePeriodType = () => {
  // 这里假设有个字段表示授权期限类型，如果没有则返回"有效期限"
  return t('pages.licenses.detail.authorizationInfo.validityPeriod')
}

// 格式化日期区间
const formatDateRange = () => {
  if (!props.licenseData?.start_date || !props.licenseData?.end_date) {
    return '-'
  }
  const startDate = formatDateUtil(props.licenseData.start_date)
  const endDate = formatDateUtil(props.licenseData.end_date)
  return `${startDate} ${t('pages.licenses.detail.authorizationInfo.to')} ${endDate}`
}
</script>

<style lang="scss" scoped>
.authorization-info-container {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-card {
  background: #FFFFFF;
  border: 1px solid #E2E2E2;
  border-radius: 4px;
  padding: 16px 20px;

  .card-title {
    font-family: 'Source Han Sans CN', sans-serif;
    font-size: 16px;
    font-weight: 700;
    line-height: 18px;
    color: #1D1D1D;
    margin-bottom: 16px;
  }

  .card-title-row {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-bottom: 16px;

    .card-title {
      margin-bottom: 0;
    }

    .card-subtitle {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 18px;
      color: #666666;
    }
  }

  .info-content {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
    gap: 16px 24px;
  }

  &.authorization-code-card {
    .info-content {
      display: flex;
      flex-direction: column;
      gap: 16px;
    }
  }

  .info-row {
    display: flex;
    align-items: flex-start;
    gap: 8px;

    .info-label {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      color: #1D1D1D;
      white-space: nowrap;
    }

    .info-value {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      color: #1D1D1D;
      flex: 1;

      &.code-value {
        font-weight: 500;
      }
    }
  }

  .info-tip {
    display: flex;
    align-items: flex-start;
    gap: 8px;

    .tip-icon {
      flex-shrink: 0;
      margin-top: 1px;
    }

    .tip-text {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      color: #4763FF;
      flex: 1;
    }
  }
}

// 响应式布局
@media (max-width: 1024px) {
  .info-card {
    .info-content {
      grid-template-columns: 1fr 1fr;
      gap: 16px 20px;
    }
  }
}

@media (max-width: 768px) {
  .info-card {
    .info-content {
      grid-template-columns: 1fr;
      gap: 16px;
    }
    
    .info-row {
      flex-direction: column;
      gap: 4px;
    }
  }
}
</style>
