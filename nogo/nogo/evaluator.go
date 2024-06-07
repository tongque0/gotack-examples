package nogo

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	"github.com/tongque0/gotack"
)

func EvaluateFunc(opt *gotack.EvalOptions) float64 {
	NoGoBoard, ok := opt.Board.(*NoGoBoard)
	if !ok {
		fmt.Println("EvaluateFunc called with a board type that is not *NoGoBoard")
		return 0.0
	}

	step := opt.Step
	countKey := fmt.Sprintf("count_%d", step) // 为当前步骤定义唯一的键名

	// 清除所有以 "count_" 开头的键，保证不保留过去的步数统计
	for k := range opt.Extra {
		if strings.HasPrefix(k, "count_") && k != countKey {
			delete(opt.Extra, k)
		}
	}

	// 检查当前步骤的键是否已存在于映射中，如果不存在，初始化为0
	if _, exists := opt.Extra[countKey]; !exists {
		opt.Extra[countKey] = 0
	}

	// 递增当前步骤的计数值
	opt.Extra[countKey] = opt.Extra[countKey].(int) + 1

	// 调用NoGoBoard的valuepoint方法进行评分
	value := NoGoBoard.valuepoint()
	return value
}

func (board *NoGoBoard) valuepoint() float64 {
	Pb, Pw := 0, 0
	for x, row := range board.Board {
		for y := range row {
			if !board.judgeAvailable(x, y, true) {
				Pb++
			}
			if !board.judgeAvailable(x, y, false) {
				Pw++
			}
		}
	}
	// return float64(Pb + Pw - step*2)
	return 1.0
}

// GetMHDMove 使用曼哈顿距离来选择最佳走法，并优先考虑位置 (5,5)
func (board *NoGoBoard) GetMHDMove(moves []gotack.Move) gotack.Move {
	if len(moves) == 0 {
		return nil // 如果没有最佳走法，返回 nil
	}

	// 检查中心位置 (5,5) 是否空闲
	centerPos := Position{X: 4, Y: 4} // 0-based index, so (5,5) is (4,4)
	if board.Board[centerPos.X][centerPos.Y] == 0 {
		for _, move := range moves {
			nogoMove, ok := move.(NoGoMove)
			if !ok {
				continue
			}
			if nogoMove.Pos == centerPos {
				return nogoMove
			}
		}
	}

	// 初始化随机数生成器
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 计算曼哈顿距离的最小值
	minDistances := make(map[Position]float64)
	for _, move := range moves {
		nogoMove, ok := move.(NoGoMove)
		if !ok {
			continue
		}

		minDistance := math.MaxFloat64
		for x := 0; x < 9; x++ {
			for y := 0; y < 9; y++ {
				if board.Board[x][y] != 0 {
					distance := math.Abs(float64(nogoMove.Pos.X-x)) + math.Abs(float64(nogoMove.Pos.Y-y))
					if distance < minDistance {
						minDistance = distance
					}
				}
			}
		}
		minDistances[nogoMove.Pos] = minDistance
	}

	// 找出最小值中的最大值
	var maxMinDistance float64
	for _, distance := range minDistances {
		if distance > maxMinDistance {
			maxMinDistance = distance
		}
	}

	// 找出所有具有最大最小值的点
	var bestMoves []gotack.Move
	for pos, distance := range minDistances {
		if distance == maxMinDistance {
			bestMoves = append(bestMoves, NoGoMove{Pos: pos})
		}
	}

	// 随机选择一个最优解
	return bestMoves[rng.Intn(len(bestMoves))]
}
