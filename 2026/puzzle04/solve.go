package main

import (
	"goff/utils"
	"strings"
)

func Part1(input string) any {
	ans := 0
	stalk := strings.Split(input, "\n")[3:]
	stalk = stalk[:len(stalk)-400]
	for _, v := range stalk {
		ans += utils.Ter(strings.Contains(v, "o"), 1, 0)
	}
	return ans
}

func Part2(input string) any {
	ans := 0
	stalk := strings.Split(input, "\n")[3:]
	stalk = stalk[:len(stalk)-1]
	side := strings.HasPrefix(stalk[len(stalk)-1], "o")
	stalk = stalk[:len(stalk)-1]
	length := len(stalk) - 1
	for i := 0; i < len(stalk); i++ {
		if side {
			if strings.HasSuffix(stalk[length-i], "o") {
				ans++
				side = !side
				continue
			}
			continue
		}
		if strings.HasPrefix(stalk[length-i], "o") {
			ans++
			side = !side
		}
	}
	return ans
}

func Part3(input string) any {
	stalk := strings.Split(input, "\n")[3:]
	stalk = stalk[:len(stalk)-1]
	endLeft, endRight := 0, 0
	for _, line := range stalk {
		if !strings.Contains(line, "o") {
			continue
		}
		if strings.HasPrefix(line, "o") {
			if endRight > 0 {
				endRight--
			}
			endLeft++
		} else {
			if endLeft > 0 {
				endLeft--
			}
			endRight++
		}
	}
	return endLeft + endRight
}
