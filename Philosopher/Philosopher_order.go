package main

import (
	"fmt"
	"time"
)

const (
	FORK      = 4
	SECONDS   = 1 * 1e9
	EATTIME   = 2
	THINKTIME = 5
	SLEEPTIME = 6
)

type table struct {
	forks []chan int
}

func createNewTable() table {
	var forks = []chan int{make(chan int, 1), make(chan int, 1), make(chan int, 1), make(chan int, 1)}
	return table{forks}
}

func philosopher(id int, table table) {
	for {
		fmt.Printf("Filozofcu %d ac ve yemek yemek istiyor\n", id)
		var success = false
		success = eating(id, table)
		if success {
			fmt.Printf("Filozofcu %d basariliydi ve yemege basliyo\n", id)
			time.Sleep(EATTIME * SECONDS)
			fmt.Printf("Filozofcu %d yemegini bitirdi\n", id)
			layDownAllForks(table, id, true, true, getLeftHandIndex(id), getRightHandIndex(id))
			time.Sleep(SLEEPTIME * SECONDS)
		} else {
			fmt.Printf("Filozofcu %d basarili degildi ve ise devam ediyor\n", id)
			time.Sleep(THINKTIME * SECONDS)
		}
	}

}

func getLeftHandIndex(id int) int {
	return id - 1
}

func getRightHandIndex(id int) int {
	if id < FORK {
		return id
	} else {
		return 0
	}
}

func eating(id int, table table) bool {
	var leftHandIndex = getLeftHandIndex(id)
	var rightHandIndex = getRightHandIndex(id)
	var leftHand = false
	var rightHand = false
	for i := 0; i < 2; i++ {
		select {
		case table.forks[leftHandIndex] <- 1:
			fmt.Printf("Filozofcu %d solundaki catali aldi %d\n", id, leftHandIndex)
			leftHand = true
		case table.forks[rightHandIndex] <- 1:
			fmt.Printf("Filozofcu %d sagindaki catali aldi %d\n", id, rightHandIndex)
			rightHand = true
		default:
			layDownAllForks(table, id, leftHand, rightHand, leftHandIndex, rightHandIndex)
			return false
		}
	}
	return true
}

func layDownAllForks(table table, id int, leftHand bool, rightHand bool, leftHandIndex int, rightHandIndex int) {
	if leftHand {
		<-table.forks[leftHandIndex]
		fmt.Printf("Filozofcu %d sol elindeki catali birakti %d \n", id, leftHandIndex)
	}
	if rightHand {
		<-table.forks[rightHandIndex]
		fmt.Printf("Filozofcu %d sag elindeki catali birakti %d\n", id, rightHandIndex)
	}
}

func main() {
	var table = createNewTable()
	go philosopher(1, table)
	go philosopher(2, table)
	go philosopher(3, table)
	go philosopher(4, table)
	for {
		time.Sleep(SECONDS)
	}
}
