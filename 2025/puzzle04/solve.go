package main

import (
	"fmt"
	"goff/utils"
	"sort"
	"strings"
)

func Part1(input string) any {
	ans := 0
	cX, cY := 0, 0
	for _, v := range strings.Split(input, "\n") {
		var X, Y int
		fmt.Sscanf(v, "%d,%d", &X, &Y)
		ans += utils.DiffNum(cX, X) + utils.DiffNum(cY, Y)
		cX, cY = X, Y
	}

	return ans
}

func Part2(input string) any {
	ans := 0
	cX, cY := 0, 0
	for _, v := range strings.Split(input, "\n") {
		var X, Y int
		fmt.Sscanf(v, "%d,%d", &X, &Y)
		if utils.DiffNum(cX, X) == utils.DiffNum(cY, Y) {
			ans += utils.DiffNum(cX, X)
			cX, cY = X, Y
			continue
		}
		ans += utils.Biggest(utils.DiffNum(cX, X), utils.DiffNum(cY, Y))
		cX, cY = X, Y
	}
	return ans
}

type Point struct {
	X, Y int
}

type Trash struct {
	Point
	distance int
}

func Part3(input string) any {
	ans := 0
	cX, cY := 0, 0
	trash := []Trash{}
	for _, v := range strings.Split(input, "\n") {
		var X, Y int
		fmt.Sscanf(v, "%d,%d", &X, &Y)
		trash = append(trash, Trash{Point{X, Y}, utils.DiffNum(cX, X) + utils.DiffNum(cY, Y)})
	}
	sort.Slice(trash, func(i, j int) bool {
		return trash[i].distance < trash[j].distance
	})

	for _, t := range trash {
		if utils.DiffNum(cX, t.X) == utils.DiffNum(cY, t.Y) {
			ans += utils.DiffNum(cX, t.X)
			cX, cY = t.X, t.Y
			continue
		}
		ans += utils.Biggest(utils.DiffNum(cX, t.X), utils.DiffNum(cY, t.Y))
		cX, cY = t.X, t.Y
	}

	return ans
}
