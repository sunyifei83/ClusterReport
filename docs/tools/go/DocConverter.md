# DocConverter - 文档转PDF工具

## 概述

`DocConverter` 是一个功能强大的文档转PDF转换工具，支持将本地文件（HTML、Markdown）或在线网页内容转换为PDF格式。该工具特别适合文档归档、离线阅读、技术文档生成等场景。

**版本**: 1.1.0  
**作者**: sunyifei83@gmail.com  
**项目**: https://github.com/sunyifei83/devops-toolkit  
**更新日期**: 2025-09-18

## 核心特性

- 📄 **多格式支持**: 支持HTML、Markdown文件转换
- 🌐 **网页抓取**: 支持单个URL或整站内容抓取转换，支持图片下载
- 📁 **批量处理**: 支持目录递归处理，批量转换文档
- 🎨 **样式定制**: 支持自定义CSS样式，内置优化的默认样式
- 📖 **目录生成**: 自动生成PDF目录（多文件时）
- 🔧 **灵活配置**: 丰富的页面设置选项
- 🚀 **高效转换**: 基于wkhtmltopdf引擎，支持Chrome headless作为备选
- 🕷️ **智能爬取**: 支持深度爬取网站内容，自动下载和转换图片
- 🔀 **智能合并**: 支持Ghostscript合并多个PDF

## 主要功能

### 1. 本地文件转换
- 单个文件转换（Markdown、HTML）
- 目录批量转换
- 文件过滤（包含/排除模式）
- 递归处理子目录
- 临时文件自动清理

### 2. 网页内容转换
- 单个网页URL转换
- 网站内容爬取（可设置深度）
- 自动跟踪同域名链接
- 保持页面原始样式
- 图片自动下载和本地化
- 支持自定义User-Agent

### 3. PDF定制选项
- 页面大小（A4、A3、Letter等）
- 页面方向（纵向/横向）
- 页边距设置（上下左右独立配置）
- 页眉页脚定制
- 目录生成（可配置深度）

### 4. 高级特性
- Markdown扩展语法支持（表格、自动标题ID等）
- 代码语法高亮
- 自定义CSS注入
- 错误容错处理
- Chrome headless备选方案
- 多URL稳定性处理

## 安装部署

### 系统要求
- Go 1.16或更高版本（编译时需要）
- wkhtmltopdf（核心转换引擎）
- Ghostscript（可选，用于PDF合并）
- Chrome/Chromium（可选，备用转换引擎）

### 依赖安装

#### 1. 安装wkhtmltopdf
```bash
# macOS
brew install --cask wkhtmltopdf

# Ubuntu/Debian
sudo apt-get install wkhtmltopdf

# CentOS/RHEL
sudo yum install wkhtmltopdf

# 或从官网下载安装包
# https://wkhtmltopdf.org/downloads.html
```

#### 2. 安装Ghostscript（可选，推荐）
```bash
# macOS
brew install ghostscript

# Ubuntu/Debian
sudo apt-get install ghostscript

# CentOS/RHEL
sudo yum install ghostscript
```

#### 3. 安装Go依赖
```bash
# 进入工具目录
cd tools/go

# 初始化Go模块（如果没有go.mod）
go mod init devops-toolkit/tools

# 安装依赖
go get github.com/PuerkitoBio/goquery
go get github.com/gomarkdown/markdown
go get github.com/gomarkdown/markdown/html
go get github.com/gomarkdown/markdown/parser
```

### 编译安装

```bash
# 1. 克隆项目
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit

# 2. 编译DocConverter
cd tools/go
go build -o docconverter DocConverter.go

# 3. 设置执行权限
chmod +x docconverter

# 4. (可选) 移动到系统路径
sudo mv docconverter /usr/local/bin/

# 5. 验证安装
docconverter -version
```

### 跨平台编译

```bash
# macOS编译
GOOS=darwin GOARCH=amd64 go build -o docconverter-darwin-amd64 DocConverter.go

# Linux编译
GOOS=linux GOARCH=amd64 go build -o docconverter-linux-amd64 DocConverter.go

# Windows编译
GOOS=windows GOARCH=amd64 go build -o docconverter.exe DocConverter.go
```

## 使用方法

### 基础用法

```bash
# 查看帮助
docconverter -h

# 查看版本
docconverter -version

# 转换单个Markdown文件
docconverter -i README.md -o readme.pdf

# 转换单个HTML文件
docconverter -i index.html -o output.pdf

# 转换整个目录
docconverter -i ./docs -o documentation.pdf
```

### 网页转换

```bash
# 转换单个网页
docconverter -i https://example.com -o website.pdf

# 爬取并转换整个网站（深度3）
docconverter -i https://example.com/docs/ -o docs.pdf --max-depth 3

# 下载网页中的图片
docconverter -i https://example.com -o site.pdf --download-images

# 自定义User-Agent
docconverter -i https://example.com -o site.pdf --user-agent "Mozilla/5.0"

# 设置超时时间
docconverter -i https://example.com -o site.pdf --timeout 60
```

### 高级用法

