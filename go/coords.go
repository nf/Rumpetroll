package main

import "math"

type Coords struct {
	X, Y float
}

func Circle(centerX, centerY, radius float, count int) chan Coords {
	ch := make(chan Coords)
	go func() {
		for i := 0; i < count; i++ {
			sin, cos := math.Sincos(float64(i)/float64(count)*2*math.Pi)
			ch <- Coords{centerX + float(sin)*radius, centerY + float(cos)*radius}
		}
		close(ch)
	}()
	return ch
}
