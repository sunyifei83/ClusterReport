package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// configCmd 代表 config 命令
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "管理配置文件",
	Long: `配置命令用于管理 ClusterReport 的配置文件。

支持的操作：
  • init    - 生成默认配置文件模板
  • show    - 显示当前配置
  • validate - 验证配置文件

配置文件支持：
  • 集群定义
  • 节点清单
  • SSH 配置
  • 输出设置
  • 收集器配置`,
}

// configInitCmd 生成配置文件模板
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "生成默认配置文件模板",
	Long: `生成一个默认的配置文件模板到当前目录。

如果文件已存在，会提示是否覆盖。

示例:
  clusterreport config init
  clusterreport config init --output custom-config.yaml`,
	RunE: runConfigInit,
}

// configShowCmd 显示当前配置
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前配置",
	Long: `显示当前加载的配置信息。

会显示：
  • 配置文件路径
  • 所有配置项的值
  • 配置来源（文件/环境变量/默认值）

示例:
  clusterreport config show
  clusterreport config show --format yaml`,
	RunE: runConfigShow,
}

// configValidateCmd 验证配置文件
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "验证配置文件",
	Long: `验证配置文件的语法和内容是否正确。

检查项：
  • YAML 语法正确性
  • 必需字段是否存在
  • 字段值是否有效
  • SSH 配置是否完整

示例:
  clusterreport config validate
  clusterreport config validate --config /path/to/config.yaml`,
	RunE: runConfigValidate,
}

var (
	configInitOutput   string
	configInitForce    bool
	configShowFormat   string
	configValidateFile string
)

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)

	// init 命令标志
	configInitCmd.Flags().StringVarP(&configInitOutput, "output", "o", "config.yaml", "输出文件路径")
	configInitCmd.Flags().BoolVarP(&configInitForce, "force", "f", false, "强制覆盖已存在的文件")

	// show 命令标志
	configShowCmd.Flags().StringVarP(&configShowFormat, "format", "f", "yaml", "输出格式: yaml, json")

	// validate 命令标志
	configValidateCmd.Flags().StringVar(&configValidateFile, "file", "", "要验证的配置文件路径（默认使用 --config 指定的文件）")
}

