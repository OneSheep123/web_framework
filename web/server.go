// create by chencanhua in 2023/04/12

package web

import (
	"log"
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
	log  func(msg string, args ...interface{})
}

type HTTPServerOptions func(s *HTTPServer)

// ServerWithMiddleware 注册中间件
func ServerWithMiddleware(mdls ...Middleware) HTTPServerOptions {
	return func(s *HTTPServer) {
		s.mdls = mdls
	}
}

func NewHTTPServer(options ...HTTPServerOptions) *HTTPServer {
	httpServer := &HTTPServer{
		router: newRouter(),
		// 设置默认log函数
		log: func(msg string, args ...interface{}) {
			log.Fatalln(msg)
		},
	}
	for _, opt := range options {
		opt(httpServer)
	}
	return httpServer
}

// ServeHTTP HTTPServer 处理请求的入口
func (s *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:  request,
		Resp: writer,
	}

	root := s.serve
	// 这里并还没开始执行，只是从后面开始遍历，注册前面一个middleware中的next
	for i := len(s.mdls) - 1; i >= 0; i-- {
		root = s.mdls[i](root)
	}
	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			s.flashResp(ctx)
		}
	}
	root = m(root)
	// 这里就是开始从前往后执行了
	root(ctx)
}

func (s *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode != 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	if ctx.Resp == nil {
		s.log("当前Resp为nil")
		return
	}
	n, err := ctx.Resp.Write(ctx.RespData)
	if err != nil || n != len(ctx.RespData) {
		s.log("写入响应失败 %v", err)
	}
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
		ctx.Resp.WriteHeader(404)
		ctx.Resp.Write([]byte("Not Found"))
		return
	}
	ctx.PathParams = mi.pathParams
	ctx.MatchRoute = mi.n.matchRoute
	mi.n.handler(ctx)
}
