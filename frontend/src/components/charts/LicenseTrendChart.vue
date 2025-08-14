<template>
  <div class="license-trend-chart">
    <!-- 卡片头部 -->
    <div class="chart-header">
      <h3 class="chart-title">授权趋势</h3>
      
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
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
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
        ref="chartRef"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { use } from 'echarts/core'
import { useDevice } from '@/utils/useDevice'
import MobileDateRange from '@/components/common/MobileDateRange.vue'
import { Calendar } from '@element-plus/icons-vue'
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

// 快捷选择选项
const quickOptions = [
  { label: '本日', value: 'today' },
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' }
]

// 响应式数据
const { isMobile } = useDevice()
const selectedQuick = ref('week')
const dateRange = ref<[string, string]>(['2024-05-13', '2024-05-17']) // 初始值会被handleQuickSelect('week')覆盖
const chartRef = ref()
const datePickerRef = ref()

// 模拟授权趋势数据
const trendData = ref([
  { date: '2024-05-13', value: 8 },
  { date: '2024-05-14', value: 12 },
  { date: '2024-05-15', value: 15 },
  { date: '2024-05-16', value: 10 },
  { date: '2024-05-17', value: 18 }
])

// 图表配置
const chartOption = computed(() => {
  const dates = trendData.value.map(item => item.date.substring(5)) // 只显示月-日
  const values = trendData.value.map(item => item.value)
  
  return {
    grid: {
      left: 60,
      right: 30,
      top: 30,
      bottom: 50,
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
        color: '#1D1D1D',
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
        color: '#1D1D1D',
        fontSize: 12,
        fontFamily: 'Inter',
        align: 'right'
      },
      splitLine: {
        lineStyle: {
          color: '#EBEBEB',
          width: 1
        }
      }
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: '#FFFFFF',
      borderColor: '#FFFFFF',
      borderWidth: 0.5,
      borderRadius: 4,
      textStyle: {
        color: '#1D1D1D',
        fontSize: 12
      },
      extraCssText: 'box-shadow: 0px 4px 12px 0px rgba(59, 210, 180, 0.2); backdrop-filter: blur(4px);',
      formatter: (params: any) => {
        const data = params[0]
        return `授权数量: ${data.value}`
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
    return '请选择日期范围'
  }
  
  const startDate = new Date(range[0])
  const endDate = new Date(range[1])
  
  const formatDate = (date: Date) => {
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return `${month}-${day}`
  }
  
  return `${formatDate(startDate)} 至 ${formatDate(endDate)}`
}

// 更新图表数据
const updateChartData = () => {
  // 这里可以根据选择的日期范围调用API获取实际数据
  // 现在使用模拟数据
  const mockData = generateMockData(dateRange.value[0], dateRange.value[1])
  trendData.value = mockData
}

// 生成模拟数据
const generateMockData = (startDate: string, endDate: string) => {
  const start = new Date(startDate)
  const end = new Date(endDate)
  const data = []
  
  const currentDate = new Date(start)
  while (currentDate <= end) {
    data.push({
      date: currentDate.toISOString().split('T')[0],
      value: Math.floor(Math.random() * 15) + 5 // 5-20之间的随机数
    })
    currentDate.setDate(currentDate.getDate() + 1)
  }
  
  return data
}

onMounted(() => {
  // 初始化时直接调用快捷选择，设置为本周
  handleQuickSelect('week')
})
</script>

<style lang="scss" scoped>
.license-trend-chart {
  background: #FFFFFF;
  border-radius: 8px;
  box-shadow: 0px 0px 4px 0px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 24px 0;
  
  .chart-title {
    font-family: 'OPPOSans', sans-serif;
    font-size: 20px;
    font-weight: 400;
    color: #1D1D1D;
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
    background: #F7F8FA;
    border-color: #F7F8FA;
    color: #666;
    
    &:hover {
      background: #E9ECEF;
      border-color: #E9ECEF;
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
      border: 1px solid #dcdfe6;
      border-radius: 4px;
      background: #fff;
      font-size: 14px;
      color: #606266;
      display: flex;
      align-items: center;
      justify-content: space-between;
      cursor: pointer;
      transition: all 0.2s ease;
      
      &:hover {
        border-color: #c0c4cc;
      }
      
      &:focus {
        border-color: #019C7C;
        outline: none;
      }
      
      .date-icon {
        color: #c0c4cc;
        font-size: 14px;
        transition: color 0.2s ease;
      }
      
      &:hover .date-icon {
        color: #909399;
      }
    }
  }
}

.chart-container {
  padding: 26px 24px 24px;
  
  .trend-chart {
    width: 100%;
    height: 246px;
  }
}

// 响应式设计
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
        background: #f7f8fa;
        border-color: #f7f8fa;
        color: #666;
        
        &:hover {
          background: #e9ecef;
          border-color: #e9ecef;
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