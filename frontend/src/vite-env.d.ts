/// <reference types="vite/client" />
declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
  }
  // 解决项目中拿不到 import.meta.env.VITE_API_BASE_URL的值，
  // 你也可以在项目中使用ts断言 来获取值 如： import.meta.env.VITE_API_BASE_UR as string