# DocConverter PDF生成问题解决方案

## 问题描述

目前 `docconverter` 工具在转换网页到PDF时遇到三个主要问题：

1. **右侧内容缺失**：GitBook类型的页面有侧边栏导航，标准工具无法捕获
2. **滚动内容缺失**：页面需要滚动才能看到的内容没有被包含
3. **自动化工具限制**：wkhtmltopdf和Chrome headless的视口限制

## 备选方案

如果目前版本工具不满足需求，可以使用以下手动方法：

### 1. Chrome开发者工具截图

1. 在Chrome中打开目标网页
2. 按 `F12` 打开开发者工具
3. 按 `Cmd+Shift+P`（Mac）或 `Ctrl+Shift+P`（Windows/Linux）
4. 输入并选择 "Capture full size screenshot"
5. 将生成的PNG图片转换为PDF

### 2. Chrome打印功能

1. 在Chrome中打开网页
2. 按 `Cmd+P`（Mac）或 `Ctrl+P`（Windows/Linux）
3. 设置：
   - 目标：另存为PDF
   - 布局：横向
   - 纸张尺寸：A3或Tabloid
   - 边距：无
   - 缩放：70-80%
   - 勾选"背景图形"

### 3. 专业工具

- **Playwright**: 更强大的自动化工具
  ```bash
  pip3 install playwright
  playwright install chromium
  ```

- **Puppeteer**: Node.js环境下的选择
  ```bash
  npm install puppeteer
  ```

- **Prince**: 专业的HTML转PDF工具
  ```bash
  brew install --cask prince
  prince [URL] -o output.pdf
  ```

## 性能对比

| 方案 | 速度 | 完整性 | 文件大小 | 稳定性 |
|-----|------|--------|---------|--------|
| wkhtmltopdf (原始) | 慢 | 差 | 小 | 低 |
| Chrome Headless | 快 | 好 | 中 | 高 |
| Chrome 超宽格式 | 快 | 最好 | 大 | 高 |
| Playwright | 中 | 最好 | 中 | 最高 |
