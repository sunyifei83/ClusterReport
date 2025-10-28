# ClusterReport - 企业级集群管理和报告生成平台

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

> 🚀 **集群综合分析与报告生成工具**  
> 集成 NodeProbe 和 PerfSnap 的一站式集群监控解决方案

[English](README.md) | [中文](README-zh.md)

## ✨ 核心特性

- **📊 自动化数据采集** - 集群节点配置和性能数据自动采集
- **🔍 深度系统分析** - CPU、内存、磁盘、网络等全方位分析
- **⚡ 实时性能监控** - 集成 PerfSnap 引擎的实时监控能力
- **🔌 可扩展插件系统** - 支持自定义采集器和分析器
- **📄 多格式报告** - HTML、JSON、Markdown 等多种格式
- **🤖 智能分析引擎** - 自动健康评分和问题检测

## 🚀 快速开始

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit

# 编译
go build -o clusterreport ./cmd/cli

# 运行
./clusterreport --help
```

### 基础使用

```bash
# 查看版本
./clusterreport version

# 采集数据（开发中）
./clusterreport collect --node localhost

# 生成报告（开发中）
./clusterreport generate --output report.html
```

## 📁 项目结构

```
devops-toolkit/                     # ClusterReport 项目根目录
├── README.md                       # 项目主页
├── README-zh.md                    # 中文文档
├── LICENSE                         # MIT 许可证
├── ROADMAP.md                      # 开发路线图
├── NEXT_STEPS.md                   # 下一步计划
├── config.yaml                     # 配置文件
├── go.mod / go.sum                 # Go 依赖
│
├── cmd/                            # 命令行入口
│   ├── cli/                       # CLI 模式
│   ├── server/                    # Server 模式（规划中）
│   └── agent/                     # Agent 模式（规划中）
│
├── pkg/                            # 核心包
│   ├── collector/                 # 数据采集器
│   ├── analyzer/                  # 数据分析器
│   └── generator/                 # 报告生成器
│
├── plugins/                        # 插件系统
│   ├── collectors/                # 采集插件（MySQL、Redis等）
│   └── analyzers/                 # 分析插件
│
├── web/                            # Web 界面（开发中）
│   └── dashboard/                 # 管理界面
│
├── deployments/                    # 部署配置
│   ├── docker/                    # Docker 配置
│   ├── kubernetes/                # K8s 配置
│   └── ansible/                   # Ansible 脚本
│
├── docs/                           # 项目文档
│   ├── getting-started/           # 快速开始
│   ├── tools/go/                  # 工具文档
│   └── archive/                   # 归档文档
│
├── tools/                          # 辅助工具
│   └── utils/                     # 工具脚本
│       └── DocConverter.go        # 文档转换
│
├── scripts/                        # 脚本工具
│   └── installation/              # 安装脚本（规划中）
│
└── legacy/                         # 旧版工具
    ├── NodeProbe.go               # 旧版节点探测
    ├── PerfSnap.go                # 旧版性能快照
    └── tools/                     # 旧版工具脚本
```

## 📚 文档

- [快速入门](docs/getting-started/quick-start.md)
- [架构设计](docs/tools/go/ClusterReport_Architecture.md)
- [详细设计](docs/tools/go/ClusterReport_Design.md)
- [开发路线图](ROADMAP.md)
- [下一步计划](NEXT_STEPS.md)
- [旧版工具迁移指南](legacy/README.md)

## 🚦 项目状态

**当前版本**: v0.7.0 (70% 完成度) 🚧  
**目标版本**: v1.0.0  
**预计发布**: 2025年12月

### 已完成 ✅
- ✅ 项目架构设计
- ✅ 核心代码框架
- ✅ 数据采集器（collector）
- ✅ 数据分析器（analyzer）
- ✅ 报告生成器（generator）
- ✅ 插件系统基础
- ✅ Web Dashboard UI

### 开发中 🚧
- 🚧 CLI 命令行工具
- 🚧 配置文件管理
- 🚧 远程节点采集（SSH）
- 🚧 完整的报告格式
- 🚧 测试用例

### 规划中 📋
- 📋 Server/Agent 架构
- 📋 数据持久化存储
- 📋 定时任务调度
- 📋 告警通知系统
- 📋 Docker/K8s 部署

详见 [ROADMAP.md](ROADMAP.md)

## 🏗️ 架构

```
┌─────────────────────────────────────────────────────────────┐
│                   ClusterReport Platform                     │
├─────────────────────────────────────────────────────────────┤
│  📦 数据采集引擎 (Built-in)                                  │
│  • NodeProbe Engine - 系统配置采集                           │
│  • PerfSnap Engine - 性能数据采集                            │
│  • Plugin System - 自定义采集器                              │
├─────────────────────────────────────────────────────────────┤
│  📊 Collector → 📈 Analyzer → 📝 Generator                   │
│  数据采集      智能分析        报告生成                       │
├─────────────────────────────────────────────────────────────┤
│  🔌 插件系统                                                  │
│  • MySQL/Redis/Custom 采集器                                 │
│  • 异常检测/趋势分析                                          │
│  • HTML/PDF/Excel 输出                                       │
└─────────────────────────────────────────────────────────────┘
```

## 🛠️ 技术栈

- **语言**: Go 1.21+
- **配置**: YAML
- **CLI**: Cobra + Viper
- **Web**: HTML + CSS + JavaScript
- **报告**: HTML、JSON、Markdown
- **部署**: Docker、Kubernetes (规划中)

## 🤝 贡献

欢迎贡献代码、报告问题或提出建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/YourFeature`)
3. 提交更改 (`git commit -m 'Add YourFeature'`)
4. 推送到分支 (`git push origin feature/YourFeature`)
5. 开启 Pull Request

详见 [CONTRIBUTING.md](CONTRIBUTING.md)

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源协议

## 📮 联系方式

- **GitHub**: https://github.com/sunyifei83/devops-toolkit
- **Issues**: https://github.com/sunyifei83/devops-toolkit/issues
- **Email**: sunyifei83@gmail.com

## 🎯 设计理念

1. **一体化平台** - 整合 NodeProbe 和 PerfSnap 功能
2. **简单易用** - 一条命令完成复杂任务
3. **高度可扩展** - 插件系统支持自定义扩展
4. **轻量级部署** - 单一二进制文件，无外部依赖
5. **云原生友好** - 容器化部署，K8s 原生支持

---

⭐ **如果这个项目对你有帮助，请给我们一个 Star！**

**注意**: 项目正在积极开发中(70%完成)，部分功能尚未实现。欢迎贡献代码和提出建议！🚀
