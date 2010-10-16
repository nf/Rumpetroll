package main

import "math"

type Point struct {
	X, Y, Angle float
}

func Circle(center Point, radius float, count int) chan Point {
	ch := make(chan Point)
	go func() {
		for i := 0; i < count; i++ {
			sin, cos := math.Sincos(float64(i) / float64(count) * 2 * math.Pi)
			ch <- Point{
				X: center.X + float(sin)*radius,
				Y: center.Y + float(cos)*radius,
			}
		}
		close(ch)
	}()
	return ch
}
