package main

type director struct {
	builder HouseBuilder
}

func newDirector(builder HouseBuilder) *director {
	return &director{builder: builder}
}

func (d *director) SetBuilder(builder HouseBuilder) {
	d.builder = builder
}

func (d *director) buildWoodHouse() *House {
	d.builder.BuildDoor()
	d.builder.BuildHeating()
	d.builder.BuildWall()
	d.builder.BuildWindow()
	d.builder.BuildYard()
	return d.builder.GetHouse()
}
