package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/analyzer"
	"github.com/spf13/cobra"
)

// generateCmd 代表 generate 命令
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "从分析结果生成报告",
	Long: `generate 命令用于从分析结果生成各种格式的报告。

支持的报告格式：
  • HTML  - 交互式网页报告
  • Markdown - 文本格式报告

报告内容：
  • 执行摘要
  • 健康评分和状态
  • 详细的问题列表
  • 系统指标
  • 优化建议

示例:
  # 生成 HTML 报告
  clusterreport generate --input analysis.json --format html --output report.html

  # 生成 Markdown 报告
  clusterreport generate --input analysis.json --format markdown --output report.md

  # 生成所有格式
  clusterreport generate --input analysis.json --format all --output-dir ./reports`,
	RunE: runGenerate,
}

var (
	generateInput     string
	generateOutput    string
	generateOutputDir string
	generateFormat    string
	generateTitle     string
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&generateInput, "input", "i", "", "输入分析结果文件（JSON 格式）")
	generateCmd.MarkFlagRequired("input")

	generateCmd.Flags().StringVarP(&generateOutput, "output", "o", "", "输出文件路径")
	generateCmd.Flags().StringVar(&generateOutputDir, "output-dir", "./reports", "输出目录")
	generateCmd.Flags().StringVarP(&generateFormat, "format", "f", "html", "报告格式: html, markdown, all")
	generateCmd.Flags().StringVar(&generateTitle, "title", "", "报告标题")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	if !quiet {
		fmt.Println("📝 ClusterReport - 开始生成报告")
		fmt.Println("================================================")
		fmt.Printf("输入文件: %s\n", generateInput)
		fmt.Printf("报告格式: %s\n\n", generateFormat)
	}

	data, err := os.ReadFile(generateInput)
	if err != nil {
		return fmt.Errorf("读取输入文件失败: %w", err)
	}

	var result analyzer.AnalysisResult
	if err := json.Unmarshal(data, &result); err != nil {
		return fmt.Errorf("解析分析结果失败: %w", err)
	}

	if !quiet {
		fmt.Printf("✅ 分析结果加载成功 (状态: %s, 评分: %.1f)\n\n", result.Status, result.Score)
	}

	if generateTitle == "" {
		generateTitle = fmt.Sprintf("集群健康分析报告 - %s", time.Now().Format("2006-01-02"))
	}

	switch generateFormat {
	case "html":
		return generateHTMLReport(&result)
	case "markdown", "md":
		return generateMarkdownReport(&result)
	case "all":
		return generateAllReports(&result)
	default:
		return fmt.Errorf("不支持的报告格式: %s", generateFormat)
	}
}

func generateHTMLReport(result *analyzer.AnalysisResult) error {
	if !quiet {
		fmt.Println("📄 正在生成 HTML 报告...")
	}

	outputFile := generateOutput
	if outputFile == "" {
		if err := os.MkdirAll(generateOutputDir, 0755); err != nil {
			return err
		}
		outputFile = filepath.Join(generateOutputDir, "report.html")
	}

	html := buildHTMLReport(result)
	if err := os.WriteFile(outputFile, []byte(html), 0644); err != nil {
		return err
	}

	if !quiet {
		fmt.Printf("✅ HTML 报告已生成: %s\n", outputFile)
	}
	return nil
}

func generateMarkdownReport(result *analyzer.AnalysisResult) error {
	if !quiet {
		fmt.Println("📄 正在生成 Markdown 报告...")
	}

	outputFile := generateOutput
	if outputFile == "" {
		if err := os.MkdirAll(generateOutputDir, 0755); err != nil {
			return err
		}
		outputFile = filepath.Join(generateOutputDir, "report.md")
	}

	markdown := buildMarkdownReport(result)
	if err := os.WriteFile(outputFile, []byte(markdown), 0644); err != nil {
		return err
	}

	if !quiet {
		fmt.Printf("✅ Markdown 报告已生成: %s\n", outputFile)
	}
	return nil
}

