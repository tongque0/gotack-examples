package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tnogo/nogo"

	"github.com/tongque0/gotack"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	isMaxPlayer := true
	board := nogo.NewNoGoBoard()
	var e *gotack.Evaluator
	for {
		fmt.Print("Ready for input: ")
		if !scanner.Scan() {
			break // 退出循环
		}
		message := scanner.Text()
		args := strings.Fields(message)

		switch args[0] {
		case "move":
			move := args[1]
			yourmove := nogo.NoGoMove{
				Pos:         nogo.Position{X: int(move[0] - 'A'), Y: int(move[1] - 'A')},
				IsMaxPlayer: !isMaxPlayer,
			}
			// 处理对手行棋
			board.Move(yourmove)
			e.EvalOptions.Step++
			// 生成着法，并保存在step结构中
			// 注意：你需要根据你的AI逻辑来填充这一部分
			// 示例：step = generateMove(board, computerSide)
			// 处理己方行棋
			bestmove := e.GetBestMove()
			mm := board.GetMHDMove(bestmove)
			m, _ := mm.(nogo.NoGoMove)
			// 输出着法
			board.Move(mm)
			e.EvalOptions.Step++
			fmt.Printf("move %c%c\n", 'A'+m.Pos.X, 'A'+m.Pos.Y)
		case "new":
			if len(args) < 2 {
				fmt.Println("Error: invalid new command")
				continue
			}

			if args[1] == "black" {
				isMaxPlayer = true
				e = gotack.NewEvaluator(
					gotack.AlphaBeta,
					gotack.NewEvaluatorOptions(
						gotack.WithBoard(board),
						gotack.WithDepth(3),
						gotack.WithIsMaxPlayer(isMaxPlayer),
						gotack.WithStep(1),
						gotack.WithIsDetail(true),
					),
					func(opts *gotack.EvalOptions) float64 {
						return nogo.EvaluateFunc(opts)
					},
				)
			} else {
				isMaxPlayer = false
				e = gotack.NewEvaluator(
					gotack.AlphaBeta,
					gotack.NewEvaluatorOptions(
						gotack.WithBoard(board),
						gotack.WithDepth(3),
						gotack.WithIsMaxPlayer(isMaxPlayer),
						gotack.WithStep(1),
						gotack.WithIsDetail(true),
					),
					func(opts *gotack.EvalOptions) float64 {
						return nogo.EvaluateFunc(opts)
					},
				)
			}
			if isMaxPlayer {
				// 处理己方行棋
				bestmove := e.GetBestMove()
				board.Move(bestmove[0])
				mm := board.GetMHDMove(bestmove)
				m, _ := mm.(nogo.NoGoMove)
				// 输出着法
				board.Move(mm)
				e.EvalOptions.Step++
				fmt.Printf("move %c%c\n", 'A'+m.Pos.X, 'A'+m.Pos.Y)
			}

		case "error":
			fmt.Println("error")
			return

		case "name?":
			// 输出引擎名
			fmt.Println("name GoTack-NoGo")

		case "end":
			// 对局结束处理
			fmt.Println("Thanks")

		case "quit":
			// 退出引擎
			fmt.Println("Quit!")
			return
		}
	}
}
