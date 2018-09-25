package mux

import (
	"net/http"
)

type trie struct {
	stc   StaticMux
	rest  *RestMux
	chain *Chain
}

type group struct {
	pattern string
	chain   *Chain
}

type ServeMux struct {
	tries   map[string]*trie
	notFond http.Handler
	groups  []*group
}

func New() *ServeMux {
	return &ServeMux{
		tries: make(map[string]*trie, 7),
		notFond: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, r.URL.Path+" not fond", http.StatusNotFound)
		}),
	}
}

func isRestURL(path string) bool {
	for i := 0; i < len(path)-1; i++ {
		if path[i] == '/' && path[i+1] == '*' || path[i] == ':' {
			return true
		}
	}
	return false
}

func (mux *ServeMux) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if root, ok := mux.tries[req.Method]; ok {
		root.chain.ServeHTTP(NewResponseWriter(rw), req)
	}
}

func (mux *ServeMux) ServeFiles(path string, root http.FileSystem) {
	t := mux.getTries(http.MethodGet)
	t.stc.ServeFiles(path, root)
}

func (mux *ServeMux) NotFound(handler http.HandlerFunc) {
	if handler != nil {
		mux.notFond = handler
	}
}

func (mux *ServeMux) Group(pattern string, fn func(), handlers ...http.HandlerFunc) {
	gr := &group{pattern: CleanPath(pattern)}
	if len(handlers) > 0 {
		gr.chain = NewChain()
		gr.chain.AddFunc(handlers...)
	}
	mux.groups = append(mux.groups, gr)
	fn()
	mux.groups = mux.groups[:len(mux.groups)-1]
}

func (mux *ServeMux) Head(pattern string, handlers ...http.HandlerFunc) {
	mux.handle(http.MethodHead, pattern, handlers)
}

func (mux *ServeMux) Get(pattern string, handlers ...http.HandlerFunc) {
	mux.handle(http.MethodGet, pattern, handlers)
	mux.handle(http.MethodHead, pattern, handlers)
}

func (mux *ServeMux) Post(pattern string, handlers ...http.HandlerFunc) {
	mux.handle(http.MethodPost, pattern, handlers)
}

func (mux *ServeMux) Put(pattern string, handlers ...http.HandlerFunc) {
	mux.handle(http.MethodPut, pattern, handlers)
}

func (mux *ServeMux) Delete(pattern string, handlers ...http.HandlerFunc) {
	mux.handle(http.MethodDelete, pattern, handlers)
}

func (mux *ServeMux) Patch(pattern string, handlers ...http.HandlerFunc) {
	mux.handle(http.MethodPatch, pattern, handlers)
}

func (mux *ServeMux) Options(pattern string, handlers ...http.HandlerFunc) {
	mux.handle(http.MethodOptions, pattern, handlers)
}

func (mux *ServeMux) Any(pattern string, handlers ...http.HandlerFunc) {
	mux.handle(http.MethodGet, pattern, handlers)
	mux.handle(http.MethodPost, pattern, handlers)
	mux.handle(http.MethodHead, pattern, handlers)
	mux.handle(http.MethodDelete, pattern, handlers)
	mux.handle(http.MethodPatch, pattern, handlers)
	mux.handle(http.MethodPut, pattern, handlers)
	mux.handle(http.MethodOptions, pattern, handlers)
}

func (mux *ServeMux) handle(method, pattern string, handlers []http.HandlerFunc) {
	pattern = CleanPath(pattern)
	c := NewChain()
	if len(mux.groups) > 0 {
		groupPattern := ""
		for _, g := range mux.groups {
			groupPattern += g.pattern
			c.Add(g.chain)
		}
		pattern = groupPattern + pattern
	}
	c.AddFunc(handlers...)

	root := mux.getTries(method)
	if isRestURL(pattern) {
		root.rest.Handle(pattern, c)
		return
	}
	root.stc.Handle(pattern, c)
}

func (mux *ServeMux) getTries(method string) *trie {
	root, ok := mux.tries[method]
	if !ok {
		root = &trie{
			stc:   NewStaticMux(),
			rest:  NewRestMux(),
			chain: NewChain(),
		}
		root.chain.Add(root.stc, root.rest, mux.notFond)
		mux.tries[method] = root
	}
	return root
}
