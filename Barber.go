package main

import (
	"fmt"
	"time"
)

func barber(work chan int) {
	for {
		var servingCustomer int = <-work
		fmt.Printf("Serving customer %d \n", servingCustomer)
		var timer int
		if servingCustomer < 4 {
			timer = 1 * 1e9
		} else {
			timer = 4 * 1e9
		}
		time.Sleep(time.Duration(timer))
		fmt.Printf("Customer %d served\n", servingCustomer)
	}
}

func customer(id int, barber1 chan int, barber2 chan int) {
	fmt.Printf("Customer %d entered shop \n", id)
	select {
	case barber1<-id:
		fmt.Printf("Customer %d chose barber1 \n", id)
		break
	case barber2<-id:
		fmt.Printf("Customer %d chose barber2 \n", id)
		break
	}
}

func main() {
	var barber1 chan int = make(chan int)
	var barber2 chan int = make(chan int)
	go barber(barber1)
	go barber(barber2)
	go customer(1, barber1, barber2)
	go customer(2, barber1, barber2)
	go customer(3, barber1, barber2)
	customer(4, barber1, barber2)
	time.Sleep(10*1e9)
}