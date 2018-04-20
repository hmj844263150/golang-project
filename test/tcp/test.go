package main

import (
	"fmt"
)

func main() {
	a := 0x444444000000
	fmt.Println(a)
	b := fmt.Sprintf("%x", a)
	fmt.Println(b)
}
