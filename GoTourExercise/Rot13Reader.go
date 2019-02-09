package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rot13 rot13Reader) Read(array []byte) (int, error) {
	reader := rot13.r
	var length, end = reader.Read(array)
	for i := 0; i < length; i++ {
		var elem = array[i]
		if elem != 32 {
			if elem > 65 && elem < 90 {
				if elem+13 < 90 {
					elem = elem + 13
				} else {
					elem = elem - 13
				}
			} else {
				if elem+13 < 122 {
					elem = elem + 13
				} else {
					elem = elem - 13
				}
			}
		}
		array[i] = byte(elem)
	}
	return length, end
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
