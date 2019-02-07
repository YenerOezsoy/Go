package main

import (
	"fmt"
	"time"
)

const (
	SECOND = 1 * 1e9
)

type Mutex struct {
	channel chan int
	value   int
}

func lock(mutex Mutex) int {
	mutex.channel <- 1
	return mutex.value
}

func unlock(mutex Mutex, value int) {
	mutex.value = value
	<-mutex.channel
}

func createMutex(value int) Mutex {
	var channel = make(chan int, 1)
	var mutex = Mutex{channel, value}
	return mutex
}

func increment(mutex Mutex) {
	var value int
	for {
		value = lock(mutex)
		value++
		fmt.Printf("increment value to: %d\n", value)
		time.Sleep(SECOND)
		unlock(mutex, value)
	}
}

func decrement(mutex Mutex) {
	var value int
	for {
		value = lock(mutex)
		value--
		fmt.Printf("decrement value to: %d\n", value)
		time.Sleep(SECOND)
		unlock(mutex, value)
	}
}

func main() {
	var value int = 10
	var mutex Mutex = createMutex(value)
	go increment(mutex)
	go decrement(mutex)
	for i := 0; i < 10; i++ {
		time.Sleep(SECOND)
	}
}
