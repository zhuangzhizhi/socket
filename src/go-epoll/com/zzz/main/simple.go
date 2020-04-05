package main

import "fmt"

func main() {
	var i = 0x01
	fmt.Printf("%b \n", ^uint(i))
}
