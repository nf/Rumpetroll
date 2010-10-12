package main

import "log"

type MessageChannel chan interface{}

var Incoming = make(MessageChannel)

func init() {
	go Muxer()
}

func Muxer() {
	count := 0
	chans := make(map[int]MessageChannel)
	for m := range Incoming {
		var id int
		switch n := m.(type) {
		case MessageChannel:
			count++
			n <- count
			chans[count] = n
		case Update:
			id = n.Id
		case Message:
			id = n.Id
		case Closed:
			id = n.Id
			close(chans[id])
			chans[id] = nil, false
		default:
			log.Stderr("unrecognized message:", m)
		}
		if id == 0 {
			continue
		}
		for chid, ch := range chans {
			if id == chid {
				continue
			}
			ch <- m
		}
	}
}
