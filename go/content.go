package main

import (
	"log"
	"math"
	"powerhouse"
	"sync"
	"time"
)

var _ = log.Print

const (
	startRadius = 200
	startThemes = 25
	minItems = 10
	displayDelay = 100e6
	triggerDistance = 20
	themeSize = 10
	itemSize = 7
	itemSpread = 30
)

var (
	contentIds       = make(chan int)
	contentGroupIds  = make(chan int)
	rootContentGroup *ContentGroup
)

func init() {
	powerhouse.ApiKey = "a5863c45a7818ed"
	go func() {
		for i := 0; ; i++ {
			contentIds <- i
		}
	}()
	go func() {
		for i := 0; ; i++ {
			contentGroupIds <- i
		}
	}()
	go func() {
		// FIXME: this could really be up in the var decl above
		// but that causes throw: init sleeping for no good reason
		rootContentGroup = NewContentGroup(loadStartContent)
	}()
}

func ContentLayer(ch MessageChannel) (inch MessageChannel) {
	inch = make(MessageChannel)
	go serveContent(inch, ch)
	return
}

func serveContent(inch, ch MessageChannel) {
	visible := make(map[int]*ContentGroup)

	// serve initial content items
	rootContentGroup.Send(ch, displayDelay)
	visible[rootContentGroup.Id] = rootContentGroup

	var lastSent *ContentItem

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
		up := Point{X: u.X, Y: u.Y}
		for _, cg := range visible {
			ci := cg.Closest(up, triggerDistance)
			if ci == nil {
				continue
			}
			if ci != lastSent {
				ci.Send(ch)
				lastSent = ci
			}
			// load children
			child := ci.Children()
			if _, ok := visible[child.Id]; ok {
				continue // already loaded
			}
			go child.Send(ch, displayDelay)
			visible[child.Id] = child
		}
	}
}

type ContentItem struct {
	mu       sync.Mutex
	content  *Content
	children *ContentGroup
}

func (ci *ContentItem) Children() *ContentGroup {
	ci.mu.Lock()
	if ci.children == nil {
		ci.children = NewContentGroup(func() []*Content {
			return loadChildren(ci.content)
		})
	}
	ci.mu.Unlock()
	return ci.children
}

func (ci *ContentItem) Send(ch MessageChannel) {
	d := Display{Id: ci.content.Id}
	switch data := ci.content.data.(type) {
	case *powerhouse.Theme:
		d.Title = data.Title
	case *powerhouse.Item:
		if data.Summary != nil {
			d.Body = *data.Summary
		}
		d.URL = data.Permanent_URL
		if data.Num_Multimedia == 0 {
			break
		}
		go func() {
			time.Sleep(1e6) // delay so that Display is sent first
			ch <- Image{data.Multimedia()}
		}()
	}
	ch <- d
}

type ContentGroup struct {
	Id	int
	mu      sync.Mutex
	content []*ContentItem
}

func NewContentGroup(loadFn func() []*Content) *ContentGroup {
	cg := &ContentGroup{Id: <-contentGroupIds}
	cg.mu.Lock()
	go func() {
		content := loadFn()
		cg.content = make([]*ContentItem, len(content))
		for i, c := range content {
			cg.content[i] = &ContentItem{content: c}
		}
		cg.mu.Unlock()
	}()
	return cg
}

func (cg *ContentGroup) loaded() bool {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	return cg.content != nil
}

func (cg *ContentGroup) Closest(u Point, max float) *ContentItem {
	if !cg.loaded() {
		return nil
	}
	var smallest float = max + 1
	var content *ContentItem
	for _, ci := range cg.content {
		d := Distance(u, Point{X:ci.content.X, Y:ci.content.Y})
		if d <= max && d < smallest {
			smallest = d
			content = ci
		}
	}
	return content
}

func (cg *ContentGroup) Send(ch MessageChannel, ns int64) {
	for !cg.loaded() {
		return
	}
	for _, c := range cg.content {
		ch <- *c.content
		time.Sleep(ns)
	}
}

func loadStartContent() []*Content {
	themes := powerhouse.GetThemes(startThemes, minItems)
	c := make([]*Content, len(themes))
	circ := Circle(Point{}, startRadius, len(themes))
	colors := Colors(Grey, 255, false, len(themes))
	for i, theme := range themes {
		id := <-contentIds
		coord := <-circ
		color := <-colors
		c[i] = &Content{
			Id: id,
			X: coord.X,
			Y: coord.Y,
			Angle: coord.Angle,
			Size: themeSize,
			Color: color.String(),
			color: color,
			data: theme,
		}
	}
	return c
}

func loadChildren(c *Content) (content []*Content) {
	switch data := c.data.(type) {
	case *powerhouse.Theme:
		items := data.Items()
		content = make([]*Content, len(items))
		spread := Spread(Point{c.X, c.Y, c.Angle}, math.Pi/startThemes, itemSpread, len(items))
		colors := Colors(c.color, 50, true, len(items))
		for i, item := range items {
			id := <- contentIds
			coord := <-spread
			color := <-colors
			content[i] = &Content{
				Id: id,
				X: coord.X,
				Y: coord.Y,
				Angle: coord.Angle,
				Size: itemSize,
				Color: color.String(),
				color: color,
				data: item,
			}
		}
	}
	return
}
