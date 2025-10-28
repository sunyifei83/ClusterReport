# DevOps Toolkit - ClusterReport 平台

[![GitHub stars](https://img.shields.io/github/stars/sunyifei83/devops-toolkit)](https://github.com/sunyifei83/devops-toolkit/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/sunyifei83/devops-toolkit)](https://github.com/sunyifei83/devops-toolkit/network)
[![GitHub issues](https://img.shields.io/github/issues/sunyifei83/devops-toolkit)](https://github.com/sunyifei83/devops-toolkit/issues)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

> 🚀 **企业级集群管理和报告生成平台**  
> 一站式集群健康检查、性能分析和报告生成解决方案

[English](README.md) | [中文](README-zh.md)

## ✨ 核心特性

### 🎯 ClusterReport - 智能集群分析平台

ClusterReport 是一个功能强大的集群管理和报告生成平台，提供：

- **📊 自动化集群报告** - 一键生成全面的集群健康检查报告
- **🔍 深度系统分析** - 内置 NodeProbe 引擎，采集详细的系统配置信息
- **⚡ 实时性能监控** - 内置 PerfSnap 引擎，实时性能数据采集和分析
- **🔌 可扩展插件系统** - 支持自定义采集器、分析器和输出格式
- **📈 可视化仪表板** - 实时监控集群状态的 Web 界面
- **📄 多格式报告** - 支持 HTML、PDF、Excel、Markdown 等多种格式
- **🤖 智能分析引擎** - 自动识别性能瓶颈和配置问题
- **📅 定时任务调度** - 支持定期自动生成报告

## 🚀 快速开始

### 安装

```bash
# 方式 1: 使用安装脚本（推荐）
curl -sSL https://raw.githubusercontent.com/sunyifei83/devops-toolkit/main/scripts/installation/install.sh | bash

# 方式 2: 从源码编译
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit
make install

# 方式 3: 使用 Docker
docker pull sunyifei83/clusterreport:latest
docker run -it sunyifei83/clusterreport:latest
```

### 5分钟快速体验

```bash
# 1. 采集单个节点的数据
clusterreport collect --node localhost

# 2. 生成 HTML 报告
clusterreport generate --format html --output report.html

# 3. 查看报告
open report.html  # macOS
# 或 xdg-open report.html  # Linux
```

### 完整集群分析

```yaml
# config.yaml
cluster:
  name: production
  nodes:
    - name: web-01
      host: 192.168.1.10
      user: admin
    - name: web-02
      host: 192.168.1.11
      user: admin
    - name: db-01
      host: 192.168.1.20
      user: admin

collectors:
  - type: system      # 系统配置（NodeProbe 引擎）
    enabled: true
  - type: performance # 性能数据（PerfSnap 引擎）
    enabled: true
  - type: mysql       # MySQL 监控
    enabled: true
  - type: redis       # Redis 监控
    enabled: true
```

```bash
# 执行集群分析
clusterreport collect --config config.yaml

# 生成报告
clusterreport generate --input cluster-data.json --format html --output cluster-report.html

# 启动 Web 仪表板
clusterreport serve --port 8080
```

## 📚 文档

### 快速开始
- [安装指南](docs/getting-started/installation.md)
- [快速入门](docs/getting-started/quick-start.md)
- [配置说明](docs/getting-started/configuration.md)

### 用户指南
- [集群管理](docs/user-guide/cluster-management.md)
- [数据采集](docs/user-guide/data-collection.md)
- [分析报告](docs/user-guide/analysis-reports.md)
- [Web 仪表板](docs/user-guide/web-dashboard.md)

### 开发者指南
- [架构设计](docs/developer-guide/architecture.md)
- [插件开发](docs/developer-guide/plugin-development.md)
- [API 参考](docs/developer-guide/api-reference.md)
- [贡献指南](docs/developer-guide/contributing.md)

### 参考文档
- [CLI 命令参考](docs/reference/cli-reference.md)
- [配置文件参考](docs/reference/configuration-reference.md)
- [插件参考](docs/reference/plugin-reference.md)

## 🏗️ 架构

```
┌─────────────────────────────────────────────────────────────┐
│                   ClusterReport Platform                     │
│          企业级集群管理和报告生成平台                          │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐   │
│  │ 调度器   │  │ 采集器   │  │ 分析器   │  │ 生成器   │   │
│  │Scheduler │  │Collector │  │ Analyzer │  │Generator │   │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              数据采集引擎 (Built-in)                  │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │                                                       │  │
│  │  📦 NodeProbe Engine    ⚡ PerfSnap Engine          │  │
│  │  • CPU/内存/磁盘配置      • 实时性能指标              │  │
│  │  • 网络接口信息           • 进程监控                  │  │
│  │  • 软件环境检查           • 火焰图生成                │  │
│  │  • 系统优化建议           • 性能瓶颈分析              │  │
│  │                                                       │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              插件系统 (Extensible)                    │  │
│  ├──────────────────────────────────────────────────────┤  │
│  │  🔌 采集器插件    📊 分析器插件    📄 输出插件        │  │
│  │  • MySQL         • 异常检测        • HTML           │  │
│  │  • Redis         • 趋势分析        • PDF            │  │
│  │  • PostgreSQL    • 容量规划        • Excel          │  │
│  │  • Elasticsearch • 自定义分析      • Markdown       │  │
│  │  • 自定义采集器                    • 自定义输出      │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │          Web Dashboard & API                          │  │
│  │  🌐 实时监控   📊 数据可视化   🔗 RESTful API        │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

## 🎯 使用场景

### 场景 1: 日常运维监控
```bash
# 定期生成集群健康报告
clusterreport collect --config prod-cluster.yaml
clusterreport generate --format html --output daily-report.html
clusterreport send --email ops-team@company.com
```

### 场景 2: 性能问题排查
```bash
# 采集性能数据和火焰图
clusterreport collect --node problem-server --enable-flamegraph
clusterreport analyze --focus performance
clusterreport generate --format pdf --output perf-analysis.pdf
```

### 场景 3: 容量规划
```bash
# 收集历史数据并进行趋势分析
clusterreport collect --cluster production --history 30d
clusterreport analyze --type capacity-planning
clusterreport generate --format excel --output capacity-plan.xlsx
```

### 场景 4: 配置审计
```bash
# 检查集群配置合规性
clusterreport collect --cluster all --mode audit
clusterreport analyze --policy security-baseline.yaml
clusterreport generate --format markdown --output audit-report.md
```

## 💡 示例

查看 [examples/](examples/) 目录获取更多示例：

- [基础集群报告](examples/basic-report/) - 最简单的使用示例
- [多集群管理](examples/multi-cluster/) - 管理多个集群
- [自定义插件](examples/custom-plugin/) - 开发自定义采集器
- [CI 集成](examples/ci-integration/) - 集成到 CI/CD 流程

## 🔌 插件生态

ClusterReport 提供丰富的插件系统：

### 内置采集器
- **System** (NodeProbe 引擎) - 系统配置信息
- **Performance** (PerfSnap 引擎) - 性能数据
- **MySQL** - MySQL 数据库监控
- **Redis** - Redis 监控
- **PostgreSQL** - PostgreSQL 监控
- **Elasticsearch** - ES 集群监控

### 内置分析器
- **Anomaly Detection** - 异常检测
- **Trend Analysis** - 趋势分析
- **Capacity Planning** - 容量规划
- **Security Audit** - 安全审计

### 输出格式
- **HTML** - 交互式 Web 报告
- **PDF** - 可打印文档
- **Excel** - 数据分析表格
- **Markdown** - 可编辑文档
- **JSON** - 原始数据

### 开发自定义插件
参见 [插件开发指南](docs/developer-guide/plugin-development.md)

## 🛠️ 辅助工具

除了核心的 ClusterReport 平台，项目还包含：

### 独立工具
- **[DocConverter](tools/go/DocConverter.go)** - 文档格式转换工具
- **Python 工具集** - 日志分析、指标采集等
- **Shell 脚本** - 自动化运维脚本

### 配置模板
- **[Docker](configs/docker/)** - Docker 配置模板
- **[Kubernetes](configs/kubernetes/)** - K8s 部署配置
- **[Nginx](configs/nginx/)** - Nginx 配置示例
- **[Terraform](configs/terraform/)** - 基础设施即代码

### Ansible Playbooks
- **[Setup](playbooks/setup/)** - 环境配置 playbooks
- **[Maintenance](playbooks/maintenance/)** - 维护 playbooks

## 📊 与其他工具对比

| 特性 | ClusterReport | Prometheus + Grafana | Zabbix | Nagios |
|------|--------------|---------------------|--------|--------|
| 一键报告生成 | ✅ | ❌ | ⚠️ 部分 | ❌ |
| 系统配置采集 | ✅ | ❌ | ⚠️ 部分 | ⚠️ 部分 |
| 性能分析 | ✅ | ✅ | ✅ | ⚠️ 部分 |
| 插件系统 | ✅ | ✅ | ✅ | ✅ |
| 学习曲线 | 低 | 中 | 高 | 中 |
| 部署复杂度 | 低 | 中 | 高 | 中 |
| 报告生成 | ✅ 多格式 | ❌ | ⚠️ 简单 | ⚠️ 简单 |

## 🚦 项目状态

- ✅ **核心功能** - 数据采集、分析、报告生成
- ✅ **NodeProbe 引擎** - 系统配置信息采集
- ✅ **PerfSnap 引擎** - 性能数据采集
- 🚧 **Web 仪表板** - 开发中
- 🚧 **插件市场** - 规划中
- 📋 **云原生支持** - 规划中

查看 [ROADMAP](tools/go/ClusterReport/ROADMAP.md) 了解详细开发计划

## 🤝 贡献

我们欢迎所有形式的贡献！

### 如何贡献
1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

详见 [贡献指南](CONTRIBUTING.md)

### 贡献者
<!-- ALL-CONTRIBUTORS-LIST:START -->
感谢所有贡献者！
<!-- ALL-CONTRIBUTORS-LIST:END -->

## 📄 许可证

本项目采用 [MIT License](LICENSE) 开源协议

## 🌟 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=sunyifei83/devops-toolkit&type=Date)](https://star-history.com/#sunyifei83/devops-toolkit&Date)

## 📮 联系方式

- **Issues**: [GitHub Issues](https://github.com/sunyifei83/devops-toolkit/issues)
- **Email**: sunyifei83@gmail.com
- **Twitter**: [@sunyifei83](https://twitter.com/sunyifei83)

## 🙏 致谢

特别感谢：
- 所有贡献者和社区成员
- 开源社区的支持
- SRE 和 DevOps 从业者的反馈

## 📚 相关资源

- [技术文档](docs/)
- [最佳实践](docs/best-practices/)
- [教程](docs/tutorials/)
- [API 文档](docs/reference/api-reference.md)

---

⭐ **如果这个项目对你有帮助，请给我们一个 Star！**

**使用 ClusterReport，让集群管理变得简单！** 🚀
