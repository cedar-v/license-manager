# 技术契约与架构设计书：登录页国际化补全

**版本**：v1.0  
**日期**：2026-03-13  
**作者**：system-architect  
**状态**：待 task-coder 实施

---

## 一、问题诊断

`src/views/Login.vue` 中存在 **8 处**硬编码英文字符串，未通过 `t()` 调用 i18n 系统，导致切换语言后这些文案不随语言切换。

### 硬编码清单（按 DOM 位置）

| 位置 | 当前硬编码值 | 目标 key |
|------|-------------|---------|
| `.hero-title` | `Secure License<br/>Management` | `login.heroTitleLine1` / `login.heroTitleLine2` |
| `.hero-subtitle` | `Streamline your software licensing...` | `login.heroSubtitle` |
| stat-label "Active Licenses" | `Active Licenses` | `login.stats.activeLicenses` |
| stat-label "Uptime" | `Uptime` | `login.stats.uptime` |
| stat-label "Support" | `Support` | `login.stats.support` |
| `.title` (card) | `Welcome back` | `login.welcomeTitle` |
| `.subtitle` (card) | `Sign in to your account to continue` | `login.welcomeSubtitle` |
| `.login-footer p` | `© 2025 Cedar-V. All rights reserved.` | `login.copyright` |

> **保持原样（不需国际化）**：品牌名 `Cedar-V`、数字统计值 `10K+` / `99.9%` / `24/7`、语言选项标签（自语言名称自描述惯例）。

---

## 二、数据模型层（Schema）— i18n 键值设计

### 新增键（追加至三端 `login` 对象末尾，位于 `error` 块之后）

```jsonc
// en.json
"heroTitleLine1": "Secure License",
"heroTitleLine2": "Management",
"heroSubtitle": "Streamline your software licensing workflow with our enterprise-grade platform.",
"stats": {
  "activeLicenses": "Active Licenses",
  "uptime": "Uptime",
  "support": "Support"
},
"welcomeTitle": "Welcome back",
"welcomeSubtitle": "Sign in to your account to continue",
"copyright": "© 2025 Cedar-V. All rights reserved."

// zh.json
"heroTitleLine1": "安全许可证",
"heroTitleLine2": "管理平台",
"heroSubtitle": "使用我们的企业级平台，简化您的软件许可证管理工作流程。",
"stats": {
  "activeLicenses": "活跃许可证",
  "uptime": "可用率",
  "support": "全天支持"
},
"welcomeTitle": "欢迎回来",
"welcomeSubtitle": "登录您的账户以继续",
"copyright": "© 2025 Cedar-V. 保留所有权利。"

// ja.json
"heroTitleLine1": "セキュアな",
"heroTitleLine2": "ライセンス管理",
"heroSubtitle": "エンタープライズグレードのプラットフォームで、ソフトウェアライセンス管理を効率化。",
"stats": {
  "activeLicenses": "アクティブライセンス",
  "uptime": "稼働率",
  "support": "サポート"
},
"welcomeTitle": "おかえりなさい",
"welcomeSubtitle": "アカウントにサインインして続けてください",
"copyright": "© 2025 Cedar-V. All rights reserved."
```

---

## 三、前端交互层（View）— 修改意图

### 目标文件

| 文件 | 操作 |
|------|------|
| `src/i18n/locales/en.json` | 追加 8 个新键至 `login` 对象 |
| `src/i18n/locales/zh.json` | 追加 8 个新键至 `login` 对象 |
| `src/i18n/locales/ja.json` | 追加 8 个新键至 `login` 对象 |
| `src/views/Login.vue` | 模板层 8 处硬编码 → `{{ t('login.*') }}` |

### Login.vue 模板变更对照

```html
<!-- BEFORE -->
<h1 class="hero-title">Secure License<br/>Management</h1>

<!-- AFTER -->
<h1 class="hero-title">{{ t('login.heroTitleLine1') }}<br/>{{ t('login.heroTitleLine2') }}</h1>
```

```html
<!-- BEFORE -->
<p class="hero-subtitle">Streamline your software licensing workflow with our enterprise-grade platform.</p>

<!-- AFTER -->
<p class="hero-subtitle">{{ t('login.heroSubtitle') }}</p>
```

```html
<!-- BEFORE (3 stat labels) -->
<span class="stat-label">Active Licenses</span>
<span class="stat-label">Uptime</span>
<span class="stat-label">Support</span>

<!-- AFTER -->
<span class="stat-label">{{ t('login.stats.activeLicenses') }}</span>
<span class="stat-label">{{ t('login.stats.uptime') }}</span>
<span class="stat-label">{{ t('login.stats.support') }}</span>
```

```html
<!-- BEFORE -->
<h2 class="title">Welcome back</h2>
<p class="subtitle">Sign in to your account to continue</p>

<!-- AFTER -->
<h2 class="title">{{ t('login.welcomeTitle') }}</h2>
<p class="subtitle">{{ t('login.welcomeSubtitle') }}</p>
```

```html
<!-- BEFORE -->
<p>© 2025 Cedar-V. All rights reserved.</p>

<!-- AFTER -->
<p>{{ t('login.copyright') }}</p>
```
