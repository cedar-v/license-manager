<template>
  <div class="json-editor">
    <!-- 操作按钮 -->
    <div class="json-editor-actions">
      <el-button size="small" @click="addRootEntry">
        {{ t('jsonEditor.addItem') }}
      </el-button>
      <el-button size="small" text :disabled="!hasData" @click="clearAll">
        {{ t('jsonEditor.clearAll') }}
      </el-button>
      <el-button size="small" text :disabled="!hasData" @click="copyJson">
        {{ t('jsonEditor.copyJson') }}
      </el-button>
      <el-button size="small" text @click="openImportDialog">
        {{ t('jsonEditor.importJson') }}
      </el-button>
    </div>

    <!-- JSON 树形结构 -->
    <div class="json-tree">
      <JsonNode
        v-for="(node, index) in nodes"
        :key="node.id"
        :node="node"
        :depth="0"
        @update="handleNodeUpdate"
        @delete="handleNodeDelete(index)"
        @add-child="handleAddChild"
        @toggle-expand="handleToggleExpand"
      />

      <div v-if="!nodes.length" class="json-empty">
        {{ t('jsonEditor.emptyState') }}
      </div>
    </div>

    <!-- JSON 预览 -->
    <div class="json-preview">
      <div class="json-preview-label">{{ t('jsonEditor.preview') }}</div>
      <el-input
        type="textarea"
        :rows="4"
        :model-value="jsonPreview"
        readonly
      />
    </div>

    <!-- 导入弹窗 -->
    <el-dialog
      v-model="importDialogVisible"
      :title="t('jsonEditor.importTitle')"
      width="600px"
      destroy-on-close
    >
      <p class="import-dialog-tip">
        {{ t('jsonEditor.importDescription') }}
      </p>
      <el-input
        v-model="importDialogContent"
        type="textarea"
        :rows="10"
        :placeholder="t('jsonEditor.importPlaceholder')"
      />
      <p v-if="importDialogError" class="json-editor-error">
        {{ importDialogError }}
      </p>
      <template #footer>
        <el-button @click="importDialogVisible = false">
          {{ t('jsonEditor.importCancel') }}
        </el-button>
        <el-button type="primary" @click="handleImportConfirm">
          {{ t('jsonEditor.importConfirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import JsonNode from './JsonNode.vue'

export type JsonValueType = 'string' | 'number' | 'boolean' | 'object' | 'array'

export interface JsonNodeData {
  id: string
  key: string
  value: string
  type: JsonValueType
  children: JsonNodeData[]
  expanded: boolean
  keyError?: string
  valueError?: string
}

const props = defineProps<{
  modelValue?: Record<string, any> | string | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: Record<string, any> | string | null): void
}>()

const { t } = useI18n()

// 响应式数据
const nodes = ref<JsonNodeData[]>([])
const importDialogVisible = ref(false)
const importDialogContent = ref('')
const importDialogError = ref('')

// 计算属性
const hasData = computed(() => nodes.value.length > 0)

const jsonPreview = computed(() => {
  const data = buildJsonData(nodes.value)
  return data ? JSON.stringify(data, null, 2) : ''
})

// 生成唯一ID
const createNodeId = () => `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`

// 创建空节点
const createEmptyNode = (type: JsonValueType = 'string'): JsonNodeData => ({
  id: createNodeId(),
  key: '',
  value: type === 'object' || type === 'array' ? '' : '',
  type,
  children: [],
  expanded: true,
  keyError: '',
  valueError: ''
})

// 构建 JSON 数据
const buildJsonData = (nodeList: JsonNodeData[]): Record<string, any> | null => {
  if (!nodeList.length) return null

  const result: Record<string, any> = {}
  for (const node of nodeList) {
    const trimmedKey = node.key.trim()
    if (!trimmedKey) continue

    let value: any

    switch (node.type) {
      case 'number':
        const num = Number(node.value)
        value = Number.isNaN(num) ? node.value : num
        break
      case 'boolean':
        value = node.value === 'true'
        break
      case 'object':
        value = buildJsonData(node.children)
        break
      case 'array':
        value = node.children.map(child => {
          if (child.type === 'object') {
            return buildJsonData(child.children)
          }
          if (child.type === 'number') {
            return Number(child.value)
          }
          if (child.type === 'boolean') {
            return child.value === 'true'
          }
          return child.value
        })
        break
      default:
        value = node.value
    }

    result[trimmedKey] = value
  }

  return Object.keys(result).length ? result : null
}

// 解析 JSON 为节点
const parseJsonToNodes = (data: unknown): JsonNodeData[] => {
  if (!data) return []
  if (typeof data === 'string') {
    try {
      data = JSON.parse(data)
    } catch {
      return []
    }
  }

  if (typeof data !== 'object' || data === null) return []

  if (Array.isArray(data)) {
    return data.map((item, index) => {
      const isObj = item !== null && typeof item === 'object'
      const isArr = Array.isArray(item)
      const itemType: JsonValueType = item === null ? 'string' : (isArr ? 'array' : (isObj ? 'object' : typeof item as JsonValueType))
      return {
        id: createNodeId(),
        key: String(index),
        value: isObj || isArr ? '' : String(item),
        type: itemType,
        children: (isObj || isArr) ? parseJsonToNodes(item) : [],
        expanded: true,
        keyError: '',
        valueError: ''
      }
    })
  }

  return Object.entries(data as Record<string, any>).map(([key, value]) => {
    let type: JsonValueType = 'string'
    let children: JsonNodeData[] = []

    if (value === null) {
      type = 'string'
    } else if (Array.isArray(value)) {
      type = 'array'
      children = value.map((item, index) => {
        const isObj = item !== null && typeof item === 'object'
        const isArr = Array.isArray(item)
        const itemType: JsonValueType = item === null ? 'string' : (isArr ? 'array' : (isObj ? 'object' : typeof item as JsonValueType))
        return {
          id: createNodeId(),
          key: String(index),
          value: isObj || isArr ? '' : String(item),
          type: itemType,
          children: (isObj || isArr) ? parseJsonToNodes(item) : [],
          expanded: true,
          keyError: '',
          valueError: ''
        }
      })
    } else if (typeof value === 'object') {
      type = 'object'
      children = parseJsonToNodes(value)
    } else if (typeof value === 'number') {
      type = 'number'
    } else if (typeof value === 'boolean') {
      type = 'boolean'
    }

    return {
      id: createNodeId(),
      key,
      value: type === 'object' || type === 'array' ? '' : String(value),
      type,
      children,
      expanded: true,
      keyError: '',
      valueError: ''
    }
  })
}

// 初始化节点
const initializeNodes = (value: unknown) => {
  nodes.value = parseJsonToNodes(value)
}

// 添加根节点
const addRootEntry = () => {
  nodes.value.push(createEmptyNode('string'))
  emitUpdate()
}

// 清空
const clearAll = () => {
  nodes.value = []
  emitUpdate()
}

// 复制 JSON
const fallbackCopy = (value: string) => {
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

const copyJson = async () => {
  if (!jsonPreview.value) {
    ElMessage.info(t('jsonEditor.copyEmpty'))
    return
  }
  try {
    if (navigator?.clipboard?.writeText) {
      await navigator.clipboard.writeText(jsonPreview.value)
    } else {
      fallbackCopy(jsonPreview.value)
    }
    ElMessage.success(t('jsonEditor.copySuccess'))
  } catch (error) {
    console.error('复制JSON失败:', error)
    ElMessage.error(t('jsonEditor.copyError'))
  }
}

// 导入弹窗
const openImportDialog = () => {
  importDialogContent.value = jsonPreview.value
  importDialogError.value = ''
  importDialogVisible.value = true
}

const handleImportConfirm = () => {
  importDialogError.value = ''
  const content = importDialogContent.value.trim()

  if (!content) {
    nodes.value = []
    importDialogVisible.value = false
    ElMessage.success(t('jsonEditor.importSuccess'))
    emitUpdate()
    return
  }

  try {
    const parsed = JSON.parse(content)
    nodes.value = parseJsonToNodes(parsed)
    importDialogVisible.value = false
    ElMessage.success(t('jsonEditor.importSuccess'))
    emitUpdate()
  } catch (error: any) {
    importDialogError.value = error?.message || t('jsonEditor.importFailed')
  }
}

// 节点更新（子节点输入 key/value 等时同步到 v-model）
const handleNodeUpdate = () => {
  emitUpdate()
}

// 节点删除
const handleNodeDelete = (index: number) => {
  nodes.value.splice(index, 1)
  emitUpdate()
}

// 添加子节点
const handleAddChild = (parentId: string) => {
  const findAndAdd = (nodeList: JsonNodeData[]): boolean => {
    for (const node of nodeList) {
      if (node.id === parentId) {
        if (node.type === 'array') {
          node.children.push(createEmptyNode('string'))
        } else if (node.type === 'object') {
          node.children.push(createEmptyNode('string'))
        }
        return true
      }
      if (node.children.length && findAndAdd(node.children)) {
        return true
      }
    }
    return false
  }
  findAndAdd(nodes.value)
  emitUpdate()
}

// 切换展开
const handleToggleExpand = (nodeId: string) => {
  const findAndToggle = (nodeList: JsonNodeData[]): boolean => {
    for (const node of nodeList) {
      if (node.id === nodeId) {
        node.expanded = !node.expanded
        return true
      }
      if (node.children.length && findAndToggle(node.children)) {
        return true
      }
    }
    return false
  }
  findAndToggle(nodes.value)
}

// 验证所有节点
const validateAll = (): boolean => {
  let isValid = true
  const existingKeys = new Map<string, number>()

  const validate = (nodeList: JsonNodeData[], parentPath: string = '') => {
    for (const node of nodeList) {
      node.keyError = ''
      node.valueError = ''

      const trimmedKey = node.key.trim()
      const fullPath = parentPath ? `${parentPath}.${trimmedKey}` : trimmedKey

      // Key 验证
      if (!trimmedKey) {
        node.keyError = t('jsonEditor.keyRequired')
        isValid = false
      } else {
        const keyCount = existingKeys.get(fullPath) || 0
        if (keyCount > 0) {
          node.keyError = t('jsonEditor.keyDuplicate')
          isValid = false
        }
        existingKeys.set(fullPath, keyCount + 1)
      }

      // Value 验证
      if (node.type === 'number') {
        if (node.value && Number.isNaN(Number(node.value))) {
          node.valueError = t('jsonEditor.numberInvalid')
          isValid = false
        }
      }

      // 递归验证子节点
      if (node.children.length) {
        validate(node.children, fullPath)
      }
    }
  }

  validate(nodes.value)
  return isValid
}

// 暴露方法
defineExpose({
  validate: validateAll,
  getValue: () => buildJsonData(nodes.value),
  setValue: initializeNodes
})

// 内部 emit 时跳过 modelValue 回灌，否则 buildJsonData 会丢掉「空 key」占位行，
// watch 再 initializeNodes 会把刚添加的行立刻清掉（表现为「添加字段」无效）。
const skipModelSync = ref(false)

watch(
  () => props.modelValue,
  (newValue) => {
    if (skipModelSync.value) return
    initializeNodes(newValue)
  },
  { immediate: true }
)

const emitUpdate = () => {
  nextTick(() => {
    skipModelSync.value = true
    const newValue = buildJsonData(nodes.value)
    emit('update:modelValue', newValue)
    nextTick(() => {
      skipModelSync.value = false
    })
  })
}
</script>

<style scoped>
.json-editor {
  width: 100%;
}

.json-editor-actions {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.json-tree {
  border: 1px solid var(--app-border-color);
  border-radius: 4px;
  padding: 12px;
  margin-bottom: 16px;
  min-height: 100px;
  background: var(--app-content-bg);
}

.json-empty {
  text-align: center;
  color: var(--app-text-secondary);
  padding: 20px;
  font-size: 13px;
}

.json-preview {
  margin-top: 16px;
}

.json-preview-label {
  font-size: 12px;
  color: var(--app-text-secondary);
  margin-bottom: 8px;
}

.import-dialog-tip {
  font-size: 13px;
  color: var(--app-text-secondary);
  margin-bottom: 12px;
}

.json-editor-error {
  color: #f56c6c;
  font-size: 12px;
  margin-top: 8px;
}
</style>