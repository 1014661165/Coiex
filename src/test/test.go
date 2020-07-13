package main

import "fmt"

func main() {
	a := make([]string, 0)
	fmt.Println(len(a))
	fmt.Println(cap(a))
	a = append(a, "1")
	fmt.Println(len(a))
	fmt.Println(cap(a))
}
