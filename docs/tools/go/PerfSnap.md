# PerfSnap - Linux系统性能快照分析工具

## 概述

`PerfSnap` 是一个专业的Linux系统性能快照分析工具，能够快速采集和分析CPU、内存、磁盘IO、网络等关键性能指标。v1.1.0版本新增了火焰图生成功能，可以自动或手动指定进程生成CPU火焰图，帮助深入分析性能瓶颈。

**版本**: 1.1.0  
**作者**: sunyifei83@gmail.com  
**项目**: https://github.com/sunyifei83/devops-toolkit  
**更新日期**: 2025-09-15

## 核心特性

- 🚀 **全面性能分析**: 一键获取系统性能全貌
- 📊 **多维度指标**: CPU、内存、磁盘、网络、进程全覆盖
- ⚡ **快速诊断**: 自动识别性能问题并提供优化建议
- 📈 **实时监控**: 支持持续监控模式，实时刷新数据
- 🎯 **智能分析**: 自动检测异常指标并预警
- 🔥 **火焰图生成**: 支持CPU火焰图生成，深入分析热点函数
- 🤖 **自动化**: 可自动选择CPU占用最高的进程生成火焰图
- 🌐 **可视化**: 生成HTML查看器，方便浏览火焰图

## 主要功能模块

### 1. 系统概况
- 运行时间和在线用户数
- 系统负载（1/5/15分钟）
- CPU核心数和负载评估

### 2. CPU分析
- 多核CPU使用率详情
- 用户态/系统态/IO等待分析
- 上下文切换和中断统计
- 高CPU进程排行

### 3. 内存监控
- 内存使用率和可用内存
- 缓存和缓冲区统计
- 交换分区使用情况

### 4. 磁盘IO
- 设备级读写IOPS
- 吞吐量统计（KB/s）
- 队列长度和利用率

### 5. 网络流量
- 接口级收发包速率
- 网络吞吐量（KB/s）
- TCP连接状态统计

### 6. 进程分析
- TOP进程（按CPU/内存排序）
- 高CPU进程详情（pidstat）
- 进程级性能指标

### 7. 火焰图生成（v1.1.0新增）
- 自动选择CPU最高进程
- 支持指定PID采样
- 可配置采样时长和频率
- 生成交互式SVG火焰图
- 提供HTML查看器

## 安装部署

### 系统要求
- Linux操作系统（CentOS、Ubuntu、Debian等）
- Go 1.15或更高版本（编译时需要）
- Root权限（推荐，获取完整数据）

### 依赖工具

#### 基础依赖（性能分析）
```bash
# CentOS/RHEL
sudo yum install -y sysstat

# Ubuntu/Debian  
sudo apt-get install -y sysstat

# 验证安装
which sar mpstat pidstat iostat
```

#### 火焰图依赖（可选）
```bash
# 安装perf工具
# CentOS/RHEL
sudo yum install -y perf

# Ubuntu/Debian
sudo apt-get install -y linux-tools-common linux-tools-generic

# 安装git（用于下载FlameGraph）
sudo yum install -y git  # 或 apt-get install -y git

# 安装perl（运行FlameGraph脚本）
sudo yum install -y perl  # 或 apt-get install -y perl
```

注：FlameGraph工具会在首次使用时自动下载安装

### 编译安装

```bash
# 1. 克隆项目
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit

# 2. 编译PerfSnap（注意：需要单独编译）
go build -o perfsnap tools/go/PerfSnap.go

# 3. 设置执行权限
chmod +x perfsnap

# 4. (可选) 移动到系统路径
sudo mv perfsnap /usr/local/bin/

# 5. 验证安装
perfsnap -version
```

## 使用方法

### 基础用法

```bash
# 生成一次性能快照报告
sudo perfsnap

# 实时监控模式（默认2秒刷新，持续60秒）
sudo perfsnap -m

# 自定义监控参数
sudo perfsnap -m -interval 5 -duration 120  # 每5秒刷新，持续120秒

# 查看帮助
perfsnap -h
```

