package main

import (
	"bytes"
	"os"
)

type Dir [2]int

var (
	Down  Dir = [2]int{0, 1}
	Left      = [2]int{-1, 0}
	Right     = [2]int{1, 0}
	Up        = [2]int{0, -1}
)

type Packet struct {
	dir Dir
	pos Point2
}

func NewPacket(pos Point2, dir Dir) *Packet {
	return &Packet{
		pos: pos, dir: dir,
	}
}

func turnLeft(d Dir) Dir {
	switch d {
	case Up:
		return Left
	case Right:
		return Up
	case Down:
		return Right
	case Left:
		return Down
	default:
		panic("whut")
	}
}

func turnRight(d Dir) Dir {
	switch d {
	case Up:
		return Right
	case Right:
		return Down
	case Down:
		return Left
	case Left:
		return Up
	default:
		panic("wuuuut")
	}
}

func (p *Packet) NextStep(route [][]byte) byte {
	if v, ok := p.moveIfYouCan(route, p.dir); ok {
		return v
	} else if v, ok := p.moveIfYouCan(route, turnRight(p.dir)); ok {
		return v
	} else if v, ok := p.moveIfYouCan(route, turnLeft(p.dir)); ok {
		return v
	}
	return EOR
}

func (p *Packet) moveIfYouCan(route [][]byte, dir Dir) (byte, bool) {
	next := Point2{p.pos.x + dir[0], p.pos.y + dir[1]}
	if next.x >= 0 && next.x < len(route[0]) && next.y >= 0 && next.y < len(route) {
		if v := route[next.y][next.x]; v != ' ' {
			p.pos = next
			p.dir = dir
			return v, true
		}
	}
	return 0, false
}

const (
	EOR byte = 255
)

func NewRoute(lines []string) [][]byte {
	route := make([][]byte, 0, len(lines))
	for _, line := range lines {
		route = append(route, []byte(line))
	}
	return route
}

func findStart(route [][]byte) Point2 {
	for j := 0; j < len(route[0]); j++ {
		if route[0][j] == '|' {
			return Point2{x: j, y: 0}
		}
	}
	panic("malformed input")
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	route := NewRoute(lines)
	start := findStart(route)
	packet := Packet{
		pos: start,
		dir: Down,
	}

	var trav bytes.Buffer
	steps := 1
	for {
		v := packet.NextStep(route)
		if v == EOR {
			break
		}
		steps++
		if (v >= 'a' && v <= 'z') || (v >= 'A' && v <= 'Z') {
			trav.WriteByte(v)
		}
	}

	printf("traverse: %s", trav.String())
	printf("steps: %d", steps)
}
