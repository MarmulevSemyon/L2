package main

func main() {
	ch := make(chan int)
	go func() {
		defer close(ch) // закрытие каналаx
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}