// runConfigInit 执行 config init 命令
func runConfigInit(cmd *cobra.Command, args []string) error {
	// 检查文件是否已存在
	if _, err := os.Stat(configInitOutput); err == nil && !configInitForce {
		return fmt.Errorf("配置文件 %s 已存在，使用 --force 强制覆盖", configInitOutput)
	}

	// 生成配置模板
	template := getConfigTemplate()

	// 写入文件
	if err := os.WriteFile(configInitOutput, []byte(template), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	if !quiet {
		fmt.Printf("✅ 配置文件模板已生成: %s\n", configInitOutput)
		fmt.Println("\n下一步:")
		fmt.Println("  1. 编辑配置文件，填入您的集群信息")
		fmt.Println("  2. 使用 'clusterreport config validate' 验证配置")
		fmt.Println("  3. 使用 'clusterreport collect --cluster <name>' 开始收集")
	}

	return nil
}

// runConfigShow 执行 config show 命令
func runConfigShow(cmd *cobra.Command, args []string) error {
	if !quiet {
		if viper.ConfigFileUsed() != "" {
			fmt.Printf("配置文件: %s\n\n", viper.ConfigFileUsed())
		} else {
			fmt.Println("未使用配置文件\n")
		}
	}

	// 获取所有配置
	settings := viper.AllSettings()

	var output []byte
	var err error

	switch configShowFormat {
	case "yaml":
		output, err = yaml.Marshal(settings)
	case "json":
		// JSON 输出将在后续实现
		return fmt.Errorf("JSON 格式暂未实现")
	default:
		return fmt.Errorf("不支持的输出格式: %s", configShowFormat)
	}

	if err != nil {
		return fmt.Errorf("格式化配置失败: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

// runConfigValidate 执行 config validate 命令
func runConfigValidate(cmd *cobra.Command, args []string) error {
	// 确定要验证的文件
	validateFile := configValidateFile
	if validateFile == "" {
		validateFile = viper.ConfigFileUsed()
		if validateFile == "" {
			validateFile = cfgFile
		}
	}

	if validateFile == "" {
		return fmt.Errorf("请指定要验证的配置文件路径")
	}

	if !quiet {
		fmt.Printf("🔍 验证配置文件: %s\n\n", validateFile)
	}

	// 读取配置文件
	data, err := os.ReadFile(validateFile)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 YAML
	var config AppConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("❌ YAML 语法错误: %w", err)
	}

	// 验证配置
	errors := validateAppConfig(&config)

	if len(errors) > 0 {
		fmt.Println("❌ 配置验证失败:\n")
		for i, e := range errors {
			fmt.Printf("  %d. %s\n", i+1, e)
		}
		return fmt.Errorf("发现 %d 个配置错误", len(errors))
	}

	if !quiet {
		fmt.Println("✅ 配置文件验证通过！")
		fmt.Printf("\n配置摘要:\n")
		fmt.Printf("  集群数量: %d\n", len(config.Clusters))
		totalNodes := 0
		for _, cluster := range config.Clusters {
			totalNodes += len(cluster.Nodes)
		}
		fmt.Printf("  节点总数: %d\n", totalNodes)
		fmt.Printf("  输出目录: %s\n", config.Output.Directory)
		fmt.Printf("  输出格式: %v\n", config.Output.Formats)
		fmt.Printf("  并发数: %d\n", config.Collector.Parallel)
	}

	return nil
}

// AppConfig 应用配置文件结构
type AppConfig struct {
	Clusters  []AppClusterConfig `yaml:"clusters"`
	Output    AppOutputConfig    `yaml:"output"`
	Collector AppCollectorConfig `yaml:"collector"`
	SSH       AppSSHConfig       `yaml:"ssh,omitempty"`
}

// AppClusterConfig 集群配置
type AppClusterConfig struct {
	Name     string            `yaml:"name"`
	Nodes    []string          `yaml:"nodes"`
	SSHKey   string            `yaml:"ssh_key,omitempty"`
	Username string            `yaml:"username,omitempty"`
	Port     int               `yaml:"port,omitempty"`
	Tags     map[string]string `yaml:"tags,omitempty"`
}

// AppOutputConfig 输出配置
type AppOutputConfig struct {
	Directory string   `yaml:"directory"`
	Formats   []string `yaml:"formats"`
}

// AppCollectorConfig 收集器配置
type AppCollectorConfig struct {
	Parallel int    `yaml:"parallel"`
	Timeout  string `yaml:"timeout"`
	Retry    int    `yaml:"retry"`
}

// AppSSHConfig SSH 全局配置
type AppSSHConfig struct {
	DefaultKey  string `yaml:"default_key,omitempty"`
	DefaultUser string `yaml:"default_user,omitempty"`
	DefaultPort int    `yaml:"default_port,omitempty"`
}

// validateAppConfig 验证配置
func validateAppConfig(config *AppConfig) []string {
	var errors []string

	// 验证集群配置
	if len(config.Clusters) == 0 {
		errors = append(errors, "至少需要定义一个集群")
	}

	clusterNames := make(map[string]bool)
	for i, cluster := range config.Clusters {
		// 检查集群名称
		if cluster.Name == "" {
			errors = append(errors, fmt.Sprintf("集群 #%d 缺少名称", i+1))
		} else if clusterNames[cluster.Name] {
			errors = append(errors, fmt.Sprintf("集群名称重复: %s", cluster.Name))
		} else {
			clusterNames[cluster.Name] = true
		}

		// 检查节点列表
		if len(cluster.Nodes) == 0 {
			errors = append(errors, fmt.Sprintf("集群 %s 没有定义任何节点", cluster.Name))
		}

		// 检查 SSH 配置
		if cluster.SSHKey != "" {
			if _, err := os.Stat(cluster.SSHKey); os.IsNotExist(err) {
				errors = append(errors, fmt.Sprintf("集群 %s 的 SSH 密钥文件不存在: %s", cluster.Name, cluster.SSHKey))
			}
		}

		// 检查端口
		if cluster.Port != 0 && (cluster.Port < 1 || cluster.Port > 65535) {
			errors = append(errors, fmt.Sprintf("集群 %s 的端口号无效: %d", cluster.Name, cluster.Port))
		}
	}

	// 验证输出配置
	if config.Output.Directory == "" {
		errors = append(errors, "输出目录未指定")
	}

	if len(config.Output.Formats) == 0 {
		errors = append(errors, "至少需要指定一种输出格式")
	}

	validFormats := map[string]bool{"html": true, "json": true, "markdown": true, "yaml": true}
	for _, format := range config.Output.Formats {
		if !validFormats[format] {
			errors = append(errors, fmt.Sprintf("不支持的输出格式: %s", format))
		}
	}

	// 验证收集器配置
	if config.Collector.Parallel < 1 {
		errors = append(errors, "并发数必须大于 0")
	}

	if config.Collector.Parallel > 100 {
		errors = append(errors, "并发数过大（最大 100）")
	}

	if config.Collector.Retry < 0 {
		errors = append(errors, "重试次数不能为负数")
	}

	return errors
}

// getConfigTemplate 获取配置文件模板
func getConfigTemplate() string {
	return `# ClusterReport 配置文件
# 文档: https://github.com/sunyifei83/ClusterReport

# 集群定义
clusters:
  # 生产环境集群
  - name: production
    nodes:
      - prod-node1.example.com
      - prod-node2.example.com
      - prod-node3.example.com
    ssh_key: ~/.ssh/id_rsa
    username: admin
    port: 22
    tags:
      env: production
      region: us-east-1

  # 测试环境集群
  - name: staging
    nodes:
      - staging-node1.example.com
      - staging-node2.example.com
    ssh_key: ~/.ssh/staging_key
    username: deploy
    tags:
      env: staging
      region: us-west-1

  # 开发环境（本地）
  - name: development
    nodes:
      - localhost
    tags:
      env: development

# 输出配置
output:
  # 报告输出目录
  directory: ./reports
  
  # 支持的输出格式: html, json, markdown, yaml
  formats:
    - html
    - json
    - markdown

# 收集器配置
collector:
  # 并发收集的节点数
  parallel: 10
  
  # 收集超时时间
  timeout: 5m
  
  # 失败重试次数
  retry: 3

# SSH 全局配置（可选，会被集群配置覆盖）
ssh:
  default_key: ~/.ssh/id_rsa
  default_user: root
  default_port: 22
`
}
