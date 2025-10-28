# Playbooks 目录说明

本目录包含 Ansible playbooks 用于自动化系统配置、应用部署和基础设施管理。

## 📁 目录结构

```
playbooks/
└── setup/               # 系统初始化和设置
```

## 🚧 开发状态

当前此目录处于规划阶段，将逐步添加各种 Ansible playbooks。

### 计划添加的内容

#### setup/ - 系统初始化
- [ ] 基础系统配置
- [ ] 用户和权限设置
- [ ] 安全加固
- [ ] 软件包安装
- [ ] 时区和本地化配置
- [ ] SSH 配置优化

**示例规划**：
```
setup/
├── base-system.yml          # 基础系统配置
├── security-hardening.yml   # 安全加固
├── user-management.yml      # 用户管理
├── ssh-config.yml           # SSH 配置
└── timezone-setup.yml       # 时区设置
```

#### deployment/ - 应用部署 (规划中)
- [ ] Web 应用部署
- [ ] 数据库部署
- [ ] 容器化应用部署
- [ ] 微服务部署
- [ ] 负载均衡配置

#### monitoring/ - 监控部署 (规划中)
- [ ] Prometheus 部署
- [ ] Grafana 部署
- [ ] Node Exporter 部署
- [ ] AlertManager 部署
- [ ] 日志收集配置

#### maintenance/ - 系统维护 (规划中)
- [ ] 系统更新
- [ ] 日志轮转
- [ ] 备份任务
- [ ] 清理任务
- [ ] 健康检查

## 📝 Playbook 编写规范

### 基础结构
```yaml
---
# Playbook 名称和描述
# 功能：描述这个 playbook 的作用
# 作者：xxx
# 版本：x.x.x
# 最后更新：YYYY-MM-DD

- name: Playbook 描述
  hosts: target_hosts
  become: yes  # 如需 root 权限
  vars:
    # 变量定义
    variable_name: value
  
  tasks:
    - name: 任务描述
      module_name:
        parameter: value
      tags:
        - tag_name
```

### 完整示例

#### 基础系统配置 Playbook
```yaml
---
# 基础系统配置 playbook
# 功能：配置新服务器的基础设置
# 作者：DevOps Team
# 版本：1.0.0

- name: Configure Base System
  hosts: all
  become: yes
  vars:
    timezone: "Asia/Shanghai"
    ntp_servers:
      - "ntp1.aliyun.com"
      - "ntp2.aliyun.com"
  
  tasks:
    - name: Update apt cache
      apt:
        update_cache: yes
        cache_valid_time: 3600
      when: ansible_os_family == "Debian"
      tags:
        - packages
    
    - name: Install essential packages
      apt:
        name:
          - vim
          - curl
          - wget
          - git
          - htop
        state: present
      when: ansible_os_family == "Debian"
      tags:
        - packages
    
    - name: Set timezone
      timezone:
        name: "{{ timezone }}"
      tags:
        - timezone
    
    - name: Configure NTP
      template:
        src: ntp.conf.j2
        dest: /etc/ntp.conf
        owner: root
        group: root
        mode: '0644'
      notify: restart ntp
      tags:
        - ntp
  
  handlers:
    - name: restart ntp
      service:
        name: ntp
        state: restarted
```

#### Inventory 示例
```ini
# inventory/hosts
[webservers]
web1 ansible_host=192.168.1.10
web2 ansible_host=192.168.1.11

[databases]
db1 ansible_host=192.168.1.20
db2 ansible_host=192.168.1.21

[all:vars]
ansible_user=admin
ansible_ssh_private_key_file=~/.ssh/id_rsa
ansible_python_interpreter=/usr/bin/python3
```

#### 变量文件示例
```yaml
# group_vars/all.yml
---
common_packages:
  - vim
  - curl
  - wget
  - git
  - htop
  - net-tools

timezone: "Asia/Shanghai"

# group_vars/webservers.yml
---
nginx_version: "1.20.2"
document_root: "/var/www/html"
```

