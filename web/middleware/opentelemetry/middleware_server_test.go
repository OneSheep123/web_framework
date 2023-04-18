// create by chencanhua in 2023/4/17
package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"testing"
	"time"
	"web_framework/web"
)

func TestMiddlewareBuilder(t *testing.T) {
	tracer := otel.GetTracerProvider().Tracer(instrumentationName)
	// 初始化zipkin
	initZipkin(t)
	builder := Builder{Tracer: tracer}
	server := web.NewHTTPServer(web.ServerWithMiddleware(builder.Build()))

	server.Get("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	server.Get("/user", func(ctx *web.Context) {
		c, span := tracer.Start(ctx.Req.Context(), "first_layer")
		defer span.End()

		c, second := tracer.Start(c, "second_layer")
		time.Sleep(time.Second)
		defer second.End()

		c, third1 := tracer.Start(c, "third_layer_1")
		time.Sleep(100 * time.Millisecond)
		third1.End()

		c, third2 := tracer.Start(c, "third_layer_1")
		time.Sleep(300 * time.Millisecond)
		third2.End()

		ctx.RespStatusCode = 200
		ctx.RespData = []byte("hello, aiweier")
	})

	server.Start(":8081")
}
