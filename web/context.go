package web

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	RespStatusCode int
	RespData       []byte
	MatchRoute     string

	// 页面渲染的引擎
	tplEngine TemplateEngine

	cacheQueryValues url.Values

	// 用户可以自由决定在这里存储什么，比如存储Session等
	// 主要用于解决在不同 Middleware 之间数据传递的问题
	// 但是要注意
	// 1. UserValues 在初始状态的时候总是 nil，你需要自己手动初始化
	UserValues map[string]any
}

// Render 渲染静态数据
func (c *Context) Render(tpl string, data any) error {
	var err error
	c.RespData, err = c.tplEngine.Render(c.Req.Context(), tpl, data)
	c.RespStatusCode = 200
	if err != nil {
		c.RespStatusCode = 500
	}
	return err
}

func (c *Context) QueryValue(key string) StringValue {
	if c.cacheQueryValues == nil {
		c.cacheQueryValues = c.Req.URL.Query()
	}
	vals, ok := c.cacheQueryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: 找不到这个 key")}
	}
	return StringValue{val: vals[0]}
}

type StringValue struct {
	val string
	err error
}

func (s StringValue) String() (string, error) {
	return s.val, s.err
}

func (s StringValue) ToInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
