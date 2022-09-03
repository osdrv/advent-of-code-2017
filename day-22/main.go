package main

import (
	"bytes"
	"os"
)

const (
	UP int = 1 << iota
	RIGHT
	DOWN
	LEFT
)

var (
	MOVES = map[int][2]int{
		UP:    {-1, 0},
		RIGHT: {0, 1},
		DOWN:  {1, 0},
		LEFT:  {0, -1},
	}
)

const (
	Clean int = iota
	Weakened
	Infected
	Flagged
)

func turnLeft(d int) int {
	return ((d >> 1) & 0xF) | (((d << 4) >> 1) & 0xF)
}

func turnRight(d int) int {
	return ((d << 1) & 0xF) | (((d << 1) >> 4) & 0xF)
}

func reverse(d int) int {
	return ((d << 2) & 0xF) | (((d << 2) >> 4) & 0xF)
}

func parseField(lines []string) map[[2]int]int {
	field := make(map[[2]int]int)
	h, w := len(lines), len(lines[0])
	for i := -h / 2; i <= h/2; i++ {
		for j := -w / 2; j <= w/2; j++ {
			if lines[h/2+i][w/2+j] == '#' {
				field[[2]int{i, j}] = Infected
			}
		}
	}
	return field
}

func printField(field map[[2]int]int) string {
	var b bytes.Buffer
	mini, minj, maxi, maxj := ALOT, ALOT, -ALOT, -ALOT
	for p := range field {
		if p[0] < mini {
			mini = p[0]
		}
		if p[0] > maxi {
			maxi = p[0]
		}
		if p[1] < minj {
			minj = p[1]
		}
		if p[1] > maxj {
			maxj = p[1]
		}
	}

	for i := mini; i <= maxi; i++ {
		for j := minj; j <= maxj; j++ {
			p := [2]int{i, j}
			if field[p] == Infected {
				b.WriteByte('#')
			} else if field[p] == Weakened {
				b.WriteByte('W')
			} else if field[p] == Flagged {
				b.WriteByte('F')
			} else {
				b.WriteByte('.')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func burst(field map[[2]int]int, curr [3]int) ([3]int, int) {
	inf := 0
	p := [2]int{curr[0], curr[1]}
	switch field[p] {
	case Clean:
		field[p] = Weakened
		curr[2] = turnLeft(curr[2])
	case Weakened:
		field[p] = Infected
		inf = 1
	case Infected:
		field[p] = Flagged
		curr[2] = turnRight(curr[2])
	case Flagged:
		curr[2] = reverse(curr[2])
		delete(field, p)
	}
	move := MOVES[curr[2]]
	curr[0] += move[0]
	curr[1] += move[1]

	return curr, inf
}

func printDir(d int) string {
	switch d {
	case UP:
		return "up"
	case RIGHT:
		return "right"
	case DOWN:
		return "down"
	case LEFT:
		return "left"
	default:
		panic("wtf")
	}
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	field := parseField(lines)
	curr := [3]int{0, 0, UP}

	N_BURSTS := 10000000
	infs := 0
	var inf int
	for i := 0; i < N_BURSTS; i++ {
		curr, inf = burst(field, curr)
		infs += inf
		//println(printField(field))
	}

	printf("infections: %d", infs)
}
