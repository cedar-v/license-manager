<template>
 <div>
     <el-dialog
    v-model="visible"
    :title="t('enterpriseLeads.edit.title', { company: data?.company })"
    width="800px"
    class="lead-edit-dialog"
    destroy-on-close
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-position="top"
      class="edit-form"
    >
      <div class="form-row">
        <el-form-item :label="t('enterpriseLeads.table.status')" prop="status" class="flex-1">
          <el-select v-model="form.status" class="w-full" :placeholder="t('enterpriseLeads.filter.statusPlaceholder')">
            <el-option :label="t('enterpriseLeads.status.pending')" value="pending" />
            <el-option :label="t('enterpriseLeads.status.contacting')" value="contacting" />
            <el-option :label="t('enterpriseLeads.status.completed')" value="completed" />
            <el-option :label="t('enterpriseLeads.status.rejected')" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('enterpriseLeads.detail.followUpDate')" prop="followUpDate" class="flex-1">
          <el-date-picker
            v-model="form.followUpDate"
            type="date"
            class="w-full"
            :placeholder="t('enterpriseLeads.edit.datePlaceholder')"
            format="YYYY/MM/DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
      </div>

      <el-form-item :label="t('enterpriseLeads.detail.followUpRecord')" prop="followUpRecord">
        <el-input
          v-model="form.followUpRecord"
          type="textarea"
          :rows="4"
          maxlength="500"
          show-word-limit
          :placeholder="t('enterpriseLeads.edit.recordPlaceholder')"
        />
      </el-form-item>

      <el-form-item :label="t('enterpriseLeads.detail.internalRemark')" prop="internalRemark">
        <el-input
          v-model="form.internalRemark"
          type="textarea"
          :rows="4"
          maxlength="500"
          show-word-limit
          :placeholder="t('enterpriseLeads.edit.remarkPlaceholder')"
        />
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="visible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" class="btn-save" @click="handleSave">
          {{ t('enterpriseLeads.edit.save') }}
        </el-button>
      </div>
    </template>
  </el-dialog>
 </div>
</template>

<script setup lang="ts">
import { ref, watch, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
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
  status: '',
  followUpDate: '',
  followUpRecord: '',
  internalRemark: ''
})

const rules = reactive<FormRules>({
  status: [{ required: true, message: t('enterpriseLeads.edit.statusRequired'), trigger: 'change' }]
})

watch(() => props.modelValue, (val) => {
  visible.value = val
  if (val && props.data) {
    form.status = props.data.status
    form.followUpDate = props.data.followUpDate || '2026-01-14'
    form.followUpRecord = props.data.followUpRecord || ''
    form.internalRemark = props.data.internalRemark || ''
  }
})

watch(visible, (val) => {
  emit('update:modelValue', val)
})

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
      &:hover .el-dialog__close {
        color: rgba(255, 255, 255, 0.8) !important;
      }
    }
  }

  .el-dialog__body {
    padding: 24px;
  }

  .el-dialog__footer {
    padding: 0 24px 24px;
    border-top: none;
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
  }

  :deep(.el-input__wrapper), :deep(.el-textarea__inner) {
    background-color: #fff;
    border-radius: 4px;
  }
}

.w-full {
  width: 100%;
}

.dialog-footer {
  display: flex;
  justify-content: center;
  gap: 16px;

  .el-button {
    padding: 10px 32px;
    height: 40px;
    font-size: 14px;
    border-radius: 4px;
  }

  .btn-save {
    background-color: #00a870 !important;
    border-color: #00a870 !important;
    color: #fff !important;
    &:hover {
      background-color: #008f5d !important;
      border-color: #008f5d !important;
    }
  }
}
</style>
