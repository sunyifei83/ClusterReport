# ClusterReport → 集群交付自动化平台 架构演进方案

## 执行摘要

**日期**: 2025-10-28
**目标**: 从报告工具演进为完整的集群交付自动化平台
**范围**: 全链条工具化 - 从物料检查到交付报告

---

## 一、需求分析：全链条工具化

### 1.1 完整的交付流程

您提出的完整链条：

```
1. 集群节点物料检查
   ├── 硬件配置验证
   ├── 网络连通性检查
   ├── 存储容量确认
   └── 基础软件版本检查

2. 物料基准性能测试
   ├── CPU 性能基准
   ├── 内存带宽测试
   ├── 磁盘 IO 测试
   ├── 网络吞吐量测试
   └── 综合性能评分

3. 生产应用部署
   ├── 环境准备
   ├── 应用安装
   ├── 配置管理
   ├── 服务启动
   └── 健康检查

4. 交付物功能验证
   ├── 功能测试
   ├── 性能验证
   ├── 高可用测试
   ├── 故障恢复测试
   └── 安全检查

5. 交付报告输出
   ├── 综合评估报告
   ├── 测试数据汇总
   ├── 问题清单
   ├── 改进建议
   └── 验收文档
```

### 1.2 当前项目覆盖范围

**已覆盖**（约30%）：
- ✅ 节点物料检查（NodeProbe）
- ✅ 性能数据采集（PerfSnap）
- ✅ 部分性能测试（iotest.sh）
- ✅ 报告生成

**未覆盖**（约70%）：
- ❌ 物料基准性能测试（系统化）
- ❌ 应用部署管理
- ❌ 功能验证框架
- ❌ 完整的Web UI平台
- ❌ 工作流编排
- ❌ 交付管理

---

## 二、架构调整方案

### 2.1 项目重新定位

#### 当前定位
```
ClusterReport = 集群报告生成工具
- 专注于数据收集和报告
- CLI 为主
```

#### 新定位
```
ClusterDelivery = 集群交付自动化平台
或
DevOpsDelivery = DevOps交付平台
或
InfraDelivery = 基础设施交付平台

- 全链条自动化
- Web UI 为主，CLI 为辅
- 涵盖从检查到交付的完整流程
```

### 2.2 推荐项目名称

基于您的需求，推荐以下名称：

| 项目名称 | 优势 | 劣势 | 推荐度 |
|---------|------|------|--------|
| **ClusterDelivery** | 明确表达"交付"概念 | 范围可能过于聚焦集群 | ⭐⭐⭐⭐⭐ |
| **DevOpsHub** | 通用、现代 | 过于宽泛 | ⭐⭐⭐⭐ |
| **InfraFlow** | 强调流程 | 不够直观 | ⭐⭐⭐ |
| **DeliveryPlatform** | 精准 | 过于通用 | ⭐⭐⭐ |
| **AutoDelivery** | 强调自动化 | 略显普通 | ⭐⭐⭐⭐ |

**最终推荐**: `ClusterDelivery` 或 `DevOpsDelivery`

理由：
1. 明确表达"交付"核心概念
2. 与现有名称有延续性
3. 专业且易于理解
4. 适合企业级定位

### 2.3 新架构设计

#### 架构层次

```
┌─────────────────────────────────────────────────────┐
│              Web UI 层（主要界面）                    │
│  - 工作流管理  - 任务监控  - 报告展示  - 系统配置      │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│              API 服务层（核心引擎）                   │
│  REST API + WebSocket + GraphQL                      │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│              工作流引擎层（新增）                     │
│  - 流程定义  - 任务调度  - 状态管理  - 事件通知       │
└─────────────────────────────────────────────────────┘
                        ↓
┌──────────────┬──────────────┬──────────────┬──────────────┐
│  检查模块    │  测试模块    │  部署模块    │  验证模块    │
│ (Checker)   │ (Tester)    │ (Deployer)  │ (Validator) │
└──────────────┴──────────────┴──────────────┴──────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│              执行器层（Agent/Executor）               │
│  - 本地执行  - 远程执行  - 批量执行  - 结果收集       │
└─────────────────────────────────────────────────────┘
                        ↓
┌─────────────────────────────────────────────────────┐
│              数据层                                   │
│  PostgreSQL + Redis + S3/MinIO                       │
└─────────────────────────────────────────────────────┘
```

