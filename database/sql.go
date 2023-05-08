package database

import (
	"fmt"
	"time"

	"github.com/wlf92/torch/internal/launch"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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

// NewMySQL 使用给定的选项创建一个新的 gorm 数据库实例.
func NewMySQL(dbName string) (*gorm.DB, error) {
	cfg := launch.Config.GetSqlByDb(dbName)

	o := defaultOptions()
	o.Host = cfg.Host
	o.Username = cfg.Username
	o.Password = cfg.Password
	o.Database = cfg.Db
	o.MaxIdleConnections = cfg.MaxIdleConnections
	o.MaxOpenConnections = cfg.MaxOpenConnections
	o.MaxConnectionLifeTime = time.Duration(cfg.MaxConnectionLifeTime) * time.Second
	o.LogLevel = cfg.LogLevel

	logLevel := logger.Silent
	if o.LogLevel != 0 {
		logLevel = logger.LogLevel(o.LogLevel)
	}

	db, err := gorm.Open(mysql.Open(o.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns 设置到数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(o.MaxOpenConnections)

	// SetConnMaxLifetime 设置连接可重用的最长时间
	sqlDB.SetConnMaxLifetime(o.MaxConnectionLifeTime)

	// SetMaxIdleConns 设置空闲连接池的最大连接数
	sqlDB.SetMaxIdleConns(o.MaxIdleConnections)

	return db, nil
}
