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
	handler := func(http.ResponseWriter, *http.Request, *Params) {}
	node := NewTrie()
	for i := range addUrl {
		node.addNode(addUrl[i], handler)
	}
	byt, _ := json.Marshal(node)
	t.Log(string(byt))

	for i := range matched {
		h, m := node.match(matched[i])
		if h != nil {
			t.Log(m.Get("account"))
			t.Log(m.Get("project"))
			t.Log(m.Get("file"))
		} else {
			t.Log("not matched")
		}
	}
}

func BenchmarkNewNode(b *testing.B) {
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
	handler := func(http.ResponseWriter, *http.Request, *Params) {}
	node := NewTrie()
	for i := range addUrl {
		node.addNode(addUrl[i], handler)
	}
	for i := 0; i < b.N; i++ {
		for i := range matched {
			node.match(matched[i])
		}
	}
}
