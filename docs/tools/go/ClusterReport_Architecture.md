# ClusterReport 架构设计文档

## 1. 需求分析与整理

### 1.1 功能需求维度

#### 数据采集层
- **工具集成需求**
  - NodeProbe: 硬件配置信息
  - PerfSnap: 性能监控数据
  - iotest.sh: IO性能测试
  - 自定义脚本: 特定业务指标
  - 第三方监控: Prometheus、Zabbix等

#### 数据处理层
- **数据标准化**: 统一不同工具的输出格式
- **数据存储**: 历史数据保存和查询
- **数据分析**: 趋势分析、异常检测
- **数据聚合**: 多节点数据汇总

#### 报告生成层
- **模板引擎**: 灵活的报告模板系统
- **图表生成**: 数据可视化
- **格式输出**: HTML、PDF、Excel、Markdown
- **定制化**: 企业品牌、自定义章节

#### 管理控制层
- **任务调度**: 定时任务、触发器
- **权限管理**: 用户认证、角色权限
- **配置管理**: 集群配置、报告模板
- **通知告警**: 邮件、webhook通知

### 1.2 非功能需求

- **性能**: 支持1000+节点数据处理
- **可扩展**: 插件化架构，易于添加新功能
- **可靠性**: 失败重试、断点续传
- **安全性**: 数据加密、访问控制
- **易用性**: 简单部署、最小依赖

## 2. 架构选型分析

### 2.1 架构模式对比

| 架构模式 | 优点 | 缺点 | 适用场景 |
|---------|------|------|---------|
| **单体应用** | 简单部署、低延迟、易调试 | 扩展性差、技术栈单一 | 小规模集群(<100节点) |
| **C/S架构** | 集中管理、易维护、安全性好 | 需要服务端部署、网络依赖 | 中等规模(100-500节点) |
| **微服务** | 高扩展性、技术栈灵活、容错性好 | 复杂度高、运维成本大 | 大规模(>500节点) |
| **Serverless** | 无需管理服务器、按需付费 | 冷启动、厂商锁定 | 不定期报告生成 |
| **混合架构** | 结合多种优点、渐进式演进 | 架构复杂、需要经验 | 逐步发展的系统 |

### 2.2 推荐架构方案

基于需求分析，推荐采用**渐进式混合架构**：

```
第一阶段：CLI工具（单体）
第二阶段：C/S架构（API服务）
第三阶段：微服务架构（按需）
```

## 3. 模块化设计方案

### 3.1 核心模块架构

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
    ├── ansible/          # ansible配置
    ├── docker/          # Docker配置
    ├── k8s/             # Kubernetes配置
    └── helm/            # Helm Charts
```

### 3.2 模块化接口设计

#### 3.2.1 采集器接口
```go
// Collector 数据采集器接口
type Collector interface {
    // 名称标识
    Name() string
    
    // 采集数据
    Collect(ctx context.Context, node Node) (*Data, error)
    
    // 验证配置
    Validate(config Config) error
    
    // 支持的数据类型
    SupportedTypes() []DataType
}

// 插件注册
type CollectorRegistry struct {
    collectors map[string]Collector
}

func (r *CollectorRegistry) Register(c Collector) {
    r.collectors[c.Name()] = c
}
```

#### 3.2.2 分析器接口
```go
// Analyzer 数据分析器接口
type Analyzer interface {
    // 分析数据
    Analyze(ctx context.Context, data []Data) (*Analysis, error)
    
    // 支持的分析类型
    Type() AnalysisType
    
    // 配置选项
    Options() map[string]interface{}
}

// 分析链
type AnalysisChain struct {
    analyzers []Analyzer
}

func (c *AnalysisChain) Process(data []Data) (*Report, error) {
    // 链式处理
}
```

#### 3.2.3 生成器接口
```go
// Generator 报告生成器接口
type Generator interface {
    // 生成报告
    Generate(ctx context.Context, analysis *Analysis) ([]byte, error)
    
    // 输出格式
    Format() OutputFormat
    
    // 模板支持
    SetTemplate(tmpl Template) error
}

// 生成器工厂
type GeneratorFactory struct {
    generators map[OutputFormat]Generator
}
```

### 3.3 插件系统设计

#### 3.3.1 插件加载机制
```go
// Plugin 插件接口
type Plugin interface {
    // 插件元信息
    Info() PluginInfo
    
    // 初始化
    Init(config Config) error
    
    // 启动
    Start(ctx context.Context) error
    
    // 停止
    Stop(ctx context.Context) error
}

