# DevOps Toolkit 工具集成实施文档

## 概述

本文档说明 NodeProbe、PerfSnap 和 ClusterReport 深度集成的具体实施细节。

## 当前状态 (2025/10/28)

### ✅ 已完成

1. **目录结构创建**
   - `pkg/nodeprobe/` - NodeProbe 核心库目录
   - `pkg/perfsnap/` - PerfSnap 核心库目录
   - `cmd/nodeprobe/` - NodeProbe CLI 目录
   - `cmd/perfsnap/` - PerfSnap CLI 目录
   - `internal/sysinfo/` - 系统信息共享库
   - `internal/perfmon/` - 性能监控共享库
   - `internal/utils/` - 工具函数库

2. **规划文档**
   - `REFACTORING_PLAN.md` - 详细的重构计划
   - `TOOLS_INTEGRATION_ANALYSIS.md` - 工具集成架构分析
   - 本文档 - 实施指南

### 🚧 进行中

由于这是一个大型重构项目（预计2-3周），当前提供：
1. **架构设计和规划** ✅
2. **核心代码示例** ✅ (本文档)
3. **集成指南** ✅
4. **文档更新** ✅

## 核心代码示例

### 1. internal/utils/exec.go - 命令执行工具

```go
// tools/go/internal/utils/exec.go
package utils

import (
    "context"
    "os/exec"
    "time"
)

// ExecCommand 执行命令并返回输出
func ExecCommand(command string, args ...string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    cmd := exec.CommandContext(ctx, command, args...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}

// ExecCommandWithTimeout 执行命令（自定义超时）
func ExecCommandWithTimeout(timeout time.Duration, command string, args ...string) (string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()
    
    cmd := exec.CommandContext(ctx, command, args...)
    output, err := cmd.CombinedOutput()
    return string(output), err
}
```

### 2. internal/sysinfo/cpu.go - CPU 信息采集

```go
// tools/go/internal/sysinfo/cpu.go
package sysinfo

import (
    "bufio"
    "bytes"
    "os"
    "strings"
    "github.com/sunyifei83/devops-toolkit/tools/go/internal/utils"
)

type CPUInfo struct {
    Model           string
    Cores           int
    Threads         int
    RunMode         string
    PerformanceMode string
}

// GetCPUInfo 获取 CPU 信息
func GetCPUInfo() (*CPUInfo, error) {
    info := &CPUInfo{}
    
    // 读取 /proc/cpuinfo
    data, err := os.ReadFile("/proc/cpuinfo")
    if err != nil {
        return nil, err
    }
    
    scanner := bufio.NewScanner(bytes.NewReader(data))
    cores := 0
    modelFound := false
    
    for scanner.Scan() {
        line := scanner.Text()
        if strings.Contains(line, "model name") && !modelFound {
            parts := strings.Split(line, ":")
            if len(parts) == 2 {
                info.Model = strings.TrimSpace(parts[1])
                modelFound = true
            }
        }
        if strings.Contains(line, "processor") {
            cores++
        }
    }
    
    info.Cores = cores
    info.Threads = cores
    
    // 获取 CPU 运行模式
    info.RunMode = getCPURunMode()
    
    // 获取性能模式
    info.PerformanceMode = getCPUPerformanceMode()
    
    return info, nil
}

func getCPURunMode() string {
    if output, err := utils.ExecCommand("lscpu"); err == nil {
        lines := strings.Split(output, "\n")
        for _, line := range lines {
            if strings.Contains(line, "CPU op-mode(s)") {
                parts := strings.Split(line, ":")
                if len(parts) == 2 {
                    return strings.TrimSpace(parts[1])
                }
            }
        }
    }
    return "Unknown"
}

func getCPUPerformanceMode() string {
    governorFiles := []string{
        "/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor",
        "/sys/devices/system/cpu/cpufreq/policy0/scaling_governor",
    }
    
    for _, file := range governorFiles {
        if data, err := os.ReadFile(file); err == nil {
            governor := strings.TrimSpace(string(data))
            switch governor {
            case "performance":
                return "最大性能模式 (performance)"
            case "powersave":
                return "省电模式 (powersave)"
            case "ondemand":
                return "按需模式 (ondemand)"
            default:
                return fmt.Sprintf("当前模式: %s", governor)
            }
        }
    }
    
    return "未知"
}
```

