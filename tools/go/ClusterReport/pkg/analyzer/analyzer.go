package analyzer

import (
	"context"
	"fmt"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/devops-toolkit/clusterreport/pkg/collector"
)

// AnalysisType 分析类型
type AnalysisType string

const (
	AnalysisTypeConfig      AnalysisType = "config"
	AnalysisTypePerformance AnalysisType = "performance"
	AnalysisTypeAnomaly     AnalysisType = "anomaly"
	AnalysisTypeTrend       AnalysisType = "trend"
	AnalysisTypeComparison  AnalysisType = "comparison"
)

// Analysis 分析结果
type Analysis struct {
	Type            AnalysisType           `json:"type"`
	Timestamp       time.Time              `json:"timestamp"`
	Results         map[string]interface{} `json:"results"`
	Insights        []Insight              `json:"insights"`
	Recommendations []Recommendation       `json:"recommendations"`
	Score           float64                `json:"score"`
}

// Insight 洞察
type Insight struct {
	Level       string                 `json:"level"` // info, warning, critical
	Category    string                 `json:"category"`
	Description string                 `json:"description"`
	Details     map[string]interface{} `json:"details"`
	Timestamp   time.Time              `json:"timestamp"`
}

// Recommendation 建议
type Recommendation struct {
	Priority  int       `json:"priority"` // 1-5, 1最高
	Category  string    `json:"category"`
	Action    string    `json:"action"`
	Reason    string    `json:"reason"`
	Impact    string    `json:"impact"`
	Effort    string    `json:"effort"` // low, medium, high
	Timestamp time.Time `json:"timestamp"`
}

// Analyzer 分析器接口
type Analyzer interface {
	// Analyze 分析数据
	Analyze(ctx context.Context, data []collector.Data) (*Analysis, error)

	// Type 分析类型
	Type() AnalysisType

	// Options 配置选项
	Options() map[string]interface{}
}

// Registry 分析器注册表
type Registry struct {
	analyzers map[AnalysisType]Analyzer
	mu        sync.RWMutex
}

// NewRegistry 创建新的注册表
func NewRegistry() *Registry {
	return &Registry{
		analyzers: make(map[AnalysisType]Analyzer),
	}
}

// Register 注册分析器
func (r *Registry) Register(analyzer Analyzer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	aType := analyzer.Type()
	if _, exists := r.analyzers[aType]; exists {
		return fmt.Errorf("analyzer %s already registered", aType)
	}

	r.analyzers[aType] = analyzer
	return nil
}

// Get 获取分析器
func (r *Registry) Get(aType AnalysisType) (Analyzer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	analyzer, exists := r.analyzers[aType]
	if !exists {
		return nil, fmt.Errorf("analyzer %s not found", aType)
	}

	return analyzer, nil
}

// Chain 分析链
type Chain struct {
	analyzers []Analyzer
	options   map[string]interface{}
}

// NewChain 创建分析链
func NewChain(analyzers ...Analyzer) *Chain {
	return &Chain{
		analyzers: analyzers,
		options:   make(map[string]interface{}),
	}
}

// Process 处理数据
func (c *Chain) Process(ctx context.Context, data []collector.Data) (*Report, error) {
	report := &Report{
		Timestamp: time.Now(),
		Analyses:  make([]*Analysis, 0),
		Data:      data,
	}

	// 执行所有分析器
	for _, analyzer := range c.analyzers {
		analysis, err := analyzer.Analyze(ctx, data)
		if err != nil {
			// 记录错误但继续
			fmt.Printf("Analyzer %s failed: %v\n", analyzer.Type(), err)
			continue
		}
		report.Analyses = append(report.Analyses, analysis)
	}

	// 计算总体评分
	report.calculateOverallScore()

	return report, nil
}

// Report 分析报告
type Report struct {
	Timestamp    time.Time              `json:"timestamp"`
	Analyses     []*Analysis            `json:"analyses"`
	Data         []collector.Data       `json:"data"`
	OverallScore float64                `json:"overall_score"`
	Summary      map[string]interface{} `json:"summary"`
}

// calculateOverallScore 计算总体评分
func (r *Report) calculateOverallScore() {
	if len(r.Analyses) == 0 {
		r.OverallScore = 0
		return
	}

	totalScore := 0.0
	for _, analysis := range r.Analyses {
		totalScore += analysis.Score
	}

	r.OverallScore = totalScore / float64(len(r.Analyses))
}

