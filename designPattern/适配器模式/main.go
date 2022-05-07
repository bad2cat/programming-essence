package main

import "fmt"

type Target interface {
	Print(s string)
}

type Adapter struct {
	at *adaptee
}

func (a *Adapter) Print(msg string) {
	a.at.PrintMsg(msg)
}

type adaptee struct {
}

func (a *adaptee) PrintMsg(msg string) {
	fmt.Println(msg)
}

func NewAdapter() *Adapter {
	return &Adapter{at:&adaptee{}}
}

func main() {
	a := NewAdapter()
	a.Print("this is adapter pattern")
}