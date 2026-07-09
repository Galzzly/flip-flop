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
	side := utils.Ter(strings.HasPrefix(stalk[len(stalk)-1], "o"), true, false)
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
	ans := 0
	stalk := strings.Split(input, "\n")[3:]
	stalk = stalk[:len(stalk)-1]
	length := len(stalk) - 1
	leaves := []bool{}
	for i := 0; i < len(stalk); i++ {
		if strings.Contains(stalk[length-i], "o") {
			leaves = append(leaves, strings.HasPrefix(stalk[length-i], "o"))
		}
	}
	for len(leaves) > 0 {
		ans++
		cur := leaves[0]
		leaves = leaves[1:]
		for i := 0; i < len(leaves); i++ {
			if leaves[i] == cur {
				continue
			}
			cur = leaves[i]
			leaves = append(leaves[:i], leaves[i+1:]...)
			i--
		}

	}
	return ans
}
