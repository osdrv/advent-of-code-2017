package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Instr struct {
	cmd          string
	arg1, arg2   int
	argR1, argR2 bool
}

func printArg(v int, isRef bool) string {
	if isRef {
		return string([]byte{byte(v) + 'a'})
	}
	return strconv.Itoa(v)
}

func (i *Instr) String() string {
	return fmt.Sprintf("[%s %s %s]", i.cmd, printArg(i.arg1, i.argR1), printArg(i.arg2, i.argR2))
}

func looksLikeRef(s string) bool {
	return len(s) == 1 && (s[0] >= 'a' && s[0] <= 'z')
}

func parseArgv(s string) (int, bool) {
	if looksLikeRef(s) {
		return int(s[0] - 'a'), true
	}
	return parseInt(s), false
}

func parseInstr(s string) *Instr {
	chs := strings.Split(s, " ")
	arg1, argR1 := parseArgv(chs[1])
	arg2, argR2 := parseArgv(chs[2])
	return &Instr{
		cmd:   chs[0],
		arg1:  arg1,
		argR1: argR1,
		arg2:  arg2,
		argR2: argR2,
	}
}

func getVal(regs *[8]int, v int, isRef bool) int {
	if !isRef {
		return v
	}
	return regs[v]
}

func setVal(regs *[8]int, v int, isRef bool, newV int) {
	if !isRef {
		panic("should be a ref!")
	}
	regs[v] = newV
}

func interpret(ii []*Instr, regs *[8]int) map[string]int {
	pc := 0
	cnts := make(map[string]int)
INTERPRET:
	for pc < len(ii) {
		printf("pc: %d, regs: %+v", pc, regs)
		i := ii[pc]
		cnts[i.cmd]++
		switch i.cmd {
		case "set":
			setVal(regs, i.arg1, i.argR1, getVal(regs, i.arg2, i.argR2))
		case "sub":
			setVal(regs, i.arg1, i.argR1, getVal(regs, i.arg1, i.argR1)-getVal(regs, i.arg2, i.argR2))
		case "mul":
			setVal(regs, i.arg1, i.argR1, getVal(regs, i.arg1, i.argR1)*getVal(regs, i.arg2, i.argR2))
		case "jnz":
			if getVal(regs, i.arg1, i.argR1) != 0 {
				pc += getVal(regs, i.arg2, i.argR2)
				continue INTERPRET
			}
		}
		pc++
	}

	return cnts
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	ii := make([]*Instr, 0, len(lines))

	for _, s := range lines {
		ii = append(ii, parseInstr(s))
	}

	printf("instrs: %+v", ii)

	var regs [8]int

	cnts := interpret(ii, &regs)

	printf("op cnts: %+v", cnts)

	printf("regs: %+v", regs)

	program()
}

func program() {
	a := 1
	b := 93
	c := b
	b *= 100
	b += 100_000
	c = b
	c += 17_000
	var d, e, f, g, h int

	for {
		f = 1
		d = 2
		e = 2
		for d := 2; d*d < b; d++ {
			if b%d == 0 {
				f = 0
				break
			}
		}
		if f == 0 {
			h++
		}
		g = b - c
		b += 17
		//TODO
		if g == 0 {
			break
		}
	}

	printf("regs: %+v", []int{a, b, c, d, e, f, g, h})
}
