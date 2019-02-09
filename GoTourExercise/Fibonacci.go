package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	var fibonacci int = 0
	var fn1 int = 1
	var fn2 int = 0
	return func() int {
		if fibonacci == 0 {
			fibonacci = fn1
			return fn2
		} else if fibonacci == 1 && fn2 != 0 {
			fibonacci = fn1
			fibonacci++
			return fn1
		} else {
			var value int = fn1 + fn2
			fn2 = fn1
			fibonacci = value
			fn1 = fibonacci
			return fibonacci
		}
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
