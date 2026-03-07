package main

func (g *Game) find(x int) int {
	root := x
	for g.Parent[root] != root {
		root = g.Parent[root]
	}

	// now root stores the root. we can point each
	// node along this path to the root for faster
	// lookups in the future!
	for g.Parent[x] != root {
		next := g.Parent[x]
		g.Parent[x] = root
		x = next
	}
	return root
}

func (g *Game) union(x, y int) {
	rootX := g.find(x)
	rootY := g.find(y)
	g.Parent[rootY] = rootX
}