## 🚀 使用方法

### 基本命令

```bash
# 测试连接
ansible all -i inventory/hosts -m ping

# 运行 playbook
ansible-playbook -i inventory/hosts playbooks/setup/base-system.yml

# 使用特定标签
ansible-playbook -i inventory/hosts playbooks/setup/base-system.yml --tags packages

# 检查模式（不实际执行）
ansible-playbook -i inventory/hosts playbooks/setup/base-system.yml --check

# 显示详细输出
ansible-playbook -i inventory/hosts playbooks/setup/base-system.yml -v

# 限制目标主机
ansible-playbook -i inventory/hosts playbooks/setup/base-system.yml --limit web1

# 使用 vault 加密的变量
ansible-playbook -i inventory/hosts playbooks/setup/base-system.yml --ask-vault-pass
```

### 最佳实践

#### 1. 目录结构
```
playbooks/
├── inventory/
│   ├── production/
│   │   ├── hosts
│   │   └── group_vars/
│   └── staging/
│       ├── hosts
│       └── group_vars/
├── roles/
│   ├── common/
│   ├── webserver/
│   └── database/
├── setup/
├── deployment/
└── maintenance/
```

#### 2. 使用 Roles
```yaml
---
- name: Deploy Web Application
  hosts: webservers
  become: yes
  roles:
    - common
    - webserver
    - application
```

#### 3. 变量优先级
- 命令行变量 (最高)
- play vars
- play vars_files
- role defaults (最低)

#### 4. 错误处理
```yaml
tasks:
  - name: Task that might fail
    command: /bin/false
    register: result
    ignore_errors: yes
  
  - name: Handle failure
    debug:
      msg: "Previous task failed"
    when: result.failed
```

#### 5. 幂等性
```yaml
# 好的示例 - 幂等
- name: Ensure nginx is installed
  apt:
    name: nginx
    state: present

# 避免 - 非幂等
- name: Install nginx
  command: apt-get install -y nginx
```

## 🔒 安全建议

### 1. 使用 Ansible Vault 加密敏感信息
```bash
# 创建加密文件
ansible-vault create secrets.yml

# 编辑加密文件
ansible-vault edit secrets.yml

# 查看加密文件
ansible-vault view secrets.yml

# 加密现有文件
ansible-vault encrypt vars/passwords.yml
```

### 2. 敏感变量示例
```yaml
# secrets.yml (加密)
---
db_password: "secure_password"
api_key: "secret_api_key"
```

### 3. SSH 密钥管理
```bash
# 生成专用密钥
ssh-keygen -t ed25519 -f ~/.ssh/ansible_key

# 在 inventory 中使用
ansible_ssh_private_key_file=~/.ssh/ansible_key
```

## 📚 相关资源

### 官方文档
- [Ansible 官方文档](https://docs.ansible.com/)
- [Ansible Galaxy](https://galaxy.ansible.com/) - Role 仓库
- [Ansible 最佳实践](https://docs.ansible.com/ansible/latest/user_guide/playbooks_best_practices.html)

### 推荐 Roles
- `geerlingguy.nginx` - Nginx 配置
- `geerlingguy.docker` - Docker 安装
- `geerlingguy.mysql` - MySQL 配置

### 项目文档
- [用户指南](/docs/UserGuide.md)
- [最佳实践](/docs/BestPractices.md)
- [工具文档](/docs/ToolsDocumentation.md)

## 🤝 贡献

欢迎贡献 playbooks 和最佳实践！

贡献时请确保：
1. Playbook 有清晰的注释和文档
2. 使用标准的 Ansible 最佳实践
3. 测试 playbook 的幂等性
4. 不包含任何敏感信息
5. 提供使用示例

## 📮 反馈

如有问题或建议，请通过以下方式联系：
- GitHub Issues: https://github.com/sunyifei83/devops-toolkit/issues
- Email: sunyifei83@gmail.com

---

*此文档会随着项目发展持续更新*
