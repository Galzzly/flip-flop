package main

import (
	"fmt"
	"strings"
)

func shortestPaths(a, b int) int {
	if a <= 0 || b <= 0 {
		return 0
	}
	type point struct{ x, y int }
	start := point{0, 0}
	end := point{a - 1, b - 1}
	// BFS out from the start. dist holds the shortest distance to each cell;
	// ways holds how many distinct shortest paths reach it. When we first see
	// a cell we inherit its discoverer's path count; when we reach it again at
	// the same distance (an alternative shortest path) we add to that count.
	dist := map[point]int{start: 0}
	ways := map[point]int{start: 1}
	queue := []point{start}
	dirs := [4]point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			n := point{p.x + d.x, p.y + d.y}
			if n.x < 0 || n.y < 0 || n.x >= a || n.y >= b {
				continue
			}
			nd := dist[p] + 1
			if cur, ok := dist[n]; !ok {
				dist[n] = nd
				ways[n] = ways[p]
				queue = append(queue, n)
			} else if cur == nd {
				ways[n] += ways[p]
			}
		}
	}
	return ways[end]
}

func Part1(input string) any {
	count := 0
	for _, line := range strings.Split(input, "\n") {
		var a, b int
		fmt.Sscanf(line, "%d %d", &a, &b)
		count += shortestPaths(a, b)
	}
	return count
}

func shortestPaths3D(a, b, c int) int {
	if a <= 0 || b <= 0 || c <= 0 {
		return 0
	}
	type point struct{ x, y, z int }
	start := point{0, 0, 0}
	end := point{a - 1, b - 1, c - 1}
	dist := map[point]int{start: 0}
	ways := map[point]int{start: 1}
	queue := []point{start}
	dirs := [6]point{
		{1, 0, 0}, {-1, 0, 0},
		{0, 1, 0}, {0, -1, 0},
		{0, 0, 1}, {0, 0, -1},
	}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, d := range dirs {
			n := point{p.x + d.x, p.y + d.y, p.z + d.z}
			if n.x < 0 || n.y < 0 || n.z < 0 || n.x >= a || n.y >= b || n.z >= c {
				continue
			}
			nd := dist[p] + 1
			if cur, ok := dist[n]; !ok {
				dist[n] = nd
				ways[n] = ways[p]
				queue = append(queue, n)
			} else if cur == nd {
				ways[n] += ways[p]
			}
		}
	}
	return ways[end]
}

func Part2(input string) any {
	count := 0
	for _, line := range strings.Split(input, "\n") {
		var a, b int
		fmt.Sscanf(line, "%d %d", &a, &b)
		count += shortestPaths3D(a, b, a) // z-dimension equals x
	}
	return count
}

func shortestPathsND(dims, length int) int {
	if dims <= 0 || length <= 0 {
		return 0
	}
	// Encode an N-dimensional coordinate (each component in [0,length)) as a
	// single base-`length` integer, so the same dist+ways BFS works for any
	// number of dimensions. Total cells = length^dims; the opposite corner is
	// the last index (every component = length-1).
	total := 1
	for i := 0; i < dims; i++ {
		total *= length
	}
	dist := make([]int, total)
	ways := make([]int, total)
	for i := range dist {
		dist[i] = -1
	}
	dist[0] = 0
	ways[0] = 1
	queue := []int{0}
	relax := func(from, to int) {
		nd := dist[from] + 1
		switch dist[to] {
		case -1:
			dist[to] = nd
			ways[to] = ways[from]
			queue = append(queue, to)
		case nd:
			ways[to] += ways[from]
		}
	}
	for len(queue) > 0 {
		idx := queue[0]
		queue = queue[1:]
		stride := 1
		for d := 0; d < dims; d++ {
			digit := (idx / stride) % length
			if digit > 0 {
				relax(idx, idx-stride) // step -1 along dimension d
			}
			if digit < length-1 {
				relax(idx, idx+stride) // step +1 along dimension d
			}
			stride *= length
		}
	}
	return ways[total-1]
}

func Part3(input string) any {
	count := 0
	for _, line := range strings.Split(input, "\n") {
		var dims, length int
		fmt.Sscanf(line, "%d %d", &dims, &length)
		count += shortestPathsND(dims, length)
	}
	return count
}
