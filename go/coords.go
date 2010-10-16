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

func Spread(origin Point, angle, spread float, count int) chan Point {
	ch := make(chan Point)
	go func() {
		a, b := Ident(origin.Angle - angle/2), Ident(origin.Angle + angle/2)
		a, b = a.Mul(spread), b.Mul(spread)
		p, q := origin, origin
		for i := 0; i < count; i += 2 {
			p, q = p.Add(a), q.Add(b)
			ch <- p
			ch <- q
		}
	}()
	return ch
}

func Ident(angle float) Point {
	sin, cos := math.Sincos(float64(angle))
	return Point{
		X: float(sin),
		Y: float(cos),
		Angle: angle,
	}
}

func (p Point) Add(q Point) Point {
	return Point{
		X: p.X + q.X,
		Y: p.Y + q.Y,
		Angle: p.Angle, // unchanged
	}
}

func (p Point) Mul(m float) Point {
	return Point{
		X: p.X*m,
		Y: p.Y*m,
		Angle: p.Angle, // unchanged
	}
}

func Distance(a, b Point) float {
	aa := a.X-b.X
	aa = aa * aa
	bb := a.Y-b.Y
	bb = bb * bb
	return float(math.Sqrt(float64(aa+bb)))
}
