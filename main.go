package main

import (
	"fmt"
	"go-concurrent/my_rw_mutex"
)

type A struct {
	a string
	b int
}

type B struct {
	a string
}

func main() {
	fmt.Println(my_rw_mutex.GetValue())
}
