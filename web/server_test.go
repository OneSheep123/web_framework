package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"path"
	"testing"
)

func TestServer(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Start(":8081")
}

func TestMiddlerware(t *testing.T) {
	var logFunc = func(log string) {
		fmt.Println(log)
	}
	var logMiddleware = func() Middleware {
		return func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				defer func() {
					l := struct {
						Host string `json:"host,omitempty"`
						// 命中的路由
						Route      string `json:"route,omitempty"`
						HTTPMethod string `json:"http_method,omitempty"`
						Path       string `json:"path,omitempty"`
					}{
						Host:       ctx.Req.Host,
						Route:      ctx.MatchRoute,
						HTTPMethod: ctx.Req.Method,
						Path:       ctx.Req.URL.Path,
					}
					data, _ := json.Marshal(l)
					logFunc(string(data))
				}()
				next(ctx)
			}
		}
	}
	var testMiddleWare = func(str string) Middleware {
		return func(next HandleFunc) HandleFunc {
			return func(ctx *Context) {
				fmt.Println(str + "头")
				next(ctx)
				fmt.Println(str + "尾")
			}
		}
	}
	server := NewHTTPServer(ServerAddMiddleware(logMiddleware()))
	server.Use(http.MethodGet, "/", testMiddleWare("/"))
	server.Use(http.MethodGet, "/a", testMiddleWare("a"))
	server.Use(http.MethodGet, "/b", testMiddleWare("b"))
	server.Use(http.MethodGet, "/a/b", testMiddleWare("ab"))
	server.Use(http.MethodGet, "/a/c", testMiddleWare("ac"))
	server.Use(http.MethodGet, "/b/d", testMiddleWare("bd"))

	server.Get("/a/b", func(ctx *Context) {
		fmt.Println("我是核心内容")
		ctx.RespStatusCode = 200
		ctx.RespData = []byte(`OK`)
	})

	server.Start(":8081")
}

func Test_TemplateEngine(t *testing.T) {
	// 加载测试文件
	tpl, err := template.ParseGlob("testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	tplEngine := &GoTemplateEngine{tpl}
	server := NewHTTPServer(ServerAddTemplate(tplEngine))
	server.Get("/login", func(ctx *Context) {
		er := ctx.Render("login.gohtml", nil)
		if er != nil {
			t.Fatal(er)
		}
	})
	err = server.Start(":8081")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_FileUpload(t *testing.T) {
	panicMiddle := func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			next(ctx)
		}
	}
	s := NewHTTPServer(ServerAddMiddleware(panicMiddle))
	s.Get("/upload_page", func(ctx *Context) {
		tpl := template.New("upload")
		tpl, err := tpl.Parse(`
<html>
<body>
	<form action="/upload" method="post" enctype="multipart/form-data">
		 <input type="file" name="myfile" />
		 <button type="submit">上传</button>
	</form>
</body>
<html>
`)
		if err != nil {
			t.Fatal(err)
		}

		page := &bytes.Buffer{}
		err = tpl.Execute(page, nil)
		if err != nil {
			t.Fatal(err)
		}

		ctx.RespStatusCode = 200
		ctx.RespData = page.Bytes()
	})

	s.Post("/upload", NewFileUploader(
		SetFileField("myfile"),
		SetDstPathFunc(func(header *multipart.FileHeader) string {
			return path.Join("testdata", "upload", header.Filename)
		}),
	).Handle())
	s.Start(":8081")
}

func Test_FileDownload(t *testing.T) {
	s := NewHTTPServer()
	s.Get("/download", (&FileDownloader{
		// 下载的文件所在目录
		Dir: "./testdata/download",
	}).Handle())
	// 在浏览器里面输入 localhost:8081/download?file=test.txt
	s.Start(":8081")
}
