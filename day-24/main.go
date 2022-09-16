package main

import (
	"os"
	"strings"
)

func parsePair(s string) [2]int {
	ss := strings.SplitN(s, "/", 2)
	assert(len(ss) == 2, "should be of len 2")
	v1, v2 := parseInt(ss[0]), parseInt(ss[1])
	return makePair(v1, v2)
}

func makePair(v1, v2 int) [2]int {
	return [2]int{min(v1, v2), max(v1, v2)}
}

func buildMaxStrengthBridge(index map[int]map[int]int, pos int, used map[[2]int]int) int {
	maxStren := 0
	for next, cnt := range index[pos] {
		p := makePair(pos, next)
		if used[p] >= cnt {
			continue
		}
		used[p]++
		maxStren = max(maxStren, pos+next+buildMaxStrengthBridge(index, next, used))
		used[p]--
	}
	return maxStren
}

func buildMaxLenBridge(index map[int]map[int]int, pos int, used map[[2]int]int) (int, int) {
	maxLen := 0
	maxStren := 0
	for next, cnt := range index[pos] {
		p := makePair(pos, next)
		if used[p] >= cnt {
			continue
		}
		used[p]++
		length, strength := buildMaxLenBridge(index, next, used)
		length += 1
		strength += pos + next
		if length >= maxLen {
			if length == maxLen {
				if strength > maxStren {
					maxStren = strength
				}
			} else {
				maxLen = length
				maxStren = strength
			}
		}
		used[p]--
	}
	return maxLen, maxStren
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	index := make(map[int]map[int]int)
	for _, s := range lines {
		p := parsePair(s)
		v1, v2 := p[0], p[1]
		if _, ok := index[v1]; !ok {
			index[v1] = make(map[int]int)
		}
		index[v1][v2]++

		if v1 != v2 {
			if _, ok := index[v2]; !ok {
				index[v2] = make(map[int]int)
			}
			index[v2][v1]++
		}
	}

	printf("index: %+v", index)

	maxStren := buildMaxStrengthBridge(index, 0, make(map[[2]int]int))

	printf("max strength: %d", maxStren)

	maxLen, stren := buildMaxLenBridge(index, 0, make(map[[2]int]int))
	printf("max length: %d, stren: %d", maxLen, stren)
}
