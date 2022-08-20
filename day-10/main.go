package main

import (
	"bytes"
	"encoding/hex"
	"os"
	"strconv"
)

func makeSequence(n int) []byte {
	list := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		list = append(list, byte(i))
	}
	return list
}

func printList(nums []byte, curr, length int) string {
	var b bytes.Buffer
	from, to := curr, (curr+length-1)%len(nums)
	for ix, num := range nums {
		s := strconv.Itoa(int(num))
		if ix == curr {
			s = "[" + s + "]"
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		if ix == from {
			b.WriteByte('(')
		}
		b.WriteString(s)
		if ix == to {
			b.WriteByte(')')
		}
	}
	return b.String()
}

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

func knotHash(input []byte) []byte {
	hash := makeSequence(256)
	curr, skip := 0, 0
	for _, num := range input {
		println(">>> " + printList(hash, curr, int(num)))
		reverse(hash, curr, int(num))
		curr = (curr + int(num) + skip) % len(hash)
		skip++
	}

	return hash
}

var (
	SALT = []byte{17, 31, 73, 47, 23}
)

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

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	println("===== part 1 =====")

	for _, line := range lines {
		nums := parseInts(line)
		input := make([]byte, 0, len(nums))
		for _, num := range nums {
			input = append(input, byte(num))
		}
		printf("Lengths: %+v", nums)
		hash := knotHash(input)

		printf("n0 * n1 = %d", int(hash[0])*int(hash[1]))
	}

	println("===== part2 =====")

	for _, line := range lines {
		input := parseInput(line)
		hash := knotHash64(input)
		println(hashToStr(hash))
	}
}

func hashToStr(hash []byte) string {
	return hex.EncodeToString(hash)
}

func parseInput(s string) []byte {
	return []byte(s)
}
