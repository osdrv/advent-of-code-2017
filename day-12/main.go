package main

import (
	"os"
	"strings"
)

func parseEdge(s string) (int, []int) {
	ss := strings.SplitN(s, " <-> ", 2)
	from := parseInt(ss[0])
	to := parseInts(ss[1])
	return from, to
}

func countAncestorsOf(graph map[int][]int, pid int, visited map[int]bool) int {

	var visit func(pid int)
	visit = func(pid int) {
		printf("visiting %d", pid)
		if _, ok := visited[pid]; ok {
			return
		}
		visited[pid] = true
		for _, ch := range graph[pid] {
			if _, ok := visited[ch]; !ok {
				visit(ch)
			}
		}
	}

	visit(pid)

	printf("visited: %+v", visited)

	return len(visited)
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	printf("===== part 1 =====")

	graph := make(map[int][]int)

	for _, line := range lines {
		from, to := parseEdge(line)
		if _, ok := graph[from]; ok {
			fatalf("duplicate definition of a vertex")
		}
		graph[from] = to
	}

	res := countAncestorsOf(graph, 0, make(map[int]bool))
	printf("ancestors of 0: %d", res)

	visited := make(map[int]bool)
	groups := 0
	for vert := range graph {
		if _, ok := visited[vert]; ok {
			continue
		}
		if size := countAncestorsOf(graph, vert, visited); size > 0 {
			groups++
		}
	}
	printf("groups: %d", groups)
}
