# NodeProbe - Linux节点配置探测工具

## 概述

`NodeProbe` 是一个专业的Linux服务器节点配置信息收集工具，能够全面探测和采集服务器的硬件配置、系统状态、软件环境等关键信息。该工具不仅能够收集配置信息，还具备自动优化系统设置的能力，是运维人员进行服务器管理、资产清点、故障排查的得力助手。

**版本**: 1.0.0  
**作者**: sunyifei@qiniu.com  
**项目**: https://github.com/sunyifei83/devops-toolkit

## 核心特性

- 🔍 **全面探测**: 深度采集CPU、内存、磁盘、网络等硬件信息
- 🚀 **自动优化**: 智能识别并自动优化系统性能设置
- 📊 **清晰展示**: 格式化输出，信息层次分明
- ⚡ **快速执行**: 秒级完成全部信息采集
- 🔧 **智能调优**: 自动调整CPU性能模式和时区设置
- 🛡️ **内核检查**: 自动检测并加载必要的内核模块

## 主要功能模块

### 1. 系统基础信息
- **主机名**: 服务器标识
- **系统负载**: 1/5/15分钟平均负载
- **时区设置**: 自动校准至Asia/Shanghai
- **操作系统**: 发行版和版本信息
- **内核版本**: Linux内核版本

### 2. CPU信息采集
- **CPU型号**: 处理器具体型号
- **核心数量**: 物理和逻辑核心数
- **运行模式**: 32位/64位支持
- **性能模式**: 自动优化至最大性能
  - powersave → performance 自动切换
  - 支持多种调度器模式

### 3. 内存配置
- **总容量**: 系统总内存大小
- **插槽信息**: 内存条数量和容量
- **内存类型**: DDR3/DDR4等规格

### 4. 存储系统
- **系统盘**: 根分区使用情况
- **数据盘**: 自动识别大容量数据盘
- **磁盘统计**: 总磁盘数和数据盘数量

### 5. 网络配置
- **接口列表**: 所有网络接口
- **接口状态**: UP/DOWN状态
- **传输速率**: 网卡速度（需要ethtool）
- **IP地址**: 各接口IP配置

### 6. 软件环境
- **Python**: 版本和安装路径
- **Java**: JDK/JRE版本和JAVA_HOME
- **内核模块**: nf_conntrack、br_netfilter等

### 7. 自动优化功能
- **CPU性能优化**: 自动切换至performance模式
- **时区校准**: 自动设置为Asia/Shanghai
- **内核模块加载**: 自动加载必要模块

## 安装部署

### 系统要求
- Linux操作系统（CentOS、Ubuntu、Debian等）
- Go 1.15或更高版本（编译时需要）

### 依赖工具
```bash
# 基础工具（通常已预装）
- /proc文件系统
- ip命令（iproute2包）
- lsblk命令

# 可选工具（增强功能）
- dmidecode  # 内存详细信息
- ethtool    # 网卡速度信息
- lscpu      # CPU详细信息
```

### 编译安装

```bash
# 1. 克隆项目
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit/tools/go

# 2. 编译NodeProbe
go build -o nodeprobe NodeProbe.go

# 3. 设置执行权限
chmod +x nodeprobe

# 4. (可选) 移动到系统路径
sudo mv nodeprobe /usr/local/bin/

# 5. 验证安装
nodeprobe
```

## 使用方法

### 基本用法

```bash
# 普通用户运行（部分功能受限）
./nodeprobe

# 推荐：使用root权限运行（完整功能）
sudo nodeprobe
```

### 权限说明

#### 普通用户权限
- ✅ 基础系统信息
- ✅ CPU基本信息
- ✅ 内存总量
- ✅ 磁盘列表
- ✅ 网络接口信息
- ❌ 内存插槽详情（需要dmidecode）
- ❌ CPU性能模式调整
- ❌ 时区自动校准
- ❌ 内核模块加载

#### Root权限
- ✅ 所有信息完整采集
- ✅ 自动优化CPU性能模式
- ✅ 自动校准系统时区
- ✅ 自动加载内核模块
- ✅ 内存硬件详细信息

## 输出示例

```
NodeProbe v1.0.0 - Linux节点配置探测工具
==================================================================

正在探测节点配置信息...

╔════════════════════════════════════════════════════════════════╗
║ 主机名:              prod-server-01                           ║
║ 系统负载:            0.15 0.12 0.10                          ║
║ 时区:                已校准至 Asia/Shanghai (原: UTC)         ║
╠════════════════════════════════════════════════════════════════╣
║ 操作系统:            CentOS Linux 7.9.2009 (Core)            ║
║ 内核版本:            3.10.0-1160.71.1.el7.x86_64             ║
╠════════════════════════════════════════════════════════════════╣
║ CPU型号:             Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz║
║ CPU核心数:           32                                       ║
║ CPU运行模式:         32-bit, 64-bit                          ║
║ CPU性能模式:         已自动调整至最大性能模式 (原: powersave)  ║
╠════════════════════════════════════════════════════════════════╣
║ 总内存:              62.79 GB                                 ║
║ 内存插槽:            4个插槽已使用                            ║
║   插槽1: 16384 MB                                            ║
║   插槽2: 16384 MB                                            ║
║   插槽3: 16384 MB                                            ║
║   插槽4: 16384 MB                                            ║
╠════════════════════════════════════════════════════════════════╣
║ 系统盘:              /dev/sda1 45G/200G (23%)                ║
║ 磁盘数量:            总计: 5, 数据盘: 4                       ║
║ 数据盘:              /dev/sdb 2.0T                           ║
║                      /dev/sdc 2.0T                           ║
║                      /dev/sdd 4.0T                           ║
║                      /dev/sde 4.0T                           ║
╠════════════════════════════════════════════════════════════════╣
║ 网络接口数:          3                                        ║
║   eth0 [UP] 1000Mb/s 192.168.1.100/24                        ║
║   eth1 [UP] 10000Mb/s 10.0.0.50/24                          ║
║   docker0 [UP] Unknown 172.17.0.1/16                         ║
╠════════════════════════════════════════════════════════════════╣
║ Python版本:          Python 3.6.8                            ║
║ Python路径:          /usr/bin/python3                        ║
╠════════════════════════════════════════════════════════════════╣
║ Java版本:            Java 1.8.0_292                          ║
║ Java路径:            /usr/bin/java (JAVA_HOME: /usr/java/jdk)║
╠════════════════════════════════════════════════════════════════╣
║ 内核模块状态:        nf_conntrack: 已自动加载,                ║
║                      br_netfilter: 已自动加载                ║
╚════════════════════════════════════════════════════════════════╝

由 NodeProbe 生成 | https://github.com/sunyifei83/devops-toolkit
```

