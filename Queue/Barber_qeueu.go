package main

import (
	"fmt"
	"time"
)

const (
	NUMBER_OF_CHAIRS = 4
	NUMBER_OF_CUSTOMER = 19
)

type Request struct {
	id int
	ack chan int
}

func Barber(queue chan Request) {
	for {
		var req Request = <- queue
		var timer int = <- req.ack
		fmt.Printf("Barber is cutting hair of customer %d \n", req.id)
		time.Sleep(time.Duration(timer * 1e9))
		fmt.Printf("Barber done cutting hair of customer %d \n", req.id)
		req.ack <- 1
	}
}

func customer(queue chan Request, id int, timer int, totalReq Request) {
	fmt.Printf("Customer %d entered the shop \n", id)
	var ack chan int = make(chan int)
	var req = Request{id, ack}
	queue <- req
	fmt.Printf("Customer %d is in queue waiting \n", id)
	ack <- timer
	totalReq.ack <- <-ack
	fmt.Println("")
}

func main() {
	var queue chan Request = make(chan Request, NUMBER_OF_CHAIRS)
	var totalAck = make(chan int)
	var totalReq = Request {-2, totalAck}
	go Barber(queue)
	go customer(queue, 0, 1, totalReq)
	go customer(queue, 1, 1, totalReq)
	go customer(queue, 2, 1, totalReq)
	go customer(queue, 3, 2, totalReq)
	go customer(queue, 4, 3, totalReq)
	go customer(queue, 5, 7, totalReq)
	go customer(queue, 6, 7, totalReq)
	go customer(queue, 7, 4, totalReq)
	go customer(queue, 8, 3, totalReq)
	go customer(queue, 9, 1, totalReq)
	go customer(queue, 10, 1, totalReq)
	go customer(queue, 11, 2, totalReq)
	go customer(queue, 12, 2, totalReq)
	go customer(queue, 13, 1, totalReq)
	go customer(queue, 14, 3, totalReq)
	go customer(queue, 15, 4, totalReq)
	go customer(queue, 16, 6, totalReq)
	go customer(queue, 17, 7, totalReq)
	go customer(queue, 18, 3, totalReq)
	var counter int = 0
	for counter != NUMBER_OF_CUSTOMER {
		<-totalReq.ack
		counter++
	}
}
