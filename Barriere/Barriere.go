package main

import (
	"fmt"
	"time"
)

const (
	SECONDS = 1 * 1e9
	OFFSET  = 0 * 1e9
)

func barrier() {
	t1 := task(5, 1)
	t2 := task(1, 2)
	t3 := task(2, 3)
	t4 := task(5, 4)

	var res1 bool = <-t1
	var res2 bool = <-t2
	var res3 bool = <-t3
	var res4 bool = <-t4

	var set = []bool{res1, res2, res3, res4}
	printResult(set)
}

func task(timer int, id int) chan bool {
	var doneSignale = make(chan int)
	var res = make(chan bool)
	go func() {
		fmt.Printf("task %d is working\n", id)
		time.Sleep(time.Duration(timer * SECONDS))
		doneSignale <- 1
	}()
	t := time.After(time.Duration((timer * SECONDS) + OFFSET))
	go func() {
		select {
		case <-doneSignale:
			res <- true
		case <-t:
			res <- false
			fmt.Printf("Task %d timed out\n", id)
		}
	}()
	return res
}

func printResult(set []bool) {
	time.Sleep(SECONDS)
	fmt.Println()
	for i, element := range set {
		if element {
			fmt.Printf("Task %d completed job in time\n", i)
		} else {
			fmt.Printf("Task %d did not complete job in time\n", i)
		}

	}
}

func main() {
	barrier()
}
