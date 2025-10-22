<template>
  <div class="basic-info-tab">
    <!-- 基础信息部分 -->
    <div class="info-section">
      <div class="info-row">
        <div class="info-label">客户ID</div>
        <div class="info-value">{{ licenseData?.customer_id || '-' }}</div>
      </div>
      <div class="info-row">
        <div class="info-label">客户名称</div>
        <div class="info-value">{{ licenseData?.customer_name || '-' }}</div>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-cards">
      <!-- 已激活设备数量 -->
      <div class="stat-card">
        <div class="card-background">
          <img src="@/assets/images/licenseIcon1.png" alt="" class="bg-image" />
        </div>
        <div class="card-content">
          <div class="stat-value">{{ licenseData?.current_activations || 0 }}</div>
          <div class="stat-label">已激活设备数量</div>
        </div>
      </div>

      <!-- 创建时间 -->
      <div class="stat-card">
        <div class="card-background">
          <img src="@/assets/images/licenseIcon2.png" alt="" class="bg-image" />
        </div>
        <div class="card-content">
          <div class="stat-value">{{ formatDate(licenseData?.start_date) }}</div>
          <div class="stat-label">创建时间</div>
        </div>
      </div>

      <!-- 到期时间 -->
      <div class="stat-card">
        <div class="card-background">
          <img src="@/assets/images/licenseIcon3.png" alt="" class="bg-image" />
        </div>
        <div class="card-content">
          <div class="stat-value">{{ formatDate(licenseData?.end_date) }}</div>
          <div class="stat-label">到期时间</div>
        </div>
      </div>

      <!-- 剩余时间（天） -->
      <div class="stat-card">
        <div class="card-background">
          <img src="@/assets/images/licenseIcon4.png" alt="" class="bg-image" />
        </div>
        <div class="card-content">
          <div class="stat-value">{{ getRemainingDays(licenseData?.end_date) }}</div>
          <div class="stat-label">剩余时间（天）</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { AuthorizationCode } from '@/api/license'
import { formatDate as formatDateUtil } from '@/utils/date'

interface Props {
  licenseData: AuthorizationCode | null
}

defineProps<Props>()

const formatDate = (date?: string) => {
  if (!date) return '-'
  return formatDateUtil(date)
}

const getRemainingDays = (endDate?: string) => {
  if (!endDate) return '-'
  const now = new Date()
  const end = new Date(endDate)
  const diff = end.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))
  return days > 0 ? days : 0
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
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 2.2857142857142856em;
      color: rgba(29, 29, 29, 0.87);
      text-align: right;
      flex-shrink: 0;
    }

    .info-value {
      flex: 1;
      font-family: 'Source Han Sans CN', sans-serif;
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
  min-width: 250px;
  height: 193.58px;
  border-radius: 11.504px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  padding: 12px 25px;
  box-sizing: border-box;

  @media (max-width: 768px) {
    min-width: 100%;
    flex: none;
  }

  .card-background {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 0;

    .bg-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .card-content {
    position: relative;
    z-index: 1;
    display: flex;
    flex-direction: column;
    gap: 12px;

    .stat-value {
      font-family: 'DINPro', sans-serif;
      font-size: 30px;
      font-weight: 700;
      line-height: 1.287999979654948em;
      color: #2C2F33;
    }

    .stat-label {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 1.307861600603376em;
      color: rgba(99, 110, 131, 1);
    }
  }
}
</style>
