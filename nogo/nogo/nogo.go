package nogo

import (
	"fmt"

	"github.com/tongque0/gotack"
)

type Position struct {
	X, Y int
}

type NoGoMove struct {
	Pos         Position
	IsMaxPlayer bool
}

// NoGoBoard represents the game board state.
type NoGoBoard struct {
	Board         [9][9]int  // 0表示空位，1表示黑子，2表示白子
	dfs_air_visit [9][9]bool // 记录每个点的气
}

// Print implements gotack.Board.
func (board *NoGoBoard) Print() {
	panic("unimplemented")
}

func NewNoGoBoard() *NoGoBoard {
	board := NoGoBoard{}
	for x := range board.Board {
		for y := range board.Board[x] {
			board.dfs_air_visit[x][y] = false
		}
	}
	return &board
}

func (m NoGoMove) String() string {
	return fmt.Sprintf("To (%d,%d)", m.Pos.X, m.Pos.Y)
}

func (b *NoGoBoard) GetAllMoves(isMaxPlayer bool) []gotack.Move {
	var moves []gotack.Move
	for x, row := range b.Board {
		for y := range row {
			if b.judgeAvailable(x, y, isMaxPlayer) {
				move := NoGoMove{
					Pos:         Position{X: x, Y: y},
					IsMaxPlayer: isMaxPlayer,
				}
				moves = append(moves, move)
			}
		}
	}
	return moves
}

func (b *NoGoBoard) Move(move gotack.Move) {
	m, ok := move.(NoGoMove)
	if !ok {
		fmt.Println("Invalid move type1")
		return
	}
	player := map[bool]int{true: 1, false: 2}[m.IsMaxPlayer]
	b.Board[m.Pos.X][m.Pos.Y] = player
}

func (b *NoGoBoard) UndoMove(move gotack.Move) {
	m, ok := move.(NoGoMove)
	if !ok {
		fmt.Println("Invalid move type2")
		return
	}
	b.Board[m.Pos.X][m.Pos.Y] = 0
}

func (b *NoGoBoard) IsGameOver() bool {
	for x := range b.Board {
		for y := range b.Board[x] {
			if b.judgeAvailable(x, y, true) || b.judgeAvailable(x, y, false) {
				return false
			}
		}
	}
	return true
}

func (b *NoGoBoard) Hash() uint64 {
	// 在这里实现生成棋盘状态的哈希值的逻辑。
	// 对于不需要置换表的游戏，可以返回0或固定值作为哈希值。
	return 0 // 示例返回值，实际开发中应根据棋盘状态生成唯一的哈希值。
}
func (b *NoGoBoard) Clone() gotack.Board {
	// 创建一个新的NoGoBoard实例
	cloned := NewNoGoBoard()

	// 复制棋盘上的棋子位置
	cloned.Board = b.Board

	// 复制访问状态数组
	cloned.dfs_air_visit = b.dfs_air_visit

	// 返回复制的棋盘
	return cloned
}
