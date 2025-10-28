# ClusterReport 开发路线图

**项目状态**: 🚧 开发中  
**当前版本**: v0.7.0 (70% 完成)  
**目标版本**: v1.0.0  
**预计发布**: 6周后

---

## 📊 当前进展

### ✅ 已完成的功能（v0.7.0）

#### 核心模块
- ✅ **采集器框架** (`pkg/collector/`)
  - 完整的指标数据结构
  - 系统指标采集（CPU、内存、磁盘、网络）
  - 单元测试框架
  - 可扩展的采集器接口

- ✅ **智能分析器** (`pkg/analyzer/`)
  - 多维度指标分析
  - 智能健康评分算法（0-100分）
  - 自动问题检测
  - 智能优化建议生成
  - 可配置阈值

- ✅ **多格式报告生成器** (`pkg/generator/`)
  - JSON 格式支持
  - HTML 精美报告（含完整CSS样式）
  - Markdown 文档格式
  - 模板化报告生成

#### 插件系统
- ✅ **采集器插件**
  - 自定义采集器示例
  - MySQL 数据库采集器
  - Redis 缓存采集器
  
- ✅ **分析器插件**
  - 异常检测分析器
  - 可扩展分析器框架

#### 前端界面
- ✅ **Web Dashboard** (`web/dashboard/`)
  - 现代化响应式设计
  - 实时监控视图
  - 交互式数据展示
  - 自动刷新功能
  - 状态可视化

#### 基础设施
- ✅ 项目结构设计
- ✅ 配置文件框架
- ✅ 基础 CLI 框架
- ✅ 完整的文档

**当前完成度: 70%**

---

## 🎯 开发路线图

### 🔥 阶段 1: CLI 模式增强（v0.8.0）
**优先级**: ⭐⭐⭐⭐⭐ 最高  
**预计时间**: 3-5 天  
**完成后进度**: 85%

#### 目标
完善命令行工具，使其成为可独立使用的强大工具

#### 任务列表

**1.1 命令实现**
- [ ] `clusterreport collect` - 数据采集
  ```bash
  # 本地采集
  clusterreport collect --local
  
  # 远程采集（SSH）
  clusterreport collect --host 192.168.1.100 --user root
  
  # 批量采集
  clusterreport collect --hosts-file hosts.txt
  
  # 自定义采集项
  clusterreport collect --metrics cpu,memory,disk
  ```

- [ ] `clusterreport analyze` - 数据分析
  ```bash
  # 分析采集的数据
  clusterreport analyze --input metrics.json
  
  # 自定义阈值
  clusterreport analyze --config custom-thresholds.yaml
  
  # 指定分析器
  clusterreport analyze --analyzers performance,security
  ```

- [ ] `clusterreport generate` - 报告生成
  ```bash
  # 生成 HTML 报告
  clusterreport generate --format html --output report.html
  
  # 生成 JSON 报告
  clusterreport generate --format json --output report.json
  
  # 生成 Markdown 报告
  clusterreport generate --format markdown --output report.md
  
  # 使用自定义模板
  clusterreport generate --template custom.tmpl
  ```

- [ ] `clusterreport report` - 一键生成完整报告
  ```bash
  # 一键生成报告（采集 + 分析 + 生成）
  clusterreport report --output report.html
  
  # 远程主机报告
  clusterreport report --host 192.168.1.100 --output report.html
  
  # 批量主机报告
  clusterreport report --hosts-file hosts.txt --output-dir reports/
  
  # 发送邮件
  clusterreport report --email admin@example.com
  ```

**1.2 配置管理**
- [ ] 完善 `config.yaml` 配置项
  ```yaml
  # 采集配置
  collector:
    interval: 60s
    timeout: 30s
    metrics:
      - cpu
      - memory
      - disk
      - network
  
  # 分析配置
  analyzer:
    thresholds:
      cpu_warning: 70
      cpu_critical: 90
      memory_warning: 80
      memory_critical: 95
  
  # 报告配置
  generator:
    format: html
    template: default
    output: ./reports
  
  # SSH 配置
  ssh:
    port: 22
    timeout: 10s
    key_file: ~/.ssh/id_rsa
  ```