func generateAllReports(result *analyzer.AnalysisResult) error {
	if err := os.MkdirAll(generateOutputDir, 0755); err != nil {
		return err
	}

	// HTML
	htmlFile := filepath.Join(generateOutputDir, "report.html")
	if err := os.WriteFile(htmlFile, []byte(buildHTMLReport(result)), 0644); err != nil {
		return err
	}

	// Markdown
	mdFile := filepath.Join(generateOutputDir, "report.md")
	if err := os.WriteFile(mdFile, []byte(buildMarkdownReport(result)), 0644); err != nil {
		return err
	}

	if !quiet {
		fmt.Printf("✅ HTML: %s\n", htmlFile)
		fmt.Printf("✅ Markdown: %s\n", mdFile)
		fmt.Printf("\n✨ 所有报告已生成到: %s\n", generateOutputDir)
	}
	return nil
}

func buildHTMLReport(result *analyzer.AnalysisResult) string {
	var html strings.Builder
	statusText := getStatusTextChinese(result.Status)

	html.WriteString(fmt.Sprintf(`<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s</title>
<style>
*{margin:0;padding:0;box-sizing:border-box}
body{font-family:-apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,Arial,sans-serif;line-height:1.6;color:#333;background:#f5f5f5;padding:20px}
.container{max-width:1200px;margin:0 auto;background:white;border-radius:8px;box-shadow:0 2px 8px rgba(0,0,0,0.1);overflow:hidden}
.header{background:linear-gradient(135deg,#667eea 0%%,#764ba2 100%%);color:white;padding:40px;text-align:center}
.header h1{font-size:32px;margin-bottom:10px}
.header .timestamp{opacity:0.9;font-size:14px}
.summary{display:grid;grid-template-columns:repeat(auto-fit,minmax(250px,1fr));gap:20px;padding:40px;background:#f8f9fa}
.summary-card{background:white;padding:20px;border-radius:8px;box-shadow:0 2px 4px rgba(0,0,0,0.05);text-align:center}
.summary-card h3{color:#666;font-size:14px;font-weight:normal;margin-bottom:10px}
.summary-card .value{font-size:36px;font-weight:bold;margin-bottom:5px}
.status-healthy{color:#28a745}
.status-warning{color:#ffc107}
.status-critical{color:#dc3545}
.content{padding:40px}
.section{margin-bottom:40px}
.section h2{font-size:24px;margin-bottom:20px;padding-bottom:10px;border-bottom:2px solid #eee}
.metrics-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(200px,1fr));gap:15px}
.metric-item{background:#f8f9fa;padding:15px;border-radius:6px;border-left:4px solid #667eea}
.metric-item .key{font-size:12px;color:#666;margin-bottom:5px}
.metric-item .value{font-size:18px;font-weight:bold}
.issue-item{background:white;border:1px solid #eee;border-left:4px solid;padding:20px;margin-bottom:15px;border-radius:4px}
.issue-item.critical{border-left-color:#dc3545;background:#fff5f5}
.issue-item.warning{border-left-color:#ffc107;background:#fffbf0}
.issue-item .severity{display:inline-block;padding:4px 12px;border-radius:12px;font-size:12px;font-weight:bold;text-transform:uppercase;margin-bottom:10px}
.severity.critical{background:#dc3545;color:white}
.severity.warning{background:#ffc107;color:#333}
.suggestion-item{background:#e7f3ff;border-left:4px solid #0066cc;padding:15px 20px;margin-bottom:10px;border-radius:4px}
.no-issues{text-align:center;padding:40px;color:#28a745;font-size:18px}
.footer{background:#f8f9fa;padding:20px;text-align:center;color:#666;font-size:14px;border-top:1px solid #eee}
</style>
</head>
<body>
<div class="container">
<div class="header">
<h1>%s</h1>
<div class="timestamp">生成时间: %s</div>
</div>
<div class="summary">
<div class="summary-card">
<h3>健康状态</h3>
<div class="value status-%s">%s</div>
</div>
<div class="summary-card">
<h3>健康评分</h3>
<div class="value">%.1f</div>
</div>
<div class="summary-card">
<h3>检测到的问题</h3>
<div class="value">%d</div>
</div>
</div>
<div class="content">
`, generateTitle, generateTitle, time.Now().Format("2006-01-02 15:04:05"), result.Status, statusText, result.Score, len(result.Issues)))

	// 指标
	if len(result.Metrics) > 0 {
		html.WriteString("<div class=\"section\"><h2>📊 系统指标</h2><div class=\"metrics-grid\">")
		for key, value := range result.Metrics {
			html.WriteString(fmt.Sprintf("<div class=\"metric-item\"><div class=\"key\">%s</div><div class=\"value\">%v</div></div>", key, value))
		}
		html.WriteString("</div></div>")
	}

	// 问题
	html.WriteString("<div class=\"section\"><h2>⚠️ 问题列表</h2>")
	if len(result.Issues) > 0 {
		for _, issue := range result.Issues {
			html.WriteString(fmt.Sprintf("<div class=\"issue-item %s\"><span class=\"severity %s\">%s</span><div class=\"description\">%s</div><div class=\"details\">当前值: %s | 阈值: %s</div></div>",
				issue.Severity, issue.Severity, issue.Severity, issue.Description, issue.Value, issue.Threshold))
		}
	} else {
		html.WriteString("<div class=\"no-issues\">✅ 未检测到问题</div>")
	}
	html.WriteString("</div>")

	// 建议
	if len(result.Suggestions) > 0 {
		html.WriteString("<div class=\"section\"><h2>💡 优化建议</h2>")
		for _, suggestion := range result.Suggestions {
			html.WriteString(fmt.Sprintf("<div class=\"suggestion-item\">%s</div>", suggestion))
		}
		html.WriteString("</div>")
	}

	html.WriteString("</div><div class=\"footer\">由 ClusterReport 生成</div></div></body></html>")
	return html.String()
}

