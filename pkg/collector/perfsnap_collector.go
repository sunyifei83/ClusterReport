package collector

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// PerfSnapCollector 包装 PerfSnap 的功能
type PerfSnapCollector struct {
	duration      int  // 采集持续时间（秒）
	generateFlame bool // 是否生成火焰图
}

// PerfSnapData 存储 PerfSnap 收集的数据
type PerfSnapData struct {
	Timestamp       string                `json:"timestamp" yaml:"timestamp"`
	Hostname        string                `json:"hostname" yaml:"hostname"`
	Uptime          string                `json:"uptime" yaml:"uptime"`
	LoadAverage     PerfSnapLoadAvg       `json:"load_average" yaml:"load_average"`
	VMStat          PerfSnapVMStat        `json:"vmstat" yaml:"vmstat"`
	CPUStats        []PerfSnapCPUStat     `json:"cpu_stats" yaml:"cpu_stats"`
	ProcessStats    []PerfSnapProcessStat `json:"process_stats" yaml:"process_stats"`
	DiskIOStats     []PerfSnapDiskIOStat  `json:"disk_io_stats" yaml:"disk_io_stats"`
	MemoryStats     PerfSnapMemoryStat    `json:"memory_stats" yaml:"memory_stats"`
	NetworkStats    []PerfSnapNetworkStat `json:"network_stats" yaml:"network_stats"`
	TCPStats        PerfSnapTCPStat       `json:"tcp_stats" yaml:"tcp_stats"`
	TopProcessesCPU []PerfSnapTopProcess  `json:"top_processes_cpu" yaml:"top_processes_cpu"`
	TopProcessesMem []PerfSnapTopProcess  `json:"top_processes_mem" yaml:"top_processes_mem"`
	DmesgErrors     []string              `json:"dmesg_errors" yaml:"dmesg_errors"`
	Issues          []PerfSnapIssue       `json:"issues" yaml:"issues"`
	Recommendations []string              `json:"recommendations" yaml:"recommendations"`
	FlameGraphPath  string                `json:"flame_graph_path,omitempty" yaml:"flame_graph_path,omitempty"`
	Version         string                `json:"perfsnap_version" yaml:"perfsnap_version"`
}

type PerfSnapLoadAvg struct {
	OneMin     float64 `json:"1min" yaml:"1min"`
	FiveMin    float64 `json:"5min" yaml:"5min"`
	FifteenMin float64 `json:"15min" yaml:"15min"`
}

type PerfSnapVMStat struct {
	RunQueue         int `json:"run_queue" yaml:"run_queue"`
	BlockedProcesses int `json:"blocked_processes" yaml:"blocked_processes"`
	ContextSwitches  int `json:"context_switches" yaml:"context_switches"`
	Interrupts       int `json:"interrupts" yaml:"interrupts"`
}

type PerfSnapCPUStat struct {
	CPU    string  `json:"cpu" yaml:"cpu"`
	User   float64 `json:"user" yaml:"user"`
	System float64 `json:"system" yaml:"system"`
	IOWait float64 `json:"iowait" yaml:"iowait"`
	Idle   float64 `json:"idle" yaml:"idle"`
}

type PerfSnapProcessStat struct {
	PID     int     `json:"pid" yaml:"pid"`
	User    string  `json:"user" yaml:"user"`
	Command string  `json:"command" yaml:"command"`
	CPUPct  float64 `json:"cpu_pct" yaml:"cpu_pct"`
}

type PerfSnapDiskIOStat struct {
	Device    string  `json:"device" yaml:"device"`
	TPS       float64 `json:"tps" yaml:"tps"`
	ReadKBps  float64 `json:"read_kbps" yaml:"read_kbps"`
	WriteKBps float64 `json:"write_kbps" yaml:"write_kbps"`
	AvgWait   float64 `json:"avg_wait" yaml:"avg_wait"`
	Util      float64 `json:"util_pct" yaml:"util_pct"`
}

