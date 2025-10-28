# ClusterReport 下一阶段操作计划

**制定日期**: 2025-10-27  
**当前版本**: v0.7.0 (70%)  
**目标版本**: v0.8.0 (85%)  
**预计时间**: 3-5 天

---

## 🎯 阶段目标：CLI 模式增强

完善命令行工具，使其成为可独立使用的强大工具，实现本地和远程系统的采集、分析和报告生成。

## 📋 详细任务清单

### 第 1 天：Collect 命令实现

#### 任务 1.1：本地采集功能
**文件**: `cmd/cli/collect.go`

```go
实现内容：
- [x] 创建 collect 子命令
- [ ] 实现本地系统指标采集
  - 调用 pkg/collector MetricsCollector
  - 采集 CPU、内存、磁盘、网络指标
- [ ] 输出 JSON 格式数据
- [ ] 添加进度显示
- [ ] 错误处理和日志记录

预计时间：3-4 小时
```

**示例代码框架**:
```go
// cmd/cli/collect.go
func collectLocal(output string) error {
    collector := collector.NewMetricsCollector(config)
    
    // 显示进度
    fmt.Println("🔍 Collecting system metrics...")
    
    metrics, err := collector.CollectMetrics()
    if err != nil {
        return fmt.Errorf("failed to collect metrics: %w", err)
    }
    
    // 保存到文件
    return saveMetrics(metrics, output)
}
```

#### 任务 1.2：远程采集功能（SSH）
**文件**: `pkg/collector/remote.go`

```go
实现内容：
- [ ] SSH 连接管理
- [ ] 远程命令执行
- [ ] 数据传输和解析
- [ ] 连接超时处理
- [ ] 认证（密钥/密码）

预计时间：4-5 小时
```

#### 任务 1.3：批量采集功能
**文件**: `cmd/cli/collect.go`

```go
实现内容：
- [ ] 读取主机列表文件
- [ ] 并发采集（goroutine pool）
- [ ] 进度条显示
- [ ] 失败重试机制
- [ ] 结果汇总

预计时间：3-4 小时
```

---

### 第 2 天：Analyze 命令实现

#### 任务 2.1：数据分析命令
**文件**: `cmd/cli/analyze.go`

```go
实现内容：
- [ ] 创建 analyze 子命令
- [ ] 从文件加载采集数据
- [ ] 调用智能分析器
- [ ] 输出分析结果
- [ ] 支持自定义阈值

预计时间：2-3 小时
```

**示例代码**:
```go
// cmd/cli/analyze.go
func analyzeMetrics(inputFile, outputFile string, config AnalyzeConfig) error {
    // 加载指标数据
    metrics, err := loadMetrics(inputFile)
    if err != nil {
        return err
    }
    
    // 执行分析
    analyzer := analyzer.NewAnalyzer(config)
    result, err := analyzer.Analyze(metrics)
    if err != nil {
        return err
    }
    
    // 输出结果
    return saveAnalysisResult(result, outputFile)
}
```

#### 任务 2.2：配置文件支持
**文件**: `pkg/config/config.go`

```go
实现内容：
- [ ] 读取 YAML 配置文件
- [ ] 配置验证
- [ ] 默认值处理
- [ ] 环境变量覆盖

预计时间：2-3 小时
```

**配置文件示例**:
```yaml
# config.yaml
collector:
  timeout: 30s
  retry: 3
  
analyzer:
  thresholds:
    cpu_warning: 70
    cpu_critical: 90
    memory_warning: 80
    memory_critical: 95
    disk_warning: 80
    disk_critical: 90
  
  scoring:
    enabled: true
    weight:
      cpu: 30
      memory: 30
      disk: 25
      network: 15
```

---

### 第 3 天：Generate 命令实现

#### 任务 3.1：报告生成命令
**文件**: `cmd/cli/generate.go`

