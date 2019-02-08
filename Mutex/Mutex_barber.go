package main

import (
	"fmt"
	"time"
)

const (
	SEC = 1 * 1e9
)

type Mtex struct {
	channel           chan int
	customer_in_queue *bool
}

func createMtex() Mtex {
	var mutex = make(chan int, 1)
	var customer_in_queue bool = false
	return Mtex{mutex, &customer_in_queue}
}

func lockM(mutex Mtex, isBarber bool) bool {
	if isBarber && *mutex.customer_in_queue {
		mutex.channel <- 1
		*mutex.customer_in_queue = false
		return true
	} else if !isBarber {
		if !*mutex.customer_in_queue {
			*mutex.customer_in_queue = true
			mutex.channel <- 1
			return true
		}
	}
	return false
}

func unlockM(mutex Mtex) {
	<-mutex.channel
}

func barber(service chan int, mutex Mtex) {
	for {
		var isLock bool = false
		for isLock != true {
			//time.Sleep(SEC)
			isLock = lockM(mutex, true)
		}
		var customerID int = <-service
		fmt.Printf("Barber is cutting hair of customer %d\n", customerID)
		time.Sleep(SEC)
		unlockM(mutex)
	}
}

func customer(service chan int, mutex Mtex, id int) {
	var isCut bool = false
	fmt.Printf("Customer %d in shop\n", id)
	for isCut != true {
		isCut = lockM(mutex, false)
		if !isCut {
			time.Sleep(SEC)
		}
	}
	fmt.Printf("Customer %d is next\n", id)
	service <- id
	unlockM(mutex)
}

func main() {
	var mutex = createMtex()
	var service = make(chan int, 1)
	go barber(service, mutex)
	go customer(service, mutex, 0)
	go customer(service, mutex, 1)
	go customer(service, mutex, 2)

	for {
		time.Sleep(5 * SEC)
	}
}
