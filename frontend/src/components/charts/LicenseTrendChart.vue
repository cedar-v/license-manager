<template>
  <div class="license-trend-chart">
    <!-- 卡片头部 -->
    <div class="chart-header">
      <h3 class="chart-title">{{ t('chart.licenseTrend.title') }}</h3>
      
      <!-- 时间选择器 -->
      <div class="time-selector">
        <!-- 桌面端：横向布局 -->
        <template v-if="!isMobile">
          <!-- 快捷选择按钮 -->
          <div class="quick-selector">
            <el-button 
              v-for="option in quickOptions" 
              :key="option.value"
              :type="selectedQuick === option.value ? 'primary' : 'default'"
              size="small"
              @click="handleQuickSelect(option.value)"
            >
              {{ option.label }}
            </el-button>
          </div>
          
          <!-- 日期范围选择器 - 条件渲染模式 -->
          <!-- 显示模式：当有快捷选择时，显示只读的文本输入框 -->
          <div v-if="selectedQuick" class="date-display-wrapper">
            <div class="date-display" @click="switchToEditMode">
              {{ formatDateRange(dateRange) }}
              <el-icon class="date-icon"><Calendar /></el-icon>
            </div>
          </div>
          
          <!-- 编辑模式：实际的日期选择器 -->
          <el-date-picker
            v-else
            ref="datePickerRef"
            v-model="dateRange"
            type="daterange"
            :range-separator="t('chart.licenseTrend.datePicker.rangeSeparator')"
            :start-placeholder="t('chart.licenseTrend.datePicker.startPlaceholder')"
            :end-placeholder="t('chart.licenseTrend.datePicker.endPlaceholder')"
            size="small"
            format="MM-DD"
            value-format="YYYY-MM-DD"
            :clearable="false"
            :editable="false"
            @change="handleDateChange"
            style="width: 200px;"
          />
        </template>
        
        <!-- 移动端：上下布局 -->
        <template v-else>
          <!-- 快捷选择按钮 -->
          <div class="quick-selector mobile-quick">
            <el-button 
              v-for="option in quickOptions" 
              :key="option.value"
              :type="selectedQuick === option.value ? 'primary' : 'default'"
              size="small"
              @click="handleQuickSelect(option.value)"
            >
              {{ option.label }}
            </el-button>
          </div>
          
          <!-- 日期范围选择器 -->
          <MobileDateRange
            v-model="dateRange"
            :show-quick-buttons="false"
            @change="handleDateChange"
          />
        </template>
      </div>
    </div>

    <!-- 图表容器 -->
    <div class="chart-container">
      <v-chart
        class="trend-chart"
        :option="chartOption"
        :autoresize="true"
        :loading="chartLoading"
        ref="chartRef"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { use } from 'echarts/core'
import { useDevice } from '@/utils/useDevice'
import MobileDateRange from '@/components/common/MobileDateRange.vue'
import { Calendar } from '@element-plus/icons-vue'
import { getAuthorizationTrend, type TrendDataItem } from '@/api/dashboard'
import { ElMessage } from 'element-plus'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  GridComponent,
  TooltipComponent,
  LegendComponent
} from 'echarts/components'
import VChart from 'vue-echarts'

// 注册必要的组件
use([
  CanvasRenderer,
  LineChart,
  GridComponent,
  TooltipComponent,
  LegendComponent
])

// 使用国际化
const { t } = useI18n()

// 快捷选择选项
const quickOptions = computed(() => [
  { label: t('chart.licenseTrend.quickOptions.today'), value: 'today' },
  { label: t('chart.licenseTrend.quickOptions.week'), value: 'week' },
  { label: t('chart.licenseTrend.quickOptions.month'), value: 'month' }
])

// 响应式数据
const { isMobile } = useDevice()
const selectedQuick = ref('week')
const dateRange = ref<[string, string]>(['2024-05-13', '2024-05-17']) // 初始值会被handleQuickSelect('week')覆盖
const chartRef = ref()
const datePickerRef = ref()

// 授权趋势数据
const trendData = ref<{ date: string; value: number }[]>([])
const chartLoading = ref(false)

