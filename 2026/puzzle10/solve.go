package main

import (
	"strings"
)

// wrap keeps x within uint16 bounds, wrapping on overflow/underflow:
// 65535 + 10 -> 9, 0 - 2 -> 65534.
func wrap(x int) int {
	return int(uint16(x))
}

// instr is a decoded Banenalang instruction: an opcode plus up to three
// numeric arguments (register indices, an immediate value, or a label id).
type instr struct {
	op   int
	args [3]int
}

// parse decodes a program into a flat instruction list plus a label -> program
// counter map. Each label points at the next instruction after it, so a jump
// lands directly on the instruction to run next.
func parse(input string) ([]instr, map[int]int) {
	na := func(s string) int { return strings.Count(s, "na") }
	var prog []instr
	labels := make(map[int]int)
	for _, line := range strings.Split(input, "\n") {
		if strings.HasPrefix(line, "be") {
			labels[na(line[2:])] = len(prog) // next instruction after the label
			continue
		}
		S := strings.Split(line[2:], "ne")
		in := instr{op: na(S[0])}
		for j := 1; j < len(S) && j <= 3; j++ {
			in.args[j-1] = na(S[j])
		}
		prog = append(prog, in)
	}
	return prog, labels
}

// execute runs a decoded program with registers r0 and r1 seeded (all others 0).
// If limit is > 0 it stops once more than limit instructions have run. It
// returns the final registers and the number of instructions executed.
func execute(prog []instr, labels map[int]int, r0, r1, limit int) (register []int, steps int) {
	register = make([]int, 16)
	register[0] = r0
	register[1] = r1
	for pc := 0; pc >= 0 && pc < len(prog); {
		steps++
		if limit > 0 && steps > limit {
			return register, steps
		}
		in := prog[pc]
		next := pc + 1
		switch in.op {
		case 0: // Load immediate value into register. (val, dest)
			register[in.args[1]] = in.args[0]
		case 1: // Copy value from one register to another. (src, dest)
			register[in.args[1]] = register[in.args[0]]
		case 2: // Add. (a, b, dest)
			register[in.args[2]] = wrap(register[in.args[0]] + register[in.args[1]])
		case 3: // Subtract. (a, b, dest)
			register[in.args[2]] = wrap(register[in.args[0]] - register[in.args[1]])
		case 4: // Multiply. (a, b, dest)
			register[in.args[2]] = wrap(register[in.args[0]] * register[in.args[1]])
		case 5: // Modulo. (a, b, dest); modulo by zero yields 0.
			if b := register[in.args[1]]; b == 0 {
				register[in.args[2]] = 0
			} else {
				register[in.args[2]] = register[in.args[0]] % b
			}
		case 6: // Increment. (reg)
			register[in.args[0]] = wrap(register[in.args[0]] + 1)
		case 7: // Decrement. (reg)
			register[in.args[0]] = wrap(register[in.args[0]] - 1)
		case 8: // Jump to label. (label)
			next = labels[in.args[0]]
		case 9: // Jump to label if register is zero. (reg, label)
			if register[in.args[0]] == 0 {
				next = labels[in.args[1]]
			}
		case 10: // Jump to label if register is not zero. (reg, label)
			if register[in.args[0]] != 0 {
				next = labels[in.args[1]]
			}
		}
		pc = next
	}
	return register, steps
}

func Part1(input string) any {
	prog, labels := parse(input)
	register, _ := execute(prog, labels, 0, 0, 0)
	return register[0]
}

func Part2(input string) any {
	prog, labels := parse(input)
	const limit = 5_000_000
	count := 0
	for v := 0; v <= 99; v++ { // starting r0 values 0..99 inclusive
		if _, steps := execute(prog, labels, v, 0, limit); steps > limit {
			count++
		}
	}
	return count
}

// Part3 counts, over every combination of r0 in 0..65535 and r1 in 0..15, how
// many exceed 5,000,000 instructions. Line 7 seeds r15=16, making the halting
// behaviour depend only on r0 mod 16 (verified: it is perfectly periodic). So
// one representative per (r1, residue) decides all 65536/16 = 4096 values that
// share that residue, turning ~1M runs into 256.
func Part3(input string) any {
	prog, labels := parse(input)
	const (
		limit  = 5_000_000
		period = 16
		blocks = 65536 / period // r0 values sharing each residue
	)
	count := 0
	for r1 := 0; r1 < period; r1++ {
		for r0 := 0; r0 < period; r0++ {
			if _, steps := execute(prog, labels, r0, r1, limit); steps > limit {
				count += blocks
			}
		}
	}
	return count
}
