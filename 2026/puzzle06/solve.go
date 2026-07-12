package main

import (
	"image"
	"math/big"
	"strings"
)

func Part1(input string) any {
	grid := make(map[image.Point]rune, 0)
	var S image.Point
	for y, line := range strings.Split(input, "\n") {
		for x, c := range line {
			if c != 'S' && c != '#' && c != '*' {
				continue
			}
			if c == 'S' {
				S = image.Point{X: x, Y: y}
			}
			grid[image.Point{X: x, Y: y}] = c
		}
	}

	// Flood-fill the gear mesh from S. S spins counter-clockwise (L); every
	// adjacent gear spins the opposite way, so spin alternates with distance.
	// spin value: false = L (counter-clockwise), true = R (clockwise).
	dirs := []image.Point{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}}
	spin := map[image.Point]bool{S: false}
	queue := []image.Point{S}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		opposite := !spin[p]
		for _, d := range dirs {
			n := p.Add(d)
			if r := grid[n]; r == '#' || r == 'S' {
				if _, ok := spin[n]; !ok {
					spin[n] = opposite
					queue = append(queue, n)
				}
			}
		}
	}

	// Read lights top-to-bottom, left-to-right into a binary string (first
	// light = most-significant bit). A light is high (1) if a neighbouring gear
	// spins R, low (0) if L, and off (skipped) if no neighbour spins. All four
	// neighbours of a light share the same x+y parity, so they can never
	// disagree on spin direction.
	var bits strings.Builder
	for y, line := range strings.Split(input, "\n") {
		for x, c := range line {
			if c != '*' {
				continue
			}
			p := image.Point{X: x, Y: y}
			for _, d := range dirs {
				if r, ok := spin[p.Add(d)]; ok {
					if r {
						bits.WriteByte('1')
					} else {
						bits.WriteByte('0')
					}
					break
				}
			}
		}
	}

	answer, ok := new(big.Int).SetString(bits.String(), 2)
	if !ok {
		return 0 // no lights were on
	}
	return answer
}

func Part2(input string) any {
	lines := strings.Split(input, "\n")
	isGear := func(c rune) bool { return c == '#' || c == 'S' || c == '3' }

	kind := make(map[image.Point]rune)
	var S image.Point
	for y, line := range lines {
		for x, c := range line {
			if isGear(c) || c == '*' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				kind[image.Point{X: x, Y: y}] = c
				if c == 'S' {
					S = image.Point{X: x, Y: y}
				}
			}
		}
	}

	dirs := []image.Point{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}}

	// Build the gear graph: 4-directional meshing plus Bluetooth links that join
	// an input's adjacent gear to its matching output's adjacent gears. Every
	// edge (meshing or Bluetooth) flips the rotation direction.
	adj := make(map[image.Point][]image.Point)
	inputGear := make(map[rune]image.Point)
	outputGears := make(map[rune][]image.Point)
	for p, c := range kind {
		switch {
		case isGear(c):
			for _, d := range dirs {
				if isGear(kind[p.Add(d)]) {
					adj[p] = append(adj[p], p.Add(d))
				}
			}
		case c >= 'a' && c <= 'z':
			for _, d := range dirs {
				if isGear(kind[p.Add(d)]) {
					inputGear[c] = p.Add(d)
				}
			}
		case c >= 'A' && c <= 'Z':
			for _, d := range dirs {
				if isGear(kind[p.Add(d)]) {
					outputGears[c] = append(outputGears[c], p.Add(d))
				}
			}
		}
	}
	for letter, ig := range inputGear {
		for _, og := range outputGears[letter-'a'+'A'] {
			adj[ig] = append(adj[ig], og)
			adj[og] = append(adj[og], ig)
		}
	}

	// 2-colour from S: false = L (counter-clockwise, like S), true = R.
	spin := map[image.Point]bool{S: false}
	queue := []image.Point{S}
	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		for _, n := range adj[p] {
			if _, ok := spin[n]; !ok {
				spin[n] = !spin[p]
				queue = append(queue, n)
			}
		}
	}

	// Read lights top-to-bottom, left-to-right (first = MSB). High (1) if a
	// neighbouring gear spins R, low (0) if only L, off (skipped) if none spin.
	var bits strings.Builder
	for y, line := range lines {
		for x, c := range line {
			if c != '*' {
				continue
			}
			p := image.Point{X: x, Y: y}
			on, high := false, false
			for _, d := range dirs {
				if r, ok := spin[p.Add(d)]; ok {
					on = true
					high = high || r
				}
			}
			if on {
				if high {
					bits.WriteByte('1')
				} else {
					bits.WriteByte('0')
				}
			}
		}
	}
	answer, ok := new(big.Int).SetString(bits.String(), 2)
	if !ok {
		return 0
	}
	return answer
}

