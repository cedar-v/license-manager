<template>
  <div class="license-form-container">
    <!-- 顶部横向区域 - 与表单平级 -->
    <div class="top-section">
      <!-- 面包屑导航 -->
      <div class="breadcrumb">
        <span class="breadcrumb-item">{{
          t('pages.licenses.form.breadcrumb.licenseManagement')
        }}</span>
        <span class="breadcrumb-separator">></span>
        <span class="breadcrumb-current">{{
          isEdit
            ? t('pages.licenses.form.breadcrumb.editLicense')
            : t('pages.licenses.form.breadcrumb.createLicense')
        }}</span>
      </div>

      <!-- 操作按钮 -->
      <div class="form-actions">
        <el-button @click="handleCancel">{{ t('pages.licenses.form.actions.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">{{
          isEdit ? t('pages.licenses.form.actions.update') : t('pages.licenses.form.actions.create')
        }}</el-button>
      </div>
    </div>

    <!-- 主表单区域 -->
    <div>
      <el-form :model="formData" :rules="formRules" ref="formRef" label-position="top">
        <!-- 基本信息 -->
        <div class="license-form">
          <h3 class="section-title">{{ t('pages.licenses.form.sections.basicInfo') }}</h3>

          <!-- 第一行：关联客户名称，客户ID -->
          <div class="fields-row">
            <el-form-item
              :label="t('pages.licenses.form.fields.customerName')"
              prop="customer_id"
              required
              class="field-item"
            >
              <el-select
                v-model="formData.customer_id"
                :placeholder="t('pages.licenses.form.placeholders.selectCustomer')"
                filterable
                remote
                :remote-method="searchCustomers"
                :loading="customerLoading"
                style="width: 100%"
              >
                <el-option
                  v-for="customer in customerOptions"
                  :key="customer.id"
                  :label="customer.customer_name"
                  :value="customer.id"
                />
              </el-select>
            </el-form-item>

            <el-form-item
              :label="t('pages.licenses.form.fields.customerId')"
              prop="customer_code"
              required
              class="field-item"
            >
              <el-input
                v-model="formData.customer_code"
                :placeholder="t('pages.licenses.form.placeholders.autoGenerateOrManual')"
                disabled
              />
            </el-form-item>
          </div>

          <!-- 第二行：备注单独一行 -->
          <div class="fields-row">
            <el-form-item
              :label="t('pages.licenses.form.fields.description')"
              prop="description"
              required
              class="field-item-full"
            >
              <el-input
                v-model="formData.description"
                type="textarea"
                :rows="4"
                :placeholder="t('pages.licenses.form.placeholders.enterDescription')"
              />
            </el-form-item>
          </div>
        </div>

        <!-- 授权配置 -->
        <div class="license-form">
          <h3 class="section-title">{{ t('pages.licenses.form.sections.licenseConfig') }}</h3>

          <!-- 第三行：授权期限类型，起止时间 -->
          <div class="fields-row">
            <el-form-item
              :label="t('pages.licenses.form.fields.validityPeriod')"
              prop="validity_type"
              required
              class="field-item"
            >
              <el-select
                v-model="formData.validity_type"
                :placeholder="t('pages.licenses.form.placeholders.selectValidityType')"
                style="width: 100%"
                @change="handleValidityTypeChange"
              >
                <el-option
                  :label="t('pages.licenses.form.validityTypes.permanent')"
                  value="permanent"
                />
                <el-option
                  :label="t('pages.licenses.form.validityTypes.limited')"
                  value="limited"
                />
              </el-select>
            </el-form-item>

            <el-form-item
              v-if="formData.validity_type === 'limited'"
              :label="t('pages.licenses.form.fields.dateRange')"
              prop="date_range"
              required
              class="field-item"
            >
              <el-date-picker
                v-model="dateRange"
                type="daterange"
                :range-separator="t('chart.licenseTrend.datePicker.rangeSeparator')"
                :start-placeholder="t('pages.licenses.form.placeholders.startDate')"
                :end-placeholder="t('pages.licenses.form.placeholders.endDate')"
                format="YYYY-MM-DD"
                value-format="YYYY-MM-DD"
                style="width: 100%"
                :disabled-date="disabledDate"
                @change="handleDateRangeChange"
              />
            </el-form-item>

            <el-form-item
              v-if="formData.validity_type === 'permanent'"
              :label="t('pages.licenses.form.fields.validityDays')"
              class="field-item"
            >
              <el-input
                :value="t('pages.licenses.form.validityTypes.permanent') + '（365000天）'"
                readonly
                style="width: 100%"
              />
            </el-form-item>
          </div>

          <!-- 第四行：最大激活设备数，部署类型，加密类型 -->
          <div class="fields-row">
            <el-form-item
              :label="t('pages.licenses.form.fields.maxActivations')"
              prop="max_activations"
              required
              class="field-item"
            >
              <el-input-number
                v-model="formData.max_activations"
                :min="1"
                :placeholder="t('pages.licenses.form.placeholders.enterMaxActivations')"
                style="width: 100%"
              />
            </el-form-item>

            <el-form-item
              :label="t('pages.licenses.form.fields.deploymentType')"
              prop="deployment_type"
              required
              class="field-item"
            >
              <el-select
                v-model="formData.deployment_type"
                :placeholder="t('pages.licenses.form.placeholders.selectDeploymentType')"
                style="width: 100%"
                @change="handleDeploymentTypeChange"
              >
                <el-option
                  v-for="option in deploymentTypeOptions"
                  :key="option.key"
                  :label="option.display"
                  :value="option.key"
                />
              </el-select>
            </el-form-item>

            <el-form-item
              :label="t('pages.licenses.form.fields.encryptionType')"
              prop="encryption_type"
              required
              class="field-item"
            >
              <el-select
                v-model="formData.encryption_type"
                :placeholder="t('pages.licenses.form.placeholders.selectEncryptionType')"
                style="width: 100%"
              >
                <el-option
                  v-for="option in encryptionTypeOptions"
                  :key="option.key"
                  :label="option.display"
                  :value="option.key"
                />
              </el-select>
            </el-form-item>
          </div>
        </div>

        <!-- 功能配置 -->
        <div class="license-form">
          <div class="key-value-section-header">
            <h3 class="section-title">{{ t('pages.licenses.form.sections.featureConfig') }}</h3>
            <div class="key-value-section-actions">
              <el-button size="small" @click="addFeatureConfigEntry">
                {{ t('pages.licenses.form.keyValue.addItem') }}
              </el-button>
              <el-button
                size="small"
                text
                :disabled="!featureConfigEntries.length"
                @click="clearFeatureConfigEntries"
              >
                {{ t('pages.licenses.form.keyValue.clearAll') }}
              </el-button>
              <el-button
                size="small"
                text
                :disabled="!featureConfigPreview"
                @click="copyJson(featureConfigPreview)"
              >
                {{ t('pages.licenses.form.keyValue.copyJson') }}
              </el-button>
              <el-button size="small" text @click="openImportDialog('feature')">
                {{ t('pages.licenses.form.keyValue.importJson') }}
              </el-button>
            </div>
          </div>

          <div class="key-value-list">
            <div v-for="(item, index) in featureConfigEntries" :key="item.id" class="key-value-row">
              <div class="key-value-field">
                <label>{{ t('pages.licenses.form.keyValue.keyLabel') }}</label>
                <el-input
                  v-model="item.key"
                  :placeholder="t('pages.licenses.form.keyValue.keyPlaceholder')"
                  @input="item.keyError = ''"
                />
                <p v-if="item.keyError" class="key-value-error">{{ item.keyError }}</p>
              </div>
              <div class="key-value-field type-field">
                <label>{{ t('pages.licenses.form.keyValue.typeLabel') }}</label>
                <el-select v-model="item.type" @change="handleEntryTypeChange(item)">
                  <el-option
                    v-for="type in typeOptions"
                    :key="type"
                    :label="t(`pages.licenses.form.keyValue.typeOptions.${type}`)"
                    :value="type"
                  />
                </el-select>
              </div>
              <div class="key-value-field">
                <label>{{ t('pages.licenses.form.keyValue.valueLabel') }}</label>
                <template v-if="item.type === 'bool'">
                  <el-select v-model="item.value" @change="item.valueError = ''">
                    <el-option
                      v-for="option in boolOptions"
                      :key="option.value"
                      :label="option.label"
                      :value="option.value"
                    />
                  </el-select>
                </template>
                <el-input
                  v-else
                  v-model="item.value"
                  :placeholder="t('pages.licenses.form.keyValue.valuePlaceholder')"
                  @input="item.valueError = ''"
                  :inputmode="item.type === 'number' ? 'decimal' : 'text'"
                />
                <p v-if="item.valueError" class="key-value-error">{{ item.valueError }}</p>
              </div>
              <el-button link type="danger" @click="removeFeatureConfigEntry(index)">
                {{ t('pages.licenses.form.keyValue.remove') }}
              </el-button>
            </div>

            <div v-if="!featureConfigEntries.length" class="key-value-empty">
              {{ t('pages.licenses.form.keyValue.emptyState') }}
            </div>
          </div>

          <div class="key-value-preview">
            <div class="key-value-preview-label">
              {{ t('pages.licenses.form.keyValue.preview') }}
            </div>
            <el-input type="textarea" :rows="3" :model-value="featureConfigPreview" readonly />
          </div>
        </div>

        <!-- 限制设置 -->
        <div class="license-form">
          <div class="key-value-section-header">
            <h3 class="section-title">{{ t('pages.licenses.form.sections.usageLimits') }}</h3>
            <div class="key-value-section-actions">
              <el-button size="small" @click="addLimitEntry">
                {{ t('pages.licenses.form.keyValue.addItem') }}
              </el-button>
              <el-button size="small" text :disabled="!limitEntries.length" @click="clearLimitEntries">
                {{ t('pages.licenses.form.keyValue.clearAll') }}
              </el-button>
              <el-button
                size="small"
                text
                :disabled="!usageJsonPreview"
                @click="copyJson(usageJsonPreview)"
              >
                {{ t('pages.licenses.form.keyValue.copyJson') }}
              </el-button>
              <el-button size="small" text @click="openImportDialog('limit')">
                {{ t('pages.licenses.form.keyValue.importJson') }}
              </el-button>
            </div>
          </div>

          <div class="key-value-list">
            <div v-for="(item, index) in limitEntries" :key="item.id" class="key-value-row">
              <div class="key-value-field">
                <label>{{ t('pages.licenses.form.keyValue.keyLabel') }}</label>
                <el-input
                  v-model="item.key"
                  :placeholder="t('pages.licenses.form.keyValue.keyPlaceholder')"
                  @input="item.keyError = ''"
                />
                <p v-if="item.keyError" class="key-value-error">{{ item.keyError }}</p>
              </div>
              <div class="key-value-field type-field">
                <label>{{ t('pages.licenses.form.keyValue.typeLabel') }}</label>
                <el-select v-model="item.type" @change="handleEntryTypeChange(item)">
                  <el-option
                    v-for="type in typeOptions"
                    :key="type"
                    :label="t(`pages.licenses.form.keyValue.typeOptions.${type}`)"
                    :value="type"
                  />
                </el-select>
              </div>
              <div class="key-value-field">
                <label>{{ t('pages.licenses.form.keyValue.valueLabel') }}</label>
                <template v-if="item.type === 'bool'">
                  <el-select v-model="item.value" @change="item.valueError = ''">
                    <el-option
                      v-for="option in boolOptions"
                      :key="option.value"
                      :label="option.label"
                      :value="option.value"
                    />
                  </el-select>
                </template>
                <el-input
                  v-else
                  v-model="item.value"
                  :placeholder="t('pages.licenses.form.keyValue.valuePlaceholder')"
                  @input="item.valueError = ''"
                  :inputmode="item.type === 'number' ? 'decimal' : 'text'"
                />
                <p v-if="item.valueError" class="key-value-error">{{ item.valueError }}</p>
              </div>
              <el-button link type="danger" @click="removeLimitEntry(index)">
                {{ t('pages.licenses.form.keyValue.remove') }}
              </el-button>
            </div>

            <div v-if="!limitEntries.length" class="key-value-empty">
              {{ t('pages.licenses.form.keyValue.emptyState') }}
            </div>
          </div>

          <div class="key-value-preview">
            <div class="key-value-preview-label">
              {{ t('pages.licenses.form.keyValue.preview') }}
            </div>
            <el-input type="textarea" :rows="3" :model-value="usageJsonPreview" readonly />
          </div>
        </div>

        <!-- 自定义参数 -->
        <div class="license-form">
          <div class="key-value-section-header">
            <h3 class="section-title">{{ t('pages.licenses.form.sections.customParameters') }}</h3>
            <div class="key-value-section-actions">
              <el-button size="small" @click="addCustomEntry">
                {{ t('pages.licenses.form.keyValue.addItem') }}
              </el-button>
              <el-button size="small" text :disabled="!customEntries.length" @click="clearCustomEntries">
                {{ t('pages.licenses.form.keyValue.clearAll') }}
              </el-button>
              <el-button
                size="small"
                text
                :disabled="!customJsonPreview"
                @click="copyJson(customJsonPreview)"
              >
                {{ t('pages.licenses.form.keyValue.copyJson') }}
              </el-button>
              <el-button size="small" text @click="openImportDialog('custom')">
                {{ t('pages.licenses.form.keyValue.importJson') }}
              </el-button>
            </div>
          </div>

          <div class="key-value-list">
            <div v-for="(item, index) in customEntries" :key="item.id" class="key-value-row">
              <div class="key-value-field">
                <label>{{ t('pages.licenses.form.keyValue.keyLabel') }}</label>
                <el-input
                  v-model="item.key"
                  :placeholder="t('pages.licenses.form.keyValue.keyPlaceholder')"
                  @input="item.keyError = ''"
                />
                <p v-if="item.keyError" class="key-value-error">{{ item.keyError }}</p>
              </div>
              <div class="key-value-field type-field">
                <label>{{ t('pages.licenses.form.keyValue.typeLabel') }}</label>
                <el-select v-model="item.type" @change="handleEntryTypeChange(item)">
                  <el-option
                    v-for="type in typeOptions"
                    :key="type"
                    :label="t(`pages.licenses.form.keyValue.typeOptions.${type}`)"
                    :value="type"
                  />
                </el-select>
              </div>
              <div class="key-value-field">
                <label>{{ t('pages.licenses.form.keyValue.valueLabel') }}</label>
                <template v-if="item.type === 'bool'">
                  <el-select v-model="item.value" @change="item.valueError = ''">
                    <el-option
                      v-for="option in boolOptions"
                      :key="option.value"
                      :label="option.label"
                      :value="option.value"
                    />
                  </el-select>
                </template>
                <el-input
                  v-else
                  v-model="item.value"
                  :placeholder="t('pages.licenses.form.keyValue.valuePlaceholder')"
                  @input="item.valueError = ''"
                  :inputmode="item.type === 'number' ? 'decimal' : 'text'"
                />
                <p v-if="item.valueError" class="key-value-error">{{ item.valueError }}</p>
              </div>
              <el-button link type="danger" @click="removeCustomEntry(index)">
                {{ t('pages.licenses.form.keyValue.remove') }}
              </el-button>
            </div>

            <div v-if="!customEntries.length" class="key-value-empty">
              {{ t('pages.licenses.form.keyValue.emptyState') }}
            </div>
          </div>

          <div class="key-value-preview">
            <div class="key-value-preview-label">
              {{ t('pages.licenses.form.keyValue.preview') }}
            </div>
            <el-input type="textarea" :rows="3" :model-value="customJsonPreview" readonly />
          </div>
        </div>
      </el-form>
    </div>
  </div>

  <el-dialog
    v-model="importDialogVisible"
    :title="importDialogTitle"
    width="520px"
    destroy-on-close
  >
    <p class="import-dialog-tip">
      {{ t('pages.licenses.form.keyValue.importDescription') }}
    </p>
    <el-input
      v-model="importDialogContent"
      type="textarea"
      :rows="8"
      :placeholder="t('pages.licenses.form.keyValue.importPlaceholder')"
    />
    <p v-if="importDialogError" class="key-value-error import-error">{{ importDialogError }}</p>
    <template #footer>
      <el-button @click="importDialogVisible = false">
        {{ t('pages.licenses.form.keyValue.importCancel') }}
      </el-button>
      <el-button type="primary" @click="handleImportConfirm">
        {{ t('pages.licenses.form.keyValue.importConfirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue'
import type { Ref } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  createLicense,
  updateLicense,
  getLicenseDetail,
  type AuthorizationCodeCreateRequest,
  type LicenseUpdateRequest
} from '@/api/license'
import { getCustomers, type Customer } from '@/api/customer'
import { getEnumOptions, type RawEnumItem } from '@/api/enum'

type KeyValueType = 'string' | 'number' | 'bool'

interface KeyValueItem {
  id: string
  key: string
  value: string
  type: KeyValueType
  keyError?: string
  valueError?: string
}

const KEY_PATTERN = /^[a-z0-9-]{1,32}$/

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const boolOptions = computed(() => [
  { label: t('pages.licenses.form.keyValue.booleanTrue'), value: 'true' },
  { label: t('pages.licenses.form.keyValue.booleanFalse'), value: 'false' }
])

// 表单引用
const formRef = ref<FormInstance>()

// 响应式数据
const submitting = ref(false)
const customerLoading = ref(false)
const customerOptions = ref<Customer[]>([])

const featureConfigEntries = ref<KeyValueItem[]>([])
const limitEntries = ref<KeyValueItem[]>([])
const customEntries = ref<KeyValueItem[]>([])
const importDialogVisible = ref(false)
const importDialogTarget = ref<'feature' | 'limit' | 'custom'>('feature')
const importDialogTitle = computed(() => {
  if (importDialogTarget.value === 'limit') {
    return t('pages.licenses.form.keyValue.importTitleLimit')
  }
  if (importDialogTarget.value === 'custom') {
    return t('pages.licenses.form.keyValue.importTitleCustom')
  }
  return t('pages.licenses.form.keyValue.importTitleFeature')
})
const importDialogContent = ref('')
const importDialogError = ref('')

// 枚举选项
const deploymentTypeOptions = ref<RawEnumItem[]>([])
const encryptionTypeOptions = ref<RawEnumItem[]>([])

const typeOptions: KeyValueType[] = ['string', 'number', 'bool']
// 表单数据
const formData = reactive<
  AuthorizationCodeCreateRequest & {
    customer_code: string
    validity_type: 'permanent' | 'limited'
  }
>({
  customer_id: '',
  customer_code: '',
  validity_type: 'limited',
  description: '',
  validity_days: 365,
  deployment_type: 'standalone',
  encryption_type: 'standard',
  max_activations: 1,
  feature_config: {},
  usage_limits: '',
  custom_parameters: ''
})

// 独立的响应式变量
const dateRange = ref<[string, string] | null>(null)

const featureConfigPreview = computed(() => buildJsonString(featureConfigEntries.value))
const usageJsonPreview = computed(() => buildJsonString(limitEntries.value))
const customJsonPreview = computed(() => buildJsonString(customEntries.value))

// 计算属性
const isEdit = computed(() => {
  return route.name === 'licenses-edit' && route.params.id
})

// 禁止选择今天之前的日期
const disabledDate = (time: Date) => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return time.getTime() < today.getTime()
}

// 表单验证规则
const formRules: FormRules = {
  customer_id: [
    {
      required: true,
      message: t('pages.licenses.form.validation.customerRequired'),
      trigger: 'change'
    },
    {
      validator: (_rule: any, value: any, callback: any) => {
        if (value && customerOptions.value.length > 0) {
          const selectedCustomer = customerOptions.value.find(c => c.id === value)
          if (!selectedCustomer) {
            callback(new Error(t('pages.licenses.form.validation.customerInvalid')))
          } else {
            callback()
          }
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  description: [
    {
      required: true,
      message: t('pages.licenses.form.validation.descriptionRequired'),
      trigger: 'blur'
    },
    {
      min: 1,
      max: 500,
      message: t('pages.licenses.form.validation.descriptionLength'),
      trigger: 'blur'
    }
  ],
  validity_type: [
    {
      required: true,
      message: t('pages.licenses.form.validation.validityTypeRequired'),
      trigger: 'change'
    }
  ],
  date_range: [
    {
      required: true,
      message: t('pages.licenses.form.validation.dateRangeRequired'),
      trigger: 'change',
      validator: (_rule: any, _value: any, callback: any) => {
        if (
          formData.validity_type === 'limited' &&
          (!dateRange.value || dateRange.value.length !== 2)
        ) {
          callback(new Error(t('pages.licenses.form.validation.dateRangeRequired')))
        } else {
          callback()
        }
      }
    },
    {
      validator: (_rule: any, _value: any, callback: any) => {
        if (
          formData.validity_type === 'limited' &&
          dateRange.value &&
          dateRange.value.length === 2
        ) {
          const startDate = new Date(dateRange.value[0])
          const endDate = new Date(dateRange.value[1])
          const today = new Date()
          today.setHours(0, 0, 0, 0)

          if (startDate < today) {
            callback(new Error(t('pages.licenses.form.validation.startDateNotPast')))
          } else if (endDate <= startDate) {
            callback(new Error(t('pages.licenses.form.validation.endDateAfterStart')))
          } else {
            // 检查日期范围是否超过合理范围（比如10年）
            const maxDays = 3650 // 10年
            const diffTime = endDate.getTime() - startDate.getTime()
            const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
            if (diffDays > maxDays) {
              callback(new Error(t('pages.licenses.form.validation.validityPeriodTooLong')))
            } else {
              callback()
            }
          }
        } else {
          callback()
        }
      },
      trigger: 'change'
    }
  ],
  deployment_type: [
    {
      required: true,
      message: t('pages.licenses.form.validation.deploymentTypeRequired'),
      trigger: 'change'
    }
  ],
  encryption_type: [
    {
      required: true,
      message: t('pages.licenses.form.validation.encryptionTypeRequired'),
      trigger: 'change'
    }
  ],
  max_activations: [
    {
      required: true,
      message: t('pages.licenses.form.validation.maxActivationsRequired'),
      trigger: 'blur'
    },
    {
      type: 'number',
      min: 1,
      max: 999999,
      message: t('pages.licenses.form.validation.maxActivationsRange'),
      trigger: 'blur'
    },
    {
      validator: (_rule: any, value: any, callback: any) => {
        if (value && !Number.isInteger(Number(value))) {
          callback(new Error(t('pages.licenses.form.validation.maxActivationsInteger')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const createEntryId = () => `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`

const createEmptyEntry = (): KeyValueItem => ({
  id: createEntryId(),
  key: '',
  value: '',
  type: 'string',
  keyError: '',
  valueError: ''
})

const convertValueByType = (value: string, type: KeyValueType) => {
  if (type === 'number') {
    const num = Number(value)
    if (Number.isNaN(num)) {
      throw new Error('invalid number')
    }
    return num
  }
  if (type === 'bool') {
    return value === 'true'
  }
  return value
}

const buildJsonRecord = (entries: KeyValueItem[]) => {
  if (!entries.length) {
    return null
  }
  const result: Record<string, string | number | boolean> = {}
  entries.forEach(item => {
    const trimmedKey = item.key.trim()
    if (trimmedKey) {
      try {
        result[trimmedKey] = convertValueByType(item.value, item.type)
      } catch {
        // 如果转换失败，则忽略该条目，等待校验提示
      }
    }
  })
  return Object.keys(result).length ? result : null
}

const buildJsonString = (entries: KeyValueItem[]) => {
  const record = buildJsonRecord(entries)
  return record ? JSON.stringify(record) : ''
}

const addFeatureConfigEntry = () => {
  featureConfigEntries.value.push(createEmptyEntry())
}

const addLimitEntry = () => {
  limitEntries.value.push(createEmptyEntry())
}

const addCustomEntry = () => {
  customEntries.value.push(createEmptyEntry())
}

const removeFeatureConfigEntry = (index: number) => {
  featureConfigEntries.value.splice(index, 1)
}

const removeLimitEntry = (index: number) => {
  limitEntries.value.splice(index, 1)
}

const removeCustomEntry = (index: number) => {
  customEntries.value.splice(index, 1)
}

const clearFeatureConfigEntries = () => {
  featureConfigEntries.value = []
}

const clearLimitEntries = () => {
  limitEntries.value = []
}

const clearCustomEntries = () => {
  customEntries.value = []
}

const handleEntryTypeChange = (entry: KeyValueItem) => {
  if (entry.type === 'bool') {
    entry.value = entry.value === 'false' ? 'false' : 'true'
  } else if (entry.type === 'number') {
    if (entry.value && Number.isNaN(Number(entry.value))) {
      entry.value = ''
    }
  }
  entry.valueError = ''
}

const fallbackCopy = (value: string) => {
  if (typeof document === 'undefined') {
    throw new Error('document is not available')
  }
  const textarea = document.createElement('textarea')
  textarea.value = value
  textarea.style.position = 'fixed'
  textarea.style.opacity = '0'
  document.body.appendChild(textarea)
  textarea.focus()
  textarea.select()
  document.execCommand('copy')
  document.body.removeChild(textarea)
}

const copyJson = async (value: string) => {
  if (!value) {
    ElMessage.info(t('pages.licenses.form.keyValue.copyEmpty'))
    return
  }
  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(value)
    } else {
      fallbackCopy(value)
    }
    ElMessage.success(t('pages.licenses.form.keyValue.copySuccess'))
  } catch (error) {
    console.error('复制JSON失败:', error)
    ElMessage.error(t('pages.licenses.form.keyValue.copyError'))
  }
}

const validateKeyValueEntries = (entries: KeyValueItem[]) => {
  let isValid = true
  const existingKeys = new Set<string>()

  entries.forEach(item => {
    item.keyError = ''
    item.valueError = ''

    const trimmedKey = item.key.trim()
    if (!trimmedKey) {
      item.keyError = t('pages.licenses.form.keyValue.keyRequired')
      isValid = false
    } else if (!KEY_PATTERN.test(trimmedKey)) {
      item.keyError = t('pages.licenses.form.keyValue.keyFormat')
      isValid = false
    } else if (existingKeys.has(trimmedKey)) {
      item.keyError = t('pages.licenses.form.keyValue.keyDuplicate')
      isValid = false
    } else {
      existingKeys.add(trimmedKey)
    }

    if (item.value.trim() === '') {
      item.valueError = t('pages.licenses.form.keyValue.valueRequired')
      isValid = false
    } else if (item.type === 'number') {
      const num = Number(item.value)
      if (Number.isNaN(num)) {
        item.valueError = t('pages.licenses.form.keyValue.numberRequired')
        isValid = false
      }
    } else if (item.type === 'bool') {
      if (item.value !== 'true' && item.value !== 'false') {
        item.valueError = t('pages.licenses.form.keyValue.boolRequired')
        isValid = false
      }
    }
  })

  return isValid
}

const openImportDialog = (target: 'feature' | 'limit' | 'custom') => {
  importDialogTarget.value = target
  importDialogContent.value =
    target === 'feature'
      ? featureConfigPreview.value
      : target === 'limit'
        ? usageJsonPreview.value
        : customJsonPreview.value
  importDialogError.value = ''
  importDialogVisible.value = true
}

const handleImportConfirm = () => {
  importDialogError.value = ''
  try {
    const content = importDialogContent.value.trim()
    if (!content) {
      if (importDialogTarget.value === 'feature') {
        featureConfigEntries.value = []
      } else if (importDialogTarget.value === 'limit') {
        limitEntries.value = []
      } else {
        customEntries.value = []
      }
    } else {
      const entries = parseJsonToEntries(content, { strict: true })
      if (!entries.length) {
        throw new Error(t('pages.licenses.form.keyValue.importEmpty'))
      }
      if (importDialogTarget.value === 'feature') {
        featureConfigEntries.value = entries
      } else if (importDialogTarget.value === 'limit') {
        limitEntries.value = entries
      } else {
        customEntries.value = entries
      }
    }
    importDialogVisible.value = false
    ElMessage.success(t('pages.licenses.form.keyValue.importSuccess'))
  } catch (error: any) {
    importDialogError.value =
      error?.message || t('pages.licenses.form.keyValue.importFailed')
  }
}

const mapObjectToEntries = (data: Record<string, any>): KeyValueItem[] => {
  return Object.entries(data).map(([key, value]) => {
    let type: KeyValueType = 'string'
    let normalizedValue = ''
    if (typeof value === 'number') {
      type = 'number'
      normalizedValue = String(value)
    } else if (typeof value === 'boolean') {
      type = 'bool'
      normalizedValue = value ? 'true' : 'false'
    } else {
      normalizedValue = value === undefined || value === null ? '' : String(value)
    }

    return {
      id: createEntryId(),
      key,
      value: normalizedValue,
      type,
      keyError: '',
      valueError: ''
    }
  })
}

const parseJsonToEntries = (raw: unknown, options?: { strict?: boolean }): KeyValueItem[] => {
  if (!raw) return []

  let parsed: unknown = raw
  if (typeof raw === 'string') {
    if (!raw.trim()) return []
    try {
      parsed = JSON.parse(raw)
    } catch (error) {
      if (options?.strict) {
        throw error
      }
      console.warn('解析JSON失败:', error)
      return []
    }
  }

  if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
    if (options?.strict) {
      throw new Error('JSON payload must be an object with key-value pairs')
    }
    return []
  }

  return mapObjectToEntries(parsed as Record<string, any>)
}

const hydrateEntriesFromField = (raw: unknown, target: Ref<KeyValueItem[]>) => {
  target.value = parseJsonToEntries(raw)
}

// 方法
const loadEnumOptions = async () => {
  try {
    // 加载部署类型选项
    const deploymentResponse = await getEnumOptions('deployment_type')
    if (deploymentResponse.code === '000000') {
      deploymentTypeOptions.value = deploymentResponse.data.items
    } else {
      throw new Error(deploymentResponse.message || '加载部署类型选项失败')
    }

    // 加载加密类型选项
    const encryptionResponse = await getEnumOptions('encryption_type')
    if (encryptionResponse.code === '000000') {
      encryptionTypeOptions.value = encryptionResponse.data.items
    } else {
      throw new Error(encryptionResponse.message || '加载加密类型选项失败')
    }
  } catch (error: any) {
    console.error('加载枚举选项失败:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('pages.licenses.form.messages.loadEnumError') + ': ' + errorMessage)
    } else {
      ElMessage.error(t('pages.licenses.form.messages.loadEnumErrorRetry'))
    }
  }
}

// 搜索或加载所有客户
const searchCustomers = async (query: string) => {
  try {
    customerLoading.value = true
    const response = await getCustomers({
      page: 1,
      page_size: 100,
      customer_name: query || undefined,
      status: 'active'
    })
    customerOptions.value = response.data.list
  } catch (error: any) {
    console.error('搜索客户失败:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('pages.licenses.form.messages.searchCustomerError') + ': ' + errorMessage)
    } else {
      ElMessage.error(t('pages.licenses.form.messages.searchCustomerErrorRetry'))
    }
    customerOptions.value = []
  } finally {
    customerLoading.value = false
  }
}

// 表单初始化
const loadCustomerInfo = () => {
  // 从路由参数获取客户信息
  const customerId = route.query.customerId as string
  const customerName = route.query.customerName as string

  if (customerId && customerName) {
    formData.customer_id = customerId
    customerOptions.value = [{
      id: customerId,
      customer_name: customerName,
      customer_code: customerId,
      customer_type: '',
      customer_type_display: '',
      contact_person: '',
      email: '',
      customer_level: '',
      customer_level_display: '',
      status: 'active',
      status_display: '',
      created_at: '',
      updated_at: '',
      created_by: '',
      updated_by: ''
    } as Customer]
  }

}

// 处理授权期限类型变化
const handleValidityTypeChange = (type: 'permanent' | 'limited') => {
  if (type === 'permanent') {
    formData.validity_days = 365000
    dateRange.value = null
  } else {
    // 有限期限时，如果已有日期范围，重新计算天数
    if (dateRange.value && dateRange.value.length === 2) {
      const startDate = new Date(dateRange.value[0])
      const endDate = new Date(dateRange.value[1])
      const diffTime = endDate.getTime() - startDate.getTime()
      formData.validity_days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    }
  }
}

// 处理部署类型变化
const handleDeploymentTypeChange = () => {
  // 部署类型变化时的处理逻辑（如需要可以在这里添加）
}

// 处理日期范围变化
const handleDateRangeChange = (dates: [string, string] | null) => {
  if (dates && dates.length === 2) {
    const startDate = new Date(dates[0])
    const endDate = new Date(dates[1])
    const diffTime = endDate.getTime() - startDate.getTime()
    formData.validity_days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  }
}

// 根据id获取授权详情
const loadLicenseDetail = async () => {
  if (!isEdit.value) return

  try {
    const id = route.params.id as string
    const response = await getLicenseDetail(id)

    if (response.code === '000000' && response.data) {
      const data = response.data

      // 填充表单数据
      formData.customer_id = data.customer_id || ''
      formData.description = data.description || ''
      formData.deployment_type = data.deployment_type || 'standalone'
      formData.encryption_type = data.encryption_type || 'standard'
      formData.max_activations = data.max_activations || 1

      // 设置表单字段
      formData.customer_code = (data as any).customer_code || data.customer_id || ''

      // 计算有效期天数和设置日期范围
      if (data.start_date && data.end_date) {
        const startDate = new Date(data.start_date)
        const endDate = new Date(data.end_date)
        const diffTime = endDate.getTime() - startDate.getTime()
        const days = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
        formData.validity_days = days

        // 根据天数判断是永久还是有限
        if (days >= 365000) {
          formData.validity_type = 'permanent'
          dateRange.value = null
        } else {
          formData.validity_type = 'limited'
          dateRange.value = [data.start_date, data.end_date]
        }
      } else {
        // 如果没有日期信息，根据validity_days判断
        const validityDays = (data as any).validity_days
        if (validityDays && validityDays >= 365000) {
          formData.validity_type = 'permanent'
          formData.validity_days = 365000
        } else {
          formData.validity_type = 'limited'
          formData.validity_days = validityDays || 365
        }
      }

      // 设置客户选项
      if (data.customer_name) {
        customerOptions.value = [
          {
            id: data.customer_id,
            customer_name: data.customer_name,
            customer_code: (data as any).customer_code || '',
            customer_type: '',
            customer_type_display: '',
            contact_person: '',
            email: '',
            customer_level: '',
            customer_level_display: '',
            status: '',
            status_display: '',
            created_at: '',
            updated_at: '',
            created_by: '',
            updated_by: ''
          } as Customer
        ]
      }

      hydrateEntriesFromField(data.feature_config, featureConfigEntries)
      hydrateEntriesFromField(data.usage_limits, limitEntries)
      hydrateEntriesFromField(data.custom_parameters, customEntries)
    } else {
      throw new Error(response.message || t('pages.licenses.form.messages.loadDetailError'))
    }
  } catch (error: any) {
    console.error('加载授权详情失败:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('pages.licenses.form.messages.loadDetailError') + ': ' + errorMessage)
    } else {
      ElMessage.error(t('pages.licenses.form.messages.loadDetailErrorRetry'))
    }
  }
}

// 提交（新增/更新）授权数据
const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    // 执行表单验证
    await formRef.value.validate()

    const featureConfigValid = validateKeyValueEntries(featureConfigEntries.value)
    const limitValid = validateKeyValueEntries(limitEntries.value)
    const customValid = validateKeyValueEntries(customEntries.value)
    if (!featureConfigValid || !limitValid || !customValid) {
      ElMessage.error(t('pages.licenses.form.keyValue.validationFailed'))
      return
    }

    formData.feature_config = buildJsonRecord(featureConfigEntries.value) || {}
    formData.usage_limits = usageJsonPreview.value
    formData.custom_parameters = customJsonPreview.value

    submitting.value = true

    if (isEdit.value) {
      // 更新授权
      const id = route.params.id as string
      const updateData: LicenseUpdateRequest = {
        description: formData.description,
        validity_days: formData.validity_days,
        deployment_type: formData.deployment_type,
        encryption_type: formData.encryption_type,
        max_activations: formData.max_activations,
        feature_config: formData.feature_config,
        usage_limits: formData.usage_limits,
        custom_parameters: formData.custom_parameters,
        change_type: 'other',
        reason: t('pages.licenses.form.messages.updateSuccess')
      }

      const response = await updateLicense(id, updateData)
      if (response.code === '000000') {
        ElMessage.success(response.message || t('pages.licenses.form.messages.updateSuccess'))
        // 返回列表页面
        router.back()
      } else {
        throw new Error(response.message || t('pages.licenses.form.messages.submitError'))
      }
    } else {
      // 创建授权
      const response = await createLicense(formData)
      if (response.code === '000000') {
        ElMessage.success(response.message || t('pages.licenses.form.messages.createSuccess'))
        // 返回列表页面
        router.back()
      } else {
        throw new Error(response.message || t('pages.licenses.form.messages.submitError'))
      }
    }
  } catch (error: any) {
    console.error('提交失败:', error)
    let errorMessage = error.backendMessage || error.response?.data?.message || error.message
    if (errorMessage) {
      ElMessage.error(t('pages.licenses.form.messages.submitError') + ': ' + errorMessage)
    } else {
      ElMessage.error(t('pages.licenses.form.messages.submitErrorRetry'))
    }
  } finally {
    submitting.value = false
  }
}

// 取消函数
const handleCancel = () => {
  ElMessageBox.confirm(
    t('pages.licenses.form.messages.cancelConfirm'),
    t('pages.licenses.form.messages.cancelTitle'),
    {
      confirmButtonText: t('pages.licenses.form.messages.cancelConfirmButton'),
      cancelButtonText: t('pages.licenses.form.messages.cancelCancelButton'),
      type: 'warning'
    }
  )
    .then(() => {
      router.back()
    })
    .catch(() => {
      // 用户选择继续编辑
    })
}

// 监听客户选择变化，自动填充客户ID
watch(
  () => formData.customer_id,
  newCustomerId => {
    if (newCustomerId) {
      const selectedCustomer = customerOptions.value.find(c => c.id === newCustomerId)
      if (selectedCustomer) {
        formData.customer_code = selectedCustomer.customer_code || selectedCustomer.id
      } else {
        // 如果没有找到客户，清空客户代码
        formData.customer_code = ''
      }
    } else {
      formData.customer_code = ''
    }
  }
)

watch(
  limitEntries,
  () => {
    formData.usage_limits = usageJsonPreview.value
  },
  { deep: true }
)

watch(
  featureConfigEntries,
  () => {
    formData.feature_config = buildJsonRecord(featureConfigEntries.value) || {}
  },
  { deep: true }
)

watch(
  customEntries,
  () => {
    formData.custom_parameters = customJsonPreview.value
  },
  { deep: true }
)

// 生命周期
onMounted(async () => {
  await loadEnumOptions()
  loadCustomerInfo()
  // 如果没有从路由获取到客户信息，则加载所有客户
  if (!route.query.customerId) {
    await searchCustomers('')
  }
  await loadLicenseDetail()
})
</script>

<style scoped>
/* 主容器 */
.license-form-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: var(--app-bg-color);
  padding: 15px;
}

/* 顶部横向区域 */
.top-section {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: transparent;
  margin-bottom: 10px;
}

.breadcrumb {
  font-size: 14px;
  color: var(--app-text-secondary);
}

.breadcrumb-item {
  color: var(--app-text-secondary);
}

.breadcrumb-separator {
  margin: 0 8px;
  color: var(--app-text-secondary);
}

.breadcrumb-current {
  color: var(--app-text-primary);
}

.form-actions {
  display: flex;
  gap: 20px;
}

/* 主表单区域 */
.license-form {
  padding: 20px;
  background: var(--app-content-bg);
  border-radius: 4px;
  width: 100%;
  margin-bottom: 20px;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--app-text-primary);
  margin: 0 0 20px 0;
}

/* 字段行布局 */
.fields-row {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 20px;
}

/* 基础字段：按照Figma设计的325px宽度 */
.fields-row .field-item {
  flex: 0 0 325px; /* 固定宽度325px */
  min-width: 280px; /* 移动端最小宽度 */
  margin-bottom: 0;
}

/* 备注字段：占满整行 */
.fields-row .field-item-full {
  flex: 1 1 100%;
  min-width: 100%;
  margin-bottom: 0;
}

.key-value-section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.key-value-section-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.key-value-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.key-value-row {
  display: flex;
  gap: 12px;
  align-items: flex-start;
  flex-wrap: wrap;
}

.key-value-field {
  flex: 1;
}

.key-value-field.type-field {
  flex: 0 0 150px;
}

.key-value-field label {
  display: block;
  font-size: 12px;
  color: var(--app-text-secondary);
  margin-bottom: 6px;
}

.key-value-error {
  color: #f56c6c;
  font-size: 12px;
  margin-top: 6px;
}

.key-value-empty {
  border: 1px dashed var(--app-border-color);
  border-radius: 4px;
  padding: 16px;
  text-align: center;
  color: var(--app-text-secondary);
  font-size: 13px;
}

.key-value-preview {
  margin-top: 16px;
}

.key-value-preview-label {
  font-size: 12px;
  color: var(--app-text-secondary);
  margin-bottom: 8px;
}

.key-value-row > .el-button {
  margin-top: 22px;
}

.import-dialog-tip {
  font-size: 13px;
  color: var(--app-text-secondary);
  margin-bottom: 12px;
}

.import-error {
  margin-top: 8px;
}

/* Element Plus 样式覆盖 */
:deep(.el-form-item__label) {
  font-weight: 400;
  color: var(--app-text-primary);
  font-size: 14px;
  line-height: 1.5714285714285714em;
  margin-bottom: 8px;
}

:deep(.el-form-item.is-required .el-form-item__label::before) {
  content: '*';
  color: #f56c6c;
  margin-right: 4px;
}

:deep(.el-input__wrapper) {
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  height: 32px;
}

:deep(.el-input__inner) {
  color: var(--app-text-primary);
  font-size: 14px;
  line-height: 1.5em;
}

:deep(.el-input__inner::placeholder) {
  color: var(--app-text-secondary);
  font-size: 14px;
}

:deep(.el-select .el-input__wrapper) {
  border: 1px solid var(--app-border-color);
  height: 32px;
}

:deep(.el-input-number) {
  width: 100%;
}

:deep(.el-input-number .el-input__wrapper) {
  border: 1px solid var(--app-border-color);
  height: 32px;
}

:deep(.el-date-editor) {
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  height: 32px;
}

:deep(.el-date-editor .el-input__wrapper) {
  border: none;
  height: 30px;
}

:deep(.el-button--primary) {
  background-color: var(--el-color-primary-dark-2);
  border-color: var(--el-color-primary-dark-2);
  font-size: 14px;
  padding: 6px 16px;
}

:deep(.el-button--primary:hover) {
  background-color: var(--el-color-primary);
  border-color: var(--el-color-primary);
}

:deep(.el-button) {
  font-size: 14px;
  padding: 6px 16px;
}

/* 响应式设计 */
@media (max-width: 1440px) {
  .fields-row {
    flex-wrap: wrap;
  }

  .fields-row .field-item {
    min-width: 200px;
  }
}

@media (max-width: 768px) {
  .top-section {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }

  .fields-row {
    flex-direction: column;
    gap: 16px;
  }

  .fields-row .field-item,
  .fields-row .field-item-full {
    width: 100%;
    flex: none;
  }
}
</style>
