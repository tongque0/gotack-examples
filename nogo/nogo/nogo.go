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

func NewNoGoBoard() *NoGoBoard {
	board := NoGoBoard{}
	for x := range board.Board {
		for y := range board.Board[x] {
			board.dfs_air_visit[x][y] = false
		}
	}
	return &board
}

func (b *NoGoBoard) Print() {
	for _, row := range b.Board {
		for _, cell := range row {
			fmt.Printf("%d ", cell)
		}
		fmt.Println()
	}
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

// func (b *NoGoBoard) IsMoveLegal(move NoGoMove) bool {
// 	player := map[bool]int{true: 1, false: 2}[move.IsMaxPlayer]
// 	return b.judgeAvailable(move.Pos.X, move.Pos.Y, player)
// }
// func (b *NoGoBoard) isMoveLegal(move NoGoMove) bool {
// 	player := map[bool]int{true: 1, false: 2}[move.IsMaxPlayer]
// 	return b.judgeAvailable(move.Pos.X, move.Pos.Y, player)
// }
// func (b *NoGoBoard) airJudge(x, y, color int) bool {
// 	b.dfs_air_visit[x][y] = true
// 	hasLiberty := false
// 	directions := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

// 	for _, d := range directions {
// 		dx := x + d[0]
// 		dy := y + d[1]
// 		if b.inBorder(dx, dy) {
// 			if b.Board[dx][dy] == 0 {
// 				hasLiberty = true
// 			}
// 			if b.Board[dx][dy] == color && !b.dfs_air_visit[dx][dy] {
// 				if b.airJudge(dx, dy, color) {
// 					hasLiberty = true
// 				}
// 			}
// 		}
// 	}
// 	return hasLiberty
// }
// func (b *NoGoBoard) inBorder(x, y int) bool {
// 	return x >= 0 && y >= 0 && x < 9 && y < 9
// }

// // 判断一个落子是否合法
// func (b *NoGoBoard) judgeAvailable(x, y, color int) bool {
// 	// 如果指定位置已经有棋子，那么这个位置就不可落子
// 	if b.Board[x][y] != 0 {
// 		return false
// 	}

// 	// 模拟在指定位置落子
// 	b.Board[x][y] = color

// 	// 在使用dfs_air_visit数组前重置它，为新的搜索准备
// 	for i := range b.dfs_air_visit {
// 		for j := range b.dfs_air_visit[i] {
// 			b.dfs_air_visit[i][j] = false
// 		}
// 	}

// 	// 判断落子后，这个位置是否有气
// 	if !b.airJudge(x, y, color) {
// 		b.Board[x][y] = 0 // 如果没有气，则撤销模拟落子
// 		return false
// 	}

// 	// 检查落子位置周围的对方棋子是否因为这次落子而没有气，即被我方棋子包围
// 	for _, d := range [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
// 		dx := x + d[0]
// 		dy := y + d[1]
// 		opponentColor := 3 - color // 根据当前颜色确定对方颜色
// 		if b.inBorder(dx, dy) && b.Board[dx][dy] == opponentColor && !b.dfs_air_visit[dx][dy] {
// 			// 如果周围的对方棋子没有气，那么这个落子也是非法的
// 			if !b.airJudge(dx, dy, opponentColor) {
// 				b.Board[x][y] = 0 // 撤销模拟落子
// 				return false
// 			}
// 		}
// 	}

// 	// 如果所有的检查都通过了，那么这个落子位置是合法的，撤销模拟落子，并返回真
// 	b.Board[x][y] = 0 // 回溯
// 	return true
// }
