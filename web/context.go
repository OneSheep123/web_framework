// create by chencanhua in 2023/04/12

package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
)

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	MatchRoute     string
	RespStatusCode int
	RespData       []byte

	queryValues url.Values
}

func (c *Context) QueryValue(key string) StringValue {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}
	val, ok := c.queryValues[key]
	if !ok {
		return StringValue{err: errors.New("web: key 不存在")}
	}
	return StringValue{
		val: val[0],
	}
}

func (c *Context) RespJSONOK(val interface{}) error {
	return c.RespJSON(http.StatusOK, val)
}

func (c *Context) RespJSON(code int, val interface{}) error {
	jsonStr, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.RespStatusCode = code
	c.RespData = jsonStr
	return nil
}

type StringValue struct {
	val string
	err error
}

func (s *StringValue) AsInt64() (int64, error) {
	if s.err != nil {
		return 0, s.err
	}
	return strconv.ParseInt(s.val, 10, 64)
}