### 火焰图生成

```bash
# 自动选择CPU最高的进程生成火焰图
sudo perfsnap -flame

# 指定进程PID生成火焰图
sudo perfsnap -flame -pid 1234

# 自定义采样参数
sudo perfsnap -flame -flame-duration 60 -flame-freq 99

# 生成火焰图并自动打开
sudo perfsnap -flame -flame-open

# 指定输出目录
sudo perfsnap -flame -flame-output /tmp/flamegraph

# 完整示例：生成性能报告+火焰图
sudo perfsnap -flame -flame-duration 30 -flame-open
```

### 命令行参数

#### 基础参数
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-m, --monitor` | 实时监控模式 | false |
| `-interval` | 监控刷新间隔(秒) | 2 |
| `-duration` | 监控持续时间(秒) | 60 |
| `-h, --help` | 显示帮助信息 | - |
| `-version` | 显示版本信息 | - |

#### 火焰图参数
| 参数 | 说明 | 默认值 |
|------|------|--------|
| `-flame` | 生成CPU火焰图 | false |
| `-pid` | 指定进程PID(0=自动选择) | 0 |
| `-flame-duration` | 采样时长(秒) | 30 |
| `-flame-freq` | 采样频率(Hz) | 99 |
| `-flame-output` | 输出目录 | 自动生成 |
| `-flame-open` | 自动打开火焰图 | false |

## 输出示例

### 性能快照报告

```
================================================================================
                    PerfSnap - 系统性能快照分析报告
                    2025-09-15 15:30:00
================================================================================

【系统概况】
运行时间: 5 days | 在线用户: 2
系统负载: 0.52 (1分钟) | 0.48 (5分钟) | 0.45 (15分钟)
CPU核心数: 8

【内存状态】
总内存: 16384 MB | 已使用: 8192 MB (50.0%) | 可用: 8192 MB
缓存: 4096 MB | 缓冲: 512 MB

【CPU统计】
CPU    User%   System%  IOWait%   Idle%   状态
------------------------------------------------------------
all      15.2      5.3      2.1     77.4   正常
0        20.5      8.2      1.5     69.8   正常
1        18.3      6.1      2.8     72.8   正常
...

【磁盘IO统计】
设备         读IOPS   写IOPS    读KB/s     写KB/s   队列   使用率%
----------------------------------------------------------------------
sda            25.3     45.2      512.5     1024.8    1.2    15.5%
sdb            10.1     20.3      256.2      512.4    0.5     8.2%

【网络流量】
接口          接收pps    发送pps    接收KB/s    发送KB/s
------------------------------------------------------------
eth0           1250.5     1180.3       125.5       118.2
eth1            850.2      920.1        85.3        92.5

【TCP连接】
已建立: 1250 | TIME_WAIT: 523 | CLOSE_WAIT: 12

【TOP进程 (按CPU排序)】
PID      用户          CPU%     内存%   命令
------------------------------------------------------------
12345    www          25.3      8.5    java
23456    mysql        18.2     15.3    mysqld
34567    root         12.5      2.1    nginx

【性能分析总结】
------------------------------------------------------------
✅ 系统性能状态良好

【优化建议】
------------------------------------------------------------
• 监控高CPU进程java(PID:12345)的资源使用情况
```

### 火焰图输出

```
================================================================================
                    开始生成火焰图
================================================================================

自动选择CPU使用率最高的进程: PID=12345, Command=java, CPU=25.3%

开始生成火焰图...
目标进程: PID=12345 (java)
采样时长: 30秒
采样频率: 99Hz
输出目录: flamegraph_20250915_153000

[1/4] 正在采集性能数据 (30秒)...
[2/4] 正在解析调用栈...
[3/4] 正在折叠调用栈...
[4/4] 正在生成火焰图...

