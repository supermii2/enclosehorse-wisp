# EncloseHorse Wisp

A Go CLI client for the Enclose.Horse puzzle game.

## Features

- **Fetch Daily Puzzle**: Automatically loads the daily puzzle from the Enclose.Horse API.
- **Bonus Puzzle Support**: Load the bonus variant of the daily puzzle with the `-bonus` / `-B` flag.
- **Historical Puzzles**: Load a puzzle for any past date with the `-date` flag.
- **Board Rendering**: Displays the puzzle board in the terminal with Unicode symbols for walls, horses, cherries, apples, bees, and more.
- **Add/Remove Walls**: Place or remove walls by specifying a location (e.g., `A1`, `b2`). Both uppercase and lowercase column letters are accepted.
- **Score Calculation**: Calculates the score for the current wall configuration, including any bonus-type modifiers. If the area is not enclosed, the score is shown as `N/A`.
- **Best Score Tracking**: Tracks the best score and wall configuration found so far. Updated automatically after each add/remove operation.
- **Reload Best Walls**: Use the `b` command to reload the best wall configuration found so far.
- **Wall Budget**: Enforces a wall budget and displays the number of walls used.
- **Submit Solution**: Submit your current wall configuration to the Enclose.Horse API and see your score compared to the optimal score.
- **Input Flexibility**: Accepts both uppercase and lowercase letters for column input.
- **Clear Output**: Shows current score, best score, and wall usage after every change.

## Usage

```
horse [flags]
```

### Flags

| Flag | Description |
|------|-------------|
| `-date YYYY-MM-DD`, `-D YYYY-MM-DD` | Load the puzzle for a specific date. Defaults to today. |
| `-bonus`, `-B` | Load the bonus variant of the daily puzzle. Errors if no bonus is available. |

### Examples

```
horse                        # today's puzzle
horse -date 2026-05-01       # puzzle from May 1 2026
horse -D 2026-05-01          # same, using short flag
horse -bonus                 # today's bonus puzzle
horse -B -date 2026-05-01    # bonus puzzle from May 1 2026
```

## Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `add [location]` | `a` | Add a wall at the given location (e.g. `add A1`, `a b2`). |
| `remove [location]` | `r` | Remove a wall at the given location. |
| `best` |`b` | Reload the best wall configuration found so far. |
| `submit` | `s` | Submit your current solution and display your score vs. optimal. |

## Symbols

| Symbol | Meaning |
|--------|---------|
| `🞮` | Wall |
| `⬛` | Water / Obstacle |
| ` .` | Empty cell |
| `🍒` | Cherry (`C`) |
| `🐎` | Horse (`H`) |
| `🦄` | Unicorn (`U`) |
| `🐝` | Bee (`S`) |
| `🍏` | Apple (`G`) |

## Scoring

Base scoring (all puzzle types):

| Cell | Points |
|------|--------|
| `.`, `H`, `U`, portals | +1 each |
| `C` (cherry) | +4 total |
| `S` (bee) | −4 total |
| `G` (apple) | +11 total |

If the enclosed area touches the edge of the grid, the score is `N/A`.

### Bonus Types

| Bonus Type | Additional Rules |
|------------|-----------------|
| `costlywalls` | Each wall placed costs −6 points. |
| `lovebirds` | Both `H` (horse) and `U` (unicorn) must be enclosed in the same connected area. If either is missing or not co-enclosed, the score is `N/A`. |

## Requirements

- Go 1.18 or later

## Build Instructions

### Linux/macOS

```
go build -o horse
```

Move to PATH (optional):

```
sudo mv horse /usr/local/bin/
```

### Windows

```
go build -o horse.exe
```

Run from any location by adding the directory containing `horse.exe` to your system PATH.

## Run Without Building

```
go run .
```

---

This project is not affiliated with enclose.horse, but is a fan-made CLI for puzzle enthusiasts.

