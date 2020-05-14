package gibon

import (
	"net/http"
)

type Middleware func(http.Handler) http.Handler

type Chain struct {
	mid []Middleware
}

func Add(m Middleware) Chain {
	c := Chain{}
	return c.Add(m)
}

func (c Chain) Add(m Middleware) Chain {
	nmid := make([]Middleware, len(c.mid))
	copy(nmid, c.mid)
	nmid = append(nmid, m)
	return Chain{mid: nmid}
}

func (c Chain) Build(h http.Handler) http.Handler {
	var newh http.Handler = h
	for i := len(c.mid) - 1; i >= 0; i-- {
		newh = c.mid[i](newh)
	}
	return newh
}

func (c Chain) BuildFunc(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return c.Build(http.HandlerFunc(f))
}
