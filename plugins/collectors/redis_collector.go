package collectors

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// RedisCollector Redis 数据采集器插件示例
type RedisCollector struct {
	name   string
	config RedisConfig
	// 实际使用时需要 redis client: client *redis.Client
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// RedisMetrics Redis 指标
type RedisMetrics struct {
	Timestamp              time.Time             `json:"timestamp"`
	Info                   map[string]string     `json:"info"`
	ConnectedClients       int64                 `json:"connected_clients"`
	UsedMemory             int64                 `json:"used_memory"`
	UsedMemoryPeak         int64                 `json:"used_memory_peak"`
	UsedMemoryRSS          int64                 `json:"used_memory_rss"`
	MemFragmentationRatio  float64               `json:"mem_fragmentation_ratio"`
	TotalCommandsProcessed int64                 `json:"total_commands_processed"`
	InstantaneousOpsPerSec int64                 `json:"instantaneous_ops_per_sec"`
	KeyspaceHits           int64                 `json:"keyspace_hits"`
	KeyspaceMisses         int64                 `json:"keyspace_misses"`
	EvictedKeys            int64                 `json:"evicted_keys"`
	ExpiredKeys            int64                 `json:"expired_keys"`
	Keyspace               map[string]KeyspaceDB `json:"keyspace"`
}

// KeyspaceDB Redis 数据库键空间信息
type KeyspaceDB struct {
	Keys    int64 `json:"keys"`
	Expires int64 `json:"expires"`
	AvgTTL  int64 `json:"avg_ttl"`
}

// NewRedisCollector 创建 Redis 采集器
func NewRedisCollector(config RedisConfig) (*RedisCollector, error) {
	// 实际使用时需要创建 Redis 客户端
	// client := redis.NewClient(&redis.Options{
	//     Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
	//     Password: config.Password,
	//     DB:       config.DB,
	// })
	//
	// if err := client.Ping(context.Background()).Err(); err != nil {
	//     return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	// }

	return &RedisCollector{
		name:   "redis-collector",
		config: config,
		// client: client,
	}, nil
}

// Name 返回采集器名称
func (c *RedisCollector) Name() string {
	return c.name
}

// Collect 采集 Redis 指标
func (c *RedisCollector) Collect() (interface{}, error) {
	metrics := &RedisMetrics{
		Timestamp: time.Now(),
		Info:      make(map[string]string),
		Keyspace:  make(map[string]KeyspaceDB),
	}

	// 采集 INFO 命令输出
	if err := c.collectInfo(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect info: %w", err)
	}

	return metrics, nil
}

// collectInfo 采集 Redis INFO
func (c *RedisCollector) collectInfo(metrics *RedisMetrics) error {
	// 实际使用时调用 Redis INFO 命令
	// info, err := c.client.Info(context.Background()).Result()
	// if err != nil {
	//     return err
	// }

	// 这里使用示例数据演示解析逻辑
	info := c.getMockInfo()

	lines := strings.Split(info, "\n")
	currentSection := ""

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "# ") {
				currentSection = strings.TrimSpace(strings.TrimPrefix(line, "# "))
			}
			continue
		}

		// 解析键值对
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		metrics.Info[key] = value

		// 解析关键指标
		c.parseMetric(metrics, key, value, currentSection)
	}

	return nil
}

// parseMetric 解析指标
func (c *RedisCollector) parseMetric(metrics *RedisMetrics, key, value, section string) {
	switch key {
	case "connected_clients":
		metrics.ConnectedClients, _ = strconv.ParseInt(value, 10, 64)
	case "used_memory":
		metrics.UsedMemory, _ = strconv.ParseInt(value, 10, 64)
	case "used_memory_peak":
		metrics.UsedMemoryPeak, _ = strconv.ParseInt(value, 10, 64)
	case "used_memory_rss":
		metrics.UsedMemoryRSS, _ = strconv.ParseInt(value, 10, 64)
	case "mem_fragmentation_ratio":
		metrics.MemFragmentationRatio, _ = strconv.ParseFloat(value, 64)
	case "total_commands_processed":
		metrics.TotalCommandsProcessed, _ = strconv.ParseInt(value, 10, 64)
	case "instantaneous_ops_per_sec":
		metrics.InstantaneousOpsPerSec, _ = strconv.ParseInt(value, 10, 64)
	case "keyspace_hits":
		metrics.KeyspaceHits, _ = strconv.ParseInt(value, 10, 64)
	case "keyspace_misses":
		metrics.KeyspaceMisses, _ = strconv.ParseInt(value, 10, 64)
	case "evicted_keys":
		metrics.EvictedKeys, _ = strconv.ParseInt(value, 10, 64)
	case "expired_keys":
		metrics.ExpiredKeys, _ = strconv.ParseInt(value, 10, 64)
	}

	// 解析 keyspace 信息
	if strings.HasPrefix(key, "db") {
		c.parseKeyspace(metrics, key, value)
	}
}

// parseKeyspace 解析键空间信息
func (c *RedisCollector) parseKeyspace(metrics *RedisMetrics, dbName, value string) {
	// 格式: keys=xxx,expires=xxx,avg_ttl=xxx
	parts := strings.Split(value, ",")
	db := KeyspaceDB{}

	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		switch kv[0] {
		case "keys":
			db.Keys, _ = strconv.ParseInt(kv[1], 10, 64)
		case "expires":
			db.Expires, _ = strconv.ParseInt(kv[1], 10, 64)
		case "avg_ttl":
			db.AvgTTL, _ = strconv.ParseInt(kv[1], 10, 64)
		}
	}

	metrics.Keyspace[dbName] = db
}

// getMockInfo 返回模拟的 Redis INFO 输出（仅用于示例）
func (c *RedisCollector) getMockInfo() string {
	return `# Server
redis_version:7.0.0
redis_mode:standalone
os:Linux 5.10.0-21-amd64 x86_64

# Clients
connected_clients:10
blocked_clients:0

# Memory
used_memory:1048576
used_memory_human:1.00M
used_memory_rss:2097152
used_memory_peak:1572864
mem_fragmentation_ratio:2.0

# Stats
total_connections_received:1000
total_commands_processed:50000
instantaneous_ops_per_sec:100
keyspace_hits:45000
keyspace_misses:5000
evicted_keys:10
expired_keys:100

# Keyspace
db0:keys=1000,expires=100,avg_ttl=3600000
db1:keys=500,expires=50,avg_ttl=7200000`
}

// Close 关闭连接
func (c *RedisCollector) Close() error {
	// 实际使用时关闭 Redis 客户端
	// if c.client != nil {
	//     return c.client.Close()
	// }
	return nil
}

// Ping 测试连接
func (c *RedisCollector) Ping(ctx context.Context) error {
	// 实际使用时调用 Redis PING
	// return c.client.Ping(ctx).Err()
	return nil
}
