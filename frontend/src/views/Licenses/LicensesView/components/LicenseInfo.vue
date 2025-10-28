<template>
  <div class="license-info-container" v-loading="loading">
    <!-- 新增许可证按钮 -->
    <div class="action-bar">
      <el-button type="primary" @click="handleAddLicense">新增许可证</el-button>
    </div>

    <!-- 设备列表 -->
    <div v-if="!loading && devices.length > 0" class="devices-list">
      <div v-for="(device, index) in devices" :key="device.id" class="device-card">
        <!-- 设备头部 -->
        <div class="device-header">
          <div class="device-title-row">
            <span class="device-title">设备{{ index + 1 }}</span>
            <el-tag :type="device.is_online ? 'success' : 'danger'" size="small">
              {{ device.is_online_display}}
            </el-tag>
          </div>
          <div class="device-info-row">
            <span class="device-info-item">设备地址：{{ device.last_online_ip || device.activation_ip || '-' }}</span>
            <span class="device-info-item">最后活跃时间：{{ formatDateTime(device.last_heartbeat || '') }}</span>
          </div>
        </div>

        <!-- 硬件信息 -->
        <div class="info-section hardware-section">
          <div class="section-header">
            <span class="section-title">硬件信息</span>
            <span class="section-subtitle">（绑定信息）</span>
          </div>
          <div class="info-grid">
            <div class="info-row">
              <div class="info-label">机器序列号</div>
              <div class="info-value">{{ device.device_info?.machine_code || '-' }}</div>
            </div>
            <div class="info-row">
              <div class="info-label">MAC地址</div>
              <div class="info-value">{{ device.device_info?.mac_address || '-' }}</div>
            </div>
            <div class="info-row">
              <div class="info-label">CPU ID</div>
              <div class="info-value">{{ device.device_info?.cpu_id || '-' }}</div>
            </div>
            <div class="info-row">
              <div class="info-label">硬件哈希</div>
              <div class="info-value">{{ device.hardware_fingerprint || '-' }}</div>
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
              <div class="license-label">许可证密钥</div>
              <div class="license-value">{{ device.license_key || '-' }}</div>
            </div>
            <div class="license-row">
              <div class="license-label">许可证状态</div>
              <div class="license-value status-text">{{ device.status_display || '此许可证已生成并绑定到指定硬件，可用于系统运行验证' }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-if="!loading && devices.length === 0" class="empty-state">
      <el-empty description="暂无设备许可证信息">
        <el-button type="primary" @click="handleAddLicense">新增许可证</el-button>
      </el-empty>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { AuthorizationCode, LicenseDevice } from '@/api/license'
import { getLicenseDevices, getLicenseDeviceDetail } from '@/api/license'
import { formatDate } from '@/utils/date'

interface Props {
  licenseData: AuthorizationCode | null
}

const props = defineProps<Props>()

const devices = ref<LicenseDevice[]>([])
const loading = ref(false)

// 获取设备列表
const fetchDevices = async () => {
  if (!props.licenseData?.id) {
    devices.value = []
    return
  }

  loading.value = true
  try {
    const response = await getLicenseDevices({
      authorization_code_id: props.licenseData.id,
      page: 1,
      page_size: 100
    })

    if (response.data?.list && response.data.list.length > 0) {
      // 检查第一个设备是否有 device_info，如果没有则需要获取详情
      const needDetails = !response.data.list[0].device_info

      if (needDetails) {
        // 并发获取所有设备的详细信息
        const detailPromises = response.data.list.map(device =>
          getLicenseDeviceDetail(device.id)
            .then(res => res.data)
            .catch(err => {
              console.error(`获取设备 ${device.id} 详情失败:`, err)
              return device // 如果获取失败，使用原数据
            })
        )

        const detailedDevices = await Promise.all(detailPromises)
        devices.value = detailedDevices as LicenseDevice[]
      } else {
        devices.value = response.data.list
      }
    } else {
      devices.value = []
    }
  } catch (error) {
    console.error('获取设备列表失败:', error)
    ElMessage.error('获取设备列表失败')
    devices.value = []
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchDevices()
})

// 监听licenseData变化，重新获取数据
watch(() => props.licenseData?.id, () => {
  fetchDevices()
})

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
    flex-direction: column;
    gap: 12px;

    .info-row {
      display: flex;
      align-items: center;
      gap: 30px;

      .info-label {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        white-space: nowrap;
        min-width: 100px;
      }

      .info-value {
        flex: 1;
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
}

.license-section {
  .license-content {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .license-row {
      display: flex;
      align-items: center;
      gap: 30px;

      .license-label {
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        white-space: nowrap;
        min-width: 100px;
      }

      .license-value {
        flex: 1;
        font-family: 'Source Han Sans CN', sans-serif;
        font-size: 14px;
        font-weight: 400;
        line-height: 22px;
        color: #1D1D1D;
        padding: 8px 12px;
        background: rgba(136, 165, 209, 0.2);
        border-radius: 8px;
        word-break: break-all;

        &.status-text {
          color: #4763FF;
          background: transparent;
          padding: 0;
        }
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
      .info-row {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;

        .info-label {
          min-width: auto;
        }

        .info-value {
          width: 100%;
        }
      }
    }
  }

  .license-section {
    .license-content {
      .license-row {
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;

        .license-label {
          min-width: auto;
        }

        .license-value {
          width: 100%;
        }
      }
    }
  }
}
</style>
