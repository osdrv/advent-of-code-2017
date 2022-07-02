package main

import (
	"os"
	"sort"
)

func conformWord(s string) string {
	b := []byte(s)
	sort.Slice(b, func(i, j int) bool {
		return b[i] < b[j]
	})
	return string(b)
}

func isValidPassNoAnagrams(s string) bool {
	if !isValidPass(s) {
		return false
	}
	ws := words(s)
	wset := make(map[string]bool)
	for _, w := range ws {
		cw := conformWord(w)
		if _, ok := wset[cw]; ok {
			return false
		}
		wset[cw] = true
	}
	return true
}

func isValidPass(s string) bool {
	ws := words(s)
	wset := make(map[string]bool)
	for _, w := range ws {
		if _, ok := wset[w]; ok {
			return false
		}
		wset[w] = true
	}
	return true
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	printf("===== part 1 =====")

	cnt := 0
	for _, line := range lines {
		isValid := isValidPass(line)
		if isValid {
			cnt++
		}
	}

	printf("Total valid count: %d", cnt)

	printf("===== part 2 =====")

	cnt = 0
	for _, line := range lines {
		isValid := isValidPassNoAnagrams(line)
		if isValid {
			cnt++
		}
	}

	printf("Total refined valid count: %d", cnt)
}
