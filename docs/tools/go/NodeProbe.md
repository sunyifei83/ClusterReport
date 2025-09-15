# NodeProbe - Linux节点配置探测工具

## 概述

`NodeProbe` 是一个专业的Linux服务器节点配置信息收集工具，能够全面探测和采集服务器的硬件配置、系统状态、软件环境等关键信息。该工具不仅能够收集配置信息，还具备自动优化系统设置的能力，支持多种输出格式，是运维人员进行服务器管理、资产清点、故障排查的得力助手。

**版本**: 1.1.1  
**作者**: sunyifei83@gmail.com  
**项目**: https://github.com/sunyifei83/devops-toolkit  
**更新日期**: 2025-09-15

## 核心特性

- 🔍 **全面探测**: 深度采集CPU、内存、磁盘、网络等硬件信息
- 🚀 **自动优化**: 智能识别并自动优化系统性能设置
- 📊 **多格式输出**: 支持表格、JSON、YAML三种输出格式
- ⚡ **快速执行**: 秒级完成全部信息采集
- 🔧 **智能调优**: 自动调整CPU性能模式和时区设置
- 🛡️ **内核检查**: 自动检测并加载必要的内核模块
- 🌍 **UTF-8支持**: 完美处理中文字符，输出格式整齐美观
- 💾 **文件导出**: 支持将结果导出到文件，便于自动化处理
- 🎯 **精准硬件定位**: 内存插槽显示实际物理位置编号

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
- **插槽信息**: 实际硬件插槽位置和容量
  - 显示物理插槽编号（如 DIMM_A1, DIMM_B2）
  - 只显示已安装内存的插槽
  - 过滤空插槽信息
- **智能统计**: 准确统计使用的插槽数量

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
cd devops-toolkit

# 2. 编译NodeProbe（注意：需要单独编译）
go build -o nodeprobe tools/go/NodeProbe.go

# 3. 设置执行权限
chmod +x nodeprobe

# 4. (可选) 移动到系统路径
sudo mv nodeprobe /usr/local/bin/

# 5. 验证安装
nodeprobe -version
```

### 注意事项

## 使用方法

### 基本用法

```bash
# 默认表格格式输出
sudo nodeprobe

# JSON格式输出
sudo nodeprobe -format json

# YAML格式输出
sudo nodeprobe -format yaml

# 输出到文件
sudo nodeprobe -format json -output server_info.json

# 静默模式（减少提示信息）
sudo nodeprobe -quiet

# 查看版本
nodeprobe -version

# 查看帮助
nodeprobe -h
```

### 命令行参数

| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| `-format` | 输出格式 | table | `-format json` |
| `-output` | 输出文件路径 | 无(输出到终端) | `-output report.json` |
| `-quiet` | 静默模式 | false | `-quiet` |
| `-version` | 显示版本信息 | - | `-version` |
| `-h` | 显示帮助信息 | - | `-h` |

### 输出格式说明

#### 1. 表格格式（默认）
最直观的展示方式，适合人工查看：
```bash
sudo nodeprobe
```

#### 2. JSON格式
适合程序处理和API集成：
```bash
# 输出到终端
sudo nodeprobe -format json

# 输出到文件
sudo nodeprobe -format json -output config.json

# 使用jq处理JSON输出
sudo nodeprobe -format json | jq '.cpu'
```

#### 3. YAML格式
适合配置管理和Ansible等工具：
```bash
# 输出到终端
sudo nodeprobe -format yaml

# 输出到文件
sudo nodeprobe -format yaml -output config.yaml
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

### 表格格式输出