## 自动优化详解

### CPU性能模式优化

NodeProbe会自动检测CPU调度器模式，并在root权限下自动优化：

| 模式 | 说明 | 自动处理 |
|------|------|----------|
| powersave | 省电模式，降低性能 | ✅ 自动切换到performance |
| performance | 最大性能模式 | 保持不变 |
| ondemand | 按需调节 | 保持不变 |
| conservative | 保守调节 | 保持不变 |
| schedutil | 调度器控制 | 保持不变 |

### 时区校准

自动检测并校准系统时区：
- 目标时区：Asia/Shanghai
- 检测方式：timedatectl、/etc/timezone、/etc/localtime
- 自动校准：需要root权限

### 内核模块管理

自动检测并加载重要内核模块：

| 模块 | 用途 | 自动加载 |
|------|------|----------|
| nf_conntrack | 连接跟踪，防火墙必需 | ✅ |
| br_netfilter | 网桥过滤，容器网络必需 | ✅ |

## 信息采集来源

| 信息类型 | 数据来源 | 备注 |
|----------|---------|------|
| 主机名 | /etc/hostname | 系统主机名 |
| 系统负载 | /proc/loadavg | 实时负载 |
| 操作系统 | /etc/os-release | 发行版信息 |
| 内核版本 | /proc/version | 内核信息 |
| CPU信息 | /proc/cpuinfo | 处理器详情 |
| 内存总量 | /proc/meminfo | 内存统计 |
| 内存插槽 | dmidecode -t 17 | 需要root |
| 磁盘信息 | lsblk、df | 存储设备 |
| 网络接口 | ip addr | 网络配置 |
| 网卡速度 | ethtool | 可选工具 |

## 最佳实践

### 服务器基线建立
```bash
# 收集新服务器配置基线
sudo nodeprobe > server_baseline_$(hostname)_$(date +%Y%m%d).txt
```

### 批量节点信息收集
```bash
# 使用脚本批量收集
for host in server1 server2 server3; do
    ssh root@$host 'nodeprobe' > ${host}_info.txt
done
```

### 配置变更对比
```bash
# 收集当前配置
sudo nodeprobe > current_config.txt

# 与基线对比
diff baseline_config.txt current_config.txt
```

### 自动化巡检
```bash
# 添加到crontab定期执行
0 9 * * * /usr/local/bin/nodeprobe > /var/log/nodeprobe/$(date +\%Y\%m\%d).log
```

## 故障排查

### 常见问题

1. **部分信息显示为"Unknown"或"N/A"**
   - 原因：权限不足或缺少相关工具
   - 解决：使用sudo运行，安装dmidecode、ethtool等工具

2. **CPU性能模式无法自动调整**
   - 原因：需要root权限或系统不支持
   - 解决：`sudo nodeprobe`

3. **内存插槽信息为空**
   - 原因：dmidecode未安装或权限不足
   - 解决：
   ```bash
   # CentOS/RHEL
   sudo yum install -y dmidecode
   
   # Ubuntu/Debian
   sudo apt-get install -y dmidecode
   ```

4. **网卡速度显示Unknown**
   - 原因：ethtool未安装
   - 解决：
   ```bash
   # CentOS/RHEL
   sudo yum install -y ethtool
   
   # Ubuntu/Debian
   sudo apt-get install -y ethtool
   ```

## 扩展开发

### 未来功能规划

1. **远程节点采集**
   ```bash
   nodeprobe --remote host1,host2,host3
   ```

2. **多格式输出**
   ```bash
   nodeprobe --format json > config.json
   nodeprobe --format yaml > config.yaml
   ```

3. **硬件基准测试**
   ```bash
   nodeprobe --benchmark cpu
   nodeprobe --benchmark memory
   nodeprobe --benchmark disk
   ```

4. **配置对比**
   ```bash
   nodeprobe --compare baseline.txt
   ```

5. **Web界面**
   - 集中管理多节点
   - 历史数据展示
   - 配置变更追踪

## 与PerfSnap配合使用

NodeProbe和PerfSnap是配套的服务器管理工具：

| 工具 | 定位 | 数据类型 | 使用场景 |
|------|------|---------|----------|
| NodeProbe | 配置探测 | 静态信息 | 资产管理、配置审计 |
| PerfSnap | 性能分析 | 动态数据 | 性能监控、故障诊断 |

### 组合使用示例

```bash
# 完整的服务器检查
echo "=== 配置信息 ===" > server_check.txt
sudo nodeprobe >> server_check.txt
echo -e "\n
