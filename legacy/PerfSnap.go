/*
PerfSnap - Linux服务器性能快照分析工具
快速采集和分析系统性能数据，包括CPU、内存、磁盘、网络等关键指标
支持生成高CPU进程的火焰图进行性能分析
Author: sunyifei83@gmail.com
Version: 1.1.1
Tips:
- PerfSnap需要系统安装sysstat包（提供sar、mpstat、pidstat、iostat等命令）
- 火焰图功能需要安装perf工具和FlameGraph工具包
- v1.1.1修复了磁盘IO统计字段解析错误，利用率现在正确显示0-100%
使用:
1. go build -o perfsnap PerfSnap.go
2. chmod +x perfsnap
3. ./perfsnap 或 sudo ./perfsnap (推荐使用root权限)
*/
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PerformanceData 存储性能数据
type PerformanceData struct {
	Uptime      UptimeInfo
	DmesgErrors []string
	VMStat      VMStatInfo
	MPStat      MPStatInfo
	PIDStat     []ProcessInfo
	IOStat      IOStatInfo
	Memory      MemoryUsage
	Network     NetworkStats
	TCPStats    TCPConnectionInfo
	TopProcs    []TopProcess
}

// UptimeInfo 系统运行时间和负载
type UptimeInfo struct {
	Uptime    string
	Users     int
	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64
}

// VMStatInfo 虚拟内存统计
type VMStatInfo struct {
	RunQueue    int
	BlockedProc int
	SwapIn      int
	SwapOut     int
	BlockIn     int
	BlockOut    int
	Interrupts  int
	CtxSwitches int
	UserCPU     int
	SystemCPU   int
	IdleCPU     int
	WaitIO      int
	StolenCPU   int
}

// MPStatInfo 多核CPU统计
type MPStatInfo struct {
	CPUStats []CPUStat
}

type CPUStat struct {
	CPU    string
	User   float64
	Nice   float64
	System float64
	IOWait float64
	Idle   float64
}

// ProcessInfo 进程统计
type ProcessInfo struct {
	PID     int
	Command string
	CPU     float64
	Memory  float64
	DiskIO  float64
}

// IOStatInfo 磁盘IO统计
type IOStatInfo struct {
	Devices []DeviceIO
}

type DeviceIO struct {
	Device   string
	RRPS     float64 // 每秒读请求
	WRPS     float64 // 每秒写请求
	RkBPS    float64 // 每秒读KB
	WkBPS    float64 // 每秒写KB
	AvgQueue float64 // 平均队列长度
	AvgWait  float64 // 平均等待时间
	SvcTime  float64 // 平均服务时间
	Util     float64 // 利用率
}

// MemoryUsage 内存使用情况
type MemoryUsage struct {
	TotalMB     int
	UsedMB      int
	FreeMB      int
	SharedMB    int
	BuffersMB   int
	CachedMB    int
	AvailableMB int
	SwapTotalMB int
	SwapUsedMB  int
	SwapFreeMB  int
}

// NetworkStats 网络统计
type NetworkStats struct {
	Interfaces []InterfaceStats
}

type InterfaceStats struct {
	Interface string
	RxPPS     float64 // 接收包/秒
	TxPPS     float64 // 发送包/秒
	RxKBPS    float64 // 接收KB/秒
	TxKBPS    float64 // 发送KB/秒
	RxErrors  int
	TxErrors  int
}

// TCPConnectionInfo TCP连接信息
type TCPConnectionInfo struct {
	Active      int
	Passive     int
	Failed      int
	Resets      int
	Established int
	TimeWait    int
	CloseWait   int
	Retrans     int
}

// TopProcess TOP进程信息
type TopProcess struct {
	PID     int
	User    string
	CPU     float64
	Memory  float64
	Command string
}

// FlameGraphConfig 火焰图配置
type FlameGraphConfig struct {
	Enabled   bool
	PID       int    // 指定进程ID，0表示自动选择CPU最高的进程
	Duration  int    // 采样时长（秒）
	Frequency int    // 采样频率（Hz）
	OutputDir string // 输出目录
	AutoOpen  bool   // 是否自动打开生成的火焰图
}

// 执行命令并返回输出
func execCommand(command string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// 执行命令并返回输出（带超时）
func execCommandWithTimeout(timeout time.Duration, command string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// 获取系统运行时间和负载
func getUptime() UptimeInfo {
	info := UptimeInfo{}

	if output, err := execCommand("uptime"); err == nil {
		// 解析uptime输出
		// 示例: 14:26:39 up 5 days, 21:05,  2 users,  load average: 0.16, 0.05, 0.06
		re := regexp.MustCompile(`up\s+(.+?),\s+(\d+)\s+users?,\s+load average:\s+([\d.]+),\s+([\d.]+),\s+([\d.]+)`)
		if matches := re.FindStringSubmatch(output); len(matches) > 5 {
			info.Uptime = strings.TrimSpace(matches[1])
			info.Users, _ = strconv.Atoi(matches[2])
			info.LoadAvg1, _ = strconv.ParseFloat(matches[3], 64)
			info.LoadAvg5, _ = strconv.ParseFloat(matches[4], 64)
			info.LoadAvg15, _ = strconv.ParseFloat(matches[5], 64)
		}
	}

	return info
}

// 获取最近的dmesg错误
func getDmesgErrors() []string {
	var errors []string

	if output, err := execCommand("dmesg", "-T", "--level=err,warn", "--since", "1 minute ago"); err == nil {
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				errors = append(errors, line)
				if len(errors) >= 10 { // 限制显示最近10条
					break
				}
			}
		}
	}

	return errors
}

