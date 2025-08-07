<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-08-05 14:46:50
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-05 15:31:22
 * @FilePath: /frontend/src/App.vue
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->
<template>
  <el-config-provider :locale="elementLocale">
    <router-view />
  </el-config-provider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import en from 'element-plus/es/locale/lang/en'
import ja from 'element-plus/es/locale/lang/ja'

const { locale } = useI18n()

 const localeMap = {
    zh: zhCn,
    en: en,
    ja: ja
  }
const elementLocale = computed(() => {
  return localeMap[locale.value as keyof typeof localeMap] || zhCn
})
</script>

<style>
#app {
  height: 100%;
}
</style>
