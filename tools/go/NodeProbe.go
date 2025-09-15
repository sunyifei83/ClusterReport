/*
NodeProbe - Linux服务器节点配置信息收集工具
全面采集服务器硬件配置、系统状态和软件环境信息，支持自动优化系统设置
Author: sunyifei83@gmail.com
Version: 1.1.0
项目: https://github.com/sunyifei83/devops-toolkit

Features:
- 支持多格式输出：表格(默认)、JSON、YAML
- 自动优化系统设置
- 完美支持中文字符对齐

TODO:
1. 支持远程节点信息采集
2. 支持节点硬件基准测试（CPU/内存/磁盘/网络）
*/
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ServerInfo 存储服务器信息
type ServerInfo struct {
	Hostname      string             `json:"hostname" yaml:"hostname"`
	LoadAverage   string             `json:"load_average" yaml:"load_average"`
	Timezone      string             `json:"timezone" yaml:"timezone"`
	OS            string             `json:"os" yaml:"os"`
	Kernel        string             `json:"kernel" yaml:"kernel"`
	CPU           CPUInfo            `json:"cpu" yaml:"cpu"`
	Memory        MemoryInfo         `json:"memory" yaml:"memory"`
	Disks         DiskInfo           `json:"disks" yaml:"disks"`
	Network       []NetworkInterface `json:"network" yaml:"network"`
	Python        PythonInfo         `json:"python" yaml:"python"`
	Java          JavaInfo           `json:"java" yaml:"java"`
	KernelModules KernelModuleStatus `json:"kernel_modules" yaml:"kernel_modules"`
	Timestamp     string             `json:"timestamp" yaml:"timestamp"`
	Version       string             `json:"nodeprobe_version" yaml:"nodeprobe_version"`
}

type CPUInfo struct {
	Model           string `json:"model" yaml:"model"`
	Cores           int    `json:"cores" yaml:"cores"`
	RunMode         string `json:"run_mode" yaml:"run_mode"`
	PerformanceMode string `json:"performance_mode" yaml:"performance_mode"`
}

type MemoryInfo struct {
	TotalGB float64  `json:"total_gb" yaml:"total_gb"`
	Slots   []string `json:"slots" yaml:"slots"`
}

type DiskInfo struct {
	SystemDisk  string   `json:"system_disk" yaml:"system_disk"`
	DataDisks   []string `json:"data_disks" yaml:"data_disks"`
	TotalDisks  int      `json:"total_disks" yaml:"total_disks"`
	DataDiskNum int      `json:"data_disk_num" yaml:"data_disk_num"`
}

type NetworkInterface struct {
	Name   string `json:"name" yaml:"name"`
	Status string `json:"status" yaml:"status"`
	Speed  string `json:"speed" yaml:"speed"`
	IP     string `json:"ip" yaml:"ip"`
}

type PythonInfo struct {
	Version string `json:"version" yaml:"version"`
	Path    string `json:"path" yaml:"path"`
}

type JavaInfo struct {
	Version string `json:"version" yaml:"version"`
	Path    string `json:"path" yaml:"path"`
}

type KernelModuleStatus struct {
	NfConntrack bool   `json:"nf_conntrack" yaml:"nf_conntrack"`
	BrNetfilter bool   `json:"br_netfilter" yaml:"br_netfilter"`
	Message     string `json:"message" yaml:"message"`
}

// 执行命令并返回输出
func execCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// 获取主机名
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// 获取系统负载
func getLoadAverage() string {
	data, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return "N/A"
	}
	fields := strings.Fields(string(data))
	if len(fields) >= 3 {
		return fmt.Sprintf("%s %s %s", fields[0], fields[1], fields[2])
	}
	return "N/A"
}

// 获取操作系统信息
func getOSInfo() (string, string) {
	osInfo := "Unknown"
	kernel := "Unknown"

	// 获取发行版信息
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "PRETTY_NAME=") {
				osInfo = strings.Trim(strings.TrimPrefix(line, "PRETTY_NAME="), "\"")
				break
			}
		}
	}

	// 获取内核版本
	if data, err := os.ReadFile("/proc/version"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			kernel = fields[2]
		}
	}

	return osInfo, kernel
}

