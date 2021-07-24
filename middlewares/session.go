package middlewares

import (
	"WebRouter/context"
	"WebRouter/redis"
	context2 "context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"
)

var cookieName = "_gin"

var lifeTime = 3600

func Session(c *context.Context) {

	cookie, err := c.Cookie(cookieName)
	if err == nil {
		sessionString, err := redis.Client().Get(context2.TODO(), cookie).Result()
		if err == nil {
			var session context.Session
			json.Unmarshal([]byte(sessionString), &session)
			c.Set("_session", session)
			return
		}
	}

	sessionKey := uuid.NewV4().String()
	c.SetCookie(cookieName, sessionKey, 3600, "/", c.Domain(), false, true)

	session := context.Session{
		Cookie:      sessionKey,
		ExpireTime:  time.Now().Unix() + int64(lifeTime),
		SessionList: make(map[string]interface{}),
	}
	c.Set("_session", session)
	jsonString, _ := json.Marshal(session)
	redis.Client().Set(context2.TODO(), sessionKey, jsonString, time.Second*time.Duration(lifeTime))
}
