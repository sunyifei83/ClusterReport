package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/analyzer"
	"github.com/devops-toolkit/clusterreport/pkg/collector"
)

// Generator 报告生成器接口
type Generator interface {
	Generate(data interface{}) ([]byte, error)
	Format() string
}

// ReportData 报告数据
type ReportData struct {
	Title       string                   `json:"title"`
	GeneratedAt time.Time                `json:"generated_at"`
	SystemInfo  *SystemInfo              `json:"system_info"`
	Metrics     *collector.SystemMetrics `json:"metrics,omitempty"`
	Analysis    *analyzer.AnalysisResult `json:"analysis,omitempty"`
	Summary     *Summary                 `json:"summary"`
	Custom      map[string]interface{}   `json:"custom,omitempty"`
}

// SystemInfo 系统信息
type SystemInfo struct {
	Hostname      string `json:"hostname"`
	OS            string `json:"os"`
	Architecture  string `json:"architecture"`
	KernelVersion string `json:"kernel_version"`
}

// Summary 报告摘要
type Summary struct {
	Status      string   `json:"status"`
	Score       float64  `json:"score"`
	TotalIssues int      `json:"total_issues"`
	Critical    int      `json:"critical"`
	Warning     int      `json:"warning"`
	Healthy     int      `json:"healthy"`
	Highlights  []string `json:"highlights"`
}

// JSONGenerator JSON 格式生成器
type JSONGenerator struct{}

// NewJSONGenerator 创建 JSON 生成器
func NewJSONGenerator() *JSONGenerator {
	return &JSONGenerator{}
}

// Generate 生成 JSON 报告
func (g *JSONGenerator) Generate(data interface{}) ([]byte, error) {
	return json.MarshalIndent(data, "", "  ")
}

// Format 返回格式名称
func (g *JSONGenerator) Format() string {
	return "json"
}

// HTMLGenerator HTML 格式生成器
type HTMLGenerator struct {
	template *template.Template
}

// NewHTMLGenerator 创建 HTML 生成器
func NewHTMLGenerator() (*HTMLGenerator, error) {
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &HTMLGenerator{
		template: tmpl,
	}, nil
}

// Generate 生成 HTML 报告
func (g *HTMLGenerator) Generate(data interface{}) ([]byte, error) {
	var buf bytes.Buffer

	if err := g.template.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.Bytes(), nil
}

// Format 返回格式名称
func (g *HTMLGenerator) Format() string {
	return "html"
}

// MarkdownGenerator Markdown 格式生成器
type MarkdownGenerator struct{}

// NewMarkdownGenerator 创建 Markdown 生成器
func NewMarkdownGenerator() *MarkdownGenerator {
	return &MarkdownGenerator{}
}

