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
	var hasPrintedAA bool = false
	isMaxPlayer := true
	board := nogo.NewNoGoBoard()
	step := 1
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
			// step.point.x = int(move[0] - 'A')
			// step.point.y = int(move[1] - 'A')
			yourmove := nogo.NoGoMove{
				Pos:         nogo.Position{X: int(move[0] - 'A'), Y: int(move[1] - 'A')},
				IsMaxPlayer: !isMaxPlayer,
			}
			// 处理对手行棋
			board.Move(yourmove)
			step++
			// 生成着法，并保存在step结构中
			// 注意：你需要根据你的AI逻辑来填充这一部分
			// 示例：step = generateMove(board, computerSide)
			// 处理己方行棋
			bestmove := e.GetBestMove(board)
			m, _ := bestmove.(nogo.NoGoMove)
			// 输出着法
			if !(m.Pos.X == 0 && m.Pos.Y == 0 && hasPrintedAA) {
				board.Move(bestmove)
				step++
				fmt.Printf("move %c%c\n", 'A'+m.Pos.X, 'A'+m.Pos.Y)
				hasPrintedAA = true
			}
		case "new":
			hasPrintedAA = false
			if len(args) < 2 {
				fmt.Println("Error: invalid new command")
				continue
			}

			if args[1] == "black" {
				isMaxPlayer = true
				e = gotack.NewEvaluator(gotack.AlphaBeta, 1, isMaxPlayer, func(board gotack.Board, isMaxPlayer bool, opts ...interface{}) float64 {
					return nogo.EvaluateFunc(board.(*nogo.NoGoBoard), isMaxPlayer, step)
				})
			} else {
				isMaxPlayer = false
				e = gotack.NewEvaluator(gotack.AlphaBeta, 1, isMaxPlayer, func(board gotack.Board, isMaxPlayer bool, opts ...interface{}) float64 {
					return nogo.EvaluateFunc(board.(*nogo.NoGoBoard), isMaxPlayer, step)
				})
			}
			if isMaxPlayer {
				// 处理己方行棋
				bestmove := e.GetBestMove(board)
				board.Move(bestmove)
				m, _ := bestmove.(nogo.NoGoMove)
				// 输出着法
				fmt.Printf("move %c%c\n", 'A'+m.Pos.X, 'A'+m.Pos.Y)
			}

		case "error":
			// 处理错误情况
			// ...

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
