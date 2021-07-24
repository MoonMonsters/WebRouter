package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type connect struct {
	client *redis.Client
}

var once = sync.Once{}

var _connect *connect

func connectRedis() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	conf := &redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	}
	c := redis.NewClient(conf)
	re := c.Ping(ctx)
	if re.Err() != nil {
		panic(re.Err())
	}
	_connect = &connect{
		client: c,
	}
}

func Client() *redis.Client {
	if _connect == nil {
		once.Do(func() {
			connectRedis()
		})
	}

	return _connect.client
}
