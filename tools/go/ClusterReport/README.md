# ClusterReport - 集群综合报告生成工具

**当前版本**: v0.7.0 (70% 完成) 🚧  
**目标版本**: v1.0.0  
**开发路线图**: 📋 [ROADMAP.md](./ROADMAP.md)  
**项目状态**: 积极开发中

## 概述

ClusterReport 是一个功能强大的集群分析和报告生成工具，能够自动收集、分析集群节点的配置和性能数据，并生成多种格式的综合报告。

### 🎯 核心价值
- **智能化**：自动采集、分析、评分，无需人工干预
- **可视化**：精美的 HTML 报告和实时 Web Dashboard
- **可扩展**：灵活的插件系统，支持自定义采集器和分析器
- **分布式**：Server/Agent 架构，支持大规模集群监控（开发中）

## 特性

- 🔍 **多源数据采集**: 支持 NodeProbe、PerfSnap 等多种数据采集工具
- 📊 **智能分析**: 自动进行配置分析、性能分析和异常检测
- 📈 **多格式报告**: 支持 HTML、PDF、Markdown、Excel 等多种报告格式
- 🔌 **插件系统**: 灵活的插件架构，易于扩展新的采集器和分析器
- ⚡ **并发处理**: 支持并发采集，提高大规模集群的处理效率
- 📅 **定时调度**: 支持定时任务，自动生成周期性报告

## 架构设计

```
ClusterReport/
├── cmd/                    # 命令行入口
│   ├── cli/               # CLI模式
│   ├── server/            # Server模式
│   └── agent/             # Agent模式
├── pkg/                    # 核心包
│   ├── collector/         # 数据采集器
│   ├── analyzer/          # 数据分析器
│   ├── generator/         # 报告生成器
│   ├── storage/           # 存储接口
│   └── scheduler/         # 任务调度器
├── internal/              # 内部实现
│   ├── config/           # 配置管理
│   ├── models/           # 数据模型
│   └── utils/            # 工具函数
├── plugins/               # 插件系统
│   ├── collectors/       # 采集插件
│   ├── analyzers/        # 分析插件
│   └── outputs/          # 输出插件
├── api/                   # API定义
│   ├── rest/            # REST API
│   ├── grpc/            # gRPC API
│   └── graphql/         # GraphQL API
├── web/                   # Web界面
│   ├── dashboard/       # 管理界面
│   └── reports/         # 报告展示
└── deployments/          # 部署配置
```

## 快速开始

### 安装

```bash
# 克隆仓库
git clone https://github.com/devops-toolkit/clusterreport.git
cd clusterreport

# 安装依赖
go mod download

# 构建
go build -o clusterreport cmd/cli/main.go
```

### 基本使用

1. **配置文件**

创建 `config.yaml` 配置文件：

```yaml
clusters:
  - name: production
    nodes:
      - 192.168.1.10
      - 192.168.1.11
      - 192.168.1.12
    username: admin
    ssh_key: ~/.ssh/id_rsa
    port: 22

output:
  directory: ./reports
  formats:
    - html
    - pdf
    - markdown
```

2. **一键生成报告**

```bash
# 生成指定集群的报告
./clusterreport report --cluster production

# 指定输出格式
./clusterreport report --cluster production --formats html,pdf

# 指定输出目录
./clusterreport report --cluster production --output ./my-reports
```

3. **分步执行**

```bash
# 数据采集
./clusterreport collect --cluster production --output data.json

# 数据分析
./clusterreport analyze --input data.json --output analysis.json

# 生成报告
./clusterreport generate --input analysis.json --format html --output report.html
```

## 命令说明

### collect - 数据采集

```bash
clusterreport collect [flags]

Flags:
  -C, --cluster string     集群名称（从配置文件读取）
  -n, --nodes strings      要采集的节点列表
  -o, --output string      输出文件路径
  -p, --parallel int       并发工作数 (默认: 10)
  -t, --timeout duration   采集超时时间 (默认: 5m)
```

### analyze - 数据分析

