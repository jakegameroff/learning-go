package main

import "fmt"

type Game struct {
	board  [totalNodes]Node
	parent [totalNodes]int
}

type Node struct {
	Index    int
	Color    string // "red" | "blue" | ""
	NodeType string
}

type Pair struct{ R, C int }

const size = 11
const totalNodes = size*size + 4
const r1 = size * size
const r2 = size*size + 1
const b1 = size*size + 2
const b2 = size*size + 3

func initGame() Game {
	var g Game
	for i := range g.board {
		nodeType := "board_node"
		var color string
		if i >= r1 {
			nodeType = "win_node"
			if i == r1 || i == r2 {
				color = "red"
			} else {
				color = "blue"
			}
		}
		g.board[i] = Node{
			Index:    i,
			Color:    color,
			NodeType: nodeType,
		}
		g.parent[i] = i
	}
	return g
}

func isOnBoard(r, c int) bool {
	if r < 0 || r > size-1 || c < 0 || c > size-1 {
		return false
	}
	return true
}

func getIndex(r, c int) int { return size*r + c }

func (g *Game) getNeighborIndices(index int) []int {
	var candidates []int
	row := index / size
	col := index % size

	switch row {
	case 0:
		candidates = append(candidates, r1)
	case size - 1:
		candidates = append(candidates, r2)
	}

	switch col {
	case 0:
		candidates = append(candidates, b1)
	case size - 1:
		candidates = append(candidates, b2)
	}

	offsets := [6]Pair{
		{0, -1}, // left
		{0, 1},  // right
		{-1, 0}, // down left
		{-1, 1}, // down right
		{1, -1}, // up left
		{1, 0},  // up right
	}

	for _, o := range offsets {
		newRow := row + o.R
		newCol := col + o.C

		if isOnBoard(newRow, newCol) {
			index := getIndex(newRow, newCol)
			candidates = append(candidates, index)
		}
	}
	fmt.Printf("%+v\n", candidates)
	return candidates
}

func (g *Game) getMonochromaticNeighbors(node Node) []int {
	var monochromaticNeighbors []int

	index := node.Index
	color := node.Color

	candidateIndices := g.getNeighborIndices(index)
	for _, index := range candidateIndices {
		node := g.board[index]
		if node.Color == color {
			monochromaticNeighbors = append(monochromaticNeighbors, node.Index)
		}
	}
	return monochromaticNeighbors
}

func (g *Game) isWinningMove(node Node) bool {
	monochromaticNeighbors := g.getMonochromaticNeighbors(node)
	for _, neighbor := range monochromaticNeighbors {
		g.union(node.Index, neighbor)
	}
	switch node.Color {
	case "red":
		return g.find(r1) == g.find(r2)
	case "blue":
		return g.find(b1) == g.find(b2)
	}
	return false
}

func (g *Game) move(node Node) {
	g.board[node.Index] = node
}
