package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tnogo/nogo"

	"github.com/tongque0/gotack"
)

const INF = 0x3f3f3f3f

var (
	line        string
	step        int
	board       *nogo.NoGoBoard
	IsMaxPlayer bool
)

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line = sc.Text()
		if line == "name?" {
			fmt.Println("name NoGoTack")
		} else if line == "quit" {
			os.Exit(0)
		} else if strings.HasPrefix(line, "new") {
			step = 1
			words := strings.Split(line, " ")
			board = nogo.NewNoGoBoard()
			if words[1] == "black" {
				IsMaxPlayer = true
				runSearch()
			} else {
				IsMaxPlayer = false
			}
		} else if strings.HasPrefix(line, "move") {
			words := strings.Split(line, " ")
			move := words[1]
			oppMove := nogo.NoGoMove{
				Pos:         nogo.Position{X: int(move[0] - 'A'), Y: int(move[1] - 'A')},
				IsMaxPlayer: !IsMaxPlayer,
			}
			board.Move(oppMove)
			step++
			if !board.IsGameOver() {
				runSearch()
			}
		} else if line == "end" {
			continue
		} else {
			continue
		}
	}
}

func runSearch() {
	var evaluator *gotack.Evaluator
	// Determine search depth based on the number of moves made so far
	var depth int
	switch {
	case step < 23:
		depth = 3
	case step < 50:
		depth = 4
	case step < 70:
		depth = 5
	default:
		depth = 6
	}

	// Initialize the evaluator with the current game state and parameters
	evaluator = gotack.NewEvaluator(
		gotack.AlphaBeta,
		gotack.NewEvaluatorOptions(
			gotack.WithBoard(board),
			gotack.WithDepth(depth),
			gotack.WithIsMaxPlayer(IsMaxPlayer),
			gotack.WithStep(step),
			gotack.WithIsDetail(true),
		),
	)
	selectedMove := evaluator.GetBestMove()
	validMove := board.GetMHDMove(selectedMove)
	board.Move(validMove)
	noGoMove, valid := validMove.(nogo.NoGoMove)
	if !valid {
		fmt.Println("Invalid move")
		return
	}
	fmt.Printf("move %c%c\n", 'A'+noGoMove.Pos.X, 'A'+noGoMove.Pos.Y)
	step++
}
