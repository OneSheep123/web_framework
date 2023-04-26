// create by chencanhua in 2023/4/26
package errmdl

import (
	"fmt"
	"web_framework/web"
)

type Buidler struct {
}

func (*Buidler) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			next(ctx)
		}
	}
}
