package mux

import (
	"path"
	"strings"
)

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

func countPartsAndParams(path string) (uint8, uint8) {
	length := len(path)
	if length == 0 {
		return 0, 0
	}
	var nPart, nParams uint
	for i := 0; i < length-1; i++ {
		if path[i] != '/' {
			continue
		}

		nPart++
		if path[i+1] != ':' && path[i+1] != '*' {
			continue
		}
		nParams++
	}

	if nPart == 0 {
		return 1, 0
	}

	if nPart > 255 {
		nPart = 255
	}

	if nParams > 255 {
		nParams = 255
	}
	return uint8(nPart), uint8(nParams)
}

func countParts(path string) uint8 {
	length := len(path)
	if length == 0 {
		return 0
	}

	var n uint
	for i := 0; i < length; i++ {
		if path[i] == '/' {
			n++
		}
	}

	if n > 255 {
		n = 255
	}
	return uint8(n)
}

func between(path string, start int) (part string, next int) {
	length := len(path) - start
	if length == 0 {
		return "",0
	}

	if path[start] == '/'{
		start = start + 1
	}

	for next = start; next < len(path); next++ {
		if path[next] == '/' {
			break
		}
	}
	part = path[start:next]
	return
}

func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isAlnum(ch byte) bool {
	return isAlpha(ch) || ('0' <= ch && ch <= '9')
}

func isPath(path string) bool {
	l := len(path)
	switch l {
	case 0:
		return false
	case 1:
		return isAlpha(path[0])
	default:
		l = l - 1
		if !isAlpha(path[0]) || !isAlnum(path[l]) {
			return false
		}
		for i := 1; i < l; i++ {
			if !isAlnum(path[i]) && path[i] != '.' && path[i] != '_' && path[i] != '-' {
				return false
			}
		}
	}
	return true
}
