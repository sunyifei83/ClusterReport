# 深度清理计划 - ClusterReport 为核心

**日期**: 2025/10/28  
**目的**: 深度清理与 ClusterReport 无关的内容，保持项目纯粹性

## 🎯 清理原则

**核心原则**: 只保留与 ClusterReport 平台直接相关的内容

- ✅ **保留**: ClusterReport 核心功能、文档、配置
- ✅ **保留**: 辅助 ClusterReport 的工具（DocConverter 等）
- ❌ **移除**: 与 ClusterReport 无关的独立功能
- ❌ **移除**: 纯示例性质的空目录

## 📋 需要清理的目录和文件

### 1. scripts/ 目录分析

#### 当前状态
```
scripts/
├── automation/
│   ├── backup/backup_databases.sh    ❌ 数据库备份脚本（与 ClusterReport 无关）
│   ├── cleanup/                      📁 空目录
│   └── deployment/                   📁 空目录
├── cloud/
│   └── qiniu/                        📁 七牛云相关（与 ClusterReport 无关）
├── monitoring/
│   ├── alerting/                     📁 空目录
│   ├── grafana/                      📁 空目录
│   └── prometheus/                   📁 空目录
├── security/
│   ├── audit/                        📁 空目录
│   ├── compliance/                   📁 空目录
│   └── scanning/                     📁 空目录
├── troubleshooting/
│   ├── logs/                         📁 空目录
│   ├── network/网络性能调优.pdf      ❌ 文档（应在 docs/）
│   └── performance/                  📁 空目录
└── installation/                      ✅ 保留（ClusterReport 安装）
```

#### 清理决策
- ❌ **删除**: `scripts/automation/` - 与 ClusterReport 无关的通用自动化脚本
- ❌ **删除**: `scripts/cloud/qiniu/` - 特定云服务相关
- ❌ **删除**: `scripts/monitoring/` - 空目录，监控功能应在 ClusterReport 内
- ❌ **删除**: `scripts/security/` - 空目录
- ❌ **移动**: `scripts/troubleshooting/network/网络性能调优.pdf` → `docs/archive/`
- ❌ **删除**: `scripts/troubleshooting/` - 空目录
- ✅ **保留**: `scripts/installation/` - ClusterReport 安装脚本

### 2. tests/ 目录分析

#### 当前状态
```
tests/
├── cosbench/      ❌ COSBench 性能测试（存储测试，与 ClusterReport 无关）
├── io500/         ❌ IO500 测试（存储测试，与 ClusterReport 无关）
├── unit/          ✅ 单元测试（空，待用）
├── integration/   ✅ 集成测试（空，待用）
├── e2e/           ✅ 端到端测试（空，待用）
└── fixtures/      ✅ 测试数据（空，待用）
```

#### 清理决策
- ❌ **移动**: `tests/cosbench/` → `legacy/tests/cosbench/` - 存储性能测试
- ❌ **移动**: `tests/io500/` → `legacy/tests/io500/` - 存储性能测试
- ✅ **保留**: `tests/unit/`, `tests/integration/`, `tests/e2e/`, `tests/fixtures/` - ClusterReport 测试框架

### 3. playbooks/ 目录分析

#### 当前状态
```
playbooks/
├── setup/
│   ├── docker-setup.yml      ❓ Docker 环境配置（通用）
│   └── web-server-setup.yml  ❓ Web 服务器配置（通用）
└── README.md
```

#### 清理决策
这些 Ansible playbooks 需要评估：
- 如果是为 ClusterReport 部署准备环境 → ✅ 保留
- 如果是通用环境配置 → ❌ 移除或移到 examples/

**建议**: 移动到 `deployments/ansible/` 并重命名为 ClusterReport 特定的部署脚本

### 4. configs/ 目录分析

#### 当前状态
```
configs/
├── clusterreport/    ✅ ClusterReport 配置
├── docker/           ❓ 通用 Docker 配置
├── kubernetes/       ❓ 通用 K8s 配置
├── nginx/            ❓ 通用 Nginx 配置
└── terraform/        ❓ 通用 Terraform 配置
```

#### 清理决策
需要区分：
- 如果是 ClusterReport 部署相关 → ✅ 保留在 `deployments/`
- 如果是通用配置示例 → ❌ 移到 `examples/` 或删除

### 5. tools/ 目录分析

