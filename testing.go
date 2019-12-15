package main

import (
	"fmt"
)

var global string

func main() {
	global = "Test"
	get()
}

func get() {
	fmt.Println(global)
}
