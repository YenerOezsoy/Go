package main

//to test code copy and paste to https://tour.golang.org/moretypes/23

import "strings"

func WordCount(s string) map[string]int {
	mapped := make(map[string]int)
	var split = strings.Split(s, " ")
	for _, word := range split {
		elem, ok := mapped[word]
		if ok {
			elem++
			mapped[word] = elem
		} else {
			mapped[word] = 1
		}
	}
	return mapped
}

func main() {
	wc.Test(WordCount)
}
