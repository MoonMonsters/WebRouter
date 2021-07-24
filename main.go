package main

import (
	"WebRouter/kernel"
	"WebRouter/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	kernel.Load()
	routes.Load(r)

	r.Run()
}
