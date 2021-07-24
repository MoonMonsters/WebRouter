package controller

import (
	"WebRouter/component/limiter"
	"WebRouter/context"
	"WebRouter/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"strings"
	"time"
)

func Index(ctx *context.Context) *response.Response {
	ctx.Session().Set("msg", "golang是最好的语言")
	ctx.Session().Set("hello", "php也是最好的语言")

	l := limiter.NewLimiter(rate.Every(time.Second*1), 1, ctx.ClientIP())
	if !l.Allow() {
		return response.Resp().String("您的访问过于频繁")
	}

	return response.Resp().Json(gin.H{"msg": "Hello Gin"})
}

func Host(ctx *context.Context) *response.Response {
	return response.Resp().String(ctx.Request.Host[:strings.Index(ctx.Request.Host, ":")])
}
