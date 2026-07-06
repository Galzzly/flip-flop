package main

import (
	"goff/utils"
	"strings"
)

func Part1(input string) any {
	ans := 0
	for _, v := range strings.Split(input, "\n") {
		n := utils.Atoi(v)
		if n < 60 {
			ans += 60 - n
		}
	}
	return ans
}

func Part2(input string) any {
	ans := 0
	for _, v := range strings.Split(input, "\n") {
		n := utils.Atoi(v)
		if n < 60 {
			ans += 60 - n
			continue
		}
		ans += (n - 60) * 5
	}
	return ans
}

func Part3(input string) any {
	ans := 0
	nums := []int{}
	for _, v := range strings.Split(input, "\n") {
		n := utils.Atoi(v)
		nums = append(nums, n)
	}
	half := len(nums) / 2
	for i := 0; i < half; i++ {
		a := nums[i]
		b := nums[half+i]
		if a < b {
			ans += b - a
			continue
		}
		ans += (a - b) * 5
	}
	return ans
}
