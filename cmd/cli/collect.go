package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/collector"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	// collect 命令标志
	collectNodes        string
	collectConfig       bool
	collectPerf         bool
	collectAll          bool
	collectFlameGraph   bool
	collectAutoOptimize bool
	collectOutput       string
	collectFormat       string
	collectDuration     int
)

// collectCmd 代表 collect 命令
var collectCmd = &cobra.Command{
	Use:   "collect",
	Short: "收集节点配置和性能数据",
	Long: `收集命令用于从本地或远程节点收集配置和性能数据。

支持的收集类型：
  • 配置收集 (NodeProbe): 系统信息、硬件配置、网络、服务等
  • 性能收集 (PerfSnap): CPU、内存、磁盘IO、网络流量、进程统计等

示例:
  # 收集本地节点所有数据
  clusterreport collect --nodes localhost

  # 仅收集配置信息
  clusterreport collect --nodes localhost --collect-config

  # 仅收集性能数据
  clusterreport collect --nodes localhost --collect-perf

  # 收集性能数据并生成火焰图
  clusterreport collect --nodes localhost --collect-perf --flame-graph

  # 输出为 YAML 格式
  clusterreport collect --nodes localhost --format yaml

  # 保存到文件
  clusterreport collect --nodes localhost --output report.json`,
	RunE: runCollect,
}

func init() {
	rootCmd.AddCommand(collectCmd)

	// 必需标志
	collectCmd.Flags().StringVar(&collectNodes, "nodes", "localhost", "要收集的节点列表（逗号分隔）")

	// 收集类型标志
	collectCmd.Flags().BoolVar(&collectConfig, "collect-config", false, "仅收集配置信息 (NodeProbe)")
	collectCmd.Flags().BoolVar(&collectPerf, "collect-perf", false, "仅收集性能数据 (PerfSnap)")
	collectCmd.Flags().BoolVar(&collectAll, "collect-all", true, "收集所有数据（配置+性能）")

	// 高级选项
	collectCmd.Flags().BoolVar(&collectFlameGraph, "flame-graph", false, "生成 CPU 火焰图（需要 perf 工具）")
	collectCmd.Flags().BoolVar(&collectAutoOptimize, "auto-optimize", false, "自动执行系统优化（需要 root 权限）")
	collectCmd.Flags().IntVar(&collectDuration, "duration", 5, "性能采集持续时间（秒）")

	// 输出选项
	collectCmd.Flags().StringVarP(&collectOutput, "output", "o", "", "输出文件路径（默认输出到标准输出）")
	collectCmd.Flags().StringVarP(&collectFormat, "format", "f", "json", "输出格式: json, yaml, table")
}

// runCollect 执行 collect 命令
func runCollect(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	if !quiet {
		fmt.Println("🚀 ClusterReport - 开始数据收集")
		fmt.Println("================================================")
		fmt.Printf("节点: %s\n", collectNodes)
		fmt.Printf("时间: %s\n", startTime.Format("2006-01-02 15:04:05"))
		fmt.Println()
	}

	// 确定收集类型
	shouldCollectConfig := collectConfig || collectAll
	shouldCollectPerf := collectPerf || collectAll

	// 如果用户明确指定了某一种，则不收集另一种
	if collectConfig && !collectPerf {
		shouldCollectPerf = false
		shouldCollectConfig = true
	}
	if collectPerf && !collectConfig {
		shouldCollectConfig = false
		shouldCollectPerf = true
	}

	var result CollectResult

	// 收集配置信息（NodeProbe）
	if shouldCollectConfig {
		if !quiet {
			fmt.Println("📋 正在收集配置信息...")
		}
		configData, err := collectNodeProbeData()
		if err != nil {
			return fmt.Errorf("配置收集失败: %w", err)
		}
		result.NodeProbe = configData
		if !quiet {
			fmt.Println("✅ 配置信息收集完成")
		}
	}

	// 收集性能数据（PerfSnap）
	if shouldCollectPerf {
		if !quiet {
			fmt.Println("📊 正在收集性能数据...")
		}
		perfData, err := collectPerfSnapData()
		if err != nil {
			return fmt.Errorf("性能收集失败: %w", err)
		}
		result.PerfSnap = perfData
		if !quiet {
			fmt.Println("✅ 性能数据收集完成")
		}
	}

	// 添加元数据
	result.Metadata = CollectMetadata{
		Timestamp:    time.Now(),
		Node:         collectNodes,
		CollectTypes: getCollectTypes(shouldCollectConfig, shouldCollectPerf),
		Duration:     time.Since(startTime).Seconds(),
		Version:      "0.8.0-dev",
	}

	// 输出结果
	if err := outputResult(&result); err != nil {
		return fmt.Errorf("输出结果失败: %w", err)
	}

	if !quiet {
		fmt.Println()
		fmt.Println("================================================")
		fmt.Printf("✨ 收集完成！耗时: %.2f 秒\n", time.Since(startTime).Seconds())
		if collectOutput != "" {
			fmt.Printf("📁 结果已保存到: %s\n", collectOutput)
		}
	}

	return nil
}

// collectNodeProbeData 收集 NodeProbe 数据
func collectNodeProbeData() (*collector.NodeProbeData, error) {
	nodeProbe := collector.NewNodeProbeCollector(collectAutoOptimize)

	if verbose {
		fmt.Println("  - 收集系统信息...")
		fmt.Println("  - 收集硬件信息...")
		fmt.Println("  - 收集网络配置...")
		fmt.Println("  - 检查内核模块...")
		if collectAutoOptimize {
			fmt.Println("  - 执行自动优化...")
		}
	}

	return nodeProbe.Collect()
}

