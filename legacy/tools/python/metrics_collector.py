#!/usr/bin/env python3
"""
指标收集器

功能：
- 收集系统性能指标
- 自定义指标采集
- 数据格式转换
- 时序数据处理
- 导出到监控系统

Author: DevOps Team
Version: 1.0.0
License: MIT
"""

import argparse
import json
import sys
import time
import psutil
from datetime import datetime
from typing import Dict, List, Optional
from pathlib import Path


class MetricsCollector:
    """指标收集器"""
    
    def __init__(self, interval: int = 60):
        self.interval = interval
        self.metrics: Dict = {}
        
    def collect_system_metrics(self) -> Dict:
        """收集系统指标"""
        metrics = {
            'timestamp': datetime.now().isoformat(),
            'cpu': self._collect_cpu_metrics(),
            'memory': self._collect_memory_metrics(),
            'disk': self._collect_disk_metrics(),
            'network': self._collect_network_metrics(),
            'process': self._collect_process_metrics(),
        }
        return metrics
    
    def _collect_cpu_metrics(self) -> Dict:
        """收集 CPU 指标"""
        cpu_percent = psutil.cpu_percent(interval=1, percpu=True)
        cpu_freq = psutil.cpu_freq()
        cpu_stats = psutil.cpu_stats()
        load_avg = psutil.getloadavg() if hasattr(psutil, 'getloadavg') else (0, 0, 0)
        
        return {
            'percent_total': psutil.cpu_percent(interval=1),
            'percent_per_cpu': cpu_percent,
            'count_logical': psutil.cpu_count(logical=True),
            'count_physical': psutil.cpu_count(logical=False),
            'frequency_current': cpu_freq.current if cpu_freq else 0,
            'frequency_min': cpu_freq.min if cpu_freq else 0,
            'frequency_max': cpu_freq.max if cpu_freq else 0,
            'load_avg_1': load_avg[0],
            'load_avg_5': load_avg[1],
            'load_avg_15': load_avg[2],
            'ctx_switches': cpu_stats.ctx_switches,
            'interrupts': cpu_stats.interrupts,
        }
    
    def _collect_memory_metrics(self) -> Dict:
        """收集内存指标"""
        mem = psutil.virtual_memory()
        swap = psutil.swap_memory()
        
        return {
            'total': mem.total,
            'available': mem.available,
            'used': mem.used,
            'free': mem.free,
            'percent': mem.percent,
            'active': getattr(mem, 'active', 0),
            'inactive': getattr(mem, 'inactive', 0),
            'buffers': getattr(mem, 'buffers', 0),
            'cached': getattr(mem, 'cached', 0),
            'swap_total': swap.total,
            'swap_used': swap.used,
            'swap_free': swap.free,
            'swap_percent': swap.percent,
        }
    
    def _collect_disk_metrics(self) -> List[Dict]:
        """收集磁盘指标"""
        disk_metrics = []
        
        for partition in psutil.disk_partitions():
            try:
                usage = psutil.disk_usage(partition.mountpoint)
                io_counters = psutil.disk_io_counters(perdisk=False)
                
                disk_metrics.append({
                    'device': partition.device,
                    'mountpoint': partition.mountpoint,
                    'fstype': partition.fstype,
                    'total': usage.total,
                    'used': usage.used,
                    'free': usage.free,
                    'percent': usage.percent,
                    'read_count': io_counters.read_count if io_counters else 0,
                    'write_count': io_counters.write_count if io_counters else 0,
                    'read_bytes': io_counters.read_bytes if io_counters else 0,
                    'write_bytes': io_counters.write_bytes if io_counters else 0,
                })
            except (PermissionError, OSError):
                continue
        
        return disk_metrics
    
    def _collect_network_metrics(self) -> List[Dict]:
        """收集网络指标"""
        network_metrics = []
        net_io = psutil.net_io_counters(pernic=True)
        
        for interface, counters in net_io.items():
            network_metrics.append({
                'interface': interface,
                'bytes_sent': counters.bytes_sent,
                'bytes_recv': counters.bytes_recv,
                'packets_sent': counters.packets_sent,
                'packets_recv': counters.packets_recv,
                'errin': counters.errin,
                'errout': counters.errout,
                'dropin': counters.dropin,
                'dropout': counters.dropout,
            })
        
        return network_metrics
    
    def _collect_process_metrics(self) -> Dict:
        """收集进程指标"""
        processes = list(psutil.process_iter(['status']))
        
        status_count = {
            'running': 0,
            'sleeping': 0,
            'zombie': 0,
            'stopped': 0,
            'other': 0,
        }
        
        for proc in processes:
            try:
                status = proc.info['status']
                if status == psutil.STATUS_RUNNING:
                    status_count['running'] += 1
                elif status in [psutil.STATUS_SLEEPING, psutil.STATUS_DISK_SLEEP]:
                    status_count['sleeping'] += 1
                elif status == psutil.STATUS_ZOMBIE:
                    status_count['zombie'] += 1
                elif status == psutil.STATUS_STOPPED:
                    status_count['stopped'] += 1
                else:
                    status_count['other'] += 1
            except (psutil.NoSuchProcess, psutil.AccessDenied):
                continue
        
        return {
            'total': len(processes),
            **status_count
        }
    
    def collect_custom_metrics(self, custom_commands: List[str]) -> Dict:
        """收集自定义指标"""
        import subprocess
        
        custom_metrics = {}
        for i, command in enumerate(custom_commands):
            try:
                result = subprocess.run(
                    command,
                    shell=True,
                    capture_output=True,
                    text=True,
                    timeout=10
                )
                custom_metrics[f'custom_{i}'] = {
                    'command': command,
                    'output': result.stdout.strip(),
                    'exit_code': result.returncode,
                }
            except subprocess.TimeoutExpired:
                custom_metrics[f'custom_{i}'] = {
                    'command': command,
                    'error': 'Timeout',
                }
            except Exception as e:
                custom_metrics[f'custom_{i}'] = {
                    'command': command,
                    'error': str(e),
                }
        
        return custom_metrics
    
    def export_prometheus(self, metrics: Dict) -> str:
        """导出为 Prometheus 格式"""
        lines = []
        
        # CPU 指标
        lines.append(f"# HELP cpu_percent CPU usage percentage")
        lines.append(f"# TYPE cpu_percent gauge")
        lines.append(f"cpu_percent {metrics['cpu']['percent_total']}")
        
        for i, percent in enumerate(metrics['cpu']['percent_per_cpu']):
            lines.append(f'cpu_percent{{cpu="{i}"}} {percent}')
        
        lines.append(f"# HELP load_average System load average")
        lines.append(f"# TYPE load_average gauge")
        lines.append(f'load_average{{period="1m"}} {metrics["cpu"]["load_avg_1"]}')
        lines.append(f'load_average{{period="5m"}} {metrics["cpu"]["load_avg_5"]}')
        lines.append(f'load_average{{period="15m"}} {metrics["cpu"]["load_avg_15"]}')
        
        # 内存指标
        lines.append(f"# HELP memory_bytes Memory usage in bytes")
        lines.append(f"# TYPE memory_bytes gauge")
        lines.append(f'memory_bytes{{type="total"}} {metrics["memory"]["total"]}')
        lines.append(f'memory_bytes{{type="used"}} {metrics["memory"]["used"]}')
        lines.append(f'memory_bytes{{type="available"}} {metrics["memory"]["available"]}')
        
        lines.append(f"# HELP memory_percent Memory usage percentage")
        lines.append(f"# TYPE memory_percent gauge")
        lines.append(f"memory_percent {metrics['memory']['percent']}")
        
        # 磁盘指标
        lines.append(f"# HELP disk_usage_bytes Disk usage in bytes")
        lines.append(f"# TYPE disk_usage_bytes gauge")
        for disk in metrics['disk']:
            mountpoint = disk['mountpoint'].replace('/', '_root') if disk['mountpoint'] == '/' else disk['mountpoint']
            lines.append(f'disk_usage_bytes{{mountpoint="{mountpoint}",type="total"}} {disk["total"]}')
            lines.append(f'disk_usage_bytes{{mountpoint="{mountpoint}",type="used"}} {disk["used"]}')
        
        # 网络指标
        lines.append(f"# HELP network_bytes Network bytes transferred")
        lines.append(f"# TYPE network_bytes counter")
        for net in metrics['network']:
            lines.append(f'network_bytes{{interface="{net["interface"]}",direction="sent"}} {net["bytes_sent"]}')
            lines.append(f'network_bytes{{interface="{net["interface"]}",direction="recv"}} {net["bytes_recv"]}')
        
        return '\n'.join(lines)
    
    def export_json(self, metrics: Dict) -> str:
        """导出为 JSON 格式"""
        return json.dumps(metrics, indent=2)
    
    def export_influxdb(self, metrics: Dict) -> List[str]:
        """导出为 InfluxDB 行协议格式"""
        lines = []
        timestamp = int(time.time() * 1000000000)  # 纳秒
        
        # CPU 指标
        lines.append(f'cpu,host=localhost percent={metrics["cpu"]["percent_total"]} {timestamp}')
        lines.append(f'load,host=localhost avg1={metrics["cpu"]["load_avg_1"]},avg5={metrics["cpu"]["load_avg_5"]},avg15={metrics["cpu"]["load_avg_15"]} {timestamp}')
        
        # 内存指标
        lines.append(f'memory,host=localhost total={metrics["memory"]["total"]},used={metrics["memory"]["used"]},percent={metrics["memory"]["percent"]} {timestamp}')
        
        # 磁盘指标
        for disk in metrics['disk']:
            lines.append(f'disk,host=localhost,mountpoint={disk["mountpoint"]} total={disk["total"]},used={disk["used"]},percent={disk["percent"]} {timestamp}')
        
        return lines
    
    def run_continuous(self, duration: Optional[int] = None, output_file: Optional[str] = None):
        """持续采集指标"""
        print(f"Starting continuous metrics collection (interval: {self.interval}s)")
        
        start_time = time.time()
        collection_count = 0
        
        try:
            while True:
                metrics = self.collect_system_metrics()
                collection_count += 1
                
                print(f"\n[{collection_count}] Collected metrics at {metrics['timestamp']}")
                print(f"  CPU: {metrics['cpu']['percent_total']:.1f}%")
                print(f"  Memory: {metrics['memory']['percent']:.1f}%")
                print(f"  Load: {metrics['cpu']['load_avg_1']:.2f}")
                
                if output_file:
                    self._append_to_file(output_file, metrics)
                
                # 检查是否达到持续时间
                if duration and (time.time() - start_time) >= duration:
                    print(f"\nCompleted {collection_count} collections in {duration}s")
                    break
                
                time.sleep(self.interval)
                
        except KeyboardInterrupt:
            print(f"\n\nStopped after {collection_count} collections")
    
    def _append_to_file(self, filename: str, metrics: Dict):
        """追加指标到文件"""
        with open(filename, 'a') as f:
            f.write(json.dumps(metrics) + '\n')