// Generate 生成 Markdown 报告
func (g *MarkdownGenerator) Generate(data interface{}) ([]byte, error) {
	reportData, ok := data.(*ReportData)
	if !ok {
		return nil, fmt.Errorf("invalid data type")
	}

	var buf bytes.Buffer

	// 标题
	buf.WriteString(fmt.Sprintf("# %s\n\n", reportData.Title))
	buf.WriteString(fmt.Sprintf("**生成时间**: %s\n\n", reportData.GeneratedAt.Format("2006-01-02 15:04:05")))

	// 系统信息
	if reportData.SystemInfo != nil {
		buf.WriteString("## 系统信息\n\n")
		buf.WriteString(fmt.Sprintf("- **主机名**: %s\n", reportData.SystemInfo.Hostname))
		buf.WriteString(fmt.Sprintf("- **操作系统**: %s\n", reportData.SystemInfo.OS))
		buf.WriteString(fmt.Sprintf("- **架构**: %s\n", reportData.SystemInfo.Architecture))
		buf.WriteString(fmt.Sprintf("- **内核版本**: %s\n\n", reportData.SystemInfo.KernelVersion))
	}

	// 摘要
	if reportData.Summary != nil {
		buf.WriteString("## 报告摘要\n\n")
		buf.WriteString(fmt.Sprintf("- **状态**: %s\n", getStatusEmoji(reportData.Summary.Status)))
		buf.WriteString(fmt.Sprintf("- **健康评分**: %.1f/100\n", reportData.Summary.Score))
		buf.WriteString(fmt.Sprintf("- **问题总数**: %d\n", reportData.Summary.TotalIssues))
		buf.WriteString(fmt.Sprintf("  - 严重: %d\n", reportData.Summary.Critical))
		buf.WriteString(fmt.Sprintf("  - 警告: %d\n", reportData.Summary.Warning))
		buf.WriteString("\n")

		if len(reportData.Summary.Highlights) > 0 {
			buf.WriteString("### 关键发现\n\n")
			for _, highlight := range reportData.Summary.Highlights {
				buf.WriteString(fmt.Sprintf("- %s\n", highlight))
			}
			buf.WriteString("\n")
		}
	}

	// 分析结果
	if reportData.Analysis != nil {
		buf.WriteString("## 详细分析\n\n")

		if len(reportData.Analysis.Issues) > 0 {
			buf.WriteString("### 发现的问题\n\n")
			buf.WriteString("| 严重度 | 类别 | 描述 | 当前值 | 阈值 |\n")
			buf.WriteString("|--------|------|------|--------|------|\n")

			for _, issue := range reportData.Analysis.Issues {
				buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
					getSeverityEmoji(issue.Severity),
					issue.Category,
					issue.Description,
					issue.Value,
					issue.Threshold,
				))
			}
			buf.WriteString("\n")
		}

		if len(reportData.Analysis.Suggestions) > 0 {
			buf.WriteString("### 优化建议\n\n")
			for i, suggestion := range reportData.Analysis.Suggestions {
				buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, suggestion))
			}
			buf.WriteString("\n")
		}
	}

	// 指标详情
	if reportData.Metrics != nil {
		buf.WriteString("## 系统指标\n\n")

		// CPU
		buf.WriteString("### CPU\n\n")
		buf.WriteString(fmt.Sprintf("- 核心数: %d\n", reportData.Metrics.CPU.Cores))
		buf.WriteString(fmt.Sprintf("- 使用率: %.2f%%\n", reportData.Metrics.CPU.Usage))
		buf.WriteString(fmt.Sprintf("- 负载平均值: %.2f (1分钟), %.2f (5分钟), %.2f (15分钟)\n\n",
			reportData.Metrics.CPU.LoadAvg1,
			reportData.Metrics.CPU.LoadAvg5,
			reportData.Metrics.CPU.LoadAvg15,
		))

		// 内存
		buf.WriteString("### 内存\n\n")
		buf.WriteString(fmt.Sprintf("- 总内存: %.2f GB\n", float64(reportData.Metrics.Memory.Total)/1024/1024/1024))
		buf.WriteString(fmt.Sprintf("- 已用内存: %.2f GB\n", float64(reportData.Metrics.Memory.Used)/1024/1024/1024))
		buf.WriteString(fmt.Sprintf("- 可用内存: %.2f GB\n", float64(reportData.Metrics.Memory.Available)/1024/1024/1024))
		buf.WriteString(fmt.Sprintf("- 使用率: %.2f%%\n\n", reportData.Metrics.Memory.UsedPercent))

		// 磁盘
		if len(reportData.Metrics.Disk) > 0 {
			buf.WriteString("### 磁盘\n\n")
			buf.WriteString("| 挂载点 | 总容量 | 已用 | 可用 | 使用率 |\n")
			buf.WriteString("|--------|--------|------|------|--------|\n")

			for _, disk := range reportData.Metrics.Disk {
				buf.WriteString(fmt.Sprintf("| %s | %.2f GB | %.2f GB | %.2f GB | %.1f%% |\n",
					disk.MountPoint,
					float64(disk.Total)/1024/1024/1024,
					float64(disk.Used)/1024/1024/1024,
					float64(disk.Available)/1024/1024/1024,
					disk.UsedPercent,
				))
			}
			buf.WriteString("\n")
		}
	}

	// 页脚
	buf.WriteString("---\n\n")
	buf.WriteString("*此报告由 ClusterReport 自动生成*\n")

	return buf.Bytes(), nil
}

// Format 返回格式名称
func (g *MarkdownGenerator) Format() string {
	return "markdown"
}

// Helper functions

func getStatusEmoji(status string) string {
	switch status {
	case "healthy":
		return "✅ 健康"
	case "warning":
		return "⚠️ 警告"
	case "critical":
		return "🔴 严重"
	default:
		return "❓ 未知"
	}
}

func getSeverityEmoji(severity string) string {
	switch severity {
	case "critical":
		return "🔴"
	case "warning":
		return "⚠️"
	case "low":
		return "ℹ️"
	default:
		return "•"
	}
}

