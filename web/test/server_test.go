package test

import (
	"bytes"
	"fmt"
	"html/template"
	"mime/multipart"
	"net/http"
	"path"
	"testing"
	"web_framework/web"
	"web_framework/web/middleware/accesslog"
	"web_framework/web/middleware/errmdl"
)

func TestServer(t *testing.T) {
	s := web.NewHTTPServer()
	s.Get("/", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, world"))
	})
	s.Get("/user", func(ctx *web.Context) {
		ctx.Resp.Write([]byte("hello, user"))
	})

	s.Start(":8081")
}

func TestMiddlerware(t *testing.T) {
	var logFunc = func(log string) {
		fmt.Println(log)
	}
	logBuilder := &accesslog.Builder{}
	logBuilder.SetLogFunc(logFunc)
	var logMiddleware = logBuilder.Build()
	var testMiddleWare = func(str string) web.Middleware {
		return func(next web.HandleFunc) web.HandleFunc {
			return func(ctx *web.Context) {
				fmt.Println(str + "头")
				next(ctx)
				fmt.Println(str + "尾")
			}
		}
	}
	server := web.NewHTTPServer(web.ServerAddMiddleware(logMiddleware))
	server.Use(http.MethodGet, "/", testMiddleWare("/"))
	server.Use(http.MethodGet, "/a", testMiddleWare("a"))
	server.Use(http.MethodGet, "/b", testMiddleWare("b"))
	server.Use(http.MethodGet, "/a/b", testMiddleWare("ab"))
	server.Use(http.MethodGet, "/a/c", testMiddleWare("ac"))
	server.Use(http.MethodGet, "/b/d", testMiddleWare("bd"))

	server.Get("/a/b", func(ctx *web.Context) {
		fmt.Println("我是核心内容")
		ctx.RespStatusCode = 200
		ctx.RespData = []byte(`OK`)
	})

	server.Start(":8081")
}

func Test_TemplateEngine(t *testing.T) {
	// 加载测试文件
	tpl, err := template.ParseGlob("../testdata/tpls/*.gohtml")
	if err != nil {
		t.Fatal(err)
	}
	tplEngine := &web.GoTemplateEngine{tpl}
	server := web.NewHTTPServer(web.ServerAddTemplate(tplEngine))
	server.Get("/login", func(ctx *web.Context) {
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
	builder := errmdl.Builder{}
	panicMiddle := builder.Build()
	s := web.NewHTTPServer(web.ServerAddMiddleware(panicMiddle))
	s.Get("/upload_page", func(ctx *web.Context) {
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

	s.Post("/upload", web.NewFileUploader(
		web.SetFileField("myfile"),
		web.SetDstPathFunc(func(header *multipart.FileHeader) string {
			return path.Join("testdata", "upload", header.Filename)
		}),
	).Handle())
	s.Start(":8081")
}

func Test_FileDownload(t *testing.T) {
	s := web.NewHTTPServer()
	s.Get("/download", (&web.FileDownloader{
		// 下载的文件所在目录
		Dir: "../testdata/download",
	}).Handle())
	// 在浏览器里面输入 localhost:8081/download?file=test.txt
	s.Start(":8081")
}
