package redisc

import (
	"context"
	"fmt"
	"ginchat/config"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

//	Init() redis连接初始化
func Init(c config.Config) *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port),
		Password: c.Redis.Password,
		DB:       c.Redis.DB,
		PoolSize: c.Redis.PoolSize,
		MinIdleConns: c.Redis.MinIdleConn,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	return rdb
}
