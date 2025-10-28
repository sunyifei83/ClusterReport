package collector

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/crypto/ssh"
)

// SystemInfo 系统信息
type SystemInfo struct {
	Hostname     string            `json:"hostname"`
	OS           string            `json:"os"`
	Kernel       string            `json:"kernel"`
	CPUInfo      CPUInfo           `json:"cpu_info"`
	MemoryInfo   MemoryInfo        `json:"memory_info"`
	DiskInfo     []DiskInfo        `json:"disk_info"`
	NetworkInfo  []NetworkInfo     `json:"network_info"`
	LoadAverage  LoadAverage       `json:"load_average"`
	Uptime       string            `json:"uptime"`
	CollectedAt  time.Time         `json:"collected_at"`
	CollectError string            `json:"collect_error,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// CPUInfo CPU 信息
type CPUInfo struct {
	Model       string  `json:"model"`
	Cores       int     `json:"cores"`
	Threads     int     `json:"threads"`
	Usage       float64 `json:"usage"`
	Temperature float64 `json:"temperature,omitempty"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	Available uint64  `json:"available"`
	UsageRate float64 `json:"usage_rate"`
}

// DiskInfo 磁盘信息
type DiskInfo struct {
	Device     string  `json:"device"`
	MountPoint string  `json:"mount_point"`
	FSType     string  `json:"fs_type"`
	Total      uint64  `json:"total"`
	Used       uint64  `json:"used"`
	Free       uint64  `json:"free"`
	UsageRate  float64 `json:"usage_rate"`
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	Interface string `json:"interface"`
	IPAddress string `json:"ip_address"`
	MACAddr   string `json:"mac_address"`
	RxBytes   uint64 `json:"rx_bytes"`
	TxBytes   uint64 `json:"tx_bytes"`
	Status    string `json:"status"`
}

// LoadAverage 系统负载
type LoadAverage struct {
	Load1  float64 `json:"load1"`
	Load5  float64 `json:"load5"`
	Load15 float64 `json:"load15"`
}

// SystemCollector 系统信息采集器
type SystemCollector struct {
	config  Config
	verbose bool
}

// NewSystemCollector 创建系统采集器
func NewSystemCollector(config Config, verbose bool) *SystemCollector {
	return &SystemCollector{
		config:  config,
		verbose: verbose,
	}
}

// CollectLocal 采集本地系统信息
func (sc *SystemCollector) CollectLocal(ctx context.Context) (*SystemInfo, error) {
	info := &SystemInfo{
		CollectedAt: time.Now(),
		Metadata:    make(map[string]string),
	}

	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}
	info.Hostname = hostname
	info.OS = runtime.GOOS
	info.Metadata["arch"] = runtime.GOARCH

	// 获取内核版本
	if kernel, err := sc.getKernelVersion(); err == nil {
		info.Kernel = kernel
	}

	// 获取 CPU 信息
	if cpuInfo, err := sc.getCPUInfo(); err == nil {
		info.CPUInfo = cpuInfo
	}

	// 获取内存信息
	if memInfo, err := sc.getMemoryInfo(); err == nil {
		info.MemoryInfo = memInfo
	}

	// 获取磁盘信息
	if diskInfo, err := sc.getDiskInfo(); err == nil {
		info.DiskInfo = diskInfo
	}

	// 获取网络信息
	if netInfo, err := sc.getNetworkInfo(); err == nil {
		info.NetworkInfo = netInfo
	}

	// 获取系统负载
	if load, err := sc.getLoadAverage(); err == nil {
		info.LoadAverage = load
	}

	// 获取运行时间
	if uptime, err := sc.getUptime(); err == nil {
		info.Uptime = uptime
	}

	return info, nil
}

// CollectRemote 采集远程系统信息
func (sc *SystemCollector) CollectRemote(ctx context.Context, host string, sshConfig *ssh.ClientConfig) (*SystemInfo, error) {
	client, err := ssh.Dial("tcp", host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", host, err)
	}
	defer client.Close()

	info := &SystemInfo{
		CollectedAt: time.Now(),
		Metadata:    make(map[string]string),
	}
	info.Metadata["remote_host"] = host

	// 执行远程命令获取系统信息
	commands := map[string]*string{
		"hostname":  &info.Hostname,
		"uname -s":  &info.OS,
		"uname -r":  &info.Kernel,
		"uptime -p": &info.Uptime,
	}

	for cmd, target := range commands {
		output, err := sc.executeRemoteCommand(client, cmd)
		if err == nil {
			*target = strings.TrimSpace(output)
		}
	}

	// 获取更详细的信息
	if cpuInfo, err := sc.getRemoteCPUInfo(client); err == nil {
		info.CPUInfo = cpuInfo
	}

	if memInfo, err := sc.getRemoteMemoryInfo(client); err == nil {
		info.MemoryInfo = memInfo
	}

	if diskInfo, err := sc.getRemoteDiskInfo(client); err == nil {
		info.DiskInfo = diskInfo
	}

	return info, nil
}