```bash
# 自定义页面设置
docconverter -i ./docs -o docs.pdf \
  --page-size A3 \
  --orientation Landscape \
  --margin-top 20mm \
  --margin-bottom 20mm

# 添加页眉页脚
docconverter -i ./docs -o report.pdf \
  --header "项目文档" \
  --footer "[page] / [topage]"

# 使用自定义CSS
docconverter -i ./docs -o styled.pdf --css custom.css

# 文件过滤
docconverter -i ./src -o code.pdf \
  --include "*.md,*.html" \
  --exclude "test*,temp*"

# 详细输出模式
docconverter -i ./docs -o docs.pdf -v

# 不生成目录
docconverter -i ./docs -o docs.pdf --toc=false

# 保留临时文件（调试用）
docconverter -i ./docs -o docs.pdf --clean=false
```

### 命令行参数

#### 基础参数
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-i, --input` | 输入路径(文件、目录或URL) | 必需 |
| `-o, --output` | 输出PDF文件名 | output.pdf |
| `-r, --recursive` | 递归处理子目录 | true |
| `-v, --verbose` | 详细输出模式 | false |
| `--version` | 显示版本信息 | - |
| `-h, --help` | 显示帮助信息 | - |

#### 文件过滤
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--include` | 包含的文件模式(逗号分隔) | *.html,*.md,*.markdown |
| `--exclude` | 排除的文件模式(逗号分隔) | 空 |

#### 页面设置
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--page-size` | 页面大小 | A4 |
| `--orientation` | 页面方向(Portrait/Landscape) | Portrait |
| `--margin-top` | 顶部边距 | 15mm |
| `--margin-bottom` | 底部边距 | 15mm |
| `--margin-left` | 左边距 | 15mm |
| `--margin-right` | 右边距 | 15mm |

#### 页眉页脚
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--header` | 页眉标题 | 空 |
| `--footer` | 页脚文本 | [page] / [topage] |

#### 目录设置
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--toc` | 生成目录 | true |
| `--toc-depth` | 目录深度 | 3 |

#### Web爬取设置
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--max-depth` | URL爬取最大深度 | 3 |
| `--timeout` | HTTP请求超时时间(秒) | 30 |
| `--user-agent` | HTTP User-Agent | Mozilla/5.0 (DocConverter/1.1.0) |
| `--download-images` | 下载并转换网页中的图片 | true |

#### 其他选项
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--css` | 自定义CSS文件路径 | 空 |
| `--clean` | 转换后清理临时文件 | true |

## 使用场景

### 1. 技术文档生成
```bash
# 将项目文档转换为PDF
docconverter -i ./docs -o project-docs.pdf --header "项目文档 v1.0"
```

### 2. API文档归档
```bash
# 抓取在线API文档
docconverter -i https://api.example.com/docs/ -o api-docs.pdf --max-depth 5
```

### 3. 博客文章备份
```bash
# 转换Markdown博客文章
docconverter -i ./posts -o blog-backup.pdf --include "*.md"
```

### 4. 报告生成
```bash
# 生成格式化报告
docconverter -i report.md -o report.pdf \
  --header "月度报告" \
  --footer "机密文件 - [page]" \
  --css report-style.css
```

### 5. 电子书制作
```bash
# 批量转换章节文件
docconverter -i ./chapters -o ebook.pdf \
  --toc \
  --toc-depth 2 \
  --page-size A5