// PluginManager 插件管理器
type PluginManager struct {
    plugins   map[string]Plugin
    loader    PluginLoader
    registry  PluginRegistry
}

// 动态加载插件
func (m *PluginManager) LoadPlugin(path string) error {
    // 使用 plugin.Open() 加载 .so 文件
    // 或使用 go-plugin 库
}
```

#### 3.3.2 插件通信
```go
// 使用 HashiCorp go-plugin
type CollectorPlugin struct {
    plugin.NetRPCUnsupportedPlugin
}

func (p *CollectorPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
    // 注册 gRPC 服务
}

func (p *CollectorPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
    // 返回 gRPC 客户端
}
```

## 4. 实现方案详细设计

### 4.1 第一阶段：CLI工具实现

#### 4.1.1 架构设计
```yaml
架构模式: 单体CLI应用
部署方式: 单一二进制文件
数据存储: 本地文件系统
配置方式: YAML配置文件
```

#### 4.1.2 核心功能
```go
// 主要命令
clusterreport collect    # 数据采集
clusterreport analyze    # 数据分析
clusterreport generate   # 报告生成
clusterreport schedule   # 任务调度

// 工作流程
1. 读取配置文件
2. SSH连接到各节点
3. 执行采集命令
4. 本地分析数据
5. 生成报告文件
```

#### 4.1.3 实现示例
```go
package main

import (
    "github.com/spf13/cobra"
    "github.com/devops-toolkit/clusterreport/pkg/collector"
    "github.com/devops-toolkit/clusterreport/pkg/analyzer"
    "github.com/devops-toolkit/clusterreport/pkg/generator"
)

func main() {
    app := &Application{
        Collector: collector.NewMultiCollector(
            collector.NewNodeProbeCollector(),
            collector.NewPerfSnapCollector(),
            collector.NewCustomCollector(),
        ),
        Analyzer: analyzer.NewChain(
            analyzer.NewConfigAnalyzer(),
            analyzer.NewPerfAnalyzer(),
            analyzer.NewAnomalyDetector(),
        ),
        Generator: generator.NewMultiFormat(
            generator.NewHTMLGenerator(),
            generator.NewPDFGenerator(),
            generator.NewExcelGenerator(),
        ),
    }
    
    rootCmd := &cobra.Command{
        Use:   "clusterreport",
        Short: "Cluster comprehensive report generator",
    }
    
    rootCmd.AddCommand(
        app.CollectCommand(),
        app.AnalyzeCommand(),
        app.GenerateCommand(),
    )
    
    rootCmd.Execute()
}
```

### 4.2 第二阶段：C/S架构实现

#### 4.2.1 架构设计
```yaml
架构模式: Client-Server
通信协议: REST API + WebSocket
数据存储: PostgreSQL/MongoDB
缓存层: Redis
消息队列: RabbitMQ/Kafka
```

#### 4.2.2 Server端设计
```go
// API Server
type Server struct {
    // HTTP服务
    httpServer *http.Server
    
    // WebSocket管理
    wsHub *websocket.Hub
    
    // 数据库连接
    db *gorm.DB
    
    // 缓存
    cache *redis.Client
    
    // 任务队列
    queue *amqp.Channel
    
    // 服务注册
    services map[string]Service
}

// REST API路由
func (s *Server) setupRoutes() {
    r := gin.New()
    
    // API版本控制
    v1 := r.Group("/api/v1")
    {
        // 集群管理
        v1.GET("/clusters", s.ListClusters)
        v1.POST("/clusters", s.CreateCluster)
        v1.GET("/clusters/:id", s.GetCluster)
        
        // 数据采集
        v1.POST("/collect", s.TriggerCollection)
        v1.GET("/collect/:id/status", s.GetCollectionStatus)
        
        // 报告管理
        v1.POST("/reports", s.GenerateReport)
        v1.GET("/reports", s.ListReports)
        v1.GET("/reports/:id", s.GetReport)
        
        // WebSocket
        v1.GET("/ws", s.WebSocketHandler)
    }
}
```

#### 4.2.3 Client端设计
```go
// CLI Client
type Client struct {
    baseURL string
    token   string
    client  *http.Client
}

// SDK封装
func (c *Client) CollectData(nodes []string) (*CollectionJob, error) {
    // 调用服务端API
}

func (c *Client) GenerateReport(params ReportParams) (*Report, error) {
    // 调用服务端API
}

