package main

import (
	"image"
	"sort"
	"strconv"
	"strings"
)

// dna holds, for each node id, the child ids grown to the left, above and
// right. A value of -1 means no child (XX).
type dna struct {
	above, left, right []int
}

// parseTrees decodes the blank-line-separated DNA blocks. Each block is two
// lines: the "above" row (one token per id) and the "left id right" triples row.
func parseTrees(input string) []dna {
	val := func(s string) int {
		if s == "XX" {
			return -1
		}
		v, _ := strconv.Atoi(s)
		return v
	}
	var trees []dna
	for _, block := range strings.Split(input, "\n\n") {
		lines := strings.Split(strings.Trim(block, "\n"), "\n")
		if len(lines) < 2 {
			continue
		}
		aboveF := strings.Fields(lines[0])
		tripF := strings.Fields(lines[1])
		n := len(aboveF)
		d := dna{above: make([]int, n), left: make([]int, n), right: make([]int, n)}
		for i := 0; i < n; i++ {
			d.left[i] = val(tripF[3*i]) // tripF[3*i+1] is the id (== i)
			d.right[i] = val(tripF[3*i+2])
			d.above[i] = val(aboveF[i])
		}
		trees = append(trees, d)
	}
	return trees
}

// P is a grid position (x horizontal, y = height above ground, 0 = row 1).
type P = image.Point

type tree struct {
	d       dna
	sprouts map[P]int  // current living sprouts: position -> node id
	own     map[P]bool // every cell this tree occupies (sprouts + stems)
	alive   bool
	mass    int // segment count, recorded at death
}

// seed is a tree waiting to be planted: its ground x-position and the DNA it
// grows by (offspring always start at node id 00).
type seed struct {
	x int
	d dna
}

// runGeneration grows all seeds together on a shared grid (planted on row 0, in
// left-to-right / input order). It returns the combined biological mass once
// every tree has died, plus the seeds for the next generation: each surviving
// sprout dropped to row 0, keeping one per column (the originally-highest wins),
// carrying its parent's DNA, ordered left to right.
func runGeneration(seeds []seed) (mass int, next []seed) {
	globalOcc := map[P]bool{}
	trees := make([]*tree, len(seeds))
	for i, s := range seeds {
		start := P{X: s.x, Y: 0}
		trees[i] = &tree{d: s.d, sprouts: map[P]int{start: 0}, own: map[P]bool{start: true}, alive: true}
		globalOcc[start] = true
	}

	for age := 1; age <= 100; age++ {
		// 1. Grow every living tree, in order.
		for _, t := range trees {
			if !t.alive {
				continue
			}
			proposals := map[P]int{}
			add := func(c P, id int) {
				if id < 0 || globalOcc[c] {
					return // XX, or blocked by any tree's cell
				}
				if ex, ok := proposals[c]; !ok || id > ex {
					proposals[c] = id // highest id wins a contested cell
				}
			}
			for p, id := range t.sprouts {
				add(P{X: p.X - 1, Y: p.Y}, t.d.left[id])
				add(P{X: p.X, Y: p.Y + 1}, t.d.above[id])
				add(P{X: p.X + 1, Y: p.Y}, t.d.right[id])
			}
			t.sprouts = proposals
			for c := range proposals {
				t.own[c] = true
				globalOcc[c] = true
			}
		}
		if age < 5 {
			continue // energy only matters from age 5
		}

		// 2. Column heights of every stem (any owner) — these block sunlight.
		colYs := map[int][]int{}
		for _, t := range trees {
			for c := range t.own {
				if _, isSprout := t.sprouts[c]; !isSprout {
					colYs[c.X] = append(colYs[c.X], c.Y)
				}
			}
		}
		for x := range colYs {
			sort.Ints(colYs[x])
		}

		// 3. Per-tree energy check.
		for _, t := range trees {
			if !t.alive {
				continue
			}
			e := 0
			for c := range t.own {
				if _, isSprout := t.sprouts[c]; isSprout {
					continue // sprouts harness nothing
				}
				ys := colYs[c.X]
				above := len(ys) - sort.SearchInts(ys, c.Y+1) // stems strictly above
				h := c.Y + 1
				if h > 10 {
					h = 10
				}
				if mult := 3 - above; mult > 0 {
					e += h * mult
				}
			}
			if e < 3*len(t.own) {
				t.alive = false
				t.mass = len(t.own) // starved
			}
		}
	}

	// Tally mass and the surviving sprouts. Per column, the originally-highest
	// sprout wins and reseeds with its parent's DNA.
	type winner struct {
		y int
		d dna
	}
	best := map[int]winner{}
	for _, t := range trees {
		if t.alive {
			t.mass = len(t.own) // died of old age at 100
		}
		mass += t.mass
		for c := range t.sprouts {
			if b, ok := best[c.X]; !ok || c.Y > b.y {
				best[c.X] = winner{c.Y, t.d}
			}
		}
	}
	xs := make([]int, 0, len(best))
	for x := range best {
		xs = append(xs, x)
	}
	sort.Ints(xs)
	for _, x := range xs {
		next = append(next, seed{x: x, d: best[x].d})
	}
	return mass, next
}

func Part1(input string) any {
	total := 0
	for _, d := range parseTrees(input) {
		m, _ := runGeneration([]seed{{x: 0, d: d}}) // each tree in its own space
		total += m
	}
	return total
}

func Part2(input string) any {
	return generationMass(parseTrees(input), 1) // one generation, competing
}

// Part3 runs three competing generations — the trees, their offspring, and the
// offspring of the offspring — and returns the final generation's mass.
func Part3(input string) any {
	return generationMass(parseTrees(input), 3)
}

// generationMass plants the initial trees 10 apart, then runs `gens` competing
// generations (each reseeded from the previous one's survivors) and returns the
// combined biological mass of the final generation.
func generationMass(dnas []dna, gens int) int {
	seeds := make([]seed, len(dnas))
	for i, d := range dnas {
		seeds[i] = seed{x: 10 * i, d: d}
	}
	mass := 0
	for g := 0; g < gens; g++ {
		mass, seeds = runGeneration(seeds)
	}
	return mass
}
