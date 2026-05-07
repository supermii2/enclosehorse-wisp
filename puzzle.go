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
	BonusType   *string   `json:"bonusType"`
	walls       []int
	bestScore    int
	bestWalls    []int
	isBonus      bool

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

	// Initialize bestScore to negative infinity and bestWalls
	puzzle.bestScore = -1 << 31
	puzzle.bestWalls = []int{}
	
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
				fmt.Print("🞮 ")
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
				case 'U':
					fmt.Print("🦄")
				default:
					portalColors := []string{"🔴", "🟠", "🟡", "🟢", "🔵", "🟣", "🟤", "⚫", "⚪", "🔶"}
					if cell >= '0' && cell <= '9' {
						fmt.Print(portalColors[cell-'0'])
					} else if cell >= 'a' && cell <= 'z' {
						fmt.Printf(" %c", 'ⓐ'+rune(cell-'a'))
					} else {
						fmt.Printf(" %c", cell)
					}
				}
			}
		}
		fmt.Println()
	}

	score := p.calculateBonusScore()
	scoreStr := fmt.Sprintf("%d", score)
	if score == -1<<31 {
		scoreStr = "N/A"
	}
	bestScoreStr := fmt.Sprintf("%d", p.bestScore)
	if p.bestScore == -1<<31 {
		bestScoreStr = "N/A"
	}
	fmt.Printf("Score (enclosed area): %s | Walls used: %d/%d | Best score: %s\n", scoreStr, len(p.walls), p.Budget, bestScoreStr)
}

func (p *Puzzle) placeWall(loc string) {
	if len(p.walls) >= p.Budget {
		fmt.Println("Wall budget reached. Cannot add more walls.")
		return
	}
	// Convert to uppercase for column letters
	loc = strings.ToUpper(loc)
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

	// Calculate score and update bestScore/bestWalls if improved
	score := p.calculateBonusScore()
	if score > p.bestScore || len(p.bestWalls) == 0 {
		p.bestScore = score
		p.bestWalls = append([]int(nil), p.walls...)
		fmt.Printf("New best score: %d\n", p.bestScore)
	}
}

func (p *Puzzle) removeWall(loc string) {
	// Convert to uppercase for column letters
	loc = strings.ToUpper(loc)
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
	// Calculate score and update bestScore/bestWalls if improved
	score := p.calculateBonusScore()
	if score > p.bestScore || len(p.bestWalls) == 0 {
		p.bestScore = score
		p.bestWalls = append([]int(nil), p.walls...)
		fmt.Printf("New best score: %d\n", p.bestScore)
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

// verifyPuzzleSolve returns true if H is enclosed (calculateScore != negative infinity)
func (p *Puzzle) verifyPuzzleSolve() bool {
	return p.calculateScore() != -1<<31
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
		return -1 << 31
	}

	// Build a set of wall indices for fast lookup
	wallSet := make(map[int]struct{}, len(p.walls))
	for _, w := range p.walls {
		wallSet[w] = struct{}{}
	}

	// Build portal map: character -> all positions with that character
	type point struct{ r, c int }
	isPortal := func(cell rune) bool {
		return (cell >= '0' && cell <= '9') || (cell >= 'a' && cell <= 'z')
	}
	portalPositions := make(map[rune][]point)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			cell := p.MapData[i][j]
			if isPortal(cell) {
				portalPositions[cell] = append(portalPositions[cell], point{i, j})
			}
		}
	}

	// BFS to find all cells in the enclosed area
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}
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

		// If this cell is a portal, teleport to all matching portal cells
		currCell := p.MapData[curr.r][curr.c]
		if isPortal(currCell) {
			for _, dest := range portalPositions[currCell] {
				if !visited[dest.r][dest.c] {
					visited[dest.r][dest.c] = true
					queue = append(queue, dest)
				}
			}
		}

		for _, d := range directions {
			nr, nc := curr.r+d[0], curr.c+d[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] {
				idx := nr*cols + nc
				if _, isWall := wallSet[idx]; isWall {
					continue
				}
				cell := p.MapData[nr][nc]
				if cell == '.' || cell == 'C' || cell == 'G' || cell == 'S' || cell == 'U' || isPortal(cell) {
					visited[nr][nc] = true
					queue = append(queue, point{nr, nc})
				}
			}
		}
	}

	if onEdge {
		return -1 << 31 // Not enclosed, negative infinity
	}

	score := 0
	for _, pt := range area {
		cell := p.MapData[pt.r][pt.c]
		score += 1
		if cell == 'C' { score += 3 }
		if cell == 'S' { score -= 5 }
		if cell == 'G' { score += 10 }
	}
	return score
}

