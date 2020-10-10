package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	content,_ := ioutil.ReadFile("./testcase/Test.java")
	fmt.Println(len(content))
	for _,b := range content{
		if b == '\\'{
			fmt.Println(true)
		}
	}
}
