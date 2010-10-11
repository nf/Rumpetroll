package main

import (
	"http"
	"io"
	"json"
	"os"
	"path"
	"websocket"
)

const (
	httpListen = ":8080"
	staticDir = "../public"
)

func StaticServer(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path[1:]
	if p == "" {
		p = "index.html"
	}
	p = path.Join(staticDir, p)
	http.ServeFile(w, r, p)
}

func SockServer(ws *websocket.Conn) {
	enc := json.NewEncoder(ws)
	// send welcome message
	welcome := map[string]interface{}{ "type": "welcome", "id": 1, }
	enc.Encode(welcome)
	io.Copy(os.Stdout, ws)
}

type Update struct {
	x, y, angle, momentum float
}

type Message struct {
	message string
}

func main() {
	http.Handle("/sock", websocket.Handler(SockServer))
	http.HandleFunc("/", StaticServer)
	http.ListenAndServe(httpListen, nil)
}

