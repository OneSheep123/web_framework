// create by chencanhua in 2023/04/12

package web

import (
	"fmt"
	"net/http"
	"testing"
)

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

func TestHTTPServer_ServeHTTP(t *testing.T) {
	server := NewHTTPServer()
	server.mdls = []Middleware{
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第一个before")
				next(ctx)
				fmt.Println("第一个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第二个before")
				next(ctx)
				fmt.Println("第二个after")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第三个中断")
			}
		},
		func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println("第四个，你看不到这句话")
			}
		},
	}
	server.ServeHTTP(nil, &http.Request{})
}
