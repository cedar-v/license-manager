<template>
  <div class="change-history-tab">
    <div class="history-section">
      <h3 class="section-title">Øô†ò</h3>
      <el-empty v-if="!historyData || historyData.length === 0" description="‚àØô°U" />
      <div v-else class="history-list">
        <div v-for="item in historyData" :key="item.id" class="history-item">
          <div class="history-time">{{ formatDateTime(item.created_at) }}</div>
          <div class="history-content">
            <div class="history-type">{{ item.change_type }}</div>
            <div class="history-reason">{{ item.reason }}</div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
interface HistoryItem {
  id: string
  change_type: string
  reason: string
  created_at: string
}

interface Props {
  historyData: HistoryItem[]
}

defineProps<Props>()

const formatDateTime = (date?: string) => {
  if (!date) return '-'
  return new Date(date).toLocaleString('zh-CN')
}
</script>

<style lang="scss" scoped>
.change-history-tab {
  width: 100%;
}

.history-section {
  .section-title {
    font-size: 16px;
    font-weight: 500;
    color: var(--app-text-primary);
    margin: 0 0 16px 0;
  }

  .history-list {
    .history-item {
      padding: 16px;
      border-left: 3px solid var(--el-color-primary);
      background: var(--app-bg-color);
      border-radius: 4px;
      margin-bottom: 12px;

      &:last-child {
        margin-bottom: 0;
      }

      .history-time {
        font-size: 12px;
        color: var(--app-text-secondary);
        margin-bottom: 8px;
      }

      .history-content {
        .history-type {
          font-size: 14px;
          font-weight: 500;
          color: var(--app-text-primary);
          margin-bottom: 4px;
        }

        .history-reason {
          font-size: 14px;
          color: var(--app-text-secondary);
        }
      }
    }
  }
}
</style>
