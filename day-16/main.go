package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"lukechampine.com/uint128"
)

var (
	MSK = uint64(0xFFFFFFFFFFFFFFFF)
)

func printDance(dance uint128.Uint128) string {
	var b bytes.Buffer
	for ix := 0; ix < 16; ix++ {
		ch := dance.And(uint128.From64(uint64(0xFF)).Lsh(uint((15 - ix) * 8))).Rsh(uint((15 - ix) * 8)).Lo
		b.WriteByte(byte(ch))
	}
	return b.String()
}

type MoveType uint8

const (
	_ MoveType = iota
	MoveSpin
	MoveExchg
	MovePartn
)

type Move struct {
	typ    MoveType
	v1, v2 int
}

func danceAll(dance uint128.Uint128, moves []Move) uint128.Uint128 {
	for _, move := range moves {
		switch move.typ {
		case MoveSpin:
			dance = dance.RotateRight(move.v1 * 8)
		case MoveExchg:
			p1, p2 := move.v1, move.v2
			v1 := (dance.Rsh(uint((15 - p1) * 8)).And64(uint64(0xFF)).Lo) & 0xFF
			v2 := (dance.Rsh(uint((15 - p2) * 8)).And64(uint64(0xFF)).Lo) & 0xFF
			dst1, dst2 := &(dance.Hi), &(dance.Hi)
			if p1 > 7 {
				dst1 = &(dance.Lo)
				p1 -= 8
			}
			if p2 > 7 {
				dst2 = &(dance.Lo)
				p2 -= 8
			}
			msk1 := MSK & (^((0xFF) << ((7 - p1) * 8)))
			msk2 := MSK & (^((0xFF) << ((7 - p2) * 8)))

			*dst1 = (*dst1 & msk1) | (v2 << ((7 - p1) * 8))
			*dst2 = (*dst2 & msk2) | (v1 << ((7 - p2) * 8))
		case MovePartn:
			v1, v2 := uint64(move.v1), uint64(move.v2)
			p1 := 0
			dst1 := &(dance.Hi)
			off1 := 0
			for {
				if p1 > 7 {
					dst1 = &(dance.Lo)
					off1 = 8
				}
				if ((*dst1)>>((7-(p1-off1))*8))&0xFF == v1 {
					break
				}
				p1++
			}
			p2 := 0
			dst2 := &(dance.Hi)
			off2 := 0
			for {
				if p2 > 7 {
					dst2 = &(dance.Lo)
					off2 = 8
				}
				if ((*dst2)>>((7-(p2-off2))*8))&0xFF == v2 {
					break
				}
				p2++
			}
			msk1 := MSK & (^((0xFF) << ((7 - (p1 - off1)) * 8)))
			msk2 := MSK & (^((0xFF) << ((7 - (p2 - off2)) * 8)))

			*dst1 = (*dst1 & msk1) | (v2 << ((7 - (p1 - off1)) * 8))
			*dst2 = (*dst2 & msk2) | (v1 << ((7 - (p2 - off2)) * 8))
		}
	}

	return dance
}

func parseMoves(ss []string) []Move {
	moves := make([]Move, 0, len(ss))
	for _, s := range ss {
		move := Move{}
		switch s[0] {
		case 's':
			rot := parseInt(s[1:])
			move.typ = MoveSpin
			move.v1 = rot
		case 'x':
			chs := strings.SplitN(s[1:], "/", 2)
			assert(len(chs) == 2, "command `x` expects 2 args")
			v1, v2 := parseInt(chs[0]), parseInt(chs[1])
			move.typ = MoveExchg
			move.v1 = v1
			move.v2 = v2
		case 'p':
			v1, v2 := int(s[1]), int(s[3])
			move.typ = MovePartn
			move.v1 = v1
			move.v2 = v2
		default:
			panic(fmt.Sprintf("Unrecognized command: %q", s))
		}
		moves = append(moves, move)
	}

	return moves
}

func newDance() uint128.Uint128 {
	var dance uint128.Uint128

	for ch := 'a'; ch <= 'p'; ch++ {
		off := int(ch - 'a')
		dance = dance.Or(uint128.From64(uint64(ch)).Lsh(uint((15 - off) * 8)))
	}

	return dance
}

func part1(moves []Move) {
	println("===== part 1 =====")

	dance := newDance()

	dance = danceAll(dance, moves)

	printf("dance state: %q", printDance(dance))
}

func part2(moves []Move) {
	println("===== part 2 =====")

	dance := newDance()

	memo := make(map[uint128.Uint128]struct{})
	idx := make(map[int]uint128.Uint128)
	memo[dance] = struct{}{}
	idx[0] = dance

	// find the cycle size
	for i := 0; i < 100; i++ {
		dance = danceAll(dance, moves)
		printf("dance %d: %s", i, printDance(dance))
		if _, ok := memo[dance]; ok {
			printf("detected a cycle after %d interations", i+1)
			break
		}
		memo[dance] = struct{}{}
		idx[i+1] = dance
	}

	cycleSize := len(memo)
	printf("cycle size: %d", cycleSize)

	lookupIx := 1_000_000_000 % cycleSize
	printf("the lookup result is: %s", printDance(idx[lookupIx]))
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)
	moves := parseMoves(words(lines[0]))

	part1(moves)

	part2(moves)
}