// collectPerfSnapData 收集 PerfSnap 数据
func collectPerfSnapData() (*collector.PerfSnapData, error) {
	perfSnap := collector.NewPerfSnapCollectorWithOptions(collectDuration, collectFlameGraph)

	if verbose {
		fmt.Println("  - 收集系统负载...")
		fmt.Println("  - 收集 CPU 统计...")
		fmt.Println("  - 收集内存使用...")
		fmt.Println("  - 收集磁盘 IO...")
		fmt.Println("  - 收集网络流量...")
		fmt.Println("  - 收集进程统计...")
		fmt.Println("  - 分析性能问题...")
		if collectFlameGraph {
			fmt.Println("  - 生成火焰图...")
		}
	}

	return perfSnap.Collect()
}

// CollectResult 收集结果
type CollectResult struct {
	Metadata  CollectMetadata          `json:"metadata" yaml:"metadata"`
	NodeProbe *collector.NodeProbeData `json:"nodeprobe,omitempty" yaml:"nodeprobe,omitempty"`
	PerfSnap  *collector.PerfSnapData  `json:"perfsnap,omitempty" yaml:"perfsnap,omitempty"`
}

// CollectMetadata 收集元数据
type CollectMetadata struct {
	Timestamp    time.Time `json:"timestamp" yaml:"timestamp"`
	Node         string    `json:"node" yaml:"node"`
	CollectTypes []string  `json:"collect_types" yaml:"collect_types"`
	Duration     float64   `json:"duration_seconds" yaml:"duration_seconds"`
	Version      string    `json:"version" yaml:"version"`
}

// getCollectTypes 获取收集类型列表
func getCollectTypes(config, perf bool) []string {
	var types []string
	if config {
		types = append(types, "config")
	}
	if perf {
		types = append(types, "performance")
	}
	return types
}

// outputResult 输出结果
func outputResult(result *CollectResult) error {
	var output []byte
	var err error

	switch collectFormat {
	case "json":
		output, err = json.MarshalIndent(result, "", "  ")
	case "yaml":
		output, err = yaml.Marshal(result)
	case "table":
		return outputTable(result)
	default:
		return fmt.Errorf("不支持的输出格式: %s", collectFormat)
	}

	if err != nil {
		return fmt.Errorf("格式化输出失败: %w", err)
	}

	// 输出到文件或标准输出
	if collectOutput != "" {
		if err := os.WriteFile(collectOutput, output, 0644); err != nil {
			return fmt.Errorf("写入文件失败: %w", err)
		}
	} else {
		fmt.Println(string(output))
	}

	return nil
}

// outputTable 以表格形式输出（简化版）
func outputTable(result *CollectResult) error {
	fmt.Println("\n📊 收集结果摘要")
	fmt.Println("================================================")

	// 元数据
	fmt.Printf("节点: %s\n", result.Metadata.Node)
	fmt.Printf("时间: %s\n", result.Metadata.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("耗时: %.2f 秒\n", result.Metadata.Duration)
	fmt.Printf("收集类型: %v\n", result.Metadata.CollectTypes)

	// NodeProbe 摘要
	if result.NodeProbe != nil {
		fmt.Println("\n📋 配置信息摘要:")
		fmt.Printf("  主机名: %s\n", result.NodeProbe.Hostname)
		fmt.Printf("  操作系统: %s\n", result.NodeProbe.OS)
		fmt.Printf("  内核版本: %s\n", result.NodeProbe.Kernel)
		fmt.Printf("  CPU: %s (%d 核心)\n", result.NodeProbe.CPU.Model, result.NodeProbe.CPU.Cores)
		fmt.Printf("  内存: %.2f GB\n", result.NodeProbe.Memory.TotalGB)
		fmt.Printf("  网络接口: %d 个\n", len(result.NodeProbe.Network))
		fmt.Printf("  磁盘: %d 个 (数据盘: %d)\n", result.NodeProbe.Disks.TotalDisks, result.NodeProbe.Disks.DataDiskNum)
	}

	// PerfSnap 摘要
	if result.PerfSnap != nil {
		fmt.Println("\n📊 性能数据摘要:")
		fmt.Printf("  主机名: %s\n", result.PerfSnap.Hostname)
		fmt.Printf("  运行时间: %s\n", result.PerfSnap.Uptime)
		fmt.Printf("  负载: %.2f, %.2f, %.2f\n",
			result.PerfSnap.LoadAverage.OneMin,
			result.PerfSnap.LoadAverage.FiveMin,
			result.PerfSnap.LoadAverage.FifteenMin)
		fmt.Printf("  CPU 统计: %d 个\n", len(result.PerfSnap.CPUStats))
		fmt.Printf("  磁盘 IO: %d 个设备\n", len(result.PerfSnap.DiskIOStats))
		fmt.Printf("  网络接口: %d 个\n", len(result.PerfSnap.NetworkStats))
		fmt.Printf("  检测到的问题: %d 个\n", len(result.PerfSnap.Issues))
		fmt.Printf("  优化建议: %d 条\n", len(result.PerfSnap.Recommendations))

		// 显示问题
		if len(result.PerfSnap.Issues) > 0 {
			fmt.Println("\n⚠️  检测到的问题:")
			for i, issue := range result.PerfSnap.Issues {
				if i >= 5 { // 最多显示5个
					fmt.Printf("  ... 还有 %d 个问题\n", len(result.PerfSnap.Issues)-5)
					break
				}
				fmt.Printf("  [%s] %s: %s\n", issue.Severity, issue.Category, issue.Description)
			}
		}
	}

	fmt.Println("\n💡 提示: 使用 --format json 或 --format yaml 获取完整数据")

	return nil
}
