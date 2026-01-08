package main

import "fmt"

func test() (x int) {
	defer func() {
		x++
		fmt.Println("выполнился defer test")
	}()
	x = 1
	return
}

func anotherTest() int {
	var x int
	defer func() {
		x++
		fmt.Println("выполнился defer anotherTest")
	}()
	x = 1
	return x
}

func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