- [ ] 支持多环境配置
  ```bash
  clusterreport --config dev.yaml
  clusterreport --config prod.yaml
  ```

- [ ] 配置验证功能
  ```bash
  clusterreport config validate
  clusterreport config show
  ```

**1.3 用户体验优化**
- [ ] 彩色终端输出
  - 成功信息：绿色
  - 警告信息：黄色
  - 错误信息：红色
  - 信息提示：蓝色

- [ ] 进度条显示
  ```
  Collecting metrics... [████████████████████] 100% (5/5 hosts)
  Analyzing data...     [████████████░░░░░░░░]  60%
  ```

- [ ] 详细日志选项
  ```bash
  clusterreport --log-level debug
  clusterreport --verbose
  ```

- [ ] 静默模式
  ```bash
  clusterreport --silent
  clusterreport --quiet
  ```

**交付成果**
- ✅ 功能完整的 CLI 工具
- ✅ 完善的配置系统
- ✅ 友好的用户界面
- ✅ 详细的命令文档

---

### 🚀 阶段 2: Server/Agent 架构（v0.9.0）
**优先级**: ⭐⭐⭐⭐ 高  
**预计时间**: 7-10 天  
**完成后进度**: 95%

#### 目标
实现分布式采集架构，支持大规模集群监控

#### 系统架构

```
┌─────────────────────────────────────┐
│         Web Dashboard               │
│      (Browser / Mobile App)         │
└─────────────┬───────────────────────┘
              │ HTTP/WebSocket
┌─────────────▼───────────────────────┐
│         ClusterReport Server         │
│  ┌──────────────────────────────┐   │
│  │   HTTP/WebSocket Server      │   │
│  │   - Web UI Serving           │   │
│  │   - REST API                 │   │
│  │   - Real-time Data Stream    │   │
│  └──────────┬───────────────────┘   │
│             │                        │
│  ┌──────────▼───────────────────┐   │
│  │      gRPC Server             │   │
│  │   - Agent Management         │   │
│  │   - Metrics Collection       │   │
│  │   - Heartbeat Monitoring     │   │
│  └──────────┬───────────────────┘   │
│             │                        │
│  ┌──────────▼───────────────────┐   │
│  │    Storage Layer             │   │
│  │   - Time Series DB           │   │
│  │   - Report Archive           │   │
│  │   - Configuration Storage    │   │
│  └──────────────────────────────┘   │
└─────────────┬───────────────────────┘
              │ gRPC (TLS)
      ┌───────┴────────┬──────────┐
┌─────▼─────┐  ┌───────▼─────┐  ┌▼──────────┐
│  Agent 1  │  │   Agent 2   │  │  Agent N  │
│ (Node 1)  │  │  (Node 2)   │  │ (Node N)  │
│           │  │             │  │           │
│ Collector │  │  Collector  │  │ Collector │
│ Analyzer  │  │  Analyzer   │  │ Analyzer  │
└───────────┘  └─────────────┘  └───────────┘
```

#### 任务列表

**2.1 Server 端实现**

- [ ] **HTTP/WebSocket Server**
  ```go
  // api/rest/server.go
  - GET  /api/v1/agents              // 获取所有 Agent 列表
  - GET  /api/v1/agents/:id          // 获取单个 Agent 详情
  - GET  /api/v1/agents/:id/metrics  // 获取 Agent 指标
  - POST /api/v1/reports             // 生成报告
  - GET  /api/v1/reports             // 获取报告列表
  - GET  /api/v1/reports/:id         // 获取单个报告
  - GET  /api/v1/health              // 健康检查
  - WS   /api/v1/ws                  // WebSocket 实时数据
  ```

