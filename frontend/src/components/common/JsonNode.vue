<template>
  <div class="json-node" :class="{ 'is-child': depth > 0 }">
    <!-- 节点内容 -->
    <div class="json-node-content">
      <!-- 展开/折叠按钮 (对象和数组) -->
      <el-button
        v-if="isContainer"
        link
        class="expand-btn"
        @click="$emit('toggle-expand', node.id)"
      >
        <el-icon>
          <ArrowRight v-if="!node.expanded" />
          <ArrowDown v-else />
        </el-icon>
      </el-button>
      <span v-else class="expand-placeholder"></span>

      <!-- Key 输入 -->
      <div class="node-field key-field">
        <el-input
          v-model="node.key"
          :placeholder="t('jsonEditor.keyPlaceholder')"
          size="small"
          @input="node.keyError = ''"
        />
        <p v-if="node.keyError" class="node-error">{{ node.keyError }}</p>
      </div>

      <!-- 类型选择 -->
      <div class="node-field type-field">
        <el-select
          v-model="node.type"
          size="small"
          @change="handleTypeChange"
        >
          <el-option
            v-for="type in typeOptions"
            :key="type"
            :label="t(`jsonEditor.types.${type}`)"
            :value="type"
          />
        </el-select>
      </div>

      <!-- Value 输入 -->
      <div class="node-field value-field">
        <!-- 布尔类型 -->
        <el-select
          v-if="node.type === 'boolean'"
          v-model="node.value"
          size="small"
        >
          <el-option :label="'true'" value="true" />
          <el-option :label="'false'" value="false" />
        </el-select>

        <!-- 对象/数组类型 - 显示子节点 -->
        <template v-else-if="isContainer">
          <span class="container-label">
            <template v-if="node.type === 'array'">[{{ node.children.length }}]</template>
            <template v-else>{ {{ node.children.length }} }</template>
          </span>
        </template>

        <!-- 数字类型 -->
        <el-input
          v-else-if="node.type === 'number'"
          v-model="node.value"
          :placeholder="t('jsonEditor.valuePlaceholder')"
          size="small"
          type="number"
          @input="node.valueError = ''"
        />

        <!-- 字符串类型 -->
        <el-input
          v-else
          v-model="node.value"
          :placeholder="t('jsonEditor.valuePlaceholder')"
          size="small"
          @input="node.valueError = ''"
        />
        <p v-if="node.valueError" class="node-error">{{ node.valueError }}</p>
      </div>

      <!-- 操作按钮 -->
      <div class="node-actions">
        <el-button
          v-if="isContainer"
          link
          type="primary"
          size="small"
          @click="$emit('add-child', node.id)"
        >
          {{ t('jsonEditor.addChild') }}
        </el-button>
        <el-button
          v-if="depth > 0 || node.type === 'array'"
          link
          type="danger"
          size="small"
          @click="$emit('delete')"
        >
          {{ t('jsonEditor.remove') }}
        </el-button>
      </div>
    </div>

    <!-- 子节点容器 -->
    <div v-if="isContainer && node.expanded" class="json-node-children">
      <JsonNode
        v-for="(child, index) in node.children"
        :key="child.id"
        :node="child"
        :depth="depth + 1"
        @update="$emit('update')"
        @delete="removeChild(index)"
        @add-child="(id: string) => $emit('add-child', id)"
        @toggle-expand="(id: string) => $emit('toggle-expand', id)"
      />

      <div v-if="!node.children.length" class="child-empty">
        {{ node.type === 'array' ? t('jsonEditor.arrayEmpty') : t('jsonEditor.objectEmpty') }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick } from 'vue'
import { ArrowRight, ArrowDown } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import type { JsonNodeData, JsonValueType } from './JsonEditor.vue'

const props = defineProps<{
  node: JsonNodeData
  depth: number
}>()

const emit = defineEmits<{
  (e: 'update'): void
  (e: 'delete'): void
  (e: 'add-child', id: string): void
  (e: 'toggle-expand', id: string): void
}>()

const { t } = useI18n()

const typeOptions: JsonValueType[] = ['string', 'number', 'boolean', 'object', 'array']

const isContainer = computed(() => props.node.type === 'object' || props.node.type === 'array')

const handleTypeChange = (newType: JsonValueType) => {
  // 清空 children 如果切换到非容器类型
  if (newType !== 'object' && newType !== 'array') {
    props.node.children = []
  }
  // 清空 value 如果切换到容器类型
  if (newType === 'object' || newType === 'array') {
    props.node.value = ''
  }
  // 初始化空 children 如果切换到容器类型
  if ((newType === 'object' || newType === 'array') && !props.node.children.length) {
    props.node.children = []
  }
  nextTick(() => {
    emit('update')
  })
}

const removeChild = (index: number) => {
  nextTick(() => {
    props.node.children.splice(index, 1)
    emit('update')
  })
}
</script>

<style scoped>
.json-node {
  margin-bottom: 8px;
}

.json-node.is-child {
  margin-left: 24px;
  padding-left: 12px;
  border-left: 1px dashed var(--app-border-color);
}

.json-node-content {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  flex-wrap: wrap;
}

.expand-btn {
  padding: 2px;
  min-width: 20px;
}

.expand-placeholder {
  width: 20px;
  display: inline-block;
}

.node-field {
  display: flex;
  flex-direction: column;
}

.key-field {
  flex: 1;
  min-width: 120px;
}

.type-field {
  flex: 0 0 100px;
}

.value-field {
  flex: 1;
  min-width: 120px;
}

.container-label {
  display: inline-block;
  padding: 4px 12px;
  background: var(--app-bg-color);
  border-radius: 4px;
  font-size: 12px;
  color: var(--app-text-secondary);
}

.node-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
}

.node-error {
  color: #f56c6c;
  font-size: 11px;
  margin-top: 2px;
}

.json-node-children {
  margin-top: 8px;
}

.child-empty {
  text-align: center;
  color: var(--app-text-secondary);
  font-size: 12px;
  padding: 8px;
}
</style>