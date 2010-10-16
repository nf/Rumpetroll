package main

import "math"

type Point struct {
	X, Y, Angle float
}

func Circle(center Point, radius float, count int) chan Point {
	ch := make(chan Point)
	go func() {
		for i := 0; i < count; i++ {
			angle := float64(i) / float64(count) * 2 * math.Pi + float64(center.Angle)
			sin, cos := math.Sincos(angle)
			ch <- Point{
				X: center.X + float(sin)*radius,
				Y: center.Y + float(cos)*radius,
				Angle: float(angle),
			}
		}
		close(ch)
	}()
	return ch
}

func Distance(a, b Point) float {
	aa := a.X-b.X
	aa = aa * aa
	bb := a.Y-b.Y
	bb = bb * bb
	return float(math.Sqrt(float64(aa+bb)))
}