#### 新增核心模块

```go
// 1. 工作流引擎
pkg/workflow/
├── engine.go           // 工作流引擎
├── pipeline.go         // 流水线定义
├── stage.go           // 阶段定义
├── task.go            // 任务定义
└── executor.go        // 执行器

// 2. 检查模块（扩展）
pkg/checker/
├── hardware.go        // 硬件检查
├── network.go         // 网络检查
├── storage.go         // 存储检查
├── software.go        // 软件检查
└── compliance.go      // 合规性检查

// 3. 测试模块（新增）
pkg/tester/
├── benchmark.go       // 性能基准测试
├── stress.go          // 压力测试
├── io.go              // IO测试
├── network.go         // 网络测试
└── integration.go     // 集成测试

// 4. 部署模块（新增）
pkg/deployer/
├── ansible.go         // Ansible部署
├── kubernetes.go      // K8s部署
├── docker.go          // Docker部署
├── script.go          // 脚本部署
└── service.go         // 服务管理

// 5. 验证模块（新增）
pkg/validator/
├── functional.go      // 功能验证
├── performance.go     // 性能验证
├── ha.go              // 高可用验证
├── security.go        // 安全验证
└── compliance.go      // 合规验证

// 6. 嵌入式服务工具
tools/embedded-service/
├── main.go            // 独立服务入口
├── api/               // API定义
├── handlers/          // 处理器
├── static/            // 静态文件
└── templates/         // 模板
```

---

## 三、项目结构重组

### 3.1 新的目录结构

```
ClusterDelivery/  (或 DevOpsDelivery/)
├── cmd/
│   ├── cli/                    # CLI工具（保留，作为辅助）
│   ├── server/                 # API服务器（主要）
│   ├── agent/                  # Agent程序
│   └── embedded-service/       # 嵌入式服务工具 ✨
│
├── pkg/
│   ├── checker/               # 检查模块 ✨
│   ├── tester/                # 测试模块 ✨
│   ├── deployer/              # 部署模块 ✨
│   ├── validator/             # 验证模块 ✨
│   ├── workflow/              # 工作流引擎 ✨
│   ├── collector/             # 数据收集（已有）
│   ├── analyzer/              # 数据分析（已有）
│   ├── generator/             # 报告生成（已有）
│   ├── storage/               # 存储层 ✨
│   └── notification/          # 通知系统 ✨
│
├── web/
│   ├── frontend/              # React前端 ✨
│   │   ├── src/
│   │   │   ├── pages/
│   │   │   │   ├── Dashboard/      # 仪表盘
│   │   │   │   ├── Workflow/       # 工作流管理
│   │   │   │   ├── Checker/        # 检查管理
│   │   │   │   ├── Tester/         # 测试管理
│   │   │   │   ├── Deployer/       # 部署管理
│   │   │   │   ├── Validator/      # 验证管理
│   │   │   │   └── Reports/        # 报告查看
│   │   │   ├── components/
│   │   │   ├── services/
│   │   │   └── utils/
│   │   └── package.json
│   └── dashboard/             # 旧的静态页面（可删除）
│
├── api/
│   ├── openapi/              # OpenAPI规范
│   ├── proto/                # gRPC定义
│   └── graphql/              # GraphQL schema
│
├── internal/
│   ├── config/               # 配置管理
│   ├── models/               # 数据模型
│   ├── middleware/           # 中间件
│   └── utils/                # 工具函数
│
├── workflows/                # 工作流定义 ✨
│   ├── templates/            # 流程模板
│   │   ├── standard-delivery.yaml
│   │   ├── quick-check.yaml
│   │   └── full-validation.yaml
│   └── custom/               # 自定义流程
│
├── scripts/
│   ├── benchmark/            # 基准测试脚本 ✨
│   ├── validation/           # 验证脚本 ✨
│   └── deployment/           # 部署脚本 ✨
│
├── docs/
│   ├── architecture/         # 架构文档
│   ├── api/                  # API文档
│   ├── workflow/             # 工作流文档 ✨
│   └── user-guide/           # 用户指南
│
├── deployments/
│   ├── docker/
│   ├── kubernetes/
│   └── ansible/
│
└── tests/
    ├── unit/
    ├── integration/
    └── e2e/
```