type PerfSnapMemoryStat struct {
	TotalMB     int     `json:"total_mb" yaml:"total_mb"`
	UsedMB      int     `json:"used_mb" yaml:"used_mb"`
	FreeMB      int     `json:"free_mb" yaml:"free_mb"`
	BufferMB    int     `json:"buffer_mb" yaml:"buffer_mb"`
	CacheMB     int     `json:"cache_mb" yaml:"cache_mb"`
	UsedPercent float64 `json:"used_percent" yaml:"used_percent"`
}

type PerfSnapNetworkStat struct {
	Interface string  `json:"interface" yaml:"interface"`
	RxKBps    float64 `json:"rx_kbps" yaml:"rx_kbps"`
	TxKBps    float64 `json:"tx_kbps" yaml:"tx_kbps"`
	RxPckps   float64 `json:"rx_pckps" yaml:"rx_pckps"`
	TxPckps   float64 `json:"tx_pckps" yaml:"tx_pckps"`
}

type PerfSnapTCPStat struct {
	ActiveOpens  int `json:"active_opens" yaml:"active_opens"`
	PassiveOpens int `json:"passive_opens" yaml:"passive_opens"`
	Established  int `json:"established" yaml:"established"`
	TimeWait     int `json:"time_wait" yaml:"time_wait"`
	CloseWait    int `json:"close_wait" yaml:"close_wait"`
}

type PerfSnapTopProcess struct {
	PID     int     `json:"pid" yaml:"pid"`
	User    string  `json:"user" yaml:"user"`
	Command string  `json:"command" yaml:"command"`
	Value   float64 `json:"value" yaml:"value"` // CPU% 或 MEM%
}

type PerfSnapIssue struct {
	Severity    string `json:"severity" yaml:"severity"`
	Category    string `json:"category" yaml:"category"`
	Description string `json:"description" yaml:"description"`
}

// NewPerfSnapCollector 创建新的 PerfSnap 收集器
func NewPerfSnapCollector() *PerfSnapCollector {
	return &PerfSnapCollector{
		duration:      5, // 默认采集5秒
		generateFlame: false,
	}
}

// NewPerfSnapCollectorWithOptions 创建带选项的 PerfSnap 收集器
func NewPerfSnapCollectorWithOptions(duration int, generateFlame bool) *PerfSnapCollector {
	return &PerfSnapCollector{
		duration:      duration,
		generateFlame: generateFlame,
	}
}

// Collect 执行 PerfSnap 数据收集
func (c *PerfSnapCollector) Collect() (*PerfSnapData, error) {
	data := &PerfSnapData{
		Version:   "1.1.1",
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 收集各项性能数据
	data.Hostname = c.getHostname()
	data.Uptime = c.getUptime()
	data.LoadAverage = c.getLoadAverage()
	data.VMStat = c.getVMStat()
	data.CPUStats = c.getCPUStats()
	data.ProcessStats = c.getProcessStats()
	data.DiskIOStats = c.getDiskIOStats()
	data.MemoryStats = c.getMemoryStats()
	data.NetworkStats = c.getNetworkStats()
	data.TCPStats = c.getTCPStats()
	data.TopProcessesCPU = c.getTopProcessesByCPU()
	data.TopProcessesMem = c.getTopProcessesByMem()
	data.DmesgErrors = c.getDmesgErrors()

	// 分析性能问题
	data.Issues = c.analyzeIssues(data)
	data.Recommendations = c.generateRecommendations(data.Issues)

	// 生成火焰图（如果启用）
	if c.generateFlame {
		flamePath, err := c.generateFlameGraph()
		if err == nil {
			data.FlameGraphPath = flamePath
		}
	}

	return data, nil
}

// getHostname 获取主机名
func (c *PerfSnapCollector) getHostname() string {
	output, err := exec.Command("hostname").Output()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(output))
}

// getUptime 获取系统运行时间
func (c *PerfSnapCollector) getUptime() string {
	output, err := exec.Command("uptime", "-p").Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(output))
}