// 获取vmstat数据（1秒采样）
func getVMStat() VMStatInfo {
	info := VMStatInfo{}

	if output, err := execCommandWithTimeout(2*time.Second, "vmstat", "1", "2"); err == nil {
		lines := strings.Split(output, "\n")
		// 取第4行数据（第3行是第一次采样，第4行是第二次采样）
		if len(lines) >= 4 {
			fields := strings.Fields(lines[3])
			if len(fields) >= 17 {
				info.RunQueue, _ = strconv.Atoi(fields[0])
				info.BlockedProc, _ = strconv.Atoi(fields[1])
				info.SwapIn, _ = strconv.Atoi(fields[6])
				info.SwapOut, _ = strconv.Atoi(fields[7])
				info.BlockIn, _ = strconv.Atoi(fields[8])
				info.BlockOut, _ = strconv.Atoi(fields[9])
				info.Interrupts, _ = strconv.Atoi(fields[10])
				info.CtxSwitches, _ = strconv.Atoi(fields[11])
				info.UserCPU, _ = strconv.Atoi(fields[12])
				info.SystemCPU, _ = strconv.Atoi(fields[13])
				info.IdleCPU, _ = strconv.Atoi(fields[14])
				info.WaitIO, _ = strconv.Atoi(fields[15])
				if len(fields) >= 17 {
					info.StolenCPU, _ = strconv.Atoi(fields[16])
				}
			}
		}
	}

	return info
}

// 获取多核CPU统计
func getMPStat() MPStatInfo {
	info := MPStatInfo{}

	if output, err := execCommandWithTimeout(2*time.Second, "mpstat", "-P", "ALL", "1", "1"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Average:") && !strings.Contains(line, "CPU") {
				fields := strings.Fields(line)
				if len(fields) >= 12 {
					stat := CPUStat{
						CPU: fields[1],
					}
					stat.User, _ = strconv.ParseFloat(fields[2], 64)
					stat.Nice, _ = strconv.ParseFloat(fields[3], 64)
					stat.System, _ = strconv.ParseFloat(fields[4], 64)
					stat.IOWait, _ = strconv.ParseFloat(fields[5], 64)
					stat.Idle, _ = strconv.ParseFloat(fields[10], 64)
					info.CPUStats = append(info.CPUStats, stat)
				}
			}
		}
	}

	return info
}

// 获取进程统计
func getPIDStat() []ProcessInfo {
	var processes []ProcessInfo

	if output, err := execCommandWithTimeout(2*time.Second, "pidstat", "1", "1"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Average:") && !strings.Contains(line, "PID") {
				fields := strings.Fields(line)
				if len(fields) >= 8 {
					proc := ProcessInfo{}
					proc.PID, _ = strconv.Atoi(fields[1])
					proc.CPU, _ = strconv.ParseFloat(fields[6], 64)
					if len(fields) >= 9 {
						proc.Command = fields[8]
					}
					// 只保留CPU使用率大于1%的进程
					if proc.CPU > 1.0 {
						processes = append(processes, proc)
					}
				}
			}
		}
	}

	return processes
}

// 获取磁盘IO统计
func getIOStat() IOStatInfo {
	info := IOStatInfo{}

	if output, err := execCommandWithTimeout(2*time.Second, "iostat", "-xz", "1", "2"); err == nil {
		lines := strings.Split(output, "\n")
		inDeviceSection := false

		for _, line := range lines {
			if strings.Contains(line, "Device") {
				inDeviceSection = true
				continue
			}

			if inDeviceSection && line != "" && !strings.Contains(line, "avg-cpu:") {
				fields := strings.Fields(line)
				// iostat -xz 输出格式:
				// Device r/s w/s rkB/s wkB/s rrqm/s wrqm/s %rrqm %wrqm r_await w_await aqu-sz rareq-sz wareq-sz svctm %util
				if len(fields) >= 16 {
					dev := DeviceIO{
						Device: fields[0],
					}
					dev.RRPS, _ = strconv.ParseFloat(fields[1], 64)
					dev.WRPS, _ = strconv.ParseFloat(fields[2], 64)
					dev.RkBPS, _ = strconv.ParseFloat(fields[3], 64)
					dev.WkBPS, _ = strconv.ParseFloat(fields[4], 64)
					dev.AvgQueue, _ = strconv.ParseFloat(fields[11], 64) // aqu-sz

					// 计算平均等待时间 (读写等待时间的加权平均)
					rAwait, _ := strconv.ParseFloat(fields[9], 64)  // r_await
					wAwait, _ := strconv.ParseFloat(fields[10], 64) // w_await
					if dev.RRPS+dev.WRPS > 0 {
						dev.AvgWait = (rAwait*dev.RRPS + wAwait*dev.WRPS) / (dev.RRPS + dev.WRPS)
					}

					dev.SvcTime, _ = strconv.ParseFloat(fields[14], 64) // svctm
					dev.Util, _ = strconv.ParseFloat(fields[15], 64)    // %util

					// 只显示有活动的设备
					if dev.RRPS > 0 || dev.WRPS > 0 || dev.Util > 0 {
						info.Devices = append(info.Devices, dev)
					}
				}
			}
		}
	}

	return info
}

// 获取内存使用情况
func getMemoryUsage() MemoryUsage {
	info := MemoryUsage{}

	if output, err := execCommand("free", "-m"); err == nil {
		lines := strings.Split(output, "\n")
		for i, line := range lines {
			fields := strings.Fields(line)
			if i == 1 && len(fields) >= 7 { // Mem行
				info.TotalMB, _ = strconv.Atoi(fields[1])
				info.UsedMB, _ = strconv.Atoi(fields[2])
				info.FreeMB, _ = strconv.Atoi(fields[3])
				info.SharedMB, _ = strconv.Atoi(fields[4])
				info.BuffersMB, _ = strconv.Atoi(fields[5])
				info.CachedMB, _ = strconv.Atoi(fields[6])
				if len(fields) >= 7 {
					info.AvailableMB, _ = strconv.Atoi(fields[6])
				}
			} else if i == 2 && len(fields) >= 4 { // Swap行
				info.SwapTotalMB, _ = strconv.Atoi(fields[1])
				info.SwapUsedMB, _ = strconv.Atoi(fields[2])
				info.SwapFreeMB, _ = strconv.Atoi(fields[3])
			}
		}
	}

	return info
}

