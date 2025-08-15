<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-01 09:32:42
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-14 09:12:33
 * @FilePath: /frontend/src/views/Dashboard.vue
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
<template>
  <Layout app-name="Cedar-V" :page-title="t('dashboard.title')">
    <!-- 页面内容 -->
    <div class="dashboard">
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
            <div class="stat-value">{{ stat.value }}</div>
            <div class="stat-label">{{ stat.label }}</div>
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
        <!-- 授权趋势图表 -->
        <LicenseTrendChart />

        <!-- 最近授权表格 -->
        <div class="table-card">
          <div class="card-header">
            <h3 class="card-title">{{ t('dashboard.recentLicenses.title') }}</h3>
          </div>
          <div class="table-container">
            <el-table
              :data="recentData"
              style="width: 100%"
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
.dashboard {
  padding: 0;
}

// 统计卡片区域
.stats-section {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 20px;
  margin-bottom: 40px;
  padding: 24px;
  background: linear-gradient(135deg, #019C7C 0%, #5AD8A6 100%);
  border-radius: 8px;
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
    background-size: 20px 20px;
    background-position: 0 0, 0 10px, 10px -10px, -10px 0px;
    opacity: 0.15;
    border-radius: 8px;
  }
}

.stats-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px 0;
  position: relative;
  z-index: 1;
  
  &:not(:last-child)::after {
    content: '';
    position: absolute;
    right: -10px;
    top: 50%;
    transform: translateY(-50%);
    width: 1px;
    height: 64px;
    background: linear-gradient(to bottom, 
      rgba(255, 255, 255, 0.08) 0%, 
      rgba(255, 255, 255, 0.6) 49%, 
      rgba(255, 255, 255, 0.08) 100%);
  }
}

.stat-icon {
  .icon-circle {
    width: 56px;
    height: 56px;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    
    .stat-icon-svg {
      color: rgba(255, 255, 255, 0.9);
      width: 24px;
      height: 24px;
    }
  }
}

.stat-content {
  flex: 1;
  
  .stat-value {
    font-size: 28px;
    font-weight: 600;
    color: #FFFFFF;
    line-height: 1.2;
  }
  
  .stat-label {
    font-size: 14px;
    color: rgba(255, 255, 255, 0.9);
    margin-top: 4px;
  }
}

.stat-trend {
  display: flex;
  align-items: center;
  gap: 4px;
  color: rgba(255, 255, 255, 0.8);
  
  .trend-icon {
    color: rgba(255, 255, 255, 0.8);
  }
  
  .trend-value {
    font-size: 14px;
    font-weight: 500;
  }
}

// 内容区域 - 默认上下布局，充满屏幕宽度
.content-section {
  display: grid;
  grid-template-columns: 1fr;
  gap: 24px;
}

// 卡片通用样式
.table-card {
  background: #FFFFFF;
  border-radius: 8px;
  box-shadow: 0px 0px 4px 0px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.card-header {
  padding: 24px 24px 0;
  
  .card-title {
    font-size: 20px;
    font-weight: 400;
    color: #1D1D1D;
    margin: 0;
    line-height: 1.3;
  }
}

// 表格卡片
.table-container {
  padding: 26px 24px 24px;
  
  // Element Plus 表格样式重写
  :deep(.el-table) {
    border: 1px solid #F5F7FA;
    
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

// 响应式设计已在默认样式中设置为上下布局

@media (max-width: 768px) {
  .stats-section {
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
    padding: 16px;
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
  
  .stat-content {
    .stat-value {
      font-size: 24px;
    }
  }
  
  // 图表和表格区域保持上下布局（已在默认样式中设置）
}

@media (max-width: 480px) {
  .stats-section {
    grid-template-columns: repeat(2, 1fr);
  }
  
  .recent-table {
    font-size: 12px;
    
    th,
    td {
      padding: 8px;
    }
  }
}
</style>