// HTML 模板
const htmlTemplate = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            background-color: #f5f7fa;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px 20px;
            border-radius: 10px;
            margin-bottom: 30px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
        }
        
        .header h1 {
            font-size: 32px;
            margin-bottom: 10px;
        }
        
        .header .subtitle {
            opacity: 0.9;
            font-size: 14px;
        }
        
        .card {
            background: white;
            border-radius: 10px;
            padding: 25px;
            margin-bottom: 20px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.08);
        }
        
        .card h2 {
            color: #667eea;
            margin-bottom: 20px;
            padding-bottom: 10px;
            border-bottom: 2px solid #f0f0f0;
            font-size: 24px;
        }
        
        .summary-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-bottom: 20px;
        }
        
        .summary-item {
            background: #f8f9fa;
            padding: 15px;
            border-radius: 8px;
            border-left: 4px solid #667eea;
        }
        
        .summary-item .label {
            font-size: 12px;
            color: #666;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        .summary-item .value {
            font-size: 24px;
            font-weight: bold;
            color: #333;
            margin-top: 5px;
        }
        
        .status-badge {
            display: inline-block;
            padding: 6px 12px;
            border-radius: 20px;
            font-size: 14px;
            font-weight: 600;
        }
        
        .status-healthy {
            background: #d4edda;
            color: #155724;
        }
        
        .status-warning {
            background: #fff3cd;
            color: #856404;
        }
        
        .status-critical {
            background: #f8d7da;
            color: #721c24;
        }
        
        .progress-bar {
            background: #e0e0e0;
            border-radius: 10px;
            height: 20px;
            overflow: hidden;
            margin: 10px 0;
        }
        
        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #667eea 0%, #764ba2 100%);
            transition: width 0.3s ease;
        }
        
        table {
            width: 100%;
            border-collapse: collapse;
            margin-top: 15px;
        }
        
        th, td {
            padding: 12px;
            text-align: left;
            border-bottom: 1px solid #e0e0e0;
        }
        
        th {
            background: #f8f9fa;
            font-weight: 600;
            color: #666;
            font-size: 12px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }
        
        tr:hover {
            background: #f8f9fa;
        }
        
        .metric-value {
            font-weight: 600;
            color: #667eea;
        }
        
        .issue-critical {
            color: #dc3545;
        }
        
        .issue-warning {
            color: #ffc107;
        }
        
        .footer {
            text-align: center;
            color: #999;
            font-size: 12px;
            margin-top: 40px;
            padding-top: 20px;
            border-top: 1px solid #e0e0e0;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.Title}}</h1>
            <div class="subtitle">生成时间: {{.GeneratedAt.Format "2006-01-02 15:04:05"}}</div>
        </div>
        
        {{if .Summary}}
        <div class="card">
            <h2>报告摘要</h2>
            <div class="summary-grid">
                <div class="summary-item">
                    <div class="label">状态</div>
                    <div class="value">
                        <span class="status-badge status-{{.Summary.Status}}">
                            {{if eq .Summary.Status "healthy"}}✅ 健康{{end}}
                            {{if eq .Summary.Status "warning"}}⚠️ 警告{{end}}
                            {{if eq .Summary.Status "critical"}}🔴 严重{{end}}
                        </span>
                    </div>
                </div>
                <div class="summary-item">
                    <div class="label">健康评分</div>
                    <div class="value">{{printf "%.1f" .Summary.Score}}/100</div>
                    <div class="progress-bar">
                        <div class="progress-fill" style="width: {{.Summary.Score}}%"></div>
                    </div>
                </div>
                <div class="summary-item">
                    <div class="label">严重问题</div>
                    <div class="value issue-critical">{{.Summary.Critical}}</div>
                </div>
                <div class="summary-item">
                    <div class="label">警告</div>
                    <div class="value issue-warning">{{.Summary.Warning}}</div>
                </div>
            </div>
        </div>
        {{end}}
        
        {{if .Analysis}}
        {{if .Analysis.Issues}}
        <div class="card">
            <h2>发现的问题</h2>
            <table>
                <thead>
                    <tr>
                        <th>严重度</th>
                        <th>类别</th>
                        <th>描述</th>
                        <th>当前值</th>
                        <th>阈值</th>
                    </tr>
                </thead>
                <tbody>
                    {{range .Analysis.Issues}}
                    <tr>
                        <td>
                            {{if eq .Severity "critical"}}🔴{{end}}
                            {{if eq .Severity "warning"}}⚠️{{end}}
                            {{.Severity}}
                        </td>
                        <td>{{.Category}}</td>
                        <td>{{.Description}}</td>
                        <td class="metric-value">{{.Value}}</td>
                        <td>{{.Threshold}}</td>
                    </tr>
                    {{end}}
                </tbody>
            </table>
        </div>
        {{end}}
        
        {{if .Analysis.Suggestions}}
        <div class="card">
            <h2>优化建议</h2>
            <ul>
                {{range .Analysis.Suggestions}}
                <li>{{.}}</li>
                {{end}}
            </ul>
        </div>
        {{end}}
        {{end}}
        
        {{if .Metrics}}
        <div class="card">
            <h2>系统指标</h2>
            <h3>CPU</h3>
            <p>核心数: <span class="metric-value">{{.Metrics.CPU.Cores}}</span></p>
            <p>使用率: <span class="metric-value">{{printf "%.2f" .Metrics.CPU.Usage}}%</span></p>
            
            <h3 style="margin-top: 20px;">内存</h3>
            <p>使用率: <span class="metric-value">{{printf "%.2f" .Metrics.Memory.UsedPercent}}%</span></p>
            <p>总内存: {{printf "%.2f" (divf .Metrics.Memory.Total 1073741824)}} GB</p>
            <p>已用: {{printf "%.2f" (divf .Metrics.Memory.Used 1073741824)}} GB</p>
        </div>
        {{end}}
        
        <div class="footer">
            此报告由 ClusterReport 自动生成
        </div>
    </div>
</body>
</html>`
