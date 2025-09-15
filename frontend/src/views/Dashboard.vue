<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-01 09:32:42
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-09-12 10:36:03
 * @FilePath: /frontend/src/views/Dashboard.vue
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
<template>
  <Layout app-name="Cedar-V" :page-title="t('dashboard.title')">
    <!-- 页面内容 -->
    <div class="content-container dashboard">
      <!-- 统计卡片区域 -->
      <div class="stats-section">
        <div class="stats-card" v-for="stat in statsData" :key="stat.id">
          <div class="stat-icon">
            <div class="icon-circle">
              <svg class="stat-icon-svg" width="24" height="24" viewBox="0 0 24 24" fill="none">
                <path v-if="stat.id === 1" d="M3 13h8V3H3v10zm0 8h8v-6H3v6zm10 0h8V11h-8v10zm0-18v6h8V3h-8z" fill="currentColor"/>
                <path v-else-if="stat.id === 2" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" fill="currentColor"/>
                <path v-else-if="stat.id === 3" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" fill="currentColor"/>
                <path v-else-if="stat.id === 4" d="M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z" fill="currentColor"/>
                <path v-else-if="stat.id === 5" d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-5 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z" fill="currentColor"/>
                <path v-else="stat.id === 6" d="M16 4c0-1.11.89-2 2-2s2 .89 2 2-.89 2-2 2-2-.89-2-2zM4 18v-1c0-1.1.9-2 2-2h2.5c.91 0 1.73.35 2.35.95C11.85 15.35 12.09 15 12.5 15H15c1.1 0 2 .9 2 2v1h3v-6H4v6z" fill="currentColor"/>
              </svg>
            </div>
          </div>
          <div class="stat-content">
            <div class="stat-label">{{ stat.label }}</div>
            <div class="stat-value">{{ stat.value }}</div>
          </div>
          <div class="stat-trend">
            <svg class="trend-icon" width="16" height="16" viewBox="0 0 16 16">
              <path d="M4.5 7.5L7 5L9.5 7.5L14 3" stroke="currentColor" stroke-width="1.5" fill="none"/>
            </svg>
            <span class="trend-value">{{ stat.trend }}</span>
          </div>
        </div>
      </div>

      <!-- 图表和表格区域 -->
      <div class="content-section">
        <!-- 授权趋势图表 - 占据40%高度 -->
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
              style="width: 100%; height: 100%;"
              :header-row-class-name="'table-header'"
              :row-class-name="(params: any) => params.rowIndex % 2 === 1 ? 'stripe-row' : ''"
            >
              <el-table-column :label="t('dashboard.recentLicenses.columns.serialNumber')"  min-width="90">
                <template #default="{ $index }" style="min-width: 90px;">
                  {{ $index + 1 }}
                </template>
              </el-table-column>
              <el-table-column prop="customer" :label="t('dashboard.recentLicenses.columns.customerName')" min-width="200" />
              <el-table-column prop="description" :label="t('dashboard.recentLicenses.columns.description')" min-width="150" />
              <el-table-column :label="t('dashboard.recentLicenses.columns.status')" width="120" align="center">
                <template #default="{ row }">
                  <span 
                    class="status-badge" 
                    :class="row.status === 1 ? 'status-valid' : 'status-invalid'"
                  >
                    {{ row.status === 1 ? t('dashboard.recentLicenses.statusLabels.valid') : t('dashboard.recentLicenses.statusLabels.invalid') }}
                  </span>
                </template>
              </el-table-column>
              <el-table-column prop="expiry" :label="t('dashboard.recentLicenses.columns.expiryTime')" width="301" />
              <el-table-column prop="createTime" :label="t('dashboard.recentLicenses.columns.createTime')" width="301" />
            </el-table>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Layout from '@/components/common/layout/Layout.vue';
import LicenseTrendChart from '@/components/charts/LicenseTrendChart.vue';

// 使用国际化
const { t } = useI18n()

// 统计卡片数据
const statsData = computed(() => [
  {
    id: 1,
    value: '1,234',
    label: t('dashboard.stats.totalLicenses'),
    trend: '8.5%',
  },
  {
    id: 2,
    value: '856',
    label: t('dashboard.stats.activeLicenses'),
    trend: '8.5%',
  },
  {
    id: 3,
    value: '123',
    label: t('dashboard.stats.expiringSoon'),
    trend: '8.5%',
  },
  {
    id: 4,
    value: '45',
    label: t('dashboard.stats.expired'),
    trend: '8.5%',
  },
  {
    id: 5,
    value: '678',
    label: t('dashboard.stats.newThisMonth'),
    trend: '8.5%',
  },
  {
    id: 6,
    value: '234',
    label: t('dashboard.stats.totalCustomers'),
    trend: '8.5%',
  }
]);

