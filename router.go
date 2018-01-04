package cobweb

import (
	"net/http"
	"github.com/chinx/utils/strutil"
)

type group struct {
	pattern  string
	handlers []Handle
}

type Router struct {
	tries   map[string]*trie
	notFond Handle
	groups  []*group
}

func NewRouter() *Router {
	return &Router{
		tries: make(map[string]*trie, 7),
	}
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if root, ok := r.tries[req.Method]; ok {
		handler, params := root.match(req.URL.Path)
		if handler != nil {
			handler(NewResponseWriter(rw), req, params)
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

	r.handle(http.MethodGet, path, func(w http.ResponseWriter, req *http.Request, ps *Params) {
		if filePath, _ := ps.Get("filepath"); filePath != "" {
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
		r.notFond(rw, req, nil)
		return
	}
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte(req.URL.Path + " not fond"))
	rw = nil
}

func (r *Router) NotFound(handler Handle) {
	r.notFond = handler
}

func (r *Router) Group(pattern string, fn func(), handlers ...Handle) {
	r.groups = append(r.groups, &group{"/" + pattern, handlers})
	fn()
	r.groups = r.groups[:len(r.groups)-1]
}

func (r *Router) Head(pattern string, handlers ...Handle) {
	r.handle(http.MethodHead, pattern, handlers...)
}

func (r *Router) Get(pattern string, handlers ...Handle) {
	r.handle(http.MethodGet, pattern, handlers...)
	r.handle(http.MethodHead, pattern, handlers...)
}

func (r *Router) Post(pattern string, handlers ...Handle) {
	r.handle(http.MethodPost, pattern, handlers...)
}

func (r *Router) Put(pattern string, handlers ...Handle) {
	r.handle(http.MethodPut, pattern, handlers...)
}

func (r *Router) Delete(pattern string, handlers ...Handle) {
	r.handle(http.MethodDelete, pattern, handlers...)
}

func (r *Router) Patch(pattern string, handlers ...Handle) {
	r.handle(http.MethodPatch, pattern, handlers...)
}

func (r *Router) Options(pattern string, handlers ...Handle) {
	r.handle(http.MethodOptions, pattern, handlers...)
}

func (r *Router) Any(pattern string, handlers ...Handle) {
	r.handle(http.MethodGet, pattern, handlers...)
	r.handle(http.MethodPost, pattern, handlers...)
	r.handle(http.MethodHead, pattern, handlers...)
	r.handle(http.MethodDelete, pattern, handlers...)
	r.handle(http.MethodPatch, pattern, handlers...)
	r.handle(http.MethodPut, pattern, handlers...)
	r.handle(http.MethodOptions, pattern, handlers...)
}

func (r *Router) handle(method, pattern string, handlers ...Handle) {
	pattern = "/" + strutil.Trim(pattern, '/')
	if len(r.groups) > 0 {
		groupPattern := ""
		h := make([]Handle, 0)
		for _, g := range r.groups {
			groupPattern += g.pattern
			h = append(h, g.handlers...)
		}
		pattern = groupPattern + pattern
		h = append(h, handlers...)
		handlers = h
	}

	root, ok := r.tries[method]
	if !ok {
		root = NewTrie()
		r.tries[method] = root
	}
	root.addNode(pattern, handlersChain(handlers))
}

func handlersChain(handlers []Handle) Handle {
	nHandlers := make([]Handle, 0, 10)
	l := 0
	for i := 0; i < len(handlers); i++ {
		if handlers[i] != nil {
			nHandlers = append(nHandlers, handlers[i])
			l++
		}
	}
	if l == 0 {
		return nil
	}

	if l == 1 {
		return nHandlers[0]
	}

	return func(rw http.ResponseWriter, req *http.Request, params *Params) {
		length := len(handlers)
		for i := 0; i < length; i++ {
			if _, ok := rw.(ResponseWriter); ok && rw.(ResponseWriter).Written() {
				return
			}
			handlers[i](rw, req, params)
		}
		if _, ok := rw.(ResponseWriter); ok && !rw.(ResponseWriter).Written() {
			rw.Write([]byte("Mohist is OK"))
		}
	}
}
