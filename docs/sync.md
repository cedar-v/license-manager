## 一次性设置（首次推送）

### 第一步：克隆开源项目
```bash
# 克隆开源项目到本地
git clone https://github.com/cedar-v/license-manager.git

# 进入项目目录
cd license-manager
```

### 第二步：添加私有仓库作为新的远程仓库
```bash
# 查看当前远程仓库（应该只有 origin 指向开源项目）
git remote -v

# 添加你的私有仓库（命名为 enterprise）
git remote add enterprise https://github.com/cedar-v/license-manager-enterprise.git

# 验证（现在应该有两个远程仓库）
git remote -v
```

应该显示：
```
origin    https://github.com/cedar-v/license-manager.git (fetch)
origin    https://github.com/cedar-v/license-manager.git (push)
enterprise   https://github.com/cedar-v/license-manager-enterprise.git (fetch)
enterprise   https://github.com/cedar-v/license-manager-enterprise.git (push)
```

### 第三步：创建并切换到 license-manager 分支
```bash
# 创建 license-manager 分支
git checkout -b license-manager

# 推送到私有仓库
git push -u enterprise license-manager
```

---

## 后续更新流程

当开源项目有更新时，你可以这样同步：

```bash
# 1. 切换到 license-manager 分支
git checkout license-manager

# 2. 从开源项目拉取最新代码（origin 是开源项目）
git pull origin main

# 或者如果开源项目的默认分支是 master
# git pull origin master

# 3. 推送到你的私有仓库
git push enterprise license-manager
```

---

## 一键同步脚本（可选）

你可以创建一个脚本方便以后使用：

**sync.sh** (Mac/Linux) 或 **sync.bat** (Windows)
```bash
#!/bin/bash
git checkout license-manager
git pull origin main
git push enterprise license-manager
echo "同步完成！"
```
