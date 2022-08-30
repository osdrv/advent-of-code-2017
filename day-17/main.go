package main

import (
	"bytes"
	"strconv"
)

func printBuf(buf []int, curr int) string {
	var b bytes.Buffer
	for ix, v := range buf {
		s := strconv.Itoa(v)
		if b.Len() > 0 {
			b.WriteRune(' ')
		}
		if ix == curr {
			s = "(" + s + ")"
		}
		b.WriteString(s)
	}

	return b.String()
}

func insertWithStep(buf []int, maxVal int, step int) ([]int, int) {
	curr := 0

	for i := 1; i <= maxVal; i++ {
		curr = (curr + step) % len(buf)
		bufcp := make([]int, len(buf)+1)
		copy(bufcp[:curr+1], buf[:curr+1])
		bufcp[curr+1] = i
		copy(bufcp[curr+2:], buf[curr+1:])
		buf = bufcp
		curr++
	}
	return buf, curr
}

func insertNoMemoGetSecond(maxVal int, step int) int {
	curr := 0
	lastInsAt0 := 0
	ll := 1
	for i := 1; i <= maxVal; i++ {
		curr = (curr + step) % ll
		if curr == 0 {
			lastInsAt0 = i
		}
		ll++
		curr++
	}
	return lastInsAt0
}

func main() {
	printf("===== part 1 =====")
	//step := 3
	step := 349

	maxVal := 2017

	buf := make([]int, 0, maxVal+1)
	buf = append(buf, 0)

	res, curr := insertWithStep(buf, maxVal, step)

	println(printBuf(res, curr))

	printf("The item after %d is %d", maxVal, res[(curr+1)%len(res)])

	printf("===== part 2 =====")

	v := insertNoMemoGetSecond(50000000, step)
	printf("The item after 0 is %d", v)
}
