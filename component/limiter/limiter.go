package limiter

/**
接口限流
*/

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type Limiters struct {
	limiters *sync.Map
}

// 限流器
type Limiter struct {
	limiter *rate.Limiter
	lastGet time.Time
	key     string
}

var GlobalLimiters = &Limiters{
	limiters: &sync.Map{},
}

var once = sync.Once{}

// 获取限流器
// r: 限流规则
// b: 时间间隔, r + b规定了访问间隔
// key: 标志, 例如ip, 手机号等
func NewLimiter(r rate.Limit, b int, key string) *Limiter {

	// 获取限流器时, 开启定时任务, 定时清除无效的限流器
	once.Do(func() {
		go GlobalLimiters.clearLimiter()
	})

	keyLimiter := GlobalLimiters.getLimiter(r, b, key)
	return keyLimiter
}

// 判断是否允许访问
func (l *Limiter) Allow() bool {
	l.lastGet = time.Now()
	return l.limiter.Allow()
}

func (ls *Limiters) getLimiter(r rate.Limit, b int, key string) *Limiter {
	limiter, ok := ls.limiters.Load(key)
	if ok {
		return limiter.(*Limiter)
	}

	l := &Limiter{
		limiter: rate.NewLimiter(r, b),
		lastGet: time.Now(),
		key:     key,
	}
	ls.limiters.Store(key, l)
	return l
}

// 清除限流器
func (ls *Limiters) clearLimiter() {
	for {
		time.Sleep(time.Minute * 1)
		ls.limiters.Range(func(key, value interface{}) bool {
			if time.Now().Unix()-value.(*Limiter).lastGet.Unix() > 60 {
				ls.limiters.Delete(key)
			}

			return true
		})
	}
}
