package analyzer

import (
	"fmt"
	"math"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/collector"
)

// Analyzer 分析器接口
type Analyzer interface {
	Analyze(data interface{}) (*AnalysisResult, error)
	Name() string
}

// AnalysisResult 分析结果
type AnalysisResult struct {
	Timestamp   time.Time              `json:"timestamp"`
	Analyzer    string                 `json:"analyzer"`
	Status      string                 `json:"status"` // healthy, warning, critical
	Score       float64                `json:"score"`  // 0-100
	Issues      []Issue                `json:"issues"`
	Metrics     map[string]interface{} `json:"metrics"`
	Suggestions []string               `json:"suggestions"`
}

// Issue 问题描述
type Issue struct {
	Severity    string `json:"severity"` // low, medium, high, critical
	Category    string `json:"category"` // cpu, memory, disk, network, etc.
	Description string `json:"description"`
	Value       string `json:"value"`
	Threshold   string `json:"threshold"`
}

// BaseAnalyzer 基础分析器
type BaseAnalyzer struct {
	name   string
	config AnalyzerConfig
}

// AnalyzerConfig 分析器配置
type AnalyzerConfig struct {
	// CPU 阈值
	CPUWarningThreshold  float64 `yaml:"cpu_warning_threshold"`
	CPUCriticalThreshold float64 `yaml:"cpu_critical_threshold"`

	// 内存阈值
	MemoryWarningThreshold  float64 `yaml:"memory_warning_threshold"`
	MemoryCriticalThreshold float64 `yaml:"memory_critical_threshold"`

	// 磁盘阈值
	DiskWarningThreshold  float64 `yaml:"disk_warning_threshold"`
	DiskCriticalThreshold float64 `yaml:"disk_critical_threshold"`

	// Load Average 阈值（相对于 CPU 核心数）
	LoadAvgWarningMultiplier  float64 `yaml:"load_avg_warning_multiplier"`
	LoadAvgCriticalMultiplier float64 `yaml:"load_avg_critical_multiplier"`
}

// DefaultAnalyzerConfig 默认配置
func DefaultAnalyzerConfig() AnalyzerConfig {
	return AnalyzerConfig{
		CPUWarningThreshold:       70.0,
		CPUCriticalThreshold:      90.0,
		MemoryWarningThreshold:    80.0,
		MemoryCriticalThreshold:   95.0,
		DiskWarningThreshold:      80.0,
		DiskCriticalThreshold:     90.0,
		LoadAvgWarningMultiplier:  0.7,
		LoadAvgCriticalMultiplier: 1.0,
	}
}

// NewBaseAnalyzer 创建基础分析器
func NewBaseAnalyzer(name string, config AnalyzerConfig) *BaseAnalyzer {
	return &BaseAnalyzer{
		name:   name,
		config: config,
	}
}

// Name 返回分析器名称
func (ba *BaseAnalyzer) Name() string {
	return ba.name
}

// Analyze 实现基础分析逻辑
func (ba *BaseAnalyzer) Analyze(data interface{}) (*AnalysisResult, error) {
	metrics, ok := data.(*collector.SystemMetrics)
	if !ok {
		return nil, fmt.Errorf("invalid data type, expected *collector.SystemMetrics")
	}

	result := &AnalysisResult{
		Timestamp:   time.Now(),
		Analyzer:    ba.name,
		Status:      "healthy",
		Score:       100.0,
		Issues:      []Issue{},
		Metrics:     make(map[string]interface{}),
		Suggestions: []string{},
	}

	// 分析 CPU
	ba.analyzeCPU(metrics, result)

	// 分析内存
	ba.analyzeMemory(metrics, result)

	// 分析磁盘
	ba.analyzeDisk(metrics, result)

	// 分析 Load Average
	ba.analyzeLoadAvg(metrics, result)

	// 计算总体评分和状态
	ba.calculateOverallStatus(result)

	return result, nil
}

// analyzeCPU 分析 CPU 使用情况
func (ba *BaseAnalyzer) analyzeCPU(metrics *collector.SystemMetrics, result *AnalysisResult) {
	cpuUsage := metrics.CPU.Usage

	result.Metrics["cpu_usage"] = cpuUsage
	result.Metrics["cpu_cores"] = metrics.CPU.Cores

	if cpuUsage >= ba.config.CPUCriticalThreshold {
		result.Issues = append(result.Issues, Issue{
			Severity:    "critical",
			Category:    "cpu",
			Description: "CPU 使用率过高",
			Value:       fmt.Sprintf("%.2f%%", cpuUsage),
			Threshold:   fmt.Sprintf("%.2f%%", ba.config.CPUCriticalThreshold),
		})
		result.Score -= 30
		result.Suggestions = append(result.Suggestions, "检查 CPU 密集型进程，考虑优化或扩容")
	} else if cpuUsage >= ba.config.CPUWarningThreshold {
		result.Issues = append(result.Issues, Issue{
			Severity:    "warning",
			Category:    "cpu",
			Description: "CPU 使用率偏高",
			Value:       fmt.Sprintf("%.2f%%", cpuUsage),
			Threshold:   fmt.Sprintf("%.2f%%", ba.config.CPUWarningThreshold),
		})
		result.Score -= 15
		result.Suggestions = append(result.Suggestions, "监控 CPU 使用趋势，准备扩容计划")
	}
}

