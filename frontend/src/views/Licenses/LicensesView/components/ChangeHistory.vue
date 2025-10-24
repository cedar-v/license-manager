<template>
  <div class="change-history-container">
    <!-- 筛选条件 -->
    <div class="filter-bar">
      <div class="filter-item">
        <span class="filter-label">时间范围：</span>
        <el-select v-model="timeRange" placeholder="请选择" class="filter-select">
          <el-option label="最近7天" value="7" />
          <el-option label="最近30天" value="30" />
          <el-option label="最近90天" value="90" />
          <el-option label="全部" value="all" />
        </el-select>
      </div>
      <div class="filter-item">
        <span class="filter-label">操作类型：</span>
        <el-select v-model="operationType" placeholder="请选择" class="filter-select">
          <el-option label="全部操作" value="all" />
          <el-option label="创建授权" value="create" />
          <el-option label="设备激活" value="activate" />
          <el-option label="授权锁定" value="lock" />
          <el-option label="解绑设备" value="unbind" />
        </el-select>
      </div>
    </div>

    <!-- 历史记录列表 -->
    <div class="history-container">
      <div v-if="historyList.length > 0" class="history-timeline">
        <div v-for="item in historyList" :key="item.id" class="history-item">
          <!-- 时间轴圆点 -->
          <div class="timeline-dot"></div>

          <!-- 历史记录内容 -->
          <div class="history-card">
            <div class="history-header">
              <div class="header-left">
                <div class="title-section">
                  <span class="operation-title">{{ item.title }}</span>
                </div>
              </div>
              <div class="header-right">
                <el-tag
                  :type="getStatusType(item.status)"
                  size="small"
                  class="status-tag"
                >
                  {{ item.statusText }}
                </el-tag>
              </div>
            </div>

            <div class="history-meta">
              <span class="operator">{{ item.operator }}</span>
              <span class="time">{{ formatDateTime(item.time) }}</span>
            </div>

            <div class="detail-type">
              {{ item.detailType === 'compare' ? '变更对比' : '变更详情' }}
            </div>

            <div class="history-details">
              <div v-for="(detail, idx) in item.details" :key="idx" class="detail-item">
                {{ detail }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="empty-state">
        <el-empty description="暂无变更历史记录" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { formatDateTime } from '@/utils/date'

interface HistoryDetail {
  id: string
  title: string
  status: 'success' | 'locked' | 'failed'
  statusText: string
  operator: string
  time: string
  detailType: 'compare' | 'detail'
  details: string[]
}

interface Props {
  licenseData?: any
}

defineProps<Props>()

// 筛选条件
const timeRange = ref('30')
const operationType = ref('all')

// 模拟历史记录数据
const historyList = ref<HistoryDetail[]>([
  {
    id: '1',
    title: '创建新授权码',
    status: 'success',
    statusText: '成功',
    operator: '管理员',
    time: '2025-04-29 18:41:20',
    detailType: 'detail',
    details: [
      '授权码：XKDA-55GD-KJFA5-UIYP',
      '客户ID：X10023'
    ]
  },
  {
    id: '2',
    title: '设备激活',
    status: 'success',
    statusText: '成功',
    operator: '张经理',
    time: '2025-04-29 18:41:20',
    detailType: 'compare',
    details: [
      '设备名称：办公室工作站',
      '状态变更：未激活➡️已激活'
    ]
  },
  {
    id: '3',
    title: '授权锁定',
    status: 'locked',
    statusText: '已锁定',
    operator: '系统自动',
    time: '2025-04-29 18:41:20',
    detailType: 'compare',
    details: [
      '授权状态：正常➡️已锁定',
      '锁定原因：检测到异常使用行为'
    ]
  },
  {
    id: '4',
    title: '解绑设备',
    status: 'success',
    statusText: '成功',
    operator: '管理员',
    time: '2025-04-29 18:41:20',
    detailType: 'detail',
    details: [
      '设备ID：GD-KJFA5-UIYP',
      '解绑原因：设备更换'
    ]
  }
])


// 获取状态类型
const getStatusType = (status: string) => {
  const typeMap: Record<string, 'success' | 'danger' | 'warning'> = {
    success: 'success',
    locked: 'danger',
    failed: 'danger'
  }
  return typeMap[status] || 'success'
}

// 查看详情
// const handleViewDetail = (item: HistoryDetail) => {
//   console.log('查看详情', item)
//   // TODO: 实现查看详情逻辑
// }
</script>

<style lang="scss" scoped>
.change-history-container {
  width: 100%;
}

.filter-bar {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 24px;

  .filter-item {
    display: flex;
    align-items: center;
    gap: 4px;

    .filter-label {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 500;
      line-height: 22px;
      color: #1D1D1D;
      white-space: nowrap;
    }

    .filter-select {
      width: 216px;

      :deep(.el-input__wrapper) {
        padding: 2px 4px 2px 12px;
        border: 1px solid #DCDEE2;
        border-radius: 4px;
      }

      :deep(.el-input__inner) {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 12px;
        font-weight: 500;
        line-height: 24px;
        color: #B2B8C2;
      }
    }
  }
}

.history-container {
  background: #FFFFFF;
  border: 1px solid #E2E2E2;
  border-radius: 4px;
  padding: 20px;
  min-height: 400px;
}

.history-timeline {
  position: relative;
  padding-left: 25px;

  &::before {
    content: '';
    position: absolute;
    left: 5px;
    top: 0;
    bottom: 0;
    width: 1px;
    background: #DCDFE6;
  }
}

.history-item {
  position: relative;
  margin-bottom: 24px;

  &:last-child {
    margin-bottom: 0;

    .history-card::after {
      display: none;
    }
  }

  .timeline-dot {
    position: absolute;
    left: -19px;
    top: 5px;
    width: 11px;
    height: 11px;
    background: #4876FF;
    border: 2px solid #FFFFFF;
    border-radius: 50%;
    z-index: 1;
  }
}

.history-card {
  position: relative;
  padding-bottom: 24px;

  &::after {
    content: '';
    position: absolute;
    left: 28px;
    bottom: 0;
    right: 0;
    height: 1px;
    background: #DCDFE6;
  }

  .history-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 4px;

    .header-left {
      flex: 1;

      .title-section {
        display: flex;
        flex-direction: column;
        gap: 4px;

        .operation-title {
          font-family: 'Source Han Sans CN', sans-serif;
          font-size: 14px;
          font-weight: 500;
          line-height: 22px;
          color: #202332;
        }

      }
    }

    .header-right {
      display: flex;
      align-items: center;

      .status-tag {
        padding: 7px 16px;
        height: 24px;
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 500;
        line-height: 21px;
        border-radius: 4px;

        &.el-tag--success {
          background: rgba(0, 194, 124, 0.08);
          border-color: transparent;
          color: #019C7C;
        }

        &.el-tag--danger {
          background: rgba(240, 20, 47, 0.08);
          border-color: transparent;
          color: #F0142F;
        }
      }
    }
  }

  .history-meta {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 8px;

    .operator,
    .time {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 18px;
      color: #8186A5;
    }
  }

  .detail-type {
    font-family: 'Source Han Sans CN', sans-serif;
    font-size: 14px;
    font-weight: 500;
    line-height: 22px;
    color: #202332;
    margin-bottom: 8px;
  }

  .history-details {
    display: flex;
    flex-wrap: wrap;
    gap: 16px;

    .detail-item {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 18px;
      color: #8186A5;
    }
  }
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 400px;
}

// 响应式布局
@media (max-width: 768px) {
  .filter-bar {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;

    .filter-item {
      width: 100%;

      .filter-select {
        flex: 1;
      }
    }
  }

  .history-timeline {
    padding-left: 20px;
  }

  .history-card {
    .history-header {
      flex-direction: column;
      align-items: flex-start;
      gap: 8px;

      .header-right {
        align-self: flex-end;
      }
    }

    .history-details {
      flex-direction: column;
      gap: 4px;
    }
  }
}
</style>
