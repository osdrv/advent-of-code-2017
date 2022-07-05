package main

import (
	"fmt"
	"os"
)

type Op uint8

const (
	_ Op = iota
	Inc
	Dec
)

type Comp uint8

const (
	_ Comp = iota
	Lt
	Lte
	Eq
	Neq
	Gte
	Gt
)

type Cond struct {
	reg  string
	comp Comp
	arg  int
}

func (c *Cond) String() string {
	var compStr string
	switch c.comp {
	case Lt:
		compStr = "<"
	case Lte:
		compStr = "<="
	case Eq:
		compStr = "=="
	case Neq:
		compStr = "!="
	case Gte:
		compStr = ">="
	case Gt:
		compStr = ">"
	default:
		panic("should not happen")
	}
	return fmt.Sprintf("%s %s %d", c.reg, compStr, c.arg)
}

type Stmt struct {
	cond *Cond
	reg  string
	op   Op
	arg  int
}

func (s *Stmt) String() string {
	var opStr string
	switch s.op {
	case Inc:
		opStr = "inc"
	case Dec:
		opStr = "dec"
	default:
		panic("should not happen")
	}
	return fmt.Sprintf("if (%s) { %s %s %d }", s.cond, s.reg, opStr, s.arg)
}

func skipWhitespace(s string, p int) int {
	for p < len(s) {
		if s[p] == ' ' || s[p] == '\t' {
			p++
			continue
		}
		break
	}
	return p
}

func isAlpha(b byte) bool {
	return b >= 'a' && b <= 'z'
}

func readIdentifier(s string, p int) (string, int) {
	p = skipWhitespace(s, p)
	start := p
	ptr := p
	for ptr < len(s) && isAlpha(s[ptr]) {
		ptr++
	}
	return s[start:ptr], ptr
}

func readOp(s string, p int) (Op, int) {
	p = skipWhitespace(s, p)
	if s[p:p+3] == "inc" {
		return Inc, p + 3
	} else if s[p:p+3] == "dec" {
		return Dec, p + 3
	}
	panic(fmt.Sprintf("unexpected Op: %q", s[p:]))
}

func isNum(b byte) bool {
	return b >= '0' && b <= '9'
}

func isSign(b byte) bool {
	return b == '-' || b == '+'
}

func readNumber(s string, p int) (int, int) {
	p = skipWhitespace(s, p)
	start := p
	ptr := p
	for ptr < len(s) && (isNum(s[ptr]) || isSign(s[ptr])) {
		ptr++
	}
	return parseInt(s[start:ptr]), ptr
}

func readStatic(exp string, s string, p int) int {
	p = skipWhitespace(s, p)
	expix := 0
	for expix < len(exp) && p < len(s) && s[p] == exp[expix] {
		p++
		expix++
	}
	if expix != len(exp) {
		panic(fmt.Sprintf("unexpected static lexeme: %q (want: %q)", s[p:], exp))
	}
	return p
}

func readComp(s string, p int) (Comp, int) {
	p = skipWhitespace(s, p)
	if s[p] == '<' {
		if s[p+1] == '=' {
			return Lte, p + 2
		}
		return Lt, p + 1
	} else if s[p] == '=' {
		if s[p+1] == '=' {
			return Eq, p + 2
		}
	} else if s[p] == '>' {
		if s[p+1] == '=' {
			return Gte, p + 2
		}
		return Gt, p + 1
	} else if s[p] == '!' {
		if s[p+1] == '=' {
			return Neq, p + 2
		}
	}
	panic(fmt.Sprintf("unable to parse comp: %q", s[p:]))
}

func parseStmt(s string) *Stmt {
	ptr := 0
	var stmtReg string
	stmtReg, ptr = readIdentifier(s, ptr)
	var op Op
	op, ptr = readOp(s, ptr)
	var stmtArg int
	stmtArg, ptr = readNumber(s, ptr)
	ptr = readStatic("if", s, ptr)
	var condReg string
	condReg, ptr = readIdentifier(s, ptr)
	var comp Comp
	comp, ptr = readComp(s, ptr)
	var condArg int
	condArg, ptr = readNumber(s, ptr)

	return &Stmt{
		reg: stmtReg,
		op:  op,
		arg: stmtArg,
		cond: &Cond{
			reg:  condReg,
			comp: comp,
			arg:  condArg,
		},
	}
}

func condTrue(cond *Cond, regs map[string]int) bool {
	regval := regs[cond.reg]
	switch cond.comp {
	case Lt:
		return regval < cond.arg
	case Lte:
		return regval <= cond.arg
	case Eq:
		return regval == cond.arg
	case Neq:
		return regval != cond.arg
	case Gte:
		return regval >= cond.arg
	case Gt:
		return regval > cond.arg
	default:
		panic("should not happen")
	}
}

func applyStmt(stmt *Stmt, regs map[string]int) {
	switch stmt.op {
	case Inc:
		regs[stmt.reg] += stmt.arg
	case Dec:
		regs[stmt.reg] -= stmt.arg
	default:
		panic("should not happen")
	}
}

func interpret(stmts []*Stmt, regs map[string]int) map[string]int {
	absmaxreg := ""
	absmaxval := -ALOT
	for _, stmt := range stmts {
		if condTrue(stmt.cond, regs) {
			applyStmt(stmt, regs)
			maxreg, maxval := findMaxRegVal(regs)
			if absmaxreg == "" || maxval > absmaxval {
				absmaxreg = maxreg
				absmaxval = maxval
			}
		}
	}

	printf("max reg: %s (%d)", absmaxreg, absmaxval)

	return regs
}

func findMaxRegVal(regs map[string]int) (string, int) {
	maxval := -ALOT
	maxreg := ""
	for reg, val := range regs {
		if val > maxval {
			maxval = val
			maxreg = reg
		}
	}

	return maxreg, maxval
}

func main() {
	f, err := os.Open("INPUT-TST")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	printf("===== part 1 =====")

	stmts := make([]*Stmt, 0, len(lines))
	for _, line := range lines {
		stmt := parseStmt(line)
		stmts = append(stmts, stmt)
		printf(stmt.String())
	}

	regs := interpret(stmts, make(map[string]int))
	printf("regs: %+v", regs)

	maxreg, maxval := findMaxRegVal(regs)

	printf("maxreg: %s (%d)", maxreg, maxval)
}
