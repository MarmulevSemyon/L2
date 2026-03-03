package main

import (
	"fmt"
	"time"
)

// Реализовать функцию, которая будет объединять один или более каналов done (каналов сигнала завершения) в один.
// Возвращаемый канал должен закрываться, как только закроется любой из исходных каналов.
// Сигнатура функции может быть такой:
// var or func(channels ...<-chan interface{}) <-chan interface{}
func main() {
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(999*time.Millisecond),
		sig(99999999*time.Microsecond),
	)
	fmt.Printf("done after %v", time.Since(start))
}

func or(channels ...<-chan interface{}) <-chan interface{} {
	switch len(channels) { // базовое условие для рекурсии
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{}) // возвращаемый канал
	go func() {                      // конкурентно слушаем все каналы
		defer close(orDone)

		switch len(channels) {
		case 2: // если два канала, ждем ответа (закрытия) из любого из них
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default: // если каналов больше 2, слушаем 3 канала
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			// рекурсивно вызываем оставшуюся часть каналов +наш orDone,
			// чтобы завершить все последующие рекусивные вызовы функции
			// (когда верхний уровень закрыл orDone, внутренние уровни тоже “увидят” закрытие и завершатся.)
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}
