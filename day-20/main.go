package main

import (
	"fmt"
	"os"
)

type Particle struct {
	id            int
	pos, vel, acl Point3
}

func NewParticle(id int, pos, vel, acl Point3) *Particle {
	return &Particle{
		id:  id,
		pos: pos,
		vel: vel,
		acl: acl,
	}
}

func (p *Particle) String() string {
	return fmt.Sprintf("point{id: %d, pos=%s, vel=%s, acl=%s}", p.id, &p.pos, &p.vel, &p.acl)
}

func (p *Particle) Move() {
	p.vel.x += p.acl.x
	p.vel.y += p.acl.y
	p.vel.z += p.acl.z
	p.pos.x += p.vel.x
	p.pos.y += p.vel.y
	p.pos.z += p.vel.z
}

func (p *Particle) Dist(from Point3) int {
	return abs(p.pos.x-from.x) + abs(p.pos.y-from.y) + abs(p.pos.z-from.z)
}

func readStr(s string, ptr int, read string) int {
	if s[ptr:ptr+len(read)] != read {
		panic("whuuut")
	}
	return ptr + len(read)
}

func isNumber(b byte) bool {
	return b == '-' || b == '+' || (b >= '0' && b <= '9')
}

func readInt(s string, ptr int) (int, int) {
	from := ptr
	until := from
	for until < len(s) && isNumber(s[until]) {
		until++
	}
	return parseInt(s[from:until]), until
}

func parseParticle(id int, s string) *Particle {
	ptr := 0
	var pos, vel, acl Point3

	ptr = readStr(s, ptr, "p=<")
	pos.x, ptr = readInt(s, ptr)
	ptr = readStr(s, ptr, ",")
	pos.y, ptr = readInt(s, ptr)
	ptr = readStr(s, ptr, ",")
	pos.z, ptr = readInt(s, ptr)

	ptr = readStr(s, ptr, ">, v=<")
	vel.x, ptr = readInt(s, ptr)
	ptr = readStr(s, ptr, ",")
	vel.y, ptr = readInt(s, ptr)
	ptr = readStr(s, ptr, ",")
	vel.z, ptr = readInt(s, ptr)

	ptr = readStr(s, ptr, ">, a=<")
	acl.x, ptr = readInt(s, ptr)
	ptr = readStr(s, ptr, ",")
	acl.y, ptr = readInt(s, ptr)
	ptr = readStr(s, ptr, ",")
	acl.z, ptr = readInt(s, ptr)
	ptr = readStr(s, ptr, ">")

	return NewParticle(id, pos, vel, acl)
}

func part1(lines []string) {
	printf("===== part 1 =====")

	particles := make([]*Particle, 0, len(lines))
	for id, line := range lines {
		particles = append(particles, parseParticle(id, line))
	}

	base := Point3{0, 0, 0}
	minDist := ALOT
	var nearest *Particle
	lastSwap := -1
	GOOD_ENOUGH := 1_000

	for i := 0; i < 1_000_000_000; i++ {
		for _, p := range particles {
			p.Move()
			if d := p.Dist(base); d < minDist {
				d = minDist
				nearest = p
				lastSwap = i
			}
		}
		if i-lastSwap > GOOD_ENOUGH {
			printf("the nearest point is %s", nearest)
			break
		}
	}
}

func part2(lines []string) {
	printf("===== part 2 =====")

	particles := make([]*Particle, 0, len(lines))
	for id, line := range lines {
		particles = append(particles, parseParticle(id, line))
	}

	GOOD_ENOUGH := 1000

	lastCollide := 0
	ix := 0
	for {
		posMap := make(map[Point3]*Particle)
		for _, p := range particles {
			if _, ok := posMap[p.pos]; ok {
				posMap[p.pos] = nil
			} else {
				posMap[p.pos] = p
			}
		}
		newParticles := make([]*Particle, 0, len(posMap))
		for _, p := range posMap {
			if p != nil {
				newParticles = append(newParticles, p)
			}
		}
		particles = newParticles
		for _, p := range particles {
			p.Move()
		}

		ix++

		if ix-lastCollide > GOOD_ENOUGH {
			printf("%d particles left", len(particles))
			break
		}
	}
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	part1(lines)

	part2(lines)
}
