package cobweb

import "net/http"

const (
	static uint = iota
	dynamic
	elastic
)

type Handle func(http.ResponseWriter, *http.Request, Params)

type node struct {
	pattern    string
	priority   uint
	method     Handle
	children   children
	wideMethod Handle
	parent     *node
	index      int
}

type children []*node

func (n children) add(child *node) children {
	index, l := 0, len(n)
	nc := make(children, l+1)
	switch l {
	case 0:
		nc[index] = child
		return nc
	case 1:
		if child.priority >= n[0].priority {
			nc[0], nc[l] = n[0], child
		} else {
			nc[0], nc[l] = child, n[0]
		}
		return nc
	}

	switch child.priority {
	case static:
		index = 0
	case elastic:
		index = l
	case dynamic:
		for index = range n {
			if n[index].priority == child.priority {
				break
			}
		}
	}
	nc[index] = child
	if index != 0 {
		copy(nc[:index], n[:index])
	}
	if index != l {
		copy(nc[index+1:], n[index:])
	}
	return nc
}
