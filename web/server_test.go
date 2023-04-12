// create by chencanhua in 2023/04/12

package web

import "testing"

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Get("/user/:((a|b|c)name)", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, reg"))
	})

	s.Start(":8081")
}
