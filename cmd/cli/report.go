package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

// reportCmd 代表 report 命令
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "一键式生成完整报告",
	Long: `report 命令是一个一键式命令，自动执行完整的工作流程：
  1. 收集数据（collect）
  2. 分析数据（analyze）
  3. 生成报告（generate）

这个命令简化了整个流程，适合快速生成集群健康报告。

特性：
  • 自动化工作流
  • 进度跟踪
  • 错误恢复
  • 中间文件管理

示例:
  # 为本地节点生成报告
  clusterreport report --nodes localhost

  # 为远程节点生成报告
  clusterreport report --nodes server1,server2,server3

  # 自定义输出目录
  clusterreport report --nodes localhost --output-dir ./my-reports

  # 指定报告格式
  clusterreport report --nodes localhost --format all

  # 保留中间文件
  clusterreport report --nodes localhost --keep-intermediate`,
	RunE: runReport,
}

var (
	reportNodes            string
	reportOutputDir        string
	reportFormat           string
	reportKeepIntermediate bool
	reportTitle            string
)

func init() {
	rootCmd.AddCommand(reportCmd)

	// 必需标志
	reportCmd.Flags().StringVar(&reportNodes, "nodes", "", "目标节点列表（逗号分隔）")
	reportCmd.MarkFlagRequired("nodes")

	// 输出选项
	reportCmd.Flags().StringVar(&reportOutputDir, "output-dir", "./cluster-reports", "输出目录")
	reportCmd.Flags().StringVarP(&reportFormat, "format", "f", "all", "报告格式: html, markdown, all")

	// 其他选项
	reportCmd.Flags().BoolVar(&reportKeepIntermediate, "keep-intermediate", false, "保留中间文件（collect 和 analyze 的输出）")
	reportCmd.Flags().StringVar(&reportTitle, "title", "", "报告标题")
}

