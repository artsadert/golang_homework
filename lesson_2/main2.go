package main

import (
	"fmt"
	"os"
	"strconv"
)

func print(a ...string) {
	for _, x := range a {
		fmt.Println(x)
	}
}

func div(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}

	return a / b, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please provide 2 numbers")
		os.Exit(1)
	}

	a, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Println("cannot parse a amount")
		os.Exit(1)
	}

	b, err := strconv.ParseFloat(os.Args[2], 64)
	if err != nil {
		fmt.Println("cannot parse b amount")
		os.Exit(1)
	}

	res, err := div(a, b)
	if err != nil {
		fmt.Println("cannot divide a by b")
		os.Exit(1)
	}

	fmt.Printf("a / b = %f\n", res)
	print("hello", "world")
}
