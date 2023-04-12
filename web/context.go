// create by chencanhua in 2023/04/12

package web

import "net/http"

type Context struct {
	Req        *http.Request
	Resp       http.ResponseWriter
	PathParams map[string]string
}
