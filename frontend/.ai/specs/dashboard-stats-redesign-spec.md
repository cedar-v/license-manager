# 技术契约与架构设计书：Dashboard 统计卡片重设计

**版本**：v1.0  
**日期**：2026-03-13  
**作者**：system-architect  
**状态**：待评审

---

## 一、现状诊断

### 当前 API 响应结构

```json
{
  "total_auth_codes": 12,
  "active_licenses": 2,
  "expiring_soon": 0,
  "abnormal_alerts": 2,
  "growth_rate": {
    "auth_codes": 100.00,
    "licenses": 0.00
  }
}
```

### 现存问题清单

| # | 问题 | 严重性 |
|---|------|--------|
| 1 | `growth_rate` 语义不清——是"同比"上月还是"环比"？无时间标注 | 高 |
| 2 | 两个增长率字段（auth_codes / licenses）对用户的业务价值极低，不知道要对比什么 | 高 |
| 3 | 前端 `statsData` 中的 `key` 映射与 i18n label 不一致（注释与实际label错位） | 中 |
| 4 | `expiring_soon` 统计的是"30天内到期的授权码"，但卡片标签写的是"即将过期"，口径模糊 | 中 |
| 5 | 没有"今日新增"这类高频关注的运营指标 | 中 |
| 6 | 增长率算法缺陷：上月总数为0时直接返回0而非100% | 低 |

---

## 二、竞品分析（License Management SaaS 参考）

### Keygen.sh / LicenseSpring / Cryptlex 共性 Dashboard 指标

```
┌─────────────────────────────────────────────────────────┐
│  高频关注（运营视角）                                     │
│  • 今日新增授权码 / 本月新增授权码（绝对数 + 环比昨日/上月）│
│  • 活跃设备数（在线 License）                            │
│  • 即将到期（7天 / 30天）→ 需要跟进续费                   │
│  • 异常告警（心跳超时 / 被撤销）                          │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│  中频关注（管理视角）                                     │
│  • 授权码总量（存量指标）                                 │
│  • 活跃许可证总量（存量指标）                             │
│  • 月度增长率（MoM Growth Rate，附说明文字"vs 上月"）     │
└─────────────────────────────────────────────────────────┘
```

### 关键竞品设计原则

1. **增长率必须标注对比周期**：统一用「vs. 上月同期」，数字旁展示箭头 ↑ ↓ 和颜色区分正负
2. **增长率不放卡片主数值**：作为副文本展示（小字 + 颜色），主数值展示绝对数
3. **即将到期细分**：7天 为紧急（红），30天 为预警（橙）
4. **今日新增**：运营最关注的高频指标，竞品均在首位展示

---

## 三、重设计方案

### 3.1 新指标体系（6张卡片）

| 位置 | 指标名 | 计算口径 | 副文本 |
|------|--------|---------|--------|
| 卡片1 | 授权码总量 | `COUNT(authorization_codes)` | ↑ 本月新增 N 个 |
| 卡片2 | 活跃许可证 | `COUNT(licenses WHERE status='active')` | ↑ 较上月 +X% |
| 卡片3 | 今日新增授权 | `COUNT(licenses WHERE DATE(created_at)=TODAY)` | 昨日 N 个 |
| 卡片4 | 7天内到期 | `COUNT(auth_codes WHERE end_date BETWEEN now AND now+7d AND NOT locked)` | 紧急跟进 |
| 卡片5 | 30天内到期 | `COUNT(auth_codes WHERE end_date BETWEEN now AND now+30d AND NOT locked)` | 含7天内 |
| 卡片6 | 异常告警 | `COUNT(licenses WHERE status='active' AND heartbeat_timeout)` | 心跳超时 |

> 删除：`auth_codes 增长率`、`licenses 增长率` 作为独立卡片（改为副文本方式展示）

### 3.2 增长率展示方式（副文本，不占卡片）

```
┌──────────────────────────┐
│ 🔑 授权码总量             │
│    12                   │
│  ↑ 本月新增 2 个          │  ← 副文本：绝对增量
└──────────────────────────┘

┌──────────────────────────┐
│ ✅ 活跃许可证             │
│    98                   │
│  ↑ +12.5%  vs. 上月      │  ← 副文本：百分比增长率 + 标注
└──────────────────────────┘
```

### 3.3 新 API Response Schema

```go
// StatsOverviewResponse 重设计
type StatsOverviewResponse struct {
    // 存量指标
    TotalAuthCodes  int64 `json:"total_auth_codes"`   // 授权码总量
    ActiveLicenses  int64 `json:"active_licenses"`    // 活跃许可证
    
    // 流量指标
    TodayNewLicenses    int64 `json:"today_new_licenses"`     // 今日新增许可证
    YesterdayNewLicenses int64 `json:"yesterday_new_licenses"` // 昨日新增（用于对比）
    MonthNewAuthCodes   int64 `json:"month_new_auth_codes"`   // 本月新增授权码
    
    // 风险指标
    ExpiringIn7Days  int64 `json:"expiring_in_7days"`   // 7天内到期
    ExpiringIn30Days int64 `json:"expiring_in_30days"`  // 30天内到期
    AbnormalAlerts   int64 `json:"abnormal_alerts"`     // 异常告警
    
    // 增长率（作为副文本，不占独立卡片）
    GrowthRate struct {
        AuthCodesMoM  float64 `json:"auth_codes_mom"`  // 授权码月环比
        LicensesMoM   float64 `json:"licenses_mom"`    // 许可证月环比
    } `json:"growth_rate"`
}
```

---

## 四、影响面分析（全栈）

### 后端改动文件

| 文件 | 操作 | 说明 |
|------|------|------|
| `backend/internal/models/license.go` | 修改 | 更新 `StatsOverviewResponse` struct |
| `backend/internal/service/license_service.go` | 修改 | 重写 `GetStatsOverview()` 实现，增加今日/昨日新增查询 |
| `backend/internal/api/handlers/license_handler.go` | 无需修改 | Handler 层直接透传 |

### 前端改动文件

| 文件 | 操作 | 说明 |
|------|------|------|
| `frontend/src/views/Dashboard.vue` | 修改 | 重新定义 `statsData` 6张卡片，增加副文本字段渲染 |
| `frontend/src/i18n/locales/zh.json` | 修改 | 更新 `dashboard.stats.*` key |
| `frontend/src/i18n/locales/en.json` | 修改 | 同上 |
| `frontend/src/i18n/locales/ja.json` | 修改 | 同上 |
| `frontend/src/api/dashboard.ts` | 修改 | 更新 `getOverviewStats` 返回类型定义 |

---

## 五、实施优先级建议

**P0（立即）**：
- 重命名"即将过期"为"30天内到期"，语义清晰化
- 增长率卡片改为副文本方式（不移除数据，只改展示层）
- 增长率数字旁加"vs. 上月"标注

**P1（本迭代）**：
- 新增"今日新增"卡片（需要后端加一个 COUNT 查询）
- "即将过期"拆分7天/30天两档

**P2（下迭代）**：
- 卡片支持点击跳转到过滤后的列表页
