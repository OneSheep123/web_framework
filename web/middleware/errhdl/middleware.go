// Package errhdl create by chencanhua in 2023/4/19
package errhdl

import "web_framework/web"

type Builder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() Builder {
	return Builder{
		resp: map[int][]byte{},
	}
}

func (m *Builder) AddCode(statusCode int, data []byte) *Builder {
	m.resp[statusCode] = data
	return m
}

func (m *Builder) Build() web.Middleware {
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			next(ctx)
			data, ok := m.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = data
			}
		}
	}
}