```
NodeProbe v1.1.1 - Linux节点配置探测工具
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
║   DIMM_A1:           16384 MB                                ║
║   DIMM_A2:           16384 MB                                ║
║   DIMM_B1:           16384 MB                                ║
║   DIMM_B2:           16384 MB                                ║
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

### JSON格式输出

```json
{
  "hostname": "prod-server-01",
  "load_average": "0.15 0.12 0.10",
  "timezone": "Asia/Shanghai",
  "os": "CentOS Linux 7.9.2009 (Core)",
  "kernel": "3.10.0-1160.71.1.el7.x86_64",
  "cpu": {
    "model": "Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz",
    "cores": 32,
    "run_mode": "32-bit, 64-bit",
    "performance_mode": "最大性能模式 (performance)"
  },
  "memory": {
    "total_gb": 62.79,
    "slots": [
      {"location": "DIMM_A1", "size": "16384 MB"},
      {"location": "DIMM_A2", "size": "16384 MB"},
      {"location": "DIMM_B1", "size": "16384 MB"},
      {"location": "DIMM_B2", "size": "16384 MB"}
    ]
  },
  "disks": {
    "system_disk": "/dev/sda1 45G/200G (23%)",
    "data_disks": [
      "/dev/sdb 2.0T",
      "/dev/sdc 2.0T",
      "/dev/sdd 4.0T",
      "/dev/sde 4.0T"
    ],
    "total_disks": 5,
    "data_disk_num": 4
  },
  "network": [
    {
      "name": "eth0",
      "status": "UP",
      "speed": "1000Mb/s",
      "ip": "192.168.1.100/24"
    },
    {
      "name": "eth1",
      "status": "UP",
      "speed": "10000Mb/s",
      "ip": "10.0.0.50/24"
    }
  ],
  "python": {
    "version": "Python 3.6.8",
    "path": "/usr/bin/python3"
  },
  "java": {
    "version": "Java 1.8.0_292",
    "path": "/usr/bin/java (JAVA_HOME: /usr/java/jdk)"
  },
  "kernel_modules": {
    "nf_conntrack": true,
    "br_netfilter": true,
    "message": "nf_conntrack: 已加载, br_netfilter: 已加载"
  },
  "timestamp": "2025-01-15 12:00:00",
  "nodeprobe_version": "1.1.0"
}
```

### YAML格式输出

```yaml
# NodeProbe Configuration Report
# Generated at: 2025-01-15 12:00:00

hostname: prod-server-01
load_average: 0.15 0.12 0.10
timezone: Asia/Shanghai
os: CentOS Linux 7.9.2009 (Core)
kernel: 3.10.0-1160.71.1.el7.x86_64
timestamp: 2025-01-15 12:00:00
nodeprobe_version: 1.1.1

cpu:
  model: Intel(R) Xeon(R) CPU E5-2680 v4 @ 2.40GHz
  cores: 32
  run_mode: 32-bit, 64-bit
  performance_mode: 最大性能模式 (performance)

memory:
  total_gb: 62.79
  slots:
    - location: DIMM_A1
      size: 16384 MB
    - location: DIMM_A2
      size: 16384 MB
    - location: DIMM_B1
      size: 16384 MB
    - location: DIMM_B2
      size: 16384 MB

disks:
  system_disk: /dev/sda1 45G/200G (23%)
  total_disks: 5
  data_disk_num: 4
  data_disks:
    - /dev/sdb 2.0T
    - /dev/sdc 2.0T
    - /dev/sdd 4.0T
    - /dev/sde 4.0T

network:
  - name: eth0
    status: UP
    speed: 1000Mb/s
    ip: 192.168.1.100/24
  - name: eth1
    status: UP
    speed: 10000Mb/s
    ip: 10.0.0.50/24

python:
  version: Python 3.6.8
  path: /usr/bin/python3

java:
  version: Java 1.8.0_292
  path: /usr/bin/java (JAVA_HOME: /usr/java/jdk)

kernel_modules:
  nf_conntrack: true
  br_netfilter: true
  message: nf_conntrack: 已加载, br_netfilter: 已加载
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
| CPU信息 | /proc/cpuinfo, lscpu | 处理器详情 |
| CPU性能模式 | /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor | 调度器模式 |
| 内存总量 | /proc/meminfo | 内存统计 |
| 内存插槽 | dmidecode -t 17 | 需要root，包含物理位置 |
| 磁盘信息 | lsblk、df | 存储设备 |
| 网络接口 | ip addr | 网络配置 |
| 网卡速度 | ethtool | 可选工具 |
| 时区信息 | timedatectl, /etc/timezone | 系统时区 |
| Python信息 | python --version, which python | 多版本检测 |
| Java信息 | java -version, JAVA_HOME | JDK/JRE信息 |
| 内核模块 | lsmod, modprobe | 模块管理 |

## 版本更新历史

