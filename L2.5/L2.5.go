package main

import (
	"fmt"
)

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	// var e customError
	return nil // возвращает ниловую ссылку на customError
}

func main() {
	var err error
	// здесь происходит преобразование nil_ссылки_на_customError в интерфейс error
	err = test()
	fmt.Printf("%#v \n", err)

	// в таблице методов интерфейса лежит не только значение поля msg (nil), но и тип (customError)
	if err != nil { // поэтому не nil
		println("error")
		return
	}
	println("ok")
}
