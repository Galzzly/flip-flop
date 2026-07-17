package main

import (
	"strconv"
	"strings"
)

func parseInts(s string) []int {
	fields := strings.Fields(s)
	out := make([]int, len(fields))
	for i, f := range fields {
		out[i], _ = strconv.Atoi(f)
	}
	return out
}

// parseGame splits the input into the called numbers and the 5x5 cards (each
// card is one line of 25 numbers in reading order).
func parseGame(input string) (called []int, cards [][]int) {
	sections := strings.SplitN(input, "\n\n", 2)
	called = parseInts(sections[0])
	for _, line := range strings.Split(sections[1], "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		cards = append(cards, parseInts(line))
	}
	return called, cards
}

// gridLines returns every full straight line (5 cells) in a 5^dim grid: all
// axis-aligned lines plus every diagonal spanning any subset of the dimensions.
// Cells are flat base-5 indices (most-significant axis first). dim 2 -> 12
// lines, dim 3 -> 109, dim 4 -> 888.
func gridLines(dim int) [][5]int {
	pow5 := make([]int, dim)
	for i, p := dim-1, 1; i >= 0; i, p = i-1, p*5 {
		pow5[i] = p
	}
	start := func(d int) []int {
		switch d {
		case 1:
			return []int{0} // a +1 axis must start at 0
		case -1:
			return []int{4} // a -1 axis must start at 4
		default:
			return []int{0, 1, 2, 3, 4} // a 0 axis is constant
		}
	}
	var lines [][5]int
	dir := make([]int, dim)
	pos := make([]int, dim)
	var emit func(axis int) // fix each axis' start cell, then record the line
	emit = func(axis int) {
		if axis == dim {
			var line [5]int
			for k := 0; k < 5; k++ {
				idx := 0
				for i := 0; i < dim; i++ {
					idx += (pos[i] + k*dir[i]) * pow5[i]
				}
				line[k] = idx
			}
			lines = append(lines, line)
			return
		}
		for _, s := range start(dir[axis]) {
			pos[axis] = s
			emit(axis + 1)
		}
	}
	var dirs func(axis int) // enumerate canonical directions (first nonzero +ve)
	dirs = func(axis int) {
		if axis == dim {
			for _, v := range dir {
				if v != 0 {
					if v > 0 {
						emit(0)
					}
					return
				}
			}
			return // all-zero vector
		}
		for v := -1; v <= 1; v++ {
			dir[axis] = v
			dirs(axis + 1)
		}
	}
	dirs(0)
	return lines
}

// solve marks called numbers on grids of dimension dim (2 = cards, 3 = cubes,
// 4 = hypercube), each built from consecutive 25-number planes, and returns the
// number after which at least 5 bingo lines are complete across all grids.
func solve(input string, dim int) int {
	called, planes := parseGame(input)
	per := 1 // planes per grid = 5^(dim-2)
	for i := 2; i < dim; i++ {
		per *= 5
	}
	var grids [][]int
	for i := 0; i+per <= len(planes); i += per {
		flat := make([]int, 0, per*25)
		for j := 0; j < per; j++ {
			flat = append(flat, planes[i+j]...)
		}
		grids = append(grids, flat)
	}
	lines := gridLines(dim)
	marked := make(map[int]bool)
	for _, n := range called {
		marked[n] = true
		total := 0
		for _, g := range grids {
			for _, line := range lines {
				all := true
				for _, cell := range line {
					if !marked[g[cell]] {
						all = false
						break
					}
				}
				if all {
					total++
				}
			}
		}
		if total >= 5 {
			return n
		}
	}
	return -1
}

func Part1(input string) any { return solve(input, 2) } // 5x5 cards
func Part2(input string) any { return solve(input, 3) } // 5x5x5 cubes
func Part3(input string) any { return solve(input, 4) } // 5x5x5x5 hypercube
