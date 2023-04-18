package dbrds

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/wlf92/torch/internal/launch"
)

// NewMySQL 使用给定的选项创建一个新的 gorm 数据库实例.
func NewRedis(dbIdx int) (redis.UniversalClient, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:      launch.Config.Rds.Addrs,
		DB:         dbIdx,
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