### 3.2 嵌入式服务工具设计

这是一个关键的新增功能，提供独立的本地Web服务：

```go
// cmd/embedded-service/main.go
package main

import (
    "embed"
    "github.com/gin-gonic/gin"
)

//go:embed static/*
var staticFiles embed.FS

//go:embed templates/*
var templates embed.FS

type EmbeddedService struct {
    router *gin.Engine
    db     *Database
    config *Config
}

func main() {
    svc := NewEmbeddedService()
    
    // 单一可执行文件
    // 内嵌所有资源
    // 零外部依赖
    
    svc.Run(":8080")
}
```

**特点**：
- 📦 单一可执行文件
- 🚀 一键启动
- 💾 内嵌SQLite数据库
- 🎨 完整Web UI
- 🔧 无需配置

---

## 四、工作流引擎设计

### 4.1 流程定义（YAML）

```yaml
# workflows/templates/standard-delivery.yaml
name: 标准集群交付流程
version: 1.0.0
description: 完整的集群交付流程

stages:
  # 阶段1: 物料检查
  - name: material-check
    display: 集群节点物料检查
    tasks:
      - name: hardware-check
        type: checker
        module: hardware
        params:
          cpu_cores_min: 8
          memory_gb_min: 16
          disk_gb_min: 100
        
      - name: network-check
        type: checker
        module: network
        params:
          bandwidth_min: 1000  # Mbps
          latency_max: 5       # ms
      
      - name: software-check
        type: checker
        module: software
        params:
          os_version: "CentOS 7.x"
          kernel_version: ">= 3.10"
  
  # 阶段2: 性能基准测试
  - name: benchmark-test
    display: 物料基准性能测试
    depends_on: [material-check]
    tasks:
      - name: cpu-benchmark
        type: tester
        module: benchmark
        script: "sysbench cpu"
        
      - name: memory-benchmark
        type: tester
        module: benchmark
        script: "sysbench memory"
      
      - name: io-benchmark
        type: tester
        module: io
        script: "fio"
        params:
          mode: "randread"
          size: "10G"
      
      - name: network-benchmark
        type: tester
        module: network
        script: "iperf3"
  
  # 阶段3: 应用部署
  - name: application-deploy
    display: 生产应用部署
    depends_on: [benchmark-test]
    tasks:
      - name: environment-prepare
        type: deployer
        module: ansible
        playbook: "prepare.yml"
      
      - name: app-install
        type: deployer
        module: ansible
        playbook: "install.yml"
      
      - name: config-manage
        type: deployer
        module: ansible
        playbook: "configure.yml"
      
      - name: service-start
        type: deployer
        module: ansible
        playbook: "start.yml"
  
  # 阶段4: 功能验证
  - name: functional-validation
    display: 交付物功能验证
    depends_on: [application-deploy]
    tasks:
      - name: functional-test
        type: validator
        module: functional
        test_suite: "smoke_tests"
      
      - name: performance-validation
        type: validator
        module: performance
        benchmark: "load_test"
      
      - name: ha-test
        type: validator
        module: ha
        scenario: "failover"
      
      - name: security-check
        type: validator
        module: security
        scan: "vulnerability"
  
  # 阶段5: 报告生成
  - name: report-generation
    display: 交付报告输出
    depends_on: [functional-validation]
    tasks:
      - name: collect-results
        type: collector
        sources: [material-check, benchmark-test, application-deploy, functional-validation]
      
      - name: analyze-data
        type: analyzer
        
      - name: generate-report
        type: generator
        formats: [html, pdf, excel]
        template: "delivery-report"

notifications:
  on_success:
    - type: email
      to: ["pm@company.com"]
  on_failure:
    - type: email
      to: ["ops@company.com"]
    - type: webhook
      url: "https://alertmanager.company.com/alert"
```

