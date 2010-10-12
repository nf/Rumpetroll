package main

import (
	"http"
	"path"
)

const (
	httpListen = ":8080"
	staticDir  = "../public"
)

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path[1:]
	if p == "" {
		p = "index.html"
	}
	p = path.Join(staticDir, p)
	http.ServeFile(w, r, p)
}

func main() {
	http.Handle("/sock", SockHandler)
	http.HandleFunc("/", StaticHandler)
	http.ListenAndServe(httpListen, nil)
}
