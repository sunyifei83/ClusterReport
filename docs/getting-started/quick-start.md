# 快速入门指南

欢迎使用 ClusterReport！本指南将帮助您了解项目的当前状态和基本使用方法。

## ⚠️ 项目状态说明

**ClusterReport 当前处于开发阶段（v0.7.0，70%完成度）**

- ✅ 核心框架已完成
- ✅ 数据采集器、分析器、生成器基础代码已完成
- 🚧 CLI 命令行工具正在完善中
- 🚧 部分功能尚未实现

**预计 v1.0 正式版发布时间**: 2025年12月

## 📋 前置条件

- Go 1.21 或更高版本
- Linux 或 macOS 系统
- Git

## 🚀 快速开始

### 步骤 1: 克隆项目

```bash
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit
```

### 步骤 2: 编译项目

```bash
# 编译 CLI 工具
go build -o clusterreport ./cmd/cli

# 验证编译成功
./clusterreport --help
```

### 步骤 3: 查看项目结构

```bash
# 查看核心组件
ls -la cmd/        # 命令行入口
ls -la pkg/        # 核心包（collector、analyzer、generator）
ls -la plugins/    # 插件系统
ls -la web/        # Web 界面
```

## 📚 理解项目架构

ClusterReport 采用模块化设计：

```
采集 (Collector) → 分析 (Analyzer) → 生成 (Generator)
```

### 核心模块

1. **pkg/collector/** - 数据采集器
   - 系统配置采集（CPU、内存、磁盘、网络）
   - 集成 NodeProbe 引擎
   - 集成 PerfSnap 引擎

2. **pkg/analyzer/** - 数据分析器
   - 多维度指标分析
   - 智能健康评分
   - 问题检测和建议

3. **pkg/generator/** - 报告生成器
   - JSON 格式支持
   - HTML 报告生成
   - Markdown 文档生成

4. **plugins/** - 插件系统
   - MySQL 采集器示例
   - Redis 采集器示例
   - 异常检测分析器

## 🛠️ 当前可用功能

### 已实现 ✅

1. **数据结构定义**
   - 完整的指标数据模型
   - 采集器接口
   - 分析器接口
   - 生成器接口

2. **基础采集功能**
   - 系统指标采集
   - 性能数据采集

3. **分析引擎**
   - 健康评分算法（0-100分）
   - 问题检测逻辑
   - 建议生成

4. **报告生成**
   - JSON 输出
   - HTML 报告（含CSS）
   - Markdown 文档

5. **插件示例**
   - 自定义采集器模板
   - MySQL/Redis 插件示例

### 开发中 🚧

1. **CLI 命令**
   - `collect` 命令
   - `analyze` 命令
   - `generate` 命令
   - 配置文件管理

2. **远程采集**
   - SSH 连接
   - 多节点支持

3. **完整的报告格式**
   - PDF 导出
   - Excel 导出

### 规划中 📋

1. **Server/Agent 架构**
2. **数据持久化**
3. **定时任务调度**
4. **Web 仪表板**
5. **告警系统**

## 📖 查看代码示例

### 示例 1: 查看采集器代码

```bash
# 查看系统采集器实现
cat pkg/collector/system_collector.go

# 查看采集器接口定义
cat pkg/collector/collector.go
```

### 示例 2: 查看分析器代码

```bash
# 查看分析器实现
cat pkg/analyzer/analyzer.go
```

### 示例 3: 查看生成器代码

```bash
# 查看报告生成器
cat pkg/generator/generator.go
```

### 示例 4: 查看插件示例

```bash
# MySQL 采集器插件
cat plugins/collectors/mysql_collector.go

# Redis 采集器插件
cat plugins/collectors/redis_collector.go

# 异常检测分析器
cat plugins/analyzers/anomaly_analyzer.go
```

## 🔄 参与开发

如果您想参与项目开发，可以：

### 1. 查看开发计划

```bash
# 查看路线图
cat ROADMAP.md

# 查看下一步计划
cat NEXT_STEPS.md
```

### 2. 查看待完成任务

根据 [ROADMAP.md](../../ROADMAP.md)，当前优先任务是：

- **阶段 1**: 完善 CLI 模式（v0.8.0）
  - 实现完整的 collect/analyze/generate 命令
  - 配置文件管理
  - 终端输出优化

- **阶段 2**: Server/Agent 架构（v0.9.0）
  - gRPC 通信
  - REST API
  - Agent 管理

### 3. 运行测试

```bash
# 运行单元测试（如果有）
go test ./...

# 查看测试覆盖率
go test -cover ./...
```

## 📚 更多文档

- [ClusterReport 架构](../tools/go/ClusterReport_Architecture.md)
- [ClusterReport 设计](../tools/go/ClusterReport_Design.md)
- [开发路线图](../../ROADMAP.md)
- [下一步计划](../../NEXT_STEPS.md)
- [旧版工具说明](../../legacy/README.md)

## 🤝 需要帮助？

- 📖 查看[项目文档](../README.md)
- 💬 [GitHub Issues](https://github.com/sunyifei83/devops-toolkit/issues)
- 📧 Email: sunyifei83@gmail.com

## ⚠️ 重要提示

**本项目正在积极开发中**：

1. ✅ 核心框架已完成，代码结构清晰
2. 🚧 CLI 工具正在完善，部分命令尚未实现
3. 📋 完整功能预计 2025年12月发布

如果您对项目感兴趣，欢迎：
- ⭐ 给项目加星
- 👀 Watch 项目进展
- 🤝 参与贡献代码
- 💡 提出建议和想法

---

**下一步**: 查看 [ROADMAP.md](../../ROADMAP.md) 了解开发计划

**项目主页**: [README.md](../../README.md)
