package main

import (
	"math/big"
	"strings"
)

func Part1(input string) any {
	lines := strings.Split(input, "\n")
	evolutions := make(map[string]string, 0)
	for _, line := range lines {
		s := strings.Split(line, " ")
		if _, ok := evolutions[s[0]]; ok {
			continue
		}
		evolutions[s[0]] = strings.Join(s[1:], "")
	}
	S := map[string]int{"A": 1, "B": 1, "C": 0}
	for i := 0; i < 7; i++ {
		newS := map[string]int{"A": 0, "B": 0, "C": 0}
		for k, v := range S {
			for _, c := range evolutions[k] {
				newS[string(c)] += v
			}
		}
		S = newS
	}
	ans := 0
	for _, v := range S {
		ans += v
	}
	return ans
}

func pairEvolveCount(input string, generations int) *big.Int {
	evolutions := make(map[string]string)
	for _, line := range strings.Split(input, "\n") {
		s := strings.Split(line, " ")
		if len(s) < 2 {
			continue
		}
		e := strings.Join(s[2:], "")
		evolutions[s[0]+s[1]] = e
		evolutions[s[1]+s[0]] = e
	}

	pairs := map[string]*big.Int{"AB": big.NewInt(1)}
	for gen := 0; gen < generations; gen++ {
		next := make(map[string]*big.Int)
		for pair, c := range pairs {
			seq := string(pair[0]) + evolutions[pair] + string(pair[1])
			for i := 0; i+1 < len(seq); i++ {
				k := seq[i : i+2]
				if next[k] == nil {
					next[k] = new(big.Int)
				}
				next[k].Add(next[k], c)
			}
		}
		pairs = next
	}

	total := big.NewInt(1)
	for _, c := range pairs {
		total.Add(total, c)
	}
	return total
}

func Part2(input string) any {
	return pairEvolveCount(input, 7)
}

func Part3(input string) any {
	return pairEvolveCount(input, 21)
}
