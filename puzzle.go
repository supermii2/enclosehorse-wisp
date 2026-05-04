package main

import (
	"fmt"
	"io"
	"net/http"
	"encoding/json"
	"strings"
	"strconv"
	"unicode"
	"bytes"
)

type Puzzle struct {
	ID          string    `json:"id"`
	MapRaw      string    `json:"map"`
	MapData     [][]rune  `json:"-"`
	Budget      int       `json:"budget"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatorName *string   `json:"creatorName"`
	walls           []int
}

func getPuzzleFromAPI(url string) (*Puzzle, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil { return nil, err }

	req.Header.Set("Referer", "https://enclose.horse/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil { return nil, err }

	var puzzle Puzzle
	err = json.Unmarshal(body, &puzzle)
	if err != nil { return nil, err }

	lines := strings.Split(puzzle.MapRaw, "\n")
	for _, line := range lines {
		puzzle.MapData = append(puzzle.MapData, []rune(line))
	}
	return &puzzle, nil
}

func (p *Puzzle) RenderMap() {
	if len(p.MapData) == 0 {
		return
	}
	fmt.Print("  ")
	for i := 0; i < len(p.MapData[0]); i++ {
		fmt.Print(" ")
		fmt.Printf("%c", 'A'+i)
	}
	fmt.Println()
	for i, row := range p.MapData {
		fmt.Printf("%02d", i+1)
		for j, cell := range row {
			idx := i*len(p.MapData[0]) + j
			isWall := false
			for _, w := range p.walls {
				if w == idx {
					isWall = true
					break
				}
			}
			if isWall {
				fmt.Print("🞮")
			} else {
				switch cell {
				case '~':
					fmt.Print("⬛")
				case '.':
					fmt.Print(" .")
				case 'C':
					fmt.Print("🍒")
				case 'H':
					fmt.Print("🐎")
				case 'S':
					fmt.Print("🐝")
				case 'G':
					fmt.Print("🍏")
				default:
					fmt.Printf(" %c", cell)
				}
			}
		}
		fmt.Println()
	}

	score := p.calculateScore()
	fmt.Printf("Score (enclosed area): %d\n", score)
}

func (p *Puzzle) placeWall(loc string) {
	if len(p.walls) >= p.Budget {
		fmt.Println("Wall budget reached. Cannot add more walls.")
		return
	}
	// Split the string into alpha and numeric parts
	i := 0
	for ; i < len(loc); i++ {
		if unicode.IsDigit(rune(loc[i])) {
			break
		}
	}
	if i == 0 || i == len(loc) {
		fmt.Println("Invalid location format")
		return
	}
	colStr := loc[:i]
	rowStr := loc[i:]
	col := int(colStr[0] - 'A')
	row, err := strconv.Atoi(rowStr)
	if err != nil {
		fmt.Println("Invalid row number")
		return
	}
	row-- // Convert to 0-based index
	if row < 0 || row >= len(p.MapData) || col < 0 || col >= len(p.MapData[0]) {
		fmt.Println("Location out of bounds")
		return
	}
	if p.MapData[row][col] != '.' {
		fmt.Println("Can only add a wall to an empty square ('.')")
		return
	}
	idx := len(p.MapData[0])*row + col
	fmt.Printf("Placing wall at index: %d (row %d, col %d)\n", idx, row, col)
	p.walls = append(p.walls, idx)
}

func (p *Puzzle) removeWall(loc string) {
	// Split the string into alpha and numeric parts
	i := 0
	for ; i < len(loc); i++ {
		if unicode.IsDigit(
			rune(loc[i])) {
			break
		}
	}
	if i == 0 || i == len(loc) {
		fmt.Println("Invalid location format")
		return
	}
	colStr := loc[:i]
	rowStr := loc[i:]
	col := int(colStr[0] - 'A')
	row, err := strconv.Atoi(rowStr)
	if err != nil {
		fmt.Println("Invalid row number")
		return
	}
	row-- // Convert to 0-based index
	if row < 0 || row >= len(p.MapData) || col < 0 || col >= len(p.MapData[0]) {
		fmt.Println("Location out of bounds")
		return
	}
	idx := len(p.MapData[0])*row + col
	// Remove from walls array if present
	found := false
	for i, w := range p.walls {
		if w == idx {
			// Remove wall from slice
			p.walls = append(p.walls[:i], p.walls[i+1:]...)
			found = true
			break
		}
	}
	if found {
		fmt.Printf("Removed wall at index: %d (row %d, col %d)\n", idx, row, col)
	} else {
		fmt.Println("No wall found at that location to remove.")
	}
}

// submitPuzzle sends a POST request with the current walls and returns (optimalScore, error)
func (p *Puzzle) submitPuzzle() (int, error) {
	if !p.verifyPuzzleSolve() {
		return 0, fmt.Errorf("Cannot submit: puzzle is not solved (verifyPuzzleSolve is false)")
	}
	url := "https://enclose.horse/api/levels/" + p.ID + "/submit"
	body := struct {
		Walls []int `json:"walls"`
	}{
		Walls: p.walls,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Referer", "https://enclose.horse/")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("submit failed: %s", resp.Status)
	}

	// Parse response for .stats.optimalScore
	var respData struct {
		Stats struct {
			OptimalScore int `json:"optimalScore"`
		} `json:"stats"`
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&respData); err != nil {
		return 0, fmt.Errorf("submit succeeded but could not parse optimal score: %v", err)
	}
	return respData.Stats.OptimalScore, nil
}

// verifyPuzzleSolve returns true if 'H' is orthogonally connected to the edge via '.' or 'C'
func (p *Puzzle) verifyPuzzleSolve() bool {
	rows := len(p.MapData)
	if rows == 0 {
		return false
	}
	cols := len(p.MapData[0])

	// Find 'H'
	var startRow, startCol int
	found := false
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if p.MapData[i][j] == 'H' {
				startRow, startCol = i, j
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return false
	}

	// Prepare visited matrix
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	// Build a set of wall indices for fast lookup
	wallSet := make(map[int]struct{}, len(p.walls))
	for _, w := range p.walls {
		wallSet[w] = struct{}{}
	}

	// BFS from H
	type point struct{ r, c int }
	queue := []point{{startRow, startCol}}
	visited[startRow][startCol] = true

	directions := [][2]int{{-1,0},{1,0},{0,-1},{0,1}} // up, down, left, right

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]

		// If we reach the edge, H is not enclosed
		if curr.r == 0 || curr.r == rows-1 || curr.c == 0 || curr.c == cols-1 {
			return false
		}

		// Explore neighbors
		for _, d := range directions {
			nr, nc := curr.r+d[0], curr.c+d[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] {
				idx := nr*cols + nc
				if _, isWall := wallSet[idx]; isWall {
					continue // treat as blocked
				}
				cell := p.MapData[nr][nc]
				if cell == '.' || cell == 'C' {
					visited[nr][nc] = true
					queue = append(queue, point{nr, nc})
				}
			}
		}
	}
	// If we never reach the edge, H is enclosed
	return true
}

// calculateScore returns the total number of '.' in the enclosed area, with 'C' worth 3 points
func (p *Puzzle) calculateScore() int {
	rows := len(p.MapData)
	if rows == 0 {
		return 0
	}
	cols := len(p.MapData[0])

	// Find 'H'
	var startRow, startCol int
	found := false
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if p.MapData[i][j] == 'H' {
				startRow, startCol = i, j
				found = true
				break
			}
		}
		if found {
			break
		}
	}
	if !found {
		return 0
	}

	// Build a set of wall indices for fast lookup
	wallSet := make(map[int]struct{}, len(p.walls))
	for _, w := range p.walls {
		wallSet[w] = struct{}{}
	}

	// BFS to find all cells in the enclosed area
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}
	type point struct{ r, c int }
	queue := []point{{startRow, startCol}}
	visited[startRow][startCol] = true

	directions := [][2]int{{-1,0},{1,0},{0,-1},{0,1}}

	onEdge := false
	area := []point{}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		area = append(area, curr)

		if curr.r == 0 || curr.r == rows-1 || curr.c == 0 || curr.c == cols-1 {
			onEdge = true
		}

		for _, d := range directions {
			nr, nc := curr.r+d[0], curr.c+d[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] {
				idx := nr*cols + nc
				if _, isWall := wallSet[idx]; isWall {
					continue
				}
				cell := p.MapData[nr][nc]
				if cell == '.' || cell == 'C' {
					visited[nr][nc] = true
					queue = append(queue, point{nr, nc})
				}
			}
		}
	}

	if onEdge {
		return 0 // Not enclosed, no score
	}

	score := 0
	for _, pt := range area {
		cell := p.MapData[pt.r][pt.c]
		if cell == '.' || cell == 'H' {
			score += 1
		} else if cell == 'C' {
			score += 3
		} else if cell == 'S' {
			score -= 5
		} else if cell == 'G' {
			score += 10
	}
	return score
}