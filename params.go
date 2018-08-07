package cobweb

import (
	"strconv"
)

type entries []string

func (e entries) add(entry string) {
	l := len(e)
	e = e[:l+1]
	e[l] = entry
}

type Params struct {
	keys   entries
	values entries
}

func (p *Params) String(key string) string {
	for i := range p.keys {
		if p.keys[i] == key {
			return p.values[i]
		}
	}
	return ""
}

func (p *Params) Int(key string) int {
	if val := p.String(key); val != "" {
		num, _ := strconv.Atoi(val)
		return num
	}
	return 0
}

func (p *Params) Int64(key string) int64 {
	if val := p.String(key); val != "" {
		num, _ := strconv.ParseInt(val, 10, 0)
		return num
	}
	return 0
}

func (p *Params) Float64(key string) float64 {
	if val := p.String(key); val != "" {
		num, _ := strconv.ParseFloat(val, 0)
		return num
	}
	return 0
}

func (p *Params) Bool(key string) bool {
	if val := p.String(key); val != "" {
		b, _ := strconv.ParseBool(val)
		return b
	}
	return false
}
