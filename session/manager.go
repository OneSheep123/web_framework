// Package session create by chencanhua in 2023/4/30
package session

import (
	"web_framework/web"
)

type Manager struct {
	Store
	Propagator
	// 用于优化使用，多次获取session时，记录id进行缓存使用
	CtxSessionId string
}

// InitSession 初始化session
func (m *Manager) InitSession(ctx *web.Context, id string) (Session, error) {
	s, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}
	// 写入到HTTP请求里面
	if err = m.Inject(id, ctx.Resp); err != nil {
		return nil, err
	}
	return s, nil
}

// GetSession 获取Session
// GetSession调用频率比较多的情况下，可以考虑进行缓存Session
func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}
	if s, ok := ctx.UserValues[m.CtxSessionId]; ok {
		return s.(Session), nil
	}
	sessionId, err := m.Extract(ctx.Req)
	if err != nil {
		return nil, err
	}
	s, err := m.Get(ctx.Req.Context(), sessionId)
	if err != nil {
		return nil, err
	}
	ctx.UserValues[m.CtxSessionId] = s
	return s, err
}

// RefreshSession 刷新session
func (m *Manager) RefreshSession(ctx *web.Context) error {
	s, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	// 刷新存储的过期时间
	err = m.Refresh(ctx.Req.Context(), s.ID())
	if err != nil {
		return err
	}
	// 重新注入 HTTP 里面(有可能刷新之后，session ID就发生改边了)
	if err = m.Inject(s.ID(), ctx.Resp); err != nil {
		return err
	}
	return nil
}

// RemoveSession 移除Session
func (m *Manager) RemoveSession(ctx *web.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), session.ID())
	if err != nil {
		return err
	}
	// 移除HTTP请求中的Session
	return m.Propagator.Remove(ctx.Resp)
}
