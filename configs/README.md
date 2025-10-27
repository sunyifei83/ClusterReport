# Configs 目录说明

本目录包含各种常用服务和工具的配置模板，帮助快速部署和配置基础设施。

## 📁 目录结构

```
configs/
├── docker/              # Docker 配置
├── kubernetes/          # Kubernetes 配置
├── nginx/              # Nginx 配置
└── terraform/          # Terraform 配置
```

## 🚧 开发状态

当前此目录处于规划阶段，各子目录将逐步添加配置模板和最佳实践示例。

### 计划添加的内容

#### docker/ - Docker 配置
- [ ] Dockerfile 最佳实践模板
- [ ] docker-compose.yml 示例
- [ ] 多阶段构建示例
- [ ] Docker 网络配置
- [ ] Docker 卷配置
- [ ] 镜像优化示例

**示例文件规划**：
```
docker/
├── Dockerfile.golang        # Go 应用模板
├── Dockerfile.python        # Python 应用模板
├── Dockerfile.nodejs        # Node.js 应用模板
├── docker-compose.yml       # 基础编排
├── docker-compose.prod.yml  # 生产环境
└── .dockerignore           # 忽略文件模板
```

#### kubernetes/ - Kubernetes 配置
- [ ] Deployment 配置模板
- [ ] Service 配置模板
- [ ] Ingress 配置模板
- [ ] ConfigMap 示例
- [ ] Secret 管理示例
- [ ] StatefulSet 配置
- [ ] DaemonSet 配置
- [ ] Job 和 CronJob 示例
- [ ] HPA (水平自动扩缩容) 配置
- [ ] NetworkPolicy 示例

**示例文件规划**：
```
kubernetes/
├── deployment.yaml          # Deployment 模板
├── service.yaml            # Service 模板
├── ingress.yaml            # Ingress 模板
├── configmap.yaml          # ConfigMap 示例
├── secret.yaml             # Secret 示例
├── statefulset.yaml        # StatefulSet 示例
├── hpa.yaml                # HPA 配置
├── namespace.yaml          # Namespace 定义
└── complete-app/           # 完整应用示例
    ├── backend/
    ├── frontend/
    └── database/
```

#### nginx/ - Nginx 配置
- [ ] 基础配置模板
- [ ] 反向代理配置
- [ ] 负载均衡配置
- [ ] SSL/TLS 配置
- [ ] 缓存配置
- [ ] 限流配置
- [ ] 日志配置
- [ ] 性能优化配置

**示例文件规划**：
```
nginx/
├── nginx.conf              # 主配置文件模板
├── ssl.conf                # SSL 配置
├── proxy.conf              # 反向代理
├── load-balancer.conf      # 负载均衡
├── cache.conf              # 缓存配置
├── rate-limit.conf         # 限流配置
├── security.conf           # 安全配置
└── sites-available/        # 站点配置示例
    ├── static-site.conf
    ├── api-gateway.conf
    └── websocket.conf
```

#### terraform/ - Terraform 配置
- [ ] AWS 基础设施模板
- [ ] GCP 基础设施模板
- [ ] 阿里云基础设施模板
- [ ] 腾讯云基础设施模板
- [ ] 模块化配置示例
- [ ] 状态管理最佳实践
- [ ] 多环境管理

**示例文件规划**：
```
terraform/
├── aws/
│   ├── ec2/                # EC2 实例
│   ├── vpc/                # VPC 网络
│   ├── rds/                # RDS 数据库
│   └── s3/                 # S3 存储
├── gcp/
│   ├── compute/            # Compute Engine
│   ├── network/            # VPC 网络
│   └── storage/            # Cloud Storage
├── alicloud/
│   ├── ecs/                # ECS 实例
│   ├── vpc/                # VPC 网络
│   └── oss/                # OSS 存储
└── modules/                # 可重用模块
    ├── network/
    ├── compute/
    └── database/
```

## 📝 配置文件编写规范

### 通用原则
1. **注释清晰**：每个配置项都应有清晰的注释说明
2. **安全第一**：不包含敏感信息，使用环境变量或密钥管理
3. **版本控制**：配置文件应该可以安全地版本控制
4. **环境分离**：区分开发、测试、生产环境的配置

### Docker 配置示例
```dockerfile
# 使用官方基础镜像
FROM golang:1.21-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码并构建
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 生产镜像
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080

# 运行应用
CMD ["./main"]
```

### Kubernetes 配置示例
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-deployment
  labels:
    app: myapp
spec:
  replicas: 3
  selector:
    matchLabels:
      app: myapp
  template:
    metadata:
      labels:
        app: myapp
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENV
          value: "production"
        resources:
          requests:
            memory: "64Mi"
            cpu: "250m"
          limits:
            memory: "128Mi"
            cpu: "500m"
```

### Nginx 配置示例
```nginx
# 基础服务器配置
server {
    listen 80;
    server_name example.com;
    
    # 日志配置
    access_log /var/log/nginx/example.access.log;
    error_log /var/log/nginx/example.error.log;
    
    # 反向代理配置
    location / {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

### Terraform 配置示例
```hcl
# 定义提供商
terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}

provider "aws" {
  region = var.aws_region
}

# 创建 VPC
resource "aws_vpc" "main" {
  cidr_block = var.vpc_cidr
  
  tags = {
    Name        = "${var.project_name}-vpc"
    Environment = var.environment
  }
}

# 定义变量
variable "aws_region" {
  description = "AWS 区域"
  type        = string
  default     = "us-east-1"
}
```

## 🔒 安全建议

### 敏感信息管理
1. **不要**在配置文件中硬编码密码、密钥等敏感信息
2. 使用环境变量或密钥管理服务（如 AWS Secrets Manager、HashiCorp Vault）
3. 使用 `.env.example` 文件作为模板，不提交实际的 `.env` 文件

### 示例：.env.example
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=mydb
DB_USER=user
DB_PASSWORD=changeme

# API 密钥
API_KEY=your_api_key_here
API_SECRET=your_api_secret_here

# 应用配置
APP_ENV=production
APP_DEBUG=false
```

## 🤝 贡献

欢迎贡献配置模板和最佳实践！

贡献时请确保：
1. 配置文件有详细的注释
2. 包含使用说明和示例
3. 不包含任何敏感信息
4. 遵循相应工具的最佳实践

## 📚 相关文档

- [用户指南](/docs/UserGuide.md)
- [最佳实践](/docs/BestPractices.md)
- [Docker 官方文档](https://docs.docker.com/)
- [Kubernetes 官方文档](https://kubernetes.io/docs/)
- [Nginx 官方文档](https://nginx.org/en/docs/)
- [Terraform 官方文档](https://www.terraform.io/docs/)

## 📮 反馈

如有问题或建议，请通过以下方式联系：
- GitHub Issues: https://github.com/sunyifei83/devops-toolkit/issues
- Email: sunyifei83@gmail.com

---

*此文档会随着项目发展持续更新*
