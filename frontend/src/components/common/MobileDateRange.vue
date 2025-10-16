<template>
  <div class="mobile-date-range">
    <div class="date-range-inputs">
      <div class="date-input-group">
        <label class="date-label">{{ t('chart.mobileDateRange.startDate') }}</label>
        <input
          type="date"
          v-model="startDateValue"
          :max="endDateValue || maxDate"
          :min="minDate"
          class="native-date-input"
          @change="handleStartChange"
        />
      </div>
      
      <div class="date-separator">{{ t('chart.mobileDateRange.separator') }}</div>
      
      <div class="date-input-group">
        <label class="date-label">{{ t('chart.mobileDateRange.endDate') }}</label>
        <input
          type="date"
          v-model="endDateValue"
          :min="startDateValue || minDate"
          :max="maxDate"
          class="native-date-input"
          @change="handleEndChange"
        />
      </div>
    </div>
    
    <!-- 快捷选择按钮 -->
    <div class="quick-buttons" v-if="showQuickButtons">
      <button
        v-for="option in quickOptions"
        :key="option.value"
        class="quick-btn"
        :class="{ active: selectedQuick === option.value }"
        @click="handleQuickSelect(option.value)"
      >
        {{ option.label }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useI18n } from 'vue-i18n'

interface Props {
  modelValue: [string, string] | null
  minDate?: string
  maxDate?: string
  showQuickButtons?: boolean
}

interface Emits {
  (e: 'update:modelValue', value: [string, string] | null): void
  (e: 'change', value: [string, string] | null): void
}

const props = withDefaults(defineProps<Props>(), {
  minDate: '2020-01-01',
  maxDate: '2030-12-31',
  showQuickButtons: true
})

const emit = defineEmits<Emits>()

// 使用国际化
const { t } = useI18n()

const startDateValue = ref('')
const endDateValue = ref('')
const selectedQuick = ref('')

// 快捷选择选项
const quickOptions = computed(() => [
  { label: t('chart.mobileDateRange.quickOptions.today'), value: 'today' },
  { label: t('chart.mobileDateRange.quickOptions.week'), value: 'week' },
  { label: t('chart.mobileDateRange.quickOptions.month'), value: 'month' }
])

// 格式化日期为 YYYY-MM-DD
const formatDateToISO = (date: Date): string => {
  return date.toISOString().split('T')[0]
}

// 处理开始日期变化 - 手动输入时清空快捷选择
const handleStartChange = () => {
  // 清空快捷选择状态，表示用户手动选择了日期
  selectedQuick.value = ''
  emitChange()
}

// 处理结束日期变化 - 手动输入时清空快捷选择
const handleEndChange = () => {
  // 清空快捷选择状态，表示用户手动选择了日期
  selectedQuick.value = ''
  emitChange()
}

// 发射变化事件
const emitChange = () => {
  if (startDateValue.value && endDateValue.value) {
    const result: [string, string] = [startDateValue.value, endDateValue.value]
    emit('update:modelValue', result)
    emit('change', result)
  } else {
    emit('update:modelValue', null)
    emit('change', null)
  }
}

// 快捷选择处理
const handleQuickSelect = (value: string) => {
  selectedQuick.value = value
  
  const today = new Date()
  let startDate: Date
  let endDate: Date
  
  switch (value) {
    case 'today':
      // 本日：今天到今天
      startDate = new Date(today)
      endDate = new Date(today)
      break
    case 'week':
      // 本周：本周一到本周日
      const currentDay = today.getDay() // 0=周日, 1=周一, ..., 6=周六
      const mondayOffset = currentDay === 0 ? 6 : currentDay - 1 // 计算到周一的偏移
      startDate = new Date(today)
      startDate.setDate(today.getDate() - mondayOffset)
      endDate = new Date(startDate)
      endDate.setDate(startDate.getDate() + 6) // 周日
      break
    case 'month':
      // 本月：当前月1号到当前月末
      startDate = new Date(today.getFullYear(), today.getMonth(), 1)
      endDate = new Date(today.getFullYear(), today.getMonth() + 1, 0) // 下个月0号=本月最后一天
      break
    default:
      return
  }
  
  startDateValue.value = formatDateToISO(startDate)
  endDateValue.value = formatDateToISO(endDate)
  
  emitChange()
}

// 监听外部值变化
watch(() => props.modelValue, (newVal) => {
  if (newVal && newVal[0] && newVal[1]) {
    startDateValue.value = newVal[0]
    endDateValue.value = newVal[1]
  } else {
    startDateValue.value = ''
    endDateValue.value = ''
  }
}, { immediate: true })
</script>

<style scoped>
.mobile-date-range {
  width: 100%;
  padding: 16px;
  background: #f7f8fa;
  border-radius: 8px;
  border: 1px solid #e4e7ed;
}

.date-range-inputs {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}

.date-input-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.date-label {
  font-size: 12px;
  color: #666;
  font-weight: 500;
}

.native-date-input {
  width: 100%;
  height: 40px;
  padding: 8px 12px;
  border: 1px solid #dcdfe6;
  border-radius: 6px;
  font-size: 14px;
  background: #fff;
  color: #606266;
  transition: border-color 0.2s ease;
  
  &:focus {
    outline: none;
    border-color: #019C7C;
    box-shadow: 0 0 0 2px rgba(1, 156, 124, 0.1);
  }
  
  &:hover {
    border-color: #c0c4cc;
  }
}

.date-separator {
  font-size: 14px;
  color: #666;
  margin-top: 20px;
  flex-shrink: 0;
}

.quick-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-btn {
  padding: 6px 12px;
  border: 1px solid #ddd;
  border-radius: 20px;
  background: #fff;
  color: #666;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  
  &:hover {
    border-color: #019C7C;
    color: #019C7C;
  }
  
  &.active {
    background: #019C7C;
    border-color: #019C7C;
    color: #fff;
  }
  
  &:active {
    transform: scale(0.95);
  }
}

/* 移动端优化 */
@media (max-width: 480px) {
  .mobile-date-range {
    padding: 12px;
  }
  
  .date-range-inputs {
    flex-direction: column;
    gap: 8px;
  }
  
  .date-separator {
    align-self: center;
    margin-top: 0;
  }
  
  .quick-buttons {
    justify-content: center;
  }
  
  .quick-btn {
    font-size: 11px;
    padding: 5px 10px;
  }
}
</style>