// 图表配置
const chartOption = computed(() => {
  const dates = trendData.value.map(item => item.date.substring(5)) // 只显示月-日
  const values = trendData.value.map(item => item.value)
  
  return {
    grid: {
      left: 60,
      right: 30,
      top: 30,
      bottom: 60,
      containLabel: false
    },
    xAxis: {
      type: 'category',
      data: dates,
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: 'var(--app-text-primary)',
        fontSize: 12,
        fontFamily: 'Inter',
        margin: 12
      }
    },
    yAxis: {
      type: 'value',
      min: 0,
      max: 20,
      interval: 5,
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: 'var(--app-text-primary)',
        fontSize: 12,
        fontFamily: 'Inter',
        align: 'right'
      },
      splitLine: {
        lineStyle: {
          color: 'var(--app-border-color)',
          width: 1
        }
      }
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'var(--app-content-bg)',
      borderColor: 'var(--app-content-bg)',
      borderWidth: 0.5,
      borderRadius: 4,
      textStyle: {
        color: 'var(--app-text-primary)',
        fontSize: 12
      },
      extraCssText: 'box-shadow: 0px 4px 12px 0px rgba(59, 210, 180, 0.2); backdrop-filter: blur(4px);',
      formatter: (params: any) => {
        const data = params[0]
        return `${t('chart.licenseTrend.tooltip.licenseCount')}: ${data.value}`
      }
    },
    series: [
      {
        type: 'line',
        data: values,
        smooth: true,
        lineStyle: {
          color: '#00C27C',
          width: 2
        },
        itemStyle: {
          color: '#00C27C',
          borderColor: '#00C27C',
          borderWidth: 3
        },
        symbol: 'circle',
        symbolSize: 8,
        emphasis: {
          itemStyle: {
            shadowBlur: 12,
            shadowColor: 'rgba(59, 210, 180, 1)'
          }
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: [
              {
                offset: 0,
                color: 'rgba(90, 216, 166, 0.6)'
              },
              {
                offset: 1,
                color: 'rgba(90, 216, 166, 0.08)'
              }
            ]
          }
        }
      }
    ]
  }
})

// 快捷选择处理 - 简化版，不处理弹出层冲突
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
  
  // 设置日期值
  dateRange.value = [
    startDate.toISOString().split('T')[0],
    endDate.toISOString().split('T')[0]
  ]
  
  // 更新图表数据
  updateChartData()
}

// 日期范围改变处理 - 只有通过日期选择器手动选择时才触发
const handleDateChange = (dates: [string, string] | null) => {
  if (dates) {
    // 清空快捷选择状态，表示用户手动选择了日期
    selectedQuick.value = ''
    dateRange.value = dates
    updateChartData()
  }
}

// 切换到编辑模式（显示真正的日期选择器）
const switchToEditMode = () => {
  selectedQuick.value = ''
}

