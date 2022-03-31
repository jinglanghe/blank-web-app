package cache

import (
	"context"
	"github.com/apulis/bmod/aistudio-aom/configs"
	"github.com/apulis/sdk/go-utils/logging"
	"github.com/go-redis/redis/v8"
	"time"
)

// Redis cache implement
type Redis struct {
	ctx    context.Context
	client *redis.Client
}

// Connect connect to redis
func (r *Redis) Connect() {
	r.ctx = context.Background()
	redisConf := configs.Config.Redis
	r.client = redis.NewClient(&redis.Options{
		Addr:     redisConf.Host,
		Password: redisConf.Auth,
		DB:       redisConf.Database,
	})
	_, err := r.client.Ping(r.ctx).Result()
	if err != nil {
		logging.Fatal().Err(err).Msg("Could not connected to redis")
	}
	logging.Debug().Msg("Successfully connected to redis")
}

// Get from key
func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// Set value with key and expire time
func (r *Redis) Set(key string, val string, expire int) error {
	return r.client.Set(r.ctx, key, val, time.Duration(expire)).Err()
}

// Del delete key in redis
func (r *Redis) Del(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// HashGet from key
func (r *Redis) HashGet(hk, key string) (string, error) {
	return r.client.HGet(r.ctx, hk, key).Result()
}

// HashDel delete key in specify redis's hashtable
func (r *Redis) HashDel(hk, key string) error {
	return r.client.HDel(r.ctx, hk, key).Err()
}

// Increase increase value
func (r *Redis) Increase(key string) error {
	return r.client.Incr(r.ctx, key).Err()
}

// Expire set ttl
func (r *Redis) Expire(key string, dur time.Duration) error {
	return r.client.Expire(r.ctx, key, dur).Err()
}