#### 当前状态
```
tools/
├── go/
│   ├── DocConverter.go           ✅ 保留（独立辅助工具）
│   ├── ClusterReport/            🔄 需要重组
│   ├── REFACTORING_PLAN.md       ✅ 保留（文档）
│   └── INTEGRATION_IMPLEMENTATION.md  ✅ 保留（文档）
├── python/
│   ├── metrics_collector.py      ❓ 是否用于 ClusterReport？
│   └── log_analyzer.py           ❓ 是否用于 ClusterReport？
└── shell/
    ├── clear_log.sh              ❓ 通用脚本
    └── iotest.sh                 ❓ IO 测试脚本
```

#### 清理决策
- ✅ **保留**: `tools/go/DocConverter.go` - 独立工具
- 🔄 **重组**: `tools/go/ClusterReport/` - 移动到顶层目录
- ❓ **评估**: Python 和 Shell 工具是否为 ClusterReport 插件

## 🔧 清理操作步骤

### 步骤 1: 清理 scripts/ 目录

```bash
# 移动 PDF 到归档
mkdir -p docs/archive/resources
mv scripts/troubleshooting/network/网络性能调优.pdf docs/archive/resources/

# 删除与 ClusterReport 无关的脚本
rm -rf scripts/automation
rm -rf scripts/cloud
rm -rf scripts/monitoring
rm -rf scripts/security
rm -rf scripts/troubleshooting

# 只保留 installation 目录
# scripts/
# └── installation/
```

### 步骤 2: 清理 tests/ 目录

```bash
# 移动存储测试到 legacy
mkdir -p legacy/tests
mv tests/cosbench legacy/tests/
mv tests/io500 legacy/tests/

# 保留测试框架目录
# tests/
# ├── unit/
# ├── integration/
# ├── e2e/
# └── fixtures/
```

### 步骤 3: 重组 playbooks/

```bash
# 移动到 deployments
mv playbooks/setup/*.yml deployments/ansible/
mv playbooks/README.md deployments/ansible/

# 删除 playbooks 目录
rm -rf playbooks
```

### 步骤 4: 重组 configs/

```bash
# ClusterReport 配置保留
# configs/clusterreport/ ✅

# 其他配置移到 deployments 或 examples
mv configs/docker deployments/
mv configs/kubernetes deployments/
mv configs/nginx examples/nginx-config
mv configs/terraform examples/terraform-config

# 或者全部移到 deployments
# 决策：ClusterReport 部署相关 → deployments
#       纯示例 → examples
```

### 步骤 5: 清理 tools/

```bash
# 评估 Python 工具
# 如果不是 ClusterReport 插件，移到 legacy
# mv tools/python/* legacy/tools/python/

# 评估 Shell 工具
# 如果是通用工具，移到 legacy
# mv tools/shell/* legacy/tools/shell/
```

## 📊 预期清理效果

### 清理前
```
devops-toolkit/
├── scripts/           # 7个子目录，大部分为空或无关
├── tests/             # 包含存储测试工具
├── playbooks/         # 通用 Ansible playbooks
├── configs/           # 混合配置
└── tools/             # 混合工具
```

### 清理后
```
devops-toolkit/
├── scripts/
│   └── installation/  # 只保留 ClusterReport 安装脚本
├── tests/
│   ├── unit/          # ClusterReport 单元测试
│   ├── integration/   # ClusterReport 集成测试
│   ├── e2e/           # ClusterReport E2E 测试
│   └── fixtures/      # 测试数据
├── deployments/
│   ├── docker/        # ClusterReport Docker 部署
│   ├── kubernetes/    # ClusterReport K8s 部署
│   └── ansible/       # ClusterReport Ansible 部署
├── configs/
│   └── clusterreport/ # ClusterReport 配置
├── examples/          # 配置示例
└── tools/
    ├── go/
    │   └── DocConverter.go  # 辅助工具
    ├── python/        # 如果是插件则保留
    └── shell/         # 如果是插件则保留
```

## ✅ 清理原则总结

1. **与 ClusterReport 无关** → 移到 `legacy/` 或删除
2. **空目录** → 删除
3. **通用示例** → 移到 `examples/`
4. **部署相关** → 整合到 `deployments/`
5. **测试框架** → 保留但清理无关测试
6. **文档资源** → 移到 `docs/archive/`

## 🎯 清理目标

完成后项目应该：
- ✅ 100% 聚焦 ClusterReport 平台
- ✅ 无空目录和无关内容
- ✅ 清晰的目录职责
- ✅ 专业的项目结构

---

**准备执行深度清理**