```bash
clusterreport analyze [flags]

Flags:
  -i, --input string       输入数据文件
  -o, --output string      输出文件路径
  -a, --analyzer string    分析器类型 (默认: default)
  -T, --threshold float    异常阈值 (默认: 0.8)
```

### generate - 报告生成

```bash
clusterreport generate [flags]

Flags:
  -i, --input string       输入分析文件
  -o, --output string      输出文件路径
  -f, --format string      输出格式 (html, pdf, excel, markdown)
  -t, --template string    报告模板 (默认: default)
  -T, --title string       报告标题
```

### report - 一键报告

```bash
clusterreport report [flags]

Flags:
  -C, --cluster string     集群名称
  -n, --nodes strings      节点列表
  -f, --formats strings    输出格式 (默认: [html])
  -o, --output string      输出目录 (默认: ./reports)
  -p, --parallel int       并发工作数 (默认: 10)
```

### schedule - 调度管理

```bash
# 列出调度任务
clusterreport schedule list

# 添加调度任务
clusterreport schedule add --cluster production --cron "0 0 * * *"

# 删除调度任务
clusterreport schedule remove <task-id>
```

### plugin - 插件管理

```bash
# 列出插件
clusterreport plugin list

# 安装插件
clusterreport plugin install <plugin-path>

# 卸载插件
clusterreport plugin uninstall <plugin-name>
```

## 插件开发

### 创建自定义采集器

```go
package collectors

import (
    "context"
    "github.com/devops-toolkit/clusterreport/pkg/collector"
)

type MyCollector struct {
    name string
}

func NewMyCollector() *MyCollector {
    return &MyCollector{
        name: "my-collector",
    }
}

func (c *MyCollector) Name() string {
    return c.name
}

func (c *MyCollector) Collect(ctx context.Context, node collector.Node) (*collector.Data, error) {
    // 实现数据采集逻辑
    data := &collector.Data{
        Node:      node.Name,
        Type:      collector.DataTypeCustom,
        Timestamp: time.Now(),
        Metrics:   make(map[string]interface{}),
    }
    
    // 采集数据...
    
    return data, nil
}

func (c *MyCollector) Validate(config collector.Config) error {
    return nil
}

func (c *MyCollector) SupportedTypes() []collector.DataType {
    return []collector.DataType{collector.DataTypeCustom}
}
```

### 创建自定义分析器

```go
package analyzers

import (
    "context"
    "github.com/devops-toolkit/clusterreport/pkg/analyzer"
    "github.com/devops-toolkit/clusterreport/pkg/collector"
)

type MyAnalyzer struct {
    threshold float64
}

func NewMyAnalyzer() *MyAnalyzer {
    return &MyAnalyzer{
        threshold: 0.8,
    }
}

func (a *MyAnalyzer) Type() analyzer.AnalysisType {
    return "custom"
}

func (a *MyAnalyzer) Analyze(ctx context.Context, data []collector.Data) (*analyzer.Analysis, error) {
    analysis := &analyzer.Analysis{
        Type:      a.Type(),
        Timestamp: time.Now(),
        Results:   make(map[string]interface{}),
    }
    
    // 实现分析逻辑...
    
    return analysis, nil
}

func (a *MyAnalyzer) Options() map[string]interface{} {
    return map[string]interface{}{
        "threshold": a.threshold,
    }
}
```

## 配置文件详解

完整的配置文件示例见 `config.yaml`。主要配置项包括：

- **clusters**: 集群配置，包括节点列表、SSH连接信息等
- **output**: 输出配置，包括目录、格式、模板等
- **schedule**: 调度配置，支持cron表达式
- **storage**: 存储配置，支持文件、数据库等
- **plugins**: 插件配置，包括插件目录、启用列表等
- **collectors**: 采集器配置，包括超时、重试等
- **analyzers**: 分析器配置，包括阈值、窗口等
- **generators**: 生成器配置，包括模板、样式等
- **notifications**: 通知配置，支持邮件、webhook等
- **logging**: 日志配置，包括级别、格式、输出等
- **performance**: 性能配置，包括并发数、超时等
- **security**: 安全配置，包括TLS、认证等

