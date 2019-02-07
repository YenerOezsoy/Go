package main

import "fmt"

func send(ch chan int) {
	var number int = 8
	ch <- number
}

func main() {
	var channel chan int = make(chan int)
	go send(channel)
	var receive int = <- channel
	fmt.Printf("%d\n", receive)
}
