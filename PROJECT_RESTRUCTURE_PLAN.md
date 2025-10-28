# DevOps Toolkit 项目重构计划 v2.0

**重构日期**: 2025/10/28  
**重构类型**: 🔥 架构级重构 - ClusterReport 作为核心  
**影响范围**: 整个项目结构和文档

## 🎯 重构目标

将 devops-toolkit 从"工具集合"转变为"以 ClusterReport 为核心的集群管理和报告平台"。

### 核心理念转变

**之前**:
```
devops-toolkit = 各种独立工具的集合
├── NodeProbe (独立工具)
├── PerfSnap (独立工具)
├── ClusterReport (独立工具)
└── 其他脚本和配置
```

**现在**:
```
devops-toolkit = ClusterReport 集群管理平台 + 辅助工具
├── ClusterReport (核心平台)
│   ├── 内置 NodeProbe 功能
│   ├── 内置 PerfSnap 功能
│   ├── 插件系统
│   └── Web Dashboard
└── 辅助工具和配置
```

## 📂 新的项目结构

```
devops-toolkit/                              # 项目根目录
├── README.md                                # 项目主文档（重写）
├── README-zh.md                            # 中文文档（重写）
├── LICENSE
├── CONTRIBUTING.md
├── CODE_OF_CONDUCT.md
├── CHANGELOG.md                            # 新增：版本变更日志
│
├── docs/                                    # 文档目录（重组）
│   ├── README.md                           # 文档索引
│   ├── getting-started/                    # 新增：快速开始
│   │   ├── installation.md
│   │   ├── quick-start.md
│   │   └── configuration.md
│   ├── user-guide/                         # 用户指南（重组）
│   │   ├── cluster-management.md          # 集群管理
│   │   ├── data-collection.md             # 数据采集
│   │   ├── analysis-reports.md            # 分析报告
│   │   └── web-dashboard.md               # Web 界面
│   ├── developer-guide/                    # 开发者指南（新增）
│   │   ├── architecture.md                # 架构说明
│   │   ├── plugin-development.md          # 插件开发
│   │   ├── api-reference.md               # API 参考
│   │   └── contributing.md                # 贡献指南
│   ├── reference/                          # 参考文档
│   │   ├── cli-reference.md               # CLI 命令参考
│   │   ├── configuration-reference.md     # 配置参考
│   │   └── plugin-reference.md            # 插件参考
│   ├── best-practices/                     # 最佳实践
│   │   ├── cluster-monitoring.md
│   │   ├── performance-optimization.md
│   │   └── security.md
│   └── tutorials/                          # 教程（新增）
│       ├── basic-cluster-report.md
│       ├── custom-collectors.md
│       └── advanced-analysis.md
│
├── cmd/                                     # 命令行工具入口（新目录）
│   ├── clusterreport/                      # ClusterReport 主程序
│   │   └── main.go
│   ├── nodeprobe/                          # NodeProbe CLI（兼容工具）
│   │   └── main.go
│   └── perfsnap/                           # PerfSnap CLI（兼容工具）
│       └── main.go
│
├── pkg/                                     # 核心库（新目录）
│   ├── clusterreport/                      # ClusterReport 核心
│   │   ├── collector/                      # 数据采集
│   │   ├── analyzer/                       # 数据分析
│   │   ├── generator/                      # 报告生成
│   │   ├── scheduler/                      # 调度器
│   │   └── storage/                        # 存储层
│   ├── nodeprobe/                          # NodeProbe 库
│   │   ├── collector.go
│   │   ├── types.go
│   │   ├── optimizer.go
│   │   └── formatter.go
│   ├── perfsnap/                           # PerfSnap 库
│   │   ├── collector.go
│   │   ├── analyzer.go
│   │   ├── flamegraph.go
│   │   └── types.go
│   └── common/                             # 公共库
│       ├── types.go
│       ├── errors.go
│       └── config.go
│
├── internal/                                # 内部共享库
│   ├── sysinfo/                            # 系统信息采集
│   │   ├── cpu.go
│   │   ├── memory.go
│   │   ├── disk.go
│   │   ├── network.go
│   │   └── os.go
│   ├── perfmon/                            # 性能监控
│   │   ├── metrics.go
│   │   ├── process.go
│   │   └── collector.go
│   └── utils/                              # 工具函数
│       ├── exec.go
│       ├── format.go
│       ├── ssh.go
│       └── logger.go
│
├── plugins/                                 # 插件目录（提升到顶层）
│   ├── README.md                           # 插件开发指南
│   ├── collectors/                         # 采集器插件
│   │   ├── nodeprobe/                     # NodeProbe 插件
│   │   ├── perfsnap/                      # PerfSnap 插件
│   │   ├── mysql/                         # MySQL 插件
│   │   ├── redis/                         # Redis 插件
│   │   ├── postgresql/                    # PostgreSQL 插件
│   │   └── elasticsearch/                 # ES 插件
│   ├── analyzers/                          # 分析器插件
│   │   ├── anomaly/                       # 异常检测
│   │   ├── trend/                         # 趋势分析
│   │   └── capacity/                      # 容量规划
│   └── outputs/                            # 输出插件
│       ├── html/                          # HTML 报告
│       ├── pdf/                           # PDF 报告
│       ├── excel/                         # Excel 报告
│       └── json/                          # JSON 输出
│
├── web/                                     # Web 界面（提升到顶层）
│   ├── dashboard/                          # 仪表板
│   │   ├── index.html
│   │   ├── css/
│   │   ├── js/
│   │   └── assets/
│   ├── reports/                            # 报告展示
│   └── api/                                # Web API
│
├── configs/                                 # 配置文件（保持）
│   ├── clusterreport/                      # ClusterReport 配置
│   │   ├── default.yaml                   # 默认配置
│   │   ├── production.yaml                # 生产环境
│   │   └── examples/                      # 配置示例
│   ├── docker/                            # Docker 配置
│   ├── kubernetes/                        # K8s 配置
│   ├── nginx/                             # Nginx 配置
│   └── terraform/                         # Terraform 配置
│
├── scripts/                                 # 脚本（保持但重组）
│   ├── installation/                       # 安装脚本（新增）
│   │   ├── install.sh
│   │   ├── install-deps.sh
│   │   └── uninstall.sh
│   ├── automation/                         # 自动化脚本
│   ├── monitoring/                         # 监控脚本
│   ├── troubleshooting/                    # 故障排查
│   └── utilities/                          # 工具脚本
│
├── playbooks/                              # Ansible Playbooks
│   ├── setup/                             # 设置 playbooks
│   └── maintenance/                       # 维护 playbooks
│
├── tests/                                   # 测试（扩展）
│   ├── unit/                              # 单元测试
│   ├── integration/                       # 集成测试
│   ├── e2e/                               # 端到端测试
│   └── fixtures/                          # 测试数据
│
├── examples/                                # 示例（新增）
│   ├── basic-report/                      # 基础报告示例
│   ├── multi-cluster/                     # 多集群示例
│   ├── custom-plugin/                     # 自定义插件示例
│   └── ci-integration/                    # CI 集成示例
│
├── deployments/                             # 部署配置（新增）
│   ├── docker/                            # Docker 部署
│   ├── kubernetes/                        # K8s 部署
│   └── ansible/                           # Ansible 部署
│
├── api/                                     # API 定义（新增）
│   ├── rest/                              # REST API
│   ├── grpc/                              # gRPC API
│   └── graphql/                           # GraphQL API
│
├── tools/                                   # 辅助工具（精简）
│   ├── go/                                # Go 工具
│   │   └── DocConverter.go               # 保留独立工具
│   ├── python/                            # Python 工具
│   └── shell/                             # Shell 工具
│
├── Makefile                                # 构建文件（新增）
├── go.mod                                  # Go 模块（顶层）
├── go.sum
├── .github/                                # GitHub 配置
│   ├── workflows/                         # CI/CD
│   │   ├── build.yml
│   │   ├── test.yml
│   │   └── release.yml
│   ├── ISSUE_TEMPLATE/
│   └── PULL_REQUEST_TEMPLATE.md
│
└── legacy/                                  # 旧版本兼容（新增）
    ├── README.md                           # 说明文档
    ├── NodeProbe.go                       # 旧版 NodeProbe
    └── PerfSnap.go                        # 旧版 PerfSnap
```