- [ ] **gRPC Server**
  ```protobuf
  // api/grpc/proto/clusterreport.proto
  
  service AgentService {
    rpc Register(RegisterRequest) returns (RegisterResponse);
    rpc Heartbeat(HeartbeatRequest) returns (HeartbeatResponse);
    rpc Unregister(UnregisterRequest) returns (UnregisterResponse);
  }
  
  service MetricsService {
    rpc ReportMetrics(MetricsRequest) returns (MetricsResponse);
    rpc StreamMetrics(stream MetricsRequest) returns (stream MetricsResponse);
  }
  
  service ReportService {
    rpc GenerateReport(ReportRequest) returns (ReportResponse);
    rpc GetReport(GetReportRequest) returns (GetReportResponse);
    rpc ListReports(ListReportsRequest) returns (ListReportsResponse);
  }
  ```

- [ ] **Agent 管理**
  ```go
  // pkg/server/agent_manager.go
  - Agent 注册/注销
  - 心跳检测（30秒超时）
  - Agent 状态跟踪
  - 离线检测和告警
  ```

- [ ] **数据接收和存储**
  ```go
  - 接收 Agent 上报的指标
  - 数据验证和清洗
  - 持久化存储
  - 数据聚合
  ```

**2.2 Agent 端实现**

- [ ] **核心功能**
  ```go
  // cmd/agent/main.go
  - 连接到 Server（gRPC + TLS）
  - 定期数据采集（可配置间隔）
  - 数据上报（批量/流式）
  - 心跳保持（每30秒）
  - 断线重连（指数退避）
  - 本地缓存（网络中断时）
  ```

- [ ] **配置管理**
  ```yaml
  # agent-config.yaml
  server:
    address: server.example.com:9090
    tls:
      enabled: true
      cert_file: /etc/clusterreport/cert.pem
      key_file: /etc/clusterreport/key.pem
      ca_file: /etc/clusterreport/ca.pem
  
  agent:
    id: node-01  # 唯一标识
    name: Web Server 01
    labels:
      env: production
      region: us-west-1
      role: web
  
  collector:
    interval: 60s
    timeout: 30s
    retry: 3
  
  cache:
    enabled: true
    size: 1000
    ttl: 1h
  ```

**2.3 安全性**
- [ ] TLS/SSL 加密通信
- [ ] Agent 认证（Token/Certificate）
- [ ] 数据加密存储
- [ ] 访问控制（RBAC）

**交付成果**
- ✅ 可扩展的 Server/Agent 架构
- ✅ gRPC 通信协议
- ✅ 安全的数据传输
- ✅ 分布式采集能力

---

### 💾 阶段 3: 存储层实现（v0.95.0）
**优先级**: ⭐⭐⭐ 中高  
**预计时间**: 5-7 天  
**完成后进度**: 98%

#### 目标
持久化存储历史数据，支持趋势分析和历史查询

#### 技术方案

**方案对比**

| 特性 | SQLite | InfluxDB | PostgreSQL+TimescaleDB |
|------|--------|----------|----------------------|
| 部署复杂度 | ⭐ 极简 | ⭐⭐ 简单 | ⭐⭐⭐ 中等 |
| 性能 | ⭐⭐ 中等 | ⭐⭐⭐⭐⭐ 优秀 | ⭐⭐⭐⭐ 良好 |
| 查询能力 | ⭐⭐⭐ 标准SQL | ⭐⭐⭐⭐ InfluxQL | ⭐⭐⭐⭐⭐ 强大SQL |
| 扩展性 | ⭐⭐ 有限 | ⭐⭐⭐⭐ 良好 | ⭐⭐⭐⭐⭐ 优秀 |
| 适用场景 | 小规模/嵌入式 | 时序数据专用 | 企业级应用 |

**推荐方案**: 
- 默认使用 **SQLite**（零配置，开箱即用）
- 可选支持 **InfluxDB**（高性能，大规模部署）
- 未来扩展 **PostgreSQL + TimescaleDB**（企业级需求