### 4.2 工作流引擎实现

```go
// pkg/workflow/engine.go
package workflow

type Engine struct {
    pipelines map[string]*Pipeline
    executor  Executor
    storage   Storage
    notifier  Notifier
}

type Pipeline struct {
    ID          string
    Name        string
    Stages      []Stage
    Status      Status
    StartTime   time.Time
    EndTime     time.Time
    Results     map[string]interface{}
}

type Stage struct {
    Name      string
    Tasks     []Task
    DependsOn []string
    Status    Status
}

type Task struct {
    Name    string
    Type    TaskType // checker, tester, deployer, validator
    Module  string
    Params  map[string]interface{}
    Status  Status
    Result  TaskResult
}

func (e *Engine) Execute(pipelineID string) error {
    pipeline := e.pipelines[pipelineID]
    
    for _, stage := range pipeline.Stages {
        // 检查依赖
        if !e.dependenciesMet(stage) {
            return fmt.Errorf("dependencies not met for stage: %s", stage.Name)
        }
        
        // 并发执行任务
        results := make(chan TaskResult, len(stage.Tasks))
        for _, task := range stage.Tasks {
            go e.executeTask(task, results)
        }
        
        // 收集结果
        for i := 0; i < len(stage.Tasks); i++ {
            result := <-results
            if result.Error != nil {
                return e.handleFailure(pipeline, stage, result)
            }
        }
    }
    
    // 生成最终报告
    return e.generateFinalReport(pipeline)
}
```

---

## 五、Web UI 设计

### 5.1 主要页面

#### 1. Dashboard（仪表盘）
```
┌─────────────────────────────────────────────┐
│  🏠 ClusterDelivery Platform                │
├─────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │ 进行中   │  │ 已完成   │  │ 失败     │  │
│  │   15     │  │   142    │  │    3     │  │
│  └──────────┘  └──────────┘  └──────────┘  │
│                                              │
│  📊 最近交付流程                             │
│  ┌─────────────────────────────────────┐    │
│  │ ✅ prod-cluster-01  | 2h ago       │    │
│  │ 🔄 test-cluster-02  | 进行中       │    │
│  │ ❌ dev-cluster-03   | 1天前        │    │
│  └─────────────────────────────────────┘    │
└─────────────────────────────────────────────┘
```

#### 2. 工作流管理
```
┌─────────────────────────────────────────────┐
│  📋 工作流列表                               │
├─────────────────────────────────────────────┤
│  [ + 新建工作流 ]  [ 📥 导入 ]  [ 🔍 搜索 ]  │
│                                              │
│  标准流程:                                   │
│  • 标准集群交付流程                          │
│  • 快速检查流程                              │
│  • 性能压测流程                              │
│                                              │
│  自定义流程:                                 │
│  • 客户A定制流程                             │
│  • 应急验证流程                              │
└─────────────────────────────────────────────┘
```

#### 3. 执行监控
```
┌─────────────────────────────────────────────┐
│  🔄 执行中: prod-cluster-01                  │
├─────────────────────────────────────────────┤
│  进度: ████████░░ 80%                       │
│                                              │
│  ✅ 阶段1: 物料检查      (完成)              │
│  ✅ 阶段2: 性能测试      (完成)              │
│  ✅ 阶段3: 应用部署      (完成)              │
│  🔄 阶段4: 功能验证      (进行中 60%)        │
│  ⏳ 阶段5: 报告生成      (等待中)            │
│                                              │
│  当前任务: 高可用测试                        │
│  [日志] [终止] [重试]                        │
└─────────────────────────────────────────────┘
```

### 5.2 技术栈

```javascript
// web/frontend/package.json
{
  "name": "cluster-delivery-frontend",
  "dependencies": {
    "react": "^18.2.0",
    "react-router-dom": "^6.x",
    "ant-design-pro": "^6.x",      // UI组件
    "@ant-design/pro-components": "^2.x",
    "echarts": "^5.x",              // 图表
    "monaco-editor": "^0.x",        // 代码编辑器
    "react-query": "^3.x",          // 数据获取
    "zustand": "^4.x",              // 状态管理
    "axios": "^1.x"                 // HTTP客户端
  }
}
```

