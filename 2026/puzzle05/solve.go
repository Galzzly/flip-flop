package main

import (
	"runtime"
	"strings"
	"sync"
)

type Point struct {
	X, Y int
}

var path []Point

func rotateRight(c byte) byte {
	switch c {
	case '>':
		return 'v'
	case 'v':
		return '<'
	case '<':
		return '^'
	case '^':
		return '>'
	}
	return c
}

func walk(cells []byte, W, H, cx, cy int, cdir byte, maxTurns int, visited []int32, gen int32) int {
	x, y := 0, 0
	turns, distinct := 0, 0
	for {
		if x < 0 || x >= W || y < 0 || y >= H {
			return distinct + 1 // stepping off the grid counts as one more tile
		}
		idx := y*W + x
		c := cells[idx]
		if x == cx && y == cy {
			c = cdir
		}
		if visited[idx] == gen {
			if turns == maxTurns {
				return distinct
			}
			turns++
			c = rotateRight(c)
		} else {
			visited[idx] = gen
			distinct++
		}
		switch c {
		case '>':
			x++
		case '<':
			x--
		case '^':
			y--
		case 'v':
			y++
		default:
			return distinct
		}
	}
}

// buildGrid parses the input into a flat W*H byte grid (index = y*W + x), which
// is far faster to walk than a map[Point]rune.
func buildGrid(input string) (cells []byte, W, H int) {
	lines := strings.Split(input, "\n")
	H = len(lines)
	for _, l := range lines {
		if len(l) > W {
			W = len(l)
		}
	}
	cells = make([]byte, W*H)
	for y, l := range lines {
		for x := 0; x < len(l); x++ {
			cells[y*W+x] = l[x]
		}
	}
	return
}

func Part1(input string) any {
	cells, W, H := buildGrid(input)
	path = nil
	visited := make([]bool, W*H)
	x, y := 0, 0
	for {
		if x < 0 || x >= W || y < 0 || y >= H {
			path = append(path, Point{X: x, Y: y}) // off-grid exit tile
			break
		}
		idx := y*W + x
		if visited[idx] {
			break
		}
		visited[idx] = true
		path = append(path, Point{X: x, Y: y})
		switch cells[idx] {
		case '>':
			x++
		case '<':
			x--
		case '^':
			y--
		case 'v':
			y++
		default:
			return len(path)
		}
	}
	return len(path)
}

func Part2(input string) any {
	cells, W, H := buildGrid(input)
	best := len(path) // floor: the unmodified path length
	dirs := [4]byte{'<', '>', '^', 'v'}

	// Without illegal turns, only cells on the original path can change the
	// route, so those are the only change candidates worth trying.
	jobs := make(chan Point, len(path))
	for _, p := range path {
		jobs <- p
	}
	close(jobs)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			visited := make([]int32, W*H)
			var gen int32
			local := 0
			for p := range jobs {
				if p.X < 0 || p.X >= W || p.Y < 0 || p.Y >= H {
					continue // the off-grid exit tile is not a real cell
				}
				for _, d := range dirs {
					gen++
					if n := walk(cells, W, H, p.X, p.Y, d, 0, visited, gen); n > local {
						local = n
					}
				}
			}
			mu.Lock()
			if local > best {
				best = local
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	return best
}

func Part3(input string) any {
	cells, W, H := buildGrid(input)

	best := len(path) // floor: the unmodified path length
	dirs := [4]byte{'<', '>', '^', 'v'}
	rows := make(chan int, H)
	for y := 1; y < H-1; y++ { // only interior (non-edge) rows may be changed
		rows <- y
	}
	close(rows)

	var mu sync.Mutex
	var wg sync.WaitGroup
	for w := 0; w < runtime.NumCPU(); w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			visited := make([]int32, W*H) // reused across walks via the generation stamp
			var gen int32
			local := 0
			for y := range rows {
				for x := 1; x < W-1; x++ { // interior (non-edge) columns
					for _, d := range dirs {
						gen++
						if n := walk(cells, W, H, x, y, d, 3, visited, gen); n > local {
							local = n
						}
					}
				}
			}
			mu.Lock()
			if local > best {
				best = local
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	return best
}