func buildMarkdownReport(result *analyzer.AnalysisResult) string {
	var md strings.Builder

	md.WriteString(fmt.Sprintf("# %s\n\n", generateTitle))
	md.WriteString(fmt.Sprintf("**生成时间**: %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	md.WriteString("## 📊 执行摘要\n\n")
	md.WriteString(fmt.Sprintf("- **健康状态**: %s\n", getStatusTextChinese(result.Status)))
	md.WriteString(fmt.Sprintf("- **健康评分**: %.1f/100\n", result.Score))
	md.WriteString(fmt.Sprintf("- **检测到的问题**: %d 个\n\n", len(result.Issues)))

	if len(result.Metrics) > 0 {
		md.WriteString("## 📈 系统指标\n\n")
		for key, value := range result.Metrics {
			md.WriteString(fmt.Sprintf("- **%s**: %v\n", key, value))
		}
		md.WriteString("\n")
	}

	md.WriteString("## ⚠️  问题列表\n\n")
	if len(result.Issues) > 0 {
		for i, issue := range result.Issues {
			md.WriteString(fmt.Sprintf("### %d. [%s] %s\n\n", i+1, strings.ToUpper(issue.Severity), issue.Description))
			md.WriteString(fmt.Sprintf("- **类别**: %s\n", issue.Category))
			md.WriteString(fmt.Sprintf("- **当前值**: %s\n", issue.Value))
			md.WriteString(fmt.Sprintf("- **阈值**: %s\n\n", issue.Threshold))
		}
	} else {
		md.WriteString("✅ 未检测到问题\n\n")
	}

	if len(result.Suggestions) > 0 {
		md.WriteString("## 💡 优化建议\n\n")
		for i, suggestion := range result.Suggestions {
			md.WriteString(fmt.Sprintf("%d. %s\n", i+1, suggestion))
		}
		md.WriteString("\n")
	}

	md.WriteString("---\n\n*由 ClusterReport 生成*\n")
	return md.String()
}

func getStatusTextChinese(status string) string {
	switch status {
	case "healthy":
		return "健康"
	case "warning":
		return "警告"
	case "critical":
		return "严重"
	default:
		return status
	}
}
