# Legacy - 旧版工具和资源

本目录包含旧版本的独立工具、测试和配置，已被 ClusterReport 平台集成或不再维护。

## 📦 目录结构

```
legacy/
├── NodeProbe.go          # 旧版节点探测工具
├── PerfSnap.go           # 旧版性能快照工具
├── tests/                # 旧版测试工具
│   ├── cosbench/         # COSBench 存储测试
│   └── io500/            # IO500 存储测试
└── tools/                # 旧版工具脚本
    ├── python/           # Python 工具
    └── shell/            # Shell 脚本
```

## 🔧 旧版工具说明

### NodeProbe.go
**状态**: ⚠️ 已废弃，已集成到 ClusterReport

**旧版使用方式**:
```bash
go run NodeProbe.go
```

**新版使用方式**:
```bash
# NodeProbe 功能已内置到 ClusterReport
clusterreport collect --node localhost
```

NodeProbe 现在作为 ClusterReport 的内置引擎，提供系统配置信息采集功能。

### PerfSnap.go
**状态**: ⚠️ 已废弃，已集成到 ClusterReport

**旧版使用方式**:
```bash
go run PerfSnap.go -flame
```

**新版使用方式**:
```bash
# PerfSnap 功能已内置到 ClusterReport
clusterreport collect --node localhost --enable-flamegraph
```

PerfSnap 现在作为 ClusterReport 的内置引擎，提供性能数据采集和分析功能。

### tests/ - 存储测试工具
**状态**: ⚠️ 已废弃，与 ClusterReport 无关

#### COSBench
对象存储性能测试工具，用于测试 S3 兼容存储系统。

#### IO500
存储 I/O 性能基准测试套件。

**说明**: 这些是独立的存储测试工具，不属于 ClusterReport 的功能范围。如需使用，请参考各自的官方文档。

### tools/python/ - Python 工具
**状态**: ⚠️ 已废弃，通用工具

#### metrics_collector.py
系统指标采集工具，收集 CPU、内存、磁盘、网络等性能指标。

**旧版使用方式**:
```bash
python3 metrics_collector.py --continuous --interval 60
```

**新版使用方式**:
```bash
# ClusterReport 提供了更强大的指标采集功能
clusterreport collect --node localhost --interval 60
```

#### log_analyzer.py
日志分析工具，解析和分析各种格式的日志文件。

**旧版使用方式**:
```bash
python3 log_analyzer.py -f /var/log/app.log -l ERROR
```

**新版使用方式**:
```bash
# ClusterReport 提供了日志分析插件
clusterreport analyze logs --file /var/log/app.log --level ERROR
```

### tools/shell/ - Shell 脚本
**状态**: ⚠️ 已废弃，通用脚本

#### clear_log.sh
审计日志清理脚本，自动清理过期的审计日志文件。

#### iotest.sh
IO 性能测试脚本（基于 YABS），测试磁盘读写性能。

**说明**: 这些是通用的系统管理脚本，不是 ClusterReport 特定功能。

## 🔄 迁移指南

### 从 NodeProbe 迁移

如果您之前使用 NodeProbe：

```bash
# 旧方式
go run NodeProbe.go -format json > config.json

# 新方式
clusterreport collect --node localhost
clusterreport generate --format json --output config.json
```

### 从 PerfSnap 迁移

如果您之前使用 PerfSnap：

```bash
# 旧方式
go run PerfSnap.go -flame -output perf.svg

# 新方式
clusterreport collect --node localhost --enable-flamegraph
clusterreport generate --format html --output perf-report.html
```

## 📚 新架构优势

使用 ClusterReport 集成版本的优势：

1. **统一平台** - 一个工具完成所有任务
2. **更强大** - 结合了 NodeProbe 和 PerfSnap 的功能
3. **更易用** - 一致的命令行接口
4. **更多功能** - 插件系统、Web 界面、定时任务等
5. **持续维护** - 新功能和 bug 修复都在 ClusterReport 中

## 🚀 开始使用新版本

1. **安装 ClusterReport**:
   ```bash
   # 使用安装脚本
   curl -sSL https://raw.githubusercontent.com/sunyifei83/devops-toolkit/main/scripts/installation/install.sh | bash
   
   # 或从源码编译
   cd /path/to/devops-toolkit
   make install
   ```

2. **快速开始**:
   ```bash
   # 采集数据（包含 NodeProbe + PerfSnap 功能）
   clusterreport collect --node localhost
   
   # 生成报告
   clusterreport generate --format html --output report.html
   
   # 查看报告
   open report.html
   ```

3. **查看文档**:
   - [快速入门](../docs/getting-started/quick-start.md)
   - [用户指南](../docs/user-guide/)
   - [CLI 参考](../docs/reference/cli-reference.md)

## ⚠️ 注意事项

- 这些旧版本工具不再维护
- 建议所有用户迁移到 ClusterReport
- 如有迁移问题，请提交 [Issue](https://github.com/sunyifei83/devops-toolkit/issues)

## 🔗 相关链接

- [项目主页](../README.md)
- [ClusterReport 文档](../docs/README.md)
- [迁移FAQ](../docs/user-guide/migration-faq.md)

---

**最后更新**: 2025/10/28  
**迁移建议**: 请尽快迁移到 ClusterReport 平台
