package gee

import (
	"fmt"
)

func Recovery() HandlerFunc {
	return func(ctx *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				ctx.String(500 , "internal server error:%s", message)
			}
		}()
		ctx.Next()
	}
}
