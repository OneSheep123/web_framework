package web

import (
	"encoding/json"
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

	s.Start(":8081")
}

func TestMiddlerware(t *testing.T) {
	var logFunc = func(log string) {
		fmt.Println(log)
	}
	var logMiddleware = func() Middleware {
		return func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				defer func() {
					l := struct {
						Host string `json:"host,omitempty"`
						// 命中的路由
						Route      string `json:"route,omitempty"`
						HTTPMethod string `json:"http_method,omitempty"`
						Path       string `json:"path,omitempty"`
					}{
						Host:       ctx.Req.Host,
						Route:      ctx.MatchRoute,
						HTTPMethod: ctx.Req.Method,
						Path:       ctx.Req.URL.Path,
					}
					data, _ := json.Marshal(l)
					logFunc(string(data))
				}()
				next(ctx)
			}
		}
	}
	var testMiddleWare = func(str string) Middleware {
		return func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println(str + "头")
				next(ctx)
				fmt.Println(str + "尾")
			}
		}
	}
	server := NewHTTPServer(ServerAddMiddleware(logMiddleware()))
	server.Use(http.MethodGet, "/", testMiddleWare("/"))
	server.Use(http.MethodGet, "/a", testMiddleWare("a"))
	server.Use(http.MethodGet, "/b", testMiddleWare("b"))
	server.Use(http.MethodGet, "/a/b", testMiddleWare("ab"))
	server.Use(http.MethodGet, "/a/c", testMiddleWare("ac"))
	server.Use(http.MethodGet, "/b/d", testMiddleWare("bd"))

	server.Get("/a/b", func(ctx *Context) {
		fmt.Println("我是核心内容")
		ctx.RespStatusCode = 200
		ctx.RespData = []byte(`OK`)
	})

	server.Start(":8081")
}