### v1.1.1 (2025-09-15)
- 🔧 **优化内存插槽显示**
  - 显示实际硬件插槽位置（DIMM_A1 等）而非递增编号
  - 自动过滤空插槽，只显示已安装内存的插槽
  - 修复插槽数量统计错误问题
  - JSON/YAML 输出包含详细的插槽位置信息
  - 显示实际系统分区部署盘信息


### v1.1.0 (2025-09-15)
- 🎉 **新增多格式输出支持**
  - 支持 JSON 格式输出
  - 支持 YAML 格式输出
  - 添加命令行参数支持（-format, -output, -quiet, -version）
  - 支持输出到文件

### v1.0.2 (2025-09-15)
- 🐛 **修复中文字符对齐问题**
  - 正确处理 UTF-8 多字节字符
  - 优化表格输出格式

### v1.0.0
- 🚀 **初始版本发布**
  - 基础硬件信息采集
  - 自动优化系统设置
  - 表格格式输出

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

3. **内存插槽信息为空或显示错误**
   - 原因：dmidecode未安装或权限不足
   - 解决：
   ```bash
   # CentOS/RHEL
   sudo yum install -y dmidecode
   
   # Ubuntu/Debian
   sudo apt-get install -y dmidecode
   ```
   - 注意：v1.1.1 版本已修复插槽统计错误和空插槽显示问题

4. **网卡速度显示Unknown**
   - 原因：ethtool未安装
   - 解决：
   ```bash
   # CentOS/RHEL
   sudo yum install -y ethtool
   
   # Ubuntu/Debian
   sudo apt-get install -y ethtool
   ```

5. **编译错误：main redeclared in this block**
   - 原因：tools/go 目录下有多个包含 main 函数的文件
   - 解决：单独编译各个工具
   ```bash
   # 不要在 tools/go 目录下执行 go build
   # 正确的方式：
   go build -o nodeprobe tools/go/NodeProbe.go
   ```

## 实际应用场景

### 1. CI/CD集成

在自动化部署流程中集成NodeProbe：

```bash
#!/bin/bash
# 部署前检查服务器配置
nodeprobe -format json -output pre_deploy.json

# 执行部署
./deploy.sh

# 部署后验证
nodeprobe -format json -output post_deploy.json

# 对比配置变化
diff <(jq . pre_deploy.json) <(jq . post_deploy.json)
```

### 2. 监控系统集成

将NodeProbe数据推送到监控系统：

```bash
#!/bin/bash
# 定期收集配置信息并推送到监控系统
nodeprobe -format json | curl -X POST \
  -H "Content-Type: application/json" \
  -d @- \
  http://monitoring.example.com/api/node/config
```

### 3. Ansible集成

在Ansible playbook中使用：

```yaml
- name: 收集节点配置
  shell: nodeprobe -format yaml -output /tmp/node_config.yaml
  
- name: 读取配置信息
  include_vars:
    file: /tmp/node_config.yaml
    name: node_config
    
- name: 根据配置执行任务
  debug:
    msg: "CPU核心数: {{ node_config.cpu.cores }}"
```

### 4. 批量配置收集

批量收集多台服务器配置：

```bash
#!/bin/bash
# batch_collect.sh
SERVERS="server1 server2 server3"
OUTPUT_DIR="configs_$(date +%Y%m%d)"
mkdir -p $OUTPUT_DIR

for server in $SERVERS; do
    echo "收集 $server 配置..."
    ssh root@$server "nodeprobe -format json" > $OUTPUT_DIR/${server}.json
done

# 生成汇总报告
echo "生成汇总报告..."
jq -s '.[0] | {
    servers: [.[] | {
        hostname: .hostname,
        cpu_cores: .cpu.cores,
        memory_gb: .memory.total_gb,
        disk_count: .disks.total_disks
    }]
}' $OUTPUT_DIR/*.json > $OUTPUT_DIR/summary.json
```

### 5. 配置基线管理

建立和维护配置基线：

```bash
#!/bin/bash
# 建立基线
nodeprobe -format json -output baseline.json

# 定期检查配置偏移
nodeprobe -format json | jq -r --slurpfile baseline baseline.json '
    if .cpu.cores != $baseline[0].cpu.cores then
        "⚠️ CPU核心数变化: \($baseline[0].cpu.cores) -> \(.cpu.cores)"
    else empty end,
    if .memory.total_gb != $baseline[0].memory.total_gb then
        "⚠️ 内存容量变化: \($baseline[0].memory.total_gb) -> \(.memory.total_gb)"
    else empty end
'
```

