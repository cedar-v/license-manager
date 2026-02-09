<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-01 09:32:42
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-11-04 15:19:44
 * @FilePath: /frontend/src/views/Dashboard.vue
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
<template>
  <Layout app-name="Cedar-V" :page-title="t('dashboard.title')">
    <!-- 页面内容 -->
    <div class="content-container dashboard">
      <!-- 统计卡片区域 -->
      <div class="stats-section">
        <div
          class="stats-card"
          :class="{ 'compact-card': stat.id === 6, 'large-card': stat.id === 7 }"
          v-for="stat in statsData"
          :key="stat.id"
        >
          <div class="stat-icon">
            <div class="icon-circle">
              <img :src="stat.icon" :alt="stat.label" class="stat-icon-img" />
            </div>
          </div>
          <div class="stat-content">
            <div class="stat-label">{{ stat.label }}</div>
            <div class="stat-value-row">
              <div class="stat-value">{{ stat.value }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 图表和表格区域 -->
      <div class="content-section">
        <!-- 授权趋势图表-->
        <div class="chart-section">
          <LicenseTrendChart />
        </div>

        <!-- 最近授权表格 - 占据60%高度 -->
        <div class="table-card">
          <div class="card-header">
            <h3 class="card-title">{{ t('dashboard.recentLicenses.title') }}</h3>
          </div>
          <div class="table-container">
            <el-table
              :data="recentData"
              :loading="loading"
              style="width: 100%; height: 100%"
              :header-row-class-name="'table-header'"
              :row-class-name="(params: any) => (params.rowIndex % 2 === 1 ? 'stripe-row' : '')"
            >
              <el-table-column
                :label="t('dashboard.recentLicenses.columns.serialNumber')"
                min-width="90"
                align="center"
              >
                <template #default="{ $index }">
                  {{ $index + 1 }}
                </template>
              </el-table-column>
              <el-table-column
                prop="customer_name"
                :label="t('dashboard.recentLicenses.columns.customerName')"
                min-width="200"
              />
              <el-table-column
                prop="description"
                :label="t('dashboard.recentLicenses.columns.description')"
                min-width="150"
              />
              <el-table-column
                :label="t('dashboard.recentLicenses.columns.status')"
                width="120"
                align="center"
              >
                <template #default="{ row }">
                  <span
                    class="status-badge"
                    :class="getStatusClass(row.status)"
                  >
                    {{ row.status_display }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column
                :label="t('dashboard.recentLicenses.columns.expiryTime')"
                width="301"
              >
                <template #default="{ row }">
                  {{ formatDate(row.end_date) }}
                </template>
              </el-table-column>
              <el-table-column
                :label="t('dashboard.recentLicenses.columns.createTime')"
                width="301"
              >
                <template #default="{ row }">
                  {{ formatDate(row.created_at) }}
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/common/layout/Layout.vue'
import LicenseTrendChart from '@/components/charts/LicenseTrendChart.vue'
import { getRecentAuthorizations, type RecentAuthorizationItem, getOverviewStats} from '@/api/dashboard'
import { ElMessage } from 'element-plus'
import { formatDate } from '@/utils/date'

// 导入dashboard目录中的图标
import icon1 from '@/assets/icons/dashboard/icon1.png'
import icon2 from '@/assets/icons/dashboard/icon2.png'
import icon3 from '@/assets/icons/dashboard/icon3.png'
import icon4 from '@/assets/icons/dashboard/icon4.png'
import icon6 from '@/assets/icons/dashboard/icon6.png'
import icon7 from '@/assets/icons/dashboard/icon7.png'

// 使用国际化
const { t } = useI18n()

// 统计卡片数据（接口动态获取 value，其余静态配置）
const statsData = ref([

  { id: 1, key: 'total_auth_codes', value: '', label: t('dashboard.stats.totalCustomers'), icon: icon1 },//总授权码
  { id: 2, key: 'active_licenses', value: '', label: t('dashboard.stats.totalLicenses'), icon: icon2 },//活跃许可证
  { id: 3, key: 'expiring_soon', value: '', label: t('dashboard.stats.activeLicenses'), icon: icon3 },//即将过期
  { id: 4, key: 'abnormal_alerts', value: '', label: t('dashboard.stats.expired'), icon: icon4 },//异常警告
  // { id: 5, key: 'growth_rate', value: '', label: t('dashboard.stats.onlineClients'), icon: icon5 },//增长率
  { id: 6, key: 'auth_codes', value: '', label: t('dashboard.stats.expiringSoon'), icon: icon6 },//授权码增长率
  { id: 7, key: 'licenses', value: '', label: t('dashboard.stats.authSuccessRate24h'), icon: icon7 }//许可证增长率
])

// 最近授权数据 - 响应式数据
const recentData = ref<RecentAuthorizationItem[]>([])
const loading = ref(false)

// 获取状态样式类名
const getStatusClass = (status: string) => {
  switch (status) {
    case 'normal':
      return 'status-valid'
    case 'locked':
      return 'status-locked'
    case 'expired':
      return 'status-invalid'
    default:
      return 'status-invalid'
  }
}

// 获取最近授权数据
const fetchRecentAuthorizations = async () => {
  try {
    loading.value = true
    const response = await getRecentAuthorizations()
    console.log('最近授权数据:', response.data)
    recentData.value = response.data.list
  } catch (error: any) {
    console.error('获取最近授权数据失败:', error)
    ElMessage.error(error?.backendMessage || '获取最近授权数据失败')
  } finally {
    loading.value = false
  }
}

// 获取仪表盘统计数据
const fetchDashboardStats = async () => {
  try {
    const resp = await getOverviewStats()
    const data = resp.data
    statsData.value.forEach(item => {
      // 处理直接映射的字段
      if (data[item.key] !== undefined && typeof data[item.key] !== 'object') {
        item.value = data[item.key]
      }
      // 处理嵌套在 growth_rate 中的字段（授权码增长率和许可证增长率）
      if (data.growth_rate && data.growth_rate[item.key] !== undefined) {
        const rateValue = data.growth_rate[item.key]
        // 如果是数字类型，格式化为保留两位小数的百分比
        if (typeof rateValue === 'number') {
          item.value = `${rateValue.toFixed(2)}%`
        } else {
          item.value = rateValue
        }
      }
    })
  } catch (e) {
    ElMessage.error('获取仪表盘统计数据失败')
  }
}

// 组件挂载时获取数据
onMounted(() => {
  console.log('Dashboard 组件已挂载')
  fetchDashboardStats()
  fetchRecentAuthorizations()
})
</script>

<style lang="scss" scoped>
/* content-container样式 */
.content-container {
  min-height: calc(100vh - 80px);
  padding: 24px;
  width: 100%;
  margin: 0;
  box-sizing: border-box;
  background: var(--app-bg-color);
}

.dashboard {
  width: 100%;
  height: 100%; /* 使用100%高度充满父容器 */
  display: flex;
  flex-direction: column;
}

// 统计卡片区域 - 基于1920*1080设计的vw适配
.stats-section {
  display: flex;
  gap: 2.08vw; /* 卡片之间的间距 40px/1920 = 2.08vw */
  margin-bottom: 1.25vw; /* 24px/1920 = 1.25vw */
  padding: 1.25vw; /* 24px/1920 = 1.25vw */
  background: linear-gradient(135deg, #019C7C 0%, #1db584 100%);
  border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
  position: relative;
  overflow: hidden; /* 确保背景图片不溢出 */

  // 背景图片层
  &::after {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: url('@/assets/images/dashboard-bg.png') lightgray 50% / cover no-repeat;
    mix-blend-mode: soft-light;
    border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
    z-index: 0;
  }

  // 网格纹理层
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background:
      linear-gradient(45deg, rgba(255, 255, 255, 0.08) 25%, transparent 25%),
      linear-gradient(-45deg, rgba(255, 255, 255, 0.08) 25%, transparent 25%),
      linear-gradient(45deg, transparent 75%, rgba(255, 255, 255, 0.08) 75%),
      linear-gradient(-45deg, transparent 75%, rgba(255, 255, 255, 0.08) 75%);
    background-size: 1.04vw 1.04vw; /* 20px/1920 = 1.04vw */
    background-position:
      0 0,
      0 0.52vw,
      0.52vw -0.52vw,
      -0.52vw 0vw; /* 10px/1920 = 0.52vw */
    opacity: 0.15;
    border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
    z-index: 1;
  }
}

.stats-card {
  flex: 1; /* 等宽分布 */
  display: flex;
  align-items: center;
  justify-content: center; /* 内容居中 */
  gap: 0.83vw; /* 16px/1920 = 0.83vw */
  position: relative;
  z-index: 2; /* 确保内容显示在背景图片和网格纹理之上 */
  min-width: 0;

  &:not(:last-child)::after {
    content: '';
    position: absolute;
    left: 100%; /* 从卡片右边缘开始 */
    top: 50%;
    transform: translate(calc(2.08vw / 2), -50%); /* 向右移动gap的一半，垂直居中 */
    width: 1px;
    height: 3.33vw; /* 64px/1920 = 3.33vw */
    background: linear-gradient(
      to bottom,
      rgba(255, 255, 255, 0.08) 0%,
      rgba(255, 255, 255, 0.6) 49%,
      rgba(255, 255, 255, 0.08) 100%
    );
  }
}

.stat-icon {
  display: flex;
  align-items: center;
  height: 100%;

  .icon-circle {
    height: 100%;
    aspect-ratio: 1; /* 保持圆形 */
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;

    .stat-icon-img {
      width: 50%; /* icon-circle 的一半 */
      height: 50%; /* icon-circle 的一半 */
      filter: brightness(0) invert(1) opacity(0.9); /* 将图标转为白色半透明 */
      object-fit: contain;
    }
  }
}

.stat-content {
  flex: 0 1 auto;
  min-width: 4.17vw; /* 80px/1920 = 4.17vw */
  display: flex;
  flex-direction: column;
  justify-content: center;
  height: 100%;

  .stat-label {
    font-size: 0.83vw; /* 16px/1920 = 0.83vw */
    font-weight: 500;
    color: #ffffff;
    line-height: 1;
    white-space: nowrap;
  }

  .stat-value-row {
    display: flex;
    align-items: center;
    gap: 0.42vw; /* 8px/1920 = 0.42vw */
    margin-top: 0.21vw; /* 4px/1920 = 0.21vw */
  }

  .stat-value {
    font-size: 1.25vw; /* 24px/1920 = 0.73vw */
    color: rgba(255, 255, 255, 0.9);
  }
}

// 内容区域 - 默认上下布局，充满屏幕宽度和高度
.content-section {
  display: flex;
  flex-direction: column;
  gap: 1.25vw; /* 24px/1920 = 1.25vw */
  flex: 1; /* 占据剩余空间 */
  min-height: 0; /* 防止flex子元素溢出 */
}

// 图表区域 - 占据50%高度
.chart-section {
  flex: 0.5; /* 50%高度 */
  min-height: 300px; /* 确保最小高度 */
  display: flex;
  flex-direction: column;
}

// 卡片通用样式 - 表格卡片占据60%高度
.table-card {
  background: var(--app-content-bg);
  border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
  box-shadow: var(--app-shadow);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  flex: 0.6; /* 占据60%高度 */
  min-height: 0; /* 防止flex子元素溢出 */
}

.card-header {
  padding: 24px 24px 0; /* 调整为24px符合4px倍数 */

  .card-title {
    font-size: 1.04vw; /* 20px/1920 = 1.04vw */
    font-weight: 500; /* 调整为中等粗细 */
    color: var(--app-text-primary);
    margin: 0;
    line-height: 1.3;
    font-family: 'Source Han Sans CN', sans-serif;
  }
}

// 表格卡片
.table-container {
  padding: 24px; /* 使用24px统一padding，符合4px倍数 */
  flex: 1; /* 占据剩余高度 */
  overflow: auto; /* 如果内容过多，允许滚动 */
  display: flex;
  flex-direction: column;

  // Element Plus 表格样式重写
  :deep(.el-table) {
    border: 1px solid var(--app-border-light);
    height: 100%; /* 充满容器高度 */

    .el-table__body-wrapper {
      flex: 1; /* 表格主体占据剩余高度 */
    }

    .table-header th {
      background-color: var(--app-bg-color) !important;
      color: var(--app-text-primary);
      font-weight: 500;
      font-size: 16px;
      font-family: 'Source Han Sans CN', sans-serif;
      height: 48px;
      padding: 12px 20px; /* 调整为12px符合4px倍数 */
      border-bottom: none;
    }

    .el-table__body tr {
      height: 48px;

      td {
        padding: 12px 20px; /* 调整为12px符合4px倍数 */
        border-bottom: 1px solid var(--app-border-light);
        font-size: 14px;
        font-family: 'Source Han Sans CN', sans-serif;
        font-weight: 400; /* 调整为正常粗细 */
        color: var(--app-text-primary);
      }

      &.stripe-row {
        background-color: var(--app-bg-color);
      }

      &:hover > td {
        background-color: var(--app-bg-color) !important;
      }
    }

    // 序号列居中
    .el-table__body tr td:nth-child(1) {
      text-align: center;
    }

    // 状态列居中
    .el-table__body tr td:nth-child(4) {
      text-align: center;
      overflow: visible;
      text-overflow: initial;
    }

    // 其他列左对齐（客户名称、描述、时间等）
    .el-table__body tr td:nth-child(2),
    .el-table__body tr td:nth-child(3),
    .el-table__body tr td:nth-child(5),
    .el-table__body tr td:nth-child(6) {
      text-align: left;
    }
  }
}

.status-badge {
  padding: 4px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-weight: 400;
  font-family: 'Source Han Sans CN', sans-serif;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  white-space: nowrap;
  width: 46px;
  height: 22px;

  &.status-valid {
    background: #f0f5ff;
    color: #4763ff;
  }

  &.status-locked {
    background: #fff7e6;
    color: #d46b08;
  }

  &.status-invalid {
    background: #ffe5e5;
    color: #e90c0c;
  }
}

// 响应式设计 - 移动端切换回px单位确保可读性
@media (max-width: 768px) {
  .stats-section {
    display: grid; /* 在移动端切换回 grid 布局 */
    grid-template-columns: repeat(2, 1fr);
    gap: 16px; /* 移动端使用固定像素 */
    padding: 16px;
    border-radius: 8px;
  }

  .stats-card {
    flex-direction: column;
    text-align: center;
    gap: 8px;
    padding: 16px 0;

    &:not(:last-child)::after {
      display: none;
    }
  }

  .stat-icon {
    height: auto;
  }

  .stat-icon .icon-circle {
    height: 100%;
    aspect-ratio: 1; /* 保持圆形 */
    min-width: 36px; /* 移动端最小大小 */
    min-height: 36px;

    .stat-icon-img {
      width: 50%; /* icon-circle 的一半 */
      height: 50%; /* icon-circle 的一半 */
      filter: brightness(0) invert(1) opacity(0.9); /* 将图标转为白色半透明 */
      object-fit: contain;
    }
  }

  .stat-content {
    .stat-value-row {
      gap: 6px;
      margin-top: 4px;
    }

    .stat-value {
      font-size: 24px; /* 移动端固定字体大小 */
    }

    .stat-label {
      font-size: 12px;
      margin-top: 4px;
      white-space: nowrap;
    }
  }

  .content-section {
    gap: 16px;
  }

  .table-card {
    border-radius: 8px;
    box-shadow: var(--app-shadow);
  }

  .card-header {
    padding: 16px 16px 0;

    .card-title {
      font-size: 18px;
    }
  }

  .table-container {
    padding: 20px 16px 16px;
  }
}

@media (max-width: 480px) {
  .stats-section {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
    padding: 12px;
  }

  .stat-content .stat-value {
    font-size: 20px;
  }

  .card-header {
    padding: 12px 12px 0;

    .card-title {
      font-size: 16px;
    }
  }

  .table-container {
    padding: 16px 12px 12px;
  }
}

/* 大屏幕vw适配 - 基于1920*1080设计的视口单位缩放 */
@media (min-width: 1025px) {
  .content-container {
    height: calc(100vh - 4.17vw); /* 精确计算可用高度：视口高度减去顶部导航栏高度 */
    padding: 1.25vw; /* 24px/1920 = 1.25vw */
    width: 100%; /* 充满整个屏幕 */
    margin: 0;
    box-sizing: border-box;
    display: flex; /* 桌面端使用flex布局传递高度给子组件 */
    flex-direction: column;
  }

  .dashboard {
    /* vw单位确保内容在2K(2560px)和4K(3840px)下按比例缩放充满屏幕 */
    overflow-x: hidden; /* 防止水平滚动条 */
  }
}

@media (max-width: 768px) {
  .content-container {
    padding: 16px;
  }
}

@media (max-width: 480px) {
  .content-container {
    padding: 12px;
  }
}

/* 超大屏幕适配 - 确保在8K等极端屏幕上也能正常显示 */
@media (min-width: 7680px) {
  .dashboard {
    .stat-content .stat-value {
      font-size: max(1.46vw, 28px); /* 设置最小字体大小 */
    }

    .stat-content .stat-label {
      font-size: max(0.73vw, 14px);
    }

    .card-title {
      font-size: max(1.04vw, 20px);
    }
  }
}
</style>

<style lang="scss">
/* 整体布局和背景 */

[data-theme='dark'] .dashboard {
  background: rgba(22, 29, 38, 1) !important;
}

/* 统计卡片区域暗模式优化 */
[data-theme='dark'] .stats-section {
  background: linear-gradient(135deg, #019C7C 0%, #10b981 100%) !important;
  box-shadow: 0 4px 20px rgba(16, 185, 129, 0.15) !important;

  // 暗色模式下的背景图片
  &::after {
    background: url('@/assets/images/dashboard-bg.png') lightgray 50% / cover no-repeat;
    mix-blend-mode: soft-light;
  }
}

/* 表格卡片暗模式 */
[data-theme='dark'] .table-card {
  background: rgba(31, 41, 53, 1) !important;
  border: 1px solid rgba(255, 255, 255, 0.12) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3) !important;
}

[data-theme='dark'] .card-title {
  color: #f9fafb !important;
}

/* 表格暗模式样式 */
[data-theme='dark'] .table-container :deep(.el-table) {
  background: rgba(46, 59, 75, 1) !important;
  border-color: rgba(255, 255, 255, 0.12) !important;
  color: #f9fafb !important;
}

[data-theme='dark'] .table-container :deep(.el-table .table-header th) {
  background-color: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12) !important;
}

[data-theme='dark'] .table-container :deep(.el-table__body tr) {
  background: rgba(46, 59, 75, 1) !important;
}

[data-theme='dark'] .table-container :deep(.el-table__body tr td) {
  background: rgba(46, 59, 75, 1) !important;
  color: #e5e7eb !important;
  border-bottom-color: rgba(255, 255, 255, 0.12) !important;
}

/* 条纹行样式 - 使用更高优先级的选择器 */
[data-theme='dark'] .dashboard .table-container :deep(.el-table .el-table__body .stripe-row) {
  background-color: rgba(46, 59, 75, 1) !important;
  border-top: 1px solid rgba(46, 59, 75, 1) !important;
  border-bottom: 1px solid rgba(46, 59, 75, 1) !important;
}

[data-theme='dark'] .dashboard .table-container :deep(.el-table .el-table__body .stripe-row td) {
  background-color: rgba(46, 59, 75, 1) !important;
  border-top: 1px solid rgba(46, 59, 75, 1) !important;
  border-bottom: 1px solid rgba(46, 59, 75, 1) !important;
}

/* 非条纹行保持统一背景 */
[data-theme='dark'] .dashboard .table-container :deep(.el-table .el-table__body tr:not(.stripe-row)) {
  background-color: rgba(46, 59, 75, 1) !important;
}

[data-theme='dark'] .dashboard .table-container :deep(.el-table .el-table__body tr:not(.stripe-row) td) {
  background-color: rgba(46, 59, 75, 1) !important;
}

[data-theme='dark'] .table-container :deep(.el-table__body tr:hover > td) {
  background-color: rgba(255, 255, 255, 0.08) !important;
}

/* 状态标签暗模式 */
[data-theme='dark'] .status-badge.status-valid {
  background: #1e40af !important;
  color: #93c5fd !important;
}

[data-theme='dark'] .status-badge.status-locked {
  background: #92400e !important;
  color: #fbbf24 !important;
}

[data-theme='dark'] .status-badge.status-invalid {
  background: #dc2626 !important;
  color: #fecaca !important;
}
</style>

<style>
/* 暗模式表格样式 - 非scoped确保优先级 */
[data-theme='dark'] .el-table {
  background: rgba(46, 59, 75, 1) !important;
}

[data-theme='dark'] .el-table .table-header th {
  background-color: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12) !important;
}

