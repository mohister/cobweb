package router

import (
	"log"
	"testing"
)

func TestCountParts(t *testing.T) {
	str := "/a/b/c/d"
	t.Log(str, CountParts(str, '/'))
	str = "/a/b/c/d/"
	t.Log(str, CountParts(str, '/'))
	str = "a/b/c/d/"
	t.Log(str, CountParts(str, '/'))
	str = "a/b/c/d"
	t.Log(str, CountParts(str, '/'))
	str = "/aa/bb/cc/dd"
	t.Log(str, CountParts(str, '/'))
	str = "//aa///bb///cc//dd//"
	t.Log(str, CountParts(str, '/'))
	str = ""
	t.Log(str, CountParts(str, '/'))
	str = "/"
	t.Log(str, CountParts(str, '/'))
	str = "//"
	t.Log(str, CountParts(str, '/'))
}

func TestNextPart(t *testing.T) {
	path := "//abc//def//ghi//jkl"
	part, start, ended := "", 0, false
	for !ended {
		part, start, ended = NextPart(path, start,'/')
		log.Println(part, start, ended)
	}

	path = "///////////pattern/////////////"
	part, start, ended = "", 0, false
	for !ended {
		part, start, ended = NextPart(path, start,'/')
		log.Println(part, start, ended)
	}

	path = "a/b/c/d"
	part, start, ended = "", 0, false
	for !ended {
		part, start, ended = NextPart(path, start,'/')
		log.Println(part, start, ended)
	}
}
