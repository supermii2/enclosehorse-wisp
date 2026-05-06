package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	dateFlag := flag.String("date", "", "Date to fetch puzzle for (YYYY-MM-DD). Defaults to today.")
	bonusFlag := flag.Bool("bonus", false, "Load the bonus level instead of the regular level.")
	flag.BoolVar(bonusFlag, "B", false, "Load the bonus level instead of the regular level.")
	flag.Parse()

	var dateStr string
	if *dateFlag != "" {
		dateStr = *dateFlag
	} else {
		dateStr = time.Now().Format("2006-01-02")
	}
	fmt.Printf("Using date string: %s\n", dateStr)
	url := fmt.Sprintf("https://enclose.horse/api/daily/%s", dateStr)

	puzzle, err := getPuzzleFromAPI(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching puzzle: %v\n", err)
		os.Exit(1)
	}

	if *bonusFlag {
		if puzzle.BonusType == nil {
			fmt.Fprintln(os.Stderr, "Error: no bonus level available for this date.")
			os.Exit(1)
		}
		bonusURL := fmt.Sprintf("https://enclose.horse/api/daily/bonus/%s", puzzle.ID)
		puzzle, err = getPuzzleFromAPI(bonusURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching bonus puzzle: %v\n", err)
			os.Exit(1)
		}
	}

	puzzle.RenderMap()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: ")
		line, err := reader.ReadString('\n')
		if err != nil {
			continue
		}
		line = strings.TrimSpace(line)
		fields := splitFields(line)
		var command, param string
		if len(fields) > 0 {
			command = fields[0]
		}
		if len(fields) > 1 {
			param = fields[1]
		}
		n := len(fields)

		switch command {
		case "add", "a":
			if n < 2 {
				fmt.Println("Invalid input format. Use: add [location] (e.g., add A1)")
				continue
			}
			puzzle.placeWall(param)
			puzzle.RenderMap()
		case "remove", "r":
			if n < 2 {
				fmt.Println("Invalid input format. Use: remove [location] (e.g., remove A1)")
				continue
			}
			puzzle.removeWall(param)
			puzzle.RenderMap()
		case "b":
			puzzle.ReloadBestWalls()
			fmt.Println("Reloaded best walls.")
			puzzle.RenderMap()
		case "submit", "s":
			if n < 1 {
				fmt.Println("Invalid input format. Use: submit")
				continue
			}
			fmt.Println("Submitting puzzle...")
			optimalScore, err := puzzle.submitPuzzle()
			if err != nil {
				fmt.Printf("Submit failed: %v\n", err)
			} else {
				userScore := puzzle.calculateBonusScore()
				fmt.Printf("Submit successful! Your score: %d, Optimal Score: %d\n", userScore, optimalScore)
				os.Exit(0)
			}
		default:
			fmt.Println("Unknown command.")
		}
	}
}

// splitFields splits a string by spaces or tabs, handling multiple spaces
func splitFields(s string) []string {
	return strings.Fields(s)
}