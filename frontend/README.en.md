<!--
 * @Author: 13895237362 2205451508@qq.com
 * @Date: 2025-07-29 09:37:26
 * @LastEditors: 13895237362 2205451508@qq.com
 * @LastEditTime: 2025-07-29 16:55:21
 * @FilePath: /vue-demo3.0/README.md
 * @Description: This is the default setting. Please set `customMade`. Open koroFileHeader for configuration: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->

# Vue 3 + TypeScript + Vite Project Template

This project is based on Vue 3, TypeScript, and Vite, integrating mainstream front-end development tools and best practices. It is suitable for quickly building modern web applications.

## Recommended Development Environment

- **IDE**: VS Code
- **Plugins**: Volar (recommended to disable Vetur), TypeScript Vue Plugin (Volar)
- **Node Version**: Recommended 16+ (tested on 19.0.0)

## Project Features

- **Tech Stack**: Vue 3 + TypeScript + Vite
- **Routing**: Integrated vue-router, supports dynamic routing
- **State Management**: Integrated vuex
- **UI Framework**: Supports Element Plus
- **Internationalization**: Integrated vue-i18n, supports multi-language switching
- **Responsive Adaptation**: Based on px and media queries
- **Code Quality**: Integrated ESLint, Prettier
- **Build Optimization**: gzip compression, bundle analysis, code splitting
- **Clear Directory Structure**: src contains assets, components, views, router, store, utils, etc.
- **Environment Variable Management**: Supports multi-environment configuration, flexible switching, and packaging for different environments
- **Static Asset Management**: Automatically categorizes and bundles assets, adds timestamps to prevent caching
- **Unit Test Support**: Can be extended with Jest/Vitest

## Directory Structure

```
├── public/               # Static assets
├── src/
│   ├── assets/           # Images, styles, and other resources
│   ├── components/       # Common components
│   ├── views/            # Page views
│   ├── router/           # Routing configuration
│   ├── store/            # State management
│   ├── utils/            # Utility functions
│   ├── i18n/             # Internationalization configuration
│   ├── App.vue           # Root component
│   └── main.ts           # Entry file
├── .env*                 # Environment variable files
├── vite.config.ts        # Vite configuration
├── tsconfig.json         # TypeScript configuration
├── README.md             # Project documentation
```

## Quick Start

1. Install dependencies

   ```bash
   npm install
   ```

2. Start the development server

   ```bash
   npm run serve
   ```

3. Build for production

   ```bash
   npm run prod
   ```

4. Build for test environment

   ```bash
   npm run test
   ```

5. Preview production build

   ```bash
   npm run preview
   ```

## Internationalization

- Chinese and English are supported by default, and you can switch languages at the top right of the login page.
- The i18n configuration files are located in the `src/i18n/` directory.

## Other Notes

- You can extend more pages and features as needed.
- It is recommended to use VS Code for development to enjoy better type hints and code completion.
- To customize the theme or adapt to more devices, modify `src/assets/styles/variables.scss`.

## Contribution & Feedback

Feel free to submit Issues or PRs to improve features and documentation.

---

> This template is integrated with [Vue 3 + TypeScript + Vite] and mainstream plugins, suitable for enterprise-level projects or personal learning use.