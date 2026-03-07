package main

func (g *Game) find(x int) int {
	root := x
	for g.parent[root] != root {
		root = g.parent[root]
	}

	// now root stores the root. we can point each
	// node along this path to the root for faster
	// lookups in the future!
	for g.parent[x] != root {
		next := g.parent[x]
		g.parent[x] = root
		x = next
	}
	return root
}

func (g *Game) union(x, y int) {
	rootX := g.find(x)
	rootY := g.find(y)
	g.parent[rootY] = rootX
}
