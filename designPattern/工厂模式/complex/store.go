package main

import "fmt"

type Store interface {
	Store() string
}

type ColdStore struct{}

func (t *ColdStore) Store() string {
	return fmt.Sprintf("use cold store...")
}

type HotStore struct{}

func (a *HotStore) Store() string {
	return fmt.Sprintf("use hot store...")
}
