/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-07-29 09:37:26
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-07-31 15:04:59
 * @FilePath: /vue-demo3.0/src/main.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// 全局样式重置
import './assets/styles/global.css'

import i18n from './i18n' // 国际化语言配置
import { loadLocaleMessages,SupportedLocale } from './i18n'
// 加载默认语言包
loadLocaleMessages(i18n.global.locale.value as SupportedLocale)

const app = createApp(App)
app.use(router)
app.use(i18n)
app.mount('#app')