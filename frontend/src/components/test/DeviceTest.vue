<template>
  <div class="device-test">
    <h3>设备检测测试</h3>
    <div class="info-grid">
      <div class="info-item">
        <label>屏幕宽度:</label>
        <span>{{ screenWidth }}px</span>
      </div>
      <div class="info-item">
        <label>屏幕高度:</label>
        <span>{{ screenHeight }}px</span>
      </div>
      <div class="info-item">
        <label>设备类型:</label>
        <span>{{ deviceType }}</span>
      </div>
      <div class="info-item">
        <label>屏幕方向:</label>
        <span>{{ orientation }}</span>
      </div>
      <div class="info-item">
        <label>是否触摸设备:</label>
        <span>{{ isTouchDevice ? '是' : '否' }}</span>
      </div>
      <div class="info-item">
        <label>更新次数:</label>
        <span>{{ updateCount }}</span>
      </div>
    </div>
    
    <div class="status">
      <p>调整浏览器窗口大小来测试防抖功能</p>
      <p>防抖延迟: 150ms（停止调整窗口后150ms才更新）</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useDevice } from '@/utils/useDevice'

const updateCount = ref(0)

const {
  screenWidth,
  screenHeight,
  deviceType,
  orientation,
  isTouchDevice
} = useDevice()

// 监听屏幕宽度变化，统计更新次数
watch([screenWidth, screenHeight], () => {
  updateCount.value++
}, { immediate: true })
</script>

<style scoped>
.device-test {
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  background: white;
  margin: 20px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 12px;
  margin: 16px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  padding: 8px 12px;
  background: #f5f5f5;
  border-radius: 4px;
}

.info-item label {
  font-weight: 600;
  color: #333;
}

.info-item span {
  color: #666;
  font-family: monospace;
}

.status {
  margin-top: 16px;
  padding: 12px;
  background: #e3f2fd;
  border-radius: 4px;
}

.status p {
  margin: 4px 0;
  color: #1976d2;
  font-size: 14px;
}
</style>