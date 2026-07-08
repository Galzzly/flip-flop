package main

import (
	"strings"
	"sync"
	"unicode"
)

func checkpassword(pwd string) (s int) {
	if strings.ContainsFunc(pwd, unicode.IsLower) {
		s += 1
	}
	if strings.ContainsFunc(pwd, unicode.IsUpper) {
		s += 1
	}
	if strings.ContainsFunc(pwd, unicode.IsNumber) {
		s += 1
	}
	return
}

func sequential(s string) int {
	best := 0
	c := rune(s[0])
	count := 1
	for _, r := range s[1:] {
		if r == c {
			count++
			continue
		}
		if count > best {
			best = count
		}
		c = r
		count = 1
	}
	if count > best {
		best = count
	}
	return best
}

func checkruleset(pwd string) (s int) {
	s = checkpassword(pwd)
	if strings.Contains(pwd, "7") && !strings.ContainsAny(pwd, "012345689") {
		s += 7
	}
	seq := sequential(pwd)
	if seq > 2 {
		s += seq * seq
	}
	if strings.Contains(pwd, "red") || strings.Contains(pwd, "green") || strings.Contains(pwd, "blue") {
		s *= 3
	}
	return
}

func Part1(input string) any {
	best := 0
	ans := ""
	for _, pwd := range strings.Split(input, "\n") {
		s := checkpassword(pwd)
		if score := s * len(pwd); score > best {
			best = score
			ans = pwd
		}
	}
	return ans
}

func Part2(input string) any {
	best := 0
	ans := ""
	for _, pwd := range strings.Split(input, "\n") {
		s := checkruleset(pwd)
		if score := s * len(pwd); score > best {
			best = score
			ans = pwd
		}
	}
	return ans
}

func Part3(input string) any {
	best := 0
	chars := []string{"a", "b", "c", "d", "e", "f",
		"g", "h", "i", "j", "k", "l", "m", "n", "o", "p",
		"q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	res := make(chan int, len(chars))
	var wg sync.WaitGroup
	for _, c := range chars {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			total := 0
			for _, pwd := range strings.Split(input, "\n") {
				p := pwd + c
				total += checkruleset(p) * len(p)
			}
			res <- total
		}(c)
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	for r := range res {
		if r > best {
			best = r
		}
	}
	return best
}
