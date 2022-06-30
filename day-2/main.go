package main

import (
	"os"
)

func computeChechsum(matrix [][]int) int {
	res := 0
	for _, row := range matrix {
		minv, maxv := row[0], row[0]
		for ix := 1; ix < len(row); ix++ {
			minv = min(minv, row[ix])
			maxv = max(maxv, row[ix])
		}
		res += maxv - minv
	}
	return res
}

func computeEvenlyDivided(matrix [][]int) int {
	res := 0

NextRow:
	for _, row := range matrix {
		for i := 0; i < len(row); i++ {
			for j := i + 1; j < len(row); j++ {
				a, b := row[i], row[j]
				if a > b {
					a, b = b, a
				}
				if b%a == 0 {
					res += b / a
					continue NextRow
				}
			}
		}
	}
	return res
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	matrix := make([][]int, 0, len(lines))

	for _, line := range lines {
		nums := parseInts(line)
		matrix = append(matrix, nums)
	}

	printf("matrix: %+v", matrix)

	printf("===== part 1 =====")

	checksum := computeChechsum(matrix)
	printf("checksum: %d", checksum)

	printf("===== part 2 =====")
	evenlyDivided := computeEvenlyDivided(matrix)
	printf("evenly divided: %d", evenlyDivided)
}
