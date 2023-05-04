// create by chencanhua in 2023/5/1
package cookie

import (
	"net/http"
	"web_framework/session"
)

var _ session.Propagator = &Propagator{}

type Propagator struct {
	CookieName    string
	CookieOptions func(c *http.Cookie)
}

type PropagatorOptions func(*Propagator)

func WithCookieName(cname string) PropagatorOptions {
	return func(propagator *Propagator) {
		propagator.CookieName = cname
	}
}

func NewPropagator(opts ...PropagatorOptions) *Propagator {
	p := &Propagator{
		CookieName:    "sessid",
		CookieOptions: func(c *http.Cookie) {},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	c := &http.Cookie{
		Name:  p.CookieName,
		Value: id,
	}
	p.CookieOptions(c)
	http.SetCookie(writer, c)
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	cookie, err := req.Cookie(p.CookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	c := &http.Cookie{
		Name:   p.CookieName,
		MaxAge: -1,
	}
	http.SetCookie(writer, c)
	return nil
}
