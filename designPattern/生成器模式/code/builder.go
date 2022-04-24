package main

type HouseBuilder interface {
	BuildWall()
	BuildYard()
	BuildDoor()
	BuildWindow()
	BuildHeating()
	GetHouse() *House
}

// woodBuilder use build wood house
type woodBuilder struct {
	wall    string
	door    string
	window  string
	heating string
	yard    *Yard
}

func newWoodBuilder() HouseBuilder {
	return &woodBuilder{}
}

func (w *woodBuilder) BuildWall() {
	w.wall = "white wall"
}

func (w *woodBuilder) BuildDoor() {
	w.door = "red wood"
}

func (w *woodBuilder) BuildYard() {
	w.yard = newYard(10, "wood yard")
}

func (w *woodBuilder) BuildWindow() {
	w.window = "wood window"
}

func (w *woodBuilder) BuildHeating() {
	w.heating = "simonzi"
}

func (w *woodBuilder) GetHouse() *House {
	return &House{
		Wall:    w.wall,
		Door:    w.door,
		Window:  w.window,
		Heating: w.heating,
		Yard:    w.yard,
	}
}
