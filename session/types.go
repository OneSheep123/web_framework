// Package session create by chencanhua in 2023/4/29
package session

import (
	"context"
	"net/http"
)

// Store 管理Session本身
// 另外，一些Api的设计都把context.Context和error给加上
// context可以用于控制上下文，error用于返回错误
type Store interface {
	// Generate 初始化一个Session
	Generate(ctx context.Context, id string) (Session, error)
	Refresh(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	Get(ctx context.Context, id string) (Session, error)
}

type Session interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, val any) error
	ID() string
}

// Propagator 给HTTP请求中设置或者移除session id
type Propagator interface {
	// Inject 写入session id到HTTP请求
	Inject(id string, writer http.ResponseWriter) error
	// Extract 从HTTP请求中获取session id
	Extract(req *http.Request) (string, error)
	// Remove 从HTTP请求中移除
	Remove(writer http.ResponseWriter) error
}
