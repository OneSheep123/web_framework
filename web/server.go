package web

import (
	"fmt"
	"net/http"
)

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler
	// Start 启动服务器
	// addr 是监听地址。如果只指定端口，可以使用 ":8081"
	// 或者 "localhost:8082"
	Start(addr string) error

	// addRoute 注册一个路由
	// method 是 HTTP 方法
	addRoute(method string, path string, handler HandleFunc)
	// 我们并不采取这种设计方案
	// addRoute(method string, path string, handlers... HandleFunc)
}

// 确保 HTTPServer 肯定实现了 Server 接口
var _ Server = &HTTPServer{}

type HTTPServer struct {
	router
	mdls []Middleware
	Log  func(log string)
}

type HTTPServerOptions func(server *HTTPServer)

func ServerAddMiddleware(mdls ...Middleware) HTTPServerOptions {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}

func NewHTTPServer(options ...HTTPServerOptions) *HTTPServer {
	server := &HTTPServer{
		router: newRouter(),
		Log: func(log string) {
			fmt.Println(log)
		},
	}
	for _, opt := range options {
		opt(server)
	}
	return server
}

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	root := s.serve
	for index := len(s.mdls) - 1; index >= 0; index-- {
		root = s.mdls[index](root)
	}

	var m = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			ctx.Resp.WriteHeader(ctx.RespStatusCode)
			ctx.Resp.Write(ctx.RespData)
		}
	}
	root = m(root)
	root(ctx)
}

// Start 启动服务器
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *HTTPServer) Post(path string, handler HandleFunc) {
	s.addRoute(http.MethodPost, path, handler)
}

func (s *HTTPServer) Get(path string, handler HandleFunc) {
	s.addRoute(http.MethodGet, path, handler)
}

func (s *HTTPServer) serve(ctx *Context) {
	mi, ok := s.findRoute(ctx.Req.Method, ctx.Req.URL.Path)
	if !ok || mi.n == nil || mi.n.handler == nil {
		ctx.RespStatusCode = 404
		ctx.RespData = []byte(`NOT FOUND`)
		return
	}
	ctx.PathParams = mi.pathParams
	ctx.MatchRoute = mi.n.matchRoute
	mi.n.handler(ctx)
}
