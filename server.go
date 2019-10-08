/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cobweb

import (
	"context"
	"net/http"
)

type Mux struct {
	recordMap map[string][]Record
	NotFound  http.HandlerFunc
	groups    []*group
}

type group struct {
	pattern  string
	handlers []http.HandlerFunc
}

func New() *Mux {
	return &Mux{
		recordMap: make(map[string][]Record, 7),
	}
}

func (m *Mux) Group(pattern string, fn func(), handlers ...http.HandlerFunc) {
	gr := &group{pattern: CleanPath(pattern)}
	if len(handlers) > 0 {
		gr.handlers = handlers
	}
	m.groups = append(m.groups, gr)
	fn()
	m.groups = m.groups[:len(m.groups)-1]
}

func (m *Mux) Head(pattern string, handlers ...http.HandlerFunc) {
	m.handle(http.MethodHead, pattern, handlers)
}

func (m *Mux) Get(pattern string, handlers ...http.HandlerFunc) {
	m.handle(http.MethodGet, pattern, handlers)
	m.handle(http.MethodHead, pattern, handlers)
}

func (m *Mux) Post(pattern string, handlers ...http.HandlerFunc) {
	m.handle(http.MethodPost, pattern, handlers)
}

func (m *Mux) Put(pattern string, handlers ...http.HandlerFunc) {
	m.handle(http.MethodPut, pattern, handlers)
}

func (m *Mux) Delete(pattern string, handlers ...http.HandlerFunc) {
	m.handle(http.MethodDelete, pattern, handlers)
}

func (m *Mux) Patch(pattern string, handlers ...http.HandlerFunc) {
	m.handle(http.MethodPatch, pattern, handlers)
}

func (m *Mux) Options(pattern string, handlers ...http.HandlerFunc) {
	m.handle(http.MethodOptions, pattern, handlers)
}

func (m *Mux) Any(pattern string, handlers ...http.HandlerFunc) {
	m.handle(http.MethodGet, pattern, handlers)
	m.handle(http.MethodPost, pattern, handlers)
	m.handle(http.MethodHead, pattern, handlers)
	m.handle(http.MethodDelete, pattern, handlers)
	m.handle(http.MethodPatch, pattern, handlers)
	m.handle(http.MethodPut, pattern, handlers)
	m.handle(http.MethodOptions, pattern, handlers)
}

func (m *Mux) handle(method, pattern string, handlers []http.HandlerFunc) {
	if len(m.groups) > 0 {
		groupPattern := ""
		for _, g := range m.groups {
			groupPattern += g.pattern
			if len(g.handlers) > 0 {
				handlers = append(g.handlers, handlers...)
			}
		}
		pattern = groupPattern + pattern
	}

	m.recordMap[method] = append(m.recordMap[method], NewRecord(CleanPath(pattern), newHandler(handlers)))
}

func (m *Mux) Build() (http.Handler, error) {
	mux := newServeMux()
	for m, records := range m.recordMap {
		router := NewRouter()
		if err := router.Build(records); err != nil {
			return nil, err
		}
		mux.routers[m] = router
	}
	mux.notFound = m.NotFound
	return mux, nil
}

type serveMux struct {
	routers  map[string]*Router
	notFound http.HandlerFunc
}

func newServeMux() *serveMux {
	return &serveMux{
		routers: make(map[string]*Router, 7),
	}
}

func (mux *serveMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, params := mux.handler(r.Method, r.URL.Path)
	ctx := context.WithValue(r.Context(), context.Background(), params)
	handler.ServeHTTP(w, r.WithContext(ctx))
}

func (mux *serveMux) handler(method, path string) (http.Handler, Params) {
	if router, found := mux.routers[method]; found {
		if handler, params, found := router.Lookup(path); found {
			return handler.(http.Handler), params
		}
	}
	if mux.notFound != nil {
		return mux.notFound, nil
	}
	return NotFound, nil
}

var NotFound = http.HandlerFunc(http.NotFound)

func newHandler(list []http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nw, ok := w.(ResponseWriter)
		if !ok {
			nw = NewResponseWriter(w)
		}
		if nw.Written() {
			return
		}
		for i := range list {
			list[i].ServeHTTP(nw, r)
			if nw.Written() {
				return
			}
		}
	})
}
