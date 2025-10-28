package collector

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
	"time"
)

// NodeProbeCollector 包装 NodeProbe 的功能
type NodeProbeCollector struct {
	autoOptimize bool // 是否启用自动优化功能
}

// NodeProbeData 存储 NodeProbe 收集的数据
type NodeProbeData struct {
	Hostname      string                 `json:"hostname" yaml:"hostname"`
	LoadAverage   string                 `json:"load_average" yaml:"load_average"`
	Timezone      string                 `json:"timezone" yaml:"timezone"`
	OS            string                 `json:"os" yaml:"os"`
	Kernel        string                 `json:"kernel" yaml:"kernel"`
	CPU           NodeProbeCPUInfo       `json:"cpu" yaml:"cpu"`
	Memory        NodeProbeMemoryInfo    `json:"memory" yaml:"memory"`
	Disks         NodeProbeDiskInfo      `json:"disks" yaml:"disks"`
	Network       []NodeProbeNetworkIF   `json:"network" yaml:"network"`
	Python        NodeProbePythonInfo    `json:"python" yaml:"python"`
	Java          NodeProbeJavaInfo      `json:"java" yaml:"java"`
	KernelModules NodeProbeKernelModules `json:"kernel_modules" yaml:"kernel_modules"`
	Timestamp     string                 `json:"timestamp" yaml:"timestamp"`
	Version       string                 `json:"nodeprobe_version" yaml:"nodeprobe_version"`
}

type NodeProbeCPUInfo struct {
	Model           string `json:"model" yaml:"model"`
	Cores           int    `json:"cores" yaml:"cores"`
	RunMode         string `json:"run_mode" yaml:"run_mode"`
	PerformanceMode string `json:"performance_mode" yaml:"performance_mode"`
}

type NodeProbeMemoryInfo struct {
	TotalGB float64               `json:"total_gb" yaml:"total_gb"`
	Slots   []NodeProbeMemorySlot `json:"slots" yaml:"slots"`
}

type NodeProbeMemorySlot struct {
	Location string `json:"location" yaml:"location"`
	Size     string `json:"size" yaml:"size"`
}

type NodeProbeDiskInfo struct {
	SystemDisk  string   `json:"system_disk" yaml:"system_disk"`
	DataDisks   []string `json:"data_disks" yaml:"data_disks"`
	TotalDisks  int      `json:"total_disks" yaml:"total_disks"`
	DataDiskNum int      `json:"data_disk_num" yaml:"data_disk_num"`
}

type NodeProbeNetworkIF struct {
	Name   string `json:"name" yaml:"name"`
	Status string `json:"status" yaml:"status"`
	Speed  string `json:"speed" yaml:"speed"`
	IP     string `json:"ip" yaml:"ip"`
}

type NodeProbePythonInfo struct {
	Version string `json:"version" yaml:"version"`
	Path    string `json:"path" yaml:"path"`
}

type NodeProbeJavaInfo struct {
	Version string `json:"version" yaml:"version"`
	Path    string `json:"path" yaml:"path"`
}

type NodeProbeKernelModules struct {
	NfConntrack bool   `json:"nf_conntrack" yaml:"nf_conntrack"`
	BrNetfilter bool   `json:"br_netfilter" yaml:"br_netfilter"`
	Message     string `json:"message" yaml:"message"`
}

// NewNodeProbeCollector 创建新的 NodeProbe 收集器
func NewNodeProbeCollector(autoOptimize bool) *NodeProbeCollector {
	return &NodeProbeCollector{
		autoOptimize: autoOptimize,
	}
}

// Collect 执行 NodeProbe 数据收集
func (c *NodeProbeCollector) Collect() (*NodeProbeData, error) {
	data := &NodeProbeData{
		Version:   "1.1.1",
		Timestamp: getCurrentTimestamp(),
	}

	// 收集各项信息
	data.Hostname = c.getHostname()
	data.LoadAverage = c.getLoadAverage()
	data.Timezone = c.getTimezone()
	data.OS, data.Kernel = c.getOSInfo()
	data.CPU = c.getCPUInfo()
	data.Memory = c.getMemoryInfo()
	data.Disks = c.getDiskInfo()
	data.Network = c.getNetworkInfo()
	data.Python = c.getPythonInfo()
	data.Java = c.getJavaInfo()
	data.KernelModules = c.checkAndLoadKernelModules()

	return data, nil
}