def main():
    """主函数"""
    parser = argparse.ArgumentParser(
        description='指标收集器 - 收集系统性能指标',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
示例:
  # 收集一次指标
  %(prog)s
  
  # 持续收集指标
  %(prog)s --continuous --interval 60
  
  # 导出为 Prometheus 格式
  %(prog)s --format prometheus
  
  # 保存到文件
  %(prog)s --output metrics.json
  
  # 持续收集并保存
  %(prog)s --continuous --interval 30 --output metrics.jsonl
"""
    )
    
    parser.add_argument(
        '-i', '--interval',
        type=int,
        default=60,
        help='采集间隔（秒）(默认: 60)'
    )
    parser.add_argument(
        '-c', '--continuous',
        action='store_true',
        help='持续采集模式'
    )
    parser.add_argument(
        '-d', '--duration',
        type=int,
        help='持续采集的总时长（秒）'
    )
    parser.add_argument(
        '-o', '--output',
        help='输出文件路径'
    )
    parser.add_argument(
        '-f', '--format',
        default='json',
        choices=['json', 'prometheus', 'influxdb'],
        help='输出格式 (默认: json)'
    )
    parser.add_argument(
        '--custom-command',
        action='append',
        help='自定义命令（可多次使用）'
    )
    parser.add_argument(
        '-v', '--version',
        action='version',
        version='%(prog)s 1.0.0'
    )
    
    args = parser.parse_args()
    
    try:
        collector = MetricsCollector(interval=args.interval)
        
        if args.continuous:
            collector.run_continuous(duration=args.duration, output_file=args.output)
        else:
            # 单次采集
            metrics = collector.collect_system_metrics()
            
            # 添加自定义指标
            if args.custom_command:
                metrics['custom'] = collector.collect_custom_metrics(args.custom_command)
            
            # 格式化输出
            if args.format == 'prometheus':
                output = collector.export_prometheus(metrics)
            elif args.format == 'influxdb':
                output = '\n'.join(collector.export_influxdb(metrics))
            else:
                output = collector.export_json(metrics)
            
            # 输出结果
            if args.output:
                with open(args.output, 'w') as f:
                    f.write(output)
                print(f"Metrics saved to: {args.output}")
            else:
                print(output)
        
        return 0
        
    except Exception as e:
        print(f"Error: {e}", file=sys.stderr)
        import traceback
        traceback.print_exc()
        return 1


if __name__ == '__main__':
    sys.exit(main())