// CollectMultiple 并发采集多个节点
func (sc *SystemCollector) CollectMultiple(ctx context.Context, nodes []string, parallel int) map[string]*SystemInfo {
	results := make(map[string]*SystemInfo)
	var mu sync.Mutex
	var wg sync.WaitGroup

	// 创建进度条
	bar := progressbar.NewOptions(len(nodes),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription(color.CyanString("[cyan][1/4][reset] Collecting system info...")),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// 创建工作池
	semaphore := make(chan struct{}, parallel)

	for _, node := range nodes {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			var info *SystemInfo
			var err error

			// 判断是本地还是远程
			if n == "localhost" || n == "127.0.0.1" {
				info, err = sc.CollectLocal(ctx)
			} else {
				// 远程采集（这里简化处理，实际需要 SSH 配置）
				info = &SystemInfo{
					Hostname:     n,
					CollectedAt:  time.Now(),
					CollectError: "Remote collection not fully implemented yet",
				}
			}

			mu.Lock()
			if err != nil {
				info = &SystemInfo{
					Hostname:     n,
					CollectedAt:  time.Now(),
					CollectError: err.Error(),
				}
			}
			results[n] = info
			mu.Unlock()

			bar.Add(1)
		}(node)
	}

	wg.Wait()
	bar.Finish()
	fmt.Println()

	return results
}

// 辅助方法：获取内核版本
func (sc *SystemCollector) getKernelVersion() (string, error) {
	cmd := exec.Command("uname", "-r")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// 辅助方法：获取 CPU 信息
func (sc *SystemCollector) getCPUInfo() (CPUInfo, error) {
	info := CPUInfo{
		Cores:   runtime.NumCPU(),
		Threads: runtime.NumCPU(),
	}

	// 获取 CPU 型号（Linux）
	if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "cat /proc/cpuinfo | grep 'model name' | head -1 | cut -d: -f2")
		if output, err := cmd.Output(); err == nil {
			info.Model = strings.TrimSpace(string(output))
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("sysctl", "-n", "machdep.cpu.brand_string")
		if output, err := cmd.Output(); err == nil {
			info.Model = strings.TrimSpace(string(output))
		}
	}

	return info, nil
}

// 辅助方法：获取内存信息
func (sc *SystemCollector) getMemoryInfo() (MemoryInfo, error) {
	info := MemoryInfo{}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("sh", "-c", "free -b | grep Mem:")
		output, err := cmd.Output()
		if err != nil {
			return info, err
		}

		fields := strings.Fields(string(output))
		if len(fields) >= 7 {
			fmt.Sscanf(fields[1], "%d", &info.Total)
			fmt.Sscanf(fields[2], "%d", &info.Used)
			fmt.Sscanf(fields[3], "%d", &info.Free)
			fmt.Sscanf(fields[6], "%d", &info.Available)
			if info.Total > 0 {
				info.UsageRate = float64(info.Used) / float64(info.Total) * 100
			}
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("sysctl", "-n", "hw.memsize")
		if output, err := cmd.Output(); err == nil {
			fmt.Sscanf(strings.TrimSpace(string(output)), "%d", &info.Total)
		}
	}

	return info, nil
}

// 辅助方法：获取磁盘信息
func (sc *SystemCollector) getDiskInfo() ([]DiskInfo, error) {
	var disks []DiskInfo

	cmd := exec.Command("df", "-B1")
	output, err := cmd.Output()
	if err != nil {
		return disks, err
	}

	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 6 {
			disk := DiskInfo{
				Device:     fields[0],
				MountPoint: fields[5],
			}

			fmt.Sscanf(fields[1], "%d", &disk.Total)
			fmt.Sscanf(fields[2], "%d", &disk.Used)
			fmt.Sscanf(fields[3], "%d", &disk.Free)

			if disk.Total > 0 {
				disk.UsageRate = float64(disk.Used) / float64(disk.Total) * 100
			}

			disks = append(disks, disk)
		}
	}

	return disks, nil
}

// 辅助方法：获取网络信息
func (sc *SystemCollector) getNetworkInfo() ([]NetworkInfo, error) {
	var networks []NetworkInfo

	if runtime.GOOS == "linux" {
		cmd := exec.Command("ip", "-o", "addr", "show")
		output, err := cmd.Output()
		if err != nil {
			return networks, err
		}

		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= 4 {
				net := NetworkInfo{
					Interface: fields[1],
					Status:    "up",
				}

				// 提取 IP 地址
				for i, field := range fields {
					if field == "inet" && i+1 < len(fields) {
						net.IPAddress = strings.Split(fields[i+1], "/")[0]
						break
					}
				}

				if net.IPAddress != "" {
					networks = append(networks, net)
				}
			}
		}
	}

	return networks, nil
}

