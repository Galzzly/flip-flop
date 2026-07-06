package utils

import (
	"strconv"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	Check(err)
	return i
}

func DiffNum[T Number](a, b T) (res T) {
	res = b - a
	if res < 0 {
		return -res
	}
	return
}