// 获取网络统计
func getNetworkStats() NetworkStats {
	info := NetworkStats{}

	if output, err := execCommandWithTimeout(2*time.Second, "sar", "-n", "DEV", "1", "1"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Average:") && !strings.Contains(line, "IFACE") {
				fields := strings.Fields(line)
				if len(fields) >= 10 {
					iface := fields[1]
					// 跳过lo接口
					if iface == "lo" {
						continue
					}

					stat := InterfaceStats{
						Interface: iface,
					}
					stat.RxPPS, _ = strconv.ParseFloat(fields[2], 64)
					stat.TxPPS, _ = strconv.ParseFloat(fields[3], 64)
					stat.RxKBPS, _ = strconv.ParseFloat(fields[4], 64)
					stat.TxKBPS, _ = strconv.ParseFloat(fields[5], 64)

					// 只显示有活动的接口
					if stat.RxPPS > 0 || stat.TxPPS > 0 {
						info.Interfaces = append(info.Interfaces, stat)
					}
				}
			}
		}
	}

	return info
}

// 获取TCP连接统计
func getTCPStats() TCPConnectionInfo {
	info := TCPConnectionInfo{}

	if output, err := execCommandWithTimeout(2*time.Second, "sar", "-n", "TCP,ETCP", "1", "1"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Average:") && strings.Contains(line, "active/s") {
				fields := strings.Fields(line)
				if len(fields) >= 5 {
					info.Active, _ = strconv.Atoi(fields[1])
					info.Passive, _ = strconv.Atoi(fields[2])
				}
			}
		}
	}

	// 获取连接状态统计
	if output, err := execCommand("ss", "-s"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "estab") {
				re := regexp.MustCompile(`estab\s+(\d+)`)
				if matches := re.FindStringSubmatch(line); len(matches) > 1 {
					info.Established, _ = strconv.Atoi(matches[1])
				}
			}
			if strings.Contains(line, "time-wait") {
				re := regexp.MustCompile(`time-wait\s+(\d+)`)
				if matches := re.FindStringSubmatch(line); len(matches) > 1 {
					info.TimeWait, _ = strconv.Atoi(matches[1])
				}
			}
			if strings.Contains(line, "close-wait") {
				re := regexp.MustCompile(`close-wait\s+(\d+)`)
				if matches := re.FindStringSubmatch(line); len(matches) > 1 {
					info.CloseWait, _ = strconv.Atoi(matches[1])
				}
			}
		}
	}

	return info
}

