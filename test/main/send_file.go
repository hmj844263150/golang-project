package main

import (
	"fmt"
	"os"
)

func main() {
	_, err := os.Stat("/home/hongmingjie/iot/mpnCondif/test.go")
	if err == nil {
		fmt.Println("exist")
		return
	}
	fmt.Println("not exist")
	return
}
