#!/bin/bash
AUDIT_LOG_DIR="/home/service/*/_package/run/auditlog/*"

# 遍历每个匹配的目录
for dir in $AUDIT_LOG_DIR; do
    # 检查是否是目录
    if [ -d "$dir" ]; then
        echo "处理目录: $dir"
        
        # 获取目录中的总文件数
        total_files=$(find "$dir" -maxdepth 1 -type f | wc -l)
        
        # 获取超过1天的文件数
        old_files=$(find "$dir" -maxdepth 1 -type f -mtime +1 | wc -l)
        
        echo "  总文件数: $total_files"
        echo "  过期文件数: $old_files"
        
        # 如果目录中只有一个文件，不管是否过期都保留
        if [ $total_files -le 1 ]; then
            echo "  目录中只有 $total_files 个文件，保留不删除"
        elif [ $old_files -gt 0 ]; then
            # 如果有多个文件且有过期文件
            # 计算需要保留的文件数（至少保留1个）
            files_to_keep=1
            
            # 如果删除所有过期文件后会导致目录为空，则保留最新的过期文件
            new_files=$((total_files - old_files))
            if [ $new_files -eq 0 ]; then
                # 没有新文件，需要从过期文件中保留至少一个
                echo "  没有新文件，将从过期文件中保留最新的一个"
                # 删除最旧的过期文件，保留最新的
                find "$dir" -maxdepth 1 -type f -mtime +1 -printf '%T@ %p\n' | \
                sort -n | \
                head -n -1 | \
                cut -d' ' -f2- | \
                xargs -r rm -f
            else
                # 有新文件存在，可以安全删除所有过期文件
                echo "  存在 $new_files 个新文件，删除所有过期文件"
                find "$dir" -maxdepth 1 -type f -mtime +1 -exec rm -f {} \;
            fi
            echo "  清理完成"
        else
            echo "  没有过期文件需要清理"
        fi
        echo ""
    fi
done

echo "所有目录处理完成"
