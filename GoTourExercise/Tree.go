package main

import (
	"fmt"
	"math/rand"
	"sort"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *Tree, ch chan int) {
	ch <- t.Value
	for i := 0; i < 2; i++ {
		if i == 0 && t.Left != nil {
			Walk(t.Left, ch)
		} else if i == 1 && t.Right != nil {
			Walk(t.Right, ch)
		}
	}
}

func Walker(t *Tree, ch chan int) {
	Walk(t, ch)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *Tree) bool {
	var chan1 = make(chan int)
	var chan2 = make(chan int)
	var array1Done = make(chan int)
	var array2Done = make(chan int)

	go Walker(t1, chan1)
	go Walker(t2, chan2)

	var array1 []int
	var array2 []int

	go fillArray(&array1, chan1, array1Done)
	go fillArray(&array2, chan2, array2Done)

	<-array1Done
	<-array2Done

	return checkEquality(&array1, &array2)

}

func fillArray(array *[]int, ch chan int, done chan int) {
	var ok = true
	var elem int

	for ok {
		elem, ok = <-ch
		*array = append(*array, elem)
	}
	done <- 1
}

func checkEquality(array1 *[]int, array2 *[]int) bool {
	sort.Ints(*array1)
	sort.Ints(*array2)

	if len(*array1) != len(*array2) {
		return false
	} else {
		for i := 0; i < len(*array1); i++ {
			if (*array1)[i] != (*array2)[i] {
				return false
			}
		}
	}
	return true
}

func main() {
	/*var ch chan int = make(chan int)
	go Walker(New(1), ch)
	var ok bool = true
	var elem int
	for ok != false {
		elem, ok = <-ch
		fmt.Printf("Value: %d\n", elem)
	}*/

	var equal bool = Same(New(1), New(1))
	if equal {
		fmt.Printf("Both trees are same")
	} else {
		fmt.Printf("Trees are not same")
	}
}

// A Tree is a binary tree with integer values.
type Tree struct {
	Left  *Tree
	Value int
	Right *Tree
}

// New returns a new, random binary tree holding the values k, 2k, ..., 10k.
func New(k int) *Tree {
	var t *Tree
	for _, v := range rand.Perm(10) {
		t = insert(t, (1+v)*k)
	}
	return t
}

func insert(t *Tree, v int) *Tree {
	if t == nil {
		return &Tree{nil, v, nil}
	}
	if v < t.Value {
		t.Left = insert(t.Left, v)
	} else {
		t.Right = insert(t.Right, v)
	}
	return t
}

func (t *Tree) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Value)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}
