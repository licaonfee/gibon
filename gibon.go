package gibon

import (
	"net/http"
)

//Middleware is a function tom improve or extends
//http.Hanlder functionality
type Middleware func(http.Handler) http.Handler

//Chain is a middleware pipeline
//is safe to use this object with zero value
type Chain struct {
	mid []Middleware
}

//New create a new empty middleware chain
//is equivalent to var gibon.Chain
func New() Chain {
	return Chain{}
}

//With create a new Chain with all pre-existent middlewares and
//"m" as last middleware in chain
func (c Chain) With(m Middleware) Chain {
	nmid := make([]Middleware, len(c.mid))
	copy(nmid, c.mid)
	nmid = append(nmid, m)
	return Chain{mid: nmid}
}

//Build create an http.Handler with all middlewared from Chain
//wrapping h (http.Hanlder)
func (c Chain) Build(h http.Handler) http.Handler {
	var newh http.Handler = h
	for i := len(c.mid) - 1; i >= 0; i-- {
		newh = c.mid[i](newh)
	}
	return newh
}

//BuildFunc same as Build but for function objects
func (c Chain) BuildFunc(f func(http.ResponseWriter, *http.Request)) http.Handler {
	return c.Build(http.HandlerFunc(f))
}
