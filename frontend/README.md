<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-07-29 09:37:26
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-07-29 16:52:58
 * @FilePath: /vue-demo3.0/README.md
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->

# Vue 3 + TypeScript + Vite 项目模板

本项目基于 Vue 3、TypeScript 和 Vite，集成了主流前端开发工具和最佳实践，适合快速搭建现代 Web 应用。

## 推荐开发环境

- **IDE**：VS Code
- **插件**：Volar（建议禁用 Vetur）、TypeScript Vue Plugin (Volar)
- **Node 版本**：建议 16+ (n本人19.0.0)

## 项目特性

- **技术栈**：Vue 3 + TypeScript + Vite
- **路由管理**：集成 vue-router，支持动态路由
- **状态管理**：集成 vuex
- **UI 框架**：支持 Element Plus
- **国际化**：集成 vue-i18n，支持多语言切换
- **响应式适配**：基于 px 和媒体查询 
- **代码规范**：集成 ESLint、Prettier
- **打包优化**：gzip 压缩、打包分析、按需分包
- **目录结构清晰**：src 下分 assets、components、views、router、store、utils 等
- **环境变量管理**：支持多环境配置，灵活切换,并打包成不同的环境
- **静态资源管理**：自动分类打包，带时间戳防缓存
- **单元测试支持**：可扩展 Jest/Vitest

## 目录结构

```
├── public/               # 静态资源
├── src/
│   ├── assets/           # 图片、样式等资源
│   ├── components/       # 公共组件
│   ├── views/            # 页面视图
│   ├── router/           # 路由配置
│   ├── store/            # 状态管理
│   ├── utils/            # 工具函数
│   ├── i18n/             # 国际化配置
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── .env*                 # 环境变量文件
├── vite.config.ts        # Vite 配置
├── tsconfig.json         # TypeScript 配置
├── README.md             # 项目说明
```

## 快速开始

1. 安装依赖

   ```bash
   npm install
   ```

2. 启动开发服务器

   ```bash
   npm run serve
   ```

3. 生产环境打包构建

   ```bash
   npm run prod
   ```
4. 测试环境打包构建

   ```bash
   npm run test
   ```

5. 预览生产环境

   ```bash
   npm run preview
   ```

## 多语言国际化

- 默认支持中文和英文，可在登录页右上角切换语言。
- 国际化配置文件位于 `src/i18n/` 目录。


## 其他说明

- 可根据实际需求扩展更多页面和功能。
- 推荐使用 VS Code 进行开发，享受更好的类型提示和代码补全。
- 如需自定义主题或适配更多设备，可修改 `src/assets/styles/variables.scss`。

## 贡献与反馈

欢迎提交 Issue 或 PR，完善功能和文档。

---

> 本模板由 [Vue 3 + TypeScript + Vite] 结合主流插件集成，适合企业级项目或个人学习使用