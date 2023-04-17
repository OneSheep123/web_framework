// Package logging create by chencanhua in 2023/4/17
package logging

import (
	"encoding/json"
	"web_framework/web"
)

type Builder struct {
	LogFunc func(log string)
}

func NewBuilder(logFunc func(log string)) *Builder {
	return &Builder{
		LogFunc: logFunc,
	}
}

func (b *Builder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 要记录请求
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				b.LogFunc(string(data))
			}()
			next(ctx)
		}
	}
}

type accessLog struct {
	Host string `json:"host,omitempty"`
	// 命中的路由
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
