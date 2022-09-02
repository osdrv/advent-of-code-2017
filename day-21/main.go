package main

import (
	"os"
	"strings"
)

type Tile [][]byte

func (t Tile) Dim() int {
	return len(t)
}

func (t Tile) Rot() Tile {
	d := t.Dim()
	for i := 0; i < d/2; i++ {
		for j := i; j < d-i-1; j++ {
			t[j][d-i-1], t[d-i-1][d-j-1], t[d-j-1][i], t[i][j] = t[i][j], t[j][d-i-1], t[d-i-1][d-j-1], t[d-j-1][i]
		}
	}
	return t
}

func (t Tile) Flip() Tile {
	d := t.Dim()
	for i := 0; i < d; i++ {
		for j := 0; j < d/2; j++ {
			t[i][j], t[i][d-j-1] = t[i][d-j-1], t[i][j]
		}
	}
	return t
}

func (t Tile) GetAt(i, j, d int) Tile {
	st := makeByteField(d, d)
	for ii := 0; ii < d; ii++ {
		for jj := 0; jj < d; jj++ {
			st[ii][jj] = t[i+ii][j+jj]
		}
	}
	return st
}

func (t Tile) PutAt(i, j int, st Tile) {
	d := st.Dim()
	for ii := 0; ii < d; ii++ {
		for jj := 0; jj < d; jj++ {
			t[ii+i][jj+j] = st[ii][jj]
		}
	}
}

func (t Tile) Equals(other Tile) bool {
	if t.Dim() != other.Dim() {
		return false
	}
	d := t.Dim()
	for i := 0; i < d; i++ {
		for j := 0; j < d; j++ {
			if t[i][j] != other[i][j] {
				return false
			}
		}
	}
	return true
}

func (t Tile) Count(eq byte) int {
	d := t.Dim()
	cnt := 0
	for i := 0; i < d; i++ {
		for j := 0; j < d; j++ {
			if t[i][j] == eq {
				cnt++
			}
		}
	}
	return cnt
}

func strArrToTile(ss []string) Tile {
	t := make([][]byte, len(ss))
	for i, s := range ss {
		t[i] = []byte(s)
	}
	return Tile(t)
}

func parseRule(s string) (Tile, Tile) {
	ss := strings.SplitN(s, " => ", 2)
	fs := strings.Split(ss[0], "/")
	ts := strings.Split(ss[1], "/")
	return strArrToTile(fs), strArrToTile(ts)
}

func parseRulebook(lines []string) [][2]Tile {
	rulebook := make([][2]Tile, 0, len(lines))
	for _, line := range lines {
		from, to := parseRule(line)
		rulebook = append(rulebook, [2]Tile{from, to})
	}

	return rulebook
}

func findEnhancement(t Tile, rulebook [][2]Tile) Tile {
	for _, rule := range rulebook {
		from := rule[0]
		if from.Equals(t) || from.Rot().Equals(t) || from.Rot().Equals(t) || from.Rot().Equals(t) || from.Flip().Equals(t) || from.Rot().Equals(t) || from.Rot().Equals(t) || from.Rot().Equals(t) {
			return rule[1]
		}
	}
	panic("not found!")
}

func enhance(t Tile, rulebook [][2]Tile) Tile {
	td := t.Dim()
	var nt Tile
	if td%2 == 0 {
		// new dimension
		nd := td / 2 * 3
		nt = Tile(makeByteField(nd, nd))
		tcnt := td / 2
		for i := 0; i < tcnt; i++ {
			for j := 0; j < tcnt; j++ {
				from := t.GetAt(i*2, j*2, 2)
				to := findEnhancement(from, rulebook)
				nt.PutAt(i*3, j*3, to)
			}
		}
	} else if td%3 == 0 {
		nd := td / 3 * 4
		nt = Tile(makeByteField(nd, nd))
		tcnt := td / 3
		for i := 0; i < tcnt; i++ {
			for j := 0; j < tcnt; j++ {
				from := t.GetAt(i*3, j*3, 3)
				to := findEnhancement(from, rulebook)
				nt.PutAt(i*4, j*4, to)
			}
		}
	} else {
		panic("should not happen")
	}
	return nt
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	printf("lines: %+v", lines)

	tile := Tile([][]byte{
		{'.', '#', '.'},
		{'.', '.', '#'},
		{'#', '#', '#'},
	})

	rulebook := parseRulebook(lines)

	for _, rule := range rulebook {
		println("========")
		println(printTile(rule[0]))
		println(printTile(rule[1]))
	}

	ITER := 18

	for i := 0; i < ITER; i++ {
		tile = enhance(tile, rulebook)
		printf("enhancement %d", i+1)
		println(printTile(tile))
	}

	printf("Total number of pixels: %d", tile.Count('#'))
}

func printTile(tile Tile) string {
	return printByteFieldWithSubs(tile, " ", map[byte]string{'.': ".", '#': "#"})
}
