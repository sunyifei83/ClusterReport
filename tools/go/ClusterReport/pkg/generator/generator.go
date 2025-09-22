package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/analyzer"
	"github.com/jung-kurt/gofpdf"
	"github.com/tealeg/xlsx/v3"
)

// OutputFormat 输出格式
type OutputFormat string

const (
	OutputFormatHTML     OutputFormat = "html"
	OutputFormatPDF      OutputFormat = "pdf"
	OutputFormatMarkdown OutputFormat = "markdown"
	OutputFormatExcel    OutputFormat = "excel"
	OutputFormatJSON     OutputFormat = "json"
)

// Generator 报告生成器接口
type Generator interface {
	// Generate 生成报告
	Generate(ctx context.Context, report *analyzer.Report) ([]byte, error)

	// Format 输出格式
	Format() OutputFormat

	// SetTemplate 设置模板
	SetTemplate(tmpl Template) error
}

// Template 报告模板
type Template struct {
	Name     string
	Path     string
	Content  string
	Sections []Section
	Assets   []Asset
}

// Section 报告章节
type Section struct {
	Title    string
	Type     string // text, table, chart
	DataPath string
	Template string
	Order    int
}

// Asset 资源文件
type Asset struct {
	Name string
	Path string
	Type string // css, js, image
}

// Options 生成器选项
type Options struct {
	Title          string
	Author         string
	Date           time.Time
	IncludeCharts  bool
	IncludeRawData bool
	IncludeTOC     bool
	CustomCSS      string
	CustomJS       string
}

// MultiFormat 多格式生成器
type MultiFormat struct {
	generators map[OutputFormat]Generator
	options    Options
}

// NewMultiFormat 创建多格式生成器
func NewMultiFormat(generators ...Generator) *MultiFormat {
	m := &MultiFormat{
		generators: make(map[OutputFormat]Generator),
		options: Options{
			Title:         "Cluster Report",
			Author:        "ClusterReport",
			Date:          time.Now(),
			IncludeCharts: true,
			IncludeTOC:    true,
		},
	}

	for _, g := range generators {
		m.generators[g.Format()] = g
	}

	return m
}

// GenerateAll 生成所有格式的报告
func (m *MultiFormat) GenerateAll(ctx context.Context, report *analyzer.Report) (map[OutputFormat][]byte, error) {
	results := make(map[OutputFormat][]byte)

	for format, generator := range m.generators {
		data, err := generator.Generate(ctx, report)
		if err != nil {
			return nil, fmt.Errorf("failed to generate %s report: %w", format, err)
		}
		results[format] = data
	}

	return results, nil
}

// HTMLGenerator HTML报告生成器
type HTMLGenerator struct {
	template *Template
	options  Options
}

// NewHTMLGenerator 创建HTML生成器
func NewHTMLGenerator() *HTMLGenerator {
	return &HTMLGenerator{
		options: Options{
			Title:         "Cluster Report",
			IncludeCharts: true,
			IncludeTOC:    true,
		},
	}
}

// Format 返回格式
func (g *HTMLGenerator) Format() OutputFormat {
	return OutputFormatHTML
}

// SetTemplate 设置模板
func (g *HTMLGenerator) SetTemplate(tmpl Template) error {
	g.template = &tmpl
	return nil
}

