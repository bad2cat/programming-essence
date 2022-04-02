package main

import "fmt"

type TransportSimple interface {
	Delivery() string
}

type Truck struct{}

func (t *Truck) Delivery() string {
	return fmt.Sprintf("use truck delivery...")
}

type AirPlane struct{}

func (a *AirPlane) Delivery() string {
	return fmt.Sprintf("use airplane delivery...")
}
