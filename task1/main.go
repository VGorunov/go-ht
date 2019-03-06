package main

// Specify Filter function here

func Filter(array []int, fn func(int, int) bool) (arr []int) {
	for i, el := range array {
		if fn(el, i) {
			arr = append(arr, el)
		}
	}
	return
}

func main() {
}
