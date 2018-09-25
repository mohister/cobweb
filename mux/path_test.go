package mux

import "testing"

func TestBetween2(t *testing.T) {
	addUrl := []string{
		"////accounts////",
		"accounts////:account",
		"////accounts////:account/projects",
		"/accounts////:account/projects/:project////",
		"/accounts/:account/projects/:project/:user/:status////",
	}
	path := ""
	for i := range addUrl {
		pathURL := CleanPath(addUrl[i])
		start, length := 0, len(pathURL)
		t.Log(addUrl[i], " => ", pathURL)
		for start < length {
			path, start = between(pathURL, start)
			t.Log("	"+path)
		}
	}
}

func BenchmarkBetween(b *testing.B) {
	addUrl := []string{
		"////accounts////",
		"accounts////:account",
		"////accounts////:account/projects",
		"/accounts////:account/projects/:project////",
		"/accounts/:account/projects/:project/:user/:status////",
	}
	for n := 0; n < b.N; n++ {
		for i := range addUrl {
			pathURL := CleanPath(addUrl[i])
			start, length := 0, len(pathURL)
			for start < length {
				_, start = between(pathURL, start)
			}
		}
	}
}