// Generate 生成HTML报告
func (g *HTMLGenerator) Generate(ctx context.Context, report *analyzer.Report) ([]byte, error) {
	// 使用默认模板
	htmlTemplate := g.getDefaultTemplate()

	tmpl, err := template.New("report").Funcs(template.FuncMap{
		"formatTime":    formatTime,
		"formatFloat":   formatFloat,
		"formatPercent": formatPercent,
		"json":          toJSON,
	}).Parse(htmlTemplate)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	data := map[string]interface{}{
		"Title":        g.options.Title,
		"GeneratedAt":  time.Now(),
		"Report":       report,
		"OverallScore": report.OverallScore,
		"Analyses":     report.Analyses,
		"Options":      g.options,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// getDefaultTemplate 获取默认HTML模板
func (g *HTMLGenerator) getDefaultTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
            background: #f5f5f5;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 2rem;
            border-radius: 10px;
            margin-bottom: 2rem;
        }
        .score-badge {
            display: inline-block;
            background: rgba(255,255,255,0.2);
            padding: 0.5rem 1rem;
            border-radius: 20px;
            font-size: 1.2rem;
            font-weight: bold;
        }
        .section {
            background: white;
            padding: 1.5rem;
            margin-bottom: 1.5rem;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .insight {
            padding: 1rem;
            margin: 0.5rem 0;
            border-left: 4px solid;
            background: #f9f9f9;
        }
        .insight.info { border-color: #3498db; }
        .insight.warning { border-color: #f39c12; }
        .insight.critical { border-color: #e74c3c; }
        .recommendation {
            background: #e8f5e9;
            padding: 1rem;
            margin: 0.5rem 0;
            border-radius: 4px;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin: 1rem 0;
        }
        th, td {
            padding: 0.75rem;
            text-align: left;
            border-bottom: 1px solid #ddd;
        }
        th {
            background: #f5f5f5;
            font-weight: 600;
        }
        .chart-container {
            margin: 2rem 0;
            padding: 1rem;
            background: white;
            border-radius: 8px;
        }
        .footer {
            text-align: center;
            color: #666;
            margin-top: 3rem;
            padding-top: 2rem;
            border-top: 1px solid #ddd;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>{{.Title}}</h1>
        <p>Generated at: {{formatTime .GeneratedAt}}</p>
        <div class="score-badge">Overall Score: {{formatFloat .OverallScore}}%</div>
    </div>

    {{if .Options.IncludeTOC}}
    <div class="section">
        <h2>Table of Contents</h2>
        <ul>
            <li><a href="#summary">Executive Summary</a></li>
            <li><a href="#analyses">Analysis Results</a></li>
            <li><a href="#insights">Insights & Recommendations</a></li>
            <li><a href="#details">Detailed Findings</a></li>
        </ul>
    </div>
    {{end}}

    <div id="summary" class="section">
        <h2>Executive Summary</h2>
        <p>This report provides a comprehensive analysis of your cluster infrastructure.</p>
        <ul>
            <li>Total Nodes Analyzed: {{len .Report.Data}}</li>
            <li>Overall Health Score: {{formatFloat .OverallScore}}%</li>
            <li>Report Generated: {{formatTime .GeneratedAt}}</li>
        </ul>
    </div>

    <div id="analyses" class="section">
        <h2>Analysis Results</h2>
        {{range .Analyses}}
        <div class="analysis">
            <h3>{{.Type}} Analysis</h3>
            <p>Score: {{formatFloat .Score}}%</p>
            
            {{if .Insights}}
            <h4>Key Insights:</h4>
            {{range .Insights}}
            <div class="insight {{.Level}}">
                <strong>{{.Description}}</strong>
                <br>Category: {{.Category}} | Level: {{.Level}}
            </div>
            {{end}}
            {{end}}

            {{if .Recommendations}}
            <h4>Recommendations:</h4>
            {{range .Recommendations}}
            <div class="recommendation">
                <strong>{{.Action}}</strong>
                <br>Reason: {{.Reason}}
                <br>Impact: {{.Impact}} | Effort: {{.Effort}} | Priority: {{.Priority}}
            </div>
            {{end}}
            {{end}}
        </div>
        {{end}}
    </div>

    <div class="footer">
        <p>Generated by ClusterReport v0.1.0</p>
    </div>
</body>
</html>`
}

// PDFGenerator PDF报告生成器
type PDFGenerator struct {
	template *Template
	options  Options
}

// NewPDFGenerator 创建PDF生成器
func NewPDFGenerator() *PDFGenerator {
	return &PDFGenerator{
		options: Options{
			Title: "Cluster Report",
		},
	}
}

// Format 返回格式
func (g *PDFGenerator) Format() OutputFormat {
	return OutputFormatPDF
}

// SetTemplate 设置模板
func (g *PDFGenerator) SetTemplate(tmpl Template) error {
	g.template = &tmpl
	return nil
}

// Generate 生成PDF报告
func (g *PDFGenerator) Generate(ctx context.Context, report *analyzer.Report) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// 添加页面
	pdf.AddPage()

	// 设置字体
	pdf.SetFont("Arial", "B", 16)

	// 标题
	pdf.Cell(190, 10, g.options.Title)
	pdf.Ln(12)

	// 生成时间
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(190, 10, fmt.Sprintf("Generated: %s", time.Now().Format("2006-01-02 15:04:05")))
	pdf.Ln(10)

	// 总体评分
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(190, 10, fmt.Sprintf("Overall Score: %.1f%%", report.OverallScore))
	pdf.Ln(15)

	// 分析结果
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(190, 10, "Analysis Results")
	pdf.Ln(10)

	for _, analysis := range report.Analyses {
		// 分析类型
		pdf.SetFont("Arial", "B", 11)
		pdf.Cell(190, 8, string(analysis.Type))
		pdf.Ln(8)

		// 评分
		pdf.SetFont("Arial", "", 10)
		pdf.Cell(190, 6, fmt.Sprintf("Score: %.1f%%", analysis.Score))
		pdf.Ln(6)

		// 洞察
		if len(analysis.Insights) > 0 {
			pdf.SetFont("Arial", "I", 10)
			pdf.Cell(190, 6, "Key Insights:")
			pdf.Ln(6)

			for _, insight := range analysis.Insights {
				pdf.SetFont("Arial", "", 9)
				pdf.MultiCell(180, 5, fmt.Sprintf("• %s (%s)", insight.Description, insight.Level), "", "", false)
				pdf.Ln(2)
			}
		}

		pdf.Ln(5)
	}

	// 输出到缓冲区
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// MarkdownGenerator Markdown报告生成器
type MarkdownGenerator struct {
	template *Template
	options  Options
}

// NewMarkdownGenerator 创建Markdown生成器
func NewMarkdownGenerator() *MarkdownGenerator {
	return &MarkdownGenerator{
		options: Options{
			Title:      "Cluster Report",
			IncludeTOC: true,
		},
	}
}

// Format 返回格式
func (g *MarkdownGenerator) Format() OutputFormat {
	return OutputFormatMarkdown
}

// SetTemplate 设置模板
func (g *MarkdownGenerator) SetTemplate(tmpl Template) error {
	g.template = &tmpl
	return nil
}

// Generate 生成Markdown报告
func (g *MarkdownGenerator) Generate(ctx context.Context, report *analyzer.Report) ([]byte, error) {
	var sb strings.Builder

	// 标题
	sb.WriteString(fmt.Sprintf("# %s\n\n", g.options.Title))
	sb.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("**Overall Score:** %.1f%%\n\n", report.OverallScore))

	// 目录
	if g.options.IncludeTOC {
		sb.WriteString("## Table of Contents\n\n")
		sb.WriteString("- [Executive Summary](#executive-summary)\n")
		sb.WriteString("- [Analysis Results](#analysis-results)\n")
		sb.WriteString("- [Insights and Recommendations](#insights-and-recommendations)\n\n")
	}

	// 执行摘要
	sb.WriteString("## Executive Summary\n\n")
	sb.WriteString("This report provides a comprehensive analysis of your cluster infrastructure.\n\n")
	sb.WriteString(fmt.Sprintf("- **Total Nodes Analyzed:** %d\n", len(report.Data)))
	sb.WriteString(fmt.Sprintf("- **Overall Health Score:** %.1f%%\n", report.OverallScore))
	sb.WriteString(fmt.Sprintf("- **Report Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	// 分析结果
	sb.WriteString("## Analysis Results\n\n")

	for _, analysis := range report.Analyses {
		sb.WriteString(fmt.Sprintf("### %s Analysis\n\n", analysis.Type))
		sb.WriteString(fmt.Sprintf("**Score:** %.1f%%\n\n", analysis.Score))

		// 洞察
		if len(analysis.Insights) > 0 {
			sb.WriteString("#### Key Insights:\n\n")
			for _, insight := range analysis.Insights {
				sb.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", insight.Level, insight.Category, insight.Description))
			}
			sb.WriteString("\n")
		}

		// 建议
		if len(analysis.Recommendations) > 0 {
			sb.WriteString("#### Recommendations:\n\n")
			for _, rec := range analysis.Recommendations {
				sb.WriteString(fmt.Sprintf("- **%s** (Priority %d)\n", rec.Action, rec.Priority))
				sb.WriteString(fmt.Sprintf("  - Reason: %s\n", rec.Reason))
				sb.WriteString(fmt.Sprintf("  - Impact: %s\n", rec.Impact))
				sb.WriteString(fmt.Sprintf("  - Effort: %s\n", rec.Effort))
			}
			sb.WriteString("\n")
		}
	}

	return []byte(sb.String()), nil
}

// ExcelGenerator Excel报告生成器
type ExcelGenerator struct {
	template *Template
	options  Options
}

// NewExcelGenerator 创建Excel生成器
func NewExcelGenerator() *ExcelGenerator {
	return &ExcelGenerator{
		options: Options{
			Title: "Cluster Report",
		},
	}
}

// Format 返回格式
func (g *ExcelGenerator) Format() OutputFormat {
	return OutputFormatExcel
}

// SetTemplate 设置模板
func (g *ExcelGenerator) SetTemplate(tmpl Template) error {
	g.template = &tmpl
	return nil
}

// Generate 生成Excel报告
func (g *ExcelGenerator) Generate(ctx context.Context, report *analyzer.Report) ([]byte, error) {
	file := xlsx.NewFile()

	// 创建摘要sheet
	summarySheet, err := file.AddSheet("Summary")
	if err != nil {
		return nil, err
	}

	// 添加标题
	titleRow := summarySheet.AddRow()
	titleCell := titleRow.AddCell()
	titleCell.Value = g.options.Title

	// 添加生成时间
	timeRow := summarySheet.AddRow()
	timeCell := timeRow.AddCell()
	timeCell.Value = "Generated: " + time.Now().Format("2006-01-02 15:04:05")

	// 添加总体评分
	scoreRow := summarySheet.AddRow()
	scoreCell := scoreRow.AddCell()
	scoreCell.Value = fmt.Sprintf("Overall Score: %.1f%%", report.OverallScore)

	// 空行
	summarySheet.AddRow()

	// 添加分析结果
	for _, analysis := range report.Analyses {
		// 分析类型
		typeRow := summarySheet.AddRow()
		typeCell := typeRow.AddCell()
		typeCell.Value = string(analysis.Type) + " Analysis"

		// 评分
		scoreRow := summarySheet.AddRow()
		scoreCell := scoreRow.AddCell()
		scoreCell.Value = fmt.Sprintf("Score: %.1f%%", analysis.Score)

		// 洞察数量
		if len(analysis.Insights) > 0 {
			insightRow := summarySheet.AddRow()
			insightCell := insightRow.AddCell()
			insightCell.Value = fmt.Sprintf("Insights: %d", len(analysis.Insights))
		}

		// 建议数量
		if len(analysis.Recommendations) > 0 {
			recRow := summarySheet.AddRow()
			recCell := recRow.AddCell()
			recCell.Value = fmt.Sprintf("Recommendations: %d", len(analysis.Recommendations))
		}

		// 空行
		summarySheet.AddRow()
	}

	// 创建详细数据sheet
	dataSheet, err := file.AddSheet("Raw Data")
	if err != nil {
		return nil, err
	}

	// 添加表头
	headerRow := dataSheet.AddRow()
	headerRow.AddCell().Value = "Node"
	headerRow.AddCell().Value = "Type"
	headerRow.AddCell().Value = "Timestamp"
	headerRow.AddCell().Value = "Metrics"

	// 添加数据
	for _, data := range report.Data {
		dataRow := dataSheet.AddRow()
		dataRow.AddCell().Value = data.Node
		dataRow.AddCell().Value = string(data.Type)
		dataRow.AddCell().Value = data.Timestamp.Format("2006-01-02 15:04:05")

		// 将metrics转换为JSON字符串
		metricsJSON, _ := json.Marshal(data.Metrics)
		dataRow.AddCell().Value = string(metricsJSON)
	}

	// 保存到缓冲区
	var buf bytes.Buffer
	err = file.Write(&buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// 辅助函数

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func formatFloat(f float64) string {
	return fmt.Sprintf("%.1f", f)
}

func formatPercent(f float64) string {
	return fmt.Sprintf("%.1f%%", f)
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

// SaveToFile 保存报告到文件
func SaveToFile(data []byte, filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}
