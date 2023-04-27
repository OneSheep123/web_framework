// Package errmdl create by chencanhua in 2023/4/26
package errmdl

import (
	"fmt"
	"net/http"
	"web_framework/web"
)

type Builder struct {
}

func (*Builder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
					ctx.RespData = []byte(`出错拉`)
					ctx.RespStatusCode = http.StatusInternalServerError
				}
			}()
			next(ctx)
		}
	}
}
