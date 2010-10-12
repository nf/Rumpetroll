package main

type Update struct {
	Id int
	Name string
	X, Y, Angle, Momentum float
}

type Message struct {
	Id int
	Message string
}

type Welcome struct {
	Id int
}

type Closed struct {
	Id int
}
