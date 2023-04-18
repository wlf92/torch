package dbsql

import (
	"fmt"
	"time"
)

// Option 定义 MySQL 数据库的选项.
type Options struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	LogLevel              int
}

// DSN 从 Option 返回 DSN.
func (o *Options) DSN() string {
	return fmt.Sprintf(`%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s`, o.Username, o.Password, o.Host, o.Database, true, "Local")
}

type Option func(o *Options)

func defaultOptions() *Options {
	return &Options{
		MaxIdleConnections:    100,
		MaxOpenConnections:    100,
		MaxConnectionLifeTime: time.Duration(10) * time.Second,
		LogLevel:              2,
	}
}
