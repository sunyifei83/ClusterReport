# Scripts 目录说明

本目录包含各种自动化脚本和工具，用于监控、自动化、故障排查、安全和云平台管理。

## 📁 目录结构

```
scripts/
├── monitoring/          # 监控相关脚本
├── automation/          # 自动化脚本
├── troubleshooting/     # 故障排查工具
├── security/            # 安全相关脚本
└── cloud/               # 云平台管理脚本
```

## 🚧 开发状态

当前此目录处于规划阶段，各子目录将逐步添加实用脚本。

### 计划添加的内容

#### monitoring/ - 监控脚本
- **prometheus/** - Prometheus 相关
  - [ ] node_exporter 安装脚本
  - [ ] prometheus 配置模板
  - [ ] 告警规则示例
  
- **grafana/** - Grafana 相关
  - [ ] grafana 安装脚本
  - [ ] dashboard 模板
  - [ ] 数据源配置
  
- **alerting/** - 告警管理
  - [ ] alertmanager 配置
  - [ ] 通知模板
  - [ ] 告警路由规则

#### automation/ - 自动化脚本
- **deployment/** - 部署自动化
  - [ ] 应用部署脚本
  - [ ] 滚动更新脚本
  - [ ] 蓝绿部署脚本
  
- **backup/** - 备份自动化
  - [ ] 数据库备份脚本
  - [ ] 配置文件备份
  - [ ] 增量备份工具
  
- **cleanup/** - 清理脚本
  - [ ] 日志清理
  - [ ] 临时文件清理
  - [ ] Docker 镜像清理

#### troubleshooting/ - 故障排查
- **network/** - 网络诊断
  - [x] 网络性能调优文档
  - [ ] 网络连接测试脚本
  - [ ] 端口扫描工具
  - [ ] DNS 诊断脚本
  
- **performance/** - 性能分析
  - [ ] CPU 分析脚本
  - [ ] 内存分析工具
  - [ ] IO 性能诊断
  
- **logs/** - 日志分析
  - [ ] 日志聚合脚本
  - [ ] 错误日志提取
  - [ ] 日志统计分析

#### security/ - 安全工具
- **audit/** - 安全审计
  - [ ] 系统审计脚本
  - [ ] 用户权限检查
  - [ ] 配置合规检查
  
- **compliance/** - 合规检查
  - [ ] CIS 基准检查
  - [ ] 安全配置验证
  - [ ] 漏洞扫描报告
  
- **scanning/** - 安全扫描
  - [ ] 端口扫描工具
  - [ ] 漏洞扫描脚本
  - [ ] 恶意软件检测

#### cloud/ - 云平台管理
- **qiniu/** - 七牛云
  - [ ] 对象存储管理
  - [ ] CDN 配置工具
  
- **aws/** - AWS (规划中)
  - [ ] EC2 管理脚本
  - [ ] S3 备份工具
  
- **gcp/** - GCP (规划中)
  - [ ] GCE 管理工具
  - [ ] Cloud Storage 操作

## 🤝 贡献

如果您有实用的脚本想要分享，欢迎提交 PR！

贡献时请确保：
1. 脚本有清晰的注释
2. 包含使用说明
3. 遵循项目代码规范
4. 添加适当的错误处理

## 📝 脚本编写规范

### Shell 脚本
```bash
#!/bin/bash
#
# 脚本名称和用途
# 作者：xxx
# 版本：x.x.x
# 最后更新：YYYY-MM-DD

set -euo pipefail  # 错误时退出，未定义变量报错

# 函数定义
function main() {
    # 主逻辑
}

# 执行主函数
main "$@"
```

### Python 脚本
```python
#!/usr/bin/env python3
"""
脚本名称和用途

作者：xxx
版本：x.x.x
最后更新：YYYY-MM-DD
"""

import sys
import argparse

def main():
    """主函数"""
    pass

if __name__ == "__main__":
    main()
```

## 📚 相关文档

- [用户指南](/docs/UserGuide.md)
- [最佳实践](/docs/BestPractices.md)
- [工具文档](/docs/ToolsDocumentation.md)

## 📮 反馈

如有问题或建议，请通过以下方式联系：
- GitHub Issues: https://github.com/sunyifei83/devops-toolkit/issues
- Email: sunyifei83@gmail.com

---

*此文档会随着项目发展持续更新*