// getLoadAverage 获取负载平均值
func (c *PerfSnapCollector) getLoadAverage() PerfSnapLoadAvg {
	loadAvg := PerfSnapLoadAvg{}
	output, err := exec.Command("sh", "-c", "cat /proc/loadavg | awk '{print $1,$2,$3}'").Output()
	if err != nil {
		return loadAvg
	}

	fields := strings.Fields(string(output))
	if len(fields) >= 3 {
		fmt.Sscanf(fields[0], "%f", &loadAvg.OneMin)
		fmt.Sscanf(fields[1], "%f", &loadAvg.FiveMin)
		fmt.Sscanf(fields[2], "%f", &loadAvg.FifteenMin)
	}

	return loadAvg
}

// getVMStat 获取虚拟内存统计
func (c *PerfSnapCollector) getVMStat() PerfSnapVMStat {
	vmstat := PerfSnapVMStat{}

	// 使用 vmstat 命令获取数据
	output, err := exec.Command("vmstat", "1", "2").Output()
	if err != nil {
		return vmstat
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) >= 4 {
		// 取最后一行的数据
		fields := strings.Fields(lines[3])
		if len(fields) >= 11 {
			fmt.Sscanf(fields[0], "%d", &vmstat.RunQueue)
			fmt.Sscanf(fields[1], "%d", &vmstat.BlockedProcesses)
			fmt.Sscanf(fields[10], "%d", &vmstat.ContextSwitches)
			fmt.Sscanf(fields[11], "%d", &vmstat.Interrupts)
		}
	}

	return vmstat
}

// getCPUStats 获取CPU统计信息
func (c *PerfSnapCollector) getCPUStats() []PerfSnapCPUStat {
	var stats []PerfSnapCPUStat

	// 使用 mpstat 命令
	output, err := exec.Command("mpstat", "-P", "ALL", "1", "1").Output()
	if err != nil {
		return stats
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "all") || strings.Contains(line, "CPU") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 12 {
			stat := PerfSnapCPUStat{}
			stat.CPU = fields[2]
			fmt.Sscanf(fields[3], "%f", &stat.User)
			fmt.Sscanf(fields[5], "%f", &stat.System)
			fmt.Sscanf(fields[6], "%f", &stat.IOWait)
			fmt.Sscanf(fields[12], "%f", &stat.Idle)
			stats = append(stats, stat)
		}
	}

	return stats
}

// getProcessStats 获取进程CPU使用统计
func (c *PerfSnapCollector) getProcessStats() []PerfSnapProcessStat {
	var stats []PerfSnapProcessStat

	// 使用 pidstat 命令
	output, err := exec.Command("pidstat", "1", "1").Output()
	if err != nil {
		return stats
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "PID") || strings.Contains(line, "Average") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 8 {
			stat := PerfSnapProcessStat{}
			fmt.Sscanf(fields[2], "%d", &stat.PID)
			stat.User = fields[4]
			fmt.Sscanf(fields[6], "%f", &stat.CPUPct)
			stat.Command = fields[7]

			if stat.CPUPct > 1.0 { // 只记录CPU使用率>1%的进程
				stats = append(stats, stat)
			}
		}
	}

	return stats
}

// getDiskIOStats 获取磁盘IO统计
func (c *PerfSnapCollector) getDiskIOStats() []PerfSnapDiskIOStat {
	var stats []PerfSnapDiskIOStat

	// 使用 iostat 命令
	output, err := exec.Command("iostat", "-dx", "1", "2").Output()
	if err != nil {
		return stats
	}

	lines := strings.Split(string(output), "\n")
	deviceSection := false

	for _, line := range lines {
		if strings.Contains(line, "Device") {
			deviceSection = true
			continue
		}

		if deviceSection && line != "" {
			fields := strings.Fields(line)
			if len(fields) >= 14 {
				stat := PerfSnapDiskIOStat{}
				stat.Device = fields[0]
				fmt.Sscanf(fields[3], "%f", &stat.TPS)
				fmt.Sscanf(fields[4], "%f", &stat.ReadKBps)
				fmt.Sscanf(fields[5], "%f", &stat.WriteKBps)
				fmt.Sscanf(fields[9], "%f", &stat.AvgWait)
				fmt.Sscanf(fields[13], "%f", &stat.Util)
				stats = append(stats, stat)
			}
		}
	}

	return stats
}

