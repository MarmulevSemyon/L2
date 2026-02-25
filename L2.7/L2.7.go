package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)
	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			// если условия для разных case будут одновременно выполнсяться,
			// то мы не можем предсказать какой из case выполнится первым
			// (Если одна или более операций могут выполниться, выбирается
			// одна из них через равномерное псевдослучайное распределение.)
			select {
			case v, ok := <-a:
				if ok {
					c <- v
				} else {
					a = nil
				}
			case v, ok := <-b:
				if ok {
					c <- v
				} else {
					b = nil
				}
			}
			if a == nil && b == nil {
				close(c)
				return
			}
		}
	}()
	return c
}

func main() {
	rand.Seed(time.Now().Unix())
	// создаются каналы a и b (ассинхронно)
	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4, 6, 8)
	// a и b сливаются в с (ассинхронно)
	c := merge(a, b)
	// читаются значения из с
	for v := range c {
		fmt.Print(v)
	}
}
