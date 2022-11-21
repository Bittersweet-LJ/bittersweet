package redis

import (
	"bittersweet/settings"
	"fmt"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
)

////声明一个全局rdb变量
//var rdb *redis.Client

var (
	client *redis.Client
	Nil    = redis.Nil
)

func Init(cfg *settings.RedisConfig) (err error) {
	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err = client.Ping().Result()
	if err != nil {
		zap.L().Error("connect redis failed", zap.Error(err))
		return
	}
	return
}

func Close() {
	_ = client.Close()
}