// getMemoryStats 获取内存统计
func (c *PerfSnapCollector) getMemoryStats() PerfSnapMemoryStat {
	memStat := PerfSnapMemoryStat{}

	// 使用 free 命令
	output, err := exec.Command("free", "-m").Output()
	if err != nil {
		return memStat
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) >= 2 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 7 {
			fmt.Sscanf(fields[1], "%d", &memStat.TotalMB)
			fmt.Sscanf(fields[2], "%d", &memStat.UsedMB)
			fmt.Sscanf(fields[3], "%d", &memStat.FreeMB)
			fmt.Sscanf(fields[5], "%d", &memStat.BufferMB)
			fmt.Sscanf(fields[6], "%d", &memStat.CacheMB)

			if memStat.TotalMB > 0 {
				memStat.UsedPercent = float64(memStat.UsedMB) / float64(memStat.TotalMB) * 100
			}
		}
	}

	return memStat
}

// getNetworkStats 获取网络统计
func (c *PerfSnapCollector) getNetworkStats() []PerfSnapNetworkStat {
	var stats []PerfSnapNetworkStat

	// 使用 sar 命令
	output, err := exec.Command("sar", "-n", "DEV", "1", "1").Output()
	if err != nil {
		return stats
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "IFACE") || strings.Contains(line, "Average") || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 9 {
			iface := fields[1]
			if iface != "lo" { // 排除本地回环
				stat := PerfSnapNetworkStat{}
				stat.Interface = iface
				fmt.Sscanf(fields[4], "%f", &stat.RxKBps)
				fmt.Sscanf(fields[5], "%f", &stat.TxKBps)
				fmt.Sscanf(fields[2], "%f", &stat.RxPckps)
				fmt.Sscanf(fields[3], "%f", &stat.TxPckps)
				stats = append(stats, stat)
			}
		}
	}

	return stats
}

// getTCPStats 获取TCP连接统计
func (c *PerfSnapCollector) getTCPStats() PerfSnapTCPStat {
	tcpStat := PerfSnapTCPStat{}

	// 使用 ss 命令统计连接状态
	output, err := exec.Command("ss", "-s").Output()
	if err != nil {
		return tcpStat
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "estab") {
			fmt.Sscanf(line, "TCP: %d (estab %d", &tcpStat.Established, &tcpStat.Established)
		}
	}

	// 统计各种状态的连接
	output2, err := exec.Command("ss", "-tan").Output()
	if err == nil {
		lines := strings.Split(string(output2), "\n")
		for _, line := range lines {
			if strings.Contains(line, "TIME-WAIT") {
				tcpStat.TimeWait++
			} else if strings.Contains(line, "CLOSE-WAIT") {
				tcpStat.CloseWait++
			}
		}
	}

	return tcpStat
}

// getTopProcessesByCPU 获取CPU使用率最高的进程
func (c *PerfSnapCollector) getTopProcessesByCPU() []PerfSnapTopProcess {
	var processes []PerfSnapTopProcess

	output, err := exec.Command("ps", "aux", "--sort=-%cpu").Output()
	if err != nil {
		return processes
	}

	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines[1:] { // 跳过标题行
		if count >= 10 { // 只取前10个
			break
		}

		fields := strings.Fields(line)
		if len(fields) >= 11 {
			proc := PerfSnapTopProcess{}
			fmt.Sscanf(fields[1], "%d", &proc.PID)
			proc.User = fields[0]
			fmt.Sscanf(fields[2], "%f", &proc.Value)
			proc.Command = strings.Join(fields[10:], " ")

			if len(proc.Command) > 50 {
				proc.Command = proc.Command[:50] + "..."
			}

			processes = append(processes, proc)
			count++
		}
	}

	return processes
}