// 获取TOP进程
func getTopProcesses() []TopProcess {
	var processes []TopProcess

	if output, err := execCommand("top", "-bn1", "-o", "%CPU"); err == nil {
		lines := strings.Split(output, "\n")
		count := 0
		for i, line := range lines {
			// 跳过前面的头部信息
			if i < 7 {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 12 {
				proc := TopProcess{}
				proc.PID, _ = strconv.Atoi(fields[0])
				proc.User = fields[1]
				proc.CPU, _ = strconv.ParseFloat(fields[8], 64)
				proc.Memory, _ = strconv.ParseFloat(fields[9], 64)
				proc.Command = fields[11]

				// 只保留前10个高CPU使用进程
				if proc.CPU > 0 {
					processes = append(processes, proc)
					count++
					if count >= 10 {
						break
					}
				}
			}
		}
	}

	return processes
}

// 并发收集所有性能数据
func collectPerformanceData() PerformanceData {
	data := PerformanceData{}
	var wg sync.WaitGroup

	// 使用goroutine并发收集数据
	wg.Add(10)

	go func() {
		defer wg.Done()
		data.Uptime = getUptime()
	}()

	go func() {
		defer wg.Done()
		data.DmesgErrors = getDmesgErrors()
	}()

	go func() {
		defer wg.Done()
		data.VMStat = getVMStat()
	}()

	go func() {
		defer wg.Done()
		data.MPStat = getMPStat()
	}()

	go func() {
		defer wg.Done()
		data.PIDStat = getPIDStat()
	}()

	go func() {
		defer wg.Done()
		data.IOStat = getIOStat()
	}()

	go func() {
		defer wg.Done()
		data.Memory = getMemoryUsage()
	}()

	go func() {
		defer wg.Done()
		data.Network = getNetworkStats()
	}()

	go func() {
		defer wg.Done()
		data.TCPStats = getTCPStats()
	}()

	go func() {
		defer wg.Done()
		data.TopProcs = getTopProcesses()
	}()

	wg.Wait()
	return data
}

// 打印性能报告
func printPerformanceReport(data PerformanceData) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("                    PerfSnap - 系统性能快照分析报告")
	fmt.Println("                    " + time.Now().Format("2006-01-02 15:04:05"))
	fmt.Println(strings.Repeat("=", 80))

	// 系统运行时间和负载
	fmt.Printf("\n【系统概况】\n")
	fmt.Printf("运行时间: %s | 在线用户: %d\n", data.Uptime.Uptime, data.Uptime.Users)
	fmt.Printf("系统负载: %.2f (1分钟) | %.2f (5分钟) | %.2f (15分钟)\n",
		data.Uptime.LoadAvg1, data.Uptime.LoadAvg5, data.Uptime.LoadAvg15)

	// CPU使用情况
	cpuCount := len(data.MPStat.CPUStats)
	if cpuCount > 0 {
		fmt.Printf("CPU核心数: %d\n", cpuCount-1) // 减去ALL统计
	}

	// 负载评估
	if cpuCount > 0 && data.Uptime.LoadAvg1 > float64(cpuCount-1) {
		fmt.Printf("⚠️  警告: 系统负载过高 (负载 %.2f > CPU核心数 %d)\n",
			data.Uptime.LoadAvg1, cpuCount-1)
	}

	// 内存使用情况
	fmt.Printf("\n【内存状态】\n")
	usedPercent := float64(data.Memory.UsedMB) / float64(data.Memory.TotalMB) * 100
	fmt.Printf("总内存: %d MB | 已使用: %d MB (%.1f%%) | 可用: %d MB\n",
		data.Memory.TotalMB, data.Memory.UsedMB, usedPercent, data.Memory.AvailableMB)
	fmt.Printf("缓存: %d MB | 缓冲: %d MB\n", data.Memory.CachedMB, data.Memory.BuffersMB)

	if data.Memory.SwapTotalMB > 0 {
		swapUsedPercent := float64(data.Memory.SwapUsedMB) / float64(data.Memory.SwapTotalMB) * 100
		fmt.Printf("交换分区: %d MB | 已使用: %d MB (%.1f%%)\n",
			data.Memory.SwapTotalMB, data.Memory.SwapUsedMB, swapUsedPercent)
		if swapUsedPercent > 50 {
			fmt.Printf("⚠️  警告: 交换分区使用率过高 (%.1f%%)\n", swapUsedPercent)
		}
	}

	// CPU详细统计
	fmt.Printf("\n【CPU统计】\n")
	fmt.Printf("%-6s %8s %8s %8s %8s %8s\n", "CPU", "User%", "System%", "IOWait%", "Idle%", "状态")
	fmt.Println(strings.Repeat("-", 60))

	for _, cpu := range data.MPStat.CPUStats {
		status := "正常"
		if cpu.IOWait > 30 {
			status = "⚠️ IO等待高"
		} else if cpu.User+cpu.System > 90 {
			status = "⚠️ 使用率高"
		}
		fmt.Printf("%-6s %8.1f %8.1f %8.1f %8.1f %s\n",
			cpu.CPU, cpu.User, cpu.System, cpu.IOWait, cpu.Idle, status)
	}

	// VM统计
	fmt.Printf("\n【虚拟内存统计】\n")
	fmt.Printf("运行队列: %d | 阻塞进程: %d\n", data.VMStat.RunQueue, data.VMStat.BlockedProc)
	fmt.Printf("CPU使用率: 用户 %d%% | 系统 %d%% | 空闲 %d%% | IO等待 %d%%\n",
		data.VMStat.UserCPU, data.VMStat.SystemCPU, data.VMStat.IdleCPU, data.VMStat.WaitIO)
	fmt.Printf("上下文切换: %d/秒 | 中断: %d/秒\n", data.VMStat.CtxSwitches, data.VMStat.Interrupts)

	if data.VMStat.CtxSwitches > 100000 {
		fmt.Printf("⚠️  警告: 上下文切换过于频繁 (%d/秒)\n", data.VMStat.CtxSwitches)
	}

	// 磁盘IO统计
	if len(data.IOStat.Devices) > 0 {
		fmt.Printf("\n【磁盘IO统计】\n")
		fmt.Printf("%-12s %8s %8s %10s %10s %8s %8s\n",
			"设备", "读IOPS", "写IOPS", "读KB/s", "写KB/s", "队列", "使用率%")
		fmt.Println(strings.Repeat("-", 70))

		for _, dev := range data.IOStat.Devices {
			status := ""
			if dev.Util > 90 {
				status = " ⚠️"
			}
			fmt.Printf("%-12s %8.1f %8.1f %10.1f %10.1f %8.1f %7.1f%%%s\n",
				dev.Device, dev.RRPS, dev.WRPS, dev.RkBPS, dev.WkBPS,
				dev.AvgQueue, dev.Util, status)
		}
	}

	// 网络统计
	if len(data.Network.Interfaces) > 0 {
		fmt.Printf("\n【网络流量】\n")
		fmt.Printf("%-12s %10s %10s %12s %12s\n",
			"接口", "接收pps", "发送pps", "接收KB/s", "发送KB/s")
		fmt.Println(strings.Repeat("-", 60))

		for _, iface := range data.Network.Interfaces {
			fmt.Printf("%-12s %10.1f %10.1f %12.1f %12.1f\n",
				iface.Interface, iface.RxPPS, iface.TxPPS, iface.RxKBPS, iface.TxKBPS)
		}
	}

	// TCP连接统计
	fmt.Printf("\n【TCP连接】\n")
	fmt.Printf("已建立: %d | TIME_WAIT: %d | CLOSE_WAIT: %d\n",
		data.TCPStats.Established, data.TCPStats.TimeWait, data.TCPStats.CloseWait)

	if data.TCPStats.TimeWait > 10000 {
		fmt.Printf("⚠️  警告: TIME_WAIT连接过多 (%d)\n", data.TCPStats.TimeWait)
	}

	// TOP进程
	if len(data.TopProcs) > 0 {
		fmt.Printf("\n【TOP进程 (按CPU排序)】\n")
		fmt.Printf("%-8s %-10s %8s %8s %-20s\n", "PID", "用户", "CPU%", "内存%", "命令")
		fmt.Println(strings.Repeat("-", 60))

		for i, proc := range data.TopProcs {
			if i >= 5 { // 只显示前5个
				break
			}
			fmt.Printf("%-8d %-10s %8.1f %8.1f %-20s\n",
				proc.PID, proc.User, proc.CPU, proc.Memory, proc.Command)
		}
	}

	// 高CPU进程 (pidstat)
	if len(data.PIDStat) > 0 {
		fmt.Printf("\n【高CPU进程 (>1%%)】\n")
		fmt.Printf("%-8s %8s %-30s\n", "PID", "CPU%", "命令")
		fmt.Println(strings.Repeat("-", 50))

		for _, proc := range data.PIDStat {
			fmt.Printf("%-8d %8.1f %-30s\n", proc.PID, proc.CPU, proc.Command)
		}
	}

	// 最近的系统错误
	if len(data.DmesgErrors) > 0 {
		fmt.Printf("\n【最近1分钟系统日志错误/警告】\n")
		fmt.Println(strings.Repeat("-", 60))
		for i, err := range data.DmesgErrors {
			if i >= 5 { // 只显示前5条
				break
			}
			// 截断过长的错误信息
			if len(err) > 75 {
				err = err[:75] + "..."
			}
			fmt.Printf("• %s\n", err)
		}
	}

	// 性能问题总结
	fmt.Printf("\n【性能分析总结】\n")
	fmt.Println(strings.Repeat("-", 60))

	issues := []string{}

	// 检查各项指标
	if cpuCount > 0 && data.Uptime.LoadAvg1 > float64(cpuCount-1)*1.5 {
		issues = append(issues, fmt.Sprintf("系统负载过高 (%.2f)", data.Uptime.LoadAvg1))
	}

	if data.VMStat.WaitIO > 30 {
		issues = append(issues, fmt.Sprintf("IO等待过高 (%d%%)", data.VMStat.WaitIO))
	}

	if data.Memory.SwapTotalMB > 0 {
		swapUsedPercent := float64(data.Memory.SwapUsedMB) / float64(data.Memory.SwapTotalMB) * 100
		if swapUsedPercent > 50 {
			issues = append(issues, fmt.Sprintf("交换分区使用过多 (%.1f%%)", swapUsedPercent))
		}
	}

	if data.VMStat.CtxSwitches > 100000 {
		issues = append(issues, fmt.Sprintf("上下文切换过于频繁 (%d/秒)", data.VMStat.CtxSwitches))
	}

	memUsedPercent := float64(data.Memory.UsedMB) / float64(data.Memory.TotalMB) * 100
	if memUsedPercent > 90 {
		issues = append(issues, fmt.Sprintf("内存使用率过高 (%.1f%%)", memUsedPercent))
	}

	// 检查磁盘利用率
	for _, dev := range data.IOStat.Devices {
		if dev.Util > 90 {
			issues = append(issues, fmt.Sprintf("磁盘 %s 利用率过高 (%.1f%%)", dev.Device, dev.Util))
		}
	}

	if data.TCPStats.TimeWait > 10000 {
		issues = append(issues, fmt.Sprintf("TIME_WAIT连接过多 (%d)", data.TCPStats.TimeWait))
	}

	// 打印问题
	if len(issues) > 0 {
		fmt.Println("⚠️  发现以下性能问题:")
		for _, issue := range issues {
			fmt.Printf("  • %s\n", issue)
		}
	} else {
		fmt.Println("✅ 系统性能状态良好")
	}

	// 建议
	fmt.Printf("\n【优化建议】\n")
	fmt.Println(strings.Repeat("-", 60))

	if data.VMStat.WaitIO > 30 {
		fmt.Println("• IO等待过高: 检查磁盘性能，考虑使用SSD或优化IO密集型进程")
	}

	if cpuCount > 0 && data.Uptime.LoadAvg1 > float64(cpuCount-1)*1.5 {
		fmt.Println("• 系统负载过高: 检查CPU密集型进程，考虑增加CPU资源或优化代码")
	}

	if memUsedPercent > 90 {
		fmt.Println("• 内存使用率过高: 检查内存泄漏，考虑增加内存或优化内存使用")
	}

	if data.Memory.SwapTotalMB > 0 {
		swapUsedPercent := float64(data.Memory.SwapUsedMB) / float64(data.Memory.SwapTotalMB) * 100
		if swapUsedPercent > 50 {
			fmt.Println("• 交换分区使用过多: 增加物理内存或调整swappiness参数")
		}
	}

	if data.VMStat.CtxSwitches > 100000 {
		fmt.Println("• 上下文切换频繁: 检查线程数量，优化并发模型")
	}

	if data.TCPStats.TimeWait > 10000 {
		fmt.Println("• TIME_WAIT过多: 优化TCP参数，如tcp_tw_reuse和tcp_tw_recycle")
	}

	fmt.Println("\n" + strings.Repeat("=", 80))
}

