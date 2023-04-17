// create by chencanhua in 2023/4/17
package logging

import (
	"fmt"
	"net/http"
	"testing"
	"web_framework/web"
)

func TestMiddleware(t *testing.T) {
	build := NewBuilder(func(log string) {
		fmt.Println(log)
	})
	server := web.NewHTTPServer(web.ServerWithMiddleware(build.Build()))
	server.Post("/a/b/*", func(ctx *web.Context) {
		fmt.Println("hello, it's me")
	})
	req, err := http.NewRequest(http.MethodPost, "/a/b/c", nil)
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
