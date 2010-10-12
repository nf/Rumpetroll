package main

type Update struct {
	Id int
	Name, X, Y, Angle, Momentum string
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
