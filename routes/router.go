package routes

import (
	"WebRouter/context"
	"WebRouter/response"
	"github.com/gin-gonic/gin"
)

type router struct {
	engine *gin.Engine
}

type group struct {
	engine      *gin.Engine
	path        string
	middlewares []context.HandlerFunc
}

type method int

const (
	GET    method = 0x000000
	POST   method = 0x000001
	PUT    method = 0x000002
	DELETE method = 0x000003
	ANY    method = 0x000004
)

func newRouter(engine *gin.Engine) *router {
	return &router{
		engine: engine,
	}
}

func (r *router) Group(path string, callback func(group), middlewares ...context.HandlerFunc) {
	callback(group{
		engine:      r.engine,
		path:        path,
		middlewares: middlewares,
	})
}

func (g group) Group(path string, callback func(group), middlewares ...context.HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
	g.path += path
	callback(g)
}

func (g group) Registered(method method, url string, action func(ctx *context.Context) *response.Response, middlewares ...context.HandlerFunc) {
	var handlers = make([]gin.HandlerFunc, len(g.middlewares)+len(middlewares)+1)
	g.middlewares = append(g.middlewares, middlewares...)
	for key, middleware := range g.middlewares {
		temp := middleware
		handlers[key] = func(c *gin.Context) {
			temp(&context.Context{Context: c})
		}
	}
	handlers[len(g.middlewares)] = convert(action)

	finalUrl := g.path + url
	switch method {
	case GET:
		g.engine.GET(finalUrl, handlers...)
	case POST:
		g.engine.POST(finalUrl, handlers...)
	case PUT:
		g.engine.PUT(finalUrl, handlers...)
	case DELETE:
		g.engine.DELETE(finalUrl, handlers...)
	case ANY:
		g.engine.Any(finalUrl, handlers...)
	}
}

// 函数转化
// 强制要求函数返回Response数据 -> gin.HandlerFunc
// 避免忘记返回值
func convert(f func(*context.Context) *response.Response) gin.HandlerFunc {
	return func(c *gin.Context) {

		// 调用实际路由函数
		// 返回了Response数据
		resp := f(&context.Context{Context: c})
		data := resp.GetData()
		// 判断返回数据类型
		switch item := data.(type) {
		case string:
			c.String(200, item)
		case gin.H:
			c.JSON(200, item)
		}
	}
}
