# DevOps Toolkit 深度集成重构计划

## 执行方案：方案B - 深度集成

**开始时间**: 2025/10/28  
**预计完成**: 2-3周  
**当前状态**: 🚧 进行中

## 重构目标

将 NodeProbe、PerfSnap 和 ClusterReport 进行深度集成，形成统一的工具链架构。

### 核心原则

1. **保持向后兼容**: 原有的独立工具继续可用
2. **代码复用最大化**: 提取共享库，避免重复
3. **清晰的职责分离**: ClusterReport 作为编排层，NodeProbe/PerfSnap 作为数据源
4. **渐进式重构**: 分阶段实施，每个阶段都可独立运行

## 新目录结构

```
tools/go/
├── pkg/                          # 共享包（可被其他项目导入）
│   ├── nodeprobe/               # NodeProbe 核心库
│   │   ├── collector.go         # 数据采集器
│   │   ├── types.go             # 数据结构定义
│   │   ├── optimizer.go         # 系统优化功能
│   │   └── formatter.go         # 输出格式化
│   ├── perfsnap/                # PerfSnap 核心库
│   │   ├── collector.go         # 性能数据采集
│   │   ├── analyzer.go          # 性能分析
│   │   ├── flamegraph.go        # 火焰图生成
│   │   └── types.go             # 数据结构
│   └── common/                  # 公共库
│       ├── types.go             # 共享类型定义
│       └── errors.go            # 错误定义
├── internal/                     # 内部共享库（不对外暴露）
│   ├── sysinfo/                 # 系统信息采集
│   │   ├── cpu.go               # CPU 信息
│   │   ├── memory.go            # 内存信息
│   │   ├── disk.go              # 磁盘信息
│   │   ├── network.go           # 网络信息
│   │   └── os.go                # 操作系统信息
│   ├── perfmon/                 # 性能监控
│   │   ├── metrics.go           # 性能指标
│   │   ├── process.go           # 进程监控
│   │   └── collector.go         # 采集器基础
│   └── utils/                   # 工具函数
│       ├── exec.go              # 命令执行
│       ├── format.go            # 格式化工具
│       └── ssh.go               # SSH 连接
├── cmd/                          # 命令行入口
│   ├── nodeprobe/               # NodeProbe CLI
│   │   └── main.go              # 薄包装层，调用 pkg/nodeprobe
│   ├── perfsnap/                # PerfSnap CLI
│   │   └── main.go              # 薄包装层，调用 pkg/perfsnap
│   └── clusterreport/           # ClusterReport CLI
│       └── main.go              # 主入口
├── ClusterReport/               # ClusterReport 主体（保持现有结构）
│   ├── cmd/cli/                 # CLI 实现
│   ├── pkg/
│   │   ├── collector/           # 采集器（使用 pkg/nodeprobe 和 pkg/perfsnap）
│   │   ├── analyzer/            # 分析器
│   │   └── generator/           # 报告生成器
│   └── plugins/
│       └── collectors/
│           ├── nodeprobe_collector.go   # NodeProbe 插件
│           └── perfsnap_collector.go    # PerfSnap 插件
├── NodeProbe.go                 # 保留（向后兼容）-> 将迁移到 cmd/nodeprobe
├── PerfSnap.go                  # 保留（向后兼容）-> 将迁移到 cmd/perfsnap
├── go.mod                       # 根模块
└── REFACTORING_PLAN.md         # 本文档
```

## 实施阶段

### ✅ Phase 0: 准备工作 (已完成)

- [x] 创建新的目录结构
- [x] 制定重构计划
- [x] 创建集成架构分析文档

### 🚧 Phase 1: 提取共享代码库 (进行中 - 第1周)

#### 1.1 创建 internal/sysinfo 库

**目标**: 提取系统信息采集的公共代码

**文件**:
- `internal/sysinfo/cpu.go` - CPU 信息采集
- `internal/sysinfo/memory.go` - 内存信息采集  
- `internal/sysinfo/disk.go` - 磁盘信息采集
- `internal/sysinfo/network.go` - 网络信息采集
- `internal/sysinfo/os.go` - 操作系统信息采集

