package main

func main() {
	fac := CreateFactory("coldWithTruck")
	fac.CreateStore().Store()
	fac.CreateTransport().Delivery()
}