// 格式化日期范围用于显示
const formatDateRange = (range: [string, string]) => {
  if (!range || !range[0] || !range[1]) {
    return t('chart.licenseTrend.datePicker.selectDateRange')
  }
  
  const startDate = new Date(range[0])
  const endDate = new Date(range[1])
  
  const formatDate = (date: Date) => {
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${month}-${day}`
  }
  
  return `${formatDate(startDate)} ${t('chart.licenseTrend.datePicker.rangeSeparator')} ${formatDate(endDate)}`
}

// 获取授权趋势数据
const fetchAuthorizationTrend = async (startDate: string, endDate: string) => {
  console.log('开始请求授权趋势数据:', { startDate, endDate })
  try {
    chartLoading.value = true
    const response = await getAuthorizationTrend({
      type: "custom",
      start_date:startDate,
      end_date:endDate
    })

    console.log('授权趋势API响应:', response)

    // 将API返回的trend-data转换为图表需要的格式
    const apiTrendData = response.data.trend_data as TrendDataItem[]
    console.log('趋势数据:', apiTrendData)

    trendData.value = apiTrendData.map(item => ({
      date: item.date,
      value: item.total_authorizations
    }))

    console.log('转换后的图表数据:', trendData.value)
  } catch (error: any) {
    console.error('获取授权趋势数据失败:', error)
    ElMessage.error(error?.backendMessage || '获取授权趋势数据失败')
  } finally {
    chartLoading.value = false
  }
}

// 更新图表数据
const updateChartData = () => {
  console.log('updateChartData 被调用，日期范围:', dateRange.value)
  fetchAuthorizationTrend(dateRange.value[0], dateRange.value[1])
}


onMounted(() => {
  console.log('LicenseTrendChart 组件已挂载，开始初始化')
  // 初始化时直接调用快捷选择，设置为本周
  handleQuickSelect('week')
  console.log('已调用 handleQuickSelect(week)')
})
</script>

<style lang="scss" scoped>
.license-trend-chart {
  background: var(--app-content-bg);
  border-radius: 0.42vw; /* 8px/1920 = 0.42vw */
  box-shadow: var(--app-shadow);
  overflow: hidden;
  height: 100%; /* 充满容器高度 */
  display: flex;
  flex-direction: column;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25vw 1.25vw 0; /* 24px/1920 = 1.25vw */
  flex-shrink: 0; /* 头部不收缩 */
  
  .chart-title {
    font-family: 'OPPOSans', sans-serif;
    font-size: 1.04vw; /* 20px/1920 = 1.04vw */
    font-weight: 400;
    color: var(--app-text-primary);
    margin: 0;
    line-height: 1.3;
  }
}

.time-selector {
  display: flex;
  align-items: center;
  gap: 12px;
  
  .quick-selector {
    display: flex;
    gap: 4px;
  }
  
  :deep(.el-button--small) {
    height: 32px;
    padding: 0 12px;
    font-size: 14px;
    border-radius: 4px;
  }
  
  :deep(.el-button--default) {
    background: var(--app-bg-color);
    border-color: var(--app-border-light);
    color: var(--app-text-secondary);
    
    &:hover {
      background: var(--app-border-color);
      border-color: var(--app-border-color);
    }
  }
  
  :deep(.el-button--primary) {
    background: #00C27C;
    border-color: #00C27C;
    
    &:hover {
      background: #019C7C;
      border-color: #019C7C;
    }
  }
  
  :deep(.el-date-editor) {
    height: 32px;
    
    .el-input__wrapper {
      border-radius: 4px;
    }
  }
  
  // 日期显示组件样式
  .date-display-wrapper {
    display: inline-block;
    
    .date-display {
      height: 32px;
      width: 200px;
      padding: 0 12px;
      border: 1px solid var(--app-border-color);
      border-radius: 4px;
      background: var(--app-content-bg);
      font-size: 14px;
      color: var(--app-text-regular);
      display: flex;
      align-items: center;
      justify-content: space-between;
      cursor: pointer;
      transition: all 0.2s ease;
      
      &:hover {
        border-color: var(--app-border-color);
      }
      
      &:focus {
        border-color: #019C7C;
        outline: none;
      }
      
      .date-icon {
        color: var(--app-text-secondary);
        font-size: 14px;
        transition: color 0.2s ease;
      }
      
      &:hover .date-icon {
        color: var(--app-text-secondary);
      }
    }
  }
}

.chart-container {
  padding: 1.35vw 1.25vw 1.25vw; /* 26px 24px 24px → vw */
  flex: 1; /* 占据剩余高度 */
  display: flex;
  min-height: 0; /* 防止flex溢出 */
  
  .trend-chart {
    width: 100%;
    height: 100%; /* 充满容器 */
    min-height: 12.8vw; /* 最小高度246px/1920 = 12.8vw，确保图表可读性 */
  }
}

// 响应式设计 - 移动端和平板使用固定像素单位
@media (max-width: 1024px) {
  .license-trend-chart {
    border-radius: 8px; /* 移动端使用固定像素 */
    box-shadow: var(--app-shadow);
  }
  
  .chart-header {
    padding: 16px 16px 0; /* 移动端固定像素 */
    
    .chart-title {
      font-size: 18px; /* 移动端固定字体 */
    }
  }
  
  .chart-container {
    padding: 16px; /* 移动端固定间距 */
    
    .trend-chart {
      height: 200px; /* 移动端固定高度 */
      min-height: unset; /* 清除vw最小高度 */
    }
  }
}

@media (max-width: 768px) {
  .chart-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0;
    padding: 16px 16px 0;
  }
  
  .chart-title {
    margin-bottom: 16px;
    font-size: 18px;
  }
  
  .time-selector {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 16px;
    
    // 移动端快捷按钮样式
    .mobile-quick {
      width: 100%;
      display: flex;
      justify-content: space-between;
      gap: 6px;
      
      :deep(.el-button--small) {
        flex: 1;
        font-size: 12px;
        height: 32px;
        padding: 0 8px;
        border-radius: 4px;
        min-width: 0;
        
        // 确保文字不被截断
        span {
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
      
      // 主要按钮样式
      :deep(.el-button--primary) {
        background: #019C7C;
        border-color: #019C7C;
        color: #fff;
      }
      
      // 默认按钮样式
      :deep(.el-button--default) {
        background: var(--app-bg-color);
        border-color: var(--app-border-light);
        color: var(--app-text-secondary);

        &:hover {
          background: var(--app-border-color);
          border-color: var(--app-border-color);
        }
      }
    }
  }
  
  .chart-container {
    padding: 16px;
    
    .trend-chart {
      height: 200px;
    }
  }
}

// 超小屏幕优化 (手机竖屏)
@media (max-width: 480px) {
  .chart-header {
    padding: 12px 12px 0;
  }
  
  .chart-title {
    font-size: 16px;
    margin-bottom: 12px;
  }
  
  .time-selector {
    gap: 12px;
    
    .mobile-quick {
      gap: 4px;
      
      :deep(.el-button--small) {
        font-size: 11px;
        height: 28px;
        padding: 0 6px;
      }
    }
  }
  
  .chart-container {
    padding: 12px;
    
    .trend-chart {
      height: 180px;
    }
  }
}
</style>