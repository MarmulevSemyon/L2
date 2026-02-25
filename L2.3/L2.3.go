package main

import (
	"fmt"
	"os"
)

// Что выведет программа?
// Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

func Foo() error {
	var err *os.PathError = nil
	fmt.Println(err == nil)
	return err //неявное преобразование в интерфейс error с "типом" *os.PathError и "значением" nil
}

func main() {
	err := Foo()
	fmt.Println(err)
	fmt.Println(err == nil)
}
