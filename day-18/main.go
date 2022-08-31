package main

import (
	"os"
	"strconv"
)

type InstrTyp uint8

const (
	_ InstrTyp = iota
	SND
	SET
	ADD
	MUL
	MOD
	RCV
	JGZ
)

type Arg interface{}

type Val int

func (v Val) String() string {
	return strconv.Itoa(int(v))
}

func (v Val) Int() int {
	return int(v)
}

type Ref byte

func (r Ref) String() string {
	return string([]byte{byte(r)})
}

func (r Ref) Byte() byte {
	return byte(r)
}

type Instr struct {
	typ  InstrTyp
	x, y Arg
}

func looksLikeChar(s string) bool {
	return s[0] >= 'a' && s[0] <= 'z'
}

func parseArg(s string) Arg {
	if looksLikeChar(s) {
		return Ref(s[0])
	}
	return Val(parseInt(s))
}

func getVal(regs map[byte]int, arg Arg) int {
	switch v := arg.(type) {
	case Val:
		return int(v)
	case Ref:
		return regs[byte(v)]
	default:
		panic("should not happen")
	}
}

func getRef(regs map[byte]int, arg Arg) byte {
	switch v := arg.(type) {
	case Ref:
		return byte(v)
	default:
		panic("should not happen either")
	}
}

func interpret(instrs []Instr, regs map[byte]int) int {
	freq := 0

	pc := 0

INSTR:
	for pc < len(instrs) {
		instr := instrs[pc]
		switch instr.typ {
		case SND:
			freq = getVal(regs, instr.x)
		case SET:
			regs[getRef(regs, instr.x)] = getVal(regs, instr.y)
		case ADD:
			regs[getRef(regs, instr.x)] += getVal(regs, instr.y)
		case MUL:
			regs[getRef(regs, instr.x)] *= getVal(regs, instr.y)
		case MOD:
			regs[getRef(regs, instr.x)] %= getVal(regs, instr.y)
		case RCV:
			v := getVal(regs, instr.x)
			if v != 0 {
				return freq
			}
		case JGZ:
			v := getVal(regs, instr.x)
			if v > 0 {
				pc += getVal(regs, instr.y)
				continue INSTR
			}
		default:
			panic("unknown instruction")
		}
		pc++
	}

	return freq
}

func parseInstr(s string) (Instr, error) {
	printf("parsing %s", s)
	var x, y Arg
	if looksLikeChar(s[4:5]) {
		x = Ref(s[4])
	} else {
		x = Val(parseInt(s[4:5]))
	}
	var typ InstrTyp
	switch s[:3] {
	case "snd":
		typ = SND
	case "set":
		typ = SET
		y = parseArg(s[6:])
	case "add":
		typ = ADD
		y = parseArg(s[6:])
	case "mul":
		typ = MUL
		y = parseArg(s[6:])
	case "mod":
		typ = MOD
		y = parseArg(s[6:])
	case "rcv":
		typ = RCV
	case "jgz":
		typ = JGZ
		y = parseArg(s[6:])
	}
	return Instr{
		typ: typ,
		x:   x,
		y:   y,
	}, nil
}

func interpretACP(pid int, instrs []Instr, regs map[byte]int, in <-chan int, out chan<- int) {
	pc := 0
	sndCnt := 0
INSTR:
	for pc < len(instrs) {
		instr := instrs[pc]
		switch instr.typ {
		case SND:
			sndCnt++
			printf("program %d sends a value cnt: %d", pid, sndCnt)
			out <- getVal(regs, instr.x)
		case SET:
			regs[getRef(regs, instr.x)] = getVal(regs, instr.y)
		case ADD:
			regs[getRef(regs, instr.x)] += getVal(regs, instr.y)
		case MUL:
			regs[getRef(regs, instr.x)] *= getVal(regs, instr.y)
		case MOD:
			regs[getRef(regs, instr.x)] %= getVal(regs, instr.y)
		case RCV:
			regs[getRef(regs, instr.x)] = <-in
		case JGZ:
			v := getVal(regs, instr.x)
			if v > 0 {
				pc += getVal(regs, instr.y)
				continue INSTR
			}
		default:
			panic("unknown instruction")
		}
		pc++
	}
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	instrs := make([]Instr, 0, len(lines))

	for _, line := range lines {
		instr, err := parseInstr(line)
		noerr(err)
		instrs = append(instrs, instr)
	}

	printf("===== part 1 =====")

	regs := make(map[byte]int)
	freq := interpret(instrs, regs)

	printf("Recovered frequency: %d", freq)

	printf("===== part 2 =====")

	ch0to1 := make(chan int, 1024*1024)
	ch1to0 := make(chan int, 1024*1024)
	defer close(ch0to1)
	defer close(ch1to0)

	done1 := make(chan struct{})
	done2 := make(chan struct{})
	go func() {
		interpretACP(0, instrs, map[byte]int{'p': 0}, ch1to0, ch0to1)
		close(done1)
	}()
	go func() {
		interpretACP(1, instrs, map[byte]int{'p': 1}, ch0to1, ch1to0)
		close(done2)
	}()

	<-done1
	<-done2
}
