package main

import (
	"strconv"
)

var (
	SALT = []byte{17, 31, 73, 47, 23}
)

func reverse(nums []byte, curr, length int) {
	from, to := curr, curr+length-1
	toIndex := func(ix int) int {
		if ix < 0 {
			return ix + len(nums)
		}
		return ix % len(nums)
	}
	for from < to {
		fromIx, toIx := toIndex(from), toIndex(to)
		nums[fromIx], nums[toIx] = nums[toIx], nums[fromIx]
		from++
		to--
	}
}

func knotHash64(input []byte) []byte {
	sparse := make([]byte, 256)
	for i := 0; i < len(sparse); i++ {
		sparse[i] = byte(i)
	}

	lengths := make([]byte, len(input)+len(SALT))
	copy(lengths[:len(input)], input)
	copy(lengths[len(input):], SALT)

	curr, skip := 0, 0
	for i := 0; i < 64; i++ {
		for _, ll := range lengths {
			reverse(sparse, curr, int(ll))
			curr = (curr + int(ll) + skip) % len(sparse)
			skip++
		}
	}

	dense := make([]byte, 0, 16)
	off := 0
	for i := 0; i < 16; i++ {
		oct := sparse[off]
		for j := 1; j < 16; j++ {
			oct ^= sparse[off+j]
		}
		dense = append(dense, oct)
		off += 16
	}

	return dense
}

func countBits(hash []byte) int {
	cnt := 0
	for _, b := range hash {
		for b > 0 {
			b &= (b - 1)
			cnt++
		}
	}
	return cnt
}

func countRegions(hashes [][]byte) int {
	var visited [128][128]bool
	var visit func(i, j int)

	isBitSet := func(i, j int) bool {
		num := hashes[i][j/8]
		off := j % 8
		return (num & (1 << (7 - off))) > 0
	}

	visit = func(i, j int) {
		visited[i][j] = true
		for _, step := range STEPS4 {
			i1, j1 := i+step[0], j+step[1]
			if i1 < 0 || j1 < 0 || i1 >= 128 || j1 >= 128 {
				continue
			}
			if isBitSet(i1, j1) && !visited[i1][j1] {
				visit(i1, j1)
			}
		}
	}

	regions := 0

	for i := 0; i < 128; i++ {
		for j := 0; j < 128; j++ {
			if isBitSet(i, j) && !visited[i][j] {
				visit(i, j)
				regions++
			}
		}
	}

	return regions
}

func main() {
	//base := "flqrgnkx"
	base := "stpzcrnm"

	hashes := make([][]byte, 0, 128)

	for i := 0; i < 128; i++ {
		key := base + "-" + strconv.Itoa(i)
		hashes = append(hashes, knotHash64([]byte(key)))
	}

	bitCnt := 0
	for _, hash := range hashes {
		bitCnt += countBits(hash)
	}

	printf("total bit count: %d", bitCnt)

	printf("number of regions: %d", countRegions(hashes))
}