## 扩展开发

### 未来功能规划

1. **远程节点采集**
   ```bash
   nodeprobe --remote host1,host2,host3
   ```

2. **硬件基准测试**
   ```bash
   nodeprobe --benchmark cpu
   nodeprobe --benchmark memory
   nodeprobe --benchmark disk
   ```

3. **配置对比**
   ```bash
   nodeprobe --compare baseline.json
   ```

4. **Web界面**
   - 集中管理多节点
   - 历史数据展示
   - 配置变更追踪

5. **插件系统**
   - 支持自定义采集模块
   - 第三方工具集成

## 与PerfSnap配合使用

NodeProbe和PerfSnap是配套的服务器管理工具，共同构成完整的服务器状态分析解决方案：

| 工具 | 定位 | 数据类型 | 使用场景 | 执行频率 |
|------|------|---------|----------|----------|
| NodeProbe | 配置探测 | 静态信息 | 资产管理、配置审计、环境准备 | 低频（配置变更时） |
| PerfSnap | 性能分析 | 动态数据 | 性能监控、故障诊断、负载分析 | 高频（实时监控） |

### 组合使用示例

#### 1. 新服务器上线检查

```bash
#!/bin/bash
# 新服务器完整检查脚本

echo "========== 服务器上线检查 =========="
echo "时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo "主机: $(hostname)"
echo ""

# Step 1: 收集硬件配置
echo ">>> 1. 检查硬件配置和系统设置..."
sudo nodeprobe > nodeprobe_$(hostname)_$(date +%Y%m%d).txt
echo "配置信息已保存"

# Step 2: 检查性能基线
echo ">>> 2. 建立性能基线..."
sudo perfsnap > perfsnap_baseline_$(hostname)_$(date +%Y%m%d).txt
echo "性能基线已建立"

# Step 3: 简单压力测试后的性能检查
echo ">>> 3. 执行压力测试..."
stress --cpu 4 --timeout 30s 2>/dev/null || echo "跳过压力测试"

echo ">>> 4. 压力测试后性能检查..."
sudo perfsnap > perfsnap_stress_$(hostname)_$(date +%Y%m%d).txt

echo ""
echo "检查完成！生成的报告文件："
ls -lh *.txt | tail -3
```

#### 2. 故障诊断流程

```bash
#!/bin/bash
# 故障诊断组合脚本

REPORT_DIR="/var/log/diagnostics/$(date +%Y%m%d_%H%M%S)"
mkdir -p $REPORT_DIR

echo "开始故障诊断..."

# 1. 先检查配置是否有变更
echo "[1/4] 检查系统配置..."
sudo nodeprobe > $REPORT_DIR/01_nodeprobe.txt

# 2. 获取当前性能快照
echo "[2/4] 获取性能快照..."
sudo perfsnap > $REPORT_DIR/02_perfsnap_current.txt

# 3. 持续监控性能（1分钟）
echo "[3/4] 开始实时监控（60秒）..."
sudo perfsnap -m 2 60 > $REPORT_DIR/03_perfsnap_monitor.txt

# 4. 收集系统日志
echo "[4/4] 收集系统日志..."
tail -n 1000 /var/log/messages > $REPORT_DIR/04_system_logs.txt 2>/dev/null
dmesg -T | tail -n 500 > $REPORT_DIR/05_dmesg.txt

# 生成诊断摘要
cat > $REPORT_DIR/00_summary.txt << EOF
故障诊断报告
生成时间: $(date)
主机名: $(hostname)

文件列表:
- 01_nodeprobe.txt: 系统配置信息
- 02_perfsnap_current.txt: 当前性能状态
- 03_perfsnap_monitor.txt: 60秒性能监控
- 04_system_logs.txt: 系统日志
- 05_dmesg.txt: 内核日志

快速检查项:
$(grep "CPU性能模式" $REPORT_DIR/01_nodeprobe.txt)
$(grep "系统负载" $REPORT_DIR/02_perfsnap_current.txt | head -1)
$(grep "内存使用率" $REPORT_DIR/02_perfsnap_current.txt | head -1)
EOF

echo ""
echo "诊断完成！报告保存在: $REPORT_DIR"
echo "查看摘要: cat $REPORT_DIR/00_summary.txt"
```

