package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
	quiet   bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "clusterreport",
	Short: "ClusterReport - 集群综合报告生成工具",
	Long: `ClusterReport 是一个功能强大的集群分析和报告生成工具。

它能够自动收集、分析集群节点的配置和性能数据，并生成多种格式的综合报告。

特性：
  • 多源数据采集（本地/远程/批量）
  • 智能分析和评分
  • 多格式报告生成（HTML/JSON/Markdown）
  • 可扩展的插件系统

版本：v0.8.0 (开发中)
文档：https://github.com/devops-toolkit/clusterreport`,
	Version: "0.8.0-dev",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全局标志
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: ./config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "详细输出模式")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "静默模式（仅输出错误）")

	// 绑定到 viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
}

// initConfig 读取配置文件和环境变量
func initConfig() {
	if cfgFile != "" {
		// 使用指定的配置文件
		viper.SetConfigFile(cfgFile)
	} else {
		// 查找默认配置文件
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.clusterreport")
		viper.AddConfigPath("/etc/clusterreport")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	// 读取环境变量
	viper.SetEnvPrefix("CLUSTERREPORT")
	viper.AutomaticEnv()

	// 读取配置文件
	if err := viper.ReadInConfig(); err == nil {
		if !quiet {
			fmt.Fprintln(os.Stderr, "使用配置文件:", viper.ConfigFileUsed())
		}
	}
}
