package main

import (
	"os"
)

func readGarbage(s string, ix int) (int, int) {
	erased := 0
	for ix < len(s) {
		if s[ix] == '>' {
			return ix + 1, erased
		} else if s[ix] == '!' {
			ix += 2
		} else {
			erased++
			ix++
		}
	}
	panic("should not happen")
}

func parseGroups(s string) ([]string, int, int) {
	stack := make([]int, 0, 1)
	ix := 0
	groups := make([]string, 0, 1)
	score := 0
	erased := 0
	for ix < len(s) {
		if s[ix] == '{' {
			stack = append(stack, ix)
			ix++
			continue
		} else if s[ix] == '}' {
			group := s[stack[len(stack)-1] : ix+1]
			groups = append(groups, group)
			score += len(stack)
			stack = stack[:len(stack)-1]
			ix++
			continue
		} else if s[ix] == '<' {
			var er int
			ix, er = readGarbage(s, ix+1)
			erased += er
		} else if s[ix] == ',' {
			ix++
			continue
		} else {
			fatalf("unexpected character: %q", s[ix:])
		}
	}

	return groups, score, erased
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	printf("===== part 1 =====")

	for _, line := range lines {
		groups, score, erased := parseGroups(line)
		printf("line: %q, groups: %d", line, len(groups))
		printf("score: %d", score)
		printf("erased: %d", erased)
	}

}
