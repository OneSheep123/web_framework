# web_framework

## 功能

- 提供HTTP Server服务，支持路由注册和路由匹配，实现基于前缀路由树的静态路由匹配、参数匹配、通配符*匹配
- 支持AOP方案，实现全局中间件和路由中间件的注册
- 支持文件上传和下载功能
- 支持模板引擎功能， 并提供基于 Go 模板的默认实现 
- 支持Session功能，并提供基于内存的实现和redis的实现

## 安装

```
go get -u github.com/OneSheep123/web_framework
```

## 使用

```go
package main

import (
	"web_framework/web"
)

func main() {
    s := web.NewHTTPServer()
	s.Get("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Start(":8081")
}
```

## 文件介绍 


```
web_framework
├─ file_demo
│  ├─ file_test.go
│  └─ testdata
│     └─ aivier.txt
├─ go.mod
├─ go.sum
├─ main.go
├─ README.md
├─ session
│  ├─ cookie
│  │  └─ propagator.go 
│  ├─ manager.go
│  ├─ memory
│  │  └─ session.go  
│  ├─ redis
│  │  └─ session.go
│  ├─ test
│  │  └─ types_test.go
│  └─ types.go
├─ template_demo
│  └─ template_test.go
└─ web
   ├─ context.go
   ├─ files.go
   ├─ middleware
   │  ├─ accesslog
   │  │  ├─ middleware.go
   │  │  └─ middleware_web_test.go
   │  └─ errmdl
   │     ├─ middleware.go
   │     └─ middleware_web_test.go
   ├─ middleware.go
   ├─ route.go
   ├─ route_test.go
   ├─ server.go
   ├─ template.go
   ├─ test
   │  └─ server_test.go
   └─ testdata
      ├─ download
      │  └─ 123.txt
      ├─ tpls
      │  └─ login.gohtml
      └─ upload
         └─ 下载.jfif
```

- context.go：路由上下文, 包含了获取请求体、请求参数的方法
- route.go: 路由文件，内涵路由节点注册方法，支持静态路由、通配符路由、参数路由
- server.go: 整体代表服务器的抽象
- files.go: 实现文件上传下载的方法、静态文件下载方法
- middleware.go: 路由中间件抽象
- template.go: 模板引擎的实现
- session / types.go: 规范session要实现的方法
- session / redis：基于redis的session实现
- session / memory：基于内存的session实现