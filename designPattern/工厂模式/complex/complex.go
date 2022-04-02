package main

type Factory interface {
	CreateTransport() TransportSimple
	CreateStore() Store
}

type CreatorA struct{}

func (c *CreatorA) CreateTransport() TransportSimple {
	return &Truck{}
}

func (c *CreatorA) CreateStore() Store {
	return &ColdStore{}
}

type CreatorB struct{}

func (c *CreatorB) CreateTransport() TransportSimple {
	return &AirPlane{}
}

func (c *CreatorB) CreateStore() Store {
	return &HotStore{}
}

func CreateFactory(typ string) Factory {
	switch typ {
	case "coldWithTruck":
		return &CreatorA{}
	case "hotWithAirPlane":
		return &CreatorB{}
	}
	return nil
}
