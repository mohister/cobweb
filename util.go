package cobweb

import (
	"path"
	"strings"
)

// NextSeparator returns an index of next separator in path.
func NextSeparator(path string, start int) int {
	for start < len(path) {
		if c := path[start]; c == '/' || c == TerminationCharacter {
			break
		}
		start++
	}
	return start
}

func CleanPath(p string) string {
	if p == "" {
		return p
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	if p[len(p)-1] == '/' && np != "/" {
		if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
			np = p[:len(p)-1]
		}
	}
	return np
}
