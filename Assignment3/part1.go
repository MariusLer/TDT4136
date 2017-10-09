package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

const inf = int(math.MaxInt32 >> 1)

type node struct {
	y int
	x int

	g int
	f int

	block bool

	parent *node
	kids   []node // In golang []*node and []node both create a slice, so they are the same. Slices are references to an array
}

type nodes []*node

func (a nodes) Len() int {
	return len(a)
}

func (a nodes) Less(i, j int) bool {
	return a[i].f < a[j].f
}

func (a nodes) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type coordinates struct {
	y int
	x int
}

func readBoard(board string) (coordinates, coordinates, [][]node) {
	var startnode coordinates
	var endnode coordinates
	currentdir, _ := os.Getwd()
	file, err := os.Open(currentdir + "/boards/" + board)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	linecount := 0
	board2 := make([][]node, 1)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tempNodeRow := make([]node, 0)
		line := scanner.Text()
		for i, c := range line {
			tempNode := node{x: i, y: linecount}
			tempNode.g = inf
			tempNode.f = inf
			if c == '#' {
				tempNode.block = true
			} else if c == 'A' {
				startnode.x = i
				startnode.y = linecount
			} else if c == 'B' {
				endnode.x = i
				endnode.y = linecount
			}
			tempNodeRow = append(tempNodeRow, tempNode)
		}
		if linecount == 0 {
			board2[0] = tempNodeRow
		} else {
			board2 = append(board2, tempNodeRow)
		}
		linecount++
	}
	return startnode, endnode, board2
}

func calculateF(current *node, stop *node) int {
	return current.g + int(math.Abs(float64(current.x-stop.x))+math.Abs(float64(current.y-stop.y)))
}

func nodeInSet(nod node, set []*node) bool {
	for _, elem := range set {
		if nod.x == elem.x && nod.y == elem.y {
			return true
		}
	}
	return false
}

func aStarSolve(startnode coordinates, stopnode coordinates, board [][]node, height int, width int) {
	open := make([]*node, 1)
	//closed := make([]*node, 0)
	open[0] = &board[startnode.y][startnode.x]
	open[0].g = 0 // startnode cost zero
	open[0].f = calculateF(&board[startnode.y][startnode.x], &board[stopnode.y][stopnode.x])

}

func main() {
	startnode, stopnode, board := readBoard("board-1-1.txt")
	/*
		for _, row := range board {
			for _, node := range row {
				fmt.Print(node)
			}
			fmt.Println()
		}
	*/
	fmt.Println(startnode)
	fmt.Println(stopnode)
	height := len(board)
	width := len(board[1])
	aStarSolve(startnode, stopnode, board, height, width)
}
