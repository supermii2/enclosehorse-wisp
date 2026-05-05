# EncloseHorse Wisp

A Go CLI client for the Enclose.Horse puzzle game.

## Features

- **Fetch Daily Puzzle**: Automatically loads the daily puzzle from the Enclose.Horse API.
- **Board Rendering**: Displays the puzzle board in the terminal with Unicode symbols for walls, horses, cherries, apples, bees, and more.
- **Add/Remove Walls**: Place or remove walls by specifying a location (e.g., `A1`, `b2`). Both uppercase and lowercase column letters are accepted.
- **Score Calculation**: Calculates the score for the current wall configuration. If the area is not enclosed, the score is shown as `N/A`.
- **Best Score Tracking**: Tracks the best score and wall configuration found so far. The best score is updated automatically after each add/remove operation.
- **Reload Best Walls**: Use the `b` command to reload the best wall configuration found so far.
- **Wall Budget**: Enforces a wall budget and displays the number of walls used.
- **Submit Solution**: Submit your current wall configuration to the Enclose.Horse API and see your score compared to the optimal score.
- **Input Flexibility**: Accepts both uppercase and lowercase letters for column input.
- **Clear Output**: Shows current score, best score, and wall usage after every change.

## Commands

- `add [location]` or `a [location]`: Add a wall at the specified location (e.g., `add A1` or `a b2`).
- `remove [location]` or `r [location]`: Remove a wall at the specified location.
- `b`: Reload the best wall configuration found so far.
- `submit` or `s`: Submit your current solution.

## Symbols

- `🞮` — Wall
- `⬛` — Water/Obstacle
- ` .` — Empty cell
- `🍒` — Cherry
- `🐎` — Horse
- `🐝` — Bee
- `🍏` — Apple

## Scoring

- Enclosed area is scored based on the puzzle rules:
  - `.` and `H` (horse): 1 point each
  - `C` (cherry): 4 points
  - `S` (bee): -4 points
  - `G` (apple): 11 points
- If the area is not enclosed, the score is shown as `N/A`.

## Requirements

- Go 1.18 or later

## Usage

```
go run .
```

Follow the prompts to interact with the puzzle.

---

This project is not affiliated with enclose.horse, but is a fan-made CLI for puzzle enthusiasts.
