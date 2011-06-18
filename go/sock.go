package main

import (
	"io"
	"json"
	"log"
	"os"
	"reflect"
	"websocket"
)

var SockHandler = websocket.Handler(SockServer)

func SockServer(ws *websocket.Conn) {
	// create and send channel to muxer
	ch := make(MessageChannel)
	Incoming <- ch
	// get id from muxer
	id, ok := (<-ch).(int)
	if !ok {
		log.Println("got unexpected type waiting for id")
		return
	}
	// send welcome message
	go func() { ch <- Welcome{Id: id} }()
	// start read/write loops
	go readMessages(id, ws)
	writeMessages(ws, ch)
}

type inMsg struct {
	Update  *Update
	Message *Message
}

func readMessages(id int, r io.Reader) {
	dec := json.NewDecoder(r)
	for {
		var blob inMsg
		err := dec.Decode(&blob)
		if err != nil {
			if err == os.EOF {
				Incoming <- Closed{Id: id}
				return
			}
			log.Println("decode error:", err)
			continue
		}
		if blob.Update != nil {
			blob.Update.Id = id
			Incoming <- *blob.Update
		}
		if blob.Message != nil {
			blob.Message.Id = id
			Incoming <- *blob.Message
		}
	}
}

func writeMessages(w io.Writer, ch chan interface{}) {
	enc := json.NewEncoder(w)
	for m := range ch {
		t := reflect.TypeOf(m)
		enc.Encode(map[string]interface{}{
			"type": t.Name(),
			"data": m,
		})
	}
}