## 🔄 重构步骤

### Phase 1: 目录结构调整 (第1天)

1. **创建新目录结构**
   ```bash
   mkdir -p cmd/{clusterreport,nodeprobe,perfsnap}
   mkdir -p pkg/{clusterreport,nodeprobe,perfsnap,common}
   mkdir -p internal/{sysinfo,perfmon,utils}
   mkdir -p plugins/{collectors,analyzers,outputs}
   mkdir -p web/{dashboard,reports,api}
   mkdir -p docs/{getting-started,user-guide,developer-guide,reference,best-practices,tutorials}
   mkdir -p examples/{basic-report,multi-cluster,custom-plugin}
   mkdir -p tests/{unit,integration,e2e,fixtures}
   mkdir -p deployments/{docker,kubernetes,ansible}
   mkdir -p api/{rest,grpc,graphql}
   mkdir -p legacy
   ```

2. **移动现有文件**
   ```bash
   # 移动旧工具到 legacy
   mv tools/go/NodeProbe.go legacy/
   mv tools/go/PerfSnap.go legacy/
   
   # 移动 ClusterReport 核心代码到 pkg
   cp -r tools/go/ClusterReport/pkg/* pkg/clusterreport/
   cp -r tools/go/ClusterReport/plugins/* plugins/
   cp -r tools/go/ClusterReport/web/* web/
   
   # 移动配置
   mkdir -p configs/clusterreport
   mv tools/go/ClusterReport/config.yaml configs/clusterreport/default.yaml
   ```