// 获取CPU信息
func getCPUInfo() CPUInfo {
	info := CPUInfo{}

	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return info
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	cores := 0
	modelFound := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "model name") && !modelFound {
			parts := strings.Split(line, ":")
			if len(parts) == 2 {
				info.Model = strings.TrimSpace(parts[1])
				modelFound = true
			}
		}
		if strings.Contains(line, "processor") {
			cores++
		}
	}

	info.Cores = cores

	// 获取CPU运行模式
	info.RunMode = getCPURunMode()

	// 获取CPU性能模式
	info.PerformanceMode = getCPUPerformanceMode()

	return info
}

// 获取CPU运行模式
func getCPURunMode() string {
	// 尝试使用lscpu命令获取CPU运行模式
	if output, err := execCommand("lscpu"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "CPU op-mode(s)") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					return strings.TrimSpace(parts[1])
				}
			}
		}
	}

	// 如果lscpu失败，尝试从/proc/cpuinfo获取架构信息
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		if strings.Contains(string(data), "lm") {
			return "32-bit, 64-bit"
		}
		return "32-bit"
	}

	return "Unknown"
}

// 获取CPU性能模式
func getCPUPerformanceMode() string {
	// 检查CPU调度器模式
	governorFiles := []string{
		"/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor",
		"/sys/devices/system/cpu/cpufreq/policy0/scaling_governor",
	}

	var currentGovernor string

	for _, file := range governorFiles {
		if data, err := os.ReadFile(file); err == nil {
			currentGovernor = strings.TrimSpace(string(data))
			break
		}
	}

	// 如果当前是省电模式且以root权限运行，自动切换到性能模式
	if currentGovernor == "powersave" && os.Geteuid() == 0 {
		// 获取所有CPU的governor文件路径
		cpuGovernorPattern := "/sys/devices/system/cpu/cpu*/cpufreq/scaling_governor"
		files, err := filepath.Glob(cpuGovernorPattern)
		if err == nil && len(files) > 0 {
			successCount := 0
			for _, file := range files {
				if err := os.WriteFile(file, []byte("performance\n"), 0644); err == nil {
					successCount++
				}
			}
			if successCount > 0 {
				return fmt.Sprintf("已自动调整至最大性能模式 (原: powersave, 成功调整 %d 个CPU)", successCount)
			}
		}
		// 如果glob失败，尝试使用shell命令
		if _, err := execCommand("sh", "-c", `for cpu in /sys/devices/system/cpu/cpu*/cpufreq/scaling_governor; do echo "performance" > $cpu; done`); err == nil {
			return "已自动调整至最大性能模式 (原: powersave)"
		}
		return "省电模式 (powersave) - 自动调整失败"
	}

	// 返回当前模式状态
	if currentGovernor != "" {
		switch currentGovernor {
		case "performance":
			return "最大性能模式 (performance)"
		case "powersave":
			if os.Geteuid() != 0 {
				return "省电模式 (powersave) - 需要root权限才能自动调整"
			}
			return "省电模式 (powersave)"
		case "ondemand":
			return "按需模式 (ondemand)"
		case "conservative":
			return "保守模式 (conservative)"
		case "schedutil":
			return "调度器模式 (schedutil)"
		default:
			return fmt.Sprintf("当前模式: %s", currentGovernor)
		}
	}

	// 尝试使用cpupower命令
	if output, err := execCommand("cpupower", "frequency-info", "--policy"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "governor") {
				parts := strings.Fields(line)
				if len(parts) > 0 {
					governor := parts[len(parts)-1]
					return fmt.Sprintf("当前模式: %s", strings.Trim(governor, "\""))
				}
			}
		}
	}

	return "未知 (可能不支持频率调节)"
}

// 获取内存信息
func getMemoryInfo() MemoryInfo {
	info := MemoryInfo{}

	// 获取总内存
	if data, err := os.ReadFile("/proc/meminfo"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					if kb, err := strconv.ParseFloat(fields[1], 64); err == nil {
						info.TotalGB = kb / 1024 / 1024
					}
				}
				break
			}
		}
	}

	// 获取内存插槽信息
	if output, err := execCommand("dmidecode", "-t", "17"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "Size:") && !strings.Contains(line, "No Module") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					info.Slots = append(info.Slots, strings.Join(parts[1:], " "))
				}
			}
		}
	}

	return info
}

