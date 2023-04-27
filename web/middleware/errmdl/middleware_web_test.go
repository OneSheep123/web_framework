// create by chencanhua in 2023/4/20
package errmdl

import (
	"testing"
	"web_framework/web"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := Builder{}
	server := web.NewHTTPServer(web.ServerAddMiddleware(builder.Build()))
	server.Get("/user", func(ctx *web.Context) {
		panic("我是错误")
	})
	server.Start(":8081")
}