// calculateBonusScore dispatches to the appropriate scoring function based on BonusType.
// If isBonus is false or BonusType is nil, it falls back to calculateScore.
func (p *Puzzle) calculateBonusScore() int {
	if !p.isBonus || p.BonusType == nil {
		return p.calculateScore()
	}
	switch *p.BonusType {
	case "costlywalls":
		return p.calculateCostlyWallsScore()
	case "lovebirds":
		return p.calculateLoveBirdsScore()
	default:
		return p.calculateScore()
	}
}

// calculateLoveBirdsScore works like calculateScore but requires both 'H' and 'U' to be
// enclosed in the same connected area. If either is missing or not co-enclosed, returns -1<<31.
func (p *Puzzle) calculateLoveBirdsScore() int {
	rows := len(p.MapData)
	if rows == 0 {
		return 0
	}
	cols := len(p.MapData[0])

	// Find 'H' and 'U'
	type point struct{ r, c int }
	var hPos, uPos point
	foundH, foundU := false, false
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			switch p.MapData[i][j] {
			case 'H':
				hPos, foundH = point{i, j}, true
			case 'U':
				uPos, foundU = point{i, j}, true
			}
		}
	}
	if !foundH || !foundU {
		return -1 << 31
	}

	// Build a set of wall indices for fast lookup
	wallSet := make(map[int]struct{}, len(p.walls))
	for _, w := range p.walls {
		wallSet[w] = struct{}{}
	}

	// Build portal map
	isPortal := func(cell rune) bool {
		return (cell >= '0' && cell <= '9') || (cell >= 'a' && cell <= 'z')
	}
	portalPositions := make(map[rune][]point)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			cell := p.MapData[i][j]
			if isPortal(cell) {
				portalPositions[cell] = append(portalPositions[cell], point{i, j})
			}
		}
	}

	// BFS from 'H'
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}
	queue := []point{hPos}
	visited[hPos.r][hPos.c] = true

	directions := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	onEdge := false
	area := []point{}
	uEnclosed := false

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		area = append(area, curr)

		if curr.r == 0 || curr.r == rows-1 || curr.c == 0 || curr.c == cols-1 {
			onEdge = true
		}
		if curr.r == uPos.r && curr.c == uPos.c {
			uEnclosed = true
		}

		currCell := p.MapData[curr.r][curr.c]
		if isPortal(currCell) {
			for _, dest := range portalPositions[currCell] {
				if !visited[dest.r][dest.c] {
					visited[dest.r][dest.c] = true
					queue = append(queue, dest)
				}
			}
		}

		for _, d := range directions {
			nr, nc := curr.r+d[0], curr.c+d[1]
			if nr >= 0 && nr < rows && nc >= 0 && nc < cols && !visited[nr][nc] {
				idx := nr*cols + nc
				if _, isWall := wallSet[idx]; isWall {
					continue
				}
				cell := p.MapData[nr][nc]
				if cell == '.' || cell == 'C' || cell == 'U' || isPortal(cell) {
					visited[nr][nc] = true
					queue = append(queue, point{nr, nc})
				}
			}
		}
	}

	if onEdge || !uEnclosed {
		return -1 << 31
	}

	score := 0
	for _, pt := range area {
		cell := p.MapData[pt.r][pt.c]
		score += 1
		if cell == 'C' { score += 3 }
		if cell == 'S' { score -= 5 }
		if cell == 'G' { score += 10 }
	}
	return score
}

// calculateCostlyWallsScore applies the base score with each wall incurring a cost of -6.
func (p *Puzzle) calculateCostlyWallsScore() int {
	base := p.calculateScore()
	if base == -1<<31 {
		return base
	}
	return base - 6*len(p.walls)
}

// ReloadBestWalls sets the current walls to the bestWalls found so far
func (p *Puzzle) ReloadBestWalls() {
	p.walls = append([]int(nil), p.bestWalls...)
}