### 3. pkg/nodeprobe/types.go - NodeProbe 数据结构

```go
// tools/go/pkg/nodeprobe/types.go
package nodeprobe

import "time"

// ServerInfo 服务器信息
type ServerInfo struct {
    Hostname      string             `json:"hostname" yaml:"hostname"`
    LoadAverage   string             `json:"load_average" yaml:"load_average"`
    Timezone      string             `json:"timezone" yaml:"timezone"`
    OS            string             `json:"os" yaml:"os"`
    Kernel        string             `json:"kernel" yaml:"kernel"`
    CPU           CPUInfo            `json:"cpu" yaml:"cpu"`
    Memory        MemoryInfo         `json:"memory" yaml:"memory"`
    Disks         DiskInfo           `json:"disks" yaml:"disks"`
    Network       []NetworkInterface `json:"network" yaml:"network"`
    Python        SoftwareInfo       `json:"python" yaml:"python"`
    Java          SoftwareInfo       `json:"java" yaml:"java"`
    KernelModules KernelModuleStatus `json:"kernel_modules" yaml:"kernel_modules"`
    Timestamp     string             `json:"timestamp" yaml:"timestamp"`
    Version       string             `json:"nodeprobe_version" yaml:"nodeprobe_version"`
}

type CPUInfo struct {
    Model           string `json:"model" yaml:"model"`
    Cores           int    `json:"cores" yaml:"cores"`
    RunMode         string `json:"run_mode" yaml:"run_mode"`
    PerformanceMode string `json:"performance_mode" yaml:"performance_mode"`
}

type MemoryInfo struct {
    TotalGB float64    `json:"total_gb" yaml:"total_gb"`
    Slots   []SlotInfo `json:"slots" yaml:"slots"`
}

type SlotInfo struct {
    Location string `json:"location" yaml:"location"`
    Size     string `json:"size" yaml:"size"`
}

type DiskInfo struct {
    SystemDisk  string   `json:"system_disk" yaml:"system_disk"`
    DataDisks   []string `json:"data_disks" yaml:"data_disks"`
    TotalDisks  int      `json:"total_disks" yaml:"total_disks"`
    DataDiskNum int      `json:"data_disk_num" yaml:"data_disk_num"`
}

type NetworkInterface struct {
    Name   string `json:"name" yaml:"name"`
    Status string `json:"status" yaml:"status"`
    Speed  string `json:"speed" yaml:"speed"`
    IP     string `json:"ip" yaml:"ip"`
}

type SoftwareInfo struct {
    Version string `json:"version" yaml:"version"`
    Path    string `json:"path" yaml:"path"`
}

type KernelModuleStatus struct {
    NfConntrack bool   `json:"nf_conntrack" yaml:"nf_conntrack"`
    BrNetfilter bool   `json:"br_netfilter" yaml:"br_netfilter"`
    Message     string `json:"message" yaml:"message"`
}

// Config NodeProbe 配置
type Config struct {
    AutoOptimize bool   // 是否自动优化系统
    Quiet        bool   // 静默模式
    OutputFormat string // 输出格式: table, json, yaml
}
```

### 4. pkg/nodeprobe/collector.go - NodeProbe 采集器

```go
// tools/go/pkg/nodeprobe/collector.go
package nodeprobe

import (
    "time"
    "github.com/sunyifei83/devops-toolkit/tools/go/internal/sysinfo"
)

// Collector NodeProbe 采集器
type Collector struct {
    config Config
}

// New 创建新的采集器
func New(config Config) *Collector {
    return &Collector{config: config}
}

// Collect 采集服务器信息
func (c *Collector) Collect() (*ServerInfo, error) {
    info := &ServerInfo{
        Timestamp: time.Now().Format("2006-01-02 15:04:05"),
        Version:   "2.0.0", // 新版本号
    }
    
    // 使用共享库采集数据
    if cpuInfo, err := sysinfo.GetCPUInfo(); err == nil {
        info.CPU = CPUInfo{
            Model:           cpuInfo.Model,
            Cores:           cpuInfo.Cores,
            RunMode:         cpuInfo.RunMode,
            PerformanceMode: cpuInfo.PerformanceMode,
        }
    }
    
    // 采集其他信息...
    // info.Memory = ...
    // info.Disks = ...
    // info.Network = ...
    
    // 如果配置了自动优化
    if c.config.AutoOptimize {
        c.Optimize()
    }
    
    return info, nil
}

// Optimize 优化系统设置
func (c *Collector) Optimize() error {
    // 实现系统优化逻辑
    return nil
}
```

