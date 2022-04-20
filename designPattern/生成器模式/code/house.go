package main

import "fmt"

type House struct {
	Wall    string
	Door    string
	Window  string
	Heating string
	Yard    *Yard
}

type Yard struct {
	Size int
	Name string
}

func newYard(size int, name string) *Yard {
	return &Yard{
		Size: size,
		Name: name,
	}
}

func (h *House) String() string {
	s := fmt.Sprintf("Door:%s,Wall:%s,Window:%s,Heating:%s", h.Door, h.Wall, h.Window, h.Heating)
	if h.Yard != nil {
		s = fmt.Sprintf("%s,yard size:%d,yard name:%s", s, h.Yard.Size, h.Yard.Name)
	}
	return s
}
