# ClusterReport - Enterprise Cluster Management and Reporting Platform

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

> 🚀 **Comprehensive Cluster Analysis and Report Generation Tool**  
> An integrated monitoring solution combining NodeProbe and PerfSnap engines

[English](README.md) | [中文](README-zh.md)

## ✨ Core Features

- **📊 Automated Data Collection** - Automatic collection of cluster node configuration and performance data
- **🔍 Deep System Analysis** - Comprehensive analysis of CPU, memory, disk, network, and more
- **⚡ Real-time Performance Monitoring** - Integrated PerfSnap engine for real-time monitoring capabilities
- **🔌 Extensible Plugin System** - Support for custom collectors and analyzers
- **📄 Multi-format Reports** - Generate reports in HTML, JSON, Markdown, and more
- **🤖 Intelligent Analysis Engine** - Automatic health scoring and issue detection

## 🚀 Quick Start

### Building from Source

```bash
# Clone the repository
git clone https://github.com/sunyifei83/ClusterReport.git
cd ClusterReport

# Build the project
go build -o clusterreport ./cmd/cli

# Run the application
./clusterreport --help
```

### Basic Usage

```bash
# Check version
./clusterreport version

# Collect data (in development)
./clusterreport collect --nodes localhost

# Generate report (in development)
./clusterreport generate --output report.html

# One-command report generation (planned)
./clusterreport report --cluster production --formats html,pdf
```

## 📁 Project Structure

```
ClusterReport/                      # Project root directory
├── README.md                       # Project homepage (English)
├── README-zh.md                    # Chinese documentation
├── LICENSE                         # MIT License
├── CODE_OF_CONDUCT.md             # Code of conduct
├── CONTRIBUTING.md                 # Contribution guidelines
├── config.yaml                     # Configuration file
├── go.mod / go.sum                 # Go dependencies
│
├── cmd/                            # Command line entry points
│   ├── cli/                       # CLI mode
│   ├── server/                    # Server mode (planned)
│   └── agent/                     # Agent mode (planned)
│
├── pkg/                            # Core packages
│   ├── collector/                 # Data collectors
│   ├── analyzer/                  # Data analyzers
│   └── generator/                 # Report generators
│
├── plugins/                        # Plugin system
│   ├── collectors/                # Collection plugins (MySQL, Redis, etc.)
│   └── analyzers/                 # Analysis plugins
│
├── web/                            # Web interface (in development)
│   └── dashboard/                 # Management dashboard
│
├── deployments/                    # Deployment configurations
│   ├── docker/                    # Docker configs
│   ├── kubernetes/                # Kubernetes configs
│   └── ansible/                   # Ansible playbooks
│
├── docs/                           # Project documentation
│   ├── getting-started/           # Quick start guides
│   ├── tools/go/                  # Tool documentation
│   └── archive/                   # Archived documentation
│
├── tools/                          # Utility tools
│   └── utils/                     # Utility scripts
│       └── DocConverter.go        # Document converter
│
├── scripts/                        # Script tools
│   └── README.md                  # Scripts documentation
│
└── legacy/                         # Legacy tools
    ├── NodeProbe.go               # Legacy node probe
    ├── PerfSnap.go                # Legacy performance snapshot
    └── tools/                     # Legacy tool scripts
```

## 📚 Documentation

- [Quick Start Guide](docs/getting-started/quick-start.md)
- [Architecture Design](docs/tools/go/ClusterReport_Architecture.md)
- [Detailed Design](docs/tools/go/ClusterReport_Design.md)
- [Legacy Tools Migration Guide](legacy/README.md)

## 🚦 Project Status

**Current Version**: v0.7.0 (70% Complete) 🚧  
**Target Version**: v1.0.0  
**Expected Release**: December 2025

### Completed ✅
- ✅ Project architecture design
- ✅ Core code framework
- ✅ Data collectors (collector package)
- ✅ Data analyzers (analyzer package)
- ✅ Report generators (generator package)
- ✅ Plugin system foundation
- ✅ Web Dashboard UI prototype

