package main

import (
	"fmt"
	"github.com/chinx/mohism/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	port := 8080
	if len(os.Args) > 1 {
		if nport, _ := strconv.Atoi(os.Args[1]); nport > 0 {
			port = nport
		}
	}
	m :=  mux.New()
	m.Post("/pages", func(resp http.ResponseWriter, req *http.Request) {
		pageURL := fmt.Sprintf("/pages/%d.html",time.Now().UnixNano())
		content,err := ioutil.ReadAll(req.Body)
		if err == nil {
			err = ioutil.WriteFile("static"+ pageURL, content, 0640)
		}

		fullURL := req.URL.Host+"/resource"+pageURL
		if err != nil{
			resp.WriteHeader(http.StatusInternalServerError)
			resp.Write([]byte(err.Error()))
			return
		}
		resp.WriteHeader(http.StatusCreated)
		resp.Write([]byte(fullURL))
	})
	m.ServeFiles("/resource", http.Dir("./static"))
	m.ServeFiles("/source", http.Dir("./pages"))
	http.ListenAndServe(fmt.Sprintf(":%d", port), m)
}
