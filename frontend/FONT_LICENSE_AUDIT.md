# 字体版权审计报告

生成日期: 2026-02-08

## 🔴 紧急处理项

### 1. Swis721 BlkCn BT - **必须移除**
- **风险等级**: 🔴 高危
- **当前状态**: 未授权商用
- **文件位置**: `src/assets/fonts/Swis721-BlkCn-BT-Black.ttf`
- **使用位置**: `src/views/Login.vue`
- **版权方**: Bitstream Inc.
- **商用费用**: ~$35-99 USD/字重

#### 替代方案：
```css
/* 推荐替换为开源字体 */
font-family: 'Roboto Condensed', 'Noto Sans SC', sans-serif;
/* 或 */
font-family: 'Barlow Condensed', 'Noto Sans SC', sans-serif;
```

---

### 2. PangMenZhengDao (庞门正道标题体) - **需确认授权**
- **风险等级**: 🟡 中危
- **当前状态**: 免费版（疑似仅限个人使用）
- **文件位置**: `src/assets/fonts/PangMenZhengDao.ttf`  
- **使用位置**: `src/views/Licenses/LicenseSearch.vue`
- **版权方**: 庞门正道字体工作室
- **免费版限制**: 通常仅限个人非商业使用

#### 处理建议：
1. **联系作者确认商用授权**
   - 官网: https://www.zcool.com.cn/u/14777796
   - 或购买商用授权（通常 ¥199-999 RMB）

2. **或替换为开源字体**：
```css
/* 推荐中文标题字体 */
font-family: 'Noto Serif SC', 'Source Han Serif CN', serif;
/* 或使用粗体 Noto Sans */
font-family: 'Noto Sans SC', sans-serif;
font-weight: 700;
```

---

## ✅ 安全字体

### Noto Sans SC - 完全合规
- **协议**: SIL Open Font License 1.1
- **版权方**: Google Fonts
- **商用**: ✅ 完全免费
- **当前状态**: ✅ 已正确引入
- **文件**: `src/assets/fonts/noto-sans-sc-*.woff2`

---

## 📋 建议行动计划

### 立即执行（24小时内）
1. ✅ 移除 `Swis721-BlkCn-BT-Black.ttf`
2. ✅ 移除 fonts.css 中的 Swis721 声明
3. ✅ 更新 Login.vue 使用开源字体

### 短期执行（1周内）
1. 🔍 联系庞门正道确认 PangMenZhengDao 授权
2. 📝 如无法获得授权，替换为开源字体
3. 📄 添加字体授权文档到项目

### 长期优化
1. 📋 建立字体使用规范文档
2. 🔍 定期审计第三方资源版权
3. ✅ 优先使用 Google Fonts 等开源字体

---

## 🔧 技术实施清单

### 移除商业字体
```bash
# 删除文件
rm frontend/src/assets/fonts/Swis721-BlkCn-BT-Black.ttf
rm frontend/src/assets/fonts/PangMenZhengDao.ttf  # 如未获授权

# 更新 fonts.css（移除相关 @font-face 声明）
# 更新组件样式（替换 font-family）
```

### 推荐开源替代字体

#### 英文/数字字体：
- **Roboto Condensed** - 类似 Swis721 风格
- **Barlow Condensed** - 现代无衬线浓缩字体
- **Oswald** - 窄体标题字体

#### 中文字体：
- **Noto Sans SC** (已使用) - 黑体，适合正文
- **Noto Serif SC** - 宋体，适合标题
- **Source Han Sans CN** - 思源黑体（Adobe）

#### 获取方式：
```bash
# 通过 Google Fonts CDN
https://fonts.google.com/
```

---

## 📚 相关法律条款

### GPL-3.0 项目字体使用
- ✅ 可使用 SIL OFL 协议字体
- ✅ 可使用 Apache License 字体  
- ⚠️ 商业字体需额外购买授权
- ❌ 不可使用未授权商业字体

### 字体授权检查清单
- [ ] 确认字体协议类型
- [ ] 验证商用授权范围
- [ ] 保存授权证明文件
- [ ] 添加版权声明到项目

---

## 🔗 相关资源

- [Google Fonts](https://fonts.google.com/) - 免费开源字体库
- [SIL Open Font License](https://scripts.sil.org/OFL)
- [Adobe Fonts](https://fonts.adobe.com/) - 部分免费字体
- [中文字体版权指南](https://www.zcool.com.cn/article/ZMTQzNjM2.html)

---

**审计人**: GitHub Copilot  
**项目**: License Manager Frontend  
**版本**: 1.0.0  
**协议**: GPL-3.0
