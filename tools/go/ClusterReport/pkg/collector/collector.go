package collector

import (
	"encoding/json"
	"fmt"
	"time"
)

// Config 配置
type Config struct {
	Interval int      `yaml:"interval"`
	Targets  []string `yaml:"targets"`
	Enabled  bool     `yaml:"enabled"`
}

// Collector 采集器接口
type Collector interface {
	Collect() (interface{}, error)
	Name() string
}

// BaseCollector 基础采集器
type BaseCollector struct {
	name   string
	config Config
}

// NewBaseCollector 创建基础采集器
func NewBaseCollector(name string, config Config) *BaseCollector {
	return &BaseCollector{
		name:   name,
		config: config,
	}
}

// Collect 实现 Collector 接口
func (bc *BaseCollector) Collect() (interface{}, error) {
	// 使用 MetricsCollector 采集系统指标
	mc := NewMetricsCollector(bc.config)
	return mc.CollectMetrics()
}

// Name 返回采集器名称
func (bc *BaseCollector) Name() string {
	return bc.name
}

// CollectData 采集数据并返回 JSON
func (bc *BaseCollector) CollectData() ([]byte, error) {
	data, err := bc.Collect()
	if err != nil {
		return nil, fmt.Errorf("failed to collect data: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	return jsonData, nil
}

// Report 报告结构
type Report struct {
	Timestamp time.Time   `json:"timestamp"`
	Collector string      `json:"collector"`
	Data      interface{} `json:"data"`
	Status    string      `json:"status"`
	Error     string      `json:"error,omitempty"`
}

// GenerateReport 生成采集报告
func (bc *BaseCollector) GenerateReport() (*Report, error) {
	report := &Report{
		Timestamp: time.Now(),
		Collector: bc.name,
		Status:    "success",
	}

	data, err := bc.Collect()
	if err != nil {
		report.Status = "failed"
		report.Error = err.Error()
		return report, err
	}

	report.Data = data
	return report, nil
}