[data-theme='dark'] .el-table__header-wrapper .el-table__header th {
  background-color: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12) !important;
}

[data-theme='dark'] .el-table__header th {
  background-color: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12) !important;
}

[data-theme='dark'] .el-table th {
  background-color: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
}

[data-theme='dark'] .el-table th.is-leaf {
  background-color: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
}

/* 使用更通用的选择器强制覆盖所有表头样式 */
[data-theme='dark'] .dashboard .el-table thead th,
[data-theme='dark'] .dashboard .el-table .el-table__header th,
[data-theme='dark'] .dashboard .el-table .el-table__header-wrapper th,
[data-theme='dark'] .dashboard .el-table__header-wrapper .el-table__header thead th {
  background-color: rgba(46, 59, 75, 1) !important;
  background: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
  border-bottom: 1px solid rgba(255, 255, 255, 0.12) !important;
}

[data-theme='dark'] .el-table__body tr {
  background: rgba(46, 59, 75, 1) !important;
}

[data-theme='dark'] .el-table__body tr td {
  background: rgba(46, 59, 75, 1) !important;
  color: #e5e7eb !important;
  border-bottom-color: rgba(255, 255, 255, 0.12) !important;
}

[data-theme='dark'] .el-table__body .stripe-row {
  background-color: rgba(46, 59, 75, 1) !important;
  border-top: 1px solid rgba(46, 59, 75, 1) !important;
  border-bottom: 1px solid rgba(46, 59, 75, 1) !important;
}

[data-theme='dark'] .el-table__body .stripe-row td {
  background-color: rgba(46, 59, 75, 1) !important;
  border-top: 1px solid rgba(46, 59, 75, 1) !important;
  border-bottom: 1px solid rgba(46, 59, 75, 1) !important;
}

/* 最强力的表头样式覆盖 */
[data-theme='dark'] th {
  background-color: rgba(46, 59, 75, 1) !important;
  background: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
}

[data-theme='dark'] .el-table__cell {
  background-color: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
}

[data-theme='dark'] .dashboard th,
[data-theme='dark'] .dashboard .el-table th,
[data-theme='dark'] .dashboard .el-table .el-table__cell {
  background-color: rgba(46, 59, 75, 1) !important;
  background: rgba(46, 59, 75, 1) !important;
  color: #f9fafb !important;
}
</style>
