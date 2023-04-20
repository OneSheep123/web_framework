// Package accesslog create by chencanhua in 2023/4/20
package accesslog

import (
	"encoding/json"
	"web_framework/web"
)

type Builder struct {
	logFunc func(log string)
}

func (b *Builder) SetLogFunc(logFunc func(log string)) *Builder {
	b.logFunc = logFunc
	return b
}

func (b *Builder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				b.logFunc(string(data))
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