// 最近授权数据 
const recentData = [
  {
    id: 1,
    customer: '石狮市潢安有限公司',
    description: '落魄山一哥',
    status: 1, // 1代表有效，0代表失效
    expiry: '2023-05-26 12:12:00',
    createTime: '2023-05-26 12:12:00'
  },
  {
    id: 2,
    customer: '石狮市潢安有限公司',
    description: '落魄山一哥',
    status: 0,
    expiry: '2023-05-26 12:12:00',
    createTime: '2023-05-26 12:12:00'
  },
  {
    id: 3,
    customer: '石狮市潢安有限公司',
    description: '落魄山一哥',
    status: 1,
    expiry: '2023-05-26 12:12:00',
    createTime: '2023-05-26 12:12:00'
  },
  {
    id: 4,
    customer: '石狮市潢安有限公司',
    description: '落魄山一哥',
    status: 0,
    expiry: '2023-05-26 12:12:00',
    createTime: '2023-05-26 12:12:00'
  },
  {
    id: 5,
    customer: '石狮市潢安有限公司',
    description: '落魄山一哥',
    status: 1,
    expiry: '2023-05-26 12:12:00',
    createTime: '2023-05-26 12:12:00'
  },
  {
    id: 6,
    customer: '石狮市潢安有限公司',
    description: '落魄山一哥',
    status: 0,
    expiry: '2023-05-26 12:12:00',
    createTime: '2023-05-26 12:12:00'
  }
];
</script>

<style lang="scss" scoped>
/* content-container样式 */
.content-container {
  min-height: calc(100vh - 80px);
  padding: 24px;
  width: 100%;
  margin: 0;
  box-sizing: border-box;
  background: #F7F8FA;
}

.dashboard {
  width: 100%;
  height: 100%; /* 使用100%高度充满父容器 */
  display: flex;
  flex-direction: column;
}

// 统计卡片区域 - 基于1920*1080设计的vw适配
.stats-section {
  display: grid;
  grid-template-columns: repeat(6, 1fr); /* 6个统计卡片 */
  gap: 1.04vw; /* 20px/1920 = 1.04vw */
  margin-bottom: 2.08vw; /* 40px/1920 = 2.08vw */
  padding: 1.25vw; /* 24px/1920 = 1.25vw */
  background: linear-gradient(135deg, #019C7C 0%, #5AD8A6 100%);
  border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
  position: relative;
  
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(45deg, rgba(255,255,255,0.08) 25%, transparent 25%), 
                linear-gradient(-45deg, rgba(255,255,255,0.08) 25%, transparent 25%), 
                linear-gradient(45deg, transparent 75%, rgba(255,255,255,0.08) 75%), 
                linear-gradient(-45deg, transparent 75%, rgba(255,255,255,0.08) 75%);
    background-size: 1.04vw 1.04vw; /* 20px/1920 = 1.04vw */
    background-position: 0 0, 0 0.52vw, 0.52vw -0.52vw, -0.52vw 0vw; /* 10px/1920 = 0.52vw */
    opacity: 0.15;
    border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
  }
}

.stats-card {
  display: flex;
  align-items: center;
  gap: 0.83vw; /* 16px/1920 = 0.83vw */
  padding: 1.25vw 0; /* 24px/1920 = 1.25vw */
  position: relative;
  z-index: 1;
  
  &:not(:last-child)::after {
    content: '';
    position: absolute;
    right: -0.52vw; /* -10px/1920 = -0.52vw */
    top: 50%;
    transform: translateY(-50%);
    width: 1px;
    height: 3.33vw; /* 64px/1920 = 3.33vw */
    background: linear-gradient(to bottom, 
      rgba(255, 255, 255, 0.08) 0%, 
      rgba(255, 255, 255, 0.6) 49%, 
      rgba(255, 255, 255, 0.08) 100%);
  }
}

.stat-icon {
  .icon-circle {
    width: 2.92vw; /* 56px/1920 = 2.92vw */
    height: 2.92vw; /* 56px/1920 = 2.92vw */
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    
    .stat-icon-svg {
      color: rgba(255, 255, 255, 0.9);
      width: 1.25vw; /* 24px/1920 = 1.25vw */
      height: 1.25vw; /* 24px/1920 = 1.25vw */
    }
  }
}