// 获取磁盘信息
func getDiskInfo() DiskInfo {
	info := DiskInfo{}

	// 获取系统盘挂载信息
	if output, err := execCommand("df", "-h"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, " / ") && !strings.Contains(line, "udev") &&
				!strings.Contains(line, "tmpfs") && !strings.Contains(line, "/boot") {
				fields := strings.Fields(line)
				if len(fields) >= 6 {
					info.SystemDisk = fmt.Sprintf("%s %s/%s (%s)",
						fields[0], fields[3], fields[2], fields[5])
				}
			}
		}
	}

	// 获取所有磁盘信息
	if output, err := execCommand("lsblk", "-d", "-o", "NAME,SIZE,TYPE"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "disk") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					diskInfo := fmt.Sprintf("/dev/%s %s", fields[0], fields[1])
					// 判断是否为数据盘（默认大于1T的认为是数据盘）
					if strings.Contains(fields[1], "T") ||
						(strings.Contains(fields[1], "G") && parseSize(fields[1]) > 1000) {
						info.DataDisks = append(info.DataDisks, diskInfo)
						info.DataDiskNum++
					}
					info.TotalDisks++
				}
			}
		}
	}

	return info
}

// 解析磁盘大小
func parseSize(sizeStr string) float64 {
	re := regexp.MustCompile(`(\d+\.?\d*)([KMGT])`)
	matches := re.FindStringSubmatch(sizeStr)
	if len(matches) == 3 {
		size, _ := strconv.ParseFloat(matches[1], 64)
		unit := matches[2]
		switch unit {
		case "T":
			return size * 1024
		case "G":
			return size
		case "M":
			return size / 1024
		case "K":
			return size / 1024 / 1024
		}
	}
	return 0
}

// 获取网络接口信息
func getNetworkInfo() []NetworkInterface {
	var interfaces []NetworkInterface

	if output, err := execCommand("ip", "addr", "show"); err == nil {
		lines := strings.Split(output, "\n")
		var currentIface NetworkInterface

		for _, line := range lines {
			// 解析接口名称和状态
			if matches := regexp.MustCompile(`^\d+:\s+(\S+):\s+<(.+?)>`).FindStringSubmatch(line); len(matches) > 2 {
				if currentIface.Name != "" {
					interfaces = append(interfaces, currentIface)
				}
				currentIface = NetworkInterface{Name: matches[1]}
				if strings.Contains(matches[2], "UP") {
					currentIface.Status = "UP"
				} else {
					currentIface.Status = "DOWN"
				}
			}

			// 解析IP地址
			if strings.Contains(line, "inet ") && !strings.Contains(line, "inet6") {
				fields := strings.Fields(line)
				for i, field := range fields {
					if field == "inet" && i+1 < len(fields) {
						currentIface.IP = fields[i+1]
						break
					}
				}
			}

			// 获取速度信息
			if currentIface.Name != "" && currentIface.Name != "lo" {
				if speedOutput, err := execCommand("ethtool", currentIface.Name); err == nil {
					if matches := regexp.MustCompile(`Speed:\s+(\S+)`).FindStringSubmatch(speedOutput); len(matches) > 1 {
						currentIface.Speed = matches[1]
					}
				}
			}
		}

		if currentIface.Name != "" {
			interfaces = append(interfaces, currentIface)
		}
	}

	// 过滤掉lo接口
	var filtered []NetworkInterface
	for _, iface := range interfaces {
		if iface.Name != "lo" {
			filtered = append(filtered, iface)
		}
	}

	return filtered
}

