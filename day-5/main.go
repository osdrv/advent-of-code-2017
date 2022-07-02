package main

import (
	"os"
)

func solveMaze(nums []int) int {
	input := make([]int, len(nums))
	copy(input, nums)

	ptr := 0
	steps := 0

	for ptr >= 0 && ptr < len(input) {
		//printf("ptr: %d, nums: %+v", ptr, input)
		advance := input[ptr]
		if advance >= 3 {
			input[ptr]--
		} else {
			input[ptr]++
		}
		steps++
		ptr += advance
	}

	return steps
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	nums := make([]int, 0, len(lines))
	for _, line := range lines {
		nums = append(nums, parseInt(line))
	}

	printf("===== part 1 =====")

	res := solveMaze(nums)

	printf("number of steps: %d", res)
}
