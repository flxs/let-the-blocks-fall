package gamestate

import (
	"github.com/flxs/let-the-blocks-fall/field"
)

type GameState struct {
	Field        field.Field
	Block        field.Block
	LinesCleared int
}

func New(w int, h int) GameState {
	var gs GameState
	gs.Field = field.NewField(w, h)
	gs.Block = field.NewBlock()
	gs.centerBlock()
	return gs
}

func (gs *GameState) centerBlock() {
	gs.Block.X = gs.Field.Width/2 - (gs.Block.Width / 2)
}

func (gs *GameState) NudgeLeft() {
	tmpBlck := gs.Block
	tmpBlck.X--
	if gs.Field.CanPlace(tmpBlck) {
		gs.Block.X--
	}
}

func (gs *GameState) Nudge(delta int, vertical bool) {
	oldBlck := gs.Block.Copy()
	tmpBlck := oldBlck.Copy()
	if vertical {
		tmpBlck.Y += delta
	} else {
		tmpBlck.X += delta
	}
	if gs.Field.CanPlace(tmpBlck) {
		gs.Block = tmpBlck
	} else if vertical && delta > 0 {
		// moving downward and collided with something -- drop block
		gs.Field.DrawBlock(oldBlck)
		gs.Block = field.NewBlock()
		gs.centerBlock()
	}
}

func (gs *GameState) RotateBlock() {
	tmpBlck := gs.Block.Copy()
	tmpBlck.Rotate()
	if gs.Field.CanPlace(tmpBlck) {
		gs.Block = tmpBlck
	}
}

func (gs *GameState) ClearCompleteLines() {
	gs.LinesCleared += gs.Field.ClearLines()
}
