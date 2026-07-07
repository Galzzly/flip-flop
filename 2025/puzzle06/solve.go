package main

import (
	"fmt"
	"strings"
)

type Point struct {
	X, Y int
}

func Part1(input string) any {
	count := 0
	for _, bird := range strings.Split(input, "\n") {
		var x, y int
		fmt.Sscanf(bird, "%d,%d", &x, &y)
		x *= 100
		y *= 100
		x %= 1000
		y %= 1000
		if x < 0 {
			x += 1000
		}
		if y < 0 {
			y += 1000
		}
		if (x > 249 && x < 750) && (y > 249 && y < 750) {
			count++
		}
	}
	return count
}

func Part2(input string) any {
	birds := make([]Point, 0, len(input))
	for _, bird := range strings.Split(input, "\n") {
		var x, y int
		fmt.Sscanf(bird, "%d,%d", &x, &y)
		x *= 3600
		y *= 3600
		x %= 1000
		y %= 1000
		birds = append(birds, Point{X: x, Y: y})
	}

	count := 0
	for _, bird := range birds {
		X, Y := 0, 0
		for i := 0; i < 1000; i++ {
			X += bird.X
			X %= 1000
			if X < 0 {
				X += 1000
			}
			Y += bird.Y
			Y %= 1000
			if Y < 0 {
				Y += 1000
			}
			if (X > 249 && X < 750) && (Y > 249 && Y < 750) {
				count++
			}
		}
	}
	return count
}

func Part3(input string) any {
	birds := make([]Point, 0, len(input))
	for _, bird := range strings.Split(input, "\n") {
		var x, y int
		fmt.Sscanf(bird, "%d,%d", &x, &y)
		x *= 31556926
		y *= 31556926
		x %= 1000
		y %= 1000
		birds = append(birds, Point{X: x, Y: y})
	}

	count := 0
	for _, bird := range birds {
		X, Y := 0, 0
		for i := 0; i < 1000; i++ {
			X += bird.X
			X %= 1000
			if X < 0 {
				X += 1000
			}
			Y += bird.Y
			Y %= 1000
			if Y < 0 {
				Y += 1000
			}
			if (X > 249 && X < 750) && (Y > 249 && Y < 750) {
				count++
			}
		}
	}
	return count
}
