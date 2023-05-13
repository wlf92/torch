package database

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/wlf92/torch/internal/launch"
)

// NewMySQL 使用给定的选项创建一个新的 gorm 数据库实例.
func NewRedis(list ...int) (redis.UniversalClient, error) {
	if len(list) == 0 {
		list = append(list, 0)
	}

	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      launch.Config.Rds.Addrs,
		DB:         list[0],
		Username:   launch.Config.Rds.Username,
		Password:   launch.Config.Rds.Password,
		MaxRetries: launch.Config.Rds.MaxRetries,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