// Web Dashboard
// 使用 React/Vue 开发管理界面
```

### 4.3 第三阶段：微服务架构（可选）

#### 4.3.1 服务拆分
```yaml
服务列表:
  - gateway-service      # API网关
  - auth-service        # 认证服务
  - collector-service   # 采集服务
  - analyzer-service    # 分析服务
  - generator-service   # 生成服务
  - scheduler-service   # 调度服务
  - storage-service     # 存储服务
  - notification-service # 通知服务
```

#### 4.3.2 服务通信
```go
// 使用 gRPC 进行服务间通信
type CollectorService interface {
    Collect(ctx context.Context, req *CollectRequest) (*CollectResponse, error)
}

// 服务发现
// 使用 Consul/Etcd 进行服务注册与发现
```

## 5. 技术栈选择

### 5.1 核心技术栈

| 组件 | 推荐技术 | 备选技术 | 选择理由 |
|------|---------|---------|---------|
| **编程语言** | Go | Rust, Python | 性能好、部署简单、并发支持 |
| **Web框架** | Gin | Echo, Fiber | 成熟稳定、社区活跃 |
| **数据库** | PostgreSQL | MySQL, MongoDB | 功能全面、JSON支持 |
| **缓存** | Redis | Memcached | 功能丰富、持久化支持 |
| **消息队列** | RabbitMQ | Kafka, NATS | 易用性、可靠性 |
| **配置管理** | Viper | - | Go标准配置库 |
| **日志** | Zap | Logrus | 高性能、结构化日志 |
| **监控** | Prometheus | - | 云原生标准 |
| **前端** | React | Vue, Angular | 生态完善、组件丰富 |

### 5.2 开发工具链

```yaml
开发环境:
  IDE: VSCode/GoLand
  调试: Delve
  测试: Go Test + Testify
  性能分析: pprof
  代码质量: golangci-lint

CI/CD:
  版本控制: Git
  CI: GitHub Actions/GitLab CI
  构建: Makefile/Task
  容器化: Docker/Buildah
  包管理: Goreleaser

部署工具:
  容器编排: Docker Compose/Kubernetes
  配置管理: Ansible/Terraform
  服务网格: Istio (可选)
```

## 6. 扩展性设计

### 6.1 插件开发指南

#### 6.1.1 创建自定义采集器
```go
// custom_collector.go
package collectors

import (
    "context"
    "github.com/devops-toolkit/clusterreport/pkg/collector"
)

type CustomCollector struct {
    name   string
    config map[string]interface{}
}

func NewCustomCollector(name string) *CustomCollector {
    return &CustomCollector{
        name: name,
    }
}

func (c *CustomCollector) Name() string {
    return c.name
}

func (c *CustomCollector) Collect(ctx context.Context, node collector.Node) (*collector.Data, error) {
    // 实现数据采集逻辑
    // 1. SSH连接到节点
    // 2. 执行命令或脚本
    // 3. 解析输出
    // 4. 返回标准化数据
    
    data := &collector.Data{
        Node:      node.Name,
        Type:      "custom",
        Timestamp: time.Now(),
        Metrics:   make(map[string]interface{}),
    }
    
    // 采集逻辑
    cmd := "your-custom-command"
    output, err := node.Execute(ctx, cmd)
    if err != nil {
        return nil, err
    }
    
    // 解析输出
    data.Metrics = parseOutput(output)
    
    return data, nil
}

// 注册插件
func init() {
    collector.Register("custom", NewCustomCollector("custom"))
}
```

#### 6.1.2 创建自定义分析器
```go
// anomaly_analyzer.go
package analyzers

import (
    "github.com/devops-toolkit/clusterreport/pkg/analyzer"
)

type AnomalyAnalyzer struct {
    threshold float64
}

func (a *AnomalyAnalyzer) Analyze(ctx context.Context, data []analyzer.Data) (*analyzer.Analysis, error) {
    analysis := &analyzer.Analysis{
        Type:      "anomaly",
        Timestamp: time.Now(),
        Results:   make(map[string]interface{}),
    }
    
    // 异常检测逻辑
    anomalies := []string{}
    for _, d := range data {
        if isAnomaly(d, a.threshold) {
            anomalies = append(anomalies, d.Node)
        }
    }
    
    analysis.Results["anomalies"] = anomalies
    analysis.Results["count"] = len(anomalies)
    
    return analysis, nil
}
```

### 6.2 集成外部工具

#### 6.2.1 集成Prometheus
```go
// prometheus_collector.go
type PrometheusCollector struct {
    client api.Client
}

