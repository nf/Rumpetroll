package main

import (
	"fmt"
	"rand"
	"time"
)

func init() {
	rand.Seed(time.Nanoseconds())
}

var Grey = Color{128, 128, 128} 

type Color struct {
	r, g, b int
}

func Colors(seed Color, variance int, feedback bool, n int) chan Color {
	ch := make(chan Color)
	go func() {
		for i := 0; i < n; i++ {
			c := seed.Variant(variance)
			ch <- c
			if feedback {
				seed = c
			}
		}
	}()
	return ch
}

func (c Color) String() string {
	return fmt.Sprintf("%d,%d,%d", c.r, c.g, c.b)
}

func (c Color) Variant(n int) Color {
	r, g, b := rand.Intn(n), rand.Intn(n), rand.Intn(n)
	c = Color{
		r: c.r + r - r/2,
		g: c.g + g - g/2,
		b: c.b + b - b/2,
	}
	if c.r > 255 {
		c.r = 255
	}
	if c.g > 255 {
		c.g = 255
	}
	if c.b > 255 {
		c.b = 255
	}
	if c.r < 0 {
		c.r = 0
	}
	if c.g < 0 {
		c.g = 0
	}
	if c.b < 0 {
		c.b = 0
	}
	return c
}
