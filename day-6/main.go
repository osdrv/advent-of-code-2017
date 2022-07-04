package main

import (
	"bytes"
	"os"
	"strconv"
)

func computeBlockSignature(blocks []int) string {
	var buf bytes.Buffer
	for _, block := range blocks {
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(strconv.Itoa(block))
	}
	return buf.String()
}

func computeFullCycle(input []int) (int, []int) {
	blocks := make([]int, len(input))
	copy(blocks, input)

	step := 0
	trace := make(map[string]struct{})
	trace[computeBlockSignature(blocks)] = struct{}{}
	for {
		maxix := 0
		for ix := 0; ix < len(blocks); ix++ {
			if blocks[ix] > blocks[maxix] {
				maxix = ix
			}
		}
		freeblocks := blocks[maxix]
		blocks[maxix] = 0
		curix := maxix
		for freeblocks > 0 {
			curix = (curix + 1) % len(blocks)
			blocks[curix]++
			freeblocks--
		}
		step++
		signature := computeBlockSignature(blocks)
		if _, ok := trace[signature]; ok {
			break
		}
		trace[signature] = struct{}{}
	}

	return step, blocks
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	nums := parseInts(lines[0])
	printf("nums: %+v", nums)

	printf("===== part 1 =====")

	nCycles, blocks := computeFullCycle(nums)
	printf("terminated after %d cycles", nCycles)
	printf("exit blocks: %+v", blocks)

	printf("===== part2 =====")
	nLoop, blocks := computeFullCycle(blocks)
	printf("loop size: %d", nLoop)
}
