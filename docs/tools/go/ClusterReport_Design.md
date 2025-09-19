# ClusterReport - 集群综合报告生成工具

## 概述

`ClusterReport` 是一个专业的集群环境综合报告生成工具，能够汇总NodeProbe、PerfSnap等工具的输出结果，生成包含硬件配置、性能基准、功能验证和性能验证的完整报告文档。

**计划版本**: 1.0.0  
**作者**: sunyifei83@gmail.com  
**项目**: https://github.com/sunyifei83/devops-toolkit  
**状态**: 设计阶段

## 核心功能

### 1. 数据收集与汇总
- **多节点数据收集**: 支持从多个节点收集NodeProbe和PerfSnap数据
- **数据格式支持**: JSON、YAML、文本格式自动识别
- **批量处理**: 支持批量节点数据导入
- **数据验证**: 自动验证数据完整性和一致性

### 2. 报告内容模块

#### 2.1 集群节点硬件环境信息
- **硬件配置汇总表**: CPU、内存、磁盘、网络配置对比
- **配置一致性检查**: 识别配置差异和异常
- **优化建议清单**: 基于最佳实践的优化建议
- **配置基线对比**: 与标准配置基线对比

#### 2.2 物料基准性能
- **性能基准数据**: CPU、内存、磁盘IO、网络性能基准
- **性能对比分析**: 节点间性能对比
- **性能瓶颈识别**: 自动识别性能短板
- **历史趋势分析**: 性能指标变化趋势

#### 2.3 功能验证结果
- **系统功能检查**: 操作系统、内核模块、软件环境
- **服务可用性**: 关键服务运行状态
- **网络连通性**: 节点间网络测试结果
- **依赖项检查**: 软件依赖和版本兼容性

#### 2.4 性能验证结果
- **负载测试结果**: CPU、内存、IO压力测试
- **性能指标达标**: 与SLA对比
- **性能问题诊断**: 问题定位和分析
- **优化建议**: 具体的性能优化方案

### 3. 报告输出格式
- **HTML报告**: 交互式Web报告，包含图表
- **PDF文档**: 专业的PDF格式报告
- **Markdown**: 便于版本控制和在线查看
- **Excel**: 数据导出用于进一步分析

## 使用场景

### 1. 新集群验收
```bash
# 收集所有节点数据
clusterreport collect --nodes node1,node2,node3

# 生成验收报告
clusterreport generate --type acceptance \
  --baseline production.yaml \
  --output cluster_acceptance_report.html
```

### 2. 定期巡检
```bash
# 月度巡检报告
clusterreport generate --type inspection \
  --data-dir /var/log/cluster/202501/ \
  --compare-with /var/log/cluster/202412/ \
  --output monthly_inspection.pdf
```

### 3. 性能基准测试
```bash
# 性能基准报告
clusterreport benchmark --nodes all \
  --tests cpu,memory,disk,network \
  --duration 3600 \
  --output performance_baseline.html
```

### 4. 故障分析报告
```bash
# 故障期间数据分析
clusterreport analyze --type incident \
  --time-range "2025-01-15 14:00:00,2025-01-15 16:00:00" \
  --focus-on performance,logs \
  --output incident_report.pdf
```

## 命令行设计

### 基础命令
```bash
clusterreport [command] [options]

Commands:
  collect     从节点收集数据
  generate    生成报告
  benchmark   执行基准测试
  analyze     分析历史数据
  compare     对比多份报告
  validate    验证配置和性能
```

### 主要参数
```bash
Global Options:
  -c, --config      配置文件路径
  -v, --verbose     详细输出
  -f, --format      输出格式 (html/pdf/md/excel)
  -o, --output      输出文件路径
  --template        报告模板

Collect Options:
  --nodes           节点列表
  --ssh-key         SSH密钥路径
  --parallel        并行收集数
  --timeout         超时时间

Generate Options:
  --type            报告类型 (acceptance/inspection/benchmark/custom)
  --data-dir        数据目录
  --baseline        基线配置文件
  --sla             SLA要求文件
  --include         包含的章节
  --exclude         排除的章节
```

## 配置文件示例