// getTopProcessesByMem 获取内存使用率最高的进程
func (c *PerfSnapCollector) getTopProcessesByMem() []PerfSnapTopProcess {
	var processes []PerfSnapTopProcess

	output, err := exec.Command("ps", "aux", "--sort=-%mem").Output()
	if err != nil {
		return processes
	}

	lines := strings.Split(string(output), "\n")
	count := 0
	for _, line := range lines[1:] { // 跳过标题行
		if count >= 10 { // 只取前10个
			break
		}

		fields := strings.Fields(line)
		if len(fields) >= 11 {
			proc := PerfSnapTopProcess{}
			fmt.Sscanf(fields[1], "%d", &proc.PID)
			proc.User = fields[0]
			fmt.Sscanf(fields[3], "%f", &proc.Value)
			proc.Command = strings.Join(fields[10:], " ")

			if len(proc.Command) > 50 {
				proc.Command = proc.Command[:50] + "..."
			}

			processes = append(processes, proc)
			count++
		}
	}

	return processes
}

// getDmesgErrors 获取dmesg错误信息
func (c *PerfSnapCollector) getDmesgErrors() []string {
	var errors []string

	output, err := exec.Command("dmesg", "-T", "-l", "err,crit,alert,emerg").Output()
	if err != nil {
		return errors
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if line != "" && len(errors) < 20 { // 最多保留20条错误
			errors = append(errors, line)
		}
	}

	return errors
}

// analyzeIssues 分析性能问题
func (c *PerfSnapCollector) analyzeIssues(data *PerfSnapData) []PerfSnapIssue {
	var issues []PerfSnapIssue

	// 检查负载
	if data.LoadAverage.OneMin > 10.0 {
		issues = append(issues, PerfSnapIssue{
			Severity:    "high",
			Category:    "load",
			Description: fmt.Sprintf("系统负载过高: 1分钟负载=%.2f", data.LoadAverage.OneMin),
		})
	}

	// 检查内存使用
	if data.MemoryStats.UsedPercent > 90 {
		issues = append(issues, PerfSnapIssue{
			Severity:    "high",
			Category:    "memory",
			Description: fmt.Sprintf("内存使用率过高: %.1f%%", data.MemoryStats.UsedPercent),
		})
	}

	// 检查磁盘IO
	for _, disk := range data.DiskIOStats {
		if disk.Util > 80 {
			issues = append(issues, PerfSnapIssue{
				Severity:    "medium",
				Category:    "disk",
				Description: fmt.Sprintf("磁盘 %s IO使用率过高: %.1f%%", disk.Device, disk.Util),
			})
		}
	}

	// 检查CPU
	for _, cpu := range data.CPUStats {
		if cpu.IOWait > 20 {
			issues = append(issues, PerfSnapIssue{
				Severity:    "medium",
				Category:    "cpu",
				Description: fmt.Sprintf("CPU %s IO等待过高: %.1f%%", cpu.CPU, cpu.IOWait),
			})
		}
	}

	return issues
}

// generateRecommendations 生成优化建议
func (c *PerfSnapCollector) generateRecommendations(issues []PerfSnapIssue) []string {
	var recommendations []string

	categoryMap := make(map[string]bool)
	for _, issue := range issues {
		categoryMap[issue.Category] = true
	}

	if categoryMap["load"] {
		recommendations = append(recommendations, "建议检查高负载进程，考虑优化或扩容")
	}

	if categoryMap["memory"] {
		recommendations = append(recommendations, "建议检查内存泄漏，考虑增加内存或优化应用")
	}

	if categoryMap["disk"] {
		recommendations = append(recommendations, "建议优化磁盘IO，考虑使用SSD或增加磁盘缓存")
	}

	if categoryMap["cpu"] {
		recommendations = append(recommendations, "建议检查磁盘性能，优化IO密集型操作")
	}

	return recommendations
}

// generateFlameGraph 生成火焰图
func (c *PerfSnapCollector) generateFlameGraph() (string, error) {
	timestamp := time.Now().Format("20060102_150405")
	outputPath := fmt.Sprintf("/tmp/flamegraph_%s.svg", timestamp)

	// 采集性能数据
	cmd := exec.Command("perf", "record", "-F", "99", "-a", "-g", "--", "sleep", fmt.Sprintf("%d", c.duration))
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// 生成火焰图（需要FlameGraph工具）
	// 这里简化实现，实际需要安装FlameGraph工具
	return outputPath, nil
}
