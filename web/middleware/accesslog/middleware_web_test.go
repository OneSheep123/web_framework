// create by chencanhua in 2023/4/20
package accesslog

import (
	"fmt"
	"testing"
	"web_framework/web"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := Builder{}
	mdl := builder.SetLogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web.NewHTTPServer(web.ServerAddMiddleware(mdl))
	server.Get("/a/b/*", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, it's me"))
	})
	server.Start(":8081")
}
