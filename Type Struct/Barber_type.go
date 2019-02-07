package Type_Struct

import (
	"fmt"
	"time"
)

type Request struct {
	id int
	ack chan int
}

func Barber(service chan Request, id int) {
	for {
		var req Request = <- service
		var timer = <- req.ack
		fmt.Printf("Barber %d working on customer %d for %d sec\n", id, req.id, timer)
		time.Sleep(time.Duration(timer * 1e9))

		if timer < 1 {
			//send no ack
		} else {
			fmt.Printf("Customer %d got his haircut\n\n", req.id)
			req.ack <- 1
		}

	}
}

func Customer(reqMain Request, id int, barber1 chan Request, barber2 chan Request, time int) {
	var ack = make(chan int)
	var request = Request{id, ack}
	select {
	case barber1 <- request:
		fmt.Printf("Customer %d chose barber 1 \n", id)
		request.ack <- time
	case barber2 <- request:
		fmt.Printf("Customer %d chose barber 2 \n", id)
		request.ack <- time
	}
	reqMain.ack <- <-ack
}

func main() {
	var ackTotal = make(chan int)
	var allDone = Request{-2, ackTotal}
	var barber1 = make(chan Request)
	var barber2 = make(chan Request)
	go Barber(barber1, 1)
	go Barber(barber2, 2)
	go Customer(allDone, 1, barber1, barber2, 1)
	go Customer(allDone, 2, barber1, barber2, 2)
	go Customer(allDone, 3, barber1, barber2, 4)
	go Customer(allDone, 4, barber1, barber2, 3)
	go Customer(allDone, 5, barber1, barber2, 5)
	var counter int = 0
	for counter != 5 {
		<-allDone.ack
		counter++
	}

}