func (p *PrometheusCollector) Collect(ctx context.Context, query string) (*Data, error) {
    v1api := v1.NewAPI(p.client)
    result, _, err := v1api.Query(ctx, query, time.Now())
    if err != nil {
        return nil, err
    }
    
    // 转换Prometheus数据到内部格式
    return convertPrometheusData(result), nil
}
```

#### 6.2.2 集成Ansible
```go
// ansible_executor.go
type AnsibleExecutor struct {
    inventory string
    playbook  string
}

func (a *AnsibleExecutor) Execute(ctx context.Context) error {
    cmd := exec.CommandContext(ctx, "ansible-playbook",
        "-i", a.inventory,
        a.playbook,
        "--extra-vars", "@vars.json",
    )
    
    return cmd.Run()
}
```

## 7. 实施路线图

### 7.1 MVP版本 (v0.1.0) - 4周

**目标**: 基础CLI工具，支持本地报告生成

**功能清单**:
- ✅ 基础命令行框架
- ✅ NodeProbe集成
- ✅ PerfSnap集成
- ✅ Markdown报告生成
- ✅ YAML配置文件

**技术实现**:
```bash
# 项目结构
clusterreport/
├── cmd/
│   └── clusterreport/
│       └── main.go
├── pkg/
│   ├── collector/
│   │   ├── nodeprobe.go
│   │   └── perfsnap.go
│   └── generator/
│       └── markdown.go
├── config.yaml
└── go.mod
```

### 7.2 基础版本 (v0.5.0) - 8周

**目标**: 完整的单机版工具

**新增功能**:
- 📊 HTML报告生成
- 📈 基础图表支持
- 🔄 批量节点处理
- 📝 报告模板系统
- 💾 本地数据存储

**关键代码**:
```go
// 报告模板系统
type Template struct {
    Name     string
    Sections []Section
    Assets   []Asset
}

type Section struct {
    Title    string
    Type     string // text, table, chart
    DataPath string
    Template string
}

// HTML生成器
type HTMLGenerator struct {
    template *Template
    charts   ChartRenderer
}

func (h *HTMLGenerator) Generate(data *ReportData) ([]byte, error) {
    // 渲染HTML模板
    tmpl := template.Must(template.ParseFiles(h.template.Path))
    
    var buf bytes.Buffer
    err := tmpl.Execute(&buf, data)
    
    return buf.Bytes(), err
}
```

### 7.3 专业版本 (v1.0.0) - 12周

**目标**: 生产可用的C/S架构

**新增功能**:
- 🌐 REST API服务
- 🖥️ Web管理界面
- 📊 高级数据分析
- 🔐 用户认证授权
- 📧 邮件通知
- 🗄️ 数据库存储

**架构升级**:
```yaml
服务组件:
  API Server:
    - Gin Web框架
    - JWT认证
    - Swagger文档
    
  Web Dashboard:
    - React前端
    - Ant Design组件
    - ECharts图表
    
  数据存储:
    - PostgreSQL主库
    - Redis缓存
    - MinIO对象存储
```

### 7.4 企业版本 (v2.0.0) - 24周

**目标**: 大规模集群支持

**新增功能**:
- 🎯 微服务架构
- 🔄 自动故障恢复
- 📊 机器学习分析
- 🌍 多集群管理
- 📱 移动端支持
- 🔌 第三方集成

## 8. 开发最佳实践

### 8.1 代码组织原则

```go
// 领域驱动设计(DDD)
project/
├── domain/           # 领域模型
│   ├── cluster/     # 集群领域
│   ├── report/      # 报告领域
│   └── user/        # 用户领域
├── application/     # 应用服务
│   ├── command/     # 命令处理
│   └── query/       # 查询处理
├── infrastructure/  # 基础设施
│   ├── persistence/ # 数据持久化
│   └── external/    # 外部服务
└── interfaces/      # 接口层
    ├── api/        # API接口
    └── cli/        # CLI接口
```

### 8.2 测试策略

```go
// 单元测试示例
func TestNodeProbeCollector_Collect(t *testing.T) {
    // 准备测试数据
    mockNode := &MockNode{
        Name: "test-node",
        Host: "127.0.0.1",
    }
    
    // 创建采集器
    collector := NewNodeProbeCollector()
    
    // 执行测试
    data, err := collector.Collect(context.Background(), mockNode)
    
    // 断言结果
    assert.NoError(t, err)
    assert.NotNil(t, data)
    assert.Equal(t, "test-node", data.Node)
}

