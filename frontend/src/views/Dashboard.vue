<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-01 09:32:42
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-12 09:44:47
 * @FilePath: /frontend/src/views/Dashboard.vue
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
<template>
  <Layout app-name="Cedar-V" page-title="仪表盘">
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
        <div class="chart-card">
          <div class="card-header">
            <h3 class="card-title">授权趋势</h3>
          </div>
          <div class="chart-container">
            <!-- 图表占位区域 -->
            <div class="chart-placeholder">
              <svg class="trend-chart" viewBox="0 0 800 200">
                <!-- 背景渐变 -->
                <defs>
                  <linearGradient id="chartGradient" x1="0%" y1="0%" x2="0%" y2="100%">
                    <stop offset="0%" stop-color="#5AD8A6" stop-opacity="0.6"/>
                    <stop offset="100%" stop-color="#5AD8A6" stop-opacity="0.08"/>
                  </linearGradient>
                </defs>
                <!-- 趋势线和填充区域 -->
                <path d="M40 160 Q200 100 400 80 Q600 60 760 40" 
                      stroke="#00C27C" 
                      stroke-width="2" 
                      fill="none"/>
                <path d="M40 160 Q200 100 400 80 Q600 60 760 40 L760 160 L40 160 Z" 
                      fill="url(#chartGradient)"/>
              </svg>
            </div>
          </div>
        </div>

        <!-- 最近授权表格 -->
        <div class="table-card">
          <div class="card-header">
            <h3 class="card-title">最近授权</h3>
          </div>
          <div class="table-container">
            <table class="recent-table">
              <thead>
                <tr>
                  <th>客户名称</th>
                  <th>产品</th>
                  <th>授权类型</th>
                  <th>到期时间</th>
                  <th>状态</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="item in recentData" :key="item.id" :class="item.rowClass">
                  <td>{{ item.customer }}</td>
                  <td>{{ item.product }}</td>
                  <td>{{ item.type }}</td>
                  <td>{{ item.expiry }}</td>
                  <td>
                    <span class="status-badge" :class="item.statusClass">
                      {{ item.status }}
                    </span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </Layout>
</template>

<script setup lang="ts">
import Layout from '@/components/common/layout/Layout.vue';

// 统计卡片数据
const statsData = [
  {
    id: 1,
    value: '1,234',
    label: '总授权数',
    trend: '8.5%',
  },
  {
    id: 2,
    value: '856',
    label: '活跃授权',
    trend: '8.5%',
  },
  {
    id: 3,
    value: '123',
    label: '即将到期',
    trend: '8.5%',
  },
  {
    id: 4,
    value: '45',
    label: '已过期',
    trend: '8.5%',
  },
  {
    id: 5,
    value: '678',
    label: '本月新增',
    trend: '8.5%',
  },
  {
    id: 6,
    value: '234',
    label: '客户总数',
    trend: '8.5%',
  }
];

// 最近授权数据
const recentData = [
  {
    id: 1,
    customer: '阿里巴巴集团',
    product: 'Cedar-V Pro',
    type: '企业版',
    expiry: '2024-12-31',
    status: '正常',
    statusClass: 'status-active',
    rowClass: ''
  },
  {
    id: 2,
    customer: '腾讯科技',
    product: 'Cedar-V Standard',
    type: '标准版',
    expiry: '2024-11-15',
    status: '即将到期',
    statusClass: 'status-warning',
    rowClass: ''
  },
  {
    id: 3,
    customer: '字节跳动',
    product: 'Cedar-V Enterprise',
    type: '企业版',
    expiry: '2025-03-20',
    status: '正常',
    statusClass: 'status-active',
    rowClass: 'row-highlight'
  },
  {
    id: 4,
    customer: '美团科技',
    product: 'Cedar-V Pro',
    type: '专业版',
    expiry: '2024-08-10',
    status: '已过期',
    statusClass: 'status-expired',
    rowClass: ''
  },
  {
    id: 5,
    customer: '滴滴出行',
    product: 'Cedar-V Standard',
    type: '标准版',
    expiry: '2025-01-30',
    status: '正常',
    statusClass: 'status-active',
    rowClass: 'row-highlight'
  },
  {
    id: 6,
    customer: '京东集团',
    product: 'Cedar-V Enterprise',
    type: '企业版',
    expiry: '2024-10-25',
    status: '即将到期',
    statusClass: 'status-warning',
    rowClass: 'row-highlight'
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
.chart-card,
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

// 图表卡片
.chart-container {
  padding: 26px 24px 24px;
  
  .chart-placeholder {
    height: 246px;
    display: flex;
    align-items: center;
    justify-content: center;
    
    .trend-chart {
      width: 100%;
      height: 100%;
    }
  }
}

// 表格卡片
.table-container {
  padding: 26px 24px 24px;
}

.recent-table {
  width: 100%;
  border-collapse: collapse;
  
  th {
    text-align: left;
    padding: 12px 16px;
    font-size: 14px;
    font-weight: 500;
    color: #666;
    border-bottom: 1px solid #F5F7FA;
    background: transparent;
  }
  
  td {
    padding: 16px;
    font-size: 14px;
    color: #1D1D1D;
    border-bottom: 1px solid #F7F8FA;
  }
  
  tr:hover {
    background-color: #F7F8FA;
  }
  
  .row-highlight {
    background-color: #F7F8FA;
  }
}

.status-badge {
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  
  &.status-active {
    background: #E8F5E8;
    color: #00C27C;
  }
  
  &.status-warning {
    background: #FFF3E0;
    color: #FF9800;
  }
  
  &.status-expired {
    background: #FFEBEE;
    color: #F44336;
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