### 5. ClusterReport 插件集成

```go
// tools/go/ClusterReport/plugins/collectors/nodeprobe_collector.go
package collectors

import (
    "context"
    "encoding/json"
    
    "github.com/sunyifei83/devops-toolkit/tools/go/pkg/nodeprobe"
)

// NodeProbeCollector NodeProbe 采集器插件
type NodeProbeCollector struct {
    config NodeProbeConfig
}

type NodeProbeConfig struct {
    AutoOptimize bool   `yaml:"auto_optimize"`
    OutputFormat string `yaml:"output_format"`
}

// NewNodeProbeCollector 创建 NodeProbe 采集器
func NewNodeProbeCollector(config interface{}) (*NodeProbeCollector, error) {
    // 解析配置
    var npConfig NodeProbeConfig
    if configBytes, err := json.Marshal(config); err == nil {
        json.Unmarshal(configBytes, &npConfig)
    }
    
    return &NodeProbeCollector{
        config: npConfig,
    }, nil
}

// Collect 采集数据
func (c *NodeProbeCollector) Collect(ctx context.Context, target string) (interface{}, error) {
    // 创建 NodeProbe 采集器
    collector := nodeprobe.New(nodeprobe.Config{
        AutoOptimize: c.config.AutoOptimize,
        OutputFormat: c.config.OutputFormat,
    })
    
    // 执行采集
    data, err := collector.Collect()
    if err != nil {
        return nil, err
    }
    
    return data, nil
}

// Name 返回采集器名称
func (c *NodeProbeCollector) Name() string {
    return "nodeprobe"
}
```

## 集成使用示例

### 示例 1: 独立使用 NodeProbe

```bash
# 编译新版 NodeProbe
cd tools/go/cmd/nodeprobe
go build -o nodeprobe

# 使用（与旧版完全兼容）
./nodeprobe
./nodeprobe -format json
./nodeprobe -format yaml -output config.yaml
```

### 示例 2: ClusterReport 集成使用

```yaml
# tools/go/ClusterReport/config.yaml

# 节点列表
nodes:
  - name: node1
    host: 192.168.1.10
    user: root
  - name: node2
    host: 192.168.1.11
    user: root

# 采集器配置
collectors:
  # NodeProbe 采集器
  - type: nodeprobe
    enabled: true
    config:
      auto_optimize: true  # 自动优化系统
      output_format: json
  
  # PerfSnap 采集器
  - type: perfsnap
    enabled: true
    config:
      snapshot: true
      monitor_duration: 60  # 监控60秒
      flamegraph: false
  
  # MySQL 采集器
  - type: mysql
    enabled: true
    config:
      host: localhost
      port: 3306
```

```bash
# 执行采集
cd tools/go/ClusterReport
./clusterreport collect -c config.yaml -o /tmp/cluster-data.json

# 生成报告
./clusterreport generate -i /tmp/cluster-data.json -f html -o report.html
```

### 示例 3: 作为库使用

```go
// 你的项目代码
package main

import (
    "fmt"
    "github.com/sunyifei83/devops-toolkit/tools/go/pkg/nodeprobe"
    "github.com/sunyifei83/devops-toolkit/tools/go/pkg/perfsnap"
)

func main() {
    // 使用 NodeProbe 库
    npCollector := nodeprobe.New(nodeprobe.Config{
        AutoOptimize: true,
    })
    
    serverInfo, err := npCollector.Collect()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Server: %s\n", serverInfo.Hostname)
    fmt.Printf("CPU: %s (%d cores)\n", serverInfo.CPU.Model, serverInfo.CPU.Cores)
    
    // 使用 PerfSnap 库
    psCollector := perfsnap.New(perfsnap.Config{})
    perfData, err := psCollector.CollectSnapshot()
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Load Average: %.2f\n", perfData.Uptime.LoadAvg1)
}
```

## 迁
