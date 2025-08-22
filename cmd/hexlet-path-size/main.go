package main

import (
	"code"
	"fmt"
)

func main() {
	size, err := code.GetSize("C:\\Users\\Admin\\Desktop\\English\\useful-language-module4.pdf")
	if err != nil {
		return
	}
	fmt.Println(size)
}
