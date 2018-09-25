package mux

import (
	"strconv"
	"sync"
)

type Parameter interface {
	String(key string) string
	Int(key string) int
	Int64(key string) int64
	Float64(key string) float64
	Bool(key string) bool
}

var maxParams uint8

func setMaxParams(max uint8) {
	if maxParams < max {
		maxParams = max
	}
}

var paramsPool = &sync.Pool{
	New: func() interface{} {
		return &params{keys: nil, values: make(entries, 0, maxParams)}
	},
}

func getParams() *params {
	return paramsPool.Get().(*params)
}

func putParams(p *params) {
	p.keys = nil
	p.values = p.values[:0]
	paramsPool.Put(p)
}

type entries []string

func (e entries) add(entry string) entries {
	l := len(e)
	e = e[:l+1]
	e[l] = entry
	return e
}

type params struct {
	keys   entries
	values entries
}

func (p *params) addValue(val string) {
	p.values = append(p.values, val)
}

func (p *params) String(key string) string {
	for i := range p.keys {
		if p.keys[i] == key {
			return p.values[i]
		}
	}
	return ""
}

func (p *params) Int(key string) int {
	if val := p.String(key); val != "" {
		num, _ := strconv.Atoi(val)
		return num
	}
	return 0
}

func (p *params) Int64(key string) int64 {
	if val := p.String(key); val != "" {
		num, _ := strconv.ParseInt(val, 10, 0)
		return num
	}
	return 0
}

func (p *params) Float64(key string) float64 {
	if val := p.String(key); val != "" {
		num, _ := strconv.ParseFloat(val, 0)
		return num
	}
	return 0
}

func (p *params) Bool(key string) bool {
	if val := p.String(key); val != "" {
		b, _ := strconv.ParseBool(val)
		return b
	}
	return false
}
