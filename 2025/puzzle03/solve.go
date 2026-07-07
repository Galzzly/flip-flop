package main

import (
	"fmt"
	"strings"
)

func Part1(input string) any {
	res := make(map[string]int)
	for _, v := range strings.Split(input, "\n") {
		res[v]++
	}
	var ans string
	count := 0
	for k, v := range res {
		if v > count {
			ans = k
			count = v
		}
	}
	return ans
}

func Part2(input string) any {
	ans := 0
	for _, v := range strings.Split(input, "\n") {
		var R, G, B int
		fmt.Sscanf(v, "%d,%d,%d", &R, &G, &B)
		if (G > R && G > B) && !(R == G || B == G || R == B) {
			ans++
		}
	}
	return ans
}

func Part3(input string) any {
	ans := 0
	for _, v := range strings.Split(input, "\n") {
		var R, G, B int
		fmt.Sscanf(v, "%d,%d,%d", &R, &G, &B)
		if R == G || B == G || R == B {
			ans += 10
			continue
		}
		if R > G && R > B {
			ans += 5
			continue
		}
		if G > R && G > B {
			ans += 2
			continue
		}
		ans += 4
	}
	return ans
}
