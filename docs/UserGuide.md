# DevOps Toolkit 用户指南

欢迎使用 DevOps Toolkit！本指南将帮助您快速上手并充分利用工具箱中的各种工具。

## 📖 目录

- [快速开始](#快速开始)
- [环境要求](#环境要求)
- [安装指南](#安装指南)
- [工具使用](#工具使用)
- [常见场景](#常见场景)
- [故障排查](#故障排查)

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit
```

### 2. 快速体验

**使用 NodeProbe 检查服务器配置**：
```bash
cd tools/go
go run NodeProbe.go
# 或使用 sudo 获取完整信息
sudo go run NodeProbe.go
```

**使用 iotest 测试磁盘性能**：
```bash
cd tools/shell
chmod +x iotest.sh
./iotest.sh
```

## 💻 环境要求

### 操作系统
- Linux（推荐 Ubuntu 18.04+, CentOS 7+）
- macOS（部分功能）
- Windows（通过 WSL2）

### 依赖软件

#### Go 工具
```bash
# Go 1.15 或更高版本
go version
```

#### Shell 工具
```bash
# Bash 4.0+
bash --version

# 常用工具
sudo apt-get install -y curl wget  # Ubuntu/Debian
sudo yum install -y curl wget      # CentOS/RHEL
```

#### 可选依赖
- **fio**: 磁盘性能测试（iotest.sh）
- **ethtool**: 网络接口信息（NodeProbe）
- **dmidecode**: 硬件信息（NodeProbe，需要 root）

## 📦 安装指南

### 方式一：编译安装（推荐）

```bash
# 进入 Go 工具目录
cd tools/go

# 编译所有工具
go build -o nodeprobe NodeProbe.go
go build -o perfsnap PerfSnap.go
go build -o docconverter DocConverter.go

# 安装到系统路径（可选）
sudo cp nodeprobe perfsnap docconverter /usr/local/bin/

# 设置 Shell 脚本权限
cd ../shell
chmod +x *.sh
```

### 方式二：直接运行

```bash
# 无需编译，直接运行
go run tools/go/NodeProbe.go
bash tools/shell/iotest.sh
```

### 方式三：添加到 PATH

```bash
# 添加到 ~/.bashrc 或 ~/.zshrc
export PATH=$PATH:/path/to/devops-toolkit/tools/go
export PATH=$PATH:/path/to/devops-toolkit/tools/shell

# 重新加载配置
source ~/.bashrc
```

## 🛠️ 工具使用

### NodeProbe - 节点配置探测

**基本使用**：
```bash
# 查看节点配置信息
nodeprobe

# 以 root 权限运行（获取完整硬件信息）
sudo nodeprobe

# 输出为 JSON 格式
nodeprobe -format json

# 输出为 YAML 格式
nodeprobe -format yaml

# 保存到文件
nodeprobe -format json -output server-info.json
```

**使用场景**：
- 服务器上线前的配置检查
- 批量收集集群节点信息
- 系统性能基线建立
- 故障排查时的环境信息收集

**权限说明**：
- 普通用户：可以查看基本系统信息
- Root 用户：可以查看完整硬件信息（内存插槽、CPU性能模式等）

### PerfSnap - 性能快照

**基本使用**：
```bash
# 采集系统性能快照
perfsnap

# 指定采集时长（秒）
perfsnap -duration 60

# 生成火焰图
perfsnap -flamegraph

# 保存结果
perfsnap -output perf-snapshot.json
```

**使用场景**：
- 性能问题诊断
- 系统瓶颈分析
- 性能基线建立
- 容量规划参考

### DocConverter - 文档转换

**基本使用**：
```bash
# Markdown 转 PDF
docconverter -input README.md -output README.pdf

# HTML 转 PDF
docconverter -input index.html -output document.pdf

# 网页转 PDF
docconverter -url https://example.com -output webpage.pdf

# 批量转换
for file in *.md; do
  docconverter -input "$file" -output "${file%.md}.pdf"
done
```

**使用场景**：
- 技术文档归档
- 报告生成
- 离线阅读准备

### iotest - 磁盘 IO 测试

**基本使用**：
```bash
# 运行 IO 性能测试
./iotest.sh

# 测试会自动：
# 1. 检测系统架构
# 2. 下载或使用本地 fio
# 3. 执行多种块大小的测试
# 4. 显示 IOPS 和吞吐量结果
```

**使用场景**：
- 存储性能评估
- 磁盘基准测试
- 性能对比分析
- 故障排查

**注意事项**：
- 需要至少 2GB 可用空间（x86_64）或 512MB（ARM）
- 测试会产生临时文件
- ZFS 文件系统需要更多空间

### clear_log - 日志清理

**基本使用**：
```bash
# 清理系统日志
./clear_log.sh

# 查看帮助
./clear_log.sh --help
```

**使用场景**：
- 磁盘空间清理
- 日志维护
- 定期清理任务

## 🎯 常见场景

### 场景1：新服务器上线检查

```bash
# 1. 检查服务器配置
sudo nodeprobe -format json -output server-config.json

# 2. 测试磁盘性能
cd tools/shell && ./iotest.sh

# 3. 采集性能基线
cd ../go && perfsnap -duration 300 -output baseline.json

# 4. 生成配置报告
docconverter -input server-config.json -output server-report.pdf
```

### 场景2：性能问题诊断

```bash
# 1. 采集当前性能快照
perfsnap -duration 60 -flamegraph -output problem.json

# 2. 检查系统配置
sudo nodeprobe

# 3. 测试磁盘性能
cd tools/shell && ./iotest.sh

# 4. 对比历史数据
diff baseline.json problem.json
```

### 场景3：批量收集集群信息

```bash
#!/bin/bash
# collect_cluster_info.sh

NODES="node1 node2 node3"
OUTPUT_DIR="cluster_info"
mkdir -p $OUTPUT_DIR

for node in $NODES; do
  echo "Collecting info from $node..."
  ssh $node "cd /path/to/devops-toolkit/tools/go && sudo ./nodeprobe -format json" \
    > $OUTPUT_DIR/${node}_config.json
done

echo "Collection completed. Results in $OUTPUT_DIR/"
```

### 场景4：定期性能监控

```bash
# 添加到 crontab
# 每天凌晨 2 点采集性能数据
0 2 * * * cd /path/to/devops-toolkit/tools/go && ./perfsnap -output /var/log/perf/$(date +\%Y\%m\%d).json
```

### 场景5：文档批量归档

```bash
#!/bin/bash
# 批量转换项目文档为 PDF

cd /path/to/project/docs
for md in *.md; do
  echo "Converting $md..."
  docconverter -input "$md" -output "pdf/${md%.md}.pdf"
done
```

## 🔧 故障排查

### 问题1：NodeProbe 无法获取完整硬件信息

**症状**：内存插槽信息为空，CPU 性能模式显示"未知"

**解决**：
```bash
# 使用 sudo 运行
sudo nodeprobe

# 检查 dmidecode 是否安装
which dmidecode || sudo apt-get install dmidecode
```

### 问题2：iotest 提示空间不足

**症状**：`Less than 2GB of space available. Skipping disk test...`

**解决**：
```bash
# 检查可用空间
df -h .

# 清理临时文件
rm -rf /tmp/yabs_*

# 或更换到有足够空间的目录
cd /data && /path/to/iotest.sh
```

### 问题3：Go 工具编译失败

**症状**：`go: cannot find module...`

**解决**：
```bash
# 更新依赖
cd tools/go
go mod download
go mod tidy

# 重新编译
go build NodeProbe.go
```

### 问题4：权限被拒绝

**症状**：`Permission denied` 错误

**解决**：
```bash
# Shell 脚本
chmod +x tools/shell/*.sh

# 或使用 bash 运行
bash tools/shell/iotest.sh
```

### 问题5：PerfSnap 采集失败

**症状**：无法采集性能数据

**解决**：
```bash
# 检查是否有足够权限
sudo perfsnap

# 检查系统负载
uptime

# 检查依赖工具
which perf || sudo apt-get install linux-tools-common
```

## 📚 进阶使用

### 集成到监控系统

```bash
# 将 NodeProbe 输出集成到 Prometheus
nodeprobe -format json | jq -r '.cpu.cores' | \
  curl -X POST http://pushgateway:9091/metrics/job/nodeprobe \
  --data-binary @-
```

### 自动化脚本

```bash
#!/bin/bash
# auto_check.sh - 自动检查脚本

LOG_DIR="/var/log/devops-toolkit"
mkdir -p $LOG_DIR

DATE=$(date +%Y%m%d_%H%M%S)

# 采集配置信息
sudo nodeprobe -format json -output $LOG_DIR/config_$DATE.json

# 性能快照
perfsnap -duration 60 -output $LOG_DIR/perf_$DATE.json

# 磁盘测试（每周一次）
if [ $(date +%u) -eq 1 ]; then
  cd /path/to/devops-toolkit/tools/shell
  ./iotest.sh > $LOG_DIR/iotest_$DATE.log 2>&1
fi

echo "Check completed: $DATE"
```

### CI/CD 集成

```yaml
# .github/workflows/server-check.yml
name: Server Check
on:
  schedule:
    - cron: '0 0 * * *'  # 每天运行

jobs:
  check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Setup Go
        uses: actions/setup-go@v2
        with:
          go-version: '1.19'
      - name: Run NodeProbe
        run: |
          cd tools/go
          go run NodeProbe.go -format json > server-info.json
      - name: Upload artifact
        uses: actions/upload-artifact@v2
        with:
          name: server-info
          path: tools/go/server-info.json
```

## 🆘 获取帮助

### 文档资源
- [工具文档](/docs/ToolsDocumentation.md) - 所有工具的详细文档
- [最佳实践](/docs/BestPractices.md) - 使用最佳实践
- [项目 README](/README.md) - 项目概览

### 社区支持
- **GitHub Issues**: https://github.com/sunyifei83/devops-toolkit/issues
- **邮件**: sunyifei83@gmail.com
- **Twitter**: @sunyifei83

### 工具帮助命令
```bash
# 查看工具版本和帮助
nodeprobe -version
nodeprobe -help

perfsnap -help
docconverter -help
```

## 🔄 更新工具

```bash
# 更新项目
cd devops-toolkit
git pull origin main

# 重新编译 Go 工具
cd tools/go
go build NodeProbe.go
go build PerfSnap.go
go build DocConverter.go

# 检查更新日志
cat CHANGELOG.md  # 如果存在
```

## ⚖️ 许可证

本项目采用 MIT 许可证。详见 [LICENSE](../LICENSE) 文件。

---

**提示**：如果您在使用过程中遇到问题或有改进建议，欢迎提交 Issue 或 Pull Request！

*最后更新：2025-10-27*
