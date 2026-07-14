package graph

import (
	"container/heap"
	"image"
	"math"
)

type Queue[T any] []T

type Graph interface {
	Neighbours(p image.Point) []image.Point
}

type Grid[T any] struct {
	x, y      int
	state     map[image.Point]T
	movements []image.Point
}

func NewGrid[T any](x, y int, movements []image.Point) *Grid[T] {
	state := make(map[image.Point]T)
	return &Grid[T]{
		x:         x,
		y:         y,
		state:     state,
		movements: movements,
	}
}

func (g *Grid[T]) IsValid(x, y int) bool {
	switch {
	case x < 0, x >= g.x, y < 0, y >= g.y:
		return false
	default:
		return true
	}
}

func (g *Grid[T]) GetState(p image.Point) T {
	return g.state[p]
}

func (g *Grid[T]) SetState(p image.Point, value T) {
	if !g.IsValid(p.X, p.Y) {
		return
	}
	g.state[p] = value
}

func (g *Grid[T]) Neighbours(p image.Point) (res []image.Point) {
	for _, m := range g.movements {
		np := p.Add(m)
		if g.IsValid(np.X, np.Y) {
			res = append(res, np)
		}
	}
	return
}

func (q *Queue[T]) Put(x T) {
	*q = append(*q, x)
}

func (q *Queue[T]) Get() T {
	ret := (*q)[0]
	*q = (*q)[1:]
	return ret
}

func (q *Queue[T]) Empty() bool {
	return len(*q) == 0
}

func Search(g Graph, s, e image.Point) (res []image.Point) {
	var queue Queue[image.Point]
	queue.Put(s)

	from := map[image.Point]*image.Point{}
	from[s] = nil

	for !queue.Empty() {
		current := queue.Get()
		if current == e {
			break
		}
		for _, p := range g.Neighbours(current) {
			if _, ok := from[p]; !ok {
				queue.Put(p)
				from[p] = &current
			}
		}
	}
	res = []image.Point{e}
	for p := from[e]; p != nil; p = from[*p] {
		res = append(res, *p)
	}
	return
}

// PriorityQueueItem represents an item in the priority queue.
type PriorityQueueItem[N comparable] struct {
	Node     N
	Distance int
	Index    int
}

// PriorityQueue implements heap.Interface for Dijkstra's algorithm.
type PriorityQueue[N comparable] []*PriorityQueueItem[N]

func (pq PriorityQueue[N]) Len() int { return len(pq) }

func (pq PriorityQueue[N]) Less(i, j int) bool {
	return pq[i].Distance < pq[j].Distance
}

func (pq PriorityQueue[N]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue[N]) Push(x any) {
	n := len(*pq)
	item := x.(*PriorityQueueItem[N])
	item.Index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[N]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	*pq = old[0 : n-1]
	return item
}

// Dijkstra finds the shortest path from start to end node
// graph is a map of node -> list of connected nodes
// Returns the distance to the end node, or -1 if no path exists
func Dijkstra[N comparable](graph map[N][]N, start, end N) int {
	distances := make(map[N]int)

	// Initialize distances for all nodes in the graph
	for node := range graph {
		distances[node] = math.MaxInt32
	}
	distances[start] = 0

	// Also ensure end node exists in distances even if it has no outgoing edges
	if _, exists := distances[end]; !exists {
		distances[end] = math.MaxInt32
	}

	pq := make(PriorityQueue[N], 0)
	heap.Init(&pq)
	heap.Push(&pq, &PriorityQueueItem[N]{
		Node:     start,
		Distance: 0,
	})

	visited := make(map[N]bool)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*PriorityQueueItem[N])

		if current.Node == end {
			return current.Distance
		}

		if visited[current.Node] {
			continue
		}
		visited[current.Node] = true

		for _, neighbor := range graph[current.Node] {
			if visited[neighbor] {
				continue
			}

			newDist := current.Distance + 1 // Each edge has weight 1

			if newDist < distances[neighbor] {
				distances[neighbor] = newDist
				heap.Push(&pq, &PriorityQueueItem[N]{
					Node:     neighbor,
					Distance: newDist,
				})
			}
		}
	}

	// No path found
	return -1
}

// DijkstraWithPath finds the shortest path and returns both distance and the path
func DijkstraWithPath[N comparable](graph map[N][]N, start, end N) (int, []N) {
	distances := make(map[N]int)
	previous := make(map[N]N)

	// Initialize distances for all nodes in the graph
	for node := range graph {
		distances[node] = math.MaxInt32
	}
	distances[start] = 0

	// Also ensure end node exists in distances even if it has no outgoing edges
	if _, exists := distances[end]; !exists {
		distances[end] = math.MaxInt32
	}

	pq := make(PriorityQueue[N], 0)
	heap.Init(&pq)
	heap.Push(&pq, &PriorityQueueItem[N]{
		Node:     start,
		Distance: 0,
	})

	visited := make(map[N]bool)

	for pq.Len() > 0 {
		current := heap.Pop(&pq).(*PriorityQueueItem[N])

		if current.Node == end {
			// Reconstruct the path by walking predecessors back to the start.
			path := []N{}
			node := end
			for {
				path = append([]N{node}, path...)
				prev, exists := previous[node]
				if !exists {
					break
				}
				node = prev
			}
			return current.Distance, path
		}

		if visited[current.Node] {
			continue
		}
		visited[current.Node] = true

		for _, neighbor := range graph[current.Node] {
			if visited[neighbor] {
				continue
			}

			newDist := current.Distance + 1

			if newDist < distances[neighbor] {
				distances[neighbor] = newDist
				previous[neighbor] = current.Node
				heap.Push(&pq, &PriorityQueueItem[N]{
					Node:     neighbor,
					Distance: newDist,
				})
			}
		}
	}

	return -1, nil
}
