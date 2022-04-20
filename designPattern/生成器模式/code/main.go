package main

import "fmt"

func main() {
	houseBuilder := newWoodBuilder()
	director := newDirector(houseBuilder)

	house := director.buildWoodHouse()

	fmt.Printf("%s", house)
}
