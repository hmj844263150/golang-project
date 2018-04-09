package main

import (
	"crypto/md5"
	"fmt"
	// "io"
)

func main() {
	str := "12313"

	data := []byte(str)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //[]byte16

	fmt.Println(md5str1)
}
