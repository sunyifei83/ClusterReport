package main

import (
	"github.com/devops-toolkit/clusterreport/pkg/collector"
)

// SystemInfo 类型别名
type SystemInfo = collector.SystemInfo

// SystemCollector 类型别名
type SystemCollector = collector.SystemCollector

// CollectorConfig 类型别名
type CollectorConfig = collector.Config

// NewSystemCollector 创建系统采集器
func NewSystemCollector(config CollectorConfig, verbose bool) *SystemCollector {
	return collector.NewSystemCollector(config, verbose)
}
