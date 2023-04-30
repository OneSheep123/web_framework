// create by chencanhua in 2023/4/30
package test

import (
	uuid2 "github.com/google/uuid"
	"net/http"
	"testing"
	"web_framework/session"
	"web_framework/web"
)

func TestSession(t *testing.T) {
	m := &session.Manager{}
	server := web.NewHTTPServer(web.ServerAddMiddleware(func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {
			// 进行一个权限校验
			if ctx.Req.URL.Path == "/login" {
				next(ctx)
				return
			}
			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.RespData = []byte(`请重新登录`)
				ctx.RespStatusCode = http.StatusUnauthorized
				return
			}
			// 刷新 session 的过期时间
			_ = m.RefreshSession(ctx)
			next(ctx)
		}
	}))
	// 登录注册
	server.Post("/login", func(ctx *web.Context) {
		uuid := uuid2.New().String()
		s, err := m.InitSession(ctx, uuid)
		if err != nil {
			ctx.RespData = []byte(`登录失败了`)
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		err = s.Set(ctx.Req.Context(), "zhangsan", "123")
		if err != nil {
			ctx.RespData = []byte(`登录失败了`)
			ctx.RespStatusCode = http.StatusInternalServerError
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登录成功")
		return
	})
	// 退出登录
	server.Post("/logout", func(ctx *web.Context) {
		// 清理各种数据
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("退出登录")
	})

	server.Get("/user", func(ctx *web.Context) {
		sess, _ := m.GetSession(ctx)
		// 假如说我要把昵称从 session 里面拿出来
		val, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(val.(string))
	})

	server.Start(":8081")
}
