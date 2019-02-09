package main

//copy paste code to https://tour.golang.org/moretypes/18
import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	//var pic [dy][dx]uint8
	pic := make([][]uint8, dy)
	for i := 0; i < len(pic); i++ {
		row := make([]uint8, dx)
		pic[i] = row
		for j := 0; j < len(pic[i]); j++ {
			pic[i][j] = uint8(i ^ j)
		}
	}
	return pic
}

func main() {
	pic.Show(Pic)
}
