package collector

import (
	"testing"
	"time"
)

func TestNewBaseCollector(t *testing.T) {
	config := Config{
		Interval: 60,
		Targets:  []string{"localhost"},
		Enabled:  true,
	}

	collector := NewBaseCollector("test-collector", config)

	if collector == nil {
		t.Fatal("Expected collector to be created, got nil")
	}

	if collector.Name() != "test-collector" {
		t.Errorf("Expected name 'test-collector', got '%s'", collector.Name())
	}
}

func TestMetricsCollector(t *testing.T) {
	config := Config{
		Interval: 60,
		Enabled:  true,
	}

	mc := NewMetricsCollector(config)
	if mc == nil {
		t.Fatal("Expected MetricsCollector to be created, got nil")
	}
}

func TestCollectMetrics(t *testing.T) {
	config := Config{
		Interval: 60,
		Enabled:  true,
	}

	mc := NewMetricsCollector(config)
	metrics, err := mc.CollectMetrics()

	if err != nil {
		t.Logf("Collect metrics returned error (expected on non-Linux): %v", err)
		return
	}

	if metrics == nil {
		t.Fatal("Expected metrics to be non-nil")
	}

	// 验证时间戳
	if metrics.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}

	// 验证 CPU 指标
	if metrics.CPU.Cores <= 0 {
		t.Errorf("Expected positive CPU cores, got %d", metrics.CPU.Cores)
	}
}

func TestCPUMetrics(t *testing.T) {
	config := Config{}
	mc := NewMetricsCollector(config)

	metrics, err := mc.collectCPUMetrics()
	if err != nil {
		t.Fatalf("collectCPUMetrics failed: %v", err)
	}

	if metrics.Cores <= 0 {
		t.Errorf("Expected positive CPU cores, got %d", metrics.Cores)
	}
}

func TestMemoryMetrics(t *testing.T) {
	config := Config{}
	mc := NewMetricsCollector(config)

	_, err := mc.collectMemoryMetrics()
	if err != nil {
		t.Logf("collectMemoryMetrics returned error (expected on non-Linux): %v", err)
	}
}

func TestDiskMetrics(t *testing.T) {
	config := Config{}
	mc := NewMetricsCollector(config)

	metrics, err := mc.collectDiskMetrics()
	if err != nil {
		t.Logf("collectDiskMetrics returned error (expected on non-Linux): %v", err)
		return
	}

	// 在 Linux 系统上应该至少有一个磁盘
	if len(metrics) == 0 {
		t.Log("No disk metrics collected (may be expected on non-Linux)")
	}
}

func TestNetworkMetrics(t *testing.T) {
	config := Config{}
	mc := NewMetricsCollector(config)

	metrics, err := mc.collectNetworkMetrics()
	if err != nil {
		t.Logf("collectNetworkMetrics returned error (expected on non-Linux): %v", err)
		return
	}

	// 在 Linux 系统上应该至少有一个网络接口
	if len(metrics) == 0 {
		t.Log("No network metrics collected (may be expected on non-Linux)")
	}
}

func TestProcessMetrics(t *testing.T) {
	config := Config{}
	mc := NewMetricsCollector(config)

	metrics, err := mc.collectProcessMetrics()
	if err != nil {
		t.Logf("collectProcessMetrics returned error (expected on non-Linux): %v", err)
		return
	}

	// 应该有一些进程
	if metrics.Total <= 0 {
		t.Error("Expected at least some processes")
	}
}

func TestBaseCollectorCollect(t *testing.T) {
	config := Config{
		Interval: 60,
		Enabled:  true,
	}

	collector := NewBaseCollector("test", config)
	data, err := collector.Collect()

	if err != nil {
		t.Logf("Collect returned error (expected on non-Linux): %v", err)
		return
	}

	if data == nil {
		t.Error("Expected non-nil data")
	}

	// 验证数据类型
	if metrics, ok := data.(*SystemMetrics); ok {
		if metrics.Timestamp.IsZero() {
			t.Error("Expected non-zero timestamp")
		}
	} else {
		t.Error("Expected data to be *SystemMetrics type")
	}
}

func TestCollectData(t *testing.T) {
	config := Config{
		Interval: 60,
		Enabled:  true,
	}

	collector := NewBaseCollector("test", config)
	jsonData, err := collector.CollectData()

	if err != nil {
		t.Logf("CollectData returned error (expected on non-Linux): %v", err)
		return
	}

	if len(jsonData) == 0 {
		t.Error("Expected non-empty JSON data")
	}

	// 验证是否是有效的 JSON
	if jsonData[0] != '{' {
		t.Error("Expected JSON data to start with '{'")
	}
}

func TestGenerateReport(t *testing.T) {
	config := Config{
		Interval: 60,
		Enabled:  true,
	}

	collector := NewBaseCollector("test-report", config)
	report, err := collector.GenerateReport()

	if err != nil {
		t.Logf("GenerateReport returned error (expected on non-Linux): %v", err)
		// 即使有错误，report 也应该被创建
		if report == nil {
			t.Fatal("Expected report to be created even on error")
		}
		if report.Status != "failed" {
			t.Error("Expected status to be 'failed' on error")
		}
		if report.Error == "" {
			t.Error("Expected error message in report")
		}
		return
	}

	if report == nil {
		t.Fatal("Expected non-nil report")
	}

	if report.Status != "success" {
		t.Errorf("Expected status 'success', got '%s'", report.Status)
	}

	if report.Collector != "test-report" {
		t.Errorf("Expected collector name 'test-report', got '%s'", report.Collector)
	}

	if report.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp in report")
	}

	if report.Data == nil {
		t.Error("Expected non-nil data in report")
	}
}

func TestSystemMetricsStructure(t *testing.T) {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
		CPU: CPUMetrics{
			Cores:     4,
			Usage:     50.0,
			LoadAvg1:  1.5,
			LoadAvg5:  1.2,
			LoadAvg15: 1.0,
		},
		Memory: MemoryMetrics{
			Total:       16 * 1024 * 1024 * 1024,
			Used:        8 * 1024 * 1024 * 1024,
			Available:   8 * 1024 * 1024 * 1024,
			UsedPercent: 50.0,
		},
		Custom: map[string]string{
			"test": "value",
		},
	}

	if metrics.CPU.Cores != 4 {
		t.Errorf("Expected 4 cores, got %d", metrics.CPU.Cores)
	}

	if metrics.Memory.UsedPercent != 50.0 {
		t.Errorf("Expected 50%% memory usage, got %.2f", metrics.Memory.UsedPercent)
	}

	if metrics.Custom["test"] != "value" {
		t.Error("Expected custom field to be preserved")
	}
}

func BenchmarkCollectMetrics(b *testing.B) {
	config := Config{
		Interval: 60,
		Enabled:  true,
	}

	mc := NewMetricsCollector(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mc.CollectMetrics()
	}
}

func BenchmarkCollectCPUMetrics(b *testing.B) {
	config := Config{}
	mc := NewMetricsCollector(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mc.collectCPUMetrics()
	}
}
