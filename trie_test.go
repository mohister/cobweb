package cobweb

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAddNode(t *testing.T) {
	addUrl := []string{
		"/accounts/",
		"/accounts/:account",
		"/accounts/:account/projects",
		"/accounts/:account/projects/:project",
		"/accounts/:account/projects/:project/:user/:status",
		"/accounts/:account/projects/:project/files/*file",
		"/ccounts/",
		"/ccounts/:account",
		"/ccounts/:account/projects",
		"/ccounts/:account/projects/:project",
		"/ccounts/:account/projects/:project/files/*file",
	}

	matched := []string{
		//"/accounts/",
		//"/accounts/account",
		//"/accounts/account/projects",
		//"/accounts/account/projects/project",
		//"/accounts/account/projects/project/files/file",
		"/accounts/account/projects/project/files",
	}
	node := NewTrie()
	for i := range addUrl {
		node.addNode(addUrl[i], new(http.HandlerFunc))
	}
	byt, _ := json.Marshal(node)
	t.Log(string(byt))

	for i := range matched {
		h, m := node.match(matched[i])
		if h != nil {
			t.Log(m.String("account"))
			t.Log(m.String("project"))
			t.Log(m.String("file"))
		} else {
			t.Log("not matched")
		}
	}
}

func BenchmarkNewNode(b *testing.B) {
	b.StopTimer()
	addUrl := []string{
		"/accounts/",
		"/accounts/:account",
		"/accounts/:account/projects",
		"/accounts/:account/projects/:project",
		"/accounts/:account/projects/:project/files/*file",
		"/ccounts/",
		"/ccounts/:account",
		"/ccounts/:account/projects",
		"/ccounts/:account/projects/:project",
		"/ccounts/:account/projects/:project/files/*file",
	}

	matched := []string{
		"/accounts/",
		"/accounts/account",
		"/accounts/account/projects",
		"/accounts/account/projects/project",
		"/accounts/account/projects/project/files/file",
		"/accounts/account/projects/project/files/file/abc/def",
	}
	node := NewTrie()
	for i := range addUrl {
		node.addNode(addUrl[i], new(http.HandlerFunc))
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for i := range matched {
			node.match(matched[i])
		}
	}
}

func TestTrim(t *testing.T) {
	strArr := []string{
		"////////////pattern////////////",
		"///////////////////////////////",
		"////////////pat/tern////////////",
		"////////////pat//////tern////////////",
	}

	for i := range strArr{
		t.Log(strArr[i], trim(strArr[i]))
	}
}
