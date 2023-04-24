// create by chencanhua in 2023/4/24
package web

import (
	"bytes"
	"context"
	"html/template"
)

// TemplateEngine
// 单一职责原则：一个接口只做一件事
// 加入或者删除一个模板，是具体模板引擎的事情，和 Web 框架什么关系都没有
type TemplateEngine interface {
	// Render 渲染页面
	// data 是渲染页面所需要的数据
	Render(ctx context.Context, tplName string, data interface{}) ([]byte, error)
}

var _ TemplateEngine = &GoTemplateEngine{}

// GoTemplateEngine Go的一种模板引擎实现
type GoTemplateEngine struct {
	T *template.Template
}

func (g *GoTemplateEngine) Render(ctx context.Context, tplName string, data interface{}) ([]byte, error) {
	res := &bytes.Buffer{}
	err := g.T.ExecuteTemplate(res, tplName, data)
	return res.Bytes(), err
}
