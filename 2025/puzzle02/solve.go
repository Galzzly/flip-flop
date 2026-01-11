package main

import (
	"strings"
)

func Part1(input string) any {
	line := strings.TrimSpace(input)
	ans := 0
	height := 0
	for _, c := range line {
		if c == '^' {
			height++
		} else {
			height--
		}
		if height > ans {
			ans = height
		}
	}
	return ans
}

func Part2(input string) any {
	line := strings.TrimSpace(input)
	ans := 0
	height := 0
	count := 1
	for i, c := range line {
		if i == 0 {
			if c == '^' {
				height += count
			} else {
				height -= count
			}
			if height > ans {
				ans = height
			}
			continue
		}
		if line[i] == line[i-1] {
			count++
		} else {
			count = 1
		}
		if c == '^' {
			height += count
		} else {
			height -= count
		}
		if height > ans {
			ans = height
		}
	}
	return ans
}

func Part3(input string) any {
	line := strings.TrimSpace(input)
	ans := 0
	height := 0
	count := 1
	for i := range line {
		if i == 0 {
			continue
		}
		if line[i] == line[i-1] {
			count++
			if i == len(line)-1 {
				f := fibonacci(count)
				if line[i] == '^' {
					height += f
				} else {
					height -= f
				}

				if height > ans {
					ans = height
				}
			}
			continue
		}
		f := fibonacci(count)
		if line[i-1] == '^' {
			height += f
		} else {
			height -= f
		}
		if height > ans {
			ans = height
		}
		count = 1
	}
	return ans
}

// Fibonacci returns the nth Fibonacci number
func fibonacci(n int) int {
	if n <= 0 {
		return 0
	}
	if n == 1 {
		return 1
	}

	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