---

## 六、实施路线图

### 6.1 重构阶段（2周）

**Week 1: 项目重组**
- [ ] 重命名项目（ClusterDelivery）
- [ ] 重组目录结构
- [ ] 更新文档和README
- [ ] 迁移现有代码到新结构

**Week 2: 基础架构**
- [ ] 实现工作流引擎框架
- [ ] 创建新模块骨架
- [ ] 设计数据库schema
- [ ] 搭建Web UI骨架

### 6.2 核心功能开发（4-6周）

**Week 3-4: 检查和测试模块**
- [ ] 完善checker模块
- [ ] 实现benchmark测试
- [ ] 集成性能测试工具
- [ ] 开发测试报告

**Week 5-6: 部署和验证模块**
- [ ] 实现deployer模块
- [ ] Ansible集成
- [ ] 实现validator模块
- [ ] 功能测试框架

**Week 7-8: Web UI开发**
- [ ] Dashboard页面
- [ ] 工作流管理界面
- [ ] 实时监控界面
- [ ] 报告展示界面

### 6.3 集成和测试（2周）

**Week 9-10**
- [ ] 端到端流程测试
- [ ] 性能测试
- [ ] 安全测试
- [ ] 文档完善

---

## 七、架构调整关键点

### 7.1 必须调整的部分

1. **项目重命名** ⭐⭐⭐⭐⭐
   ```bash
   ClusterReport → ClusterDelivery
   # 或
   ClusterReport → DevOpsDelivery
   ```

2. **目录重组** ⭐⭐⭐⭐⭐
   ```
   新增:
   - pkg/workflow/
   - pkg/checker/
   - pkg/tester/
   - pkg/deployer/
   - pkg/validator/
   - cmd/embedded-service/
   - web/frontend/
   - workflows/
   ```

3. **数据模型扩展** ⭐⭐⭐⭐⭐
   ```sql
   新增表:
   - workflows
   - pipelines
   - stages
   - tasks
   - executions
   - validations
   - deployments
   ```

4. **Web UI为主** ⭐⭐⭐⭐⭐
   - CLI工具保留但降级为辅助工具
   - Web UI成为主要交互界面

### 7.2 可以保留的部分

1. **现有核心代码** ✅
   ```
   保留:
   - pkg/collector/
   - pkg/analyzer/
   - pkg/generator/
   - cmd/cli/ (降级为工具)
   ```

2. **数据收集逻辑** ✅
   - NodeProbe集成
   - PerfSnap集成
   - 报告生成

3. **配置系统** ✅
   - config.yaml
   - 配置管理逻辑

### 7.3 需要增强的部分

1. **工作流系统** 🆕
   - 流程定义
   - 任务编排
   - 状态管理
   - 事件通知

2. **测试框架** 🆕
   - 性能基准测试
   - 功能验证测试
   - 集成测试

3. **部署能力** 🆕
   - 应用部署
   - 配置管理
   - 服务控制

---

## 八、迁移策略

### 8.1 无缝迁移方案

```
第一步: 创建新分支
git checkout -b feature/platform-evolution

第二步: 重组目录（保留兼容性）
ClusterDelivery/
├── cmd/
│   ├── cli/          # 保留，添加别名
│   └── server/       # 新增主服务
├── legacy/           # 移动旧代码
│   └── clusterreport/
└── pkg/
    ├── core/         # 新核心功能
    └── compat/       # 兼容层

第三步: 渐进式迁移
- 保持CLI命令兼容
- 逐步添加新功能
- 文档同步更新

第四步: 发布策略
- v1.0: 现有功能（ClusterReport）
- v2.0: 平台化版本（ClusterDelivery）
```

### 8.2 向后兼容性

```go
// cmd/cli/main.go
func main() {
    // 检测是否为旧命令
    if isLegacyCommand() {
        fmt.Println("注意: 您正在使用旧的ClusterReport命令")
        fmt.Println("建议升级到新的ClusterDelivery平台")
        fmt.Println("运行 'clusterdelivery migrate' 了解更多")
    }
    
    // 执行命令（兼容旧逻辑）
    execute()
}
```

