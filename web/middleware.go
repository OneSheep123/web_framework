// Package web create by chencanhua in 2023/4/20
package web

type Middleware func(next HandleFunc) HandleFunc
