package cobweb

import (
	"context"
	"net/http"
)

type group struct {
	pattern string
	chain   *Chain
}

type Router struct {
	tries   map[string]*trie
	notFond http.Handler
	groups  []*group
}

func New() *Router {
	return &Router{
		tries:  make(map[string]*trie, 7),
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if root, ok := r.tries[req.Method]; ok {
		handler, params := root.match(req.URL.Path)
		ctx := context.WithValue(req.Context(), context.Background(), params)
		if handler != nil {
			handler.ServeHTTP(NewResponseWriter(rw), req.WithContext(ctx))
			if params !=nil {
				putParams(params)
			}
			return
		}
	}
	r.notFoundHandler(rw, req)
}

func (r *Router) ServeFiles(path string, root http.FileSystem) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}

	fileServer := http.FileServer(root)

	r.Get(path, func(w http.ResponseWriter, req *http.Request) {
		ps := req.Context().Value(context.Background()).(Params)
		if filePath := ps.String("filepath"); filePath != "" {
			req.URL.Path = filePath
			fileServer.ServeHTTP(w, req)
			return
		}
		r.notFoundHandler(w, req)
	})
}

func (r *Router) notFoundHandler(rw http.ResponseWriter, req *http.Request) {
	// Handle 404
	if r.notFond != nil {
		r.notFond.ServeHTTP(rw, req)
		return
	}
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(req.URL.Path + " not fond"))
	rw = nil
}

func (r *Router) NotFound(handler http.HandlerFunc) {
	r.notFond = handler
}

func (r *Router) Group(pattern string, fn func(), handlers ...http.HandlerFunc) {
	gr := &group{pattern: "/" + pattern}
	if len(handlers) > 0 {
		gr.chain = NewChain()
		gr.chain.AddFunc(handlers...)
	}
	r.groups = append(r.groups, gr)
	fn()
	r.groups = r.groups[:len(r.groups)-1]
}

func (r *Router) Head(pattern string, handlers ...http.HandlerFunc) {
	r.handle(http.MethodHead, pattern, handlers)
}

func (r *Router) Get(pattern string, handlers ...http.HandlerFunc) {
	r.handle(http.MethodGet, pattern, handlers)
	r.handle(http.MethodHead, pattern, handlers)
}

func (r *Router) Post(pattern string, handlers ...http.HandlerFunc) {
	r.handle(http.MethodPost, pattern, handlers)
}

func (r *Router) Put(pattern string, handlers ...http.HandlerFunc) {
	r.handle(http.MethodPut, pattern, handlers)
}

func (r *Router) Delete(pattern string, handlers ...http.HandlerFunc) {
	r.handle(http.MethodDelete, pattern, handlers)
}

func (r *Router) Patch(pattern string, handlers ...http.HandlerFunc) {
	r.handle(http.MethodPatch, pattern, handlers)
}

func (r *Router) Options(pattern string, handlers ...http.HandlerFunc) {
	r.handle(http.MethodOptions, pattern, handlers)
}

func (r *Router) Any(pattern string, handlers ...http.HandlerFunc) {
	r.handle(http.MethodGet, pattern, handlers)
	r.handle(http.MethodPost, pattern, handlers)
	r.handle(http.MethodHead, pattern, handlers)
	r.handle(http.MethodDelete, pattern, handlers)
	r.handle(http.MethodPatch, pattern, handlers)
	r.handle(http.MethodPut, pattern, handlers)
	r.handle(http.MethodOptions, pattern, handlers)
}

func (r *Router) handle(method, pattern string, handlers []http.HandlerFunc) {
	pattern = "/" + trim(pattern)
	c := NewChain()
	if len(r.groups) > 0 {
		groupPattern := ""
		for _, g := range r.groups {
			groupPattern += g.pattern
			c.Add(g.chain)
		}
		pattern = groupPattern + pattern
	}
	c.AddFunc(handlers...)

	root, ok := r.tries[method]
	if !ok {
		root = NewTrie()
		r.tries[method] = root
	}
	root.addNode(pattern, c)
}