// 获取时区信息并校准
func getTimezone() string {
	var currentTZ string

	// 尝试使用timedatectl命令
	if output, err := execCommand("timedatectl", "show", "--property=Timezone", "--value"); err == nil {
		currentTZ = strings.TrimSpace(output)
		if currentTZ != "" {
			// 如果当前时区不是Asia/Shanghai，尝试设置
			if currentTZ != "Asia/Shanghai" && os.Geteuid() == 0 {
				if _, err := execCommand("timedatectl", "set-timezone", "Asia/Shanghai"); err == nil {
					return fmt.Sprintf("已校准至 Asia/Shanghai (原: %s)", currentTZ)
				}
			}
			return currentTZ
		}
	}

	// 尝试读取/etc/timezone文件
	if data, err := os.ReadFile("/etc/timezone"); err == nil {
		currentTZ = strings.TrimSpace(string(data))
		if currentTZ != "" {
			// 如果当前时区不是Asia/Shanghai，尝试设置
			if currentTZ != "Asia/Shanghai" && os.Geteuid() == 0 {
				// 尝试写入新时区
				if err := os.WriteFile("/etc/timezone", []byte("Asia/Shanghai\n"), 0644); err == nil {
					// 更新/etc/localtime
					if _, err := execCommand("ln", "-sf", "/usr/share/zoneinfo/Asia/Shanghai", "/etc/localtime"); err == nil {
						return fmt.Sprintf("已校准至 Asia/Shanghai (原: %s)", currentTZ)
					}
				}
			}
			return currentTZ
		}
	}

	// 尝试读取/etc/localtime符号链接
	if link, err := os.Readlink("/etc/localtime"); err == nil {
		// 通常格式为 /usr/share/zoneinfo/Asia/Shanghai
		parts := strings.Split(link, "zoneinfo/")
		if len(parts) == 2 {
			currentTZ = parts[1]
			if currentTZ != "Asia/Shanghai" && os.Geteuid() == 0 {
				if _, err := execCommand("ln", "-sf", "/usr/share/zoneinfo/Asia/Shanghai", "/etc/localtime"); err == nil {
					return fmt.Sprintf("已校准至 Asia/Shanghai (原: %s)", currentTZ)
				}
			}
			return currentTZ
		}
	}

	// 尝试使用date命令
	if output, err := execCommand("date", "+%Z"); err == nil {
		return strings.TrimSpace(output)
	}

	return "Unknown"
}

// 获取Python信息
func getPythonInfo() PythonInfo {
	info := PythonInfo{}

	// 尝试多个可能的Python命令
	pythonCommands := []string{"python3", "python", "python2"}

	for _, cmd := range pythonCommands {
		// 获取Python版本
		if output, err := execCommand(cmd, "--version"); err == nil {
			// Python版本信息可能在stdout或stderr中
			versionStr := strings.TrimSpace(output)
			if versionStr != "" {
				info.Version = versionStr

				// 获取Python路径
				if pathOutput, err := execCommand("which", cmd); err == nil {
					info.Path = strings.TrimSpace(pathOutput)
				}
				break
			}
		}
	}

	// 如果没有找到Python
	if info.Version == "" {
		info.Version = "Not installed"
		info.Path = "N/A"
	}

	return info
}

// 检查和加载内核模块
func checkAndLoadKernelModules() KernelModuleStatus {
	status := KernelModuleStatus{}

	modules := []struct {
		name  string
		field *bool
	}{
		{"nf_conntrack", &status.NfConntrack},
		{"br_netfilter", &status.BrNetfilter},
	}

	var messages []string

	for _, mod := range modules {
		// 检查模块是否已加载
		if output, err := execCommand("lsmod"); err == nil {
			if strings.Contains(output, mod.name) {
				*mod.field = true
				messages = append(messages, fmt.Sprintf("%s: 已加载", mod.name))
				continue
			}
		}

		// 如果未加载且是root用户，尝试加载
		if os.Geteuid() == 0 {
			if _, err := execCommand("modprobe", mod.name); err == nil {
				*mod.field = true
				messages = append(messages, fmt.Sprintf("%s: 已自动加载", mod.name))
			} else {
				*mod.field = false
				messages = append(messages, fmt.Sprintf("%s: 加载失败", mod.name))
			}
		} else {
			*mod.field = false
			messages = append(messages, fmt.Sprintf("%s: 未加载(需要root权限)", mod.name))
		}
	}

	status.Message = strings.Join(messages, ", ")
	return status
}

// 获取Java信息
func getJavaInfo() JavaInfo {
	info := JavaInfo{}

	// 获取Java版本
	if output, err := execCommand("java", "-version"); err == nil {
		// Java版本信息通常在stderr中，但execCommand会捕获两者
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "version") {
				// 提取版本号
				if matches := regexp.MustCompile(`"([^"]+)"`).FindStringSubmatch(line); len(matches) > 1 {
					info.Version = "Java " + matches[1]
				} else {
					info.Version = strings.TrimSpace(line)
				}
				break
			}
		}

		// 获取Java路径
		if pathOutput, err := execCommand("which", "java"); err == nil {
			info.Path = strings.TrimSpace(pathOutput)

			// 尝试获取JAVA_HOME
			if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
				info.Path = fmt.Sprintf("%s (JAVA_HOME: %s)", info.Path, javaHome)
			}
		}
	}

	// 如果没有找到Java
	if info.Version == "" {
		info.Version = "Not installed"
		info.Path = "N/A"
	}

	return info
}

