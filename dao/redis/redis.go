package redis

import (
	"WebApp/settings"
	"fmt"

	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func Close() {
	_ = rdb.Close()
}
