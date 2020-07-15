package main

import (
	"fmt"
	"strings"
)

func main() {
	a := "a\"c"
	fmt.Println(strings.Count(a, "\""))
}
