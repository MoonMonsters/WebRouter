package lock

import (
	"WebRouter/redis"
	"context"
	goredis "github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Lock struct {
	key        string        //分布式锁的名称
	expiration time.Duration // 持锁时间
	requestId  string        // 当前持锁的请求id, 释放锁时, 该id需要保持一致
}

// 创建锁
// key: 标识
// expiration: 过期时间
// requestId: 请求id
func NewLock(key string, expiration time.Duration) *Lock {
	requestId := uuid.NewV4().String()
	return &Lock{
		key:        key,
		expiration: expiration,
		requestId:  requestId,
	}
}

// 获取锁
func (lk *Lock) Get() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	ok, err := redis.Client().SetNX(ctx, lk.key, lk.requestId, lk.expiration).Result()
	if err != nil {
		return false
	}

	return ok
}

// 释放锁
func (lk *Lock) Release() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	const luaScript = `
		if redis.call('get', KEYS[1])==ARGV[1] then
			return redis.call('del', KEYS[1])
		else
			return 0
		end
	`
	script := goredis.NewScript(luaScript)
	_, err := script.Run(ctx, redis.Client(), []string{lk.key}, lk.requestId).Result()
	return err
}

// 阻塞获取锁
func (lk *Lock) Block(expiration time.Duration) bool {
	t := time.Now()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		ok, err := redis.Client().SetNX(ctx, lk.key, lk.requestId, lk.expiration).Result()
		cancel()

		if err != nil {
			return false
		}

		if ok {
			return true
		}

		time.Sleep(time.Microsecond * 200)
		if time.Now().Sub(t) > expiration {
			return false
		}
	}
}

// 强制释放锁
func (lk *Lock) ForceRelease() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	_, err := redis.Client().Del(ctx, lk.key).Result()
	return err
}
