package collector

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// SystemMetrics 系统指标
type SystemMetrics struct {
	Timestamp time.Time         `json:"timestamp"`
	CPU       CPUMetrics        `json:"cpu"`
	Memory    MemoryMetrics     `json:"memory"`
	Disk      []DiskMetrics     `json:"disk"`
	Network   []NetworkMetrics  `json:"network"`
	Process   ProcessMetrics    `json:"process"`
	Custom    map[string]string `json:"custom,omitempty"`
}

// CPUMetrics CPU指标
type CPUMetrics struct {
	Cores     int     `json:"cores"`
	Usage     float64 `json:"usage"`
	LoadAvg1  float64 `json:"load_avg_1"`
	LoadAvg5  float64 `json:"load_avg_5"`
	LoadAvg15 float64 `json:"load_avg_15"`
}

// MemoryMetrics 内存指标
type MemoryMetrics struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	SwapTotal   uint64  `json:"swap_total"`
	SwapUsed    uint64  `json:"swap_used"`
}

// DiskMetrics 磁盘指标
type DiskMetrics struct {
	Device      string  `json:"device"`
	MountPoint  string  `json:"mount_point"`
	FSType      string  `json:"fs_type"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"used_percent"`
}

// NetworkMetrics 网络指标
type NetworkMetrics struct {
	Interface   string `json:"interface"`
	BytesSent   uint64 `json:"bytes_sent"`
	BytesRecv   uint64 `json:"bytes_recv"`
	PacketsSent uint64 `json:"packets_sent"`
	PacketsRecv uint64 `json:"packets_recv"`
	Errors      uint64 `json:"errors"`
}

// ProcessMetrics 进程指标
type ProcessMetrics struct {
	Total    int `json:"total"`
	Running  int `json:"running"`
	Sleeping int `json:"sleeping"`
	Zombie   int `json:"zombie"`
}

// MetricsCollector 指标采集器
type MetricsCollector struct {
	config Config
}

// NewMetricsCollector 创建指标采集器
func NewMetricsCollector(config Config) *MetricsCollector {
	return &MetricsCollector{
		config: config,
	}
}

// CollectMetrics 采集所有指标
func (mc *MetricsCollector) CollectMetrics() (*SystemMetrics, error) {
	metrics := &SystemMetrics{
		Timestamp: time.Now(),
		Custom:    make(map[string]string),
	}

	// 采集 CPU 指标
	cpuMetrics, err := mc.collectCPUMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to collect CPU metrics: %w", err)
	}
	metrics.CPU = cpuMetrics

	// 采集内存指标
	memMetrics, err := mc.collectMemoryMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to collect memory metrics: %w", err)
	}
	metrics.Memory = memMetrics

	// 采集磁盘指标
	diskMetrics, err := mc.collectDiskMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to collect disk metrics: %w", err)
	}
	metrics.Disk = diskMetrics

	// 采集网络指标
	netMetrics, err := mc.collectNetworkMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to collect network metrics: %w", err)
	}
	metrics.Network = netMetrics

	// 采集进程指标
	procMetrics, err := mc.collectProcessMetrics()
	if err != nil {
		return nil, fmt.Errorf("failed to collect process metrics: %w", err)
	}
	metrics.Process = procMetrics

	return metrics, nil
}

// collectCPUMetrics 采集CPU指标
func (mc *MetricsCollector) collectCPUMetrics() (CPUMetrics, error) {
	metrics := CPUMetrics{
		Cores: runtime.NumCPU(),
	}

	// 获取负载平均值（仅Linux）
	if runtime.GOOS == "linux" {
		out, err := exec.Command("cat", "/proc/loadavg").Output()
		if err == nil {
			fields := strings.Fields(string(out))
			if len(fields) >= 3 {
				metrics.LoadAvg1, _ = strconv.ParseFloat(fields[0], 64)
				metrics.LoadAvg5, _ = strconv.ParseFloat(fields[1], 64)
				metrics.LoadAvg15, _ = strconv.ParseFloat(fields[2], 64)
			}
		}
	}

	return metrics, nil
}

