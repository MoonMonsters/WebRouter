package middlewares

import (
	"WebRouter/context"
	"fmt"
)

func M1(c *context.Context) {
	fmt.Println("我是M1")
}

func M2(c *context.Context) {
	fmt.Println("我是M2")
}

func M3(c *context.Context) {
	fmt.Println("我是M3")
}
