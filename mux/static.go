package mux

import (
	"net/http"
)

type StaticMux map[string]stcEntry

type stcEntry struct {
	h        http.Handler
	pattern  string
	isPrefix bool
}

func prefixMatch(pattern, path string) bool {
	n, l := len(pattern), len(path)
	if n > l {
		return false
	}
	if l == n {
		return pattern == path
	}
	return path[0:n] == pattern && path[n] == '/'
}

func NewStaticMux() StaticMux { return make(map[string]stcEntry, 32) }

func (mux StaticMux) match(path string) (h http.Handler) {
	v, ok := mux[path]
	if ok {
		return v.h
	}

	for k, v := range mux {
		if v.isPrefix && prefixMatch(k, path) {
			return v.h
		}

	}
	return
}

func (mux StaticMux) handle(pattern string, handler http.Handler, isPrefix bool) {
	if pattern == "" {
		panic("http: invalid pattern")
	}
	if handler == nil {
		panic("http: nil handler")
	}
	if mux.match(pattern) != nil {
		panic("http: multiple registrations for " + pattern)
	}

	if pattern[0] != '/' {
		pattern = "/" + pattern
	}
	mux[pattern] = stcEntry{h: handler, pattern: pattern, isPrefix: isPrefix}
}

func (mux StaticMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h := mux.match(r.URL.Path); h != nil {
		h.ServeHTTP(w, r)
	}
}

func (mux StaticMux) Handle(pattern string, handler http.Handler) {
	mux.handle(pattern, handler, false)
}

func (mux StaticMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	mux.handle(pattern, handler, false)
}

func (mux StaticMux) ServeFiles(pattern string, root http.FileSystem) {
	fileServer := http.FileServer(root)

	mux.handle(pattern, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = r.URL.Path[len(pattern):]
			fileServer.ServeHTTP(w, r)
		}), true)
}