// 获取主机名
func (c *NodeProbeCollector) getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

// 获取系统负载
func (c *NodeProbeCollector) getLoadAverage() string {
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

// 获取时区信息并校准
func (c *NodeProbeCollector) getTimezone() string {
	var currentTZ string

	// 尝试使用timedatectl命令
	if output, err := execCommand("timedatectl", "show", "--property=Timezone", "--value"); err == nil {
		currentTZ = strings.TrimSpace(output)
		if currentTZ != "" {
			// 如果启用自动优化且当前时区不是Asia/Shanghai，尝试设置
			if c.autoOptimize && currentTZ != "Asia/Shanghai" && os.Geteuid() == 0 {
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
			if c.autoOptimize && currentTZ != "Asia/Shanghai" && os.Geteuid() == 0 {
				if err := os.WriteFile("/etc/timezone", []byte("Asia/Shanghai\n"), 0644); err == nil {
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
		parts := strings.Split(link, "zoneinfo/")
		if len(parts) == 2 {
			currentTZ = parts[1]
			if c.autoOptimize && currentTZ != "Asia/Shanghai" && os.Geteuid() == 0 {
				if _, err := execCommand("ln", "-sf", "/usr/share/zoneinfo/Asia/Shanghai", "/etc/localtime"); err == nil {
					return fmt.Sprintf("已校准至 Asia/Shanghai (原: %s)", currentTZ)
				}
			}
			return currentTZ
		}
	}

	return "Unknown"
}

// 获取操作系统信息
func (c *NodeProbeCollector) getOSInfo() (string, string) {
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
func (c *NodeProbeCollector) getCPUInfo() NodeProbeCPUInfo {
	info := NodeProbeCPUInfo{}

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
	info.RunMode = c.getCPURunMode()
	info.PerformanceMode = c.getCPUPerformanceMode()

	return info
}

// 获取CPU运行模式
func (c *NodeProbeCollector) getCPURunMode() string {
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

	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		if strings.Contains(string(data), "lm") {
			return "32-bit, 64-bit"
		}
		return "32-bit"
	}

	return "Unknown"
}

// 获取CPU性能模式
func (c *NodeProbeCollector) getCPUPerformanceMode() string {
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

	// 如果启用自动优化且当前是省电模式且以root权限运行，自动切换到性能模式
	if c.autoOptimize && currentGovernor == "powersave" && os.Geteuid() == 0 {
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
	}

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

	return "未知 (可能不支持频率调节)"
}

// 获取内存信息
func (c *NodeProbeCollector) getMemoryInfo() NodeProbeMemoryInfo {
	info := NodeProbeMemoryInfo{}

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
		var currentSlot NodeProbeMemorySlot
		var hasSize bool

		for _, line := range lines {
			line = strings.TrimSpace(line)

			if strings.Contains(line, "Memory Device") {
				if hasSize && currentSlot.Size != "" {
					info.Slots = append(info.Slots, currentSlot)
				}
				currentSlot = NodeProbeMemorySlot{}
				hasSize = false
			}

			if strings.HasPrefix(line, "Locator:") && !strings.HasPrefix(line, "Bank Locator:") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					currentSlot.Location = strings.TrimSpace(parts[1])
				}
			}

			if strings.HasPrefix(line, "Size:") {
				if !strings.Contains(line, "No Module") &&
					!strings.Contains(line, "None") &&
					!strings.Contains(line, "Unknown") {
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						sizeInfo := strings.Join(parts[1:], " ")
						if sizeInfo != "" && sizeInfo != "0" {
							currentSlot.Size = sizeInfo
							hasSize = true
						}
					}
				}
			}
		}

		if hasSize && currentSlot.Size != "" {
			info.Slots = append(info.Slots, currentSlot)
		}
	}

	return info
}

// 获取磁盘信息
func (c *NodeProbeCollector) getDiskInfo() NodeProbeDiskInfo {
	info := NodeProbeDiskInfo{}

	// 获取系统盘挂载信息
	if output, err := execCommand("df", "-h", "/"); err == nil {
		lines := strings.Split(output, "\n")
		if len(lines) > 1 {
			for i := 1; i < len(lines); i++ {
				line := lines[i]
				if line == "" {
					continue
				}

				fields := strings.Fields(line)
				if len(fields) >= 6 && fields[5] == "/" {
					info.SystemDisk = fmt.Sprintf("%s %s/%s (%s)",
						fields[0], fields[3], fields[2], fields[4])
					break
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
						(strings.Contains(fields[1], "G") && c.parseSize(fields[1]) > 1000) {
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
func (c *NodeProbeCollector) parseSize(sizeStr string) float64 {
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
func (c *NodeProbeCollector) getNetworkInfo() []NodeProbeNetworkIF {
	var interfaces []NodeProbeNetworkIF

	if output, err := execCommand("ip", "addr", "show"); err == nil {
		lines := strings.Split(output, "\n")
		var currentIface NodeProbeNetworkIF

		for _, line := range lines {
			if matches := regexp.MustCompile(`^\d+:\s+(\S+):\s+<(.+?)>`).FindStringSubmatch(line); len(matches) > 2 {
				if currentIface.Name != "" {
					interfaces = append(interfaces, currentIface)
				}
				currentIface = NodeProbeNetworkIF{Name: matches[1]}
				if strings.Contains(matches[2], "UP") {
					currentIface.Status = "UP"
				} else {
					currentIface.Status = "DOWN"
				}
			}

			if strings.Contains(line, "inet ") && !strings.Contains(line, "inet6") {
				fields := strings.Fields(line)
				for i, field := range fields {
					if field == "inet" && i+1 < len(fields) {
						currentIface.IP = fields[i+1]
						break
					}
				}
			}

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
	var filtered []NodeProbeNetworkIF
	for _, iface := range interfaces {
		if iface.Name != "lo" {
			filtered = append(filtered, iface)
		}
	}

	return filtered
}

// 获取Python信息
func (c *NodeProbeCollector) getPythonInfo() NodeProbePythonInfo {
	info := NodeProbePythonInfo{}

	pythonCommands := []string{"python3", "python", "python2"}

	for _, cmd := range pythonCommands {
		if output, err := execCommand(cmd, "--version"); err == nil {
			versionStr := strings.TrimSpace(output)
			if versionStr != "" {
				info.Version = versionStr

				if pathOutput, err := execCommand("which", cmd); err == nil {
					info.Path = strings.TrimSpace(pathOutput)
				}
				break
			}
		}
	}

	if info.Version == "" {
		info.Version = "Not installed"
		info.Path = "N/A"
	}

	return info
}

// 获取Java信息
func (c *NodeProbeCollector) getJavaInfo() NodeProbeJavaInfo {
	info := NodeProbeJavaInfo{}

	if output, err := execCommand("java", "-version"); err == nil {
		lines := strings.Split(output, "\n")
		for _, line := range lines {
			if strings.Contains(line, "version") {
				if matches := regexp.MustCompile(`"([^"]+)"`).FindStringSubmatch(line); len(matches) > 1 {
					info.Version = "Java " + matches[1]
				} else {
					info.Version = strings.TrimSpace(line)
				}
				break
			}
		}

		if pathOutput, err := execCommand("which", "java"); err == nil {
			info.Path = strings.TrimSpace(pathOutput)

			if javaHome := os.Getenv("JAVA_HOME"); javaHome != "" {
				info.Path = fmt.Sprintf("%s (JAVA_HOME: %s)", info.Path, javaHome)
			}
		}
	}

	if info.Version == "" {
		info.Version = "Not installed"
		info.Path = "N/A"
	}

	return info
}

// 检查和加载内核模块
func (c *NodeProbeCollector) checkAndLoadKernelModules() NodeProbeKernelModules {
	status := NodeProbeKernelModules{}

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

		// 如果启用自动优化且未加载且是root用户，尝试加载
		if c.autoOptimize && os.Geteuid() == 0 {
			if _, err := execCommand("modprobe", mod.name); err == nil {
				*mod.field = true
				messages = append(messages, fmt.Sprintf("%s: 已自动加载", mod.name))
			} else {
				*mod.field = false
				messages = append(messages, fmt.Sprintf("%s: 加载失败", mod.name))
			}
		} else {
			*mod.field = false
			if os.Geteuid() != 0 {
				messages = append(messages, fmt.Sprintf("%s: 未加载(需要root权限)", mod.name))
			} else {
				messages = append(messages, fmt.Sprintf("%s: 未加载", mod.name))
			}
		}
	}

	status.Message = strings.Join(messages, ", ")
	return status
}

// execCommand 执行命令并返回输出
func execCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// getCurrentTimestamp 获取当前时间戳
func getCurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