// ConfigAnalyzer 配置分析器
type ConfigAnalyzer struct {
	threshold map[string]interface{}
}

// NewConfigAnalyzer 创建配置分析器
func NewConfigAnalyzer() *ConfigAnalyzer {
	return &ConfigAnalyzer{
		threshold: map[string]interface{}{
			"min_memory": 8 * 1024 * 1024 * 1024, // 8GB
			"min_cpu":    4,
			"min_disk":   100 * 1024 * 1024 * 1024, // 100GB
		},
	}
}

// Type 返回分析类型
func (a *ConfigAnalyzer) Type() AnalysisType {
	return AnalysisTypeConfig
}

// Options 返回配置选项
func (a *ConfigAnalyzer) Options() map[string]interface{} {
	return a.threshold
}

// Analyze 分析配置
func (a *ConfigAnalyzer) Analyze(ctx context.Context, data []collector.Data) (*Analysis, error) {
	analysis := &Analysis{
		Type:            AnalysisTypeConfig,
		Timestamp:       time.Now(),
		Results:         make(map[string]interface{}),
		Insights:        make([]Insight, 0),
		Recommendations: make([]Recommendation, 0),
	}

	// 分析每个节点的配置
	nodeConfigs := make([]map[string]interface{}, 0)
	for _, d := range data {
		if d.Type == collector.DataTypeHardware {
			nodeConfigs = append(nodeConfigs, d.Metrics)

			// 检查配置是否满足最低要求
			a.checkNodeConfig(d.Node, d.Metrics, analysis)
		}
	}

	// 检查配置一致性
	a.checkConfigConsistency(nodeConfigs, analysis)

	// 计算评分
	analysis.Score = a.calculateScore(analysis)

	return analysis, nil
}

// checkNodeConfig 检查节点配置
func (a *ConfigAnalyzer) checkNodeConfig(node string, config map[string]interface{}, analysis *Analysis) {
	// 检查内存
	if memory, ok := config["memory"].(float64); ok {
		minMemory := a.threshold["min_memory"].(int)
		if memory < float64(minMemory) {
			analysis.Insights = append(analysis.Insights, Insight{
				Level:       "warning",
				Category:    "hardware",
				Description: fmt.Sprintf("Node %s has insufficient memory", node),
				Details: map[string]interface{}{
					"current":  memory,
					"required": minMemory,
				},
				Timestamp: time.Now(),
			})

			analysis.Recommendations = append(analysis.Recommendations, Recommendation{
				Priority:  2,
				Category:  "hardware",
				Action:    fmt.Sprintf("Upgrade memory on node %s", node),
				Reason:    "Current memory is below recommended minimum",
				Impact:    "Improved performance and stability",
				Effort:    "medium",
				Timestamp: time.Now(),
			})
		}
	}

	// 检查CPU
	if cpu, ok := config["cpu_cores"].(float64); ok {
		minCPU := a.threshold["min_cpu"].(int)
		if cpu < float64(minCPU) {
			analysis.Insights = append(analysis.Insights, Insight{
				Level:       "warning",
				Category:    "hardware",
				Description: fmt.Sprintf("Node %s has insufficient CPU cores", node),
				Details: map[string]interface{}{
					"current":  cpu,
					"required": minCPU,
				},
				Timestamp: time.Now(),
			})
		}
	}
}

// checkConfigConsistency 检查配置一致性
func (a *ConfigAnalyzer) checkConfigConsistency(configs []map[string]interface{}, analysis *Analysis) {
	if len(configs) < 2 {
		return
	}

	// 检查关键配置的一致性
	keys := []string{"os_version", "kernel_version", "cpu_model"}

	for _, key := range keys {
		values := make(map[string]int)
		for _, config := range configs {
			if val, ok := config[key].(string); ok {
				values[val]++
			}
		}

		if len(values) > 1 {
			analysis.Insights = append(analysis.Insights, Insight{
				Level:       "info",
				Category:    "consistency",
				Description: fmt.Sprintf("Inconsistent %s across nodes", key),
				Details: map[string]interface{}{
					"values": values,
				},
				Timestamp: time.Now(),
			})
		}
	}
}

