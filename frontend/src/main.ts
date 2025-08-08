/*
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-07-29 09:37:26
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-08-05 15:17:46
 * @FilePath: /vue-demo3.0/src/main.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'

// 全局样式重置
import './assets/styles/global.css'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import * as ElementPlusIconsVue from '@element-plus/icons-vue'
import i18n, { getCurrentElementLocale } from './i18n'


const app = createApp(App)

// 全局注册Element Plus图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

// 正确配置Element Plus语言包
app.use(ElementPlus, {
  locale: getCurrentElementLocale()
})
app.use(i18n)
app.use(router)
app.mount('#app')
