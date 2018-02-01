package cobweb

func between(s string, start int) (part string, next int) {
	from := -1
	next = len(s)
	for i := start; i < next; i++ {
		if from == -1 {
			if s[i] != '/' {
				from = i
			}
			continue
		}
		if part == "" {
			if s[i] == '/' {
				part = s[from:i]
			}
		} else if s[i] != '/' {
			next = i
			break
		}
	}
	if part == "" && from != -1 {
		part = s[from:next]
	}
	return
}

func trim(s string) string {
	l := len(s)
	if l == 0 {
		return s
	}
	start, end := 0, 0
	for i := 0; i < l; {
		if start == 0 && s[i] != '/' {
			start = i
		} else {
			i++
		}
		if end == 0 && s[l-1] != '/' {
			end = l
		} else {
			l--
		}
	}

	if start == end {
		return ""
	}
	return s[start:end]
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
			if !isAlnum(path[i]) && path[i] != '.' && path[i] != '_' {
				return false
			}
		}
	}
	return true
}
