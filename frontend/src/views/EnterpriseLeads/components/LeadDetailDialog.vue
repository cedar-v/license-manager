<template>
    <div class="dialog_box">
        <el-dialog v-model="visible" :title="t('enterpriseLeads.detail.title', { company: data?.company })"
            width="800px" class="lead-detail-dialog" destroy-on-close>
            <div v-if="data" class="detail-content">
                <!-- 企业基本信息 -->
                <div class="detail-section">
                    <h3 class="section-title">{{ t('enterpriseLeads.detail.basicInfo') }}</h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.company') }}：</span>
                            <span class="value">{{ data.company }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.contact') }}：</span>
                            <span class="value">{{ data.contact }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.phone') }}：</span>
                            <span class="value">{{ data.phone }}</span>
                        </div>
                    </div>
                </div>

                <!-- 需求信息 -->
                <div class="detail-section">
                    <h3 class="section-title">{{ t('enterpriseLeads.detail.requirementInfo') }}</h3>
                    <div class="info-list">
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.detail.email') }}：</span>
                            <span class="value">{{ data.email || 'wangming@sh-smart.com' }}</span>
                        </div>
                        <div class="info-item block">
                            <span class="label">{{ t('enterpriseLeads.detail.description') }}：</span>
                            <span class="value">{{ data.description ||
                                '我们需要为生产线上的500台设备部署专业版授权，同时为办公室的200台电脑部署基础版授权。需要支持离线授权功能，部署时间为下个月初。' }}</span>
                        </div>
                        <div class="info-item block">
                            <span class="label">{{ t('enterpriseLeads.detail.otherInfo') }}：</span>
                            <span class="value">{{ data.otherInfo || '预算范围在50-80万之间，决策时间约2周。' }}</span>
                        </div>
                    </div>
                </div>

                <!-- 跟进信息 -->
                <div class="detail-section">
                    <h3 class="section-title">{{ t('enterpriseLeads.detail.followUpInfo') }}</h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.submittedAt') }}：</span>
                            <span class="value">{{ data.submittedAt }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.status') }}：</span>
                            <span class="status-value" :class="data.status">{{
                                t(`enterpriseLeads.status.${data.status}`) }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.detail.followUpDate') }}：</span>
                            <span class="value">{{ data.followUpDate || '2025-12-21' }}</span>
                        </div>
                    </div>
                    <div class="info-list mt-12">
                        <div class="info-item block">
                            <span class="label">{{ t('enterpriseLeads.detail.followUpRecord') }}：</span>
                            <span class="value">{{ data.followUpRecord || '已初步沟通，客户对专业版的AI功能非常感兴趣，需要提供详细的功能演示。'
                                }}</span>
                        </div>
                        <div class="info-item block">
                            <span class="label">{{ t('enterpriseLeads.detail.internalRemark') }}：</span>
                            <span class="value">{{ data.internalRemark || '高价值客户，需要安排技术专家进行演示。' }}</span>
                        </div>
                    </div>
                </div>
            </div>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'

const props = defineProps<{
    modelValue: boolean
    data: any
}>()

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const visible = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
    visible.value = val
})

watch(visible, (val) => {
    emit('update:modelValue', val)
})
</script>

<style lang="scss" scoped>
:deep(.el-dialog__headerbtn) {
    top: 10px !important;
}

:deep(.el-dialog) {
    border-radius: 8px;
    overflow: hidden;
    padding: 0 !important;


    .el-dialog__header {
        margin-right: 0;
        padding: 20px 24px;
        background: linear-gradient(90deg, #00928A 0%, #00D19E 100%) !important;
        border-bottom: none;
        display: flex;
        align-items: center;

        .el-dialog__title {
            color: #fff !important;
            font-size: 18px;
            font-weight: 600;
        }

        .el-dialog__headerbtn {
            top: 20px;

            .el-dialog__close {
                color: #fff !important;
                font-size: 20px;
            }

            &:hover .el-dialog__close {
                color: rgba(255, 255, 255, 0.8) !important;
            }
        }
    }

    .el-dialog__body {
        padding: 24px;
    }
}


.detail-content {
    padding: 8px 0;
}

.detail-section {
    margin-bottom: 24px;
    border: 1px solid #E9E9E9;

    &:last-child {
        margin-bottom: 0;
    }
}

.section-title {
    font-size: 14px;
    font-weight: 600;
    color: #333;
    margin-bottom: 16px;
    padding-left: 0;
    position: relative;
    background-color: #FAFAFA;
    padding: 11px 24px;
    box-sizing: border-box;
}

.info-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 16px;
     padding: 11px 24px;
    box-sizing: border-box;
   
}

.info-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
     padding: 11px 24px;
    box-sizing: border-box;
}

.info-item {
    display: flex;
    font-size: 14px;
    line-height: 1.6;

    &.block {
        flex-direction: row;
        align-items: flex-start;
    }

    .label {
        color: #666;
        flex-shrink: 0;
        width: 100px;
    }

    .value {
        color: #333;
    }

    .status-value {
        font-weight: 500;

        &.contacting {
            color: #409eff;
        }

        &.pending {
            color: #e6a23c;
        }

        &.completed {
            color: #333;
        }

        &.rejected {
            color: #999;
        }
    }
}

.mt-12 {
    margin-top: 12px;
}
</style>
