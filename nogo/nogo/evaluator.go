package nogo

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/tongque0/gotack"
)

// EvaluateFunc 根据当前棋盘状态评估分数，这里只是一个简化的示例。
// 实际的评估函数可能会更加复杂。
func EvaluateFunc(board gotack.Board, isMaxPlayer bool, opts ...interface{}) float64 {

	NoGoBoard, ok := board.(*NoGoBoard)
	if !ok {
		fmt.Println("EvaluateFunc called with a board type that is not *NoGoBoard")
		return 0.0 // 或者处理这种情况的其他方式，比如返回一个默认分数或错误处理
	}
	// 解析opts，获取轮数
	var turn int
	if len(opts) > 0 {
		turnVal, ok := opts[0].(int) // 类型断言，将opts的第一个元素转换为整数
		if !ok {
			fmt.Println("EvaluateFunc: Expected an integer for the turn, but got something else.")
			// 处理错误或者使用默认值
		} else {
			turn = turnVal
		}
	}
	// TODO: 根据GameBoard的当前状态和isMaxPlayer，实现具体的评估逻辑
	// 这里需要开发者根据具体游戏的规则和逻辑来填充评分算法
	// 示例中仅调用GameBoard的Print方法来假设进行了评估过程
	// NoGoBoard.Print()

	// 返回评估分数，这个分数应该基于GameBoard的状态和isMaxPlayer计算得出
	value := NoGoBoard.valuepoint(turn)
	return value // 示例返回值，实际开发中应根据游戏逻辑修改
}

func (board *NoGoBoard) valuepoint(step int) float64 {
	// Pb, Pw := 0, 0
	// for x, row := range board.Board {
	// 	for y, _ := range row {
	// 		if !board.judgeAvailable(x, y, 1) {
	// 			Pb++
	// 		}
	// 		if !board.judgeAvailable(x, y, 2) {
	// 			Pw++
	// 		}
	// 	}
	// }
	// return float64(Pb + Pw - step*2)
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Float64()
}
