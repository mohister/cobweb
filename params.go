package cobweb

type Params []*param

type param struct {
	key string
	val string
}

func (p *Params) Get(key string) (string, bool) {
	for _, entry := range *p {
		if entry != nil && entry.key == key {
			return entry.val, true
		}
	}
	return "", false
}

func (p *Params) Set(key, val string) {
	for _, entry := range *p {
		if entry != nil && entry.key == key {
			entry.val = val
			return
		}
	}
	*p = append(*p, &param{key: key, val: val})
}

func (p *Params) setIndex(key, val string, index int) {
	(*p)[index] = &param{key: key, val: val}
}
