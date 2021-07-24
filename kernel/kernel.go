package kernel

import (
	"WebRouter/context"
	"WebRouter/exceptions"
	"WebRouter/middlewares"
)

var Middleware []context.HandlerFunc

func Load() {
	Middleware = []context.HandlerFunc{
		exceptions.Exception,
		middlewares.Session,
	}
}
