/*
v1版本暂时只支持本机执行采集cpu内存等基础信息，不支持其他节点远程执行
TODO: 1.支持其他节点远程执行 2.支持节点硬件基准测试并输出结果(cpu/内存/磁盘/网络)
Author: sunyifei@qiniu.com
Version: 1.0.0
*/
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// ServerInfo 存储服务器信息
type ServerInfo struct {
	Hostname      string
	LoadAverage   string
	Timezone      string
	OS            string
	Kernel        string
	CPU           CPUInfo
	Memory        MemoryInfo
	Disks         DiskInfo
	Network       []NetworkInterface
	Python        PythonInfo
	Java          JavaInfo
	KernelModules KernelModuleStatus
}

type CPUInfo struct {
	Model           string
	Cores           int
	RunMode         string
	PerformanceMode string
}

type MemoryInfo struct {
	TotalGB float64
	Slots   []string
}

type DiskInfo struct {
	SystemDisk  string
	DataDisks   []string
	TotalDisks  int
	DataDiskNum int
}

type NetworkInterface struct {
	Name   string
	Status string
	Speed  string
	IP     string
}

type PythonInfo struct {
	Version string
	Path    string
}

type JavaInfo struct {
	Version string
	Path    string
}

type KernelModuleStatus struct {
	NfConntrack bool
	BrNetfilter bool
	Message     string
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

// 打印服务器信息
func printServerInfo(info ServerInfo) {
	fmt.Println("╔════════════════════════════════════════════════════════════════╗")
	fmt.Printf("║ %-20s %-43s ║\n", "主机名:", info.Hostname)
	fmt.Printf("║ %-20s %-43s ║\n", "系统负载:", info.LoadAverage)
	fmt.Printf("║ %-20s %-43s ║\n", "时区:", info.Timezone)
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 系统信息
	fmt.Printf("║ %-20s %-43s ║\n", "操作系统:", truncateString(info.OS, 43))
	fmt.Printf("║ %-20s %-43s ║\n", "内核版本:", info.Kernel)
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// CPU信息
	fmt.Printf("║ %-20s %-43s ║\n", "CPU型号:", truncateString(info.CPU.Model, 43))
	fmt.Printf("║ %-20s %-43d ║\n", "CPU核心数:", info.CPU.Cores)
	fmt.Printf("║ %-20s %-43s ║\n", "CPU运行模式:", info.CPU.RunMode)
	fmt.Printf("║ %-20s %-43s ║\n", "CPU性能模式:", truncateString(info.CPU.PerformanceMode, 43))
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 内存信息
	fmt.Printf("║ %-20s %.2f GB %36s ║\n", "总内存:", info.Memory.TotalGB, "")
	if len(info.Memory.Slots) > 0 {
		fmt.Printf("║ %-20s %-43s ║\n", "内存插槽:", fmt.Sprintf("%d个插槽已使用", len(info.Memory.Slots)))
		for i, slot := range info.Memory.Slots {
			fmt.Printf("║   插槽%d: %-55s ║\n", i+1, slot)
		}
	}
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 磁盘信息
	fmt.Printf("║ %-20s %-43s ║\n", "系统盘:", truncateString(info.Disks.SystemDisk, 43))
	fmt.Printf("║ %-20s 总计: %d, 数据盘: %d %22s ║\n", "磁盘数量:",
		info.Disks.TotalDisks, info.Disks.DataDiskNum, "")

	if len(info.Disks.DataDisks) > 0 {
		for i, disk := range info.Disks.DataDisks {
			if i == 0 {
				fmt.Printf("║ %-20s %-43s ║\n", "数据盘:", truncateString(disk, 43))
			} else {
				fmt.Printf("║ %-20s %-43s ║\n", "", truncateString(disk, 43))
			}
		}
	}
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 网络信息
	fmt.Printf("║ %-20s %-43d ║\n", "网络接口数:", len(info.Network))
	for _, iface := range info.Network {
		netInfo := fmt.Sprintf("%s [%s] %s %s", iface.Name, iface.Status, iface.Speed, iface.IP)
		fmt.Printf("║   %-61s ║\n", truncateString(netInfo, 61))
	}
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// Python信息
	fmt.Printf("║ %-20s %-43s ║\n", "Python版本:", truncateString(info.Python.Version, 43))
	fmt.Printf("║ %-20s %-43s ║\n", "Python路径:", truncateString(info.Python.Path, 43))
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// Java信息
	fmt.Printf("║ %-20s %-43s ║\n", "Java版本:", truncateString(info.Java.Version, 43))
	fmt.Printf("║ %-20s %-43s ║\n", "Java路径:", truncateString(info.Java.Path, 43))
	fmt.Println("╠════════════════════════════════════════════════════════════════╣")

	// 内核模块信息
	fmt.Printf("║ %-20s %-43s ║\n", "内核模块状态:", truncateString(info.KernelModules.Message, 43))
	// 如果消息太长，换行显示
	if len(info.KernelModules.Message) > 43 {
		remaining := info.KernelModules.Message[43:]
		for len(remaining) > 0 {
			if len(remaining) > 43 {
				fmt.Printf("║ %-20s %-43s ║\n", "", truncateString(remaining[:43], 43))
				remaining = remaining[43:]
			} else {
				fmt.Printf("║ %-20s %-43s ║\n", "", remaining)
				break
			}
		}
	}

	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
}

// 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func main() {
	// 检查是否以root权限运行
	if os.Geteuid() != 0 {
		fmt.Println("⚠️  某些硬件信息需要root权限才能获取完整数据")
		fmt.Println("建议使用: sudo ./hwinfo")
		fmt.Println()
	}

	fmt.Println("正在收集硬件信息...")

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
	}

	info.OS, info.Kernel = getOSInfo()

	fmt.Print("\033[2J\033[H") // 清屏
	printServerInfo(info)
}
