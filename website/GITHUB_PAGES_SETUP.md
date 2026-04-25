# GitHub Pages 自动部署配置指南

## 📋 前置条件

确保你的项目已经推送到 GitHub 仓库。

## 🔧 配置步骤

### 步骤 1：启用 GitHub Pages

1. **访问仓库设置**
   - 打开你的 GitHub 仓库
   - 点击顶部的 **Settings** 标签

2. **进入 Pages 设置**
   - 在左侧边栏找到并点击 **Pages**
   - 或者访问：`https://github.com/your-org/carrot-agent/settings/pages`

3. **配置部署源**
   - 在 "Build and deployment" 部分
   - **Source** 选择：**GitHub Actions**
   - ⚠️ 不要选择 "Deploy from a branch"

4. **保存设置**
   - 配置会自动保存

### 步骤 2：验证工作流文件

确认 `.github/workflows/deploy-website.yml` 已存在并包含以下关键配置：

```yaml
name: Deploy Website

on:
  push:
    branches: [main]        # 监听 main 分支
    paths:
      - 'website/**'         # 仅当 website 目录变化时触发
  workflow_dispatch:         # 允许手动触发

permissions:
  contents: read
  pages: write              # 需要写入 Pages 的权限
  id-token: write
```

✅ 工作流文件已经配置完成！

### 步骤 3：提交并推送代码

```bash
# 添加所有更改
git add .

# 提交
git commit -m "Add website with auto-deployment to GitHub Pages"

# 推送到 GitHub
git push origin main
```

### 步骤 4：监控部署过程

1. **查看 Actions**
   - 访问：`https://github.com/your-org/carrot-agent/actions`
   - 你会看到 "Deploy Website" 工作流正在运行

2. **检查部署状态**
   - 点击正在进行的工作流
   - 查看详细日志
   - 等待所有步骤完成（通常需要 2-5 分钟）

3. **成功标志**
   - ✅ 所有步骤显示绿色对勾
   - ✅ "Deploy to GitHub Pages" 步骤成功
   - ✅ 显示部署 URL

### 步骤 5：访问网站

部署成功后，你的网站将在以下地址可用：

```
https://your-org.github.io/carrot-agent/
```

中文版本：
```
https://your-org.github.io/carrot-agent/zh/
```

## 🔍 故障排除

### 问题 1：工作流没有触发

**可能原因**：
- 推送的分支不是 `main`
- 更改的文件不在 `website/` 目录下

**解决方案**：
```bash
# 确认当前分支
git branch

# 确认更改的文件
git status

# 如果需要，强制触发工作流
# 在 GitHub Actions 页面点击 "Run workflow"
```

### 问题 2：部署失败

**检查日志**：
1. 访问 Actions 标签
2. 点击失败的工作流
3. 查看错误信息

**常见错误**：

#### 错误：Permission denied
```
Error: Resource not accessible by integration
```

**解决**：
- 确认 Settings → Pages 中选择了 "GitHub Actions"
- 确认工作流中有正确的 permissions 配置

#### 错误：npm install 失败
```
Error: npm ERR! code ENOENT
```

**解决**：
- 确认 `website/package.json` 存在
- 确认 `website/package-lock.json` 已提交

#### 错误：Build failed
```
Error: Build failed with exit code 1
```

**解决**：
- 本地测试构建：`cd website && npm run docs:build`
- 检查 Markdown 语法错误
- 检查配置文件语法

### 问题 3：网站显示 404

**可能原因**：
- 部署尚未完成
- 基础路径配置不正确

**解决方案**：

1. **等待几分钟**，部署可能需要时间传播

2. **检查 VitePress 配置**，如果需要自定义 base：

```typescript
// website/docs/.vitepress/config.ts
export default defineConfig({
  base: '/carrot-agent/',  // 如果仓库名不是用户名.github.io
  // ... 其他配置
})
```

3. **清除浏览器缓存**
   - Chrome: Ctrl+Shift+Delete
   - 或使用无痕模式访问

### 问题 4：网站样式丢失

**可能原因**：
- 静态资源路径问题
- Base 路径配置错误

**解决方案**：

在 `config.ts` 中添加 base 配置：

```typescript
export default defineConfig({
  base: '/carrot-agent/',  // 替换为你的仓库名
  // ...
})
```

然后重新提交：

```bash
git add website/docs/.vitepress/config.ts
git commit -m "Fix base path for GitHub Pages"
git push origin main
```

## 🎯 优化建议

### 1. 添加部署状态徽章

在 README.md 中添加：

```markdown
![Deploy Status](https://github.com/your-org/carrot-agent/actions/workflows/deploy-website.yml/badge.svg)
```

### 2. 配置自定义域名（可选）

1. **创建 CNAME 文件**
   ```bash
   echo "docs.carrot-agent.com" > website/docs/public/CNAME
   ```

2. **配置 DNS**
   - 在你的域名提供商处添加 CNAME 记录
   - 指向 `your-org.github.io`

3. **提交并推送**
   ```bash
   git add website/docs/public/CNAME
   git commit -m "Add custom domain"
   git push origin main
   ```

### 3. 添加环境变量保护

如果将来需要 API 密钥等敏感信息：

1. **在 GitHub 设置中添加 Secrets**
   - Settings → Secrets and variables → Actions
   - 添加新的 secret

2. **在工作流中使用**
   ```yaml
   env:
     API_KEY: ${{ secrets.API_KEY }}
   ```

## 📊 监控和维护

### 查看部署历史

```
https://github.com/your-org/carrot-agent/deployments
```

### 手动触发部署

如果需要重新部署：

1. 访问 Actions 标签
2. 点击 "Deploy Website" 工作流
3. 点击右上角 "Run workflow"
4. 选择分支（通常是 main）
5. 点击 "Run workflow"

### 查看网站统计

使用以下工具监控网站访问：
- Google Analytics
- Cloudflare Analytics
- GitHub Traffic（仓库 Insights 标签）

## ✅ 验证清单

部署前确认：

- [ ] 工作流文件存在于 `.github/workflows/deploy-website.yml`
- [ ] GitHub Pages 设置中选择 "GitHub Actions" 作为源
- [ ] 代码已推送到 main 分支
- [ ] website 目录包含所有必要文件
- [ ] package.json 和 package-lock.json 已提交
- [ ] 本地构建测试通过（`npm run docs:build`）

部署后确认：

- [ ] Actions 工作流成功完成
- [ ] 可以访问网站 URL
- [ ] 英文页面正常显示
- [ ] 中文页面正常显示
- [ ] 语言切换功能正常
- [ ] 所有链接工作正常
- [ ] 样式和布局正确

## 🚀 快速命令参考

```bash
# 本地测试构建
cd website
npm install
npm run docs:build

# 预览构建结果
npm run docs:preview

# 提交并触发部署
cd ..
git add .
git commit -m "Update website"
git push origin main

# 查看部署状态
# 访问: https://github.com/your-org/carrot-agent/actions
```

## 📞 获取帮助

如果遇到问题：

1. **查看工作流日志** - 最直接的错误信息
2. **检查 GitHub Status** - https://www.githubstatus.com/
3. **查阅 VitePress 文档** - https://vitepress.dev/guide/deploy
4. **搜索 GitHub Issues** - 可能有其他人遇到类似问题

---

**最后更新**: 2024年
**维护者**: Carrot Agent Team