#### 3. 日常巡检脚本

```bash
#!/bin/bash
# daily_inspection.sh - 日常巡检脚本

INSPECTION_LOG="/var/log/inspection/$(date +%Y%m%d).log"
mkdir -p $(dirname $INSPECTION_LOG)

{
    echo "====== 日常巡检报告 ======"
    echo "日期: $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    
    # 基础配置检查（每天一次）
    echo "=== 配置信息 ==="
    sudo nodeprobe | grep -E "主机名:|CPU核心数:|总内存:|磁盘数量:|网络接口数:"
    
    echo ""
    echo "=== 性能状态 ==="
    sudo perfsnap | grep -E "系统负载:|CPU使用率:|内存:|磁盘 .* 利用率|TCP连接:"
    
    echo ""
    echo "=== 异常检查 ==="
    # 检查是否有性能问题
    sudo perfsnap | grep "⚠️" || echo "✅ 无性能告警"
    
} | tee $INSPECTION_LOG

# 发送邮件通知（如果配置了邮件）
# mail -s "服务器巡检报告 $(hostname) $(date +%Y%m%d)" admin@example.com < $INSPECTION_LOG
```

#### 4. 性能对比分析

```bash
#!/bin/bash
# 性能变化对比分析

echo "=== 配置与性能对比分析 ==="

# 收集当前状态
TEMP_DIR=$(mktemp -d)
sudo nodeprobe > $TEMP_DIR/config_now.txt
sudo perfsnap > $TEMP_DIR/perf_now.txt

# 与基线对比（假设有基线文件）
BASELINE_DIR="/opt/baseline"

if [ -f "$BASELINE_DIR/nodeprobe_baseline.txt" ]; then
    echo ">>> 配置变更："
    diff -u $BASELINE_DIR/nodeprobe_baseline.txt $TEMP_DIR/config_now.txt | \
        grep "^[+-]" | grep -v "^[+-][+-][+-]" | head -20
else
    echo "未找到配置基线文件"
fi

if [ -f "$BASELINE_DIR/perfsnap_baseline.txt" ]; then
    echo ""
    echo ">>> 性能指标对比："
    # 提取关键指标进行对比
    for metric in "系统负载" "CPU使用率" "内存" "磁盘.*利用率"; do
        echo "- $metric:"
        echo "  基线: $(grep "$metric" $BASELINE_DIR/perfsnap_baseline.txt | head -1)"
        echo "  当前: $(grep "$metric" $TEMP_DIR/perf_now.txt | head -1)"
    done
else
    echo "未找到性能基线文件"
fi

# 清理临时文件
rm -rf $TEMP_DIR
```

#### 5. 批量服务器检查