✅ 火焰图生成成功: flamegraph_20250915_153000/flamegraph_java_12345_20250915_153000.svg
查看器已生成: flamegraph_20250915_153000/viewer_20250915_153000.html

【火焰图分析提示】
------------------------------------------------------------
1. 查看火焰图: 打开 flamegraph_20250915_153000/viewer_20250915_153000.html
2. 宽度越宽的函数，CPU占用时间越长
3. 寻找平顶（plateau）区域，这些可能是性能瓶颈
4. 关注递归调用和深度调用栈
5. 使用浏览器的搜索功能查找特定函数
```
示例：
./perfsnap_v1.1.0 -flame -pid 226594
![alt](/docs/img/ca5521a7059b.png)

## 火焰图分析指南

### 什么是火焰图

火焰图（Flame Graph）是一种性能分析的可视化工具，用于快速识别CPU热点函数。

### 如何阅读火焰图

1. **X轴**：表示采样数量的比例，越宽的函数执行时间越长
2. **Y轴**：表示调用栈深度，从下到上是函数调用关系
3. **颜色**：通常是随机的，仅用于区分不同函数
4. **交互**：
   - 点击某个框可以放大该部分
   - Ctrl+F搜索特定函数
   - 点击Reset按钮重置视图

### 性能问题识别

1. **宽平顶**：寻找宽且平的区域，这些通常是性能瓶颈
2. **高塔**：过深的调用栈可能存在递归或设计问题
3. **分散**：过于分散的小块表示CPU时间分散，没有明显热点
4. **集中**：某个函数占据大部分宽度，说明是主要瓶颈

### 常见模式

- **CPU密集型**：某些计算函数占据大部分宽度
- **锁竞争**：锁相关函数（如mutex_lock）占比较高
- **系统调用**：大量时间花在系统调用上
- **GC压力**：垃圾回收相关函数占比高（Java/Go等）

## 故障排查

### 常见问题

1. **sysstat工具未安装**
   ```bash
   # 错误: sar: command not found
   # 解决:
   sudo yum install -y sysstat  # CentOS/RHEL
   sudo apt-get install -y sysstat  # Ubuntu/Debian
   ```

2. **权限不足**
   ```bash
   # 错误: 部分数据无法获取
   # 解决: 使用sudo运行
   sudo perfsnap
   ```

3. **perf工具未安装**（火焰图功能）
   ```bash
   # 错误: perf工具未安装
   # 解决:
   sudo yum install -y perf  # CentOS/RHEL
   sudo apt-get install -y linux-tools-common linux-tools-generic  # Ubuntu
   ```

4. **kernel.perf_event_paranoid限制**
   ```bash
   # 错误: perf record失败
   # 解决: 调整内核参数
   sudo sysctl -w kernel.perf_event_paranoid=-1
   ```

5. **FlameGraph下载失败**
   ```bash
   # 错误: 无法克隆FlameGraph
   # 解决: 手动下载
   cd ~
   git clone https://github.com/brendangregg/FlameGraph.git
   # 或使用代理
   git config --global http.proxy http://proxy:port
   ```

6. **编译错误**
   ```bash
   # 错误: main redeclared in this block
   # 解决: 单独编译各个工具
   go build -o perfsnap tools/go/PerfSnap.go
   go build -o nodeprobe tools/go/NodeProbe.go
   ```

## 最佳实践

### 性能基线建立
```bash
# 在系统正常时建立基线
sudo perfsnap > baseline_$(date +%Y%m%d).txt

# 定期对比
sudo perfsnap > current.txt
diff baseline_*.txt current.txt
```

### 问题诊断流程
```bash
# 1. 快速概览
sudo perfsnap | head -50

# 2. 如发现CPU问题，生成火焰图
sudo perfsnap -flame -flame-duration 60

# 3. 持续监控
sudo perfsnap -m -interval 1
