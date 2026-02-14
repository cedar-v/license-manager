<template>
  <div class="basic-info-tab">
    <!-- 基础信息部分 -->
    <div class="info-section">
      <div class="info-row">
        <div class="info-label">{{ $t('customers.basicInfo.customerCode') }}</div>
        <div class="info-value">{{ licenseData?.customer_info?.customer_code || '-' }}</div>
      </div>
      <div class="info-row">
        <div class="info-label">{{ $t('customers.basicInfo.customerName') }}</div>
        <div class="info-value">{{ licenseData?.customer_info?.customer_name || '-' }}</div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <!-- 已激活设备数量 -->
      <div class="stat-card stat-card-1">
        <div class="card-content">
          <div class="stat-value">{{ getActivatedLicenses(licenseData) }}</div>
          <div class="stat-label">{{ $t('customers.basicInfo.activatedDevices') }}</div>
        </div>
      </div>

      <!-- 创建时间 -->
      <div class="stat-card stat-card-2">
        <div class="card-content">
          <div class="stat-value">{{ formatDateShort(licenseData?.start_date) }}</div>
          <div class="stat-label">{{ $t('customers.basicInfo.creationTime') }}</div>
        </div>
      </div>

      <!-- 到期时间 -->
      <div class="stat-card stat-card-3">
        <div class="card-content">
          <div class="stat-value">{{ formatDateShort(licenseData?.end_date) }}</div>
          <div class="stat-label">{{ $t('customers.basicInfo.expiryTime') }}</div>
        </div>
      </div>

      <!-- 剩余时间 -->
      <div class="stat-card stat-card-4">
        <div class="card-content">
          <div class="stat-value">{{ formatRemainingTime(licenseData?.end_date) }}</div>
          <div class="stat-label">{{ $t('customers.basicInfo.remainingDays') }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AuthorizationCode } from '@/api/license'
import { formatDateShort as formatDateUtil } from '@/utils/date'
import { useI18n } from 'vue-i18n'

interface Props {
  licenseData: AuthorizationCode | null
}

defineProps<Props>()
const { t } = useI18n()

const getActivatedLicenses = (license?: AuthorizationCode | null) => {
  if (!license) return 0
  if (typeof license.current_activations === 'number') {
    return license.current_activations
  }
  if (typeof license.activated_licenses_count === 'number') {
    return license.activated_licenses_count
  }
  return 0
}

const formatDateShort = (date?: string) => {
  if (!date) return '-'
  return formatDateUtil(date)
}

const formatRemainingTime = (endDate?: string) => {
  if (!endDate) return '-'
  const end = new Date(endDate)
  if (Number.isNaN(end.getTime())) return '-'
  const now = new Date()
  const diffMs = end.getTime() - now.getTime()
  const days = Math.max(0, Math.ceil(diffMs / (1000 * 60 * 60 * 24)))

  if (days > 36500) {
    return t('customers.basicInfo.remainingTime.permanent')
  }

  if (days >= 731) {
    const years = Math.floor(days / 365)
    return t('customers.basicInfo.remainingTime.years', { years })
  }

  if (days >= 366) {
    const years = Math.floor(days / 365)
    const remainingDays = days - years * 365
    if (remainingDays === 0) {
      return t('customers.basicInfo.remainingTime.years', { years })
    }
    const months = Math.max(1, Math.round(remainingDays / 30))
    return t('customers.basicInfo.remainingTime.yearMonth', { years, months })
  }

  return t('customers.basicInfo.remainingTime.days', { days })
}
</script>

<style lang="scss" scoped>
.basic-info-tab {
  width: 100%;
  padding: 0;
}

.info-section {
  margin-bottom: 32px;

  .info-row {
    display: flex;
    align-items: center;
    gap: 40px;
    padding: 8px 0;

    .info-label {
      width: 80px;
      
      font-size: 14px;
      font-weight: 400;
      line-height: 2.2857142857142856em;
      color: rgba(29, 29, 29, 0.87);
      text-align: right;
      flex-shrink: 0;
    }

    .info-value {
      flex: 1;
      
      font-size: 14px;
      font-weight: 500;
      line-height: 1.5em;
      color: rgba(29, 29, 29, 0.87);
    }
  }
}

.stats-cards {
  display: flex;
  gap: 16px;
  margin-top: 24px;
  flex-wrap: wrap;

  @media (max-width: 1200px) {
    flex-wrap: wrap;
  }

  @media (max-width: 768px) {
    flex-direction: column;
  }
}

.stat-card {
  position: relative;
  flex: 1;
  max-width: 280px;
  height: 192px;
  border-radius: 8px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 40px 24px;
  box-sizing: border-box;
  background-image: url('@/assets/images/licenseIcon1.png');
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  background-color: #ffffff;

  @media (max-width: 768px) {
    min-width: 100%;
    flex: none;
  }

  .card-content {
    position: relative;
    z-index: 1;
    display: flex;
    flex-direction: column;
    gap: 8px;
    max-width: 100%;
    overflow: hidden;

    .stat-value {
      font-family: 'DINPro', sans-serif;
      font-size: 30px;
      font-weight: 700;
      line-height: 1.2;
      color: #2C2F33;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }

    .stat-label {
      
      font-size: 14px;
      font-weight: 400;
      line-height: 1.3;
      color: rgba(99, 110, 131, 1);
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
  }

  // 为每个卡片设置不同的背景图片
  &.stat-card-1 {
    background-image: url('@/assets/images/licenseIcon1.png');
  }

  &.stat-card-2 {
    background-image: url('@/assets/images/licenseIcon2.png');
  }

  &.stat-card-3 {
    background-image: url('@/assets/images/licenseIcon3.png');
  }

  &.stat-card-4 {
    background-image: url('@/assets/images/licenseIcon4.png');
  }
}
</style>