```bash
#!/bin/bash
# 批量检查多台服务器

SERVERS="server1 server2 server3 server4"
REPORT_DIR="cluster_report_$(date +%Y%m%d_%H%M%S)"
mkdir -p $REPORT_DIR

echo "开始批量检查服务器集群..."

for server in $SERVERS; do
    echo ">>> 检查 $server ..."
    
    # 并行执行配置和性能检查
    ssh root@$server "sudo nodeprobe" > $REPORT_DIR/${server}_nodeprobe.txt 2>&1 &
    ssh root@$server "sudo perfsnap" > $REPORT_DIR/${server}_perfsnap.txt 2>&1 &
done

# 等待所有任务完成
wait

# 生成汇总报告
cat > $REPORT_DIR/00_cluster_summary.md << EOF
# 集群检查报告

生成时间: $(date)

## 服务器列表
$(for s in $SERVERS; do echo "- $s"; done)

## 配置汇总

| 服务器 | CPU核心 | 内存 | 磁盘数 | 状态 |
|--------|---------|------|--------|------|
$(for server in $SERVERS; do
    if [ -f "$REPORT_DIR/${server}_nodeprobe.txt" ]; then
        cpu=$(grep "CPU核心数:" $REPORT_DIR/${server}_nodeprobe.txt | awk '{print $2}')
        mem=$(grep "总内存:" $REPORT_DIR/${server}_nodeprobe.txt | awk '{print $2, $3}')
        disk=$(grep "磁盘数量:" $REPORT_DIR/${server}_nodeprobe.txt | awk -F'总计:' '{print $2}' | awk '{print $1}')
        echo "| $server | $cpu | $mem | $disk | ✅ |"
    else
        echo "| $server | - | - | - | ❌ |"
    fi
done)

## 性能状态

| 服务器 | 负载 | CPU使用 | 内存使用 | 告警 |
|--------|------|---------|----------|------|
$(for server in $SERVERS; do
    if [ -f "$REPORT_DIR/${server}_perfsnap.txt" ]; then
        load=$(grep "系统负载:" $REPORT_DIR/${server}_perfsnap.txt | head -1 | awk -F': ' '{print $2}' | awk '{print $1}')
        cpu=$(grep "CPU使用率:" $REPORT_DIR/${server}_perfsnap.txt | grep -oE '[0-9]+%' | head -1)
        mem=$(grep "内存:" $REPORT_DIR/${server}_perfsnap.txt | grep -oE '[0-9.]+%' | head -1)
        alerts=$(grep -c "⚠️" $REPORT_DIR/${server}_perfsnap.txt)
        echo "| $server | $load | $cpu | $mem | $alerts |"
    else
        echo "| $server | - | - | - | - |"
    fi
done)

## 详细报告
$(for server in $SERVERS; do
    echo "- [$server NodeProbe]($REPORT_DIR/${server}_nodeprobe.txt)"
    echo "- [$server PerfSnap]($REPORT_DIR/${server}_perfsnap.txt)"
done)
EOF

echo ""
echo "批量检查完成！"
echo "查看汇总报告: cat $REPORT_DIR/00_cluster_summary.md"
```

#### 6. 自动化运维集成

```bash
#!/bin/bash
# 集成到自动化运维流程

# 添加到crontab的定时任务
cat << 'EOF' > /etc/cron.d/server-inspection
# 每天早上9点执行配置检查
0 9 * * * root /usr/local/bin/nodeprobe > /var/log/nodeprobe/$(date +\%Y\%m\%d).log 2>&1

# 每小时执行性能快照
0 * * * * root /usr/local/bin/perfsnap > /var/log/perfsnap/$(date +\%Y\%m\%d_\%H).log 2>&1

# 每周一生成周报
0 10 * * 1 root /usr/local/bin/weekly_report.sh

# 性能告警检查（每5分钟）
*/5 * * * * root /usr/local/bin/perfsnap | grep -q "⚠️" && /usr/local/bin/send_alert.sh
EOF

# weekly_report.sh - 周报生成脚本
cat << 'EOF' > /usr/local/bin/weekly_report.sh
#!/bin/bash
REPORT_FILE="/var/reports/weekly_$(date +%Y%W).html"

{
    echo "<html><head><title>服务器周报</title></head><body>"
    echo "<h1>服务器运行周报</h1>"
    echo "<p>生成时间: $(date)</p>"
    
    echo "<h2>1. 配置信息</h2>"
    echo "<pre>"
    sudo nodeprobe
    echo "</pre>"
    
    echo "<h2>2. 本周性能趋势</h2>"
    echo "<pre>"
    for day in $(seq 0 6); do
        date -d "$day days ago" +%Y%m%d
        grep "系统负载\|CPU使用率\|内存" /var/log/perfsnap/$(date -d "$day days ago" +%Y%m%d)*.log | head -3
        echo "---"
    done
    echo "</pre>"
    
    echo "<h2>3. 告警统计</h2>"
    echo "<pre>"
    grep -h "⚠️" /var/log/perfsnap/*.log | sort | uniq -c | sort -rn
    echo "</pre>"
    
    echo "</body></html>"
} > $REPORT_FILE

# 发送周报
# mail -s "$(hostname) 服务器周报" -a $REPORT_FILE admin@example.com < /dev/null
EOF

chmod +x /usr/local/bin/weekly_report.sh
```

### 使用场景对照表

| 场景 | NodeProbe使用 | PerfSnap使用 | 组合价值 |
|------|---------------|--------------|----------|
| **新服务器验收** | ✅ 验证硬件配置是否符合采购要求 | ✅ 建立性能基线 | 完整的验收报告 |
| **故障诊断** | ✅ 检查配置是否被修改 | ✅ 定位性能瓶颈 | 快速定位问题根源 |
| **