// runReport 执行 report 命令
func runReport(cmd *cobra.Command, args []string) error {
	startTime := time.Now()

	if !quiet {
		fmt.Println("🚀 ClusterReport - 一键式报告生成")
		fmt.Println("================================================")
		fmt.Printf("节点: %s\n", reportNodes)
		fmt.Printf("输出目录: %s\n", reportOutputDir)
		fmt.Printf("报告格式: %s\n", reportFormat)
		fmt.Printf("开始时间: %s\n\n", startTime.Format("2006-01-02 15:04:05"))
	}

	// 创建输出目录
	if err := os.MkdirAll(reportOutputDir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 定义中间文件路径
	timestamp := time.Now().Format("20060102-150405")
	collectFile := filepath.Join(reportOutputDir, fmt.Sprintf("collect-%s.json", timestamp))
	analyzeFile := filepath.Join(reportOutputDir, fmt.Sprintf("analysis-%s.json", timestamp))

	// 步骤 1: 收集数据
	if !quiet {
		fmt.Println("📊 步骤 1/3: 收集数据")
		fmt.Println("------------------------------------------------")
	}

	collectCmd := fmt.Sprintf("collect --nodes %s --output %s", reportNodes, collectFile)
	if verbose {
		collectCmd += " --verbose"
	} else if quiet {
		collectCmd += " --quiet"
	}

	if err := executeSubCommand(collectCmd); err != nil {
		return fmt.Errorf("数据收集失败: %w", err)
	}

	if !quiet {
		fmt.Println("✅ 数据收集完成\n")
	}

	// 步骤 2: 分析数据
	if !quiet {
		fmt.Println("🔍 步骤 2/3: 分析数据")
		fmt.Println("------------------------------------------------")
	}

	analyzeCmd := fmt.Sprintf("analyze --input %s --output %s --format json", collectFile, analyzeFile)
	if verbose {
		analyzeCmd += " --verbose"
	} else if quiet {
		analyzeCmd += " --quiet"
	}

	if err := executeSubCommand(analyzeCmd); err != nil {
		return fmt.Errorf("数据分析失败: %w", err)
	}

	if !quiet {
		fmt.Println("✅ 数据分析完成\n")
	}

	// 步骤 3: 生成报告
	if !quiet {
		fmt.Println("📝 步骤 3/3: 生成报告")
		fmt.Println("------------------------------------------------")
	}

	reportDir := filepath.Join(reportOutputDir, fmt.Sprintf("report-%s", timestamp))
	generateCmd := fmt.Sprintf("generate --input %s --format %s --output-dir %s", analyzeFile, reportFormat, reportDir)
	if reportTitle != "" {
		generateCmd += fmt.Sprintf(" --title \"%s\"", reportTitle)
	}
	if verbose {
		generateCmd += " --verbose"
	} else if quiet {
		generateCmd += " --quiet"
	}

	if err := executeSubCommand(generateCmd); err != nil {
		return fmt.Errorf("报告生成失败: %w", err)
	}

	if !quiet {
		fmt.Println("✅ 报告生成完成\n")
	}

	// 清理中间文件（如果需要）
	if !reportKeepIntermediate {
		if !quiet {
			fmt.Println("🧹 清理中间文件...")
		}
		os.Remove(collectFile)
		os.Remove(analyzeFile)
		if !quiet {
			fmt.Println("✅ 中间文件已清理\n")
		}
	} else {
		if !quiet {
			fmt.Println("📦 保留中间文件:")
			fmt.Printf("   收集数据: %s\n", collectFile)
			fmt.Printf("   分析结果: %s\n\n", analyzeFile)
		}
	}

	// 完成总结
	duration := time.Since(startTime)
	if !quiet {
		fmt.Println("================================================")
		fmt.Println("✨ 报告生成完成！")
		fmt.Printf("总耗时: %.2f 秒\n\n", duration.Seconds())

		fmt.Println("📁 输出文件:")
		fmt.Printf("   报告目录: %s\n", reportDir)

		// 列出生成的报告文件
		files, err := os.ReadDir(reportDir)
		if err == nil {
			for _, file := range files {
				if !file.IsDir() {
					fmt.Printf("   - %s\n", file.Name())
				}
			}
		}

		fmt.Println("\n💡 提示:")
		if reportFormat == "html" || reportFormat == "all" {
			htmlFile := filepath.Join(reportDir, "report.html")
			fmt.Printf("   打开 HTML 报告: open %s\n", htmlFile)
		}
		if reportFormat == "markdown" || reportFormat == "all" {
			mdFile := filepath.Join(reportDir, "report.md")
			fmt.Printf("   查看 Markdown 报告: cat %s\n", mdFile)
		}
	}

	return nil
}

// executeSubCommand 执行子命令
func executeSubCommand(cmdStr string) error {
	// 这里我们实际上需要直接调用对应的函数，而不是创建新的进程
	// 为了简化，我们使用一个辅助函数来解析并执行命令

	// 保存当前的标志值
	oldQuiet := quiet
	oldVerbose := verbose

	defer func() {
		quiet = oldQuiet
		verbose = oldVerbose
	}()

	// 解析命令字符串并执行相应的函数
	// 注意：这是一个简化的实现，实际项目中可能需要更复杂的命令解析
	args := parseCommandString(cmdStr)
	if len(args) == 0 {
		return fmt.Errorf("无效的命令")
	}

	subCmd := args[0]
	subArgs := args[1:]

	// 根据子命令执行相应的函数
	switch subCmd {
	case "collect":
		return runCollectWithArgs(subArgs)
	case "analyze":
		return runAnalyzeWithArgs(subArgs)
	case "generate":
		return runGenerateWithArgs(subArgs)
	default:
		return fmt.Errorf("未知的子命令: %s", subCmd)
	}
}

// parseCommandString 解析命令字符串
func parseCommandString(cmdStr string) []string {
	// 简单的空格分割，实际项目中可能需要更复杂的解析
	args := []string{}
	current := ""
	inQuotes := false

	for _, char := range cmdStr {
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				current += string(char)
			} else if current != "" {
				args = append(args, current)
				current = ""
			}
		default:
			current += string(char)
		}
	}

	if current != "" {
		args = append(args, current)
	}

	return args
}

// runCollectWithArgs 使用参数运行 collect
func runCollectWithArgs(args []string) error {
	// 解析参数
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--nodes":
			if i+1 < len(args) {
				collectNodes = args[i+1]
				i++
			}
		case "--output":
			if i+1 < len(args) {
				collectOutput = args[i+1]
				i++
			}
		case "--quiet":
			quiet = true
		case "--verbose":
			verbose = true
		}
	}
	return runCollect(nil, nil)
}

// runAnalyzeWithArgs 使用参数运行 analyze
func runAnalyzeWithArgs(args []string) error {
	// 解析参数
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--input":
			if i+1 < len(args) {
				analyzeInput = args[i+1]
				i++
			}
		case "--output":
			if i+1 < len(args) {
				analyzeOutput = args[i+1]
				i++
			}
		case "--format":
			if i+1 < len(args) {
				analyzeFormat = args[i+1]
				i++
			}
		case "--quiet":
			quiet = true
		case "--verbose":
			verbose = true
		}
	}
	return runAnalyze(nil, nil)
}

// runGenerateWithArgs 使用参数运行 generate
func runGenerateWithArgs(args []string) error {
	// 解析参数
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--input":
			if i+1 < len(args) {
				generateInput = args[i+1]
				i++
			}
		case "--format":
			if i+1 < len(args) {
				generateFormat = args[i+1]
				i++
			}
		case "--output-dir":
			if i+1 < len(args) {
				generateOutputDir = args[i+1]
				i++
			}
		case "--title":
			if i+1 < len(args) {
				generateTitle = args[i+1]
				i++
			}
		case "--quiet":
			quiet = true
		case "--verbose":
			verbose = true
		}
	}
	return runGenerate(nil, nil)
}
