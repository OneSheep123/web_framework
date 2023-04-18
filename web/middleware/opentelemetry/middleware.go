// Package opentelemetry create by chencanhua in 2023/4/17
package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"web_framework/web"
)

const instrumentationName = "web_framework/web/middleware/opentelemetry"

type Builder struct {
	Tracer trace.Tracer
}

func (b *Builder) Build() web.Middleware {
	if b.Tracer == nil {
		b.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			reqCtx := ctx.Req.Context()
			// 尝试和客户端的 trace 结合在一起
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			reqCtx, span := b.Tracer.Start(reqCtx, "unknown")
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("peer.hostname", ctx.Req.Host))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.scheme", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("span.kind", "server"))
			span.SetAttributes(attribute.String("component", "web"))
			span.SetAttributes(attribute.String("peer.address", ctx.Req.RemoteAddr))
			span.SetAttributes(attribute.String("http.proto", ctx.Req.Proto))

			ctx.Req = ctx.Req.WithContext(reqCtx)

			next(ctx)

			// 这个是只有执行完 next 才可能有值
			if ctx.MatchRoute != "" {
				span.SetName(ctx.MatchRoute)
			}

			// 把响应码加上去
			span.SetAttributes(attribute.Int("http.status", ctx.RespStatusCode))
		}
	}
}
