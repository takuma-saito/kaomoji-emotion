package main
import ("fmt")

func Slice(list []int, x, y int) (main []int, rest []int) {
	before := make([]int, x)
	after := make([]int, len(list) - y)
	main = make([]int, y - x)
	copy(before, list[0:x])
	copy(main, list[x:y])
	copy(after, list[y:len(list)])
	rest = append(before, after...)
	return
}

func main() {
	k := 3
	x := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	for i := 1; i <= len(x) / i; i++ {
		main, rest := Slice(x, (i - 1) * k, i * k)
		fmt.Println(main, rest)
	}
}

