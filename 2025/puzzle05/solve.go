package main

import (
	"unicode"
)

var partner []int

func Part1(input string) any {
	// Each letter appears exactly twice: the two ends of one tunnel.
	// Precompute, for every position, the index of its matching end.
	partner = make([]int, len(input))
	seen := make(map[byte]int, len(input)/2)
	for i := 0; i < len(input); i++ {
		c := input[i]
		if j, ok := seen[c]; ok {
			partner[i], partner[j] = j, i
			delete(seen, c)
		} else {
			seen[c] = i
		}
	}

	// The train reads left-to-right. At position p it enters tunnel input[p]
	// and teleports to the matching end q, taking |q-p| steps (the characters
	// between the ends plus the exit step). It then advances one position to
	// the next tunnel's entrance (an uncounted step) and repeats until it
	// steps off the end of the line.
	steps := 0
	for p := 0; p < len(input); {
		q := partner[p]
		if q > p {
			steps += q - p
		} else {
			steps += p - q
		}
		p = q + 1
	}
	return steps
}

func Part2(input string) any {
	tunnels := make(map[rune]bool, len(input)/2)
	for _, c := range input {
		tunnels[c] = false
	}

	for p := 0; p < len(input); {
		q := partner[p]
		tunnels[rune(input[p])] = true
		p = q + 1
	}

	// List the unvisited tunnels in order of first appearance (each letter
	// appears twice, so dedupe as we scan).
	var ans string
	listed := make(map[rune]bool, len(input)/2)
	for _, c := range input {
		if !tunnels[c] && !listed[c] {
			ans += string(c)
			listed[c] = true
		}
	}
	return ans
}

func Part3(input string) any {
	steps := 0
	for p := 0; p < len(input); {
		q := partner[p]
		var step int
		if q > p {
			step = q - p
		} else {
			step = p - q
		}
		if unicode.IsLower(rune(input[p])) {
			steps += step
		} else {
			steps -= step
		}
		p = q + 1
	}
	return steps
}
