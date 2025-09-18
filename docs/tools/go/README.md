# Go工具集文档

## 概述

本目录包含了devops-toolkit项目中所有Go语言开发的工具文档。这些工具涵盖了系统监控、性能分析、文档处理等多个方面，是运维和开发人员的得力助手。

## 工具列表

### 1. [NodeProbe](./NodeProbe.md) - Linux节点配置探测工具
**版本**: 1.1.1  
**功能**: 全面采集服务器硬件配置、系统状态和软件环境信息  
**特色**:
- 自动优化系统设置（CPU性能模式、时区等）
- 支持多格式输出（表格、JSON、YAML）
- 内存插槽物理位置显示
- 智能内核模块管理

**快速使用**:
```bash
# 默认表格输出
sudo nodeprobe

# JSON格式输出到文件
sudo nodeprobe -format json -output server_info.json
```

### 2. [PerfSnap](./PerfSnap.md) - Linux系统性能快照分析工具
**版本**: 1.1.0  
**功能**: 快速采集和分析CPU、内存、磁盘IO、网络等关键性能指标  
**特色**:
- 一键生成性能报告
- 实时监控模式
- CPU火焰图生成
- 智能性能问题诊断

**快速使用**:
```bash
# 生成性能快照
sudo perfsnap

# 生成火焰图
sudo perfsnap -flame

# 实时监控
sudo perfsnap -m -interval 2 -duration 60
```

### 3. [DocConverter](./DocConverter.md) - 文档转PDF工具
**版本**: 1.1.0  
**功能**: 将本地文件（HTML、Markdown）或在线网页内容转换为PDF  
**特色**:
- 支持网站爬取和批量转换
- 图片自动下载和本地化
- 自定义页面设置和样式
- 支持多种PDF引擎

**快速使用**:
```bash
# 转换Markdown文件
docconverter -i README.md -o readme.pdf

# 爬取网站生成PDF
docconverter -i https://example.com/docs/ -o docs.pdf --max-depth 3
```

## 工具对比

| 工具 | 用途 | 数据类型 | 执行频率 | 需要Root |
|------|------|---------|----------|----------|
| NodeProbe | 硬件配置采集 | 静态信息 | 低频 | 推荐 |
| PerfSnap | 性能监控分析 | 动态数据 | 高频 | 推荐 |
| DocConverter | 文档处理 | 文件转换 | 按需 | 否 |

## 安装指南

### 系统要求
- Linux操作系统（NodeProbe、PerfSnap）
- macOS/Linux/Windows（DocConverter）
- Go 1.15+（编译时需要）

### 快速安装

```bash
# 克隆项目
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit/tools/go

# 编译所有工具
go build -o nodeprobe NodeProbe.go
go build -o perfsnap PerfSnap.go
go build -o docconverter DocConverter.go

# 安装到系统路径
sudo mv nodeprobe perfsnap docconverter /usr/local/bin/

# 验证安装
nodeprobe -version
perfsnap -version
docconverter -version
```

### 依赖安装

#### NodeProbe依赖
```bash
# 可选，增强功能
sudo yum install -y dmidecode ethtool  # CentOS/RHEL
sudo apt-get install -y dmidecode ethtool  # Ubuntu/Debian
```

#### PerfSnap依赖
```bash
# 必需
sudo yum install -y sysstat  # CentOS/RHEL
sudo apt-get install -y sysstat  # Ubuntu/Debian

# 火焰图功能（可选）
sudo yum install -y perf git perl  # CentOS/RHEL
sudo apt-get install -y linux-tools-common linux-tools-generic git perl  # Ubuntu
```

#### DocConverter依赖
```bash
# 必需
# macOS
brew install --cask wkhtmltopdf

# Linux
sudo apt-get install -y wkhtmltopdf  # Ubuntu/Debian
sudo yum install -y wkhtmltopdf  # CentOS/RHEL

# 可选
brew install ghostscript  # macOS
sudo apt-get install -y ghostscript  # Ubuntu/Debian
```

## 使用场景

### 场景1：新服务器上线检查
```bash
#!/bin/bash
# 完整的服务器检查流程

# 1. 收集配置信息
sudo nodeprobe -format json -output config.json

# 2. 建立性能基线
sudo perfsnap > performance_baseline.txt

# 3. 生成检查报告PDF
echo "# 服务器检查报告" > report.md
echo "## 配置信息" >> report.md
sudo nodeprobe >> report.md
echo "## 性能基线" >> report.md
sudo perfsnap >> report.md

docconverter -i report.md -o server_report.pdf
```

### 场景2：故障诊断
```bash
#!/bin/bash
# 快速故障诊断

# 1. 检查系统配置是否变更
sudo nodeprobe > current_config.txt
diff baseline_config.txt current_config.txt

# 2. 分析当前性能
sudo perfsnap

# 3. 如果CPU异常，生成火焰图
sudo perfsnap -flame -flame-duration 60

# 4. 持续监控
sudo perfsnap -m -interval 1 -duration 300
```

### 场景3：批量服务器管理
```bash
#!/bin/bash
# 批量收集服务器信息

SERVERS="server1 server2 server3"
DATE=$(date +%Y%m%d)

for server in $SERVERS; do
    echo "Processing $server..."
    
    # 收集配置
    ssh root@$server "nodeprobe -format json" > ${server}_config_${DATE}.json
    
    # 收集性能
    ssh root@$server "perfsnap" > ${server}_perf_${DATE}.txt
done

# 生成汇总报告
cat *_config_${DATE}.json | jq -s '.' > all_configs.json
cat *_perf_${DATE}.txt > all_performance.txt
```

