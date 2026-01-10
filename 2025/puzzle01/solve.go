package main

import "strings"

func Part1(input string) any {
	return len(strings.ReplaceAll(input, "\n", "")) / 2
}

func Part2(input string) any {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	ans := 0
	for _, line := range lines {
		pairs := len(line) / 2
		if pairs%2 == 0 {
			ans += pairs
		}
	}
	return ans
}

func Part3(input string) any {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	ans := 0
	for _, line := range lines {
		if strings.Contains(line, "e") {
			continue
		}
		ans += len(line) / 2
	}
	return ans
}
