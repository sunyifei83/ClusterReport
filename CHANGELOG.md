# 更新日志

本文档记录 DevOps Toolkit 项目的所有重要变更。

格式基于 [Keep a Changelog](https://keepachangelog.com/zh-CN/1.0.0/)，
并且本项目遵循 [语义化版本](https://semver.org/lang/zh-CN/)。

## [未发布]

### 新增
- 项目全面审计报告
- 完整的用户指南文档
- 最佳实践指南文档
- 贡献者行为准则

### 修复
- 修复文档中的链接错误

## [1.1.1] - 2025-09-16

### 工具更新

#### NodeProbe v1.1.1
- 新增多格式输出支持（表格、JSON、YAML）
- 新增自动系统优化功能（CPU性能模式、时区校准）
- 新增内核模块自动加载功能
- 改进中文字符显示对齐
- 优化内存插槽信息显示（限制显示数量）

#### PerfSnap v1.1.0
- 新增火焰图生成功能
- 改进性能数据采集准确性
- 新增更多系统指标监控

#### DocConverter v1.1.0
- 新增网页爬取并转换为PDF功能
- 改进Markdown渲染效果
- 优化PDF样式和格式

## [1.0.0] - 2025-09-14

### 新增

#### Go 工具
- **NodeProbe v1.0.0** - Linux节点配置探测工具
  - 硬件配置信息采集
  - 系统状态监控
  - 软件环境检测
  - 支持JSON输出格式

- **PerfSnap v1.0.0** - 性能快照分析工具
  - CPU性能监控
  - 内存使用分析
  - 磁盘IO监控
  - 网络流量统计

- **DocConverter v1.0.0** - 文档转换工具
  - Markdown转PDF
  - HTML转PDF
  - 批量转换支持

#### Shell 工具
- **iotest.sh v1.0.0** - 磁盘IO性能测试工具
  - 基于fio的性能测试
  - 支持多种块大小测试
  - 详细的IOPS和吞吐量报告

- **clear_log.sh v1.0.0** - 日志清理工具
  - 安全的日志清理
  - 支持多种日志格式

#### ClusterReport (开发中)
- 基础项目架构
- CLI命令框架
- 数据采集器模块
- 数据分析器模块
- 报告生成器模块
- 插件系统基础

#### 文档
- 项目README（中英文）
- 工具文档总览
- 各工具详细文档
- MIT许可证
- 贡献指南

#### 配置和模板
- .gitignore配置
- GitHub Issue模板
- GitHub PR模板
- 基础CI/CD配置结构

### 项目结构
- 建立标准化目录结构
- 创建tools、scripts、configs、docs等目录
- 设置playbooks和tests占位符

## [0.1.0] - 2025-09-01

### 新增
- 项目初始化
- 基础目录结构
- 初始README文档

---

## 版本说明

### 版本号格式
本项目使用语义化版本号：`主版本号.次版本号.修订号`

- **主版本号**：不兼容的API修改
- **次版本号**：向下兼容的功能性新增
- **修订号**：向下兼容的问题修正

### 变更类型
- **新增**：新功能
- **变更**：既有功能的变更
- **废弃**：即将移除的功能
- **移除**：已移除的功能
- **修复**：任何bug修复
- **安全**：修复安全问题

---

## 计划功能

### v2.0.0 (规划中)
- ClusterReport完整实现
  - Server模式
  - Agent模式
  - Storage层
  - Scheduler调度器
  - REST API
  - Web Dashboard
  - 通知系统

### v1.5.0 (规划中)
- Python工具集
  - log_analyzer.py
  - metrics_collector.py
  - alert_manager.py
  - config_validator.py
- 更多配置模板
- 更多自动化脚本

### v1.2.0 (近期)
- 完善CI/CD流程
- 添加单元测试
- 添加集成测试
- 性能优化
- 文档完善

---

## 贡献

感谢所有为本项目做出贡献的开发者！

如果您想为项目做出贡献，请阅读 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 支持

如有问题或建议，请通过以下方式联系：
- GitHub Issues: https://github.com/sunyifei83/devops-toolkit/issues
- Email: sunyifei83@gmail.com

---

*本更新日志由项目维护者手动更新*
