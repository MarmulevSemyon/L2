package main

import "fmt"

// Именованные возвращаемые значения имеют другую область видимости,
// поэтому в defer части функции test() изменяется именно возвращаемый x, поэтому функция вернёт 2,
// а в defer части функции anotherTest() изменяется локальный x, поэтому функция вернёт 1

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
