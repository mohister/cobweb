package strutil

func TrimRight(s string, b byte) string {
	l := len(s)
	if l == 0 {
		return s
	}
	for ; l > 0; l-- {
		if s[l-1] != b {
			break
		}
	}
	return s[:l]
}

func TrimLeft(s string, b byte) string {
	l := len(s)
	if l == 0 {
		return s
	}
	i := 0
	for ; i < l; i++ {
		if s[i] != b {
			break
		}
	}
	return s[i:]
}

func Trim(s string, b byte) string {
	return TrimRight(TrimLeft(s, b), b)
}

func TrimBytesRight(s []byte, b byte) []byte {
	l := len(s)
	if l == 0 {
		return s
	}
	for ; l > 0; l-- {
		if s[l-1] != b {
			break
		}
	}
	return s[:l]
}

func TrimBytesLeft(s []byte, b byte) []byte {
	l := len(s)
	if l == 0 {
		return s
	}
	i := 0
	for ; i < l; i++ {
		if s[i] != b {
			break
		}
	}
	return s[i:]
}

func TrimBytes(s []byte, b byte) []byte {
	return TrimBytesRight(TrimBytesLeft(s, b), b)
}
