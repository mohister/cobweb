package cobweb

func TraversePath(s string, start int) (part string, next int) {
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

func IsAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsUnderling(ch byte) bool {
	return ch == '_'
}

func IsDot(ch byte) bool {
	return ch == '.'
}

func IsAlnum(ch byte) bool {
	return IsAlpha(ch) || IsDigit(ch)
}

func CheckPath(path string) bool {
	l := len(path)
	switch l {
	case 0:
		return false
	case 1:
		return IsAlpha(path[0])
	default:
		l = l - 1
		if !IsAlpha(path[0]) || !IsAlnum(path[l]) {
			return false
		}
		for i := 1; i < l; i++ {
			if !IsAlnum(path[i]) && !IsDot(path[i]) && !IsUnderling(path[i]) {
				return false
			}
		}
	}
	return true
}
