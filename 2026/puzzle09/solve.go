package main

import (
	"container/heap"
	"image"
	"math"
	"strings"

	"goff/utils/graph"
)

func Part1(input string) any {
	lines := strings.Split(input, "\n")
	open := map[image.Point]bool{}
	var S, E image.Point
	for y, line := range lines {
		for x, c := range line {
			p := image.Point{X: x, Y: y}
			switch c {
			case '#':
				// wall
			case 'S':
				S = p
				open[p] = true
			case 'E':
				E = p
				open[p] = true
			default:
				open[p] = true
			}
		}
	}
	// Build an open-cell adjacency graph; Dijkstra (weight-1 edges) gives the
	// shortest number of steps from S to E.
	dirs := []image.Point{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}}
	adj := map[image.Point][]image.Point{}
	for p := range open {
		for _, d := range dirs {
			if n := p.Add(d); open[n] {
				adj[p] = append(adj[p], n)
			}
		}
	}
	return graph.Dijkstra(adj, S, E)
}

func Part2(input string) any {
	lines := strings.Split(input, "\n")
	open := map[image.Point]bool{}
	var S, E image.Point
	for y, line := range lines {
		for x, c := range line {
			p := image.Point{X: x, Y: y}
			switch c {
			case '#':
				// wall
			case 'S':
				S = p
				open[p] = true
			case 'E':
				E = p
				open[p] = true
			default:
				open[p] = true
			}
		}
	}

	// Still unit-weight: from a tile you may walk to an adjacent open tile
	// (1 step), or fire the portal gun in a direction to slide to the last
	// open tile before a wall — also 1 step, no matter how far.
	dirs := []image.Point{{X: 1}, {X: -1}, {Y: 1}, {Y: -1}}
	adj := map[image.Point][]image.Point{}
	for p := range open {
		for _, d := range dirs {
			n := p.Add(d)
			if !open[n] {
				continue
			}
			adj[p] = append(adj[p], n) // walk one tile
			far := n
			for open[far.Add(d)] {
				far = far.Add(d)
			}
			if far != n {
				adj[p] = append(adj[p], far) // portal-gun slide
			}
		}
	}
	return graph.Dijkstra(adj, S, E)
}

func Part3(input string) any {
	lines := strings.Split(input, "\n")
	H := len(lines)
	W := 0
	for _, l := range lines {
		if len(l) > W {
			W = len(l)
		}
	}
	openG := make([]bool, W*H)
	cidx := func(x, y int) int { return y*W + x }
	isOpen := func(x, y int) bool {
		return x >= 0 && x < W && y >= 0 && y < H && openG[cidx(x, y)]
	}
	var Sx, Sy, Ex, Ey int
	for y, line := range lines {
		for x, c := range line {
			switch c {
			case 'S':
				Sx, Sy = x, y
				openG[cidx(x, y)] = true
			case 'E':
				Ex, Ey = x, y
				openG[cidx(x, y)] = true
			case '.':
				openG[cidx(x, y)] = true
			}
		}
	}
	dx := [4]int{1, -1, 0, 0}
	dy := [4]int{0, 0, 1, -1}

	// Reindex open cells so states pack into small integer ids.
	posID := make([]int, W*H)
	for i := range posID {
		posID[i] = -1
	}
	var openCells []int
	for c := 0; c < W*H; c++ {
		if openG[c] {
			posID[c] = len(openCells)
			openCells = append(openCells, c)
		}
	}
	O := len(openCells)

	// Per open cell precompute:
	//   openNbr     - adjacent open cells (a normal walk step)
	//   hasWall     - whether a portal can be planted underfoot (needs a wall to
	//                 fire at so the portal's entry lands on this cell)
	//   leapTargets - the slide-endpoint reached by firing the gun in each open
	//                 direction (where the *other* portal lands)
	openNbr := make([][]int, O)
	hasWall := make([]bool, O)
	leapTargets := make([][]int, O)
	for pi, c := range openCells {
		x, y := c%W, c/W
		for di := 0; di < 4; di++ {
			nx, ny := x+dx[di], y+dy[di]
			if !isOpen(nx, ny) {
				hasWall[pi] = true
				continue
			}
			openNbr[pi] = append(openNbr[pi], posID[cidx(nx, ny)])
			ex, ey := nx, ny
			for isOpen(ex+dx[di], ey+dy[di]) {
				ex, ey = ex+dx[di], ey+dy[di]
			}
			leapTargets[pi] = append(leapTargets[pi], posID[cidx(ex, ey)])
		}
	}

	startPos := posID[cidx(Sx, Sy)]
	endPos := posID[cidx(Ex, Ey)]
	if startPos == endPos {
		return 0
	}

	// Weighted shortest path over reduced states (cell, mode):
	//   mode 0 = on foot, no portal underfoot
	//   mode 1 = standing on a portal, ready to leapfrog
	// Edges: walk to a neighbour (1), plant a portal underfoot (1), or leapfrog
	// = re-fire the far portal and step into the near one (2). Only the portal
	// you're standing on matters - each leapfrog re-fires the other - so the far
	// portal needn't be tracked, collapsing the state space to 2*open.
	const inf = math.MaxInt32
	dist := make([]int, O*2)
	for i := range dist {
		dist[i] = inf
	}
	dist[startPos*2] = 0
	pq := &graph.PriorityQueue[int]{}
	heap.Init(pq)
	heap.Push(pq, &graph.PriorityQueueItem[int]{Node: startPos * 2, Distance: 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(*graph.PriorityQueueItem[int])
		d, id := cur.Distance, cur.Node
		if d > dist[id] {
			continue
		}
		pos, mode := id/2, id%2
		if pos == endPos {
			return d
		}
		relax := func(nid, nd int) {
			if nd < dist[nid] {
				dist[nid] = nd
				heap.Push(pq, &graph.PriorityQueueItem[int]{Node: nid, Distance: nd})
			}
		}
		for _, np := range openNbr[pos] {
			relax(np*2, d+1) // walk
		}
		if mode == 0 {
			if hasWall[pos] {
				relax(pos*2+1, d+1) // plant a portal underfoot
			}
		} else {
			for _, b := range leapTargets[pos] {
				relax(b*2+1, d+2) // leapfrog
			}
		}
	}
	return -1
}