```go
实现内容：
- [ ] 创建 generate 子命令
- [ ] 支持多格式输出（HTML/JSON/Markdown）
- [ ] 模板选择
- [ ] 自定义标题和描述
- [ ] 美化输出

预计时间：3-4 小时
```

**示例代码**:
```go
// cmd/cli/generate.go
func generateReport(input, output, format string) error {
    // 加载分析结果
    analysis, err := loadAnalysis(input)
    if err != nil {
        return err
    }
    
    // 选择生成器
    var gen generator.Generator
    switch format {
    case "html":
        gen, _ = generator.NewHTMLGenerator()
    case "json":
        gen = generator.NewJSONGenerator()
    case "markdown":
        gen = generator.NewMarkdownGenerator()
    }
    
    // 生成报告
    report, err := gen.Generate(analysis)
    if err != nil {
        return err
    }
    
    return os.WriteFile(output, report, 0644)
}
```

#### 任务 3.2：Report 一键命令
**文件**: `cmd/cli/report.go`

```go
实现内容：
- [ ] 集成 collect + analyze + generate
- [ ] 一键生成完整报告
- [ ] 清理临时文件
- [ ] 支持邮件发送（可选）

预计时间：2-3 小时
```

---

### 第 4 天：用户体验优化

#### 任务 4.1：彩色终端输出
**文件**: `pkg/ui/color.go`

```go
实现内容：
- [ ] 使用 fatih/color 库
- [ ] 成功信息（绿色）
- [ ] 警告信息（黄色）
- [ ] 错误信息（红色）
- [ ] 信息提示（蓝色）

预计时间：2 小时
```

**示例代码**:
```go
package ui

import "github.com/fatih/color"

var (
    Success = color.New(color.FgGreen).SprintFunc()
    Warning = color.New(color.FgYellow).SprintFunc()
    Error   = color.New(color.FgRed).SprintFunc()
    Info    = color.New(color.FgBlue).SprintFunc()
)

func PrintSuccess(msg string) {
    fmt.Printf("✅ %s\n", Success(msg))
}

func PrintError(msg string) {
    fmt.Printf("❌ %s\n", Error(msg))
}
```

#### 任务 4.2：进度条显示
**文件**: `pkg/ui/progress.go`

```go
实现内容：
- [ ] 使用 schollz/progressbar 库
- [ ] 采集进度条
- [ ] 分析进度条
- [ ] 批量操作进度

预计时间：2-3 小时
```

**示例代码**:
```go
package ui

import "github.com/schollz/progressbar/v3"

func NewProgressBar(max int, description string) *progressbar.ProgressBar {
    return progressbar.NewOptions(max,
        progressbar.OptionEnableColorCodes(true),
        progressbar.OptionShowBytes(false),
        progressbar.OptionSetWidth(40),
        progressbar.OptionSetDescription(description),
        progressbar.OptionSetTheme(progressbar.Theme{
            Saucer:        "[green]=[reset]",
            SaucerHead:    "[green]>[reset]",
            SaucerPadding: " ",
            BarStart:      "[",
            BarEnd:        "]",
        }),
    )
}
```

#### 任务 4.3：日志系统
**文件**: `pkg/log/logger.go`

```go
实现内容：
- [ ] 使用 sirupsen/logrus 库
- [ ] 日志级别控制
- [ ] 文件日志输出
- [ ] 结构化日志
- [ ] 日志轮转

预计时间：2-3 小时
```

---

### 第 5 天：测试和文档

#### 任务 5.1：单元测试
**文件**: `cmd/cli/*_test.go`

```go
测试内容：
- [ ] collect 命令测试
- [ ] analyze 命令测试
- [ ] generate 命令测试
- [ ] report 命令测试
- [ ] 配置加载测试
- [ ] 错误处理测试

目标覆盖率：>80%
预计时间：4-5 小时
```

#### 任务 5.2：集成测试
**文件**: `cmd/cli/integration_test.go`

```go
测试内容：
- [ ] 端到端流程测试
- [ ] 多格式报告生成
- [ ] 批量采集测试
- [ ] 配置文件测试

预计时间：2-3 小时
```

