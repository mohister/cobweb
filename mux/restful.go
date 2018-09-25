package mux

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

var (
	rootMethodNotOnly = "Handler of root path mush be only"
	methodNil         = "Handler of pattern '%s' mush not nil"
	pathNotMatched    = "Path '%s' in pattern '%s' not matched regex `^[a-z](_*[a-z0-9]+)*$`"
	endpointErr       = "'*' in pattern '%s' must be endpoint"
	conflictedErr     = "Pattern '%s' conflicted with '%s'"
)

type RestMux struct {
	*restEntry
	buffer *bytes.Buffer
	hash   map[string]string
}

func NewRestMux() *RestMux {
	return &RestMux{
		restEntry: &restEntry{
			pattern:  "root",
			priority: static,
			children: make(entryList, 0, 32),
		},
		buffer: bytes.NewBuffer(make([]byte, 0, 1024)),
		hash:   make(map[string]string, 32),
	}
}

func (mux *RestMux) match(path string) (h http.Handler, params *params) {
	count := countParts(path)
	if mux.maxParts < count {
		return
	}

	part, target := "", mux.restEntry
	start, length := 0, len(path)
	ended := false
walk:
	count--
	part, start = between(path, start)
	if part == "" {
		h = mux.h
		return
	}
	ended = start == length
	for i := target.index; i < len(target.children); i++ {
		nn := target.children[i]
		if nn.maxParts < count {
			continue
		}
		if nn.priority == elastic {
			if nn.h != nil {
				h = nn.h
				if !ended {
					part += "/" + path[start:]
				}
				if params == nil {
					params = getParams()
				}
				params.addValue(part)
			}
			target = nn
			break
		}
		if isDynamic := nn.priority == dynamic; isDynamic || nn.pattern == part {
			if isDynamic {
				if params == nil {
					params = getParams()
				}
				params.values.add(part)
			}
			if !ended {
				target.index = i
				target = nn
				goto walk
			} else if nn.h != nil {
				h = nn.h
				target = nn
				break
			} else if nn.wide != nil {
				h = nn.wide
				break
			}
		}
	}
	if len(target.paramsKeys) > 0 {
		params.keys = target.paramsKeys
	}
	return
}

func (mux *RestMux) handle(pattern string, handler http.Handler) {
	if handler == nil {
		panic(fmt.Sprintf(methodNil, pattern))
	}

	nParts, nParams := countPartsAndParams(pattern)

	var keys entries
	if nParams != 0 {
		keys = make([]string, 0, nParams)
		setMaxParams(nParams)
	}

	path, target := "", mux.restEntry
	start, length := 0, len(pattern)
	for start < length {
		path, start = between(pattern, start)
		if path == "" {
			if _, ok := mux.hash["/"]; ok {
				panic(rootMethodNotOnly)
			}
			target.h = handler
			mux.hash["/"] = "/"
			break
		}
		mux.buffer.WriteByte('/')

		priority, ended := static, start == length
		if path[0] == '*' || path[0] == ':' {
			priority = dynamic
			mux.buffer.WriteString("%v")
			if path[0] == '*' {
				if !ended {
					panic(fmt.Sprintf(endpointErr, pattern))
				}
				priority = elastic
			}
			path = path[1:]
			keys = keys.add(path)
		} else {
			mux.buffer.WriteString(path)
		}

		if !isPath(path) {
			panic(fmt.Sprintf(pathNotMatched, path, pattern))
		}
		target.SetMaxParts(nParts)
		nParts--

		var nn *restEntry
		for i := range target.children {
			if target.children[i].priority != priority {
				continue
			}
			if target.children[i].pattern == path {
				nn = target.children[i]
				break
			}
		}

		if nn == nil {
			nn = &restEntry{
				pattern:    path,
				priority:   priority,
				children:   make(entryList, 0, 10),
				paramsKeys: keys,
			}
			if priority == elastic && target.wide == nil {
				target.wide = handler
			}
			target.children = target.children.add(nn)
		}

		if ended {
			flag := mux.buffer.String()
			mux.buffer.Reset()
			if val, ok := mux.hash[flag]; ok {
				panic(fmt.Sprintf(conflictedErr, val, pattern))
			}
			nn.h = handler
			mux.hash[flag] = pattern
			break
		}
		target = nn
	}
}

func (mux *RestMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h, params := mux.match(r.URL.Path); h != nil {
		ctx := context.WithValue(r.Context(), context.Background(), params)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (mux *RestMux) Handle(pattern string, handler http.Handler) {
	mux.handle(pattern, handler)
}

func (mux *RestMux) HandleFunc(pattern string, handler http.HandlerFunc) {
	mux.handle(pattern, handler)
}
