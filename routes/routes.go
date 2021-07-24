package routes

import (
	"WebRouter/controller"
	"WebRouter/kernel"
	"WebRouter/middlewares"
	"github.com/gin-gonic/gin"
)

func Load(r *gin.Engine) {

	router := newRouter(r)
	router.Group("", func(g group) {
		config(g)
	}, kernel.Middleware...)
}

func config(router group) {

	router.Registered(GET, "/", controller.Index)

	router.Group("/api", func(api group) {
		api.Group("/user", func(user group) {
			user.Registered(GET, "/info", controller.Index)
		}, middlewares.M2)
		api.Group("/host", func(host group) {
			host.Registered(GET, "/info", controller.Host)
		})
	})
}
