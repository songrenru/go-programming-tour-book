package main

import "fmt"

func main() {
	var t test
	fmt.Printf("t = %v", t)
}

type test struct {
	name string
	age int
}