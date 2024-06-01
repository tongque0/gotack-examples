package nogo

// inBorder 检查给定的坐标是否在棋盘范围内
func (b *NoGoBoard) inBorder(x, y int) bool {
	return x >= 0 && y >= 0 && x < 9 && y < 9
}

// judgeAvailable 检查在 (x, y) 位置落子是否合法。
func (b *NoGoBoard) judgeAvailable(x, y int, isMaxPlayer bool) bool {
	if b.Board[x][y] != 0 {
		return false // 已有棋子，不可落子
	}

	player := map[bool]int{true: 1, false: 2}[isMaxPlayer] // 假设1为Max玩家的棋子，2为Min玩家的棋子
	b.Board[x][y] = player                                 // 模拟在 (x, y) 落子

	// 重置访问标记数组
	for i := range b.dfs_air_visit {
		for j := range b.dfs_air_visit[i] {
			b.dfs_air_visit[i][j] = false
		}
	}

	// 判断是否自杀
	if !b.airJudge(x, y, player) {
		b.Board[x][y] = 0 // 撤销落子
		return false
	}

	// 检查相邻敌方棋子是否无气，可能导致合法的提子
	for _, d := range [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}} {
		dx := x + d[0]
		dy := y + d[1]
		if b.inBorder(dx, dy) && b.Board[dx][dy] == 3-player && !b.dfs_air_visit[dx][dy] { // 敌方棋子颜色为3-player
			if !b.airJudge(dx, dy, 3-player) {
				b.Board[x][y] = 0 // 撤销落子
				return false
			}
		}
	}

	b.Board[x][y] = 0 // 撤销落子
	return true
}

// airJudge 判断 (x, y) 位置的棋子是否有气
func (b *NoGoBoard) airJudge(x, y, color int) bool {
	if b.dfs_air_visit[x][y] {
		return false
	}
	b.dfs_air_visit[x][y] = true
	hasLiberty := false
	directions := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

	for _, d := range directions {
		dx := x + d[0]
		dy := y + d[1]
		if b.inBorder(dx, dy) {
			if b.Board[dx][dy] == 0 { // 直接相邻有空位，有气
				hasLiberty = true
			} else if b.Board[dx][dy] == color && !b.dfs_air_visit[dx][dy] { // 相邻同色棋子，递归检查
				if b.airJudge(dx, dy, color) {
					hasLiberty = true
				}
			}
		}
	}
	return hasLiberty
}
