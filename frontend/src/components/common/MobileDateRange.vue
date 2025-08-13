<template>
  <div class="mobile-date-range">
    <div class="date-range-inputs">
      <div class="date-input-group">
        <label class="date-label">开始日期</label>
        <input
          type="date"
          v-model="startDateValue"
          :max="endDateValue || maxDate"
          :min="minDate"
          class="native-date-input"
          @change="handleStartChange"
        />
      </div>
      
      <div class="date-separator">至</div>
      
      <div class="date-input-group">
        <label class="date-label">结束日期</label>
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

const startDateValue = ref('')
const endDateValue = ref('')
const selectedQuick = ref('')

// 快捷选择选项
const quickOptions = [
  { label: '今天', value: 'today' },
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' },
  { label: '近7天', value: '7days' },
  { label: '近30天', value: '30days' }
]

// 格式化日期为 YYYY-MM-DD
const formatDate = (date: Date): string => {
  return date.toISOString().split('T')[0]
}

// 处理开始日期变化
const handleStartChange = () => {
  selectedQuick.value = ''
  emitChange()
}

// 处理结束日期变化
const handleEndChange = () => {
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
  let endDate: Date = new Date(today)
  
  switch (value) {
    case 'today':
      startDate = new Date(today)
      break
    case 'week':
      startDate = new Date(today.getTime() - 6 * 24 * 60 * 60 * 1000)
      break
    case 'month':
      startDate = new Date(today.getFullYear(), today.getMonth(), 1)
      break
    case '7days':
      startDate = new Date(today.getTime() - 6 * 24 * 60 * 60 * 1000)
      break
    case '30days':
      startDate = new Date(today.getTime() - 29 * 24 * 60 * 60 * 1000)
      break
    default:
      return
  }
  
  startDateValue.value = formatDate(startDate)
  endDateValue.value = formatDate(endDate)
  
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