func Part3(input string) any {
	lines := strings.Split(input, "\n")
	isGear := func(c rune) bool { return c == '#' || c == 'S' || c == '3' }

	kind := make(map[image.Point]rune)
	var S image.Point
	for y, line := range lines {
		for x, c := range line {
			if isGear(c) || c == '*' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
				kind[image.Point{X: x, Y: y}] = c
				if c == 'S' {
					S = image.Point{X: x, Y: y}
				}
			}
		}
	}
	dirs := []image.Point{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}}

	// Label gear-mesh connected components (Bluetooth links excluded) and their
	// sizes. comp ids start at 1; compSize is 0-indexed by id-1.
	comp := make(map[image.Point]int)
	var compSize []int
	for p, c := range kind {
		if !isGear(c) {
			continue
		}
		if _, ok := comp[p]; ok {
			continue
		}
		id := len(compSize) + 1
		compSize = append(compSize, 0)
		comp[p] = id
		q := []image.Point{p}
		for len(q) > 0 {
			cur := q[0]
			q = q[1:]
			compSize[id-1]++
			for _, d := range dirs {
				n := cur.Add(d)
				if isGear(kind[n]) {
					if _, ok := comp[n]; !ok {
						comp[n] = id
						q = append(q, n)
					}
				}
			}
		}
	}

	// Resolve Bluetooth endpoints.
	inputGear := make(map[rune]image.Point)
	outputGears := make(map[rune][]image.Point)
	for p, c := range kind {
		switch {
		case c >= 'a' && c <= 'z':
			for _, d := range dirs {
				if isGear(kind[p.Add(d)]) {
					inputGear[c] = p.Add(d)
				}
			}
		case c >= 'A' && c <= 'Z':
			for _, d := range dirs {
				if isGear(kind[p.Add(d)]) {
					outputGears[c] = append(outputGears[c], p.Add(d))
				}
			}
		}
	}

	spin := make(map[image.Point]bool)
	rotating := make(map[int]bool)
	// fill flood-fills a component from a seed gear (spin s), meshing outward.
	fill := func(seed image.Point, s bool) {
		id := comp[seed]
		rotating[id] = true
		spin[seed] = s
		q := []image.Point{seed}
		for len(q) > 0 {
			cur := q[0]
			q = q[1:]
			for _, d := range dirs {
				n := cur.Add(d)
				if comp[n] == id {
					if _, ok := spin[n]; !ok {
						spin[n] = !spin[cur]
						q = append(q, n)
					}
				}
			}
		}
	}

	isPrime := func(n int) bool {
		if n < 2 {
			return false
		}
		for i := 2; i*i <= n; i++ {
			if n%i == 0 {
				return false
			}
		}
		return true
	}

	// The starting gear's group always rotates (L), regardless of size.
	fill(S, false)

	// Cascade Bluetooth activations: a connection fires only if its input gear
	// is spinning and the destination group is non-prime (and not yet spinning).
	// Newly-spinning groups may enable more inputs, so repeat until stable.
	for changed := true; changed; {
		changed = false
		for letter, ig := range inputGear {
			s, ok := spin[ig]
			if !ok {
				continue // input's gear isn't spinning
			}
			for _, og := range outputGears[letter-'a'+'A'] {
				id := comp[og]
				if rotating[id] || isPrime(compSize[id-1]) {
					continue
				}
				fill(og, !s) // rotate opposite the received signal
				changed = true
			}
		}
	}

	// Read lights top-to-bottom, left-to-right (first = MSB).
	var bits strings.Builder
	for y, line := range lines {
		for x, c := range line {
			if c != '*' {
				continue
			}
			p := image.Point{X: x, Y: y}
			on, high := false, false
			for _, d := range dirs {
				if r, ok := spin[p.Add(d)]; ok {
					on = true
					high = high || r
				}
			}
			if on {
				if high {
					bits.WriteByte('1')
				} else {
					bits.WriteByte('0')
				}
			}
		}
	}
	answer, ok := new(big.Int).SetString(bits.String(), 2)
	if !ok {
		return 0
	}
	return answer
}