// 填充空格到指定宽度
func padRight(s string, width int) string {
	sWidth := displayWidth(s)
	if sWidth >= width {
		return truncateByWidth(s, width)
	}
	return s + strings.Repeat(" ", width-sWidth)
}

// 打印服务器信息
func printServerInfo(info ServerInfo) {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")

	// 基础信息
	fmt.Printf("║ %-20s %-43s ║\n", "主机名:", info.Hostname)
	fmt.Printf("║ %-20s %-43s ║\n", "系统负载:", info.LoadAverage)
	fmt.Printf("║ %-20s %-43s ║\n", "时区:", info.Timezone)
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 系统信息
	fmt.Printf("║ %s %s ║\n", padRight("操作系统:", 20), padRight(info.OS, 43))
	fmt.Printf("║ %s %s ║\n", padRight("内核版本:", 20), padRight(info.Kernel, 43))
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// CPU信息
	fmt.Printf("║ %s %s ║\n", padRight("CPU型号:", 20), padRight(info.CPU.Model, 43))
	fmt.Printf("║ %s %s ║\n", padRight("CPU核心数:", 20), padRight(fmt.Sprintf("%d", info.CPU.Cores), 43))
	fmt.Printf("║ %s %s ║\n", padRight("CPU运行模式:", 20), padRight(info.CPU.RunMode, 43))
	fmt.Printf("║ %s %s ║\n", padRight("CPU性能模式:", 20), padRight(info.CPU.PerformanceMode, 43))
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 内存信息
	memStr := fmt.Sprintf("%.2f GB", info.Memory.TotalGB)
	fmt.Printf("║ %s %s ║\n", padRight("总内存:", 20), padRight(memStr, 43))
	if len(info.Memory.Slots) > 0 {
		slotInfo := fmt.Sprintf("%d个插槽已使用", len(info.Memory.Slots))
		fmt.Printf("║ %s %s ║\n", padRight("内存插槽:", 20), padRight(slotInfo, 43))
		for i, slot := range info.Memory.Slots {
			slotLabel := fmt.Sprintf("  插槽%d:", i+1)
			fmt.Printf("║ %s %s ║\n", padRight(slotLabel, 20), padRight(slot, 43))
		}
	}
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 磁盘信息
	fmt.Printf("║ %s %s ║\n", padRight("系统盘:", 20), padRight(info.Disks.SystemDisk, 43))
	diskCount := fmt.Sprintf("总计: %d, 数据盘: %d", info.Disks.TotalDisks, info.Disks.DataDiskNum)
	fmt.Printf("║ %s %s ║\n", padRight("磁盘数量:", 20), padRight(diskCount, 43))

	if len(info.Disks.DataDisks) > 0 {
		for i, disk := range info.Disks.DataDisks {
			if i == 0 {
				fmt.Printf("║ %s %s ║\n", padRight("数据盘:", 20), padRight(disk, 43))
			} else {
				fmt.Printf("║ %s %s ║\n", padRight("", 20), padRight(disk, 43))
			}
		}
	}
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 网络信息
	fmt.Printf("║ %s %s ║\n", padRight("网络接口数:", 20), padRight(fmt.Sprintf("%d", len(info.Network)), 43))
	for _, iface := range info.Network {
		netInfo := fmt.Sprintf("%s [%s] %s %s", iface.Name, iface.Status, iface.Speed, iface.IP)
		fmt.Printf("║   %s ║\n", padRight(netInfo, 61))
	}
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// Python信息
	fmt.Printf("║ %s %s ║\n", padRight("Python版本:", 20), padRight(info.Python.Version, 43))
	fmt.Printf("║ %s %s ║\n", padRight("Python路径:", 20), padRight(info.Python.Path, 43))
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// Java信息
	fmt.Printf("║ %s %s ║\n", padRight("Java版本:", 20), padRight(info.Java.Version, 43))
	fmt.Printf("║ %s %s ║\n", padRight("Java路径:", 20), padRight(info.Java.Path, 43))
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 内核模块信息
	messages := strings.Split(info.KernelModules.Message, ", ")
	for i, msg := range messages {
		if i == 0 {
			fmt.Printf("║ %s %s ║\n", padRight("内核模块状态:", 20), padRight(msg, 43))
		} else {
			fmt.Printf("║ %s %s ║\n", padRight("", 20), padRight(msg, 43))
		}
	}

	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
}

