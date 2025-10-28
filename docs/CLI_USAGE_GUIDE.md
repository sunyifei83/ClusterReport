# ClusterReport CLI 使用指南

**版本**: v0.8.0-dev  
**更新时间**: 2025年10月28日

---

## 概述

ClusterReport CLI 是一个强大的命令行工具，用于收集、分析和生成集群节点的配置和性能报告。

### 核心功能

- 📋 **配置收集** (NodeProbe): 系统信息、硬件配置、网络、服务等
- 📊 **性能收集** (PerfSnap): CPU、内存、磁盘IO、网络流量、进程统计等
- 🔍 **智能分析**: 自动检测性能问题和配置问题
- 📑 **多格式报告**: 支持 JSON、YAML、Table、HTML 等格式
- 🚀 **自动优化**: 可选的系统自动优化功能

---

## 安装

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/sunyifei83/ClusterReport.git
cd ClusterReport

# 编译 CLI 工具
go build -o clusterreport ./cmd/cli

# 安装到系统路径（可选）
sudo cp clusterreport /usr/local/bin/
```

### 验证安装

```bash
clusterreport --version
clusterreport --help
```

---

## 快速开始

### 1. 收集本地节点所有数据

```bash
clusterreport collect --nodes localhost
```

### 2. 仅收集配置信息

```bash
clusterreport collect --nodes localhost --collect-config
```

### 3. 仅收集性能数据

```bash
clusterreport collect --nodes localhost --collect-perf
```

### 4. 保存结果到文件

```bash
clusterreport collect --nodes localhost --output report.json
```

### 5. 使用不同的输出格式

```bash
# JSON 格式（默认）
clusterreport collect --nodes localhost --format json

# YAML 格式
clusterreport collect --nodes localhost --format yaml

# 表格格式（摘要）
clusterreport collect --nodes localhost --format table
```

---

## 命令详解

### collect - 数据收集

收集命令用于从本地或远程节点收集配置和性能数据。

#### 基本语法

```bash
clusterreport collect [flags]
```

#### 标志说明

##### 必需标志

| 标志 | 默认值 | 说明 |
|------|--------|------|
| `--nodes` | localhost | 要收集的节点列表（逗号分隔） |

##### 收集类型标志

| 标志 | 默认值 | 说明 |
|------|--------|------|
| `--collect-all` | true | 收集所有数据（配置+性能） |
| `--collect-config` | false | 仅收集配置信息 (NodeProbe) |
| `--collect-perf` | false | 仅收集性能数据 (PerfSnap) |

##### 高级选项

| 标志 | 默认值 | 说明 |
|------|--------|------|
| `--auto-optimize` | false | 自动执行系统优化（需要 root 权限） |
| `--flame-graph` | false | 生成 CPU 火焰图（需要 perf 工具） |
| `--duration` | 5 | 性能采集持续时间（秒） |

##### 输出选项

| 标志 | 默认值 | 说明 |
|------|--------|------|
| `-o, --output` | - | 输出文件路径（默认输出到标准输出） |
| `-f, --format` | json | 输出格式: json, yaml, table |

##### 全局标志

| 标志 | 默认值 | 说明 |
|------|--------|------|
| `--config` | ./config.yaml | 配置文件路径 |
| `-v, --verbose` | false | 详细输出模式 |
| `-q, --quiet` | false | 静默模式（仅输出错误） |

#### 使用示例

##### 示例 1: 基本收集

```bash
# 收集本地节点所有数据
clusterreport collect --nodes localhost
```

##### 示例 2: 仅配置收集

```bash
# 仅收集配置信息，不收集性能数据
clusterreport collect --nodes localhost --collect-config
```

##### 示例 3: 仅性能收集

```bash
# 仅收集性能数据，不收集配置信息
clusterreport collect --nodes localhost --collect-perf
```

##### 示例 4: 生成火焰图

```bash
# 收集性能数据并生成 CPU 火焰图
clusterreport collect --nodes localhost --collect-perf --flame-graph
```

**注意**: 火焰图生成需要安装 `perf` 工具和 FlameGraph 工具包。

##### 示例 5: 自动优化

```bash
# 收集数据并自动执行系统优化
sudo clusterreport collect --nodes localhost --auto-optimize
```

**警告**: 自动优化会修改系统配置，需要 root 权限。包括：
- CPU 调频策略调整（powersave → performance）
- 时区校准（统一为 Asia/Shanghai）
- 内核模块自动加载（nf_conntrack, br_netfilter）

##### 示例 6: 详细输出

```bash
# 显示详细的收集过程
clusterreport collect --nodes localhost --verbose
```

##### 示例 7: 保存到文件

```bash
# 保存 JSON 格式结果
clusterreport collect --nodes localhost --output report.json