// collectMemoryMetrics 采集内存指标
func (mc *MetricsCollector) collectMemoryMetrics() (MemoryMetrics, error) {
	metrics := MemoryMetrics{}

	if runtime.GOOS == "linux" {
		out, err := exec.Command("free", "-b").Output()
		if err != nil {
			return metrics, err
		}

		lines := strings.Split(string(out), "\n")
		if len(lines) >= 2 {
			// 解析内存行
			memFields := strings.Fields(lines[1])
			if len(memFields) >= 7 {
				metrics.Total, _ = strconv.ParseUint(memFields[1], 10, 64)
				metrics.Used, _ = strconv.ParseUint(memFields[2], 10, 64)
				metrics.Available, _ = strconv.ParseUint(memFields[6], 10, 64)
				if metrics.Total > 0 {
					metrics.UsedPercent = float64(metrics.Used) / float64(metrics.Total) * 100
				}
			}
		}

		// 解析交换分区
		if len(lines) >= 3 {
			swapFields := strings.Fields(lines[2])
			if len(swapFields) >= 3 {
				metrics.SwapTotal, _ = strconv.ParseUint(swapFields[1], 10, 64)
				metrics.SwapUsed, _ = strconv.ParseUint(swapFields[2], 10, 64)
			}
		}
	}

	return metrics, nil
}

// collectDiskMetrics 采集磁盘指标
func (mc *MetricsCollector) collectDiskMetrics() ([]DiskMetrics, error) {
	var metrics []DiskMetrics

	if runtime.GOOS == "linux" {
		out, err := exec.Command("df", "-B1", "-T").Output()
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(out), "\n")
		for i, line := range lines {
			if i == 0 || line == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 7 {
				total, _ := strconv.ParseUint(fields[2], 10, 64)
				used, _ := strconv.ParseUint(fields[3], 10, 64)
				avail, _ := strconv.ParseUint(fields[4], 10, 64)

				var usedPercent float64
				if total > 0 {
					usedPercent = float64(used) / float64(total) * 100
				}

				metrics = append(metrics, DiskMetrics{
					Device:      fields[0],
					FSType:      fields[1],
					Total:       total,
					Used:        used,
					Available:   avail,
					UsedPercent: usedPercent,
					MountPoint:  fields[6],
				})
			}
		}
	}

	return metrics, nil
}

// collectNetworkMetrics 采集网络指标
func (mc *MetricsCollector) collectNetworkMetrics() ([]NetworkMetrics, error) {
	var metrics []NetworkMetrics

	if runtime.GOOS == "linux" {
		out, err := exec.Command("cat", "/proc/net/dev").Output()
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(out), "\n")
		for i, line := range lines {
			if i < 2 || line == "" {
				continue
			}

			parts := strings.Split(line, ":")
			if len(parts) != 2 {
				continue
			}

			iface := strings.TrimSpace(parts[0])
			fields := strings.Fields(parts[1])
			if len(fields) >= 16 {
				bytesRecv, _ := strconv.ParseUint(fields[0], 10, 64)
				packetsRecv, _ := strconv.ParseUint(fields[1], 10, 64)
				bytesSent, _ := strconv.ParseUint(fields[8], 10, 64)
				packetsSent, _ := strconv.ParseUint(fields[9], 10, 64)

				metrics = append(metrics, NetworkMetrics{
					Interface:   iface,
					BytesRecv:   bytesRecv,
					PacketsRecv: packetsRecv,
					BytesSent:   bytesSent,
					PacketsSent: packetsSent,
				})
			}
		}
	}

	return metrics, nil
}

// collectProcessMetrics 采集进程指标
func (mc *MetricsCollector) collectProcessMetrics() (ProcessMetrics, error) {
	metrics := ProcessMetrics{}

	if runtime.GOOS == "linux" {
		out, err := exec.Command("ps", "-eo", "state").Output()
		if err != nil {
			return metrics, err
		}

		lines := strings.Split(string(out), "\n")
		for i, line := range lines {
			if i == 0 || line == "" {
				continue
			}

			state := strings.TrimSpace(line)
			metrics.Total++

			switch state[0] {
			case 'R':
				metrics.Running++
			case 'S', 'D':
				metrics.Sleeping++
			case 'Z':
				metrics.Zombie++
			}
		}
	}

	return metrics, nil
}
