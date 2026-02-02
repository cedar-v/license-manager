<template>
    <div class="dialog_box">
        <el-dialog v-model="visible" :title="t('enterpriseLeads.detail.title', { company: detailData?.company_name })"
            width="800px" class="lead-detail-dialog" destroy-on-close v-loading="loading">
            <div v-if="detailData" class="detail-content">
                <!-- 企业基本信息 -->
                <div class="detail-section">
                    <h3 class="section-title">{{ t('enterpriseLeads.detail.basicInfo') }}</h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.company') }}：</span>
                            <span class="value">{{ detailData.company_name }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.contact') }}：</span>
                            <span class="value">{{ detailData.contact_name }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.phone') }}：</span>
                            <span class="value">{{ detailData.contact_phone }}</span>
                        </div>
                    </div>
                </div>

                <!-- 需求信息 -->
                <div class="detail-section">
                    <h3 class="section-title">{{ t('enterpriseLeads.detail.requirementInfo') }}</h3>
                    <div class="info-list">
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.detail.email') }}：</span>
                            <span class="value">{{ detailData.contact_email || '-' }}</span>
                        </div>
                        <div class="info-item block">
                            <span class="label">{{ t('enterpriseLeads.detail.description') }}：</span>
                            <span class="value">{{ detailData.requirement || '-' }}</span>
                        </div>
                        <div class="info-item block">
                            <span class="label">{{ t('enterpriseLeads.detail.otherInfo') }}：</span>
                            <span class="value">{{ detailData.extra_info || '-' }}</span>
                        </div>
                    </div>
                </div>

                <!-- 跟进信息 -->
                <div class="detail-section">
                    <h3 class="section-title">{{ t('enterpriseLeads.detail.followUpInfo') }}</h3>
                    <div class="info-grid">
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.submittedAt') }}：</span>
                            <span class="value">{{ detailData.created_at }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.table.status') }}：</span>
                            <span class="status-value" :class="detailData.status">{{
                                t(`enterpriseLeads.status.${detailData.status}`) }}</span>
                        </div>
                        <div class="info-item">
                            <span class="label">{{ t('enterpriseLeads.detail.followUpDate') }}：</span>
                            <span class="value">{{ detailData.follow_up_date || '-' }}</span>
                        </div>
                    </div>
                    <div class="info-list mt-12">
                        <div class="info-item block">
                            <span class="label">{{ t('enterpriseLeads.detail.followUpRecord') }}：</span>
                            <span class="value">{{ detailData.follow_up_record || '-' }}</span>
                        </div>
                        <div class="info-item block">
                            <span class="label">{{ t('enterpriseLeads.detail.internalRemark') }}：</span>
                            <span class="value">{{ detailData.internal_note || '-' }}</span>
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
import { getLeadDetail, type Lead } from '@/api/lead'
import { ElMessage } from 'element-plus'
import { formatDateTime } from '@/utils/date'

const props = defineProps<{
    modelValue: boolean
    id: string | number | null
}>()

const emit = defineEmits(['update:modelValue'])

const { t } = useI18n()
const visible = ref(props.modelValue)
const loading = ref(false)
const detailData = ref<Lead | null>(null)

const fetchDetail = async (id: string | number) => {
    loading.value = true
    try {
        const res = await getLeadDetail(id)
        if (res.code === '000000' && res.data) {
            const data = res.data
            detailData.value = {
                ...data,
                created_at: formatDateTime(data.created_at || ''),
                follow_up_date: data.follow_up_date ? formatDateTime(data.follow_up_date) : null
            }
        }
    } catch (error: any) {
        console.error('Fetch lead detail error:', error)
        ElMessage.error(error.backendMessage || t('enterpriseLeads.messages.fetchDetailError'))
    } finally {
        loading.value = false
    }
}

watch(() => props.modelValue, (val) => {
    visible.value = val
    if (val && props.id) {
        fetchDetail(props.id)
    }
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
