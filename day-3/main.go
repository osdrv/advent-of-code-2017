package main

func computeSpiralManhattan(v int) int {
	i := 0
	dim := 0
	for {
		// the spirale square dimension
		dim = 2*i + 1
		if dim*dim >= v {
			break
		}
		i++
	}

	printf("v: %d, dim: %d", v, dim)

	low, high := (dim-2)*(dim-2)+1, dim*dim
	printf("low: %d, high: %d", low, high)

	ur := low + dim - 2
	ul := ur + dim - 1
	bl := ul + dim - 1
	br := high

	printf("ur: %d, ul: %d, br: %d, bl: %d", ur, ul, br, bl)

	dist := (dim - 1) / 2

	if v <= ur {
		dist += abs(v - (ur - (dim-1)/2))
	} else if v <= ul {
		dist += abs(v - (ul - (dim-1)/2))
	} else if v <= bl {
		dist += abs(v - (bl - (dim-1)/2))
	} else {
		dist += abs(v - (br - (dim-1)/2))
	}

	return dist
}

var (
	STEPS = [][2]int{
		{0, -1},
		{-1, 0},
		{0, 1},
		{1, 0},
	}
)

func computeAdjacentCells(v int) int {
	field := make(map[Point2]int)

	sumAround := func(p Point2) int {
		if p.x == 0 && p.y == 0 {
			return 1
		}
		res := 0
		for _, s := range STEPS8 {
			res += field[Point2{p.x + s[0], p.y + s[1]}]
		}
		return res
	}

	curr := Point2{0, 0}
	nextval := 0
	stix := len(STEPS) - 1
	off := 0
	for {
		nextval = sumAround(curr)
		field[curr] = nextval
		printf("point: %v, sum: %d", curr, nextval)
		if nextval > v {
			break
		}
		if abs(curr.x+STEPS[stix][0]) > off || abs(curr.y+STEPS[stix][1]) > off {
			stix++
			if stix >= len(STEPS) {
				stix %= len(STEPS)
				curr.x++
				off++
				continue
			}
		}
		curr.x += STEPS[stix][0]
		curr.y += STEPS[stix][1]
	}

	return nextval
}

func main() {
	input := []int{
		1,
		12,
		23,
		1024,
		312051,
	}

	printf("===== part 1 =====")

	for _, v := range input {
		res := computeSpiralManhattan(v)
		printf("input: %d, result: %d", v, res)
	}

	printf("===== part 2 =====")

	printf("adjacent cells lookup: %d", computeAdjacentCells(input[len(input)-1]))
}
