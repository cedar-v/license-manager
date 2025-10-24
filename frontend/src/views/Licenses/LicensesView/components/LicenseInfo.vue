<template>
  <div class="license-info-container">
    <!-- 新增许可证按钮 -->
    <div class="action-bar">
      <el-button type="primary" @click="handleAddLicense">新增许可证</el-button>
    </div>

    <!-- 设备列表 -->
    <div v-if="devices.length > 0" class="devices-list">
      <div v-for="(device, index) in devices" :key="device.id" class="device-card">
        <!-- 设备头部 -->
        <div class="device-header">
          <div class="device-title-row">
            <span class="device-title">设备{{ index + 1 }}</span>
            <el-tag :type="device.status === 'online' ? 'success' : 'danger'" size="small">
              {{ device.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </div>
          <div class="device-info-row">
            <span class="device-info-item">设备地址：{{ device.address || '-' }}</span>
            <span class="device-info-item">最后活跃时间：{{ formatDateTime(device.lastActiveTime) }}</span>
          </div>
        </div>

        <!-- 硬件信息 -->
        <div class="info-section hardware-section">
          <div class="section-header">
            <span class="section-title">硬件信息</span>
            <span class="section-subtitle">（绑定信息）</span>
          </div>
          <div class="info-grid">
            <div class="info-col">
              <div class="info-label">机器序列号</div>
              <div class="info-label">MAC地址</div>
              <div class="info-label">CPU ID</div>
              <div class="info-label">硬件哈希</div>
            </div>
            <div class="info-col values">
              <div class="info-value">{{ device.hardware?.serialNumber || '-' }}</div>
              <div class="info-value">{{ device.hardware?.macAddress || '-' }}</div>
              <div class="info-value">{{ device.hardware?.cpuId || '-' }}</div>
              <div class="info-value hash-value">{{ device.hardware?.hardwareHash || '-' }}</div>
            </div>
          </div>
        </div>

        <!-- 软件许可证 -->
        <div class="info-section license-section">
          <div class="section-header">
            <span class="section-title">软件许可证</span>
            <span class="section-subtitle">（用于系统验证）</span>
          </div>
          <div class="license-content">
            <div class="license-row">
              <span class="license-label">许可证</span>
              <div class="license-value-wrapper">
                <div class="license-value">{{ device.license?.certificate || '-' }}</div>
              </div>
            </div>
            <div class="license-status">
              <span class="license-label">许可证状态</span>
              <span class="license-status-text">{{ device.license?.status || '此许可证已生成并绑定到指定硬件，可用于系统运行验证' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <el-empty description="暂无设备许可证信息">
        <el-button type="primary" @click="handleAddLicense">新增许可证</el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { AuthorizationCode } from '@/api/license'
import { formatDate } from '@/utils/date'

interface Props {
  licenseData: AuthorizationCode | null
}

interface DeviceHardware {
  serialNumber: string
  macAddress: string
  cpuId: string
  hardwareHash: string
}

interface DeviceLicense {
  certificate: string
  status: string
}

interface Device {
  id: number
  status: 'online' | 'offline'
  address: string
  lastActiveTime: string
  hardware: DeviceHardware
  license: DeviceLicense
}

defineProps<Props>()

// 模拟设备数据，实际应从API获取
const devices = ref<Device[]>([
  {
    id: 1,
    status: 'online',
    address: '125.665.2.612',
    lastActiveTime: '2025-09-14 10:10:25',
    hardware: {
      serialNumber: 'VMware-56 4d 74 f6 88 a5 09 f1-81 c4 06 6b e1 ae 59 41',
      macAddress: '00:0c:29:ae:59:41',
      cpuId: 'Intel(R)-Core(TM)-i7-8700-CPU-@-3.20GHz',
      hardwareHash: 'ePZtmBWUgp8iIEKtROmt21cOwMeQqhzCsEc6j6czYzw='
    },
    license: {
      certificate: 'eyJjdXN0b21lcl9pZCI6MTAsImV4cGlyZV9kYXlzIjoxODAsImdlbmVyYXRlZF9hdCI6MTc1NDAyOTE3Nn1ONATrOXHImU+fQaZZpmPoQNjEGpemoqk',
      status: '此许可证已生成并绑定到指定硬件，可用于系统运行验证'
    }
  },
  {
    id: 2,
    status: 'offline',
    address: '125.665.2.612',
    lastActiveTime: '2025-09-14 10:10:25',
    hardware: {
      serialNumber: 'VMware-56 4d 74 f6 88 a5 09 f1-81 c4 06 6b e1 ae 59 41',
      macAddress: '00:0c:29:ae:59:41',
      cpuId: 'Intel(R)-Core(TM)-i7-8700-CPU-@-3.20GHz',
      hardwareHash: 'ePZtmBWUgp8iIEKtROmt21cOwMeQqhzCsEc6j6czYzw='
    },
    license: {
      certificate: 'eyJjdXN0b21lcl9pZCI6MTAsImV4cGlyZV9kYXlzIjoxODAsImdlbmVyYXRlZF9hdCI6MTc1NDAyOTE3Nn1ONATrOXHImU+fQaZZpmPoQNjEGpemoqk',
      status: '此许可证已生成并绑定到指定硬件，可用于系统运行验证'
    }
  }
])

// 格式化日期时间
const formatDateTime = (dateTime: string) => {
  if (!dateTime) return '-'
  return formatDate(dateTime)
}

// 处理新增许可证
const handleAddLicense = () => {
  // TODO: 实现新增许可证逻辑
  console.log('新增许可证')
}
</script>

<style lang="scss" scoped>
.license-info-container {
  width: 100%;
}

.action-bar {
  margin-bottom: 16px;

  :deep(.el-button) {
    padding: 7px 16px;
    font-family: 'Source Han Sans CN', sans-serif;
    font-size: 14px;
    font-weight: 500;
    line-height: 21px;
    border-radius: 4px;
  }
}

.devices-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.device-card {
  background: #FFFFFF;
  border: 1px solid #E2E2E2;
  border-radius: 4px;
  padding: 16px 20px;
}

.device-header {
  margin-bottom: 20px;

  .device-title-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;

    .device-title {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 16px;
      font-weight: 700;
      line-height: 18px;
      color: #1D1D1D;
    }

    :deep(.el-tag) {
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

  .device-info-row {
    display: flex;
    gap: 40px;

    .device-info-item {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 22px;
      color: #1D1D1D;
    }
  }
}

.info-section {
  background: rgba(247, 248, 250, 0.7);
  border: 1px solid rgba(226, 226, 226, 0.6);
  border-radius: 4px;
  padding: 16px;
  margin-bottom: 16px;

  &:last-child {
    margin-bottom: 0;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 4px;
    margin-bottom: 16px;

    .section-title {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 16px;
      font-weight: 700;
      line-height: 18px;
      color: #1D1D1D;
    }

    .section-subtitle {
      font-family: 'Source Han Sans CN', sans-serif;
      font-size: 14px;
      font-weight: 400;
      line-height: 18px;
      color: #666666;
    }
  }
}

.hardware-section {
  .info-grid {
    display: flex;
    gap: 30px;

    .info-col {
      display: flex;
      flex-direction: column;
      gap: 12px;

      &.values {
        flex: 1;
      }

      .info-label {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        white-space: nowrap;
      }

      .info-value {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        padding: 8px 12px;
        background: rgba(136, 165, 209, 0.2);
        border-radius: 8px;
        word-break: break-all;

        &.hash-value {
          margin-top: 5px;
        }
      }
    }
  }
}

.license-section {
  .license-content {
    .license-row {
      display: flex;
      gap: 30px;
      margin-bottom: 24px;

      .license-label {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        white-space: nowrap;
      }

      .license-value-wrapper {
        flex: 1;

        .license-value {
          font-family: 'Source Han Sans CN', sans-serif;
          font-size: 14px;
          font-weight: 400;
          line-height: 22px;
          color: #1D1D1D;
          padding: 8px 12px;
          background: rgba(136, 165, 209, 0.2);
          border-radius: 8px;
          word-break: break-all;
        }
      }
    }

    .license-status {
      display: flex;
      gap: 30px;

      .license-label {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        white-space: nowrap;
      }

      .license-status-text {
        flex: 1;
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #4763FF;
      }
    }
  }
}

.empty-state {
  padding: 40px;
  text-align: center;
}

// 响应式布局
@media (max-width: 768px) {
  .device-header {
    .device-info-row {
      flex-direction: column;
      gap: 8px;
    }
  }

  .hardware-section {
    .info-grid {
      flex-direction: column;
      gap: 16px;
    }
  }

  .license-section {
    .license-content {
      .license-row,
      .license-status {
        flex-direction: column;
        gap: 8px;
      }
    }
  }
}
</style>
