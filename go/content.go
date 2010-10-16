package main

import (
	"container/vector"
	"log"
	"powerhouse"
	"sync"
	"time"
)

var _ = log.Print

const (
	displayDelay = 100e6
)

var (
	contentIds = make(chan int)
	rootContentGroup = NewContentGroup(loadStartContent)
)

func init() { 
	powerhouse.ApiKey = "a5863c45a7818ed"
	go func() { for i := 0;; i++ { contentIds <- i } }()
}

func contentLayer(ch MessageChannel) (inch MessageChannel) {
	inch = make(MessageChannel)
	go serveContent(inch, ch)
	return
}

func serveContent(inch, ch MessageChannel) {
	// serve initial content items
	rootContentGroup.Send(ch, displayDelay)

	// 
	var content vector.Vector
	_ = content
	for m := range inch {
		// before doing anything, forward message to muxer
		Incoming <- m
		// we only care about updates
		u, ok := m.(Update)
		if !ok {
			continue
		}
		// test if close to any content blocks
		// if so, expand and display additional content
		_ = u
	}
}

type ContentGroup struct {
	mu      sync.Mutex
	content []*Content
}

func NewContentGroup(loadFn func() []*Content) *ContentGroup {
	cg := new(ContentGroup)
	cg.mu.Lock()
	go func() {
		cg.content = loadFn()
		cg.mu.Unlock()
	}()
	return cg
}

func (cg *ContentGroup) Closest(x, y, max float) *Content {
	if !cg.loaded() { return nil }
	return nil
}

func (cg *ContentGroup) loaded() bool {
	cg.mu.Lock()
	if cg.content == nil {
		return false
	}
	cg.mu.Unlock()
	return true
}

func (cg *ContentGroup) Send(ch MessageChannel, ns int64) {
	for !cg.loaded() { return }

	for _, c := range cg.content {
		ch <- *c
		time.Sleep(ns)
	}
}

func loadStartContent() []*Content {
	n := 20 
	c := make([]*Content, n)
	circ := Circle(0, 0, 200, n)
	for i := 0; i < n; i++ {
		coord := <-circ
		id := <-contentIds
		c[i] = &Content{id, coord.X, coord.Y}
	}
	return c
}