```yaml
# clusterreport.yaml
cluster:
  name: "Production Cluster"
  environment: "production"
  
nodes:
  - name: node1
    host: 192.168.1.10
    role: master
  - name: node2
    host: 192.168.1.11
    role: worker
  - name: node3
    host: 192.168.1.12
    role: worker

baseline:
  cpu:
    cores: 32
    model: "Intel Xeon Gold"
  memory:
    total_gb: 128
  disk:
    system_gb: 500
    data_tb: 2

performance_thresholds:
  cpu_usage: 80
  memory_usage: 85
  disk_usage: 90
  load_average: 2.0

report:
  template: "enterprise"
  logo: "/path/to/logo.png"
  company: "Your Company"
  author: "DevOps Team"
  
data_collection:
  tools:
    nodeprobe:
      enabled: true
      options: "-format json"
    perfsnap:
      enabled: true
      options: "-flame"
  retention_days: 90
```

## 报告模板结构

```markdown
# 集群综合分析报告

## 执行摘要
- 集群概况
- 关键发现
- 风险评估
- 建议措施

## 1. 集群硬件环境
### 1.1 节点配置清单
- 配置对比表
- 配置一致性分析
- 异常配置标识

### 1.2 优化建议
- CPU性能优化
- 内存配置优化
- 存储架构优化
- 网络配置优化

## 2. 物料基准性能
### 2.1 性能基准数据
- CPU性能基准
- 内存性能基准
- 磁盘IO基准
- 网络性能基准

### 2.2 性能对比分析
- 节点间对比
- 与标准基线对比
- 性能瓶颈分析

## 3. 功能验证
### 3.1 系统功能
- 操作系统检查
- 内核模块状态
- 系统服务状态

### 3.2 应用功能
- 应用服务状态
- 依赖项检查
- 接口可用性

## 4. 性能验证
### 4.1 负载测试结果
- CPU负载测试
- 内存压力测试
- IO性能测试
- 网络吞吐测试

### 4.2 性能达标分析
- SLA符合性
- 性能指标评分
- 改进建议

## 5. 问题与建议
### 5.1 发现的问题
- 严重问题
- 警告事项
- 优化机会

### 5.2 改进建议
- 立即行动项
- 短期改进
- 长期规划

## 附录
- A. 详细数据
- B. 测试方法
- C. 参考标准
```

## 技术架构

### 核心组件
```go
// 主要模块
- Collector    // 数据收集器
- Parser       // 数据解析器
- Analyzer     // 数据分析器
- Validator    // 验证器
- Generator    // 报告生成器
- Formatter    // 格式化器
```

### 数据流程
```
1. 数据收集
   ├── SSH远程执行
   ├── 本地文件读取
   └── API调用

2. 数据处理
   ├── 格式转换
   ├── 数据清洗
   └── 数据聚合

3. 数据分析
   ├── 统计分析
   ├── 对比分析
   └── 趋势分析

4. 报告生成
   ├── 模板渲染
   ├── 图表生成
   └── 格式输出
```

## 集成方案

### 与现有工具集成
```bash
#!/bin/bash
# 自动化报告生成脚本

# 1. 收集节点数据
for node in $(cat nodes.txt); do
    ssh $node "nodeprobe -format json" > data/${node}_node.json
    ssh $node "perfsnap" > data/${node}_perf.txt
done

# 2. 生成综合报告
clusterreport generate \
    --data-dir ./data \
    --type acceptance \
    --format html,pdf \
    --output reports/cluster_report_$(date +%Y%m%d)

# 3. 发送报告
mail -s "Cluster Report" -a reports/*.pdf team@example.com < /dev/null
```

### CI/CD集成
```yaml
# .gitlab-ci.yml
cluster-report:
  stage: report
  script:
    - clusterreport collect --nodes $CLUSTER_NODES
    - clusterreport generate --type inspection --output report.html
    - clusterreport validate --sla sla.yaml --fail-on-violation
  artifacts:
    paths:
      - report.html
    reports:
      cluster: report.html
```

## 输出示例

### HTML报告界面
```
┌─────────────────────────────────────────┐
│  Cluster Report - Production Cluster     │
│  Generated: 2025-01-19 12:00:00         │
├─────────────────────────────────────────┤
│  Executive Summary                       │
│  ├─ Total Nodes: 10                     │
│  ├─ Health Score: 92/100                │
│  └─ Critical Issues: 2                  │
├─────────────────────────────────────────┤
│  [Tab1: Hardware] [Tab2: Performance]   │
│  [Tab3: Validation] [Tab4: Issues]      │
├─────────────────────────────────────────┤
│  Node Configuration Comparison           │
│  ┌──────┬──────┬──────┬──────┐         │
│  │Node  │CPU   │Memory│Disk  │         │
│  ├──────┼──────┼──────┼──────┤         │
│  │node1 │32    │128GB │2TB   │         │
│  │node2 │32    │128GB │2TB   │         │
│  │node3 │16⚠️  │64GB⚠️│1TB   │         │
│  └──────┴──────┴──────┴──────┘         │
└─────────────────────────────────────────┘
```

