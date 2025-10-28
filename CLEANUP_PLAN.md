# 项目清理计划

**日期**: 2025/10/28  
**目的**: 清理项目仓库中的无关内容，保持项目整洁

## 🎯 清理目标

根据项目重构计划，识别并清理以下类型的内容：

1. ✅ 旧的文档文件（已被新文档替代）
2. ✅ 技术书籍等非代码资源
3. ✅ 重复的工具目录
4. ✅ 临时文件和系统文件
5. ✅ 过时的审计和分析报告

## 📋 需要清理的内容

### 1. 文档目录清理 (docs/)

#### 删除/移动的文件
- ❌ `docs/BestPractices.md` - 旧版，已被 `docs/best-practices/` 目录替代
- ❌ `docs/ToolsDocumentation.md` - 旧版工具文档
- ❌ `docs/UserGuide.md` - 旧版，已被 `docs/user-guide/` 目录替代
- ❌ `docs/Technical Books/` - 技术书籍，不应在代码仓库中

#### 保留的文件
- ✅ `docs/README.md` - 新的文档索引
- ✅ `docs/getting-started/` - 新的快速开始文档
- ✅ `docs/user-guide/` - 新的用户指南目录
- ✅ `docs/developer-guide/` - 新的开发者指南目录
- ✅ `docs/reference/` - 参考文档目录
- ✅ `docs/best-practices/` - 最佳实践目录
- ✅ `docs/tutorials/` - 教程目录
- ✅ `docs/img/` - 文档图片
- ✅ `docs/tools/` - 工具文档（待整合）

### 2. 工具目录清理 (tools/)

#### 需要移动的文件
- 🔄 `tools/go/NodeProbe.go` → `legacy/NodeProbe.go`
- 🔄 `tools/go/PerfSnap.go` → `legacy/PerfSnap.go`
- 🔄 `tools/go/ClusterReport/` → 需要重组到新结构

#### 保留的内容
- ✅ `tools/go/DocConverter.go` - 独立工具，保留
- ✅ `tools/python/` - Python 工具
- ✅ `tools/shell/` - Shell 工具
- ✅ `tools/go/REFACTORING_PLAN.md` - 重构计划文档
- ✅ `tools/go/INTEGRATION_IMPLEMENTATION.md` - 实施指南

### 3. 顶层目录清理

#### 删除的分析报告（已完成任务）
- ❌ `PROJECT_AUDIT_REPORT.md` - 审计报告（归档）
- ❌ `PROJECT_COMPREHENSIVE_REVIEW.md` - 综合评审（归档）
- ❌ `TOOLS_INTEGRATION_ANALYSIS.md` - 集成分析（归档）
- ❌ `TOOLS_INTEGRATION_STATUS.md` - 集成状态（归档）

这些报告移动到 `docs/archive/` 目录保存

#### 保留的文档
- ✅ `PROJECT_RESTRUCTURE_PLAN.md` - 重构计划（重要参考）
- ✅ `README.md` - 项目主页
- ✅ `README-zh.md` - 中文主页
- ✅ `CHANGELOG.md` - 变更日志
- ✅ `CONTRIBUTING.md` - 贡献指南
- ✅ `CODE_OF_CONDUCT.md` - 行为准则
- ✅ `LICENSE` - 许可证

### 4. 系统文件清理

#### 删除的系统文件
- ❌ `.DS_Store` - macOS 系统文件（多个位置）
- ❌ `tools/go/ClusterReport/clusterreport` - 编译产物

### 5. ClusterReport 目录重组

`tools/go/ClusterReport/` 的内容需要根据新架构重组：