**从以下代码提取**:
- NodeProbe.go 的采集函数
- ClusterReport/pkg/collector/system_collector.go

#### 1.2 创建 internal/utils 工具库

**文件**:
- `internal/utils/exec.go` - 命令执行封装
- `internal/utils/format.go` - 格式化工具
- `internal/utils/ssh.go` - SSH 连接工具

#### 1.3 创建 internal/perfmon 性能监控库

**文件**:
- `internal/perfmon/metrics.go` - 性能指标定义
- `internal/perfmon/collector.go` - 性能数据采集
- `internal/perfmon/process.go` - 进程监控

**从以下代码提取**:
- PerfSnap.go 的采集函数

### Phase 2: 重构 NodeProbe (第1-2周)

#### 2.1 创建 pkg/nodeprobe 库

**文件**:
- `pkg/nodeprobe/types.go` - 数据结构定义
- `pkg/nodeprobe/collector.go` - 使用 internal/sysinfo 实现采集
- `pkg/nodeprobe/optimizer.go` - 系统优化功能
- `pkg/nodeprobe/formatter.go` - 输出格式化（表格、JSON、YAML）

**接口设计**:
```go
package nodeprobe

type Collector struct {
    config Config
}

func New(config Config) *Collector
func (c *Collector) Collect() (*ServerInfo, error)
func (c *Collector) Optimize() error
func (c *Collector) FormatTable(info *ServerInfo) string
func (c *Collector) FormatJSON(info *ServerInfo) ([]byte, error)
func (c *Collector) FormatYAML(info *ServerInfo) ([]byte, error)
```

#### 2.2 创建 cmd/nodeprobe CLI

**文件**: `cmd/nodeprobe/main.go`

这是一个薄包装层，调用 pkg/nodeprobe:
```go
package main

import (
    "github.com/sunyifei83/devops-toolkit/tools/go/pkg/nodeprobe"
)

func main() {
    // 解析命令行参数
    // 调用 nodeprobe 库
    // 输出结果
}
```

#### 2.3 迁移策略

1. 保留原 NodeProbe.go 一段时间（添加废弃警告）
2. 文档更新推荐使用新的二进制
3. 3个月后移除旧文件

### Phase 3: 重构 PerfSnap (第2周)

#### 3.1 创建 pkg/perfsnap 库

**文件**:
- `pkg/perfsnap/types.go` - 数据结构
- `pkg/perfsnap/collector.go` - 性能数据采集（使用 internal/perfmon）
- `pkg/perfsnap/analyzer.go` - 性能分析
- `pkg/perfsnap/flamegraph.go` - 火焰图生成
- `pkg/perfsnap/monitor.go` - 实时监控

**接口设计**:
```go
package perfsnap

type Collector struct {
    config Config
}

func New(config Config) *Collector
func (c *Collector) CollectSnapshot() (*PerformanceData, error)
func (c *Collector) Analyze(data *PerformanceData) *AnalysisResult
func (c *Collector) Monitor(duration, interval int) error
func (c *Collector) GenerateFlameGraph(config FlameGraphConfig) error
```

#### 3.2 创建 cmd/perfsnap CLI

类似 NodeProbe 的处理方式

### Phase 4: ClusterReport 集成 (第2-3周)

#### 4.1 创建 NodeProbe 插件

**文件**: `ClusterReport/plugins/collectors/nodeprobe_collector.go`

```go
package collectors

import (
    "context"
    "github.com/sunyifei83/devops-toolkit/tools/go/pkg/nodeprobe"
)

type NodeProbeCollector struct {
    collector *nodeprobe.Collector
}

func (c *NodeProbeCollector) Collect(ctx context.Context, target string) (interface{}, error) {
    // 如果是远程节点，先 SSH 连接
    // 调用 nodeprobe.Collector.Collect()
    // 返回结果
}
```

#### 4.2 创建 PerfSnap 插件

