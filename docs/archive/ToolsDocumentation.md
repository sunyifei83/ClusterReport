# DevOps Toolkit - 工具文档总览

## 概述

本文档提供DevOps Toolkit中所有工具的索引和简要说明。每个工具都有详细的独立文档，可通过链接访问。

## 工具分类

### Go工具

#### 1. [NodeProbe](tools/go/NodeProbe.md) - 分布式节点探测工具
- **版本**: 1.0.0
- **用途**: 批量检测和监控分布式节点的健康状态
- **特性**: 
  - 支持HTTP/HTTPS、TCP、Ping多种探测方式
  - 并发探测，实时监控
  - 灵活的配置和报警机制
- **使用场景**: 集群健康检查、服务可用性监控、网络连通性测试

#### 2. [PerfSnap](tools/go/PerfSnap.md) - Linux系统性能快照分析工具
- **版本**: 1.1.0
- **用途**: 快速采集和分析系统性能指标
- **特性**:
  - CPU、内存、磁盘IO、网络全方位监控
  - 火焰图生成功能
  - 智能性能问题诊断
- **使用场景**: 性能基线建立、故障诊断、性能调优

#### 3. [DocConverter](tools/go/DocConverter.md) - 文档转PDF工具
- **版本**: 1.1.0
- **用途**: 将文档和网页转换为PDF格式
- **特性**:
  - 支持Markdown、HTML文件转换
  - 网页内容爬取和转换
  - 批量处理和样式定制
- **使用场景**: 技术文档归档、报告生成、离线阅读

### Shell工具

#### 1. [iotest](tools/shell/iotest.md) - 磁盘IO性能测试工具
- **版本**: 1.0.0
- **用途**: 测试磁盘读写性能
- **特性**:
  - 顺序/随机读写测试
  - 多线程并发测试
  - 详细的性能报告
- **使用场景**: 存储性能评估、基准测试、故障排查

### Python工具

*待添加*

## 快速开始

### 安装要求

#### 系统要求
- Linux/macOS/Windows (部分工具仅支持Linux)
- Go 1.15+ (用于编译Go工具)
- Python 3.6+ (用于Python工具)
- Bash 4.0+ (用于Shell脚本)

#### 通用依赖
```bash
# Go工具编译
go version  # 确认Go已安装

# 系统工具
sudo apt-get install -y git curl wget  # Ubuntu/Debian
sudo yum install -y git curl wget      # CentOS/RHEL
brew install git curl wget             # macOS
```

### 快速安装

```bash
# 克隆项目
git clone https://github.com/sunyifei83/devops-toolkit.git
cd devops-toolkit

# 编译所有Go工具
cd tools/go
for tool in *.go; do
    go build -o "${tool%.go}" "$tool"
done

# 设置Shell脚本权限
cd ../shell
chmod +x *.sh
```

## 工具对比

| 工具名称 | 语言 | 主要功能 | 适用场景 | 平台支持 |
|---------|------|---------|---------|----------|
| NodeProbe | Go | 节点健康检测 | 分布式系统监控 | 全平台 |
| PerfSnap | Go | 性能快照分析 | Linux性能诊断 | Linux |
| DocConverter | Go | 文档转PDF | 文档处理 | 全平台 |
| iotest | Shell | IO性能测试 | 存储评测 | Linux/macOS |

## 使用建议

### 监控场景
- **实时监控**: 使用NodeProbe进行服务健康检查
- **性能分析**: 使用PerfSnap定期收集性能数据
- **存储测试**: 使用iotest评估磁盘性能

### 文档场景
- **文档归档**: 使用DocConverter将Markdown文档转为PDF
- **网站备份**: 使用DocConverter爬取并保存网页内容

### 最佳实践

1. **定期更新工具**
   ```bash
   cd devops-toolkit
   git pull origin main
   # 重新编译工具
   ```

2. **配置别名**
   ```bash
   # 添加到 ~/.bashrc 或 ~/.zshrc
   export PATH=$PATH:/path/to/devops-toolkit/tools/go
   export PATH=$PATH:/path/to/devops-toolkit/tools/shell
   ```

3. **集成到CI/CD**
   - NodeProbe: 集成到健康检查流程
   - PerfSnap: 定期性能报告生成
   - DocConverter: 自动化文档发布

## 工具矩阵

### 按功能分类

#### 监控类
- NodeProbe - 服务可用性监控
- PerfSnap - 系统性能监控

#### 测试类
- iotest - 磁盘性能测试

#### 文档类
- DocConverter - 文档格式转换

### 按使用频率

#### 日常使用
- NodeProbe (服务监控)
- DocConverter (文档处理)

#### 定期使用
- PerfSnap (性能巡检)
- iotest (性能基线)

#### 故障时使用
- PerfSnap (性能诊断)
- iotest (存储问题排查)

## 版本历史

### 最新更新
- **2025-09-16**: DocConverter v1.1.0 - 新增网页爬取功能
- **2025-09-15**: PerfSnap v1.1.0 - 新增火焰图生成
- **2025-09-14**: NodeProbe v1.0.0 - 首次发布

## 贡献指南

欢迎贡献新工具或改进现有工具！

### 添加新工具
1. 在相应语言目录下创建工具
2. 编写详细的工具文档
3. 更新本文档索引
4. 提交Pull Request

### 文档规范
- 每个工具必须有独立的Markdown文档
- 包含：概述、特性、安装、使用、示例、故障排查
- 提供实际使用案例

## 支持与反馈

- **项目地址**: https://github.com/sunyifei83/devops-toolkit
- **问题反馈**: 通过GitHub Issues
- **邮件联系**: sunyifei83@gmail.com

## 许可证

本项目采用MIT许可证，详见[LICENSE](../LICENSE)文件。

## 相关链接

- [用户指南](UserGuide.md)
- [最佳实践](BestPractices.md)
- [项目README](../README.md)

---

*本文档持续更新中，最后更新时间：2025-09-16*