// 截断字符串（正确处理UTF-8字符）
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen-3]) + "..."
}

// 计算字符串显示宽度（考虑中文字符）
func displayWidth(s string) int {
	width := 0
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			// 中文字符占2个宽度
			width += 2
		} else {
			width += 1
		}
	}
	return width
}

// 按显示宽度截断字符串
func truncateByWidth(s string, maxWidth int) string {
	if displayWidth(s) <= maxWidth {
		return s
	}

	runes := []rune(s)
	width := 0
	cutIndex := 0

	for i, r := range runes {
		if r >= 0x4E00 && r <= 0x9FFF {
			width += 2
		} else {
			width += 1
		}

		if width > maxWidth-3 {
			cutIndex = i
			break
		}
	}

	if cutIndex == 0 {
		cutIndex = len(runes)
	}

	return string(runes[:cutIndex]) + "..."
}

// outputJSON 输出JSON格式
func outputJSON(info ServerInfo) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(info)
}

// outputYAML 输出YAML格式
func outputYAML(info ServerInfo) error {
	// 手动构建YAML输出
	var output strings.Builder

	output.WriteString("# NodeProbe Configuration Report\n")
	output.WriteString(fmt.Sprintf("# Generated at: %s\n\n", info.Timestamp))

	output.WriteString(fmt.Sprintf("hostname: %s\n", info.Hostname))
	output.WriteString(fmt.Sprintf("load_average: %s\n", info.LoadAverage))
	output.WriteString(fmt.Sprintf("timezone: %s\n", info.Timezone))
	output.WriteString(fmt.Sprintf("os: %s\n", info.OS))
	output.WriteString(fmt.Sprintf("kernel: %s\n", info.Kernel))
	output.WriteString(fmt.Sprintf("timestamp: %s\n", info.Timestamp))
	output.WriteString(fmt.Sprintf("nodeprobe_version: %s\n\n", info.Version))

	// CPU信息
	output.WriteString("cpu:\n")
	output.WriteString(fmt.Sprintf("  model: %s\n", info.CPU.Model))
	output.WriteString(fmt.Sprintf("  cores: %d\n", info.CPU.Cores))
	output.WriteString(fmt.Sprintf("  run_mode: %s\n", info.CPU.RunMode))
	output.WriteString(fmt.Sprintf("  performance_mode: %s\n\n", info.CPU.PerformanceMode))

	// 内存信息
	output.WriteString("memory:\n")
	output.WriteString(fmt.Sprintf("  total_gb: %.2f\n", info.Memory.TotalGB))
	if len(info.Memory.Slots) > 0 {
		output.WriteString("  slots:\n")
		for _, slot := range info.Memory.Slots {
			output.WriteString(fmt.Sprintf("    - %s\n", slot))
		}
	}
	output.WriteString("\n")

	// 磁盘信息
	output.WriteString("disks:\n")
	output.WriteString(fmt.Sprintf("  system_disk: %s\n", info.Disks.SystemDisk))
	output.WriteString(fmt.Sprintf("  total_disks: %d\n", info.Disks.TotalDisks))
	output.WriteString(fmt.Sprintf("  data_disk_num: %d\n", info.Disks.DataDiskNum))
	if len(info.Disks.DataDisks) > 0 {
		output.WriteString("  data_disks:\n")
		for _, disk := range info.Disks.DataDisks {
			output.WriteString(fmt.Sprintf("    - %s\n", disk))
		}
	}
	output.WriteString("\n")

	// 网络信息
	output.WriteString("network:\n")
	for _, iface := range info.Network {
		output.WriteString(fmt.Sprintf("  - name: %s\n", iface.Name))
		output.WriteString(fmt.Sprintf("    status: %s\n", iface.Status))
		output.WriteString(fmt.Sprintf("    speed: %s\n", iface.Speed))
		output.WriteString(fmt.Sprintf("    ip: %s\n", iface.IP))
	}
	output.WriteString("\n")

	// Python信息
	output.WriteString("python:\n")
	output.WriteString(fmt.Sprintf("  version: %s\n", info.Python.Version))
	output.WriteString(fmt.Sprintf("  path: %s\n\n", info.Python.Path))

	// Java信息
	output.WriteString("java:\n")
	output.WriteString(fmt.Sprintf("  version: %s\n", info.Java.Version))
	output.WriteString(fmt.Sprintf("  path: %s\n\n", info.Java.Path))

	// 内核模块信息
	output.WriteString("kernel_modules:\n")
	output.WriteString(fmt.Sprintf("  nf_conntrack: %v\n", info.KernelModules.NfConntrack))
	output.WriteString(fmt.Sprintf("  br_netfilter: %v\n", info.KernelModules.BrNetfilter))
	output.WriteString(fmt.Sprintf("  message: %s\n", info.KernelModules.Message))

	fmt.Print(output.String())
	return nil
}

