package simulate

import (
	"fmt"

	"github.com/chainreaction/game"
	"github.com/chainreaction/utils"
)

func updateBoard(board *[][]game.Pixel, x int, y int, color string) {
	q := utils.NewQueue()
	q.Enqueue(utils.Pair{x, y})

	m := len(*board)
	n := len((*board)[0])

	(*board)[x][y].DotCount++
	(*board)[x][y].Color = color

	for !q.IsEmpty() {
		x, y = q.Dequeue()

		cnt := (*board)[x][y].DotCount

		switch cnt {
		case 2:
			if utils.IsCorner(m, n, x, y) {
				updatePixelState(board, x, y, color, q)
			}
		case 3:
			if utils.IsOnEdge(m, n, x, y) {
				updatePixelState(board, x, y, color, q)
			}
		case 4:
			updatePixelState(board, x, y, color, q)
		}

	}

	return
}

func updatePixelState(board *[][]game.Pixel, x int, y int, color string, q *utils.Queue) {
	(*board)[x][y].DotCount = 0
	(*board)[x][y].Color = ""

	m := len(*board)
	n := len((*board)[0])

	if x > 0 {
		(*board)[x-1][y].DotCount++
		(*board)[x-1][y].Color = color
		q.Enqueue(utils.Pair{x - 1, y})
	}
	if y > 0 {
		(*board)[x][y-1].DotCount++
		(*board)[x][y-1].Color = color
		q.Enqueue(utils.Pair{x, y - 1})
	}
	if x < m-1 {
		(*board)[x+1][y].DotCount++
		(*board)[x+1][y].Color = color
		q.Enqueue(utils.Pair{x + 1, y})
	}
	if y < n-1 {
		(*board)[x][y+1].DotCount++
		(*board)[x][y+1].Color = color
		q.Enqueue(utils.Pair{x, y + 1})
	}
}

// ChainReaction is called after each move and spreads the orbs on the board
func ChainReaction(gameInstance *game.Instance, move game.Move) error {
	board := gameInstance.Board

	x, y := move.XPos, move.YPos

	if x < 0 && y < 0 && x > gameInstance.Dimension && y > gameInstance.Dimension {
		return fmt.Errorf("Given positions x %v and y %v are out of range", x, y)
	}

	updateBoard(&board, x, y, gameInstance.GetPlayerByID(move.PlayerUserName).Color)

	return nil
}