// 实时监控模式
func monitorMode(interval int, duration int) {
	fmt.Printf("开始实时监控 (间隔: %d秒, 持续: %d秒)\n", interval, duration)
	fmt.Println(strings.Repeat("=", 80))

	endTime := time.Now().Add(time.Duration(duration) * time.Second)

	for time.Now().Before(endTime) {
		data := collectPerformanceData()

		// 清屏
		fmt.Print("\033[H\033[2J")

		// 打印简化的实时数据
		fmt.Printf("时间: %s | 剩余: %ds\n",
			time.Now().Format("15:04:05"),
			int(endTime.Sub(time.Now()).Seconds()))
		fmt.Println(strings.Repeat("-", 60))

		// 显示关键指标
		fmt.Printf("负载: %.2f %.2f %.2f | ",
			data.Uptime.LoadAvg1, data.Uptime.LoadAvg5, data.Uptime.LoadAvg15)

		memUsedPercent := float64(data.Memory.UsedMB) / float64(data.Memory.TotalMB) * 100
		fmt.Printf("内存: %.1f%% | ", memUsedPercent)

		fmt.Printf("CPU: 用户%d%% 系统%d%% 空闲%d%% IO等待%d%%\n",
			data.VMStat.UserCPU, data.VMStat.SystemCPU,
			data.VMStat.IdleCPU, data.VMStat.WaitIO)

		// 显示活跃的磁盘IO
		if len(data.IOStat.Devices) > 0 {
			fmt.Printf("\n磁盘IO:\n")
			for _, dev := range data.IOStat.Devices {
				if dev.Util > 10 { // 只显示利用率>10%的设备
					fmt.Printf("  %s: 读%.1f/s 写%.1f/s 利用率%.1f%%\n",
						dev.Device, dev.RRPS, dev.WRPS, dev.Util)
				}
			}
		}

		// 显示网络流量
		if len(data.Network.Interfaces) > 0 {
			fmt.Printf("\n网络:\n")
			for _, iface := range data.Network.Interfaces {
				fmt.Printf("  %s: ↓%.1fKB/s ↑%.1fKB/s\n",
					iface.Interface, iface.RxKBPS, iface.TxKBPS)
			}
		}

		// 显示TOP进程
		if len(data.TopProcs) > 0 {
			fmt.Printf("\nTOP进程 (CPU):\n")
			for i, proc := range data.TopProcs {
				if i >= 3 { // 只显示前3个
					break
				}
				fmt.Printf("  %d %s CPU:%.1f%% MEM:%.1f%%\n",
					proc.PID, proc.Command, proc.CPU, proc.Memory)
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	fmt.Println("\n监控结束")
}

// 检查并安装FlameGraph工具
func checkAndInstallFlameGraph() (string, error) {
	// 检查常见的FlameGraph安装位置
	commonPaths := []string{
		"/opt/FlameGraph",
		"/usr/local/FlameGraph",
		os.Getenv("HOME") + "/FlameGraph",
		"./FlameGraph",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(filepath.Join(path, "flamegraph.pl")); err == nil {
			return path, nil
		}
	}

	// 如果没找到，尝试在当前用户目录下克隆
	installPath := os.Getenv("HOME") + "/FlameGraph"
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		fmt.Println("正在安装FlameGraph工具...")
		cmd := exec.Command("git", "clone", "https://github.com/brendangregg/FlameGraph.git", installPath)
		if err := cmd.Run(); err != nil {
			// 如果git clone失败，尝试使用wget下载关键脚本
			fmt.Println("Git克隆失败，尝试直接下载脚本...")
			os.MkdirAll(installPath, 0755)

			scripts := []string{"stackcollapse-perf.pl", "flamegraph.pl"}
			for _, script := range scripts {
				url := fmt.Sprintf("http://rzkkr9sg3.hd-bkt.clouddn.com/test_tools/%s", script)
				cmd := exec.Command("wget", "-O", filepath.Join(installPath, script), url)
				if err := cmd.Run(); err != nil {
					return "", fmt.Errorf("无法安装FlameGraph: %v", err)
				}
				// 添加执行权限
				os.Chmod(filepath.Join(installPath, script), 0755)
			}
		}
		fmt.Printf("FlameGraph工具已安装到: %s\n", installPath)
	}

	return installPath, nil
}

// 生成火焰图
func generateFlameGraph(config FlameGraphConfig, topProcs []TopProcess) error {
	// 检查perf命令是否可用
	if _, err := exec.LookPath("perf"); err != nil {
		return fmt.Errorf("perf工具未安装，请先安装: sudo apt-get install linux-tools-common linux-tools-generic 或 sudo yum install perf")
	}

	// 获取FlameGraph工具路径
	flameGraphPath, err := checkAndInstallFlameGraph()
	if err != nil {
		return err
	}

	// 确定目标进程
	targetPID := config.PID
	processName := ""

	if targetPID == 0 && len(topProcs) > 0 {
		// 自动选择CPU使用率最高的进程
		// 对进程按CPU使用率排序
		sort.Slice(topProcs, func(i, j int) bool {
			return topProcs[i].CPU > topProcs[j].CPU
		})
		targetPID = topProcs[0].PID
		processName = topProcs[0].Command
		fmt.Printf("\n自动选择CPU使用率最高的进程: PID=%d, Command=%s, CPU=%.1f%%\n",
			targetPID, processName, topProcs[0].CPU)
	} else if targetPID > 0 {
		// 查找进程名
		for _, proc := range topProcs {
			if proc.PID == targetPID {
				processName = proc.Command
				break
			}
		}
		if processName == "" {
			// 尝试从/proc获取进程名
			if cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", targetPID)); err == nil {
				processName = strings.TrimSpace(string(cmdline))
			} else {
				processName = fmt.Sprintf("pid_%d", targetPID)
			}
		}
	} else {
		return fmt.Errorf("没有找到合适的进程用于生成火焰图")
	}

	// 创建输出目录
	if config.OutputDir == "" {
		config.OutputDir = fmt.Sprintf("flamegraph_%s", time.Now().Format("20060102_150405"))
	}
	os.MkdirAll(config.OutputDir, 0755)

	// 生成文件名
	timestamp := time.Now().Format("20060102_150405")
	perfDataFile := filepath.Join(config.OutputDir, fmt.Sprintf("perf_%s_%d.data", processName, targetPID))
	perfStackFile := filepath.Join(config.OutputDir, fmt.Sprintf("perf_%s_%d.stacks", processName, targetPID))
	foldedFile := filepath.Join(config.OutputDir, fmt.Sprintf("perf_%s_%d.folded", processName, targetPID))
	svgFile := filepath.Join(config.OutputDir, fmt.Sprintf("flamegraph_%s_%d_%s.svg", processName, targetPID, timestamp))

	fmt.Printf("\n开始生成火焰图...\n")
	fmt.Printf("目标进程: PID=%d (%s)\n", targetPID, processName)
	fmt.Printf("采样时长: %d秒\n", config.Duration)
	fmt.Printf("采样频率: %dHz\n", config.Frequency)
	fmt.Printf("输出目录: %s\n", config.OutputDir)

	// 步骤1: 使用perf record采集数据
	fmt.Printf("\n[1/4] 正在采集性能数据 (%d秒)...\n", config.Duration)
	recordCmd := exec.Command("perf", "record",
		"-F", strconv.Itoa(config.Frequency),
		"-p", strconv.Itoa(targetPID),
		"-g",
		"-o", perfDataFile,
		"--", "sleep", strconv.Itoa(config.Duration))

	recordCmd.Stdout = os.Stdout
	recordCmd.Stderr = os.Stderr

	if err := recordCmd.Run(); err != nil {
		// 如果失败，可能是权限问题，尝试调整kernel.perf_event_paranoid
		fmt.Println("尝试调整内核参数...")
		exec.Command("sysctl", "-w", "kernel.perf_event_paranoid=-1").Run()

		// 重试
		if err := recordCmd.Run(); err != nil {
			return fmt.Errorf("perf record失败: %v", err)
		}
	}

	// 步骤2: 生成调用栈
	fmt.Println("[2/4] 正在解析调用栈...")
	scriptCmd := exec.Command("perf", "script", "-i", perfDataFile)
	stackOutput, err := scriptCmd.Output()
	if err != nil {
		return fmt.Errorf("perf script失败: %v", err)
	}

	// 写入栈文件
	if err := os.WriteFile(perfStackFile, stackOutput, 0644); err != nil {
		return fmt.Errorf("写入栈文件失败: %v", err)
	}

	// 步骤3: 折叠调用栈
	fmt.Println("[3/4] 正在折叠调用栈...")
	collapseScript := filepath.Join(flameGraphPath, "stackcollapse-perf.pl")
	collapseCmd := exec.Command("perl", collapseScript)
	collapseCmd.Stdin = strings.NewReader(string(stackOutput))

	foldedOutput, err := collapseCmd.Output()
	if err != nil {
		// 如果perl脚本执行失败，尝试使用简单的折叠方法
		fmt.Println("使用备用折叠方法...")
		foldedOutput = simpleFoldStacks(string(stackOutput))
	}

	if err := os.WriteFile(foldedFile, foldedOutput, 0644); err != nil {
		return fmt.Errorf("写入折叠文件失败: %v", err)
	}

	// 步骤4: 生成火焰图SVG
	fmt.Println("[4/4] 正在生成火焰图...")
	flameGraphScript := filepath.Join(flameGraphPath, "flamegraph.pl")
	flameCmd := exec.Command("perl", flameGraphScript,
		"--title", fmt.Sprintf("CPU Flame Graph: %s (PID %d)", processName, targetPID),
		"--width", "1200")
	flameCmd.Stdin = strings.NewReader(string(foldedOutput))

	svgOutput, err := flameCmd.Output()
	if err != nil {
		return fmt.Errorf("生成火焰图失败: %v", err)
	}

	if err := os.WriteFile(svgFile, svgOutput, 0644); err != nil {
		return fmt.Errorf("写入SVG文件失败: %v", err)
	}

	fmt.Printf("\n✅ 火焰图生成成功: %s\n", svgFile)

	// 生成简单的HTML查看器
	htmlFile := filepath.Join(config.OutputDir, fmt.Sprintf("viewer_%s.html", timestamp))
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>火焰图 - %s (PID %d)</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: #f5f5f5;
        }
        h1 {
            color: #333;
        }
        .info {
            background-color: white;
            padding: 15px;
            border-radius: 5px;
            margin-bottom: 20px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        .svg-container {
            background-color: white;
            padding: 10px;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        iframe {
            width: 100%%;
            height: 800px;
            border: none;
        }
    </style>
</head>
<body>
    <h1>CPU 火焰图分析</h1>
    <div class="info">
        <h2>进程信息</h2>
        <p><strong>进程名:</strong> %s</p>
        <p><strong>PID:</strong> %d</p>
        <p><strong>采样时长:</strong> %d 秒</p>
        <p><strong>采样频率:</strong> %d Hz</p>
        <p><strong>生成时间:</strong> %s</p>
    </div>
    <div class="svg-container">
        <h2>火焰图</h2>
        <iframe src="%s"></iframe>
    </div>
    <div class="info">
        <h3>如何阅读火焰图</h3>
        <ul>
            <li>X轴表示采样数量的比例，越宽表示执行时间越长</li>
            <li>Y轴表示调用栈深度，从下到上是调用关系</li>
            <li>颜色通常是随机的，仅用于区分不同的函数</li>
            <li>点击某个框可以放大该部分</li>
            <li>搜索功能可以高亮特定函数</li>
        </ul>
    </div>
</body>
</html>`, processName, targetPID, processName, targetPID, config.Duration, config.Frequency,
		time.Now().Format("2006-01-02 15:04:05"), filepath.Base(svgFile))

	if err := os.WriteFile(htmlFile, []byte(htmlContent), 0644); err != nil {
		return fmt.Errorf("写入HTML文件失败: %v", err)
	}

	fmt.Printf("查看器已生成: %s\n", htmlFile)

	// 自动打开火焰图
	if config.AutoOpen {
		fmt.Println("\n正在打开火焰图...")
		// 尝试不同的打开命令
		openCommands := [][]string{
			{"xdg-open", htmlFile}, // Linux
			{"open", htmlFile},     // macOS
			{"start", htmlFile},    // Windows
		}

		for _, cmd := range openCommands {
			if err := exec.Command(cmd[0], cmd[1:]...).Start(); err == nil {
				break
			}
		}
	}

	// 输出分析建议
	fmt.Printf("\n【火焰图分析提示】\n")
	fmt.Println(strings.Repeat("-", 60))
	fmt.Println("1. 查看火焰图: 打开 " + htmlFile)
	fmt.Println("2. 宽度越宽的函数，CPU占用时间越长")
	fmt.Println("3. 寻找平顶（plateau）区域，这些可能是性能瓶颈")
	fmt.Println("4. 关注递归调用和深度调用栈")
	fmt.Println("5. 使用浏览器的搜索功能查找特定函数")

	return nil
}

// 简单的栈折叠方法（备用）
func simpleFoldStacks(perfOutput string) []byte {
	stacks := make(map[string]int)
	lines := strings.Split(perfOutput, "\n")
	currentStack := []string{}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if len(currentStack) > 0 {
				// 反转栈（perf输出是从叶子到根）
				for i, j := 0, len(currentStack)-1; i < j; i, j = i+1, j-1 {
					currentStack[i], currentStack[j] = currentStack[j], currentStack[i]
				}
				stackStr := strings.Join(currentStack, ";")
				stacks[stackStr]++
				currentStack = []string{}
			}
		} else if strings.Contains(line, " ") {
			// 提取函数名
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				funcName := parts[1]
				// 清理函数名
				funcName = strings.TrimSuffix(funcName, "+0x")
				if idx := strings.Index(funcName, "+"); idx > 0 {
					funcName = funcName[:idx]
				}
				currentStack = append(currentStack, funcName)
			}
		}
	}

	// 构建输出
	var result strings.Builder
	for stack, count := range stacks {
		result.WriteString(fmt.Sprintf("%s %d\n", stack, count))
	}

	return []byte(result.String())
}

func main() {
	// 定义命令行参数
	var (
		monitorFlag  = flag.Bool("m", false, "实时监控模式")
		monitorFlag2 = flag.Bool("monitor", false, "实时监控模式")
		interval     = flag.Int("interval", 2, "监控模式刷新间隔(秒)")
		duration     = flag.Int("duration", 60, "监控模式持续时间(秒)")

		// 火焰图相关参数
		flameGraph     = flag.Bool("flame", false, "生成CPU火焰图")
		flamePID       = flag.Int("pid", 0, "指定进程PID(0=自动选择CPU最高)")
		flameDuration  = flag.Int("flame-duration", 30, "火焰图采样时长(秒)")
		flameFrequency = flag.Int("flame-freq", 99, "火焰图采样频率(Hz)")
		flameOutput    = flag.String("flame-output", "", "火焰图输出目录")
		flameAutoOpen  = flag.Bool("flame-open", false, "自动打开生成的火焰图")

		helpFlag  = flag.Bool("h", false, "显示帮助信息")
		helpFlag2 = flag.Bool("help", false, "显示帮助信息")
		version   = flag.Bool("version", false, "显示版本信息")
	)

	flag.Parse()

	// 显示版本信息
	if *version {
		fmt.Println("PerfSnap v1.1.1")
		fmt.Println("支持火焰图生成功能")
		fmt.Println("修复磁盘IO统计字段解析问题")
		os.Exit(0)
	}

	// 显示帮助信息
	if *helpFlag || *helpFlag2 {
		fmt.Println("PerfSnap v1.1.1 - Linux系统性能快照分析工具")
		fmt.Println("\n用法:")
		fmt.Println("  perfsnap [选项]")
		fmt.Println("\n基础选项:")
		fmt.Println("  无参数               生成一次性能快照报告")
		fmt.Println("  -m, --monitor       实时监控模式")
		fmt.Println("  -interval N         监控刷新间隔(秒), 默认2")
		fmt.Println("  -duration N         监控持续时间(秒), 默认60")
		fmt.Println("\n火焰图选项:")
		fmt.Println("  -flame              生成CPU火焰图")
		fmt.Println("  -pid N              指定进程PID(0=自动选择CPU最高), 默认0")
		fmt.Println("  -flame-duration N   火焰图采样时长(秒), 默认30")
		fmt.Println("  -flame-freq N       火焰图采样频率(Hz), 默认99")
		fmt.Println("  -flame-output DIR   火焰图输出目录")
		fmt.Println("  -flame-open         自动打开生成的火焰图")
		fmt.Println("\n其他选项:")
		fmt.Println("  -h, --help          显示帮助信息")
		fmt.Println("  -version            显示版本信息")
		fmt.Println("\n示例:")
		fmt.Println("  perfsnap                        # 生成性能快照")
		fmt.Println("  perfsnap -m                     # 实时监控(默认2秒间隔，60秒)")
		fmt.Println("  perfsnap -m -interval 5 -duration 120  # 每5秒刷新，持续120秒")
		fmt.Println("  perfsnap -flame                 # 生成CPU最高进程的火焰图")
		fmt.Println("  perfsnap -flame -pid 1234       # 生成指定进程的火焰图")
		fmt.Println("  perfsnap -flame -flame-duration 60 -flame-open  # 采样60秒并自动打开")
		fmt.Println("\n注意:")
		fmt.Println("  - 建议使用root权限运行以获取完整数据")
		fmt.Println("  - 火焰图功能需要安装perf工具")
		fmt.Println("  - 首次使用会自动下载FlameGraph工具")
		fmt.Println("\n项目: https://github.com/sunyifei83/devops-toolkit")
		return
	}

	// 检查是否以root权限运行
	if os.Geteuid() != 0 {
		fmt.Println("⚠️  建议以root权限运行以获取完整的性能数据")
		fmt.Println("   sudo perfsnap")
		fmt.Println()
	}

	// 实时监控模式
	if *monitorFlag || *monitorFlag2 {
		monitorMode(*interval, *duration)
		return
	}

	// 生成性能报告
	fmt.Println("正在收集系统性能数据，请稍候...")
	startTime := time.Now()
	data := collectPerformanceData()
	elapsed := time.Since(startTime)

	printPerformanceReport(data)
	fmt.Printf("\n数据收集耗时: %.2f秒\n", elapsed.Seconds())

	// 如果需要生成火焰图
	if *flameGraph {
		config := FlameGraphConfig{
			Enabled:   true,
			PID:       *flamePID,
			Duration:  *flameDuration,
			Frequency: *flameFrequency,
			OutputDir: *flameOutput,
			AutoOpen:  *flameAutoOpen,
		}

		fmt.Println("\n" + strings.Repeat("=", 80))
		fmt.Println("                    开始生成火焰图")
		fmt.Println(strings.Repeat("=", 80))

		if err := generateFlameGraph(config, data.TopProcs); err != nil {
			fmt.Printf("\n❌ 火焰图生成失败: %v\n", err)
			fmt.Println("\n可能的解决方案:")
			fmt.Println("1. 确保以root权限运行: sudo perfsnap -flame")
			fmt.Println("2. 安装perf工具:")
			fmt.Println("   Ubuntu/Debian: sudo apt-get install linux-tools-common linux-tools-generic")
			fmt.Println("   CentOS/RHEL: sudo yum install perf")
			fmt.Println("3. 检查目标进程是否存在")
		}
	}
}
