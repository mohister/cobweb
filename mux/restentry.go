package mux

import "net/http"

const (
	static uint8 = iota
	dynamic
	elastic
)

type restEntry struct {
	h          http.Handler
	pattern    string
	wide       http.Handler
	children   entryList
	paramsKeys []string
	index      int
	priority   uint8
	maxParts   uint8
}

func (n *restEntry) SetMaxParts(max uint8) {
	if n.maxParts < max {
		n.maxParts = max
	}
}

type entryList []*restEntry

func (n entryList) add(child *restEntry) entryList {
	index, l := 0, len(n)
	nc := make(entryList, l+1)
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