// analyzeMemory 分析内存使用情况
func (ba *BaseAnalyzer) analyzeMemory(metrics *collector.SystemMetrics, result *AnalysisResult) {
	memUsage := metrics.Memory.UsedPercent

	result.Metrics["memory_usage_percent"] = memUsage
	result.Metrics["memory_total"] = metrics.Memory.Total
	result.Metrics["memory_available"] = metrics.Memory.Available

	if memUsage >= ba.config.MemoryCriticalThreshold {
		result.Issues = append(result.Issues, Issue{
			Severity:    "critical",
			Category:    "memory",
			Description: "内存使用率严重过高",
			Value:       fmt.Sprintf("%.2f%%", memUsage),
			Threshold:   fmt.Sprintf("%.2f%%", ba.config.MemoryCriticalThreshold),
		})
		result.Score -= 30
		result.Suggestions = append(result.Suggestions, "立即检查内存泄漏，清理缓存或增加内存")
	} else if memUsage >= ba.config.MemoryWarningThreshold {
		result.Issues = append(result.Issues, Issue{
			Severity:    "warning",
			Category:    "memory",
			Description: "内存使用率偏高",
			Value:       fmt.Sprintf("%.2f%%", memUsage),
			Threshold:   fmt.Sprintf("%.2f%%", ba.config.MemoryWarningThreshold),
		})
		result.Score -= 15
		result.Suggestions = append(result.Suggestions, "关注内存使用趋势，优化应用内存占用")
	}
}

// analyzeDisk 分析磁盘使用情况
func (ba *BaseAnalyzer) analyzeDisk(metrics *collector.SystemMetrics, result *AnalysisResult) {
	for _, disk := range metrics.Disk {
		diskUsage := disk.UsedPercent

		metricKey := fmt.Sprintf("disk_usage_%s", disk.MountPoint)
		result.Metrics[metricKey] = diskUsage

		if diskUsage >= ba.config.DiskCriticalThreshold {
			result.Issues = append(result.Issues, Issue{
				Severity:    "critical",
				Category:    "disk",
				Description: fmt.Sprintf("磁盘 %s 使用率严重过高", disk.MountPoint),
				Value:       fmt.Sprintf("%.2f%%", diskUsage),
				Threshold:   fmt.Sprintf("%.2f%%", ba.config.DiskCriticalThreshold),
			})
			result.Score -= 20
			result.Suggestions = append(result.Suggestions,
				fmt.Sprintf("清理 %s 上的日志和临时文件，或扩展磁盘空间", disk.MountPoint))
		} else if diskUsage >= ba.config.DiskWarningThreshold {
			result.Issues = append(result.Issues, Issue{
				Severity:    "warning",
				Category:    "disk",
				Description: fmt.Sprintf("磁盘 %s 使用率偏高", disk.MountPoint),
				Value:       fmt.Sprintf("%.2f%%", diskUsage),
				Threshold:   fmt.Sprintf("%.2f%%", ba.config.DiskWarningThreshold),
			})
			result.Score -= 10
			result.Suggestions = append(result.Suggestions,
				fmt.Sprintf("监控 %s 磁盘使用，计划清理策略", disk.MountPoint))
		}
	}
}

// analyzeLoadAvg 分析系统负载
func (ba *BaseAnalyzer) analyzeLoadAvg(metrics *collector.SystemMetrics, result *AnalysisResult) {
	cores := float64(metrics.CPU.Cores)
	load1 := metrics.CPU.LoadAvg1
	load5 := metrics.CPU.LoadAvg5
	load15 := metrics.CPU.LoadAvg15

	result.Metrics["load_avg_1"] = load1
	result.Metrics["load_avg_5"] = load5
	result.Metrics["load_avg_15"] = load15
	result.Metrics["load_per_core_1"] = load1 / cores

	// 使用 1 分钟负载平均值进行评估
	loadPerCore := load1 / cores

	criticalThreshold := ba.config.LoadAvgCriticalMultiplier
	warningThreshold := ba.config.LoadAvgWarningMultiplier

	if loadPerCore >= criticalThreshold {
		result.Issues = append(result.Issues, Issue{
			Severity:    "critical",
			Category:    "load",
			Description: "系统负载过高",
			Value:       fmt.Sprintf("%.2f (%.2f per core)", load1, loadPerCore),
			Threshold:   fmt.Sprintf("%.2f per core", criticalThreshold),
		})
		result.Score -= 25
		result.Suggestions = append(result.Suggestions, "系统负载严重，检查运行进程和 I/O 等待")
	} else if loadPerCore >= warningThreshold {
		result.Issues = append(result.Issues, Issue{
			Severity:    "warning",
			Category:    "load",
			Description: "系统负载偏高",
			Value:       fmt.Sprintf("%.2f (%.2f per core)", load1, loadPerCore),
			Threshold:   fmt.Sprintf("%.2f per core", warningThreshold),
		})
		result.Score -= 10
		result.Suggestions = append(result.Suggestions, "关注系统负载趋势")
	}
}

// calculateOverallStatus 计算总体状态
func (ba *BaseAnalyzer) calculateOverallStatus(result *AnalysisResult) {
	// 确保评分在 0-100 范围内
	result.Score = math.Max(0, math.Min(100, result.Score))

	// 根据问题严重程度确定状态
	hasCritical := false
	hasWarning := false

	for _, issue := range result.Issues {
		if issue.Severity == "critical" {
			hasCritical = true
		} else if issue.Severity == "warning" {
			hasWarning = true
		}
	}

	if hasCritical || result.Score < 60 {
		result.Status = "critical"
	} else if hasWarning || result.Score < 80 {
		result.Status = "warning"
	} else {
		result.Status = "healthy"
	}
}

// SystemAnalyzer 系统分析器
type SystemAnalyzer struct {
	*BaseAnalyzer
}

// NewSystemAnalyzer 创建系统分析器
func NewSystemAnalyzer(config AnalyzerConfig) *SystemAnalyzer {
	return &SystemAnalyzer{
		BaseAnalyzer: NewBaseAnalyzer("system-analyzer", config),
	}
}
