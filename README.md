# web_framework
## 使用

```go
go mod tidy # 导入模块
```

## 文件介绍 

- context.go：路由上下文
- route.go：路由文件，内涵路由节点注册方法，支持静态路由、通配符路由、参数路由
- route_test.go: 路由测试文件
- server.go：整体代表服务器的抽象
- server_test.go: web框架测试文件
- middleware.go: Middleware 的 类型
- middleware文件夹：可以进行注册的中间件

## Middleware介绍

### AccessLog

logging文件夹下，记录所有进来的请求

### Tracing 

> Tracing：**链路追逐**，它记录从收到请求到返回响应的 整个过程。在分布式环境下，**一般代表着请求从 Web 收到，沿着后续微服务链条传递，得到响应 再返回到前端的过程**。业界上，常用的tracing工具有 SkyWalking、Zipkin、Jeager 等。 

opentelemetry文件夹下。定义一个统一的 API，允许用户注入自己的 tracing 实现。

