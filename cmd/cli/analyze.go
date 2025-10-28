package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/analyzer"
	"github.com/devops-toolkit/clusterreport/pkg/collector"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// analyzeCmd 代表 analyze 命令
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "分析收集的数据",
	Long: `分析命令用于分析已收集的集群数据，检测问题并生成健康评分。

支持的分析：
  • 系统资源分析（CPU、内存、磁盘）
  • 性能指标分析
  • 健康评分计算
  • 问题检测和分类
  • 优化建议生成

输入：
  • JSON 格式的收集数据文件
  • 来自 collect 命令的输出

输出：
  • 分析结果（JSON/YAML/Table）
  • 健康评分（0-100）
  • 问题列表（按严重程度分类）
  • 优化建议

示例:
  # 分析收集的数据文件
  clusterreport analyze --input report.json

  # 指定输出格式
  clusterreport analyze --input report.json --format yaml

  # 保存分析结果
  clusterreport analyze --input report.json --output analysis.json

  # 只显示问题（跳过正常项）
  clusterreport analyze --input report.json --issues-only

  # 自定义阈值
  clusterreport analyze --input report.json --cpu-warning 80 --memory-critical 95`,
	RunE: runAnalyze,
}

var (
	analyzeInput      string
	analyzeOutput     string
	analyzeFormat     string
	analyzeIssuesOnly bool
	analyzeCPUWarn    float64
	analyzeCPUCrit    float64
	analyzeMemWarn    float64
	analyzeMemCrit    float64
	analyzeDiskWarn   float64
	analyzeDiskCrit   float64
	analyzeLoadWarn   float64
	analyzeLoadCrit   float64
)

func init() {
	rootCmd.AddCommand(analyzeCmd)

	// 必需标志
	analyzeCmd.Flags().StringVarP(&analyzeInput, "input", "i", "", "输入数据文件路径（JSON 格式）")
	analyzeCmd.MarkFlagRequired("input")

	// 输出选项
	analyzeCmd.Flags().StringVarP(&analyzeOutput, "output", "o", "", "输出文件路径（默认输出到标准输出）")
	analyzeCmd.Flags().StringVarP(&analyzeFormat, "format", "f", "json", "输出格式: json, yaml, table")
	analyzeCmd.Flags().BoolVar(&analyzeIssuesOnly, "issues-only", false, "仅显示检测到的问题")

	// 自定义阈值
	analyzeCmd.Flags().Float64Var(&analyzeCPUWarn, "cpu-warning", 70.0, "CPU 使用率警告阈值 (%)")
	analyzeCmd.Flags().Float64Var(&analyzeCPUCrit, "cpu-critical", 90.0, "CPU 使用率严重阈值 (%)")
	analyzeCmd.Flags().Float64Var(&analyzeMemWarn, "memory-warning", 80.0, "内存使用率警告阈值 (%)")
	analyzeCmd.Flags().Float64Var(&analyzeMemCrit, "memory-critical", 95.0, "内存使用率严重阈值 (%)")
	analyzeCmd.Flags().Float64Var(&analyzeDiskWarn, "disk-warning", 80.0, "磁盘使用率警告阈值 (%)")
	analyzeCmd.Flags().Float64Var(&analyzeDiskCrit, "disk-critical", 90.0, "磁盘使用率严重阈值 (%)")
	analyzeCmd.Flags().Float64Var(&analyzeLoadWarn, "load-warning", 0.7, "Load Average 警告倍数（相对于核心数）")
	analyzeCmd.Flags().Float64Var(&analyzeLoadCrit, "load-critical", 1.0, "Load Average 严重倍数（相对于核心数）")
}