#### 任务 5.3：文档更新
**文件**: `docs/cli-guide.md`

```markdown
文档内容：
- [ ] CLI 使用指南
- [ ] 命令参考手册
- [ ] 配置文件说明
- [ ] 常见问题（FAQ）
- [ ] 示例脚本

预计时间：2-3 小时
```

---

## 🔧 技术实现细节

### 依赖库添加

```bash
# 在 go.mod 中添加以下依赖
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/fatih/color@latest
go get github.com/schollz/progressbar/v3@latest
go get github.com/sirupsen/logrus@latest
go get golang.org/x/crypto/ssh@latest
```

### 目录结构调整

```
cmd/cli/
├── main.go              # 主入口
├── collect.go           # collect 命令
├── analyze.go           # analyze 命令
├── generate.go          # generate 命令
├── report.go            # report 命令（一键）
├── config.go            # config 子命令
└── root.go              # 根命令配置

pkg/
├── ui/                  # 用户界面
│   ├── color.go        # 彩色输出
│   ├── progress.go     # 进度条
│   └── spinner.go      # 加载动画
├── log/                 # 日志系统
│   └── logger.go
└── ssh/                 # SSH 客户端
    └── client.go
```

### 命令行参数设计

```bash
# collect 命令
clusterreport collect [flags]
  --local                    # 本地采集（默认）
  --host string              # 远程主机
  --hosts-file string        # 主机列表文件
  --user string              # SSH 用户名
  --key string               # SSH 密钥路径
  --password string          # SSH 密码
  --port int                 # SSH 端口（默认22）
  --timeout duration         # 超时时间（默认30s）
  --output string            # 输出文件
  --format string            # 输出格式（json）
  --parallel int             # 并发数（默认10）
  --verbose                  # 详细输出
  --quiet                    # 静默模式

# analyze 命令
clusterreport analyze [flags]
  --input string             # 输入文件（必需）
  --output string            # 输出文件
  --config string            # 配置文件
  --threshold float          # 全局阈值
  --cpu-warning int          # CPU 警告阈值
  --cpu-critical int         # CPU 严重阈值
  --memory-warning int       # 内存警告阈值
  --memory-critical int      # 内存严重阈值
  --verbose                  # 详细输出

# generate 命令
clusterreport generate [flags]
  --input string             # 输入文件（必需）
  --output string            # 输出文件（必需）
  --format string            # 格式（html/json/markdown）
  --template string          # 模板文件
  --title string             # 报告标题
  --description string       # 报告描述
  --open                     # 生成后自动打开

# report 命令（一键生成）
clusterreport report [flags]
  --host string              # 单个主机
  --hosts-file string        # 主机列表
  --output-dir string        # 输出目录
  --formats strings          # 输出格式列表
  --email string             # 发送邮件地址
  --config string            # 配置文件
  --parallel int             # 并发数
```

---

## ✅ 验收标准

### 功能完整性
- [ ] 本地采集正常工作
- [ ] 远程采集（SSH）正常工作
- [ ] 批量采集支持至少10个主机
- [ ] 分析功能输出正确结果
- [ ] 支持 HTML、JSON、Markdown 三种格式
- [ ] 一键报告命令可用

### 性能要求
- [ ] 单机采集时间 < 5秒
- [ ] 10台主机并发采集时间 < 30秒
- [ ] HTML 报告生成时间 < 2秒
- [ ] 内存占用 < 100MB

### 用户体验
- [ ] 命令行输出有彩色显示
- [ ] 长时间操作有进度条
- [ ] 错误信息清晰明确
- [ ] 帮助文档完整

### 代码质量
- [ ] 单元测试覆盖率 > 80%
- [ ] 所有公开函数有注释
- [ ] 没有 golint 警告
- [ ] 没有已知的 bug

---

## 📊 进度跟踪

| 日期 | 任务 | 状态 | 完成度 | 备注 |
|------|------|------|---------|------|
|
