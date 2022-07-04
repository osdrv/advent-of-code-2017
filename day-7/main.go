package main

import (
	"os"
	"strings"
)

func parseTreeNode(s string) (string, int) {
	chunks := strings.SplitN(s, " ", 2)
	return chunks[0], parseInt(chunks[1][1 : len(chunks[1])-1])
}

func parseTree(lines []string) (map[string]string, map[string]int) {
	tree := make(map[string]string)
	weights := make(map[string]int)
	for _, line := range lines {
		chunks := strings.SplitN(line, " -> ", 2)
		name, weight := parseTreeNode(chunks[0])
		weights[name] = weight
		if _, ok := tree[name]; !ok {
			tree[name] = name
		}
		if len(chunks) > 1 {
			children := words(chunks[1])
			for _, node := range children {
				tree[node] = name
			}
		}
	}

	return tree, weights
}

func findDisbalance(root string, tree map[string]string, weights map[string]int) {
	adj := make(map[string][]string)
	for node, parent := range tree {
		if node != parent {
			if _, ok := adj[parent]; !ok {
				adj[parent] = make([]string, 0, 1)
			}
			adj[parent] = append(adj[parent], node)
		}
	}

	// compute cumulative weights
	cumwghts := make(map[string]int)
	var computeWeight func(node string) int
	computeWeight = func(node string) int {
		w := weights[node]
		for _, ch := range adj[node] {
			w += computeWeight(ch)
		}
		cumwghts[node] = w
		return w
	}
	computeWeight(root)

	var findDisbalance func(node string)
	findDisbalance = func(node string) {
		if len(adj[node]) == 0 {
			// no children
			return
		}
		chwght := 0
		expwght := 0
		wghts := make(map[int]int)
		for _, ch := range adj[node] {
			chwght += cumwghts[ch]
			wghts[cumwghts[ch]]++
			if wghts[cumwghts[ch]] > 1 || expwght == 0 {
				expwght = cumwghts[ch]
			}
		}
		for _, ch := range adj[node] {
			if cumwghts[ch] != expwght {
				printf("node %s seems to be disbalanced (%d vs %d(expected))", ch, cumwghts[ch], expwght)
				printf("node %s own weight is: %d", ch, weights[ch])
				printf("expected weight is: %d", weights[ch]-(cumwghts[ch]-expwght))
				findDisbalance(ch)
			}
		}
	}

	findDisbalance(root)
}

func main() {
	f, err := os.Open("INPUT")
	noerr(err)
	defer f.Close()

	lines := readLines(f)

	tree, weights := parseTree(lines)

	printf("===== part 1 =====")

	var root string
	for node, parent := range tree {
		if node == parent {
			printf("The root node is: %s", node)
			root = node
			break
		}
	}

	printf("===== part 2 =====")

	findDisbalance(root, tree, weights)
}
