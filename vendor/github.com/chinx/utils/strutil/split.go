package strutil

import "errors"

var pathErr = errors.New("path not matched regex \"^[a-zA-Z](([a-zA-Z0-9_.]*[a-zA-Z0-9]+))*$\"")

func SplitPath(path string) (arr []string, err error) {
	arr = make([]string, SplitCount(path, '/'))
	index, length := 0, len(path)
	for part, next := "", 0; next < length; {
		part, next = Between(path, '/', next)
		if !IsPath(path) {
			return arr, pathErr
		}
		arr[index] = part
		index++
	}
	return
}

func SplitCount(s string, b byte) (index int) {
	length := len(s)
	find := true
	for i := 0; i < length; i++ {
		if s[i] == b {
			if find {
				continue
			}
			find = true
		} else if find {
			find = false
			index++
		}
	}
	return
}

func Between(s string, b byte, start int) (part string, next int) {
	from := -1
	next = len(s)
	for i := start; i < next; i++ {
		if from == -1 {
			if s[i] != b {
				from = i
			}
			continue
		}
		if part == "" {
			if s[i] == b {
				part = s[from:i]
			}
		} else if s[i] != b {
			next = i
			break
		}
	}
	if part == "" && from != -1 {
		part = s[from:next]
	}
	return
}