// 辅助方法：获取系统负载
func (sc *SystemCollector) getLoadAverage() (LoadAverage, error) {
	load := LoadAverage{}

	if runtime.GOOS == "linux" {
		cmd := exec.Command("cat", "/proc/loadavg")
		output, err := cmd.Output()
		if err != nil {
			return load, err
		}

		fields := strings.Fields(string(output))
		if len(fields) >= 3 {
			fmt.Sscanf(fields[0], "%f", &load.Load1)
			fmt.Sscanf(fields[1], "%f", &load.Load5)
			fmt.Sscanf(fields[2], "%f", &load.Load15)
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("sysctl", "-n", "vm.loadavg")
		output, err := cmd.Output()
		if err != nil {
			return load, err
		}

		// macOS 格式: { 1.23 2.34 3.45 }
		str := strings.Trim(string(output), "{} \n")
		fields := strings.Fields(str)
		if len(fields) >= 3 {
			fmt.Sscanf(fields[0], "%f", &load.Load1)
			fmt.Sscanf(fields[1], "%f", &load.Load5)
			fmt.Sscanf(fields[2], "%f", &load.Load15)
		}
	}

	return load, nil
}

// 辅助方法：获取运行时间
func (sc *SystemCollector) getUptime() (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "linux" {
		cmd = exec.Command("uptime", "-p")
	} else if runtime.GOOS == "darwin" {
		cmd = exec.Command("uptime")
	} else {
		return "", fmt.Errorf("unsupported OS")
	}

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

// 远程执行命令的辅助方法
func (sc *SystemCollector) executeRemoteCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf

	if err := session.Run(command); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// 获取远程 CPU 信息
func (sc *SystemCollector) getRemoteCPUInfo(client *ssh.Client) (CPUInfo, error) {
	info := CPUInfo{}

	// 获取 CPU 核心数
	if output, err := sc.executeRemoteCommand(client, "nproc"); err == nil {
		fmt.Sscanf(strings.TrimSpace(output), "%d", &info.Cores)
		info.Threads = info.Cores
	}

	// 获取 CPU 型号
	if output, err := sc.executeRemoteCommand(client, "cat /proc/cpuinfo | grep 'model name' | head -1 | cut -d: -f2"); err == nil {
		info.Model = strings.TrimSpace(output)
	}

	return info, nil
}

// 获取远程内存信息
func (sc *SystemCollector) getRemoteMemoryInfo(client *ssh.Client) (MemoryInfo, error) {
	info := MemoryInfo{}

	output, err := sc.executeRemoteCommand(client, "free -b | grep Mem:")
	if err != nil {
		return info, err
	}

	fields := strings.Fields(output)
	if len(fields) >= 7 {
		fmt.Sscanf(fields[1], "%d", &info.Total)
		fmt.Sscanf(fields[2], "%d", &info.Used)
		fmt.Sscanf(fields[3], "%d", &info.Free)
		fmt.Sscanf(fields[6], "%d", &info.Available)
		if info.Total > 0 {
			info.UsageRate = float64(info.Used) / float64(info.Total) * 100
		}
	}

	return info, nil
}

// 获取远程磁盘信息
func (sc *SystemCollector) getRemoteDiskInfo(client *ssh.Client) ([]DiskInfo, error) {
	var disks []DiskInfo

	output, err := sc.executeRemoteCommand(client, "df -B1")
	if err != nil {
		return disks, err
	}

	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if i == 0 || line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 6 {
			disk := DiskInfo{
				Device:     fields[0],
				MountPoint: fields[5],
			}

			fmt.Sscanf(fields[1], "%d", &disk.Total)
			fmt.Sscanf(fields[2], "%d", &disk.Used)
			fmt.Sscanf(fields[3], "%d", &disk.Free)

			if disk.Total > 0 {
				disk.UsageRate = float64(disk.Used) / float64(disk.Total) * 100
			}

			disks = append(disks, disk)
		}
	}

	return disks, nil
}

// PrintSummary 打印采集摘要
func (sc *SystemCollector) PrintSummary(results map[string]*SystemInfo) {
	fmt.Println()
	color.Cyan("=== Collection Summary ===")
	fmt.Println()

	successCount := 0
	failCount := 0

	for node, info := range results {
		if info.CollectError == "" {
			successCount++
			color.Green("✓ %s: %s (%s)", node, info.Hostname, info.OS)
			if sc.verbose {
				fmt.Printf("  CPU: %s (%d cores)\n", info.CPUInfo.Model, info.CPUInfo.Cores)
				fmt.Printf("  Memory: %.2f GB (%.1f%% used)\n",
					float64(info.MemoryInfo.Total)/1024/1024/1024,
					info.MemoryInfo.UsageRate)
				fmt.Printf("  Load: %.2f, %.2f, %.2f\n",
					info.LoadAverage.Load1,
					info.LoadAverage.Load5,
					info.LoadAverage.Load15)
			}
		} else {
			failCount++
			color.Red("✗ %s: %s", node, info.CollectError)
		}
	}

	fmt.Println()
	color.Cyan("Total: %d nodes, %d succeeded, %d failed",
		len(results), successCount, failCount)
	fmt.Println()
}