// calculateScore 计算评分
func (a *ConfigAnalyzer) calculateScore(analysis *Analysis) float64 {
	score := 100.0

	// 根据洞察级别扣分
	for _, insight := range analysis.Insights {
		switch insight.Level {
		case "critical":
			score -= 20
		case "warning":
			score -= 10
		case "info":
			score -= 2
		}
	}

	// 确保分数在0-100之间
	if score < 0 {
		score = 0
	}

	return score
}

// PerfAnalyzer 性能分析器
type PerfAnalyzer struct {
	threshold map[string]float64
}

// NewPerfAnalyzer 创建性能分析器
func NewPerfAnalyzer() *PerfAnalyzer {
	return &PerfAnalyzer{
		threshold: map[string]float64{
			"cpu_usage":    80.0,
			"memory_usage": 85.0,
			"disk_usage":   90.0,
			"load_average": 2.0,
		},
	}
}

// Type 返回分析类型
func (a *PerfAnalyzer) Type() AnalysisType {
	return AnalysisTypePerformance
}

// Options 返回配置选项
func (a *PerfAnalyzer) Options() map[string]interface{} {
	options := make(map[string]interface{})
	for k, v := range a.threshold {
		options[k] = v
	}
	return options
}

// Analyze 分析性能
func (a *PerfAnalyzer) Analyze(ctx context.Context, data []collector.Data) (*Analysis, error) {
	analysis := &Analysis{
		Type:            AnalysisTypePerformance,
		Timestamp:       time.Now(),
		Results:         make(map[string]interface{}),
		Insights:        make([]Insight, 0),
		Recommendations: make([]Recommendation, 0),
	}

	// 收集性能指标
	perfMetrics := make(map[string][]float64)

	for _, d := range data {
		if d.Type == collector.DataTypePerformance {
			a.analyzeNodePerformance(d.Node, d.Metrics, analysis, perfMetrics)
		}
	}

	// 计算统计信息
	a.calculateStatistics(perfMetrics, analysis)

	// 计算评分
	analysis.Score = a.calculateScore(analysis)

	return analysis, nil
}

// analyzeNodePerformance 分析节点性能
func (a *PerfAnalyzer) analyzeNodePerformance(node string, metrics map[string]interface{},
	analysis *Analysis, perfMetrics map[string][]float64) {

	// 检查CPU使用率
	if cpuUsage, ok := metrics["cpu_usage"].(float64); ok {
		perfMetrics["cpu_usage"] = append(perfMetrics["cpu_usage"], cpuUsage)

		if cpuUsage > a.threshold["cpu_usage"] {
			analysis.Insights = append(analysis.Insights, Insight{
				Level:       "warning",
				Category:    "performance",
				Description: fmt.Sprintf("High CPU usage on node %s", node),
				Details: map[string]interface{}{
					"current":   cpuUsage,
					"threshold": a.threshold["cpu_usage"],
				},
				Timestamp: time.Now(),
			})
		}
	}

	// 检查内存使用率
	if memUsage, ok := metrics["memory_usage"].(float64); ok {
		perfMetrics["memory_usage"] = append(perfMetrics["memory_usage"], memUsage)

		if memUsage > a.threshold["memory_usage"] {
			analysis.Insights = append(analysis.Insights, Insight{
				Level:       "warning",
				Category:    "performance",
				Description: fmt.Sprintf("High memory usage on node %s", node),
				Details: map[string]interface{}{
					"current":   memUsage,
					"threshold": a.threshold["memory_usage"],
				},
				Timestamp: time.Now(),
			})

			analysis.Recommendations = append(analysis.Recommendations, Recommendation{
				Priority:  3,
				Category:  "performance",
				Action:    "Investigate memory-consuming processes",
				Reason:    fmt.Sprintf("Memory usage %.1f%% exceeds threshold", memUsage),
				Impact:    "Prevent out-of-memory issues",
				Effort:    "low",
				Timestamp: time.Now(),
			})
		}
	}
}

// calculateStatistics 计算统计信息
func (a *PerfAnalyzer) calculateStatistics(metrics map[string][]float64, analysis *Analysis) {
	stats := make(map[string]map[string]float64)

	for metric, values := range metrics {
		if len(values) == 0 {
			continue
		}

		stats[metric] = map[string]float64{
			"min":    min(values),
			"max":    max(values),
			"avg":    avg(values),
			"median": median(values),
			"stddev": stddev(values),
		}
	}

	analysis.Results["statistics"] = stats
}

