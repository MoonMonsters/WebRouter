package context

import (
	"WebRouter/redis"
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"sync"
	"time"
)

// 重写gin.Context, 方便扩展
type Context struct {
	*gin.Context
}

type Session struct {
	Cookie      string                 `json:"cookie"`
	ExpireTime  int64                  `json:"expire_time"`
	SessionList map[string]interface{} `json:"session_list"`
	Lock        *sync.Mutex            // session的访问, 需要加锁, session的操作是并发的
}

type HandlerFunc func(*Context)

// 获取域名
func (c *Context) Domain() string {
	return c.Request.Host[:strings.Index(c.Request.Host, ":")]
}

func (c *Context) Session() *Session {
	var session Session

	_session, ok := c.Get("_session")
	if !ok {
		return nil
	}

	session = _session.(Session)
	session.Lock = &sync.Mutex{}

	return &session
}

// 写入数据到session中
func (s *Session) Set(key string, value interface{}) error {
	// 加锁
	s.Lock.Lock()
	defer s.Lock.Unlock()

	sessionString, err := redis.Client().Get(context.TODO(), s.Cookie).Result()
	if err != nil {
		return err
	}

	var session Session
	err = json.Unmarshal([]byte(sessionString), &session)
	if err != nil {
		return err
	}

	// 写入具体数据, 重新序列化session
	session.SessionList[key] = value
	sessionStringNew, err := json.Marshal(session)
	e := s.ExpireTime - time.Now().Unix()
	if e < 0 {
		return errors.New("the session has expired")
	}

	redis.Client().Set(context.TODO(), s.Cookie, sessionStringNew, time.Duration(e)*time.Second)
	return nil
}

// 获取session数据
func (s *Session) Get(key string) (interface{}, error) {
	sessionString, err := redis.Client().Get(context.TODO(), s.Cookie).Result()
	if err != nil {
		return nil, err
	}

	var session Session
	err = json.Unmarshal([]byte(sessionString), &session)
	if err != nil {
		return nil, err
	}

	value, ok := session.SessionList[key]
	if ok {
		return value, nil
	}
	return nil, errors.New("not found key: " + key)
}

// 删除session数据
func (s *Session) Remove(key string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	sessionString, err := redis.Client().Get(context.TODO(), s.Cookie).Result()
	if err != nil {
		return err
	}

	var session Session
	err = json.Unmarshal([]byte(sessionString), &session)
	if err != nil {
		return err
	}
	delete(session.SessionList, key)
	sessionStringNew, err := json.Marshal(session)
	if err != nil {
		return err
	}

	e := s.ExpireTime - time.Now().Unix()
	if e < 0 {
		return errors.New("the session has expired")
	}
	redis.Client().Set(context.TODO(), s.Cookie, sessionStringNew, time.Duration(e)*time.Second)
	return nil
}
