package cobweb

import (
	"bytes"
	"log"
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
	hash map[string]string
}

func NewTrie() *trie {
	return &trie{
		node: &node{
			pattern:  "root",
			priority: static,
			children: make([]*node, 0, 10),
			parent:   nil,
		},
		hash: make(map[string]string),
	}
}
func (n *trie) addNode(pattern string, handler Handle) {
	if handler == nil {
		log.Panicf(methodNil, pattern)
	}
	buffer := &bytes.Buffer{}
	target := n.node
	length := len(pattern)
	for path, next := "", 0; next < length; {
		path, next = TraversePath(pattern, next)
		if path == "" {
			if _, ok := n.hash["/"]; ok {
				log.Panic(rootMethodNotOnly)
			}
			target.method = handler
			n.hash["/"] = "/"
			break
		}
		buffer.WriteByte('/')

		priority, ended := static, next == length
		if path[0] == '*' || path[0] == ':' {
			priority = dynamic
			buffer.WriteString("%v")
			if path[0] == '*' {
				if !ended {
					log.Panicf(endpointErr, pattern)
				}
				priority = elastic
			}
			path = path[1:]
		} else {
			buffer.WriteString(path)
		}

		if !CheckPath(path) {
			log.Panicf(pathNotMatched, path, pattern)
		}

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
				pattern:  path,
				priority: priority,
				children: make([]*node, 0, 10),
				parent:   target,
			}
			if priority == elastic && target.wideMethod == nil {
				nn.parent.wideMethod = handler
			}
			target.children = target.children.add(nn)
		}

		if ended {
			flag := buffer.String()
			if val, ok := n.hash[flag]; ok {
				log.Panicf(conflictedErr, val, pattern)
			}
			nn.method = handler
			n.hash[flag] = pattern
			break
		}
		target = nn
	}
}

func (n *trie) match(pattern string) (handler Handle, params *Params) {
	var dynamics []string
	target := n.node
	ended, length := false, len(pattern)
	path, next := "", 0
walk:
	path, next = TraversePath(pattern, next)
	if path == "" {
		handler = n.method
		return
	}
	ended = next == length
	for i := target.index; i < len(target.children); i++ {
		nn := target.children[i]
		if nn.priority == elastic {
			if nn.method != nil {
				handler = nn.method
				if !ended {
					path += "/" + pattern[next:]
				}
				dynamics = append(dynamics, path)
			}
			target = nn
			break
		}
		if isDynamic := nn.priority == dynamic; isDynamic || nn.pattern == path {
			if isDynamic {
				dynamics = append(dynamics, path)
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
	l := len(dynamics)
	p := make(Params, l)
	params = &p
	for last := target; last != nil; last = last.parent {
		last.index = 0
		if handler != nil {
			if last.priority == dynamic || last.priority == elastic {
				l--
				params.setIndex(last.pattern, dynamics[l], l)
			}
		}
	}
	return
}