### Phase 2: 文档重写 (第2-3天)

需要重写/更新的文档：

1. **README.md** - 项目主页（全新）
2. **README-zh.md** - 中文文档（全新）
3. **docs/getting-started/** - 快速开始指南
4. **docs/user-guide/** - 用户手册
5. **docs/developer-guide/** - 开发者指南
6. **CHANGELOG.md** - 版本历史

### Phase 3: 代码重构 (第4-10天)

按照之前的深度集成计划实施：
- 提取共享库
- 重构 NodeProbe 为库
- 重构 PerfSnap 为库
- ClusterReport 集成

### Phase 4: 构建和部署 (第11-12天)

1. **Makefile** - 统一构建
2. **Docker镜像** - 容器化
3. **K8s部署** - Kubernetes 支持
4. **CI/CD** - 自动化流程

### Phase 5: 示例和测试 (第13-14天)

1. 创建使用示例
2. 编写测试用例
3. 文档验证

## 📝 关键文档模板

### 新 README.md 结构

```markdown
# DevOps Toolkit - ClusterReport Platform

> 🚀 企业级集群管理和报告平台

## 快速开始

一键部署 ClusterReport...

## 核心功能

- 📊 集群健康检查和报告生成
- 🔍 系统配置信息采集 (内置 NodeProbe)
- ⚡ 性能数据分析 (内置 PerfSnap)
- 🔌 可扩展的插件系统
- 📈 实时监控仪表板
- 📄 多格式报告输出

## 架构

ClusterReport 作为核心平台...

## 文档

- [快速开始](docs/getting-started/quick-start.md)
- [用户指南](docs/user-guide/)
- [插件开发](docs/developer-guide/plugin-development.md)
- [API 参考](docs/reference/api-reference.md)

## 示例

见 [examples/](examples/) 目录

## 贡献

欢迎贡献！见 [CONTRIBUTING.md](CONTRIBUTING.md)
```

## 🎯 重点调整

### 1. 项目定位

**之前**: 工具集合  
**现在**: ClusterReport 集群管理平台

### 2. 文档重点

**之前**: 分散的工具文档  
**现在**: 统一的平台文档 + ClusterReport 为中心

### 3. 使用方式

**之前**:
```bash
# 三个独立工具
nodeprobe
perfsnap  
clusterreport
```

**现在**:
```bash
# 一个核心命令
clusterreport collect    # 包含 nodeprobe + perfsnap 功能
clusterreport analyze
clusterreport generate
clusterreport serve      # 启动 Web 界面

# 兼容命令（调用核心平台）
nodeprobe     # -> clusterreport nodeprobe
perfsnap      # -> clusterreport perfsnap
```

### 4. 开发焦点

**之前**: 维护三个独立工具  
**现在**: 发展 ClusterReport 生态系统

## 📋 文档更新清单

- [ ] README.md（重写）
- [ ] README-zh.md（重写）
- [ ] docs/getting-started/（新建）
- [ ] docs/user-guide/（重组）
- [ ] docs/developer-guide/（新建）
- [ ] docs/reference/（新建）
- [ ] CHANGELOG.md（新建）
- [ ] 所有工具文档（更新为插件文档）

## 🚀 实施时间线

- **Day 1**: 目录结构调整
- **Day 2-3**: 文档重写
- **Day 4-10**: 代码重构
- **Day 11-12**: 构建和部署
- **Day 13-14**: 示例和测试

**总计**: 约 2周

## ✅ 成功标准

1. ClusterReport 成为项目核心，所有文档反映这一点
2. 旧工具作为 ClusterReport 的一部分继续可用
3. 新用户首先了解 ClusterReport
4. 文档完整且易于理解
5. 示例丰富且可运行

## 🎉 预期成果

- 项目定位更清晰
- 文档更专业和完整
- 用户体验更统一
- 生态系统更健康
- 维护成本更低

---

**下一步**: 开始执行 Phase 1 - 目录结构调整
