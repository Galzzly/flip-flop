package main

import (
	"fmt"
	"image"
	"strings"
)

func Part1(input string) any {
	lines := strings.Split(input, "\n")
	instr := lines[0]
	sushi := make([]image.Point, 0, len(lines[2:]))
	for _, line := range lines[2:] {
		var x, y int
		fmt.Sscanf(line, "%d,%d", &x, &y)
		sushi = append(sushi, image.Point{X: x, Y: y})
	}
	P := image.Point{X: 0, Y: 0}
	S := 0
	for _, c := range instr[:len(instr)/2] {
		switch c {
		case '>':
			P.X++
		case '<':
			P.X--
		case '^':
			P.Y++
		case 'v':
			P.Y--
		}
		if P == sushi[S] {
			S++
		}
	}
	return S
}

func Part2(input string) any {
	lines := strings.Split(input, "\n")
	instr := lines[0]
	var sushi []image.Point
	for _, line := range lines[2:] {
		if line == "" {
			continue
		}
		var x, y int
		fmt.Sscanf(line, "%d,%d", &x, &y)
		sushi = append(sushi, image.Point{X: x, Y: y})
	}

	dir := map[rune]image.Point{
		'>': {X: 1}, '<': {X: -1}, '^': {Y: 1}, 'v': {Y: -1},
	}

	// Body as a deque: index 0 = head, last = tail. occupied mirrors it for
	// O(1) collision checks.
	head := image.Point{X: 0, Y: 0}
	body := []image.Point{head}
	occupied := map[image.Point]bool{head: true}
	si := 0

	for _, c := range instr {
		d, ok := dir[c]
		if !ok {
			continue
		}
		newHead := head.Add(d)
		if si < len(sushi) && newHead == sushi[si] {
			// Eat: grow by keeping the tail. (Sushi is never on the snake.)
			body = append([]image.Point{newHead}, body...)
			occupied[newHead] = true
			head = newHead
			si++
			continue
		}
		// Normal move: the tail vacates, so free it before checking collision
		// (moving into the current tail cell is allowed).
		tail := body[len(body)-1]
		delete(occupied, tail)
		if occupied[newHead] {
			return len(body) // head hit its body — snake dies at this length
		}
		body = append([]image.Point{newHead}, body[:len(body)-1]...)
		occupied[newHead] = true
		head = newHead
	}
	return len(body)
}

func Part3(input string) any {
	lines := strings.Split(input, "\n")
	instr := lines[0]
	var sushi []image.Point
	for _, line := range lines[2:] {
		if line == "" {
			continue
		}
		var x, y int
		fmt.Sscanf(line, "%d,%d", &x, &y)
		sushi = append(sushi, image.Point{X: x, Y: y})
	}
	dir := map[rune]image.Point{
		'>': {X: 1}, '<': {X: -1}, '^': {Y: 1}, 'v': {Y: -1},
	}
	head := image.Point{X: 0, Y: 0}
	body := []image.Point{head} // head..tail
	occupied := map[image.Point]bool{head: true}
	si := 0
	selfEats := 0

	for _, c := range instr {
		d, ok := dir[c]
		if !ok {
			continue
		}
		newHead := head.Add(d)

		// Eat sushi -> grow (keep tail).
		if si < len(sushi) && newHead == sushi[si] {
			body = append([]image.Point{newHead}, body...)
			occupied[newHead] = true
			head = newHead
			si++
			continue
		}

		tail := body[len(body)-1]
		switch {
		case newHead == tail && len(body) > 1:
			// Move onto the vacating tail: a normal move.
			delete(occupied, tail)
			body = append([]image.Point{newHead}, body[:len(body)-1]...)
		case occupied[newHead]:
			// Self-eat: remove the segment moved onto and every segment behind
			// it; the head advances onto that tile and the remaining front
			// shifts forward (no growth).
			k := 0
			for i, p := range body {
				if p == newHead {
					k = i
					break
				}
			}
			for j := k - 1; j < len(body); j++ {
				delete(occupied, body[j])
			}
			body = append([]image.Point{newHead}, body[:k-1]...)
			selfEats++
		default:
			// Normal move.
			delete(occupied, tail)
			body = append([]image.Point{newHead}, body[:len(body)-1]...)
		}
		occupied[newHead] = true
		head = newHead
	}
	return len(body) * selfEats
}
