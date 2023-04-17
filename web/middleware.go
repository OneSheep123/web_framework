// Package web create by chencanhua in 2023/4/16
package web

type Middleware func(next HandleFunc) HandleFunc