```

## 输出示例

### 转换本地文件
```
DocConverter - 文档转PDF工具
找到 15 个文件待转换
[1/15] 处理: docs/intro.md
[2/15] 处理: docs/chapter1.md
[3/15] 处理: docs/chapter2.md
...
执行命令: wkhtmltopdf --enable-local-file-access ...
✅ PDF生成成功: documentation.pdf
文件大小: 2.45 MB
页面数量: 15 个源文件
```

### 爬取网站
```
检测到URL输入: https://example.com/docs/
开始爬取网站内容...
爬取页面 [深度:0]: https://example.com/docs/
✓ 下载图片: /images/logo.png
爬取页面 [深度:1]: https://example.com/docs/getting-started
爬取页面 [深度:1]: https://example.com/docs/api-reference
爬取页面 [深度:2]: https://example.com/docs/api/users
...
发现 23 个页面URL
找到 23 个文件/URL待转换
✅ PDF生成成功: website-docs.pdf
文件大小: 5.12 MB
页面数量: 23 个源文件
```

### 多URL处理
```
⚠️ 检测到多个URL (10)，将单独处理每个页面以确保稳定性
[1/10] 处理页面: https://example.com/page1
[2/10] 处理页面: https://example.com/page2
⚠️ 页面 3 处理失败: connection timeout，跳过
[4/10] 处理页面: https://example.com/page4
...
合并 9 个成功的PDF文件...
✅ PDF生成成功: merged.pdf
```

## 内置CSS样式

DocConverter内置了优化的CSS样式，适合打印和PDF输出：

```css
/* 主要特性 */
- 优化的字体栈（跨平台兼容）
- 清晰的标题层级（带下划线）
- 代码块样式（背景色和边框）
- 表格样式（斑马纹）
- 引用块样式（左边框）
- 链接样式（黑色打印）
- 响应式图片（最大宽度100%）
```

## 故障排查

### 常见问题

1. **wkhtmltopdf未安装**
   ```bash
   错误: wkhtmltopdf未安装或不可用
   解决: 按照安装指南安装wkhtmltopdf
   ```

2. **Go依赖缺失**
   ```bash
   错误: cannot find package "github.com/PuerkitoBio/goquery"
   解决: 
   go get github.com/PuerkitoBio/goquery
   go get github.com/gomarkdown/markdown
   ```

3. **网页转换失败**
   ```bash
   错误: HTTP状态码: 403
   解决: 使用自定义User-Agent
   docconverter -i URL --user-agent "Mozilla/5.0..."
   ```

4. **中文乱码**
   ```bash
   问题: PDF中文字显示为方块
   解决: 安装中文字体
   # Ubuntu/Debian
   sudo apt-get install fonts-wqy-microhei fonts-wqy-zenhei
   # CentOS/RHEL
   sudo yum install wqy-microhei-fonts wqy-zenhei-fonts
   ```

5. **内存不足**
   ```bash
   问题: 转换大量文件时内存溢出
   解决: 分批处理或减少爬取深度
   ```

6. **CSS样式丢失**
   ```bash
   问题: 网页样式没有应用
   解决: 使用--css参数指定样式文件
   ```

7. **PDF合并失败**
   ```bash
   问题: 多个PDF无法合并
   解决: 安装Ghostscript
   brew install ghostscript  # macOS
   ```

8. **图片下载失败**
   ```bash
   问题: 网页中的图片未能下载
   解决: 检查网络连接，或使用--download-images=false跳过图片
   ```

## 最佳实践

### 1. 文档组织
```bash
# 按照逻辑顺序命名文件
01-introduction.md
02-installation.md
03-configuration.md
# 这样可以保证PDF中的顺序正确
```

### 2. 样式优化
```css
/* custom.css - 打印优化样式 */
@media print {
    .no-print { display: none; }
    a { color: black; text-decoration: none; }
    code { background: #f4f4f4; padding: 2px 4px; }
    pre { page-break-inside: avoid; }
    h1, h2, h3 { page-break-after: avoid; }
}
```

### 3. 批处理脚本
```bash
#!/bin/bash
# batch-convert.sh
for dir in project1 project2 project3; do
    docconverter -i ./$dir/docs -o ${dir}-docs.pdf
done
```

### 4. CI/CD集成
```yaml
# .github/workflows/docs.yml
- name: Generate PDF Documentation
  run: |
    docconverter -i ./docs -o docs.pdf --header "v${{ github.ref_name }}"
    
- name: Upload PDF
  uses: actions/upload-artifact@v2
  with:
    name: documentation
    path: docs.pdf
```

### 5. Makefile集成
```makefile
# Makefile
.PHONY: docs
docs:
	docconverter -i ./docs -o documentation.pdf \
		--header "项目文档" \
		--footer "[page] / [topage]" \
		--toc

clean:
	rm -f *.pdf
```

## 性能优化

1. **减少爬取深度**: 对于大型网站，限制爬取深度可以显著提高速度
2. **文件过滤**: 使用include/exclude参数减少处理的文件数量
3. **并行处理**: 可以将大批文件分组并行处理
4. **缓存利用**: 保留临时HTML文件可以加速重复转换
5. **禁用图片**: 如果不需要图片，使用--download-images=false

## 版本历史

### v1.1.0 (当前版本)
- 新增图片下载功能
- 支持Chrome headless作为备选引擎
- 改进多URL稳定性处理
- 支持Ghostscript PDF合并
- 优化默认CSS样式
- 改进错误处理和提示信息

### v1.0.0
- 初始版本发布
- 基础文件转换功能
- 网页爬取支持
- 自定义页面设置

## 注意事项

1. **版权问题**: 爬取网站内容时请遵守robots.txt和版权规定
2. **资源消耗**: 大量文件转换可能消耗较多CPU和内存
3. **网络限制**: 某些网站可能限制爬虫访问
4. **字体支持**: 确保系统安装了所需字体，特别是中文字体
5. **路径问题**: 使用绝对路径可以避免路径相关的问题

## 相关工具

本工具是devops-toolkit工具集的一部分，其他相关工具包括：

- **NodeProbe**: Linux节点配置探测工具
- **PerfSnap**: Linux系统性能快照分析工具
- **iotest.sh**: 磁盘IO性能测试工具

## 技术栈

- **语言**: Go 1.16+
- **转换引擎**: wkhtmltopdf / Chrome headless
- **Markdown解析**: gomarkdown
- **网页解析**: goquery
- **PDF合并**: Ghostscript

## 贡献指南

欢迎提交Issue和Pull Request。提交代码前请确保：

1. 代码通过`go fmt`格式化
2. 添加必要的注释
3. 更新相关文档
4. 测试各种场景

## 许可证

MIT License - 详见LICENSE文件

## 联系方式

- 作者: sunyifei83@gmail.com
- 项目: https://github.com/sunyifei83/devops-toolkit
- Issue: https://github.com/sunyifei83/devops-toolkit/issues
