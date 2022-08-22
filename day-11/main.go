package main

import (
	"os"
	"strings"
)

type Step uint8

func (s Step) String() string {
	return stepToStr[s]
}

const (
	_ Step = iota
	N
	NE
	SE
	S
	SW
	NW
)

type Coord struct {
	x, y, z int
}

var CUBE_STEPS = map[Step]Coord{
	N:  {-1, 1, 0},
	NE: {0, 1, -1},
	SE: {1, 0, -1},
	S:  {1, -1, 0},
	SW: {0, -1, 1},
	NW: {-1, 0, 1},
}

var (
	strToStep = map[string]Step{
		"n": N, "ne": NE, "se": SE, "s": S, "sw": SW, "nw": NW,
	}
	stepToStr = map[Step]string{
		N: "n", NE: "ne", SE: "se", S: "s", SW: "sw", NW: "nw",
	}
)

func parseSteps(s string) []Step {
	ss := strings.Split(s, ",")
	steps := make([]Step, 0, len(ss))
	for _, st := range ss {
		step, ok := strToStep[st]
		if !ok {
			fatalf("unrecognized step: %q", step)
		}
		steps = append(steps, step)
	}
	return steps
}

func followSteps(steps []Step) Coord {
	curr := Coord{0, 0, 0}
	for _, step := range steps {
		ss := CUBE_STEPS[step]
		curr.x += ss.x
		curr.y += ss.y
		curr.z += ss.z
	}
	return curr
}

func cubeDist(a, b Coord) int {
	return (abs(a.x-b.x) + abs(a.y-b.y) + abs(a.z-b.z)) / 2
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	println("===== part 1 =====")

	for _, line := range lines {
		steps := parseSteps(line)
		dest := followSteps(steps)
		printf("steps: %+v", steps)
		printf("dest: %v", dest)
		distance := cubeDist(Coord{0, 0, 0}, dest)
		printf("distance: %d", distance)
	}

	println("===== part 2 =====")

	for _, line := range lines {
		steps := parseSteps(line)
		curr := Coord{0, 0, 0}
		maxDist := 0
		maxix := 0
		for ix, step := range steps {
			st := CUBE_STEPS[step]
			curr.x += st.x
			curr.y += st.y
			curr.z += st.z
			if dist := cubeDist(Coord{0, 0, 0}, curr); dist > maxDist {
				maxix = ix
				maxDist = dist
			}
		}
		printf("max distance: %d at step %d", maxDist, maxix)
	}
}