### In Development 🚧
- 🚧 CLI command-line tools
- 🚧 Configuration file management
- 🚧 Remote node collection (SSH)
- 🚧 Complete report formats
- 🚧 Test coverage

### Planned 📋
- 📋 Server/Agent architecture
- 📋 Data persistence and storage
- 📋 Scheduled task system
- 📋 Alert and notification system
- 📋 Docker/Kubernetes deployment

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                   ClusterReport Platform                     │
├─────────────────────────────────────────────────────────────┤
│  📦 Built-in Data Collection Engines                        │
│  • NodeProbe Engine - System configuration collection       │
│  • PerfSnap Engine - Performance data collection            │
│  • Plugin System - Custom collectors                        │
├─────────────────────────────────────────────────────────────┤
│  📊 Collector → 📈 Analyzer → 📝 Generator                  │
│  Data Collection   Smart Analysis   Report Generation       │
├─────────────────────────────────────────────────────────────┤
│  🔌 Plugin System                                            │
│  • MySQL/Redis/Custom Collectors                            │
│  • Anomaly Detection/Trend Analysis                         │
│  • HTML/PDF/Excel Output                                    │
└─────────────────────────────────────────────────────────────┘
```

### Core Components

1. **Data Collection Layer**
   - System metrics collection (CPU, Memory, Disk, Network)
   - Performance data acquisition
   - Plugin-based extensibility

2. **Analysis Layer**
   - Multi-dimensional metric analysis
   - Intelligent health scoring (0-100)
   - Automated issue detection

3. **Report Generation Layer**
   - Multiple output formats (HTML, JSON, Markdown, PDF)
   - Customizable templates
   - Rich visualization and charts

4. **Plugin System**
   - Custom collector support
   - Third-party integrations (MySQL, Redis, Prometheus)
   - Extensible analyzer framework

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Configuration**: YAML
- **CLI Framework**: Cobra + Viper
- **Web**: HTML + CSS + JavaScript
- **Report Formats**: HTML, JSON, Markdown, PDF (planned)
- **Deployment**: Docker, Kubernetes (planned)

## 🎯 Use Cases

### 1. New Cluster Acceptance
```bash
# Collect data from all nodes
clusterreport collect --cluster production --nodes node1,node2,node3

# Generate acceptance report
clusterreport generate --type acceptance --baseline baseline.yaml --output acceptance_report.html
```

### 2. Regular Inspections
```bash
# Monthly inspection report
clusterreport report --cluster production --formats html,pdf --output ./reports/
```

### 3. Performance Benchmarking
```bash
# Performance baseline testing
clusterreport collect --nodes all --include-benchmarks
clusterreport analyze --input collected_data.json --type performance
clusterreport generate --format html --output performance_report.html
```

### 4. Incident Analysis
```bash
# Analyze specific time period
clusterreport analyze --time-range "2025-01-15 14:00,2025-01-15 16:00" \
  --focus performance,logs --output incident_report.pdf
