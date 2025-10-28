# 快速入门指南

欢迎使用 ClusterReport！本指南将帮助您在 5 分钟内生成第一份集群报告。

## 🎯 学习目标

完成本指南后，您将能够：
- ✅ 安装 ClusterReport
- ✅ 采集本地系统数据
- ✅ 生成 HTML 报告
- ✅ 理解基本工作流程

## 📋 前置条件

- Linux 或 macOS 系统
- 具有管理员权限（用于系统信息采集）
- 5 分钟时间

## 🚀 快速开始

### 步骤 1: 安装 ClusterReport

选择以下任一安装方式：

#### 方式 A: 使用安装脚本（推荐）

```bash
curl -sSL https://raw.githubusercontent.com/sunyifei83/devops-toolkit/main/scripts/installation/install.sh | bash
```

#### 方式 B: 从源码编译

```bash
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit
make install
```

#### 方式 C: 使用 Docker

```bash
docker pull sunyifei83/clusterreport:latest
docker run -it sunyifei83/clusterreport:latest clusterreport --version
```

### 步骤 2: 验证安装

```bash
clusterreport --version
```

预期输出：
```
ClusterReport version 2.0.0
Built with NodeProbe Engine v1.0 and PerfSnap Engine v1.0
```

### 步骤 3: 采集本地数据

运行以下命令采集本地系统数据：

```bash
clusterreport collect --node localhost
```

您将看到类似输出：
```
🔍 Starting data collection...
📊 NodeProbe Engine: Collecting system configuration...
   ✓ CPU Information
   ✓ Memory Information
   ✓ Disk Information
   ✓ Network Interfaces
   ✓ OS Information
⚡ PerfSnap Engine: Collecting performance data...
   ✓ CPU Metrics
   ✓ Memory Metrics
   ✓ Disk I/O
   ✓ Network Metrics
✅ Data collection completed!
💾 Data saved to: cluster-data.json
```

### 步骤 4: 生成 HTML 报告

```bash
clusterreport generate --format html --output my-first-report.html
```

输出：
```
📄 Generating report...
   ✓ Processing data
   ✓ Analyzing metrics
   ✓ Creating visualizations
   ✓ Rendering HTML
✅ Report generated: my-first-report.html
```

### 步骤 5: 查看报告

```bash
# macOS
open my-first-report.html

# Linux
xdg-open my-first-report.html

# 或使用浏览器直接打开
```

## 🎉 恭喜！

您已经成功生成了第一份集群报告！报告包含：

- 📊 **系统配置概览** - CPU、内存、磁盘、网络配置
- ⚡ **性能指标** - 实时性能数据和趋势
- 🔍 **系统分析** - 潜在问题和优化建议
- 📈 **可视化图表** - 直观的数据展示

## 🔄 完整工作流程

```
┌──────────────┐
│   collect    │  采集数据
│ (NodeProbe + │
│  PerfSnap)   │
└──────┬───────┘
       │
       ▼
┌──────────────┐
│    analyze   │  分析数据（可选）
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   generate   │  生成报告
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ view/share   │  查看/分享报告
└──────────────┘
```

## 📝 下一步

### 🎯 监控多个节点

创建配置文件 `cluster.yaml`：

```yaml
cluster:
  name: my-cluster
  nodes:
    - name: web-01
      host: 192.168.1.10
      user: admin
      ssh_key: ~/.ssh/id_rsa
    - name: web-02
      host: 192.168.1.11
      user: admin
      ssh_key: ~/.ssh/id_rsa
    - name: db-01
      host: 192.168.1.20
      user: admin
      ssh_key: ~/.ssh/id_rsa

collectors:
  - type: system
    enabled: true
  - type: performance
    enabled: true
```

使用配置文件：

```bash
# 采集所有节点数据
clusterreport collect --config cluster.yaml

# 生成集群报告
clusterreport generate --format html --output cluster-report.html
```

### 🔌 启用更多插件

```yaml
collectors:
  - type: system        # 系统配置（NodeProbe）
    enabled: true
  - type: performance   # 性能数据（PerfSnap）
    enabled: true
  - type: mysql         # MySQL 监控
    enabled: true
    config:
      host: localhost
      port: 3306
  - type: redis         # Redis 监控
    enabled: true
    config:
      host: localhost
      port: 6379
```

### 📊 生成不同格式的报告

```bash
# HTML 报告（交互式）
clusterreport generate --format html --output report.html

# PDF 报告（可打印）
clusterreport generate --format pdf --output report.pdf

# Excel 报告（数据分析）
clusterreport generate --format excel --output report.xlsx

# Markdown 报告（可编辑）
clusterreport generate --format markdown --output report.md

# JSON 数据（API 集成）
clusterreport generate --format json --output report.json
```

### 🌐 启动 Web 仪表板

```bash
clusterreport serve --port 8080
```

然后在浏览器访问 `http://localhost:8080`

### 📅 设置定时任务

使用 cron 设置每天生成报告：

```bash
# 编辑 crontab
crontab -e

# 添加以下行（每天凌晨 2 点执行）
0 2 * * * cd /path/to/project && clusterreport collect --config cluster.yaml && clusterreport generate --format html --output daily-report.html
```

## 💡 常见使用场景

### 场景 1: 快速健康检查

```bash
# 一键检查单个服务器
clusterreport collect --node production-server && \
clusterreport generate --format html --output health-check.html && \
open health-check.html
```

### 场景 2: 性能问题排查

```bash
# 启用详细性能分析
clusterreport collect \
  --node problem-server \
  --enable-flamegraph \
  --duration 60s

clusterreport analyze --focus performance
clusterreport generate --format pdf --output perf-analysis.pdf
```

### 场景 3: 容量规划

```bash
# 收集历史数据
clusterreport collect \
  --cluster production \
  --history 30d

clusterreport analyze --type capacity-planning
clusterreport generate --format excel --output capacity-plan.xlsx
```

## 🆘 故障排查

### 问题 1: 连接远程节点失败

**症状**: `Error: Failed to connect to node`

**解决方案**:
```bash
# 检查 SSH 连接
ssh -i ~/.ssh/id_rsa user@host

# 确保 SSH 密钥配置正确
chmod 600 ~/.ssh/id_rsa
```

### 问题 2: 权限不足

**症状**: `Error: Permission denied`

**解决方案**:
```bash
# 使用 sudo 运行（采集系统信息时需要）
sudo clusterreport collect --node localhost

# 或添加用户到必要的组
sudo usermod -aG docker,sudo $USER
```

### 问题 3: 数据采集超时

**症状**: `Error: Timeout collecting data`

**解决方案**:
```bash
# 增加超时时间
clusterreport collect --node remote-server --timeout 300s

# 或禁用耗时的采集器
clusterreport collect --node remote-server --disable performance
```

## 📚 更多资源

- [配置文件详解](configuration.md)
- [用户指南](../user-guide/)
- [CLI 命令参考](../reference/cli-reference.md)
- [示例代码](../../examples/)

## 🤝 需要帮助？

- 📖 查看[完整文档](../README.md)
- 💬 [GitHub Discussions](https://github.com/sunyifei83/devops-toolkit/discussions)
- 🐛 [报告问题](https://github.com/sunyifei83/devops-toolkit/issues)
- 📧 Email: sunyifei83@gmail.com

---

**下一步**: 阅读[配置说明](configuration.md)了解更多配置选项

**上一步**: [安装指南](installation.md)
