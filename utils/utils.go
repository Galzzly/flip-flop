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

func IntPow(n, m int) int {
	if m == 0 {
		return 1
	}
	if n == 0 {
		return 0
	}
	if m == 1 {
		return n
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
}

func Ter[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func Biggest[T Number](a, b T) T {
	return Ter(a < b, b, a)
}