### 命令行输出
```
ClusterReport v1.0.0 - Starting report generation...

[Phase 1/4] Collecting Data
✓ Connected to 10 nodes
✓ NodeProbe data collected from 10/10 nodes
✓ PerfSnap data collected from 10/10 nodes
✓ Benchmark data collected from 10/10 nodes

[Phase 2/4] Analyzing Data
✓ Hardware configuration analyzed
✓ Performance baselines calculated
✓ Anomalies detected: 3 warnings, 1 critical

[Phase 3/4] Validating Requirements
✓ CPU configuration: PASS
⚠ Memory configuration: WARNING (1 node below spec)
✓ Disk configuration: PASS
✓ Network configuration: PASS
✗ Performance SLA: FAIL (2 metrics below threshold)

[Phase 4/4] Generating Report
✓ HTML report generated: cluster_report_20250119.html
✓ PDF report generated: cluster_report_20250119.pdf
✓ Excel data exported: cluster_data_20250119.xlsx

Report Summary:
- Overall Health Score: 85/100
- Configuration Issues: 2
- Performance Issues: 3
- Recommendations: 8

Report generated successfully in 45 seconds.
```

## 开发计划

### Phase 1: MVP (v0.1.0)
- [x] 需求分析和设计
- [ ] 基础数据收集功能
- [ ] 简单报告生成（Markdown）
- [ ] 命令行界面

### Phase 2: 核心功能 (v0.5.0)
- [ ] 完整数据分析
- [ ] HTML报告生成
- [ ] 配置文件支持
- [ ] 批量节点处理

### Phase 3: 高级功能 (v1.0.0)
- [ ] PDF报告生成
- [ ] 图表和可视化
- [ ] 性能基准测试
- [ ] SLA验证
- [ ] 报告模板系统

### Phase 4: 企业特性 (v1.5.0)
- [ ] Web界面
- [ ] 调度任务
- [ ] 报告存档
- [ ] 趋势分析
- [ ] 告警集成

## 相关工具

- **NodeProbe**: 提供硬件配置数据
- **PerfSnap**: 提供性能数据和火焰图
- **DocConverter**: 用于生成PDF报告
- **iotest.sh**: 提供磁盘IO性能测试数据

## 与其他报告工具对比

| 特性 | ClusterReport | 通用监控工具 | 云服务报告 |
|------|--------------|-------------|-----------|
| 硬件配置分析 | ✅ 深度分析 | ⚠️ 基础 | ⚠️ 有限 |
| 性能基准测试 | ✅ 完整 | ✅ 部分 | ⚠️ 基础 |
| 功能验证 | ✅ 全面 | ❌ 无 | ⚠️ 有限 |
| 离线报告 | ✅ 支持 | ⚠️ 部分 | ❌ 需在线 |
| 定制化 | ✅ 高度灵活 | ⚠️ 有限 | ⚠️ 有限 |
| 批量处理 | ✅ 原生支持 | ⚠️ 需脚本 | ✅ 支持 |
| 成本 | ✅ 开源免费 | ⚠️ 部分收费 | ❌ 收费 |

## 未来规划

### 短期目标（3个月）
- 实现基础数据收集和报告生成
- 支持主流Linux发行版
- 提供Docker容器版本

### 中期目标（6个月）
- Web管理界面
- 报告调度和自动化
- 集成更多数据源（Prometheus、Grafana等）

### 长期目标（12个月）
- 机器学习驱动的异常检测
- 智能优化建议
- 多集群管理支持
- SaaS版本

## 贡献指南

欢迎贡献代码和建议！请参考[贡献指南](../../../CONTRIBUTING.md)。

## 许可证

MIT License - 详见[LICENSE](../../../LICENSE)文件

## 联系方式

- **作者**: sunyifei83@gmail.com
- **项目**: https://github.com/sunyifei83/devops-toolkit
- **Issues**: https://github.com/sunyifei83/devops-toolkit/issues
