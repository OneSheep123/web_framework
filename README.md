# web_framework
## 使用

```go
 go mode tidy # 初始化模块
```

## 文件介绍 

- context.go：路由上下文, 包含了获取请求体、请求参数的方法
- route.go: 路由文件，内涵路由节点注册方法，支持静态路由、通配符路由、参数路由
- route_test.go: 路由测试文件
- server.go: 整体代表服务器的抽象
- server_test.go: web框架测试文件
- files.go: 实现文件上传下载的方法、静态文件下载方法
- middleware.go: 路由中间件抽象
- middleware(dir): 相关中间件的实现
- template.go: 模板引擎的实现

