package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

// DataType 数据类型
type DataType string

const (
	DataTypeHardware    DataType = "hardware"
	DataTypePerformance DataType = "performance"
	DataTypeNetwork     DataType = "network"
	DataTypeStorage     DataType = "storage"
	DataTypeCustom      DataType = "custom"
)

// Node 节点信息
type Node struct {
	Name     string            `json:"name"`
	Host     string            `json:"host"`
	Port     int               `json:"port"`
	Username string            `json:"username"`
	SSHKey   string            `json:"ssh_key"`
	Tags     []string          `json:"tags"`
	Metadata map[string]string `json:"metadata"`
}

// Data 采集的数据
type Data struct {
	Node      string                 `json:"node"`
	Type      DataType               `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Metrics   map[string]interface{} `json:"metrics"`
	Raw       []byte                 `json:"raw,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// CollectionResult 采集结果
type CollectionResult struct {
	Success   bool      `json:"success"`
	Node      string    `json:"node"`
	Data      []*Data   `json:"data,omitempty"`
	Error     string    `json:"error,omitempty"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Duration  string    `json:"duration"`
}

// Config 采集器配置
type Config struct {
	Timeout  time.Duration          `json:"timeout"`
	Parallel int                    `json:"parallel"`
	Retry    int                    `json:"retry"`
	Options  map[string]interface{} `json:"options"`
}

// Collector 数据采集器接口
type Collector interface {
	// Name 返回采集器名称
	Name() string

	// Collect 采集数据
	Collect(ctx context.Context, node Node) (*Data, error)

	// Validate 验证配置
	Validate(config Config) error

	// SupportedTypes 支持的数据类型
	SupportedTypes() []DataType
}

// Registry 采集器注册表
type Registry struct {
	collectors map[string]Collector
	mu         sync.RWMutex
}

// NewRegistry 创建新的注册表
func NewRegistry() *Registry {
	return &Registry{
		collectors: make(map[string]Collector),
	}
}

// Register 注册采集器
func (r *Registry) Register(collector Collector) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := collector.Name()
	if _, exists := r.collectors[name]; exists {
		return fmt.Errorf("collector %s already registered", name)
	}

	r.collectors[name] = collector
	return nil
}

// Get 获取采集器
func (r *Registry) Get(name string) (Collector, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	collector, exists := r.collectors[name]
	if !exists {
		return nil, fmt.Errorf("collector %s not found", name)
	}

	return collector, nil
}

// List 列出所有采集器
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.collectors))
	for name := range r.collectors {
		names = append(names, name)
	}
	return names
}

// MultiCollector 多采集器组合
type MultiCollector struct {
	collectors []Collector
	config     Config
}

// NewMultiCollector 创建多采集器
func NewMultiCollector(collectors ...Collector) *MultiCollector {
	return &MultiCollector{
		collectors: collectors,
		config: Config{
			Timeout:  5 * time.Minute,
			Parallel: 10,
			Retry:    3,
		},
	}
}

// CollectAll 从所有节点采集数据
func (m *MultiCollector) CollectAll(ctx context.Context, nodes []Node) []CollectionResult {
	results := make([]CollectionResult, 0, len(nodes))
	resultChan := make(chan CollectionResult, len(nodes))

	// 使用工作池限制并发
	semaphore := make(chan struct{}, m.config.Parallel)
	var wg sync.WaitGroup

	for _, node := range nodes {
		wg.Add(1)
		go func(n Node) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := m.collectFromNode(ctx, n)
			resultChan <- result
		}(node)
	}

	// 等待所有采集完成
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果
	for result := range resultChan {
		results = append(results, result)
	}

	return results
}

// collectFromNode 从单个节点采集数据
func (m *MultiCollector) collectFromNode(ctx context.Context, node Node) CollectionResult {
	result := CollectionResult{
		Node:      node.Name,
		StartTime: time.Now(),
		Data:      make([]*Data, 0),
	}

	// 设置超时
	timeoutCtx, cancel := context.WithTimeout(ctx, m.config.Timeout)
	defer cancel()

	// 执行所有采集器
	for _, collector := range m.collectors {
		data, err := collector.Collect(timeoutCtx, node)
		if err != nil {
			// 记录错误但继续执行其他采集器
			data = &Data{
				Node:      node.Name,
				Type:      DataTypeCustom,
				Timestamp: time.Now(),
				Error:     err.Error(),
			}
		}
		result.Data = append(result.Data, data)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime).String()
	result.Success = true

	// 检查是否有错误
	for _, data := range result.Data {
		if data.Error != "" {
			result.Success = false
			break
		}
	}

	return result
}

// Execute 在节点上执行命令
func (n *Node) Execute(ctx context.Context, command string) ([]byte, error) {
	// 构建SSH命令
	sshCmd := fmt.Sprintf("ssh -o StrictHostKeyChecking=no")

	if n.SSHKey != "" {
		sshCmd += fmt.Sprintf(" -i %s", n.SSHKey)
	}

	if n.Port > 0 {
		sshCmd += fmt.Sprintf(" -p %d", n.Port)
	}

	sshCmd += fmt.Sprintf(" %s@%s '%s'", n.Username, n.Host, command)

	// 执行命令
	cmd := exec.CommandContext(ctx, "sh", "-c", sshCmd)
	return cmd.Output()
}

// NodeProbeCollector NodeProbe采集器
type NodeProbeCollector struct {
	name string
}

// NewNodeProbeCollector 创建NodeProbe采集器
func NewNodeProbeCollector() *NodeProbeCollector {
	return &NodeProbeCollector{
		name: "nodeprobe",
	}
}

// Name 返回采集器名称
func (c *NodeProbeCollector) Name() string {
	return c.name
}

// Collect 采集数据
func (c *NodeProbeCollector) Collect(ctx context.Context, node Node) (*Data, error) {
	data := &Data{
		Node:      node.Name,
		Type:      DataTypeHardware,
		Timestamp: time.Now(),
		Metrics:   make(map[string]interface{}),
	}

	// 执行NodeProbe命令
	output, err := node.Execute(ctx, "/usr/local/bin/NodeProbe")
	if err != nil {
		data.Error = err.Error()
		return data, err
	}

	// 解析输出
	data.Raw = output
	if err := json.Unmarshal(output, &data.Metrics); err != nil {
		// 如果不是JSON格式，保存原始数据
		data.Metrics["raw"] = string(output)
	}

	return data, nil
}

// Validate 验证配置
func (c *NodeProbeCollector) Validate(config Config) error {
	return nil
}

// SupportedTypes 支持的数据类型
func (c *NodeProbeCollector) SupportedTypes() []DataType {
	return []DataType{DataTypeHardware}
}

// PerfSnapCollector PerfSnap采集器
type PerfSnapCollector struct {
	name string
}

// NewPerfSnapCollector 创建PerfSnap采集器
func NewPerfSnapCollector() *PerfSnapCollector {
	return &PerfSnapCollector{
		name: "perfsnap",
	}
}

// Name 返回采集器名称
func (c *PerfSnapCollector) Name() string {
	return c.name
}

// Collect 采集数据
func (c *PerfSnapCollector) Collect(ctx context.Context, node Node) (*Data, error) {
	data := &Data{
		Node:      node.Name,
		Type:      DataTypePerformance,
		Timestamp: time.Now(),
		Metrics:   make(map[string]interface{}),
	}

	// 执行PerfSnap命令
	output, err := node.Execute(ctx, "/usr/local/bin/PerfSnap")
	if err != nil {
		data.Error = err.Error()
		return data, err
	}

	// 解析输出
	data.Raw = output
	if err := json.Unmarshal(output, &data.Metrics); err != nil {
		// 如果不是JSON格式，保存原始数据
		data.Metrics["raw"] = string(output)
	}

	return data, nil
}

// Validate 验证配置
func (c *PerfSnapCollector) Validate(config Config) error {
	return nil
}

// SupportedTypes 支持的数据类型
func (c *PerfSnapCollector) SupportedTypes() []DataType {
	return []DataType{DataTypePerformance}
}

// CustomCollector 自定义采集器
type CustomCollector struct {
	name    string
	command string
	parser  func([]byte) (map[string]interface{}, error)
}

// NewCustomCollector 创建自定义采集器
func NewCustomCollector() *CustomCollector {
	return &CustomCollector{
		name:    "custom",
		command: "echo 'custom data'",
		parser:  defaultParser,
	}
}

// Name 返回采集器名称
func (c *CustomCollector) Name() string {
	return c.name
}

// Collect 采集数据
func (c *CustomCollector) Collect(ctx context.Context, node Node) (*Data, error) {
	data := &Data{
		Node:      node.Name,
		Type:      DataTypeCustom,
		Timestamp: time.Now(),
		Metrics:   make(map[string]interface{}),
	}

	// 执行自定义命令
	output, err := node.Execute(ctx, c.command)
	if err != nil {
		data.Error = err.Error()
		return data, err
	}

	// 解析输出
	data.Raw = output
	if c.parser != nil {
		metrics, err := c.parser(output)
		if err != nil {
			data.Metrics["raw"] = string(output)
			data.Metrics["parse_error"] = err.Error()
		} else {
			data.Metrics = metrics
		}
	} else {
		data.Metrics["raw"] = string(output)
	}

	return data, nil
}

// Validate 验证配置
func (c *CustomCollector) Validate(config Config) error {
	if c.command == "" {
		return fmt.Errorf("command is required for custom collector")
	}
	return nil
}

// SupportedTypes 支持的数据类型
func (c *CustomCollector) SupportedTypes() []DataType {
	return []DataType{DataTypeCustom}
}

// SetCommand 设置命令
func (c *CustomCollector) SetCommand(command string) {
	c.command = command
}

// SetParser 设置解析器
func (c *CustomCollector) SetParser(parser func([]byte) (map[string]interface{}, error)) {
	c.parser = parser
}

// defaultParser 默认解析器
func defaultParser(data []byte) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// 尝试JSON解析
	if err := json.Unmarshal(data, &result); err == nil {
		return result, nil
	}

	// 返回原始数据
	result["raw"] = string(data)
	return result, nil
}

// CachedCollector 带缓存的采集器
type CachedCollector struct {
	collector Collector
	cache     map[string]*cacheEntry
	ttl       time.Duration
	mu        sync.RWMutex
}

type cacheEntry struct {
	data      *Data
	timestamp time.Time
}

// NewCachedCollector 创建带缓存的采集器
func NewCachedCollector(collector Collector, ttl time.Duration) *CachedCollector {
	return &CachedCollector{
		collector: collector,
		cache:     make(map[string]*cacheEntry),
		ttl:       ttl,
	}
}

// Name 返回采集器名称
func (c *CachedCollector) Name() string {
	return c.collector.Name() + "_cached"
}

// Collect 采集数据（带缓存）
func (c *CachedCollector) Collect(ctx context.Context, node Node) (*Data, error) {
	key := fmt.Sprintf("%s-%s", c.collector.Name(), node.Name)

	// 尝试从缓存获取
	c.mu.RLock()
	if entry, exists := c.cache[key]; exists {
		if time.Since(entry.timestamp) < c.ttl {
			c.mu.RUnlock()
			return entry.data, nil
		}
	}
	c.mu.RUnlock()

	// 执行实际采集
	data, err := c.collector.Collect(ctx, node)
	if err != nil {
		return nil, err
	}

	// 保存到缓存
	c.mu.Lock()
	c.cache[key] = &cacheEntry{
		data:      data,
		timestamp: time.Now(),
	}
	c.mu.Unlock()

	return data, nil
}

// Validate 验证配置
func (c *CachedCollector) Validate(config Config) error {
	return c.collector.Validate(config)
}

// SupportedTypes 支持的数据类型
func (c *CachedCollector) SupportedTypes() []DataType {
	return c.collector.SupportedTypes()
}
