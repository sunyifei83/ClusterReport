# DevOps Toolkit 最佳实践指南

本指南总结了使用 DevOps Toolkit 时的最佳实践和经验，帮助您更高效、更安全地使用工具箱中的各种工具。

## 📚 目录

- [通用最佳实践](#通用最佳实践)
- [工具特定实践](#工具特定实践)
- [安全性建议](#安全性建议)
- [性能优化](#性能优化)
- [自动化建议](#自动化建议)
- [团队协作](#团队协作)

## 🎯 通用最佳实践

### 1. 版本控制

**✅ 推荐做法**：
```bash
# 始终使用特定版本的工具
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit
git checkout v1.0.0  # 使用稳定版本

# 记录使用的工具版本
nodeprobe -version > tool-versions.txt
```

**❌ 避免**：
- 在生产环境使用未测试的最新版本
- 不同环境使用不同版本的工具

### 2. 环境隔离

**✅ 推荐做法**：
```bash
# 为不同环境创建独立配置
/opt/devops-toolkit/
├── production/
│   ├── nodeprobe
│   └── config/
├── staging/
│   ├── nodeprobe
│   └── config/
└── development/
    ├── nodeprobe
    └── config/
```

### 3. 日志管理

**✅ 推荐做法**：
```bash
# 统一的日志目录结构
/var/log/devops-toolkit/
├── nodeprobe/
│   └── 2025-10-27.log
├── perfsnap/
│   └── 2025-10-27.log
└── iotest/
    └── 2025-10-27.log

# 使用 logrotate 管理日志
cat > /etc/logrotate.d/devops-toolkit <<EOF
/var/log/devops-toolkit/*/*.log {
    daily
    rotate 7
    compress
    delaycompress
    missingok
    notifempty
}
EOF
```

### 4. 权限管理

**✅ 推荐做法**：
```bash
# 使用最小权限原则
# 创建专用用户
sudo useradd -r -s /bin/bash devops-toolkit

# 设置正确的文件权限
sudo chown -R devops-toolkit:devops-toolkit /opt/devops-toolkit
sudo chmod 750 /opt/devops-toolkit/tools/shell/*.sh
sudo chmod 644 /opt/devops-toolkit/configs/*

# 对需要 root 权限的工具使用 sudo
sudo -u devops-toolkit nodeprobe
```

## 🛠️ 工具特定实践

### NodeProbe 最佳实践

#### 1. 定期采集基线数据

```bash
#!/bin/bash
# 每月采集一次完整配置信息作为基线

BASELINE_DIR="/var/lib/devops-toolkit/baseline"
DATE=$(date +%Y%m)
mkdir -p $BASELINE_DIR

sudo nodeprobe -format json -output $BASELINE_DIR/baseline_$DATE.json

# 保留最近6个月的基线
find $BASELINE_DIR -name "baseline_*.json" -mtime +180 -delete
```

#### 2. 配置对比

```bash
#!/bin/bash
# 对比当前配置与基线

BASELINE="/var/lib/devops-toolkit/baseline/baseline_$(date +%Y%m).json"
CURRENT="/tmp/current_config.json"

sudo nodeprobe -format json -output $CURRENT

# 使用 jq 对比关键配置
echo "CPU 核心数对比："
echo "基线: $(jq -r '.cpu.cores' $BASELINE)"
echo "当前: $(jq -r '.cpu.cores' $CURRENT)"

echo "内存容量对比："
echo "基线: $(jq -r '.memory.total_gb' $BASELINE)"
echo "当前: $(jq -r '.memory.total_gb' $CURRENT)"
```

#### 3. 批量采集模式

```bash
#!/bin/bash
# 并发采集多节点信息

NODES_FILE="nodes.txt"  # 每行一个节点IP
MAX_PARALLEL=10
OUTPUT_DIR="cluster_configs"

mkdir -p $OUTPUT_DIR

cat $NODES_FILE | xargs -P $MAX_PARALLEL -I {} bash -c '
  NODE={}
  echo "Collecting from $NODE..."
  ssh -o ConnectTimeout=5 $NODE "sudo /opt/devops-toolkit/tools/go/nodeprobe -format json" \
    > '$OUTPUT_DIR'/${NODE}_config.json 2>&1 || echo "Failed: $NODE"
'

echo "采集完成，结果保存在 $OUTPUT_DIR/"
```

### PerfSnap 最佳实践

#### 1. 采集策略

```bash
# 根据不同场景使用不同的采集时长

# 快速检查（30秒）
perfsnap -duration 30 -output quick_check.json

# 常规分析（5分钟）
perfsnap -duration 300 -output normal_analysis.json

# 深度分析（30分钟）
perfsnap -duration 1800 -flamegraph -output deep_analysis.json
```

#### 2. 性能基线建立

```bash
#!/bin/bash
# 在系统空闲时建立性能基线

# 检查系统负载
LOAD=$(uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | cut -d. -f1)

if [ $LOAD -lt 2 ]; then
  echo "系统负载低，开始采集基线..."
  perfsnap -duration 600 -output /var/lib/devops-toolkit/baseline/perf_baseline.json
else
  echo "系统负载过高($LOAD)，跳过基线采集"
fi
```

#### 3. 定期性能报告

```bash
#!/bin/bash
# 每周生成性能趋势报告

WEEK=$(date +%Y-W%V)
REPORT_DIR="/var/log/devops-toolkit/reports/$WEEK"
mkdir -p $REPORT_DIR

# 采集本周数据
perfsnap -duration 300 -output $REPORT_DIR/perf_$(date +%u).json

# 周日生成汇总报告
if [ $(date +%u) -eq 7 ]; then
  echo "生成本周性能报告..."
  # 这里可以添加数据分析和可视化逻辑
fi
```

### iotest 最佳实践

#### 1. 测试前准备

```bash
#!/bin/bash
# 测试前的系统检查

echo "=== IO测试前检查 ==="

# 检查磁盘空间
AVAIL_GB=$(df -BG . | awk 'NR==2 {print $4}' | sed 's/G//')
if [ $AVAIL_GB -lt 3 ]; then
  echo "❌ 可用空间不足3GB，请清理后再测试"
  exit 1
fi

# 检查系统负载
LOAD=$(uptime | awk -F'load average:' '{print $2}' | awk '{print $1}' | cut -d. -f1)
if [ $LOAD -gt 5 ]; then
  echo "⚠️  系统负载较高($LOAD)，测试结果可能不准确"
  read -p "是否继续？(y/n) " -n 1 -r
  echo
  if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
  fi
fi

# 检查是否有其他IO密集操作
if pgrep -x "fio\|dd" > /dev/null; then
  echo "❌ 检测到其他IO测试进程，请等待完成后再测试"
  exit 1
fi

echo "✅ 系统检查通过，可以开始测试"
./iotest.sh
```

#### 2. 结果保存和对比

```bash
#!/bin/bash
# 保存测试结果并与历史数据对比

RESULT_DIR="/var/lib/devops-toolkit/iotest-results"
mkdir -p $RESULT_DIR

DATE=$(date +%Y%m%d_%H%M%S)
RESULT_FILE="$RESULT_DIR/iotest_$DATE.log"

# 运行测试并保存结果
./iotest.sh | tee $RESULT_FILE

# 提取关键指标
echo "提取测试指标..."
grep "Total" $RESULT_FILE > $RESULT_DIR/latest_summary.txt

# 与上次结果对比
if [ -f "$RESULT_DIR/previous_summary.txt" ]; then
  echo "=== 性能对比 ==="
  diff $RESULT_DIR/previous_summary.txt $RESULT_DIR/latest_summary.txt || true
fi

cp $RESULT_DIR/latest_summary.txt $RESULT_DIR/previous_summary.txt
```

### DocConverter 最佳实践

#### 1. 批量转换模板

```bash
#!/bin/bash
# 批量转换文档并保持目录结构

SRC_DIR="./docs"
DST_DIR="./pdf-docs"

find $SRC_DIR -name "*.md" | while read md_file; do
  # 计算相对路径
  rel_path=$(realpath --relative-to=$SRC_DIR "$md_file")
  pdf_file="$DST_DIR/${rel_path%.md}.pdf"
  
  # 创建目标目录
  mkdir -p $(dirname "$pdf_file")
  
  # 转换
  echo "Converting: $md_file"
  docconverter -input "$md_file" -output "$pdf_file"
done

echo "批量转换完成"
```

#### 2. 自动化文档归档

```bash
#!/bin/bash
# 每月归档项目文档

MONTH=$(date +%Y-%m)
ARCHIVE_DIR="/archive/docs/$MONTH"
mkdir -p $ARCHIVE_DIR

# 转换所有文档
for doc in docs/*.md; do
  docconverter -input "$doc" -output "$ARCHIVE_DIR/$(basename ${doc%.md}).pdf"
done

# 创建归档压缩包
tar czf "/archive/docs-$MONTH.tar.gz" -C /archive/docs $MONTH

echo "文档归档完成: docs-$MONTH.tar.gz"
```

## 🔒 安全性建议

### 1. 敏感信息处理

**✅ 推荐做法**：
```bash
# 不要在命令行中直接暴露敏感信息
# 使用配置文件或环境变量

# 配置文件方式
cat > ~/.devops-toolkit/config <<EOF
CLUSTER_SSH_KEY=/path/to/key
CLUSTER_USER=admin
EOF

chmod 600 ~/.devops-toolkit/config

# 在脚本中加载配置
source ~/.devops-toolkit/config

# 环境变量方式
export DEVOPS_TOOLKIT_KEY_PATH=/path/to/key
```

**❌ 避免**：
```bash
# 不要这样做
ssh user@host -i /path/to/private/key  # 密钥路径暴露在进程列表中
nodeprobe --api-token XXXXX            # Token暴露在进程列表中
```

### 2. 输出数据脱敏

```bash
#!/bin/bash
# 对输出进行脱敏处理

# 采集配置信息
nodeprobe -format json -output /tmp/raw_config.json

# 脱敏处理
jq 'del(.network[].ip) | del(.hostname)' /tmp/raw_config.json \
  > /tmp/sanitized_config.json

# 清理原始文件
shred -u /tmp/raw_config.json
```

### 3. 审计日志

```bash
#!/bin/bash
# 记录工具使用审计日志

AUDIT_LOG="/var/log/devops-toolkit/audit.log"

log_audit() {
  echo "$(date -Iseconds) | $(whoami) | $(hostname) | $@" >> $AUDIT_LOG
}

# 在脚本中使用
log_audit "START nodeprobe collection"
sudo nodeprobe -format json -output config.json
log_audit "END nodeprobe collection - status: $?"
```

## ⚡ 性能优化

### 1. 并发处理

```bash
#!/bin/bash
# 使用 GNU Parallel 提高并发效率

# 安装 GNU Parallel
# sudo apt-get install parallel

# 并发采集多节点
parallel -j 20 --timeout 60 \
  'ssh {} "sudo nodeprobe -format json"' \
  > configs/{}.json \
  ::: $(cat nodes.txt)
```

### 2. 结果缓存

```bash
#!/bin/bash
# 缓存不常变化的配置信息

CACHE_DIR="/var/cache/devops-toolkit"
CACHE_TTL=86400  # 24小时

mkdir -p $CACHE_DIR

get_node_config() {
  local node=$1
  local cache_file="$CACHE_DIR/${node}_config.json"
  
  # 检查缓存是否有效
  if [ -f "$cache_file" ]; then
    local cache_age=$(($(date +%s) - $(stat -c %Y "$cache_file")))
    if [ $cache_age -lt $CACHE_TTL ]; then
      echo "使用缓存数据: $node"
      cat "$cache_file"
      return
    fi
  fi
  
  # 缓存失效或不存在，重新采集
  echo "采集新数据: $node"
  ssh $node "sudo nodeprobe -format json" | tee "$cache_file"
}
```

### 3. 资源限制

```bash
#!/bin/bash
# 限制工具资源使用

# 使用 nice 降低进程优先级
nice -n 10 ./iotest.sh

# 使用 cpulimit 限制CPU使用
cpulimit -l 50 -p $(pgrep -f iotest) &

# 使用 cgroup 限制资源（需要 root）
cgcreate -g cpu,memory:/devops-toolkit
cgset -r cpu.shares=512 devops-toolkit
cgset -r memory.limit_in_bytes=1G devops-toolkit
cgexec -g cpu,memory:devops-toolkit ./iotest.sh
```

## 🤖 自动化建议

### 1. CI/CD 集成

```yaml
# .gitlab-ci.yml
