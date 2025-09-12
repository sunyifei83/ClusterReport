/*
Linux服务器性能快速分析工具
采集最近1分钟的系统性能数据，包括CPU、内存、磁盘、网络等指标
Author: sunyifei83@gmail.com
Version: 1.0.0
Tips: LinuxHostLoad工具需要系统安装sysstat包（提供sar、mpstat、pidstat、iostat等命令）
使用:
1.go build
2.chmod 700 LinuxHostLoad_v1  执行 ./LinuxHostLoad_v1
*/
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
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
				if len(fields) >= 14 {
					dev := DeviceIO{
						Device: fields[0],
					}
					dev.RRPS, _ = strconv.ParseFloat(fields[1], 64)
					dev.WRPS, _ = strconv.ParseFloat(fields[2], 64)
					dev.RkBPS, _ = strconv.ParseFloat(fields[3], 64)
					dev.WkBPS, _ = strconv.ParseFloat(fields[4], 64)
					dev.AvgQueue, _ = strconv.ParseFloat(fields[8], 64)
					dev.AvgWait, _ = strconv.ParseFloat(fields[9], 64)
					dev.SvcTime, _ = strconv.ParseFloat(fields[12], 64)
					dev.Util, _ = strconv.ParseFloat(fields[13], 64)

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
	fmt.Println("                    Linux 系统性能快速分析报告")
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

func main() {
	// 检查是否以root权限运行（某些命令需要root权限）
	if os.Geteuid() != 0 {
		fmt.Println("⚠️  建议以root权限运行以获取完整的性能数据")
		fmt.Println("   sudo " + os.Args[0])
		fmt.Println()
	}

	// 检查命令行参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "-m", "--monitor":
			// 实时监控模式
			interval := 2  // 默认2秒间隔
			duration := 60 // 默认监控60秒

			if len(os.Args) > 2 {
				if i, err := strconv.Atoi(os.Args[2]); err == nil {
					interval = i
				}
			}
			if len(os.Args) > 3 {
				if d, err := strconv.Atoi(os.Args[3]); err == nil {
					duration = d
				}
			}

			monitorMode(interval, duration)
			return

		case "-h", "--help":
			fmt.Println("Linux系统性能快速分析工具")
			fmt.Println("\n用法:")
			fmt.Println("  " + os.Args[0] + "              - 生成一次性能分析报告")
			fmt.Println("  " + os.Args[0] + " -m [间隔] [持续时间] - 实时监控模式")
			fmt.Println("  " + os.Args[0] + " -h          - 显示帮助信息")
			fmt.Println("\n示例:")
			fmt.Println("  " + os.Args[0] + "              - 生成性能报告")
			fmt.Println("  " + os.Args[0] + " -m           - 实时监控(默认2秒间隔，60秒)")
			fmt.Println("  " + os.Args[0] + " -m 5 120     - 每5秒刷新，持续120秒")
			return
		}
	}

	// 默认模式：生成一次性报告
	fmt.Println("正在收集系统性能数据，请稍候...")

	startTime := time.Now()
	data := collectPerformanceData()
	elapsed := time.Since(startTime)

	printPerformanceReport(data)
	fmt.Printf("\n数据收集耗时: %.2f秒\n", elapsed.Seconds())
}
