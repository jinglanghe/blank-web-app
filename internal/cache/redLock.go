package cache

import (
	"gitlab.apulis.com.cn/hjl/blank-web-app-2/configs"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
)

var RedLock *redsync.Redsync

func InitRedLock() {
	redisConf := configs.Config.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConf.Host,
		Password: redisConf.Auth,
		DB:       redisConf.Database,
	})
	pool := goredis.NewPool(client)

	RedLock = redsync.New(pool)
}
