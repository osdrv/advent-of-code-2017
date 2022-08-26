package main

import (
	"os"
	"strings"
)

func parse(s string) (int, int) {
	chunks := strings.SplitN(s, ": ", 2)
	return parseInt(chunks[0]), parseInt(chunks[1])
}

func passThrough(layers map[int]int, maxl int, exitEarly bool, delay int) (int, bool) {
	t := delay
	hits := 0
	hit := false
	for t <= delay+maxl {
		d, ok := layers[t-delay]
		if !ok {
			t++
			continue
		}
		pos := t % (2 * (d - 1))
		//printf("time: %d, pos: %d (depth: %d)", t, pos, d)
		if pos == 0 {
			//printf("hit at time: %d, depth: %d", t, d)
			hits += t * d
			hit = true
			if exitEarly {
				return -1, true
			}
		}
		t++
	}
	return hits, hit
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	layers := make(map[int]int, len(lines))

	maxl := 0
	for _, line := range lines {
		layer, depth := parse(line)
		layers[layer] = depth
		if layer > maxl {
			maxl = layer
		}
	}

	println("===== part 1 =====")

	hits, _ := passThrough(layers, maxl, false, 0)

	printf("hits: %d", hits)

	println("===== part 2 =====")

	d := 10

	for {
		//println("===============")
		if _, hit := passThrough(layers, maxl, true, d); !hit {
			printf("passed with no caughts with delay: %d", d)
			break
		}
		d++
	}

}
