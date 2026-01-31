<template>
    <div class="loading_box">
        <el-dialog v-model="visible" :title="t('packages.edit.title')" width="700px" class="package-edit-dialog"
            destroy-on-close>
            <el-form ref="formRef" :model="form" :rules="rules" label-position="top" class="edit-form">
                <div class="form-row">
                    <el-form-item :label="t('packages.edit.name')" prop="name" class="flex-1">
                        <el-input v-model="form.name" />
                    </el-form-item>
                    <el-form-item :label="t('packages.edit.type')" prop="type" class="flex-1">
                        <el-select v-model="form.type" class="w-full">
                            <el-option label="试用版" value="trial" />
                            <el-option label="基础版" value="basic" />
                            <el-option label="专业版" value="pro" />
                            <el-option label="定制版" value="custom" />
                        </el-select>
                    </el-form-item>
                </div>

                <div class="form-row">
                    <el-form-item :label="t('packages.edit.price')" prop="price" class="flex-1">
                        <el-input v-model="form.price" />
                    </el-form-item>
                    <el-form-item :label="t('packages.edit.cycle')" prop="cycle" class="flex-1">
                        <el-input v-model="form.cycle" />
                    </el-form-item>
                </div>

                <el-form-item :label="t('packages.edit.description')" prop="description">
                    <el-input v-model="form.description" type="textarea" :rows="4" maxlength="500" show-word-limit
                        :placeholder="t('packages.edit.descPlaceholder')" />
                </el-form-item>

                <div class="features-section">
                    <label class="section-label">{{ t('packages.edit.features') }}</label>
                    <div class="features-list">
                        <div v-for="(feature, index) in form.features" :key="index" class="feature-item">
                            <el-input v-model="form.features[index]"
                                :placeholder="t('packages.edit.featurePlaceholder')" />
                            <el-button class="btn-remove" @click="removeFeature(index)">
                                <el-icon>
                                    <Minus />
                                </el-icon>
                            </el-button>
                        </div>
                        <el-button class="btn-add" @click="addFeature">
                            <el-icon>
                                <Plus />
                            </el-icon>
                        </el-button>
                    </div>
                </div>

                <div class="form-row mt-20">
                    <el-form-item :label="t('packages.edit.status')" prop="status" class="flex-1">
                        <el-select v-model="form.status" class="w-full">
                            <el-option :label="t('packages.status.enabled')" value="enabled" />
                            <el-option :label="t('packages.status.disabled')" value="disabled" />
                        </el-select>
                    </el-form-item>
                    <el-form-item :label="t('packages.edit.sort')" prop="sort" class="flex-1">
                        <el-input-number v-model="form.sort" :min="1" class="w-full" controls-position="right" />
                    </el-form-item>
                </div>

                <el-form-item :label="t('packages.edit.remark')" prop="remark">
                    <el-input v-model="form.remark" type="textarea" :rows="4" maxlength="500" show-word-limit
                        :placeholder="t('packages.edit.remarkPlaceholder')" />
                </el-form-item>
            </el-form>

            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="visible = false">{{ t('common.cancel') }}</el-button>
                    <el-button type="primary" class="btn-save" @click="handleSave">
                        {{ t('packages.edit.save') }}
                    </el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import { Plus, Minus } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const props = defineProps<{
    modelValue: boolean
    data: any
}>()

const emit = defineEmits(['update:modelValue', 'save'])

const { t } = useI18n()
const visible = ref(props.modelValue)
const formRef = ref<FormInstance>()

const form = reactive({
    name: '',
    type: '',
    price: '',
    cycle: '',
    description: '',
    features: [] as string[],
    status: 'enabled',
    sort: 1,
    remark: ''
})

const rules = reactive<FormRules>({
    name: [{ required: true, message: t('common.required'), trigger: 'blur' }],
    type: [{ required: true, message: t('common.required'), trigger: 'change' }],
    price: [{ required: true, message: t('common.required'), trigger: 'blur' }],
    cycle: [{ required: true, message: t('common.required'), trigger: 'blur' }],
    description: [{ required: true, message: t('common.required'), trigger: 'blur' }],
    status: [{ required: true, message: t('common.required'), trigger: 'change' }],
    sort: [{ required: true, message: t('common.required'), trigger: 'blur' }]
})

watch(() => props.modelValue, (val) => {
    visible.value = val
    if (val && props.data) {
        form.name = props.data.name || ''
        form.type = props.data.type || 'trial'
        form.price = props.data.price || ''
        form.cycle = props.data.cycle || ''
        form.description = props.data.description || '体验完整功能，适合新用户试用评估'
        form.features = props.data.features || ['全部功能试用', '支持1台设备激活', '基础功能支持', '永久授权', '高级功能', '离线授权支持']
        form.status = props.data.status || 'enabled'
        form.sort = props.data.sort || 6
        form.remark = props.data.remark || ''
    }
})

watch(visible, (val) => {
    emit('update:modelValue', val)
})

const addFeature = () => {
    form.features.push('')
}

const removeFeature = (index: number) => {
    form.features.splice(index, 1)
}

const handleSave = async () => {
    if (!formRef.value) return
    await formRef.value.validate((valid) => {
        if (valid) {
            emit('save', { ...props.data, ...form })
            visible.value = false
        }
    })
}
</script>

<style lang="scss" scoped>

:deep(.el-dialog__headerbtn) {
    top: 10px !important;
}

:deep(.el-dialog) {
    border-radius: 8px;
    overflow: hidden;
    padding: 0;

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
        }
    }

    .el-dialog__body {
        padding: 24px;
    }
}

.edit-form {
    .form-row {
        display: flex;
        gap: 24px;
        margin-bottom: 8px;
    }

    .flex-1 {
        flex: 1;
    }

    :deep(.el-form-item__label) {
        font-weight: 500;
        color: #333;
        padding-bottom: 8px;

        &::before {
            color: #f56c6c;
            margin-right: 4px;
        }
    }

    :deep(.el-input__wrapper),
    :deep(.el-textarea__inner) {
        background-color: #fff;
        border-radius: 4px;
    }
}

.features-section {
    .section-label {
        display: block;
        font-weight: 500;
        color: #333;
        margin-bottom: 12px;
        font-size: 14px;

        &::before {
            content: '*';
            color: #f56c6c;
            margin-right: 4px;
        }
    }
}

.features-list {
    background: #f9f9f9;
    border: 1px solid #eee;
    border-radius: 4px;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;

    .feature-item {
        display: flex;
        gap: 12px;
        align-items: center;
    }

    .btn-remove {
        width: 32px;
        height: 32px;
        padding: 0;
        border: 1px dashed #ccc;
        color: #999;

        &:hover {
            color: #f56c6c;
            border-color: #f56c6c;
        }
    }

    .btn-add {
        width: 44px;
        height: 44px;
        padding: 0;
        border: 1px dashed #ccc;
        color: #999;
        align-self: flex-start;

        &:hover {
            color: #00a870;
            border-color: #00a870;
        }
    }
}

.w-full {
    width: 100%;
}

.mt-20 {
    margin-top: 20px;
}
:deep(.el-dialog__footer){
    padding-top: 0px !important;
}

.dialog-footer {
    display: flex;
    justify-content: center;
    gap: 16px;
    margin-top: 0px !important;
    padding-bottom: 20px !important;

    .el-button {
        padding: 6px 32px;
        height: 40px;
        font-size: 14px;
        border-radius: 4px;
    }

    .btn-save {
        background-color: #00a870 !important;
        border-color: #00a870 !important;
        color: #fff !important;
    }
}
</style>
