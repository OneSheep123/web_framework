package web

import "net/http"

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string

	RespStatusCode int
	RespData       []byte
	MatchRoute     string

	// 页面渲染的引擎
	tplEngine TemplateEngine
}

// Render 渲染静态数据
func (c *Context) Render(tpl string, data interface{}) error {
	var err error
	c.RespData, err = c.tplEngine.Render(c.Req.Context(), tpl, data)
	c.RespStatusCode = 200
	if err != nil {
		c.RespStatusCode = 500
	}
	return err
}