```

## 🤝 Contributing

We welcome contributions, bug reports, and feature suggestions!

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## 📄 License

This project is licensed under the [MIT License](LICENSE)

## 📮 Contact

- **GitHub**: https://github.com/sunyifei83/ClusterReport
- **Issues**: https://github.com/sunyifei83/ClusterReport/issues
- **Email**: sunyifei83@gmail.com

## 🎯 Design Philosophy

1. **Unified Platform** - Integrate NodeProbe and PerfSnap functionality
2. **Simple to Use** - Accomplish complex tasks with a single command
3. **Highly Extensible** - Plugin system for custom extensions
4. **Lightweight Deployment** - Single binary, no external dependencies
5. **Cloud-Native Ready** - Container deployment with native Kubernetes support

## 🗺️ Roadmap

### Phase 1: CLI Foundation (v0.8.0) - Q2 2025
- Complete CLI command implementation
- Configuration file management
- Local data collection and analysis
- Basic report generation (HTML, Markdown, JSON)

### Phase 2: Remote Collection (v0.9.0) - Q3 2025
- SSH-based remote node collection
- Batch processing of multiple nodes
- Enhanced data analysis
- Advanced reporting features

### Phase 3: Server/Agent Architecture (v1.0.0) - Q4 2025
- REST API server
- Web management dashboard
- Agent-based data collection
- Scheduled task system
- User authentication and authorization

### Phase 4: Enterprise Features (v2.0.0) - 2026
- Microservices architecture
- Data persistence and history
- Machine learning-based anomaly detection
- Multi-cluster management
- Alert and notification system
- Third-party integrations (Prometheus, Grafana, etc.)

## 📊 Current Implementation Status

| Component | Status | Completion | Notes |
|-----------|--------|------------|-------|
| **Core Framework** | ✅ Complete | 100% | Architecture and interfaces defined |
| **Data Collector** | ✅ Complete | 100% | System metrics collection implemented |
| **Data Analyzer** | ✅ Complete | 90% | Health scoring and analysis logic |
| **Report Generator** | ✅ Complete | 80% | HTML, JSON, Markdown supported |
| **CLI Commands** | 🚧 In Progress | 60% | Basic commands functional |
| **Plugin System** | ✅ Complete | 70% | Interface defined, examples provided |
| **Web Dashboard** | 🚧 In Progress | 40% | UI prototype complete |
| **Configuration** | 🚧 In Progress | 70% | YAML config support |
| **Testing** | 🚧 In Progress | 30% | Unit tests for core modules |
| **Documentation** | ✅ Complete | 85% | Architecture and design docs |

## 🔧 Development Setup

### Prerequisites
- Go 1.21 or higher
- Linux or macOS (Windows support planned)
- Git

### Build and Run
```bash
# Install dependencies
go mod download

# Build
go build -o clusterreport ./cmd/cli

# Run tests
go test ./...

# Run with verbose output
./clusterreport collect --nodes localhost --verbose
```

### Development Commands
```bash
# Format code
go fmt ./...

# Run linter
golangci-lint run

# Generate documentation
go doc -all
```

## 🌟 Comparison with Other Tools

| Feature | ClusterReport | Generic Monitoring | Cloud Services |
|---------|--------------|-------------------|----------------|
| Hardware Config Analysis | ✅ Deep | ⚠️ Basic | ⚠️ Limited |
| Performance Benchmarking | ✅ Complete | ✅ Partial | ⚠️ Basic |
| Offline Reports | ✅ Supported | ⚠️ Partial | ❌ Requires Online |
| Customization | ✅ Highly Flexible | ⚠️ Limited | ⚠️ Limited |
| Batch Processing | ✅ Native Support | ⚠️ Script Required | ✅ Supported |
| Cost | ✅ Open Source Free | ⚠️ Partially Paid | ❌ Paid |
| Multi-Cluster Support | 📋 Planned | ✅ Supported | ✅ Supported |

## 📖 Learn More

- [Architecture Documentation](docs/tools/go/ClusterReport_Architecture.md) - Detailed architectural design
- [Design Documentation](docs/tools/go/ClusterReport_Design.md) - Feature specifications
- [Quick Start Guide](docs/getting-started/quick-start.md) - Get started quickly
- [Plugin Development](docs/tools/go/README.md) - How to create plugins

## 🙏 Acknowledgments

This project integrates concepts and learnings from:
- NodeProbe - Hardware configuration collection
- PerfSnap - Performance profiling and analysis
- Various open-source monitoring and reporting tools

---

⭐ **If this project helps you, please give us a Star!**

**Note**: The project is under active development (70% complete). Some features are not yet implemented. Contributions and suggestions are welcome! 🚀
