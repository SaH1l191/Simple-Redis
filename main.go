package main

import (
	"fmt"
	"simple-redis/store"
)

func main() {
	fmt.Println("Hello, World!")
	s := store.NewStore()
	s.Set("a", "23")
	s.Set("b", "3444")

	val, err := s.Get("a")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Print("Value of a is ", val)
}
