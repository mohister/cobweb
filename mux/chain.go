package mux

import "net/http"

type Chain struct {
	list []http.Handler
}

func NewChain() *Chain {
	return &Chain{
		make([]http.Handler, 0, 10),
	}
}

func (c *Chain) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	nrw, ok := rw.(ResponseWriter)
	if !ok {
		nrw = NewResponseWriter(rw)
	}
	if nrw.Written() {
		return
	}
	for i := range c.list {
		c.list[i].ServeHTTP(nrw, r)
		if nrw.Written() {
			return
		}
	}
}

func (c *Chain) Add(handlers ...http.Handler) {
	l := len(c.list)
	c.expansion(len(handlers))
	copy(c.list[l:], handlers)
}

func (c *Chain) AddFunc(handlers ...http.HandlerFunc) {
	l := len(c.list)
	c.expansion(len(handlers))
	for i := range handlers {
		c.list[l+i] = handlers[i]
	}
}

func (c *Chain) expansion(n int) {
	l := len(c.list)
	if cl := cap(c.list); l+n > cl {
		list := make([]http.Handler, 0, cl+n<<1)
		copy(list, c.list)
		c.list = list
	}
	c.list = c.list[:l+n]
}