---

## 九、技术栈完整清单

### 后端技术栈

```go
核心框架:
- Go 1.19+
- Gin Web Framework
- GORM (数据库ORM)
- Viper (配置管理)

工作流:
- Temporal (可选，专业工作流引擎)
- 或自研轻量级引擎

数据库:
- PostgreSQL (主数据库)
- Redis (缓存+队列)
- SQLite (嵌入式服务)

消息队列:
- RabbitMQ 或 NATS

监控:
- Prometheus
- Grafana
```

### 前端技术栈

```javascript
核心框架:
- React 18
- TypeScript
- Ant Design Pro

状态管理:
- Zustand / Redux Toolkit

数据获取:
- React Query
- Axios

可视化:
- ECharts
- D3.js (高级图表)

实时通信:
- WebSocket
- Socket.io

代码编辑:
- Monaco Editor (YAML编辑)
```

---

## 十、关键决策建议

### 10.1 立即决策

1. **项目重命名** 🔴
   ```
   推荐: ClusterDelivery
   备选: DevOpsDelivery, InfraDelivery
   
   行动: 
   - 更新仓库名
   - 更新所有文档
   - 通知相关人员
   ```

2. **架构转型确认** 🔴
   ```
   从: 报告工具
   到: 交付平台
   
   影响:
   - 项目定位变化
   - 团队规模需求
   - 开发周期延长
   ```

### 10.2 技术选型

1. **工作流引擎** 🟠
   ```
   选项A: 自研轻量级引擎
   - 优点: 灵活、简单、可控
   - 缺点: 需要开发时间
   - 推荐度: ⭐⭐⭐⭐⭐ (适合MVP)
   
   选项B: Temporal
   - 优点: 成熟、功能强大
   - 缺点: 复杂、依赖多
   - 推荐度: ⭐⭐⭐ (适合长期)
   ```

2. **嵌入式服务** 🟠
   ```
   实现方式:
   - embed.FS (Go 1.16+)
   - 单一可执行文件
   - 内嵌SQLite
   - 内嵌前端资源
   
   优势:
   - 零配置启动
   - 易于分发
   - 适合POC/演示
   ```

### 10.3 实施建议

**建议顺序**:

1. **Week 1-2**: 项目重组和重命名
   - 影响: 高
   - 风险: 低
   - 优先级: 🔴 P0

2. **Week 3-4**: 完成v0.5.0剩余功能
   - 图表、PDF、持久化
   - 影响: 中
   - 风险: 低
   - 优先级: 🔴 P0

3. **Week 5-8**: 开发核心平台功能
   - 工作流引擎
   - 测试/部署/验证模块
   - 影响: 高
   - 风险: 中
   - 优先级: 🟠 P1

4. **Week 9-12**: Web UI开发
   - 前端界面
   - API集成
   - 影响: 高
   - 风险: 中
   - 优先级: 🟠 P1

---

## 十一、总结

### 核心变化

| 维度 | 当前 | 目标 |
|------|------|------|
| **定位** | 报告工具 | 交付平台 |
| **范围** | 数据收集+报告 | 全链条自动化 |
| **界面** | CLI为主 | Web UI为主 |
| **流程** | 线性执行 | 工作流编排 |
| **功能** | 30%覆盖 | 100%覆盖 |

### 建议的项目名称

🏆 **最终推荐**: `ClusterDelivery`

**理由**:
1. ✅ 明确表达"交付"核心价值
2. ✅ 与现有名称有延续性
3. ✅ 专业且易于理解
4. ✅ 适合企业级定位
5. ✅ 域名可用性好

### 关键里程碑

```
v1.0.0-rc1 ✅ - ClusterReport基础功能
v1.5.0     🔄 - 完善报告功能 + 项目重组
v2.0.0     🎯 - ClusterDelivery平台发布
v2.5.0     🚀 - 完整交付流程自动化
v3.0.0     🌟 - 企业级特性 + 多云支持
```

---

**文档版本**: v1.0
**创建日期**: 2025-10-28

