package main

import "fmt"

type TransportSimple interface {
	Delivery() string
}

func CreateTransport(typ string) TransportSimple {
	switch typ {
	case "truck":
		return &Truck{}
	case "airPlane":
		return &AirPlane{}
	default:
		return nil
	}
}

type Truck struct{}

func (t *Truck) Delivery() string {
	return fmt.Sprintf("use truck delivery...")
}

type AirPlane struct{}

func (a *AirPlane) Delivery() string {
	return fmt.Sprintf("use airplane delivery...")
}
