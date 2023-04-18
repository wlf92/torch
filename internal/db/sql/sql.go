package dbsql

import (
	"time"

	"github.com/wlf92/torch/internal/launch"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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
