#!/usr/bin/env python3
"""
日志分析工具

功能：
- 解析和分析各种格式的日志文件
- 统计错误、警告和关键事件
- 生成日志分析报告
- 支持时间范围过滤
- 异常模式识别

Author: DevOps Team
Version: 1.0.0
License: MIT
"""

import argparse
import re
import sys
from collections import Counter, defaultdict
from datetime import datetime, timedelta
from pathlib import Path
from typing import Dict, List, Optional, Tuple
import json


class LogEntry:
    """日志条目"""
    
    def __init__(self, timestamp: Optional[datetime], level: str, message: str, raw: str):
        self.timestamp = timestamp
        self.level = level.upper()
        self.message = message
        self.raw = raw


class LogAnalyzer:
    """日志分析器"""
    
    # 常见日志级别
    LOG_LEVELS = ['DEBUG', 'INFO', 'WARN', 'WARNING', 'ERROR', 'CRITICAL', 'FATAL']
    
    # 常见时间戳格式
    TIMESTAMP_PATTERNS = [
        r'\d{4}-\d{2}-\d{2}\s+\d{2}:\d{2}:\d{2}',  # 2024-01-01 12:00:00
        r'\d{2}/\w{3}/\d{4}:\d{2}:\d{2}:\d{2}',   # 01/Jan/2024:12:00:00
        r'\w{3}\s+\d{1,2}\s+\d{2}:\d{2}:\d{2}',   # Jan 01 12:00:00
    ]
    
    def __init__(self, log_file: str, log_level: str = 'ALL'):
        self.log_file = Path(log_file)
        self.log_level = log_level.upper()
        self.entries: List[LogEntry] = []
        self.stats: Dict = {}
        
    def parse_log_file(self) -> None:
        """解析日志文件"""
        if not self.log_file.exists():
            raise FileNotFoundError(f"Log file not found: {self.log_file}")
        
        print(f"Parsing log file: {self.log_file}")
        
        with open(self.log_file, 'r', encoding='utf-8', errors='ignore') as f:
            for line in f:
                line = line.strip()
                if not line:
                    continue
                
                entry = self._parse_line(line)
                if entry and self._should_include(entry):
                    self.entries.append(entry)
        
        print(f"Parsed {len(self.entries)} log entries")
    
    def _parse_line(self, line: str) -> Optional[LogEntry]:
        """解析单行日志"""
        # 提取时间戳
        timestamp = self._extract_timestamp(line)
        
        # 提取日志级别
        level = self._extract_level(line)
        
        # 提取消息
        message = line
        
        return LogEntry(timestamp, level, message, line)
    
    def _extract_timestamp(self, line: str) -> Optional[datetime]:
        """提取时间戳"""
        for pattern in self.TIMESTAMP_PATTERNS:
            match = re.search(pattern, line)
            if match:
                timestamp_str = match.group()
                try:
                    # 尝试不同的时间格式
                    for fmt in ['%Y-%m-%d %H:%M:%S', '%d/%b/%Y:%H:%M:%S', '%b %d %H:%M:%S']:
                        try:
                            return datetime.strptime(timestamp_str, fmt)
                        except ValueError:
                            continue
                except:
                    pass
        return None
    
    def _extract_level(self, line: str) -> str:
        """提取日志级别"""
        line_upper = line.upper()
        for level in self.LOG_LEVELS:
            if level in line_upper:
                return level
        return 'UNKNOWN'
    
    def _should_include(self, entry: LogEntry) -> bool:
        """判断是否应该包含此日志条目"""
        if self.log_level == 'ALL':
            return True
        return entry.level == self.log_level
    
    def analyze(self) -> None:
        """分析日志"""
        print("Analyzing logs...")
        
        self.stats = {
            'total_entries': len(self.entries),
            'level_counts': Counter(),
            'hourly_distribution': defaultdict(int),
            'error_patterns': Counter(),
            'top_messages': Counter(),
        }
        
        for entry in self.entries:
            # 统计日志级别
            self.stats['level_counts'][entry.level] += 1
            
            # 统计时间分布
            if entry.timestamp:
                hour_key = entry.timestamp.strftime('%Y-%m-%d %H:00')
                self.stats['hourly_distribution'][hour_key] += 1
            
            # 统计错误模式
            if entry.level in ['ERROR', 'CRITICAL', 'FATAL']:
                # 提取错误类型
                error_pattern = self._extract_error_pattern(entry.message)
                self.stats['error_patterns'][error_pattern] += 1
            
            # 统计常见消息
            message_key = entry.message[:100]  # 前100字符
            self.stats['top_messages'][message_key] += 1
    
    def _extract_error_pattern(self, message: str) -> str:
        """提取错误模式"""
        # 尝试提取异常类型
        exception_match = re.search(r'(\w+Exception|\w+Error)', message)
        if exception_match:
            return exception_match.group(1)
        
        # 提取第一个单词作为模式
        words = message.split()
        if words:
            return words[0]
        
        return 'Unknown'
    
    def generate_report(self, output_format: str = 'text') -> str:
        """生成分析报告"""
        if output_format == 'json':
            return self._generate_json_report()
        elif output_format == 'html':
            return self._generate_html_report()
        else:
            return self._generate_text_report()
    
    def _generate_text_report(self) -> str:
        """生成文本报告"""
        report = []
        report.append("=" * 80)
        report.append("LOG ANALYSIS REPORT")
        report.append("=" * 80)
        report.append(f"Log File: {self.log_file}")
        report.append(f"Analysis Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        report.append(f"Total Entries: {self.stats['total_entries']}")
        report.append("")
        
        # 日志级别统计
        report.append("-" * 80)
        report.append("LOG LEVEL DISTRIBUTION")
        report.append("-" * 80)
        for level, count in self.stats['level_counts'].most_common():
            percentage = (count / self.stats['total_entries']) * 100
            report.append(f"{level:12} {count:8} ({percentage:5.2f}%)")
        report.append("")
        
        # 错误模式
        if self.stats['error_patterns']:
            report.append("-" * 80)
            report.append("TOP ERROR PATTERNS")
            report.append("-" * 80)
            for pattern, count in self.stats['error_patterns'].most_common(10):
                report.append(f"{pattern:40} {count:8}")
            report.append("")
        
        # 时间分布
        if self.stats['hourly_distribution']:
            report.append("-" * 80)
            report.append("HOURLY DISTRIBUTION (Last 24 hours)")
            report.append("-" * 80)
            sorted_hours = sorted(self.stats['hourly_distribution'].items())[-24:]
            for hour, count in sorted_hours:
                bar = '#' * (count // 10 if count >= 10 else 1)
                report.append(f"{hour:20} {count:6} {bar}")
            report.append("")
        
        # 常见消息
        report.append("-" * 80)
        report.append("TOP 10 MESSAGES")
        report.append("-" * 80)
        for message, count in self.stats['top_messages'].most_common(10):
            report.append(f"[{count:5}] {message[:70]}")
        report.append("")
        
        report.append("=" * 80)
        
        return '\n'.join(report)
    
    def _generate_json_report(self) -> str:
        """生成JSON报告"""
        report_data = {
            'log_file': str(self.log_file),
            'analysis_time': datetime.now().isoformat(),
            'total_entries': self.stats['total_entries'],
            'level_counts': dict(self.stats['level_counts']),
            'error_patterns': dict(self.stats['error_patterns'].most_common(20)),
            'hourly_distribution': dict(sorted(self.stats['hourly_distribution'].items())[-24:]),
            'top_messages': [
                {'message': msg, 'count': count}
                for msg, count in self.stats['top_messages'].most_common(20)
            ]
        }
        return json.dumps(report_data, indent=2)
    
    def _generate_html_report(self) -> str:
        """生成HTML报告"""
        html = f"""
<!DOCTYPE html>
<html>
<head>
    <title>Log Analysis Report</title>
    <style>
        body {{ font-family: Arial, sans-serif; margin: 20px; background-color: #f5f5f5; }}
        .container {{ max-width: 1200px; margin: 0 auto; background-color: white; padding: 20px; border-radius: 5px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }}
        h1 {{ color: #333; border-bottom: 2px solid #007bff; padding-bottom: 10px; }}
        h2 {{ color: #555; margin-top: 30px; }}
        .stat-box {{ background-color: #f8f9fa; padding: 15px; margin: 10px 0; border-left: 4px solid #007bff; }}
        table {{ width: 100%; border-collapse: collapse; margin: 20px 0; }}
        th, td {{ padding: 12px; text-align: left; border-bottom: 1px solid #ddd; }}
        th {{ background-color: #007bff; color: white; }}
        tr:hover {{ background-color: #f5f5f5; }}
        .error {{ color: #dc3545; }}
        .warning {{ color: #ffc107; }}
        .info {{ color: #17a2b8; }}
    </style>
</head>
<body>
    <div class="container">
        <h1>Log Analysis Report</h1>
        <div class="stat-box">
            <p><strong>Log File:</strong> {self.log_file}</p>
            <p><strong>Analysis Time:</strong> {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}</p>
            <p><strong>Total Entries:</strong> {self.stats['total_entries']}</p>
        </div>
        
        <h2>Log Level Distribution</h2>
        <table>
            <tr><th>Level</th><th>Count</th><th>Percentage</th></tr>
"""
        for level, count in self.stats['level_counts'].most_common():
            percentage = (count / self.stats['total_entries']) * 100
            html += f"            <tr><td>{level}</td><td>{count}</td><td>{percentage:.2f}%</td></tr>\n"
        
        html += """
        </table>
        
        <h2>Top Error Patterns</h2>
        <table>
            <tr><th>Pattern</th><th>Count</th></tr>
"""
        for pattern, count in self.stats['error_patterns'].most_common(10):
            html += f"            <tr><td class='error'>{pattern}</td><td>{count}</td></tr>\n"
        
        html += """
        </table>
    </div>
</body>
</html>
"""
        return html


def main():
    """主函数"""
    parser = argparse.ArgumentParser(
        description='日志分析工具 - 解析和分析日志文件',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 分析所有日志
  %(prog)s -f /var/log/app.log
  
  # 只分析ERROR级别
  %(prog)s -f /var/log/app.log -l ERROR
  
  # 生成HTML报告
  %(prog)s -f /var/log/app.log -o report.html -format html
  
  # 生成JSON报告
  %(prog)s -f /var/log/app.log -o report.json -format json
"""
    )
    
    parser.add_argument(
        '-f', '--file',
        required=True,
        help='日志文件路径'
    )
    parser.add_argument(
        '-l', '--level',
        default='ALL',
        choices=['ALL', 'DEBUG', 'INFO', 'WARN', 'WARNING', 'ERROR', 'CRITICAL', 'FATAL'],
        help='过滤日志级别 (默认: ALL)'
    )
    parser.add_argument(
        '-o', '--output',
        help='输出文件路径 (默认: 输出到终端)'
    )
    parser.add_argument(
        '-format', '--format',
        default='text',
        choices=['text', 'json', 'html'],
        help='报告格式 (默认: text)'
    )
    parser.add_argument(
        '-v', '--version',
        action='version',
        version='%(prog)s 1.0.0'
    )
    
    args = parser.parse_args()
    
    try:
        # 创建分析器
        analyzer = LogAnalyzer(args.file, args.level)
        
        # 解析日志
        analyzer.parse_log_file()
        
        # 分析
        analyzer.analyze()
        
        # 生成报告
        report = analyzer.generate_report(args.format)
        
        # 输出报告
        if args.output:
            with open(args.output, 'w', encoding='utf-8') as f:
                f.write(report)
            print(f"\nReport saved to: {args.output}")
        else:
            print("\n" + report)
        
        return 0
        
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        return 1


if __name__ == '__main__':
    sys.exit(main())
