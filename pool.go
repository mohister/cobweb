package cobweb

import (
	"sync"
	)

var maxParams uint8

func setMaxParams(max uint8) {
	if maxParams < max {
		maxParams = max
	}
}

var paramsPool = &sync.Pool{
	New: func() interface{} {
		return &Params{keys:nil,values: make(entries, 0, maxParams)} },
}

func getParams() *Params {
	return paramsPool.Get().(*Params)
}

func putParams(p *Params) {
	p.keys = nil
	p.values = p.values[:0]
	paramsPool.Put(p)
}
