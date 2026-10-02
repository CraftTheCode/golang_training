package main

import (
	"errors"
	"fmt"
	"golang_training/12-practice-package"
	"golang_training/11-practice"
)

type Account struct {
	Owner string
}

func MultipleReturn() (int, string) {
	return 1, "hello"
}

func Divide(a, b int) (int, error) {
	// if b == 0 {
	// 	return 0, fmt.Errorf("cannot divide by zero")
	// }
	// return a / b, nil

	if b == 0 {
		return 0, errors.New("cannot divide by zero")
	}
	return a / b, nil

}

func main() {
	fmt.Println("Hello Go world!")

	practice.Experiment()
}
