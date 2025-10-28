<template>
  <div class="change-history-container" v-loading="loading">
    <!-- 筛选条件 -->
    <div class="filter-bar">
      <div class="filter-item">
        <span class="filter-label">操作类型：</span>
        <el-select v-model="changeType" placeholder="请选择" class="filter-select">
          <el-option label="全部操作" value="all" />
          <el-option label="续期" value="renewal" />
          <el-option label="升级" value="upgrade" />
          <el-option label="限制变更" value="limit_change" />
          <el-option label="功能切换" value="feature_toggle" />
          <el-option label="锁定" value="lock" />
          <el-option label="解锁" value="unlock" />
          <el-option label="其他" value="other" />
        </el-select>
      </div>
      <div class="filter-item">
        <span class="filter-label">时间范围：</span>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          class="date-picker"
          :shortcuts="shortcuts"
          format="YYYY-MM-DD"
          value-format="YYYY-MM-DD"
        />
      </div>
    </div>

    <!-- 历史记录列表 -->
    <div class="history-container">
      <div v-if="!loading && historyList.length > 0" class="history-timeline">
        <div v-for="item in historyList" :key="item.id" class="history-item">
          <!-- 时间轴圆点 -->
          <div class="timeline-dot"></div>

          <!-- 历史记录内容 -->
          <div class="history-card">
            <div class="history-header">
              <div class="header-left">
                <div class="title-section">
                  <span class="operation-title">{{ item.change_type_display || item.change_type }}</span>
                </div>
              </div>
              <div class="header-right">
                <el-tag
                  :type="getStatusType(item.change_type)"
                  size="small"
                  class="status-tag"
                >
                  {{ item.change_type_display || item.change_type }}
                </el-tag>
              </div>
            </div>

            <div class="history-meta">
              <span class="operator">{{ item.operator_name || item.operator_id }}</span>
              <span class="time">{{ formatDateTime(item.created_at) }}</span>
            </div>

            <div v-if="item.reason" class="history-details">
              <div class="detail-item">
                变更原因：{{ item.reason }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!loading && historyList.length === 0" class="empty-state">
        <el-empty description="暂无变更历史记录" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import type { AuthorizationCode, AuthorizationChangeItem } from '@/api/license'
import { getAuthorizationChanges } from '@/api/license'
import { formatDate } from '@/utils/date'

interface Props {
  licenseData: AuthorizationCode | null
}

const props = defineProps<Props>()

// 筛选条件
const dateRange = ref<[string, string] | null>(null)
const changeType = ref('all')
const loading = ref(false)

// 历史记录数据
const historyList = ref<AuthorizationChangeItem[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(20)

// 日期快捷选项
const shortcuts = [
  {
    text: '最近7天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    },
  },
  {
    text: '最近30天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    },
  },
  {
    text: '最近90天',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end]
    },
  },
]

// 获取日期范围参数
const getDateRangeParams = () => {
  if (!dateRange.value || dateRange.value.length !== 2) {
    return { start_date: undefined, end_date: undefined }
  }

  return {
    start_date: dateRange.value[0],
    end_date: dateRange.value[1]
  }
}

// 获取变更历史列表
const fetchChangeHistory = async () => {
  console.log('获取变更历史，许可证ID:', props.licenseData)
  if (!props.licenseData?.id) {
    historyList.value = []
    return
  }

  loading.value = true
  try {
    const dateParams = getDateRangeParams()
    const response = await getAuthorizationChanges(props.licenseData.id, {
      page: currentPage.value,
      page_size: pageSize.value,
      change_type: changeType.value === 'all' ? undefined : changeType.value,
      ...dateParams
    })

    if (response.data?.list) {
      historyList.value = response.data.list
      total.value = response.data.total
    } else {
      historyList.value = []
      total.value = 0
    }
  } catch (error) {
    console.error('获取变更历史失败:', error)
    ElMessage.error('获取变更历史失败')
    historyList.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchChangeHistory()
})

// 监听 licenseData 变化
watch(() => props.licenseData?.id, () => {
  currentPage.value = 1
  fetchChangeHistory()
})

// 监听筛选条件变化
watch([dateRange, changeType], () => {
  currentPage.value = 1
  fetchChangeHistory()
})

// 格式化日期时间
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return '-'
  return formatDate(dateTime)
}

// 获取状态类型（根据变更类型判断）
const getStatusType = (changeType: string) => {
  // 锁定相关操作显示为危险
  if (changeType.includes('lock')) {
    return 'danger'
  }
  // 其他操作显示为成功
  return 'success'
}
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

    .filter-select,
    .date-picker {
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

    .date-picker {
      width: 280px;

      :deep(.el-range-separator) {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 12px;
        color: #666666;
      }

      :deep(.el-range-input) {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 12px;
        color: #1D1D1D;
      }

      :deep(.el-range__icon) {
        font-size: 14px;
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
    margin-bottom: 12px;

    .operator,
    .time {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 18px;
      color: #8186A5;
    }
  }

  .history-details {
    display: flex;
    flex-direction: column;
    gap: 8px;

    .detail-item {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      color: #666666;
      padding: 8px 12px;
      background: rgba(136, 165, 209, 0.1);
      border-radius: 4px;
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

      .filter-select,
      .date-picker {
        flex: 1;
        width: 100%;
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