func main() {
	// 定义命令行参数
	var (
		outputFormat = flag.String("format", "table", "输出格式: table(默认), json, yaml")
		showVersion  = flag.Bool("version", false, "显示版本信息")
		quiet        = flag.Bool("quiet", false, "静默模式，减少提示信息")
		outputFile   = flag.String("output", "", "输出到文件")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "NodeProbe v1.1.0 - Linux节点配置探测工具\n\n")
		fmt.Fprintf(os.Stderr, "用法: %s [选项]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "选项:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  %s                    # 默认表格输出\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -format json       # JSON格式输出\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -format yaml       # YAML格式输出\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -output report.json -format json  # 输出到文件\n", os.Args[0])
	}

	flag.Parse()

	// 显示版本信息
	if *showVersion {
		fmt.Println("NodeProbe v1.1.0")
		os.Exit(0)
	}

	// 验证输出格式
	validFormats := map[string]bool{"table": true, "json": true, "yaml": true}
	if !validFormats[*outputFormat] {
		fmt.Fprintf(os.Stderr, "错误: 不支持的输出格式 '%s'\n", *outputFormat)
		fmt.Fprintf(os.Stderr, "支持的格式: table, json, yaml\n")
		os.Exit(1)
	}

	// 如果指定了输出文件，重定向stdout
	var originalStdout *os.File
	if *outputFile != "" {
		originalStdout = os.Stdout
		file, err := os.Create(*outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "错误: 无法创建输出文件 %s: %v\n", *outputFile, err)
			os.Exit(1)
		}
		defer file.Close()
		os.Stdout = file
	}

	// 非静默模式且表格输出时显示提示信息
	if !*quiet && *outputFormat == "table" && *outputFile == "" {
		fmt.Println("NodeProbe v1.1.0 - Linux节点配置探测工具")
		fmt.Println("=" + strings.Repeat("=", 65))

		// 检查是否以root权限运行
		if os.Geteuid() != 0 {
			fmt.Println("⚠️  某些硬件信息需要root权限才能获取完整数据")
			fmt.Println("建议使用: sudo nodeprobe")
			fmt.Println()
		}

		fmt.Println("正在探测节点配置信息...")
	}

	// 收集服务器信息
	info := ServerInfo{
		Hostname:      getHostname(),
		LoadAverage:   getLoadAverage(),
		Timezone:      getTimezone(),
		CPU:           getCPUInfo(),
		Memory:        getMemoryInfo(),
		Disks:         getDiskInfo(),
		Network:       getNetworkInfo(),
		Python:        getPythonInfo(),
		Java:          getJavaInfo(),
		KernelModules: checkAndLoadKernelModules(),
		Timestamp:     time.Now().Format("2006-01-02 15:04:05"),
		Version:       "1.1.0",
	}

	info.OS, info.Kernel = getOSInfo()

	// 根据输出格式选择输出方式
	switch *outputFormat {
	case "json":
		if err := outputJSON(info); err != nil {
			fmt.Fprintf(os.Stderr, "错误: 输出JSON失败: %v\n", err)
			os.Exit(1)
		}
	case "yaml":
		if err := outputYAML(info); err != nil {
			fmt.Fprintf(os.Stderr, "错误: 输出YAML失败: %v\n", err)
			os.Exit(1)
		}
	default: // table
		if !*quiet && *outputFormat == "table" {
			fmt.Print("\033[2J\033[H") // 清屏
		}
		printServerInfo(info)
		if !*quiet {
			fmt.Println("\n由 NodeProbe 生成 | https://github.com/sunyifei83/devops-toolkit")
		}
	}

	// 如果输出到文件，显示成功消息
	if *outputFile != "" && originalStdout != nil {
		os.Stdout = originalStdout
		fmt.Printf("✅ 结果已保存到: %s\n", *outputFile)
	}
}