**文件**: `ClusterReport/plugins/collectors/perfsnap_collector.go`

```go
package collectors

import (
    "context"
    "github.com/sunyifei83/devops-toolkit/tools/go/pkg/perfsnap"
)

type PerfSnapCollector struct {
    collector *perfsnap.Collector
}

func (c *PerfSnapCollector) Collect(ctx context.Context, target string) (interface{}, error) {
    // 调用 perfsnap.Collector.CollectSnapshot()
    // 返回性能数据
}
```

#### 4.3 更新 ClusterReport 配置

**config.yaml 添加**:
```yaml
collectors:
  - type: nodeprobe
    enabled: true
    config:
      optimize: true  # 是否自动优化系统
  
  - type: perfsnap
    enabled: true
    config:
      snapshot: true
      flamegraph: false  # 是否生成火焰图
```

#### 4.4 更新数据流

```
用户执行: clusterreport collect -c config.yaml

ClusterReport:
  1. 读取配置文件
  2. 对每个节点:
     a. 创建 NodeProbeCollector -> 调用 pkg/nodeprobe
     b. 创建 PerfSnapCollector -> 调用 pkg/perfsnap
     c. 创建其他插件 (MySQL, Redis, etc.)
  3. 合并所有数据
  4. 保存到存储层
  
分析阶段:
  1. 读取采集数据
  2. 配置分析器 (基于 NodeProbe 数据)
  3. 性能分析器 (基于 PerfSnap 数据)
  4. 生成分析报告
  
报告生成:
  1. 读取分析结果
  2. 生成 HTML/PDF/Excel/Markdown
```

### Phase 5: 文档更新 (第3周)

#### 5.1 更新 README.md

- 说明新的工具架构
- 更新安装说明
- 更新使用示例

#### 5.2 更新工具文档

- `docs/tools/go/NodeProbe.md` - 更新为使用新架构
- `docs/tools/go/PerfSnap.md` - 更新为使用新架构
- `docs/tools/go/ClusterReport.md` - 添加集成说明

#### 5.3 创建集成指南

- `docs/INTEGRATION_GUIDE.md` - 如何使用集成后的工具链
- `docs/PLUGIN_DEVELOPMENT.md` - 如何开发新插件

#### 5.4 更新 API 文档

- 为 pkg/nodeprobe 和 pkg/perfsnap 添加 GoDoc
- 生成 API 文档

### Phase 6: 测试和验证 (第3周)

#### 6.1 单元测试

- 为所有新包添加单元测试
- 覆盖率目标: 70%+

#### 6.2 集成测试

- 测试 ClusterReport 调用 NodeProbe/PerfSnap
- 测试远程节点采集
- 测试报告生成

#### 6.3 性能测试

- 对比重构前后的性能
- 确保没有性能退化

## 迁移指南

### 对于 NodeProbe 用户

**旧方式**:
```bash
# 编译
go build -o nodeprobe NodeProbe.go

# 运行
./nodeprobe
```

**新方式**:
```bash
# 编译
cd cmd/nodeprobe && go build -o nodeprobe

# 或使用 make
make nodeprobe

# 运行（完全兼容）
./nodeprobe
```

### 对于 PerfSnap 用户

**旧方式**:
```bash
go build -o perfsnap PerfSnap.go
./perfsnap
```

**新方式**:
```bash
cd cmd/perfsnap && go build -o perfsnap
# 或
make perfsnap
./perfsnap
```

### 对于 ClusterReport 用户

**新功能**:
```bash
# 现在可以利用 NodeProbe 和 PerfSnap 的完整功能
clusterreport collect --enable-nodeprobe --enable-perfsnap

# 生成包含详细配置和性能分析的报告
clusterreport generate --format html --output report.html
```

## 好处总结

### 1. 代码复用

- **重复代码减少 60%+**
  - 系统信息采集逻辑统一在 internal/sysinfo
  - 三个工具共享相同的数据采集代码

### 2. 维护性提升

- **单一责任原则**
  - ClusterReport 专注于编排和