## 报告格式

### HTML 报告
- 美观的可视化界面
- 交互式图表
- 响应式设计
- 支持导出和打印

### PDF 报告
- 专业的文档格式
- 适合存档和分享
- 包含完整的分析结果

### Markdown 报告
- 纯文本格式
- 易于版本控制
- 可以在Git中查看

### Excel 报告
- 结构化数据
- 支持进一步分析
- 包含多个工作表

## 开发状态

**当前版本**: v0.7.0 (70% 完成)  
**详细路线图**: [ROADMAP.md](./ROADMAP.md)

### ✅ 已完成模块 (v0.7.0)

#### 核心功能
- ✅ **采集器框架** (`pkg/collector/`) - 100%
  - 完整的指标数据结构
  - 系统指标采集（CPU、内存、磁盘、网络）
  - 单元测试框架

- ✅ **智能分析器** (`pkg/analyzer/`) - 100%
  - 多维度指标分析
  - 智能健康评分算法（0-100分）
  - 自动问题检测和建议生成

- ✅ **报告生成器** (`pkg/generator/`) - 100%
  - JSON 格式支持
  - HTML 精美报告（含完整CSS）
  - Markdown 文档格式

#### 插件系统
- ✅ 自定义采集器示例
- ✅ MySQL 数据库采集器
- ✅ Redis 缓存采集器
- ✅ 异常检测分析器

#### 前端界面
- ✅ **Web Dashboard** (`web/dashboard/`) - 80%
  - 现代化响应式设计
  - 实时监控视图
  - 交互式数据展示

### 🚧 开发中模块

#### 阶段 1: CLI 模式增强 (v0.8.0) - 优先级⭐⭐⭐⭐⭐
**预计时间**: 3-5 天 | **目标完成度**: 85%

- [ ] 完整的 collect/analyze/generate/report 命令
- [ ] 配置文件管理系统
- [ ] 彩色终端输出和进度条
- [ ] 远程主机采集（SSH）

#### 阶段 2: Server/Agent 架构 (v0.9.0) - 优先级⭐⭐⭐⭐
**预计时间**: 7-10 天 | **目标完成度**: 95%

- [ ] gRPC Server/Agent 通信
- [ ] REST API 实现
- [ ] Agent 管理和心跳检测
- [ ] TLS 安全通信

#### 阶段 3: 存储层 (v0.95.0) - 优先级⭐⭐⭐
**预计时间**: 5-7 天 | **目标完成度**: 98%

- [ ] SQLite 存储实现
- [ ] InfluxDB 支持（可选）
- [ ] 历史数据查询
- [ ] 数据归档和清理

#### 阶段 4: 调度和自动化 (v1.0.0) - 优先级⭐⭐
**预计时间**: 4-6 天 | **目标完成度**: 100%

- [ ] Cron 调度器
- [ ] 自动报告生成
- [ ] 告警系统
- [ ] 多渠道通知（邮件/Webhook/Slack）

### 📅 发布计划

- **v0.8.0** - 2025-11 上旬：CLI 工具功能完整
- **v0.9.0** - 2025-11 中旬：Server/Agent 架构可用
- **v0.95.0** - 2025-11 下旬：数据持久化支持
- **v1.0.0** - 2025-12 上旬：正式版发布 🎉

## 依赖说明

主要依赖库：
- `github.com/spf13/cobra`: CLI框架
- `github.com/spf13/viper`: 配置管理
- `github.com/jung-kurt/gofpdf`: PDF生成
- `github.com/tealeg/xlsx`: Excel生成
- `github.com/hashicorp/go-plugin`: 插件系统
- `github.com/robfig/cron`: 定时调度

## 贡献指南

欢迎贡献代码、报告问题或提出建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件

## 联系方式

- 项目地址: https://github.com/devops-toolkit/clusterreport
- 问题反馈: https://github.com/devops-toolkit/clusterreport/issues

## 鸣谢

感谢所有为本项目做出贡献的开发者！
