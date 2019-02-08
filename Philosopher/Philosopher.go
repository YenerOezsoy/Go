package main

import (
	"fmt"
	"time"
)

const (
	FORKS  = 4
	EAT    = 2
	THINK  = 5
	SECOND = 1 * 1e9
	SLEEP  = 6
)

func createForks() chan int {
	return make(chan int, FORKS)
}

func Philosopher(id int, forks chan int) {
	for {
		fmt.Printf("Philosopher %d is hungry and tries to eat\n", id)
		var success bool = eat(forks, id)
		if success {
			fmt.Printf("Philosopher %d was successfull\n", id)
			time.Sleep(time.Duration(EAT * SECOND))
			layDownForks(forks)
			fmt.Printf("Philosopher %d is done eating\n", id)
			time.Sleep(time.Duration(SLEEP * SECOND))
		} else {
			fmt.Printf("Philosopher %d was not successfull\n", id)
			time.Sleep(time.Duration(THINK * SECOND))
		}
	}
}

func eat(forks chan int, id int) bool {
	var tookForks int = 0
	for i := 0; i < 2; i++ {
		select {
		case forks <- 1:
			fmt.Printf("Philosopher %d took one fork\n", id)
			tookForks++
		default:
			for tookForks != 0 {
				fmt.Printf("Philosopher %d layed one fork down\n", id)
				<-forks
				tookForks--
			}
			return false
		}
	}
	return true
}

func layDownForks(forks chan int) {
	<-forks
	<-forks
}

func main() {
	var forks chan int = createForks()
	go Philosopher(1, forks)
	go Philosopher(2, forks)
	go Philosopher(3, forks)
	Philosopher(4, forks)
	for {
		time.Sleep(SECOND)
	}
}
