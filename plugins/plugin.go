package plugins

import (
	"context"
	"fmt"
	"plugin"
	"sync"

	"github.com/devops-toolkit/clusterreport/pkg/analyzer"
	"github.com/devops-toolkit/clusterreport/pkg/collector"
)

// PluginType 插件类型
type PluginType string

const (
	PluginTypeCollector PluginType = "collector"
	PluginTypeAnalyzer  PluginType = "analyzer"
	PluginTypeOutput    PluginType = "output"
)

// Plugin 插件接口
type Plugin interface {
	// Name 插件名称
	Name() string

	// Type 插件类型
	Type() PluginType

	// Version 插件版本
	Version() string

	// Init 初始化插件
	Init(config map[string]interface{}) error

	// Cleanup 清理资源
	Cleanup() error
}

// Manager 插件管理器
type Manager struct {
	plugins map[string]Plugin
	mu      sync.RWMutex
}

// NewManager 创建插件管理器
func NewManager() *Manager {
	return &Manager{
		plugins: make(map[string]Plugin),
	}
}

// Load 加载插件
func (m *Manager) Load(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 打开插件文件
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin %s: %w", path, err)
	}

	// 查找New函数
	newFunc, err := p.Lookup("New")
	if err != nil {
		return fmt.Errorf("plugin %s does not export New function: %w", path, err)
	}

	// 创建插件实例
	createPlugin, ok := newFunc.(func() Plugin)
	if !ok {
		return fmt.Errorf("New function has wrong signature in plugin %s", path)
	}

	pluginInstance := createPlugin()

	// 注册插件
	name := pluginInstance.Name()
	if _, exists := m.plugins[name]; exists {
		return fmt.Errorf("plugin %s already loaded", name)
	}

	m.plugins[name] = pluginInstance
	return nil
}

// Unload 卸载插件
func (m *Manager) Unload(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	// 清理资源
	if err := plugin.Cleanup(); err != nil {
		return fmt.Errorf("failed to cleanup plugin %s: %w", name, err)
	}

	delete(m.plugins, name)
	return nil
}

// Get 获取插件
func (m *Manager) Get(name string) (Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return plugin, nil
}

// List 列出所有插件
func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.plugins))
	for name := range m.plugins {
		names = append(names, name)
	}
	return names
}

// CollectorPlugin 采集器插件接口
type CollectorPlugin interface {
	Plugin
	collector.Collector
}

// AnalyzerPlugin 分析器插件接口
type AnalyzerPlugin interface {
	Plugin
	analyzer.Analyzer
}

// OutputPlugin 输出插件接口
type OutputPlugin interface {
	Plugin

	// Generate 生成输出
	Generate(ctx context.Context, report *analyzer.Report) ([]byte, error)

	// Format 输出格式
	Format() string
}

// BasePlugin 插件基础实现
type BasePlugin struct {
	name    string
	pType   PluginType
	version string
	config  map[string]interface{}
}

// NewBasePlugin 创建基础插件
func NewBasePlugin(name string, pType PluginType, version string) *BasePlugin {
	return &BasePlugin{
		name:    name,
		pType:   pType,
		version: version,
		config:  make(map[string]interface{}),
	}
}

// Name 返回插件名称
func (p *BasePlugin) Name() string {
	return p.name
}

// Type 返回插件类型
func (p *BasePlugin) Type() PluginType {
	return p.pType
}

// Version 返回插件版本
func (p *BasePlugin) Version() string {
	return p.version
}

// Init 初始化插件
func (p *BasePlugin) Init(config map[string]interface{}) error {
	p.config = config
	return nil
}

// Cleanup 清理资源
func (p *BasePlugin) Cleanup() error {
	return nil
}

// CollectorAdapter 采集器适配器
type CollectorAdapter struct {
	*BasePlugin
	collector collector.Collector
}

// NewCollectorAdapter 创建采集器适配器
func NewCollectorAdapter(name, version string, c collector.Collector) *CollectorAdapter {
	return &CollectorAdapter{
		BasePlugin: NewBasePlugin(name, PluginTypeCollector, version),
		collector:  c,
	}
}

// Collect 采集数据
func (a *CollectorAdapter) Collect(ctx context.Context, node collector.Node) (*collector.Data, error) {
	return a.collector.Collect(ctx, node)
}

// Validate 验证配置
func (a *CollectorAdapter) Validate(config collector.Config) error {
	return a.collector.Validate(config)
}

// SupportedTypes 支持的数据类型
func (a *CollectorAdapter) SupportedTypes() []collector.DataType {
	return a.collector.SupportedTypes()
}

// AnalyzerAdapter 分析器适配器
type AnalyzerAdapter struct {
	*BasePlugin
	analyzer analyzer.Analyzer
}

// NewAnalyzerAdapter 创建分析器适配器
func NewAnalyzerAdapter(name, version string, a analyzer.Analyzer) *AnalyzerAdapter {
	return &AnalyzerAdapter{
		BasePlugin: NewBasePlugin(name, PluginTypeAnalyzer, version),
		analyzer:   a,
	}
}

// Analyze 分析数据
func (a *AnalyzerAdapter) Analyze(ctx context.Context, data []collector.Data) (*analyzer.Analysis, error) {
	return a.analyzer.Analyze(ctx, data)
}

// AnalysisType 分析类型
func (a *AnalyzerAdapter) AnalysisType() analyzer.AnalysisType {
	return a.analyzer.Type()
}

// Options 配置选项
func (a *AnalyzerAdapter) Options() map[string]interface{} {
	return a.analyzer.Options()
}