.stat-content {
  flex: 1;
  min-width: 64px;
  
  .stat-label {
    font-size: 0.83vw; /* 16px/1920 = 1.46vw */
    font-weight: 500;
    color: #FFFFFF;
    line-height: 1;
  }
  
  .stat-value {
    font-size: 1.25vw; /* 24px/1920 = 0.73vw */
    color: rgba(255, 255, 255, 0.9);
    margin-top: 0.21vw; /* 4px/1920 = 0.21vw */
  }
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 0.21vw; /* 4px/1920 = 0.21vw */
  color: rgba(255, 255, 255, 0.8);
  
  .trend-icon {
    color: rgba(255, 255, 255, 0.8);
  }
  
  .trend-value {
    font-size: 0.73vw; /* 14px/1920 = 0.73vw */
    font-weight: 500;
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

// 图表区域 - 占据40%高度
.chart-section {
  flex: 0.4; /* 40%高度 */
  min-height: 0;
  display: flex;
  flex-direction: column;
}

// 卡片通用样式 - 表格卡片占据60%高度
.table-card {
  background: #FFFFFF;
  border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
  box-shadow: 0px 0px 0.21vw 0px rgba(0, 0, 0, 0.1); /* 4px/1920 = 0.21vw */
  overflow: hidden;
  display: flex;
  flex-direction: column;
  flex: 0.6; /* 占据60%高度 */
  min-height: 0; /* 防止flex子元素溢出 */
}

.card-header {
  padding: 1.25vw 1.25vw 0; /* 24px/1920 = 1.25vw */
  
  .card-title {
    font-size: 1.04vw; /* 20px/1920 = 1.04vw */
    font-weight: 400;
    color: #1D1D1D;
    margin: 0;
    line-height: 1.3;
  }
}

// 表格卡片
.table-container {
  padding: 1.35vw 1.25vw 1.25vw; /* 26px 24px 24px → 1.35vw 1.25vw 1.25vw */
  flex: 1; /* 占据剩余高度 */
  overflow: auto; /* 如果内容过多，允许滚动 */
  display: flex;
  flex-direction: column;
  
  // Element Plus 表格样式重写
  :deep(.el-table) {
    border: 1px solid #F5F7FA;
    height: 100%; /* 充满容器高度 */
    
    .el-table__body-wrapper {
      flex: 1; /* 表格主体占据剩余高度 */
    }
    
    .table-header th {
      background-color: #F7F8FA !important;
      color: #1D1D1D;
    font-weight: 500;
      font-size: 16px;
      font-family: 'Source Han Sans CN', sans-serif;
      height: 48px;
      padding: 13px 20px;
      border-bottom: none;
    }
    
    .el-table__body tr {
      height: 48px;
      
      td {
        padding: 13px 20px;
        border-bottom: 1px solid #F7F8FA;
    font-size: 14px;
        font-family: 'Source Han Sans CN', sans-serif;
        font-weight: 350;
    color: #1D1D1D;
  }
  
      &.stripe-row {
    background-color: #F7F8FA;
  }
  
      &:hover > td {
        background-color: #F7F8FA !important;
      }
    }
    
    // 序号列和状态列居中
    .el-table__body tr td:nth-child(1),
    .el-table__body tr td:nth-child(4) {
      text-align: center;
    }
    
    // 状态列防止文本溢出
    .el-table__body tr td:nth-child(4) {
      overflow: visible;
      text-overflow: initial;
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
    background:#F0F5FF;
    color:#4763FF;
  }
  
  &.status-invalid {
    background:#FFE5E5;
    color: #E90C0C;
  }
}

// 响应式设计 - 移动端切换回px单位确保可读性
@media (max-width: 768px) {
  .stats-section {
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
  
  .stat-icon .icon-circle {
    width: 48px; /* 移动端固定大小 */
    height: 48px;
    
    .stat-icon-svg {
      width: 20px;
      height: 20px;
    }
  }
  
  .stat-content {
    .stat-value {
      font-size: 24px; /* 移动端固定字体大小 */
    }
    
    .stat-label {
      font-size: 12px;
      margin-top: 4px;
    }
  }
  
  .stat-trend .trend-value {
    font-size: 12px;
  }
  
  .content-section {
    gap: 16px;
  }
  
  .table-card {
    border-radius: 8px;
    box-shadow: 0px 0px 4px 0px rgba(0, 0, 0, 0.1);
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