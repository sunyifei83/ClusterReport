# clear_log.sh - 审计日志清理工具

## 概述

`clear_log.sh` 是一个用于自动清理系统审计日志的 Bash 脚本。该脚本智能地管理日志文件，在删除过期日志的同时确保每个日志目录至少保留一个文件，防止日志目录完全清空导致的服务异常。

## 功能特性

### 主要功能
- **自动扫描审计日志目录**：遍历 `/home/service/*/_package/run/auditlog/*` 路径下的所有日志目录
- **智能清理策略**：删除超过1天的过期日志文件
- **安全保护机制**：确保每个目录至少保留一个日志文件
- **详细执行报告**：输出每个目录的处理情况和统计信息

### 清理策略

脚本采用以下智能清理策略：

1. **单文件保护**：如果目录中只有一个文件，无论是否过期都不会删除
2. **有新文件时**：如果存在未过期的新文件，则删除所有过期文件
3. **仅有过期文件时**：保留最新的一个过期文件，删除其余的过期文件

## 使用方法

### 基本用法

```bash
# 添加执行权限
chmod +x clear_log.sh

# 执行脚本
./clear_log.sh
```

### 定时任务配置

建议通过 crontab 设置定期执行：

```bash
# 编辑 crontab
crontab -e

# 每天凌晨2点执行清理
0 2 * * * /path/to/clear_log.sh >> /var/log/clear_log.log 2>&1

# 每周一凌晨3点执行清理
0 3 * * 1 /path/to/clear_log.sh >> /var/log/clear_log.log 2>&1
```

## 工作原理

### 执行流程

1. **目录扫描**
   - 扫描路径模式：`/home/service/*/_package/run/auditlog/*`
   - 验证每个匹配路径是否为有效目录

2. **文件统计**
   - 统计目录中的总文件数
   - 统计超过1天的过期文件数

3. **清理决策**
   - 根据文件数量和过期情况决定清理策略
   - 执行相应的清理操作

4. **结果报告**
   - 输出每个目录的处理详情
   - 显示清理前后的文件统计

### 核心逻辑解析

#### 文件统计
```bash
# 获取总文件数
total_files=$(find "$dir" -maxdepth 1 -type f | wc -l)

# 获取过期文件数（超过1天）
old_files=$(find "$dir" -maxdepth 1 -type f -mtime +1 | wc -l)
```

#### 智能保留策略
```bash
# 保留最新的过期文件，删除最旧的
find "$dir" -maxdepth 1 -type f -mtime +1 -printf '%T@ %p\n' | \
sort -n | \
head -n -1 | \
cut -d' ' -f2- | \
xargs -r rm -f
```

## 输出示例

```
处理目录: /home/service/app1/_package/run/auditlog/2024
  总文件数: 15
  过期文件数: 10
  存在 5 个新文件，删除所有过期文件
  清理完成

处理目录: /home/service/app2/_package/run/auditlog/2024
  总文件数: 1
  过期文件数: 1
  目录中只有 1 个文件，保留不删除

处理目录: /home/service/app3/_package/run/auditlog/2024
  总文件数: 8
  过期文件数: 8
  没有新文件，将从过期文件中保留最新的一个
  清理完成

所有目录处理完成
```

## 注意事项

### 使用前提
- 需要对目标日志目录有读写权限
- 建议在非业务高峰期执行
- 首次使用建议先在测试环境验证

### 安全建议
1. **备份重要日志**：在首次运行前，建议备份重要的审计日志
2. **权限控制**：确保脚本文件权限适当，避免未授权访问
3. **日志监控**：建议将脚本执行日志重定向到文件，便于追踪和审计
4. **测试验证**：修改日志路径前，先使用 `echo` 替代 `rm` 命令进行测试

### 性能考虑
- 对于包含大量文件的目录，执行时间可能较长
- 建议在系统负载较低时执行
- 可考虑使用 `nice` 命令降低进程优先级

## 自定义配置

### 修改日志路径

如需清理其他路径的日志，修改脚本开头的路径变量：

```bash
# 原始路径
AUDIT_LOG_DIR="/home/service/*/_package/run/auditlog/*"

# 自定义路径示例
AUDIT_LOG_DIR="/var/log/myapp/*/audit"
```

### 调整保留时间

修改 `find` 命令中的 `-mtime` 参数：

```bash
# 保留7天内的文件
find "$dir" -maxdepth 1 -type f -mtime +7

# 保留30天内的文件
find "$dir" -maxdepth 1 -type f -mtime +30
```

### 添加文件类型过滤

如果只想清理特定类型的文件：

```bash
# 只清理 .log 文件
find "$dir" -maxdepth 1 -type f -name "*.log" -mtime +1

# 清理多种类型
find "$dir" -maxdepth 1 -type f \( -name "*.log" -o -name "*.txt" \) -mtime +1
```

## 故障排查

### 常见问题

1. **权限不足**
   ```
   rm: cannot remove 'file': Permission denied
   ```
   解决方案：确保脚本以适当权限运行，或使用 sudo

2. **路径不存在**
   ```
   处理目录: /home/service/*/...
   ```
   解决方案：检查路径模式是否正确，确保目标目录存在

3. **没有文件被清理**
   - 检查文件的修改时间是否真的超过设定天数
   - 验证 find 命令的参数是否正确

### 调试模式

在脚本开头添加调试选项：

```bash
#!/bin/bash
set -x  # 开启调试模式
set -e  # 遇到错误立即退出
```

## 扩展功能建议

### 1. 添加日志压缩
在删除前先压缩归档：
```bash
# 压缩7天前的日志
find "$dir" -name "*.log" -mtime +7 -exec gzip {} \;
```

### 2. 发送清理报告
通过邮件发送清理统计：
```bash
# 将输出重定向并发送邮件
./clear_log.sh | mail -s "Log Cleanup Report" admin@example.com
```

### 3. 添加配置文件
创建配置文件管理不同的清理策略：
```bash
# config.conf
LOG_DIRS="/home/service/*/_package/run/auditlog/*"
RETENTION_DAYS=1
MIN_FILES_TO_KEEP=1
```

## 版本历史

- v1.0: 初始版本，基本的日志清理功能
- 支持智能保留策略
- 添加详细的执行报告

## 相关文档

- [iotest.sh - I/O 性能测试工具](./iotest.md)
- [Shell 脚本最佳实践](../../BestPractices.md)

## 许可证

MIT License

## 作者

DevOps Toolkit Team

---

*最后更新：2024年*