// 集成测试
func TestReportGeneration(t *testing.T) {
    // 端到端测试
    app := NewApplication()
    
    // 收集数据
    data := app.CollectData([]string{"node1", "node2"})
    
    // 分析数据
    analysis := app.AnalyzeData(data)
    
    // 生成报告
    report := app.GenerateReport(analysis)
    
    // 验证报告
    assert.Contains(t, report, "Executive Summary")
}
```

### 8.3 性能优化

```go
// 并发处理
func (c *Collector) CollectFromNodes(nodes []Node) []Data {
    var wg sync.WaitGroup
    dataChan := make(chan *Data, len(nodes))
    
    // 限制并发数
    semaphore := make(chan struct{}, 10)
    
    for _, node := range nodes {
        wg.Add(1)
        go func(n Node) {
            defer wg.Done()
            
            semaphore <- struct{}{}
            defer func() { <-semaphore }()
            
            data, err := c.Collect(context.Background(), n)
            if err == nil {
                dataChan <- data
            }
        }(node)
    }
    
    wg.Wait()
    close(dataChan)
    
    // 收集结果
    results := []Data{}
    for data := range dataChan {
        results = append(results, *data)
    }
    
    return results
}

// 缓存优化
type CachedCollector struct {
    collector Collector
    cache     *cache.Cache
    ttl       time.Duration
}

func (c *CachedCollector) Collect(ctx context.Context, node Node) (*Data, error) {
    key := fmt.Sprintf("%s-%s", c.collector.Name(), node.Name)
    
    // 尝试从缓存获取
    if cached, found := c.cache.Get(key); found {
        return cached.(*Data), nil
    }
    
    // 执行实际采集
    data, err := c.collector.Collect(ctx, node)
    if err != nil {
        return nil, err
    }
    
    // 保存到缓存
    c.cache.Set(key, data, c.ttl)
    
    return data, nil
}
```

## 9. 部署方案

### 9.1 Docker部署

```dockerfile
# Dockerfile
FROM golang:1.19-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o clusterreport cmd/clusterreport/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/clusterreport .
COPY --from=builder /app/config.yaml .

EXPOSE 8080
CMD ["./clusterreport", "server"]
```

### 9.2 Kubernetes部署

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: clusterreport
spec:
  replicas: 3
  selector:
    matchLabels:
      app: clusterreport
  template:
    metadata:
      labels:
        app: clusterreport
    spec:
      containers:
      - name: clusterreport
        image: clusterreport:latest
        ports:
        - containerPort: 8080
        env:
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: db-secret
              key: host
        volumeMounts:
        - name: config
          mountPath: /etc/clusterreport
      volumes:
      - name: config
        configMap:
          name: clusterreport-config
```

### 9.3 Helm Chart

```yaml
# values.yaml
replicaCount: 3

image:
  repository: clusterreport
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 8080

ingress:
  enabled: true
  host: clusterreport.example.com

postgresql:
  enabled: true
  auth:
    database: clusterreport

redis:
  enabled: true
  auth:
    enabled: false
```

## 10. 总结与建议

### 10.1 架构选择建议

根据不同规模和需求，建议采用以下架构：

| 场景 | 推荐架构 | 理由 |
|------|---------|------|
| **POC/试用** | CLI单体 | 快速验证、零部署成本 |
| **小型团队** | CLI + 本地存储 | 简单易用、维护成本低 |
| **中型企业** | C/S架构 | 集中管理、权限控制 |
| **大型企业** | 微服务架构 | 高可用、可扩展 |
| **云原生** | Serverless | 弹性伸缩、按需付费 |

### 10.2 技术债务管理

1. **定期重构**: 每个版本预留20%时间用于重构
2. **代码审查**: 所有PR必须经过审查
3. **技术文档**: 保持文档与代码同步
4. **依赖更新**: 定期更新第三方依赖
5. **性能监控**: 持续监控系统性能

### 10.3 下一步行动

1. **立即开始**: 从MVP版本开始，快速迭代
2. **用户反馈**: 尽早收集用户反馈，调整方向
3. **社区建设**: 建立开源社区，吸引贡献者
4. **商业模式**: 考虑开源核心+商业增值服务
5. **生态构建**: 与其他工具集成，形成生态
