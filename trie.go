package cobweb

import (
	"bytes"
	"fmt"
	"net/http"
)

var (
	rootMethodNotOnly = "Method of root path mush be only"
	methodNil         = "Method of pattern '%s' mush not nil"
	pathNotMatched    = "Path '%s' in pattern '%s' not matched regex `^[a-z](_*[a-z0-9]+)*$`"
	endpointErr       = "'*' in pattern '%s' must be endpoint"
	conflictedErr     = "Pattern '%s' conflicted with '%s'"
)

type trie struct {
	*node
	buffer *bytes.Buffer
	hash   map[string]string
}

func NewTrie() *trie {
	return &trie{
		node: &node{
			pattern:  "root",
			priority: static,
			children: make([]*node, 0, 32),
		},
		buffer: bytes.NewBuffer(make([]byte, 0, 1024)),
		hash:   make(map[string]string, 32),
	}
}

func (t *trie) addNode(pattern string, handler http.Handler) {
	if handler == nil {
		panic(fmt.Sprintf(methodNil, pattern))
	}

	nParts, nParams := countPartsAndParams(pattern)

	var keys entries
	if nParams != 0 {
		keys = make([]string, 0, nParams)
		setMaxParams(nParams)
	}

	path, target := "", t.node
	start, length := 0, len(pattern)
	for start < length {
		path, start = between(pattern, start)
		if path == "" {
			if _, ok := t.hash["/"]; ok {
				panic(rootMethodNotOnly)
			}
			target.method = handler
			t.hash["/"] = "/"
			break
		}
		t.buffer.WriteByte('/')

		priority, ended := static, start == length
		if path[0] == '*' || path[0] == ':' {
			priority = dynamic
			t.buffer.WriteString("%v")
			if path[0] == '*' {
				if !ended {
					panic(fmt.Sprintf(endpointErr, pattern))
				}
				priority = elastic
			}
			path = path[1:]
			keys.add(path)
		} else {
			t.buffer.WriteString(path)
		}

		if !isPath(path) {
			panic(fmt.Sprintf(pathNotMatched, path, pattern))
		}
		target.SetMaxParts(nParts)
		nParts--

		var nn *node
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
			nn = &node{
				pattern:    path,
				priority:   priority,
				children:   make([]*node, 0, 10),
				paramsKeys: keys,
			}
			if priority == elastic && target.wideMethod == nil {
				target.wideMethod = handler
			}
			target.children = target.children.add(nn)
		}

		if ended {
			flag := t.buffer.String()
			t.buffer.Reset()
			if val, ok := t.hash[flag]; ok {
				panic(fmt.Sprintf(conflictedErr, val, pattern))
			}
			nn.method = handler
			t.hash[flag] = pattern
			break
		}
		target = nn
	}
}

func (t *trie) match(pattern string) (handler http.Handler, params *Params) {
	count := countParts(pattern)
	if t.maxParts < count {
		return
	}

	path, target := "", t.node
	start, length := 0, len(pattern)
	ended := false
walk:
	count--
	path, start = between(pattern, start)
	if path == "" {
		handler = t.method
		return
	}
	ended = start == length
	for i := target.index; i < len(target.children); i++ {
		nn := target.children[i]
		if nn.maxParts < count {
			continue
		}
		if nn.priority == elastic {
			if nn.method != nil {
				handler = nn.method
				if !ended {
					path += "/" + pattern[start:]
				}
				if params == nil{
					params = getParams()
				}
				params.values.add(path)
			}
			target = nn
			break
		}
		if isDynamic := nn.priority == dynamic; isDynamic || nn.pattern == path {
			if isDynamic {
				if params == nil{
					params = getParams()
				}
				params.values.add(path)
			}
			if !ended {
				target.index = i
				target = nn
				goto walk
			} else if nn.method != nil {
				handler = nn.method
				target = nn
				break
			} else if nn.wideMethod != nil {
				handler = nn.wideMethod
				break
			}
		}
	}
	if len(target.paramsKeys) > 0{
		params.values = target.paramsKeys
	}
	return
}