#### 移动到新位置
```
tools/go/ClusterReport/pkg/ → pkg/clusterreport/
tools/go/ClusterReport/cmd/ → cmd/clusterreport/
tools/go/ClusterReport/plugins/ → plugins/
tools/go/ClusterReport/web/ → web/
tools/go/ClusterReport/config.yaml → configs/clusterreport/default.yaml
tools/go/ClusterReport/README.md → 整合到主 README
tools/go/ClusterReport/ROADMAP.md → 保留
tools/go/ClusterReport/go.mod → 顶层 go.mod
```

## 🔧 清理操作

### 步骤 1: 创建归档目录
```bash
mkdir -p docs/archive
mkdir -p legacy
```

### 步骤 2: 移动旧文档到归档
```bash
mv PROJECT_AUDIT_REPORT.md docs/archive/
mv PROJECT_COMPREHENSIVE_REVIEW.md docs/archive/
mv TOOLS_INTEGRATION_ANALYSIS.md docs/archive/
mv TOOLS_INTEGRATION_STATUS.md docs/archive/
mv docs/BestPractices.md docs/archive/
mv docs/ToolsDocumentation.md docs/archive/
mv docs/UserGuide.md docs/archive/
```

### 步骤 3: 删除技术书籍目录
```bash
rm -rf "docs/Technical Books"
```

### 步骤 4: 移动旧工具到 legacy
```bash
mv tools/go/NodeProbe.go legacy/
mv tools/go/PerfSnap.go legacy/
```

### 步骤 5: 创建 legacy/README.md
说明这些是旧版本，提供迁移指南

### 步骤 6: 删除系统文件
```bash
find . -name ".DS_Store" -type f -delete
rm -f tools/go/ClusterReport/clusterreport
```

### 步骤 7: 更新 .gitignore
添加常见的忽略规则

## 📊 清理前后对比

### 清理前
```
devops-toolkit/
├── PROJECT_AUDIT_REPORT.md          ❌ 临时报告
├── PROJECT_COMPREHENSIVE_REVIEW.md  ❌ 临时报告
├── TOOLS_INTEGRATION_ANALYSIS.md    ❌ 临时报告
├── TOOLS_INTEGRATION_STATUS.md      ❌ 临时报告
├── docs/
│   ├── BestPractices.md            ❌ 旧文档
│   ├── ToolsDocumentation.md       ❌ 旧文档
│   ├── UserGuide.md                ❌ 旧文档
│   └── Technical Books/            ❌ 非代码资源
├── tools/go/
│   ├── NodeProbe.go                ❌ 待移动到 legacy
│   ├── PerfSnap.go                 ❌ 待移动到 legacy
│   └── ClusterReport/              ❌ 待重组
└── .DS_Store                        ❌ 系统文件
```

### 清理后
```
devops-toolkit/
├── PROJECT_RESTRUCTURE_PLAN.md     ✅ 保留（重要）
├── docs/
│   ├── README.md                   ✅ 新文档索引
│   ├── getting-started/            ✅ 新结构
│   ├── user-guide/                 ✅ 新结构
│   ├── developer-guide/            ✅ 新结构
│   ├── reference/                  ✅ 新结构
│   └── archive/                    📦 旧文档归档
├── legacy/
│   ├── README.md                   ✅ 说明文档
│   ├── NodeProbe.go                📦 旧版本
│   └── PerfSnap.go                 📦 旧版本
└── [干净的项目结构]
```

## ✅ 清理后的好处

1. **项目更整洁** - 移除无关内容
2. **结构更清晰** - 新旧版本分离
3. **维护更简单** - 减少混淆
4. **仓库更小** - 移除大文件（技术书籍）
5. **导航更容易** - 清晰的文档结构

## 📝 注意事项

1. **不要直接删除** - 先移动到 archive 或 legacy
2. **保留历史记录** - Git 历史中仍可访问
3. **更新文档链接** - 确保所有链接指向新位置
4. **通知团队** - 告知结构变更

## 🎯 执行清理

执行命令后验证：
- 检查新的目录结构
- 验证文档链接
- 测试构建过程
- 更新相关文档

---

**准备执行清理操作**
