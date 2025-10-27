#!/bin/bash
#
# 数据库备份脚本
# 功能：自动备份 MySQL/PostgreSQL 数据库
# 作者：DevOps Team
# 版本：1.0.0
# 最后更新：2025-10-27

set -euo pipefail

# ============================================
# 配置部分
# ============================================

# 备份目录
BACKUP_DIR="${BACKUP_DIR:-/var/backups/databases}"
# 保留天数
RETENTION_DAYS="${RETENTION_DAYS:-7}"
# 日志文件
LOG_FILE="${LOG_FILE:-/var/log/database-backup.log}"
# 时间戳格式
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# 数据库配置（从环境变量或配置文件读取）
DB_TYPE="${DB_TYPE:-mysql}"  # mysql 或 postgres
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-}"
DB_NAMES="${DB_NAMES:-all}"  # 数据库名称，多个用逗号分隔，或 'all' 备份所有

# 压缩选项
COMPRESS="${COMPRESS:-true}"
COMPRESS_TOOL="${COMPRESS_TOOL:-gzip}"  # gzip 或 xz

# 通知配置
ENABLE_NOTIFICATION="${ENABLE_NOTIFICATION:-false}"
NOTIFICATION_EMAIL="${NOTIFICATION_EMAIL:-admin@example.com}"

# ============================================
# 函数定义
# ============================================

# 日志函数
log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*" | tee -a "$LOG_FILE"
}

log_error() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] ERROR: $*" | tee -a "$LOG_FILE" >&2
}

# 检查必要的命令
check_requirements() {
    local required_cmds=()
    
    if [ "$DB_TYPE" = "mysql" ]; then
        required_cmds+=("mysqldump" "mysql")
    elif [ "$DB_TYPE" = "postgres" ]; then
        required_cmds+=("pg_dump" "psql")
    fi
    
    if [ "$COMPRESS" = "true" ]; then
        required_cmds+=("$COMPRESS_TOOL")
    fi
    
    for cmd in "${required_cmds[@]}"; do
        if ! command -v "$cmd" &> /dev/null; then
            log_error "Required command not found: $cmd"
            exit 1
        fi
    done
}

# 创建备份目录
create_backup_dir() {
    if [ ! -d "$BACKUP_DIR" ]; then
        mkdir -p "$BACKUP_DIR"
        log "Created backup directory: $BACKUP_DIR"
    fi
}

# 获取数据库列表
get_databases() {
    local databases=()
    
    if [ "$DB_NAMES" = "all" ]; then
        if [ "$DB_TYPE" = "mysql" ]; then
            databases=$(mysql -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" \
                -e "SHOW DATABASES;" | grep -Ev "^(Database|information_schema|performance_schema|mysql|sys)$")
        elif [ "$DB_TYPE" = "postgres" ]; then
            databases=$(PGPASSWORD="$DB_PASSWORD" psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" \
                -t -c "SELECT datname FROM pg_database WHERE datistemplate = false AND datname != 'postgres';" | xargs)
        fi
    else
        IFS=',' read -ra databases <<< "$DB_NAMES"
    fi
    
    echo "${databases[@]}"
}

# 备份 MySQL 数据库
backup_mysql() {
    local db_name=$1
    local backup_file="${BACKUP_DIR}/${db_name}_${TIMESTAMP}.sql"
    
    log "Backing up MySQL database: $db_name"
    
    if mysqldump -h "$DB_HOST" -P "$DB_PORT" -u "$DB_USER" -p"$DB_PASSWORD" \
        --single-transaction --quick --lock-tables=false \
        "$db_name" > "$backup_file" 2>> "$LOG_FILE"; then
        
        log "MySQL backup completed: $backup_file"
        
        if [ "$COMPRESS" = "true" ]; then
            compress_backup "$backup_file"
        fi
        
        return 0
    else
        log_error "MySQL backup failed for database: $db_name"
        return 1
    fi
}

# 备份 PostgreSQL 数据库
backup_postgres() {
    local db_name=$1
    local backup_file="${BACKUP_DIR}/${db_name}_${TIMESTAMP}.sql"
    
    log "Backing up PostgreSQL database: $db_name"
    
    if PGPASSWORD="$DB_PASSWORD" pg_dump -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" \
        -F p -f "$backup_file" "$db_name" 2>> "$LOG_FILE"; then
        
        log "PostgreSQL backup completed: $backup_file"
        
        if [ "$COMPRESS" = "true" ]; then
            compress_backup "$backup_file"
        fi
        
        return 0
    else
        log_error "PostgreSQL backup failed for database: $db_name"
        return 1
    fi
}

# 压缩备份文件
compress_backup() {
    local file=$1
    
    log "Compressing backup file: $file"
    
    if [ "$COMPRESS_TOOL" = "gzip" ]; then
        gzip "$file"
        log "Compressed with gzip: ${file}.gz"
    elif [ "$COMPRESS_TOOL" = "xz" ]; then
        xz "$file"
        log "Compressed with xz: ${file}.xz"
    fi
}

# 清理旧备份
cleanup_old_backups() {
    log "Cleaning up backups older than $RETENTION_DAYS days"
    
    find "$BACKUP_DIR" -name "*.sql*" -type f -mtime +$RETENTION_DAYS -delete
    
    log "Old backups cleaned up"
}

# 发送通知
send_notification() {
    local status=$1
    local message=$2
    
    if [ "$ENABLE_NOTIFICATION" = "true" ]; then
        if command -v mail &> /dev/null; then
            echo "$message" | mail -s "Database Backup $status" "$NOTIFICATION_EMAIL"
        else
            log "Mail command not found, skipping notification"
        fi
    fi
}

# ============================================
# 主流程
# ============================================

main() {
    log "=========================================="
    log "Starting database backup process"
    log "=========================================="
    
    # 检查必要条件
    check_requirements
    create_backup_dir
    
    # 获取数据库列表
    databases=$(get_databases)
    
    if [ -z "$databases" ]; then
        log_error "No databases found to backup"
        send_notification "FAILED" "No databases found to backup"
        exit 1
    fi
    
    log "Databases to backup: $databases"
    
    # 备份每个数据库
    success_count=0
    fail_count=0
    
    for db in $databases; do
        if [ "$DB_TYPE" = "mysql" ]; then
            if backup_mysql "$db"; then
                ((success_count++))
            else
                ((fail_count++))
            fi
        elif [ "$DB_TYPE" = "postgres" ]; then
            if backup_postgres "$db"; then
                ((success_count++))
            else
                ((fail_count++))
            fi
        fi
    done
    
    # 清理旧备份
    cleanup_old_backups
    
    # 生成报告
    log "=========================================="
    log "Backup Summary:"
    log "  Total databases: $((success_count + fail_count))"
    log "  Successful: $success_count"
    log "  Failed: $fail_count"
    log "  Backup location: $BACKUP_DIR"
    log "=========================================="
    
    # 发送通知
    if [ $fail_count -eq 0 ]; then
        send_notification "SUCCESS" "All $success_count database(s) backed up successfully"
        exit 0
    else
        send_notification "PARTIAL" "$success_count succeeded, $fail_count failed"
        exit 1
    fi
}

# 执行主函数
main "$@"
