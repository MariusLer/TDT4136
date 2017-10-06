package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type node struct {
	x int
	y int

	g int
	h int
	f int

	status int // 1 = open, 0 = closed

	parent *node
	kids   []node // In golang []*node and []node both create a slice, so they are the same. Slices are references to an array
}

func readBoard(board string) (node, node, [][]node) {
	startnode := node{g: 0, f: 0}
	endnode := node{}
	currentdir, _ := os.Getwd()
	file, err := os.Open(currentdir + "/boards/" + board)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	linecount := 0
	width := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linecount++
		if width == 0 {
			line := scanner.Text()
			for i := range line {
				width = i + 1
			}
		}
	}
	board2 := make([][]node, linecount)
	for i := range board2 {
		board2[i] = make([]node, width)
	}
	scanner = bufio.NewScanner(file)
	linecount = 0
	for scanner.Scan() {
		line := scanner.Text()
		for i, c := range line {
			tempNode := node{x: i, y: linecount}
			if c == '.' {
				tempNode.g = 1
			} else if c == '#' {
				tempNode.g = 9001
			} else if c == 'A' {
				startnode.x = i
				startnode.y = linecount
			} else if c == 'B' {
				endnode.x = i
				endnode.y = linecount
			}
			board2[i][linecount] = tempNode
			fmt.Println(tempNode)
		}
		linecount++
	}
	return startnode, endnode, board2
}

func main() {
	startnode, stopnode, _ := readBoard("board-1-1.txt")
	//fmt.Println(board)
	fmt.Println(startnode)
	fmt.Println(stopnode)
}