// calculateScore 计算评分
func (a *PerfAnalyzer) calculateScore(analysis *Analysis) float64 {
	score := 100.0

	// 根据洞察级别扣分
	for _, insight := range analysis.Insights {
		switch insight.Level {
		case "critical":
			score -= 15
		case "warning":
			score -= 8
		case "info":
			score -= 2
		}
	}

	// 确保分数在0-100之间
	if score < 0 {
		score = 0
	}

	return score
}

// AnomalyDetector 异常检测器
type AnomalyDetector struct {
	sensitivity float64
	window      int
}

// NewAnomalyDetector 创建异常检测器
func NewAnomalyDetector() *AnomalyDetector {
	return &AnomalyDetector{
		sensitivity: 2.0, // 标准差倍数
		window:      10,  // 滑动窗口大小
	}
}

// Type 返回分析类型
func (d *AnomalyDetector) Type() AnalysisType {
	return AnalysisTypeAnomaly
}

// Options 返回配置选项
func (d *AnomalyDetector) Options() map[string]interface{} {
	return map[string]interface{}{
		"sensitivity": d.sensitivity,
		"window":      d.window,
	}
}

// Analyze 检测异常
func (d *AnomalyDetector) Analyze(ctx context.Context, data []collector.Data) (*Analysis, error) {
	analysis := &Analysis{
		Type:            AnalysisTypeAnomaly,
		Timestamp:       time.Now(),
		Results:         make(map[string]interface{}),
		Insights:        make([]Insight, 0),
		Recommendations: make([]Recommendation, 0),
	}

	// 按节点分组数据
	nodeData := make(map[string][]collector.Data)
	for _, d := range data {
		nodeData[d.Node] = append(nodeData[d.Node], d)
	}

	// 检测每个节点的异常
	anomalies := make([]map[string]interface{}, 0)
	for node, nodeMetrics := range nodeData {
		if anomaly := d.detectNodeAnomalies(node, nodeMetrics); anomaly != nil {
			anomalies = append(anomalies, anomaly)

			// 添加洞察
			analysis.Insights = append(analysis.Insights, Insight{
				Level:       "warning",
				Category:    "anomaly",
				Description: fmt.Sprintf("Anomaly detected on node %s", node),
				Details:     anomaly,
				Timestamp:   time.Now(),
			})
		}
	}

	analysis.Results["anomalies"] = anomalies
	analysis.Results["total_nodes"] = len(nodeData)
	analysis.Results["anomaly_nodes"] = len(anomalies)

	// 计算评分
	analysis.Score = 100.0 - (float64(len(anomalies))/float64(len(nodeData)))*50

	return analysis, nil
}

// detectNodeAnomalies 检测节点异常
func (d *AnomalyDetector) detectNodeAnomalies(node string, data []collector.Data) map[string]interface{} {
	anomalies := make(map[string]interface{})

	// 提取时间序列数据
	for _, d := range data {
		if d.Type == collector.DataTypePerformance {
			// 检查CPU异常
			if cpu, ok := d.Metrics["cpu_usage"].(float64); ok {
				if cpu > 95.0 {
					anomalies["high_cpu"] = cpu
				}
			}

			// 检查内存异常
			if mem, ok := d.Metrics["memory_usage"].(float64); ok {
				if mem > 95.0 {
					anomalies["high_memory"] = mem
				}
			}
		}
	}

	if len(anomalies) > 0 {
		anomalies["node"] = node
		anomalies["timestamp"] = time.Now()
		return anomalies
	}

	return nil
}

// 统计工具函数

func min(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	minVal := values[0]
	for _, v := range values[1:] {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func max(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	maxVal := values[0]
	for _, v := range values[1:] {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}

func avg(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func median(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	sorted := make([]float64, len(values))
	copy(sorted, values)
	sort.Float64s(sorted)

	n := len(sorted)
	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2
	}
	return sorted[n/2]
}

func stddev(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}

	mean := avg(values)
	sumSquares := 0.0

	for _, v := range values {
		diff := v - mean
		sumSquares += diff * diff
	}

	variance := sumSquares / float64(len(values))
	return math.Sqrt(variance)
}
