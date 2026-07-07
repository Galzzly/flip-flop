package main

func Part1(input string) any {
	wall := make([]int, 100)
	R := 0
	for _, v := range input {
		switch v {
		case '>':
			R = (R + 1) % len(wall)
		case '<':
			R = (R - 1 + len(wall)) % len(wall)
		}
		wall[R]++
	}
	idx, biggest := 0, 0
	for i, v := range wall {
		if v > biggest {
			biggest = v
			idx = i
		}
	}
	return (idx + 1) * biggest
}

func Part2(input string) any {
	R, W, ans := 1, 1, 0

	for i, v := range input {
		switch v {
		case '>':
			R++
			if R == 101 {
				R = 1
			}
		case '<':
			R--
			if R == 0 {
				R = 100
			}
		}
		switch input[len(input)-1-i] {
		case '>':
			W++
			if W == 101 {
				W = 1
			}
		case '<':
			W--
			if W == 0 {
				W = 100
			}
		}
		if R == W {
			ans++
		}
	}
	return ans
}

func Part3(input string) any {
	wall := make([]int, 100)
	R := 0
	for i, v := range input {
		switch v {
		case '>':
			R = (R + 1) % len(wall)
		case '<':
			R = (R - 1 + len(wall)) % len(wall)
		}
		switch input[len(input)-1-i] {
		case '>':
			R = (R - 1 + len(wall)) % len(wall)
		case '<':
			R++
			if R == len(wall) {
				R = 0
			}
		}
		wall[R]++
	}

	idx, biggest := 0, 0
	for i, v := range wall {
		if v > biggest {
			biggest = v
			idx = i
		}
	}
	return (idx + 1) * biggest
}
