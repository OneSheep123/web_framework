// Package memory create by chencanhua in 2023/4/30
package memory

import (
	"context"
	"fmt"
	cache "github.com/patrickmn/go-cache"
	"sync"
	"time"
	"web_framework/session"
)

var _ session.Store = &Store{}

type Store struct {
	sessions   *cache.Cache
	expiration time.Duration
}

// NewStore 创建一个 Store 的实例
// 实际上，这里也可以考虑使用 Option 设计模式，允许用户控制过期检查的间隔
func NewStore(expiration time.Duration) *Store {
	return &Store{
		sessions: cache.New(expiration, time.Second),
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	m := &memorySession{
		sessionId: id,
		values:    map[string]any{},
	}
	s.sessions.Set(m.ID(), m, s.expiration)
	return m, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	m, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	s.sessions.Set(m.ID(), m, s.expiration)
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.sessions.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	m, ok := s.sessions.Get(id)
	if !ok {
		return nil, fmt.Errorf("当前 %s 对应的session不存在", id)
	}
	// 注意：这里是*memorySession实现了session.Session接口
	return m.(*memorySession), nil
}

var _ session.Session = &memorySession{}

type memorySession struct {
	sessionId string
	values    map[string]any
	sync.RWMutex
}

func (m *memorySession) Get(ctx context.Context, key string) (any, error) {
	m.RLock()
	defer m.RUnlock()
	val, ok := m.values[key]
	if !ok {
		return nil, fmt.Errorf("当前 %s 不存在值", key)
	}
	return val, nil
}

func (m *memorySession) Set(ctx context.Context, key string, val any) error {
	m.Lock()
	defer m.Unlock()
	if m.values == nil {
		m.values = map[string]any{}
	}
	m.values[key] = val
	return nil
}

func (m *memorySession) ID() string {
	return m.sessionId
}
