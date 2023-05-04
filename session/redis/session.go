// Package redis create by chencanhua in 2023/4/30
package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"web_framework/session"
)

var (
	errSessionNotFound = errors.New("session: id 对应的 session 不存在")
)

var _ session.Store = &Store{}

type Store struct {
	client     redis.Cmdable
	expiration time.Duration
	prefix     string
}

type StoreOptions func(s *Store)

func StoreSetExpiration(expiration time.Duration) StoreOptions {
	return func(s *Store) {
		s.expiration = expiration
	}
}

// NewStore 获取一个redis session使用
// client不是可选参数，是必选
func NewStore(client redis.Cmdable, opts ...StoreOptions) *Store {
	s := &Store{
		expiration: time.Minute * 15,
		client:     client,
		prefix:     "sessid",
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	key := redisKey(s.prefix, id)
	_, err := s.client.HSet(ctx, key, id, id).Result()
	if err != nil {
		return nil, err
	}
	_, err = s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return nil, err
	}
	return &redisSession{
		id:     id,
		key:    key,
		client: s.client,
	}, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	key := redisKey(s.prefix, id)
	_, err := s.client.Expire(ctx, key, s.expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	key := redisKey(s.prefix, id)
	_, err := s.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	key := redisKey(s.prefix, id)
	// 注意：并不是把整个session放近redis里面
	result, err := s.client.Exists(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	if result != 1 {
		return nil, errSessionNotFound
	}
	return &redisSession{
		client: s.client,
		id:     id,
		key:    key,
	}, nil
}

var _ session.Session = &redisSession{}

type redisSession struct {
	client redis.Cmdable
	// session id
	id string
	// 加了前缀的session id
	key string
}

func (r *redisSession) Get(ctx context.Context, key string) (any, error) {
	return r.client.HGet(ctx, r.key, key).Result()
}

func (r *redisSession) Set(ctx context.Context, key string, val any) error {
	const lua = `
if redis.call("exists", KEYS[1])
then
	return redis.call("hset", KEYS[1], ARGV[1], ARGV[2])
else
	return -1
end
`
	res, err := r.client.Eval(ctx, lua, []string{r.key}, key, val).Int()
	if err != nil {
		return err
	}
	if res < 0 {
		return errSessionNotFound
	}
	return nil
}

func (r *redisSession) ID() string {
	return r.id
}

func redisKey(prefix, id string) string {
	return fmt.Sprintf("%s-%s", prefix, id)
}
