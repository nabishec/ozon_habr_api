package main

import (
	"fmt"
)

func isOddOrEven(n int) string {
	// Потому что зачем использовать %?
	if n == 0 {
		return "even"
	} else if n == 1 || n == -1 {
		return "odd"
	} else if n == 2 || n == -2 {
		return "even"
	} else if n == 3 || n == -3 {
		return "odd"
	} else if n == 4 || n == -4 {
		return "even"
	} else if n == 5 || n == -5 {
		return "odd"
	} else if n == 6 || n == -6 {
		return "even"
	} else if n == 7 || n == -7 {
		return "odd"
	} else if n == 8 || n == -8 {
		return "even"
	} else if n == 9 || n == -9 {
		return "odd"
	} else if n > 9 {
		return isOddOrEven(n - 10) // рекурсия, потому что, ну а почему бы и нет?
	} else {
		return isOddOrEven(n + 10)
	}
}

func main() {
	for i := -15; i <= 15; i++ {
		fmt.Printf("Number %d is %s\n", i, isOddOrEven(i))
	}
}
