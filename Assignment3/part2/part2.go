package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

const inf = int(math.MaxInt32 >> 1)

type node struct {
	y int
	x int

	g int
	f int

	weight int

	open     bool
	closed   bool
	solution bool // Is this part of best path?

	parent *node
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

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func readBoard(board string) (*node, *node, [][]node) {
	var startnodeX int
	var startnodeY int
	var endnodeX int
	var endnodeY int
	currentdir, _ := os.Getwd()
	file, err := os.Open(currentdir + "/../" + "/boards/" + board)
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
			if c == 'w' {
				tempNode.weight = 100
			} else if c == 'm' {
				tempNode.weight = 50
			} else if c == 'f' {
				tempNode.weight = 10
			} else if c == 'g' {
				tempNode.weight = 5
			} else if c == 'r' {
				tempNode.weight = 1
			} else if c == 'A' {
				startnodeX = i
				startnodeY = linecount
				tempNode.weight = 1
			} else if c == 'B' {
				endnodeX = i
				endnodeY = linecount
				tempNode.weight = 1
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
	return &board2[startnodeY][startnodeX], &board2[endnodeY][endnodeX], board2
}

func calculateF(current *node, stop *node) int {
	return current.g + int(math.Abs(float64(current.x-stop.x))+math.Abs(float64(current.y-stop.y)))
}

func nodeInSet(nod *node, set []*node) bool {
	for _, elem := range set {
		if nod.x == elem.x && nod.y == elem.y {
			return true
		}
	}
	return false
}

func findNeighbors(board [][]node, height int, width int, current *node) []*node {
	neighbors := make([]*node, 0)
	for i := max(0, current.x-1); i < min(width, current.x+2); i++ {
		if current.x != i {
			neighbors = append(neighbors, &board[current.y][i])
		}
	}
	for j := max(0, current.y-1); j < min(height, current.y+2); j++ {
		if current.y != j {
			neighbors = append(neighbors, &board[j][current.x])
		}
	}
	return neighbors
}

func printPath(board [][]node, stopnode *node, startnode *node) {
	var current *node
	current = stopnode
	for current != startnode {
		current.solution = true
		fmt.Println("Y: ", current.y, "X: ", current.x)
		current = current.parent

	}
	fmt.Println("Y: ", current.y, "X: ", current.x)
}

func aStarSolve(startnode *node, stopnode *node, board [][]node, height int, width int) bool {
	height = len(board)
	width = len(board[0])
	open := make([]*node, 1)
	closed := make([]*node, 0)
	open[0] = startnode
	open[0].g = 0 // startnode cost zero
	open[0].f = calculateF(startnode, stopnode)
	open[0].solution = true
	// Algorithm loop
	for len(open) != 0 {
		current := open[0]
		fmt.Println("Current", *current)
		if current.x == stopnode.x && current.y == stopnode.y { // done
			return true
		}
		open = open[1:] // pops open
		current.open = false
		current.closed = true
		closed = append(closed, current)

		neighbors := findNeighbors(board, height, width, current)
		for _, node := range neighbors {
			if nodeInSet(node, closed) {
				tempg := current.g + node.weight
				if tempg < node.g {
					fmt.Println("Found a shorter path to a closed node, expand program") // don't think this will happen
				}
				continue
			}
			if !nodeInSet(node, open) {
				open = append(open, node)
				node.open = true
			}
			tempg := current.g + node.weight
			if tempg >= node.g { // did not find a better path
				continue
			}
			// This path is best
			node.parent = current
			node.g = tempg
			node.f = calculateF(node, stopnode)
		}
		sort.Sort(nodes(open))
	}
	return false
}

func colorRectangle(image *image.RGBA, color color.RGBA, row int, col int) {
	for i := row*22 + 1; i < row*22+22; i++ {
		for j := col*22 + 1; j < col*22+22; j++ {
			image.Set(j, i, color) // img.Set takes x before y
		}
	}
}

func drawBorders(img *image.RGBA) {
	for i := 0; i <= img.Bounds().Max.X; i += 22 {
		for j := 0; j <= img.Bounds().Max.Y; j++ {
			img.Set(i, j, color.RGBA{0, 0, 0, 255})
		}
	}
	for j := 0; j <= img.Bounds().Max.Y; j += 22 {
		for i := 0; i <= img.Bounds().Max.X; i++ {
			img.Set(i, j, color.RGBA{0, 0, 0, 255})
		}
	}
}

func drawSmallRect(img *image.RGBA, color color.RGBA, row int, col int) {
	for i := row*22 + 1; i < row*22+5; i++ {
		for j := col*22 + 1; j < col*22+5; j++ {
			img.Set(j, i, color)
		}
	}
}

func drawDisk(img *image.RGBA, color color.RGBA, row int, col int) {
	var radius = 5.2
	var centerX = row*22 + 11
	var centerY = col*22 + 11
	for i := centerX - 5; i <= centerX+5; i++ {
		for j := centerY - 5; j <= centerY+5; j++ {
			currentRadius := math.Sqrt(math.Pow(float64(centerX-i), 2) + math.Pow(float64(centerY-j), 2))
			if currentRadius <= float64(radius) {
				img.Set(j, i, color)
			}
		}
	}
}

func drawImage(board [][]node, startnode *node, stopnode *node, filename string) {
	// Creating squares of size 20x20 with borders
	height := len(board)
	width := len(board[0])
	img := image.NewRGBA(image.Rect(0, 0, width*22+1, height*22+1))

	drawBorders(img)
	for i := range board {
		for j := range board[i] {
			if &board[i][j] == startnode {
				colorRectangle(img, color.RGBA{255, 0, 0, 255}, i, j)
				drawDisk(img, color.RGBA{0, 255, 255, 255}, i, j)
				continue
			}
			if &board[i][j] == stopnode {
				colorRectangle(img, color.RGBA{0, 255, 0, 255}, i, j)
				drawDisk(img, color.RGBA{0, 255, 255, 255}, i, j)
				continue
			}
			switch board[i][j].weight {
			case 1:
				colorRectangle(img, color.RGBA{139, 69, 9, 255}, i, j)
			case 5:
				colorRectangle(img, color.RGBA{152, 251, 152, 255}, i, j)
			case 10:
				colorRectangle(img, color.RGBA{0, 100, 0, 255}, i, j)
			case 50:
				colorRectangle(img, color.RGBA{205, 200, 177, 255}, i, j)
			case 100:
				colorRectangle(img, color.RGBA{65, 105, 225, 255}, i, j)
			}
			if board[i][j].solution {
				drawDisk(img, color.RGBA{0, 255, 255, 255}, i, j)
			}
			if board[i][j].closed {
				drawSmallRect(img, color.RGBA{255, 0, 255, 255}, i, j)
			}
			if board[i][j].open {
				drawSmallRect(img, color.RGBA{255, 255, 0, 255}, i, j)
			}
		}
	}

	file, _ := os.Create(filename)
	defer file.Close()
	png.Encode(file, img)
}

func main() {
	var file = os.Args[1]
	startnode, stopnode, board := readBoard(file)

	height := len(board)
	width := len(board[1])
	aStarSolve(startnode, stopnode, board, height, width)
	printPath(board, stopnode, startnode)
	imagename := strings.Replace(string(file), "txt", "png", -1)
	drawImage(board, startnode, stopnode, imagename)
}
