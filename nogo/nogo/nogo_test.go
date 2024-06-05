package nogo

import (
	"testing"
)

// 测试 GetAllMoves 函数，尤其关注边角和部分填充的棋盘
func TestGetAllMoves(t *testing.T) {
	board := NewNoGoBoard()

	// 设定棋盘，假设 1 表示 MaxPlayer 的棋子，-1 表示 MinPlayer 的棋子
	board.Board[0][0] = -1
	board.Board[0][1] = 1
	board.Board[1][0] = 0
	board.Board[8][8] = 1
	board.Board[7][8] = -1
	board.Board[8][7] = -1

	// 检查边角位置
	move1s := board.GetAllMoves(true) // 假设为 MaxPlayer 获取所有合法移动
	move2s := board.GetAllMoves(false) // 假设为 MaxPlayer 获取所有合法移动

	// 打印所有合法移动
	if len(move1s) == 0 {
		t.Log("No moves available")
	} else {
		t.Log("Available move1s:", move1s)
		t.Log("Available move2s:", move2s)
	}
}