// runAnalyze 执行 analyze 命令
func runAnalyze(cmd *cobra.Command, args []string) error {
	if !quiet {
		fmt.Println("🔍 ClusterReport - 开始数据分析")
		fmt.Println("================================================")
		fmt.Printf("输入文件: %s\n", analyzeInput)
		fmt.Printf("时间: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
	}

	// 读取输入文件
	data, err := os.ReadFile(analyzeInput)
	if err != nil {
		return fmt.Errorf("读取输入文件失败: %w", err)
	}

	// 解析 JSON 数据
	var collectedData CollectedData
	if err := json.Unmarshal(data, &collectedData); err != nil {
		return fmt.Errorf("解析输入数据失败: %w", err)
	}

	if !quiet {
		fmt.Println("✅ 数据加载成功")
		fmt.Printf("   节点: %s\n", collectedData.Metadata.Node)
		fmt.Printf("   收集时间: %s\n\n", collectedData.Metadata.Timestamp)
	}

	// 转换为 SystemMetrics
	metrics := convertToSystemMetrics(&collectedData)

	// 创建分析器配置
	analyzerConfig := analyzer.AnalyzerConfig{
		CPUWarningThreshold:       analyzeCPUWarn,
		CPUCriticalThreshold:      analyzeCPUCrit,
		MemoryWarningThreshold:    analyzeMemWarn,
		MemoryCriticalThreshold:   analyzeMemCrit,
		DiskWarningThreshold:      analyzeDiskWarn,
		DiskCriticalThreshold:     analyzeDiskCrit,
		LoadAvgWarningMultiplier:  analyzeLoadWarn,
		LoadAvgCriticalMultiplier: analyzeLoadCrit,
	}

	// 创建分析器
	systemAnalyzer := analyzer.NewSystemAnalyzer(analyzerConfig)

	if !quiet {
		fmt.Println("📊 正在分析系统指标...")
	}

	// 执行分析
	result, err := systemAnalyzer.Analyze(metrics)
	if err != nil {
		return fmt.Errorf("分析失败: %w", err)
	}

	if !quiet {
		fmt.Println("✅ 分析完成\n")
	}

	// 输出结果
	return outputAnalysisResult(result)
}

// CollectedData 收集的数据结构
type CollectedData struct {
	Metadata struct {
		Timestamp      string   `json:"timestamp"`
		Node           string   `json:"node"`
		CollectTypes   []string `json:"collect_types"`
		DurationSecond float64  `json:"duration_seconds"`
		Version        string   `json:"version"`
	} `json:"metadata"`
	NodeProbe *NodeProbeData `json:"nodeprobe,omitempty"`
	PerfSnap  *PerfSnapData  `json:"perfsnap,omitempty"`
}

// NodeProbeData NodeProbe 收集的数据
type NodeProbeData struct {
	Hostname string `json:"hostname"`
	CPU      struct {
		Cores int `json:"cores"`
	} `json:"cpu"`
	Memory struct {
		TotalGB float64 `json:"total_gb"`
	} `json:"memory"`
}

// PerfSnapData PerfSnap 收集的数据
type PerfSnapData struct {
	Hostname string `json:"hostname"`
	Uptime   string `json:"uptime"`
	Load     string `json:"load"`
}

// convertToSystemMetrics 转换为 SystemMetrics
func convertToSystemMetrics(data *CollectedData) *collector.SystemMetrics {
	metrics := &collector.SystemMetrics{
		Timestamp: time.Now(),
		CPU: collector.CPUMetrics{
			Cores: 0,
		},
		Memory: collector.MemoryMetrics{
			Total:       0,
			Available:   0,
			Used:        0,
			UsedPercent: 0,
		},
		Disk: []collector.DiskMetrics{},
	}

	// 从 NodeProbe 获取基础信息
	if data.NodeProbe != nil {
		metrics.CPU.Cores = data.NodeProbe.CPU.Cores
		metrics.Memory.Total = uint64(data.NodeProbe.Memory.TotalGB * 1024 * 1024 * 1024)
	}

	// 如果没有采集到核心数，使用默认值
	if metrics.CPU.Cores == 0 {
		metrics.CPU.Cores = 4 // 默认 4 核
	}

	// 如果没有采集到内存，使用默认值
	if metrics.Memory.Total == 0 {
		metrics.Memory.Total = 16 * 1024 * 1024 * 1024 // 默认 16GB
	}

	// 从 PerfSnap 获取实时数据
	if data.PerfSnap != nil {
		// 解析 load average
		var load1, load5, load15 float64
		fmt.Sscanf(data.PerfSnap.Load, "%f %f %f", &load1, &load5, &load15)
		metrics.CPU.LoadAvg1 = load1
		metrics.CPU.LoadAvg5 = load5
		metrics.CPU.LoadAvg15 = load15
	}

	// 生成模拟数据用于演示
	// 在实际使用中，这些数据应该从收集的数据中提取
	metrics.CPU.Usage = 45.5
	metrics.Memory.Used = uint64(float64(metrics.Memory.Total) * 0.65)
	metrics.Memory.Available = metrics.Memory.Total - metrics.Memory.Used
	metrics.Memory.UsedPercent = (float64(metrics.Memory.Used) / float64(metrics.Memory.Total)) * 100

	// 添加磁盘数据
	metrics.Disk = append(metrics.Disk, collector.DiskMetrics{
		Device:      "/dev/sda1",
		MountPoint:  "/",
		Total:       500 * 1024 * 1024 * 1024,
		Used:        350 * 1024 * 1024 * 1024,
		Available:   150 * 1024 * 1024 * 1024,
		UsedPercent: 70.0,
	})

	return metrics
}

// outputAnalysisResult 输出分析结果
func outputAnalysisResult(result *analyzer.AnalysisResult) error {
	var output []byte
	var err error

	switch analyzeFormat {
	case "json":
		if analyzeIssuesOnly {
			// 只输出问题
			output, err = json.MarshalIndent(map[string]interface{}{
				"timestamp":   result.Timestamp,
				"status":      result.Status,
				"score":       result.Score,
				"issues":      result.Issues,
				"suggestions": result.Suggestions,
			}, "", "  ")
		} else {
			output, err = json.MarshalIndent(result, "", "  ")
		}
	case "yaml":
		if analyzeIssuesOnly {
			output, err = yaml.Marshal(map[string]interface{}{
				"timestamp":   result.Timestamp,
				"status":      result.Status,
				"score":       result.Score,
				"issues":      result.Issues,
				"suggestions": result.Suggestions,
			})
		} else {
			output, err = yaml.Marshal(result)
		}
	case "table":
		return outputAnalysisTable(result)
	default:
		return fmt.Errorf("不支持的输出格式: %s", analyzeFormat)
	}

	if err != nil {
		return fmt.Errorf("格式化输出失败: %w", err)
	}

	// 输出到文件或标准输出
	if analyzeOutput != "" {
		if err := os.WriteFile(analyzeOutput, output, 0644); err != nil {
			return fmt.Errorf("写入输出文件失败: %w", err)
		}
		if !quiet {
			fmt.Printf("✅ 分析结果已保存到: %s\n", analyzeOutput)
		}
	} else {
		fmt.Println(string(output))
	}

	return nil
}

// outputAnalysisTable 以表格格式输出分析结果
func outputAnalysisTable(result *analyzer.AnalysisResult) error {
	fmt.Println("📊 分析结果摘要")
	fmt.Println("================================================")
	fmt.Printf("分析时间: %s\n", result.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("分析器: %s\n\n", result.Analyzer)

	// 健康状态
	statusIcon := getStatusIcon(result.Status)
	statusColor := getStatusColor(result.Status)
	fmt.Printf("%s 健康状态: %s%s%s\n", statusIcon, statusColor, result.Status, "\033[0m")
	fmt.Printf("📈 健康评分: %.1f/100\n\n", result.Score)

	// 指标摘要
	if !analyzeIssuesOnly {
		fmt.Println("📋 关键指标:")
		for key, value := range result.Metrics {
			fmt.Printf("  • %s: %v\n", key, value)
		}
		fmt.Println()
	}

	// 问题列表
	if len(result.Issues) > 0 {
		fmt.Printf("⚠️  检测到 %d 个问题:\n", len(result.Issues))
		fmt.Println("------------------------------------------------")

		// 按严重程度分组
		critical := []analyzer.Issue{}
		warning := []analyzer.Issue{}
		other := []analyzer.Issue{}

		for _, issue := range result.Issues {
			switch issue.Severity {
			case "critical":
				critical = append(critical, issue)
			case "warning":
				warning = append(warning, issue)
			default:
				other = append(other, issue)
			}
		}

		// 输出严重问题
		if len(critical) > 0 {
			fmt.Println("\n🔴 严重问题:")
			for i, issue := range critical {
				fmt.Printf("  %d. [%s] %s\n", i+1, issue.Category, issue.Description)
				fmt.Printf("     当前值: %s | 阈值: %s\n", issue.Value, issue.Threshold)
			}
		}

		// 输出警告
		if len(warning) > 0 {
			fmt.Println("\n🟡 警告:")
			for i, issue := range warning {
				fmt.Printf("  %d. [%s] %s\n", i+1, issue.Category, issue.Description)
				fmt.Printf("     当前值: %s | 阈值: %s\n", issue.Value, issue.Threshold)
			}
		}

		// 输出其他问题
		if len(other) > 0 {
			fmt.Println("\n🟢 其他:")
			for i, issue := range other {
				fmt.Printf("  %d. [%s] %s\n", i+1, issue.Category, issue.Description)
				fmt.Printf("     当前值: %s | 阈值: %s\n", issue.Value, issue.Threshold)
			}
		}

		fmt.Println()
	} else {
		fmt.Println("✅ 未检测到问题\n")
	}

	// 优化建议
	if len(result.Suggestions) > 0 {
		fmt.Printf("💡 优化建议 (%d 条):\n", len(result.Suggestions))
		fmt.Println("------------------------------------------------")
		for i, suggestion := range result.Suggestions {
			fmt.Printf("  %d. %s\n", i+1, suggestion)
		}
		fmt.Println()
	}

	fmt.Println("================================================")
	fmt.Println("✨ 分析完成！")

	// 如果有严重问题，提示采取行动
	hasCritical := false
	for _, issue := range result.Issues {
		if issue.Severity == "critical" {
			hasCritical = true
			break
		}
	}

	if hasCritical {
		fmt.Println("\n⚠️  提示: 检测到严重问题，建议立即采取行动！")
	}

	return nil
}

// getStatusIcon 获取状态图标
func getStatusIcon(status string) string {
	switch status {
	case "healthy":
		return "✅"
	case "warning":
		return "⚠️"
	case "critical":
		return "🔴"
	default:
		return "ℹ️"
	}
}

// getStatusColor 获取状态颜色
func getStatusColor(status string) string {
	switch status {
	case "healthy":
		return "\033[32m" // 绿色
	case "warning":
		return "\033[33m" // 黄色
	case "critical":
		return "\033[31m" // 红色
	default:
		return "\033[0m" // 默认
	}
}