# 保存 YAML 格式结果
clusterreport collect --nodes localhost --format yaml --output report.yaml
```

##### 示例 8: 表格摘要

```bash
# 在终端显示摘要信息
clusterreport collect --nodes localhost --format table
```

##### 示例 9: 调整采集时长

```bash
# 采集 10 秒的性能数据
clusterreport collect --nodes localhost --collect-perf --duration 10
```

---

## 配置收集 (NodeProbe)

### 收集的数据

配置收集会采集以下信息：

#### 1. 系统信息
- 主机名
- 操作系统版本
- 内核版本
- 系统负载
- 时区信息

#### 2. CPU 信息
- CPU 型号
- CPU 核心数
- CPU 运行模式（32位/64位）
- CPU 调频策略（性能模式/省电模式）

#### 3. 内存信息
- 总内存大小
- 内存插槽信息（通过 dmidecode）
  - 插槽位置
  - 内存大小

#### 4. 磁盘信息
- 系统盘信息
- 数据盘列表
- 磁盘总数和数据盘数量

#### 5. 网络信息
- 网络接口列表
- 接口状态（UP/DOWN）
- 接口速率
- IP 地址

#### 6. 环境信息
- Python 版本和路径
- Java 版本和路径

#### 7. 内核模块
- nf_conntrack 状态
- br_netfilter 状态

### 数据格式

#### JSON 输出示例

```json
{
  "metadata": {
    "timestamp": "2025-10-28T16:18:52+08:00",
    "node": "localhost",
    "collect_types": ["config"],
    "duration_seconds": 0.9,
    "version": "0.8.0-dev"
  },
  "nodeprobe": {
    "hostname": "server01",
    "load_average": "0.15 0.23 0.31",
    "timezone": "Asia/Shanghai",
    "os": "Ubuntu 20.04.3 LTS",
    "kernel": "5.4.0-84-generic",
    "cpu": {
      "model": "Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz",
      "cores": 56,
      "run_mode": "32-bit, 64-bit",
      "performance_mode": "最大性能模式 (performance)"
    },
    "memory": {
      "total_gb": 251.78,
      "slots": [
        {"location": "DIMM_A1", "size": "32 GB"},
        {"location": "DIMM_A2", "size": "32 GB"}
      ]
    },
    "disks": {
      "system_disk": "/dev/sda1 50G/200G (25%)",
      "data_disks": ["/dev/sdb 4T", "/dev/sdc 4T"],
      "total_disks": 3,
      "data_disk_num": 2
    },
    "network": [
      {
        "name": "eth0",
        "status": "UP",
        "speed": "10000Mb/s",
        "ip": "192.168.1.100/24"
      }
    ],
    "python": {
      "version": "Python 3.8.10",
      "path": "/usr/bin/python3"
    },
    "java": {
      "version": "Java 11.0.11",
      "path": "/usr/bin/java (JAVA_HOME: /usr/lib/jvm/java-11-openjdk-amd64)"
    },
    "kernel_modules": {
      "nf_conntrack": true,
      "br_netfilter": true,
      "message": "nf_conntrack: 已加载, br_netfilter: 已加载"
    }
  }
}
```

---

## 性能收集 (PerfSnap)

### 收集的数据

性能收集会采集以下实时性能指标：

#### 1. 系统负载
- 运行时间
- 1/5/15 分钟负载平均值

#### 2. 虚拟内存统计 (vmstat)
- 运行队列长度
- 上下文切换次数
- 中断次数

#### 3. CPU 统计 (mpstat)
- 每个 CPU 核心的使用率
- 用户态/系统态/空闲比例
- IO 等待时间

#### 4. 进程统计 (pidstat)
- 进程 CPU 使用率
- 进程优先级
- 进程命令

#### 5. 磁盘 IO (iostat)
- 读写速率 (MB/s)
- IOPS
- 设备利用率
- 平均等待时间

#### 6. 内存使用 (free)
- 总内存
- 已用内存
- 可用内存
- 缓存和缓冲区

#### 7. 网络统计 (sar -n DEV)
- 接收/发送速率 (KB/s)
- 接收/发送包数 (pps)

#### 8. TCP 连接统计
- 各种状态的连接数
- 新建连接速率
- 连接重传率

#### 9. Top 进程
- CPU 消耗最高的进程
- 内存消耗最高的进程

#### 10. Dmesg 错误
- 系统错误日志

#### 11. 性能问题检测
- 自动检测的性能问题
- 问题严重程度
- 优化建议

### 火焰图生成

使用 `--flame-graph` 标志可以生成 CPU 火焰图：

```bash
clusterreport collect --nodes localhost --collect-perf --flame-graph
```

火焰图会保存为 SVG 文件，可以用浏览器打开查看。

---

## 输出格式

### JSON 格式

标准的 JSON 格式，适合程序处理：

```bash
clusterreport collect --nodes localhost --format json > report.json
```

### YAML 格式

人类可读的 YAML 格式：

```bash
clusterreport collect --nodes localhost --format yaml > report.yaml
```

### Table 格式

在终端显示的摘要表格：

```bash
clusterreport collect --nodes localhost --format table
```

输出示例：

```
📊 收集结果摘要
================================================
节点: localhost
时间: 2025-10-28 16:18:52
耗时: 0.90 秒
收集类型: [config performance]

