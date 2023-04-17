// create by chencanhua in 2023/4/17
package logging

import (
	"fmt"
	"testing"
	"web_framework/web"
)

func TestMiddlewareBuilder(t *testing.T) {
	build := NewBuilder(func(log string) {
		fmt.Println(log)
	})
	server := web.NewHTTPServer(web.ServerWithMiddleware(build.Build()))
	server.Get("/a/b/*", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, it's me"))
	})
	server.Start(":8081")
}