### 场景4：文档归档
```bash
#!/bin/bash
# 项目文档归档

# 转换所有Markdown文档
docconverter -i ./docs -o project_docs.pdf \
    --header "项目文档 v1.0" \
    --footer "[page] / [topage]"

# 备份在线文档
docconverter -i https://wiki.example.com/project/ \
    -o wiki_backup.pdf \
    --max-depth 5
```

## 最佳实践

### 1. 定期巡检
```bash
# crontab配置
# 每天采集配置信息
0 9 * * * /usr/local/bin/nodeprobe -quiet -output /var/log/nodeprobe/$(date +\%Y\%m\%d).json

# 每小时性能快照
0 * * * * /usr/local/bin/perfsnap > /var/log/perfsnap/$(date +\%Y\%m\%d_\%H).txt

# 每周生成报告
0 10 * * 1 /usr/local/bin/weekly_report.sh
```

### 2. 监控集成
```python
#!/usr/bin/env python3
# 集成到监控系统

import subprocess
import json
import requests

# 收集节点信息
result = subprocess.run(['nodeprobe', '-format', 'json'], 
                       capture_output=True, text=True)
node_data = json.loads(result.stdout)

# 发送到监控系统
requests.post('http://monitoring.example.com/api/metrics', 
              json=node_data)
```

### 3. CI/CD集成
```yaml
# .gitlab-ci.yml
stages:
  - test
  - deploy
  - verify

verify-deployment:
  stage: verify
  script:
    - ssh $DEPLOY_HOST "nodeprobe -format json" > node_info.json
    - ssh $DEPLOY_HOST "perfsnap" > performance.txt
    - python3 verify_deployment.py node_info.json performance.txt
  artifacts:
    paths:
      - node_info.json
      - performance.txt
```

## 故障排查

### 通用问题

1. **权限不足**
   ```bash
   错误: Permission denied
   解决: 使用sudo运行或检查文件权限
   ```

2. **命令未找到**
   ```bash
   错误: command not found
   解决: 确保工具已安装并在PATH中
   export PATH=$PATH:/usr/local/bin
   ```

3. **依赖缺失**
   ```bash
   错误: xxx: No such file or directory
   解决: 安装相应的依赖包
   ```

### 工具特定问题

- **NodeProbe**: 查看[故障排查](./NodeProbe.md#故障排查)
- **PerfSnap**: 查看[故障排查](./PerfSnap.md#故障排查)
- **DocConverter**: 查看[故障排查](./DocConverter.md#故障排查)

## 开发指南

### 编译所有工具
```bash
#!/bin/bash
# build-all.sh

TOOLS="NodeProbe PerfSnap DocConverter"
BUILD_DIR="./build"

mkdir -p $BUILD_DIR

for tool in $TOOLS; do
    echo "Building $tool..."
    output_name=$(echo $tool | tr '[:upper:]' '[:lower:]')
    go build -o $BUILD_DIR/$output_name ${tool}.go
done

echo "Build complete. Binaries in $BUILD_DIR"
```

### 跨平台编译
```bash
#!/bin/bash
# cross-compile.sh

PLATFORMS="linux/amd64 darwin/amd64 windows/amd64"

for platform in $PLATFORMS; do
    GOOS=${platform%/*}
    GOARCH=${platform#*/}
    
    echo "Building for $GOOS/$GOARCH..."
    
    # NodeProbe和PerfSnap仅支持Linux
    if [ "$GOOS" = "linux" ]; then
        GOOS=$GOOS GOARCH=$GOARCH go build -o nodeprobe-$GOOS-$GOARCH NodeProbe.go
        GOOS=$GOOS GOARCH=$GOARCH go build -o perfsnap-$GOOS-$GOARCH PerfSnap.go
    fi
    
    # DocConverter支持所有平台
    ext=""
    [ "$GOOS" = "windows" ] && ext=".exe"
    GOOS=$GOOS GOARCH=$GOARCH go build -o docconverter-$GOOS-$GOARCH$ext DocConverter.go
done
```

## 版本管理

| 工具 | 当前版本 | 最后更新 | 主要变更 |
|------|---------|---------|---------|
| NodeProbe | 1.1.1 | 2025-09-15 | 修复内存插槽显示 |
| PerfSnap | 1.1.0 | 2025-09-15 | 新增火焰图功能 |
| DocConverter | 1.1.0 | 2025-09-18 | 新增图片下载功能 |

## 贡献指南

1. Fork项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启Pull Request

### 代码规范
- 使用`go fmt`格式化代码
- 添加适当的注释
- 更新相关文档
- 编写测试用例

## 许可证

MIT License - 详见[LICENSE](../../../LICENSE)文件

## 联系方式

- **作者**: sunyifei83@gmail.com
- **项目**: https://github.com/sunyifei83/devops-toolkit
- **Issues**: https://github.com/sunyifei83/devops-toolkit/issues

## 相关链接

- [Shell工具文档](../shell/README.md)
- [Python工具文档](../python/README.md)
- [项目主页](https://github.com/sunyifei83/devops-toolkit)
