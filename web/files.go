// Package web create by chencanhua in 2023/4/25
package web

import (
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// FileUploader 文件上传
type FileUploader struct {
	FileField   string
	DstPathFunc func(header *multipart.FileHeader) string
}

/**
  options模式实例
*/

type FileUploadOptions func(f *FileUploader)

func SetFileField(fileField string) FileUploadOptions {
	return func(f *FileUploader) {
		f.FileField = fileField
	}
}

func SetDstPathFunc(fun func(header *multipart.FileHeader) string) FileUploadOptions {
	return func(f *FileUploader) {
		f.DstPathFunc = fun
	}
}

func NewFileUploader(opts ...FileUploadOptions) *FileUploader {
	uploader := &FileUploader{}
	for _, opt := range opts {
		opt(uploader)
	}
	return uploader
}

func (f *FileUploader) Handle() HandleFunc {
	return func(ctx *Context) {
		src, header, err := ctx.Req.FormFile(f.FileField)
		defer src.Close()
		if err != nil {
			ctx.RespData = []byte(`上传失败`)
			ctx.RespStatusCode = http.StatusInternalServerError
			log.Fatalln(err)
			return
		}
		dst, err := os.OpenFile(f.DstPathFunc(header),
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
		defer dst.Close()
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			log.Fatalln(err)
			return
		}

		_, err = io.CopyBuffer(dst, src, nil)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("上传失败")
			log.Fatalln(err)
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("上次成功")
	}
}

type FileDownloader struct {
	Dir string
}

func (f *FileDownloader) Handle() HandleFunc {
	return func(ctx *Context) {
		req, _ := ctx.QueryValue("file").String()
		path := filepath.Join(f.Dir, filepath.Clean(req))
		fn := filepath.Base(path)
		header := ctx.Resp.Header()
		header.Set("Content-Disposition", "attachment;filename="+fn)
		header.Set("Content-Description", "File Transfer")
		header.Set("Content-Type", "application/octet-stream")
		header.Set("Content-Transfer-Encoding", "binary")
		header.Set("Expires", "0")
		header.Set("Cache-Control", "must-revalidate")
		header.Set("Pragma", "public")
		http.ServeFile(ctx.Resp, ctx.Req, path)
	}
}
