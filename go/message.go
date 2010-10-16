package main

import "powerhouse"

type Update struct {
	Id                    int
	Name                  string
	X, Y, Angle, Momentum float
}

type Message struct {
	Id      int
	Message string
}

type Welcome struct {
	Id int
}

type Closed struct {
	Id int
}

type Content struct {
	Id          int
	X, Y, Angle float
	data        interface{} // store for API data
}

type Display struct {
	Title, Body, Image string
}

type Image struct {
	Multimedia []*powerhouse.Multimedia
}