📋 配置信息摘要:
  主机名: server01
  操作系统: Ubuntu 20.04.3 LTS
  内核版本: 5.4.0-84-generic
  CPU: Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz (56 核心)
  内存: 251.78 GB
  网络接口: 2 个
  磁盘: 3 个 (数据盘: 2)

📊 性能数据摘要:
  主机名: server01
  运行时间: 45 days, 3:22
  负载: 0.15, 0.23, 0.31
  CPU 统计: 56 个
  磁盘 IO: 3 个设备
  网络接口: 2 个
  检测到的问题: 2 个
  优化建议: 5 条

⚠️  检测到的问题:
  [medium] cpu: CPU利用率较高: 85.2%
  [low] memory: 内存使用率: 75.3%

💡 提示: 使用 --format json 或 --format yaml 获取完整数据
```

---

## 故障排除

### 常见问题

#### 1. 权限不足

**问题**: 某些命令需要 root 权限

**解决方案**:
```bash
# 使用 sudo 运行
sudo clusterreport collect --nodes localhost
```

#### 2. 命令未找到

**问题**: `dmidecode: command not found`

**解决方案**:
```bash
# Ubuntu/Debian
sudo apt-get install dmidecode

# CentOS/RHEL
sudo yum install dmidecode
```

#### 3. 工具未安装

**问题**: 性能收集工具未安装

**解决方案**:
```bash
# 安装 sysstat 包（包含 sar, mpstat, iostat, pidstat）
# Ubuntu/Debian
sudo apt-get install sysstat

# CentOS/RHEL
sudo yum install sysstat
```

#### 4. 火焰图生成失败

**问题**: `perf: command not found`

**解决方案**:
```bash
# Ubuntu/Debian
sudo apt-get install linux-tools-common linux-tools-generic

# CentOS/RHEL
sudo yum install perf

# 安装 FlameGraph
git clone https://github.com/brendangregg/FlameGraph.git
export PATH=$PATH:$(pwd)/FlameGraph
```

#### 5. macOS 兼容性

**问题**: 很多 Linux 工具在 macOS 上不可用

**说明**: ClusterReport 主要设计用于 Linux 系统。在 macOS 上运行时，由于缺少 Linux 特定的工具（如 lscpu, dmidecode, sysstat 等），收集的数据会比较有限。

**建议**: 在 Linux 系统上使用 ClusterReport 以获得完整功能。

---

## 最佳实践

### 1. 定期收集

建议定期收集性能数据以跟踪趋势：

```bash
# 每小时收集一次性能数据
0 * * * * /usr/local/bin/clusterreport collect --nodes localhost --collect-perf --output /var/log/perf/$(date +\%Y\%m\%d\%H).json
```

### 2. 收集前检查

在执行自动优化前，先收集当前状态：

```bash
# 收集基线数据
clusterreport collect --nodes localhost --output baseline.json

# 执行优化
sudo clusterreport collect --nodes localhost --auto-optimize --output optimized.json
```

### 3. 组合使用

分别收集配置和性能数据：

```bash
# 收集配置（通常不常变化）
clusterreport collect --nodes localhost --collect-config --output config.json

# 定期收集性能数据
clusterreport collect --nodes localhost --collect-perf --output perf-$(date +%Y%m%d-%H%M).json
```

### 4. 详细诊断

遇到问题时使用详细模式：

```bash
clusterreport collect --nodes localhost --verbose --format table
```

---

## 下一步

- 查看 [配置文件指南](./CONFIG_GUIDE.md) 了解配置管理
- 查看 [API 文档](./API_DOCUMENTATION.md) 了解编程接口
- 参与 [贡献指南](../CONTRIBUTING.md) 一起改进项目

---

## 反馈与支持

- **GitHub Issues**: https://github.com/sunyifei83/ClusterReport/issues
- **GitHub Discussions**: https://github.com/sunyifei83/ClusterReport/discussions
- **Email**: sunyifei83@gmail.com

---

**文档版本**: 1.0  
**最后更新**: 2025年10月28日  
**维护者**: ClusterReport 项目团队
