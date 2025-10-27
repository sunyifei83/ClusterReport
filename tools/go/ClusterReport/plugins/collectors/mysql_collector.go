package collectors

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLCollector MySQL 数据采集器插件
type MySQLCollector struct {
	name   string
	config MySQLConfig
	db     *sql.DB
}

// MySQLConfig MySQL 配置
type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

// MySQLMetrics MySQL 指标
type MySQLMetrics struct {
	Timestamp        time.Time         `json:"timestamp"`
	Status           map[string]string `json:"status"`
	Connections      int64             `json:"connections"`
	MaxConnections   int64             `json:"max_connections"`
	ThreadsConnected int64             `json:"threads_connected"`
	ThreadsRunning   int64             `json:"threads_running"`
	QPS              float64           `json:"qps"`
	TPS              float64           `json:"tps"`
	SlowQueries      int64             `json:"slow_queries"`
	InnoDBStats      InnoDBStats       `json:"innodb_stats"`
}

// InnoDBStats InnoDB 统计信息
type InnoDBStats struct {
	BufferPoolSize    int64   `json:"buffer_pool_size"`
	BufferPoolUsed    int64   `json:"buffer_pool_used"`
	BufferPoolHitRate float64 `json:"buffer_pool_hit_rate"`
	RowsRead          int64   `json:"rows_read"`
	RowsInserted      int64   `json:"rows_inserted"`
	RowsUpdated       int64   `json:"rows_updated"`
	RowsDeleted       int64   `json:"rows_deleted"`
}

// NewMySQLCollector 创建 MySQL 采集器
func NewMySQLCollector(config MySQLConfig) (*MySQLCollector, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping MySQL: %w", err)
	}

	return &MySQLCollector{
		name:   "mysql-collector",
		config: config,
		db:     db,
	}, nil
}

// Name 返回采集器名称
func (c *MySQLCollector) Name() string {
	return c.name
}

// Collect 采集 MySQL 指标
func (c *MySQLCollector) Collect() (interface{}, error) {
	metrics := &MySQLMetrics{
		Timestamp: time.Now(),
		Status:    make(map[string]string),
	}

	// 采集全局状态
	if err := c.collectGlobalStatus(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect global status: %w", err)
	}

	// 采集 InnoDB 统计
	if err := c.collectInnoDBStats(metrics); err != nil {
		return nil, fmt.Errorf("failed to collect InnoDB stats: %w", err)
	}

	return metrics, nil
}

// collectGlobalStatus 采集全局状态
func (c *MySQLCollector) collectGlobalStatus(metrics *MySQLMetrics) error {
	rows, err := c.db.Query("SHOW GLOBAL STATUS")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name, value string
		if err := rows.Scan(&name, &value); err != nil {
			continue
		}
		metrics.Status[name] = value

		// 解析关键指标
		switch name {
		case "Connections":
			fmt.Sscanf(value, "%d", &metrics.Connections)
		case "Max_used_connections":
			fmt.Sscanf(value, "%d", &metrics.MaxConnections)
		case "Threads_connected":
			fmt.Sscanf(value, "%d", &metrics.ThreadsConnected)
		case "Threads_running":
			fmt.Sscanf(value, "%d", &metrics.ThreadsRunning)
		case "Slow_queries":
			fmt.Sscanf(value, "%d", &metrics.SlowQueries)
		}
	}

	return nil
}

// collectInnoDBStats 采集 InnoDB 统计
func (c *MySQLCollector) collectInnoDBStats(metrics *MySQLMetrics) error {
	rows, err := c.db.Query("SHOW ENGINE INNODB STATUS")
	if err != nil {
		return err
	}
	defer rows.Close()

	// 简化实现：实际需要解析 InnoDB STATUS 输出
	metrics.InnoDBStats = InnoDBStats{
		BufferPoolSize: 0,
		BufferPoolUsed: 0,
	}

	return nil
}

// Close 关闭连接
func (c *MySQLCollector) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}
