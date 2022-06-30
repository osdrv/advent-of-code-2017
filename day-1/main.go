package main

import (
	"os"
)

func computeRepeatingNextSum(s string, off int) int {
	res := 0

	ptr := 0
	lastptr := len(s) - 1

	for ptr <= lastptr {
		nextptr := ptr + off
		if nextptr > lastptr {
			nextptr %= len(s)
		}
		next := s[nextptr]

		if s[ptr] == next {
			res += int(s[ptr] - '0')
		}
		ptr++
	}

	return res
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	printf("===== part 1 =====")
	for _, line := range lines {
		res := computeRepeatingNextSum(line, 1)
		printf("line: %s, res: %d", line, res)
	}

	printf("===== part 2 =====")
	for _, line := range lines {
		res := computeRepeatingNextSum(line, len(line)/2)
		printf("line: %s, res: %d", line, res)
	}
}
