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

    <!-- 功能配置 -->
    <div class="info-card">
      <div class="card-title">{{ t('pages.licenses.detail.authorizationInfo.featureConfig') }}</div>
      <div v-if="featureConfig.length" class="key-value-table">
        <div class="key-value-header">
          <span>{{ t('pages.licenses.detail.authorizationInfo.keyColumn') }}</span>
          <span>{{ t('pages.licenses.detail.authorizationInfo.typeColumn') }}</span>
          <span>{{ t('pages.licenses.detail.authorizationInfo.valueColumn') }}</span>
        </div>
        <div v-for="item in featureConfig" :key="`feature-${item.key}`" class="key-value-row">
          <span class="key-cell">{{ item.key }}</span>
          <span class="type-cell">{{ t(`pages.licenses.detail.authorizationInfo.typeOptions.${item.type}`) }}</span>
          <span class="value-cell">{{ item.value || '-' }}</span>
        </div>
      </div>
      <div v-else class="key-value-empty">
        {{ t('pages.licenses.detail.authorizationInfo.keyValueEmpty') }}
      </div>
    </div>

    <!-- 使用限制 -->
    <div class="info-card">
      <div class="card-title">{{ t('pages.licenses.detail.authorizationInfo.usageLimits') }}</div>
      <div v-if="usageLimits.length" class="key-value-table">
        <div class="key-value-header">
          <span>{{ t('pages.licenses.detail.authorizationInfo.keyColumn') }}</span>
          <span>{{ t('pages.licenses.detail.authorizationInfo.typeColumn') }}</span>
          <span>{{ t('pages.licenses.detail.authorizationInfo.valueColumn') }}</span>
        </div>
        <div v-for="item in usageLimits" :key="`limit-${item.key}`" class="key-value-row">
          <span class="key-cell">{{ item.key }}</span>
          <span class="type-cell">{{ t(`pages.licenses.detail.authorizationInfo.typeOptions.${item.type}`) }}</span>
          <span class="value-cell">{{ item.value || '-' }}</span>
        </div>
      </div>
      <div v-else class="key-value-empty">
        {{ t('pages.licenses.detail.authorizationInfo.keyValueEmpty') }}
      </div>
    </div>

    <!-- 自定义参数 -->
    <div class="info-card">
      <div class="card-title">{{ t('pages.licenses.detail.authorizationInfo.customParameters') }}</div>
      <div v-if="customParameters.length" class="key-value-table">
        <div class="key-value-header">
          <span>{{ t('pages.licenses.detail.authorizationInfo.keyColumn') }}</span>
          <span>{{ t('pages.licenses.detail.authorizationInfo.typeColumn') }}</span>
          <span>{{ t('pages.licenses.detail.authorizationInfo.valueColumn') }}</span>
        </div>
        <div v-for="item in customParameters" :key="`custom-${item.key}`" class="key-value-row">
          <span class="key-cell">{{ item.key }}</span>
          <span class="type-cell">{{ t(`pages.licenses.detail.authorizationInfo.typeOptions.${item.type}`) }}</span>
          <span class="value-cell">{{ item.value || '-' }}</span>
        </div>
      </div>
      <div v-else class="key-value-empty">
        {{ t('pages.licenses.detail.authorizationInfo.keyValueEmpty') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { AuthorizationCode } from '@/api/license'
import { formatDateShort as formatDateUtil } from '@/utils/date'

interface Props {
  licenseData: AuthorizationCode | null
}

const props = defineProps<Props>()
const { t } = useI18n()

type KeyValueType = 'string' | 'number' | 'bool'

interface KeyValueItem {
  key: string
  value: string
  type: KeyValueType
}

const parseKeyValueData = (source: unknown): KeyValueItem[] => {
  if (!source) return []

  let parsed: Record<string, any> | null = null
  if (typeof source === 'string') {
    if (!source.trim()) return []
    try {
      parsed = JSON.parse(source)
    } catch (error) {
      console.warn('[AuthorizationInfo] Failed to parse JSON:', error)
      return []
    }
  } else if (typeof source === 'object') {
    parsed = source as Record<string, any>
  }

  if (!parsed) return []

  return Object.entries(parsed).map(([key, value]) => {
    let type: KeyValueType = 'string'
    let normalizedValue = ''
    if (typeof value === 'number') {
      type = 'number'
      normalizedValue = String(value)
    } else if (typeof value === 'boolean') {
      type = 'bool'
      normalizedValue = value ? 'true' : 'false'
    } else {
      normalizedValue = value === undefined || value === null ? '' : String(value)
    }

    return {
      key,
      value: normalizedValue,
      type
    }
  })
}

const featureConfig = computed(() => parseKeyValueData(props.licenseData?.feature_config))
const usageLimits = computed(() => parseKeyValueData(props.licenseData?.usage_limits))
const customParameters = computed(() => parseKeyValueData(props.licenseData?.custom_parameters))

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
      
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      color: #1D1D1D;
      white-space: nowrap;
    }

    .info-value {
      
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
      
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      color: #4763FF;
      flex: 1;
    }
  }

  .key-value-table {
    border: 1px solid #eef0f2;
    border-radius: 4px;
    overflow: hidden;

    .key-value-header,
    .key-value-row {
      display: grid;
      grid-template-columns: 200px 140px 1fr;
      padding: 10px 16px;
      gap: 12px;
    }

    .key-value-header {
      background: #f5f6f8;
      font-weight: 600;
      font-size: 13px;
      color: #1d1d1d;
    }

    .key-value-row {
      font-size: 13px;
      color: #1d1d1d;
      border-top: 1px solid #eef0f2;

      .key-cell {
        font-weight: 600;
      }

      .type-cell {
        color: #606266;
        font-weight: 500;
      }

      .value-cell {
        word-break: break-all;
      }
    }
  }

  .key-value-empty {
    border: 1px dashed #d7dbe2;
    border-radius: 4px;
    padding: 12px 16px;
    color: #909399;
    font-size: 13px;
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
