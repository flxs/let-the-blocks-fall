package field

import (
	"math/rand"
)

/*
Field represents the playing field, a Width x Height grid of integers,
each of which represents a color. The current block is random.
*/
type Field struct {
	Width  int
	Height int
	Matrix []int
}

/*
NewField constructs a blank Field instance
*/
func NewField(w int, h int) Field {
	var f Field
	f.Width = w
	f.Height = h
	f.Matrix = make([]int, w*h)
	for i := range f.Matrix {
		f.Matrix[i] = 0
	}
	return f
}

/*
Block represents the current block with its shape (Matrix), dimensions (of the Matrix)
and position (X, Y)
*/
type Block struct {
	Width  int
	Height int
	X      int
	Y      int
	Matrix []int
}

/*
NewBlock constructs a random tetroid Block instance
*/
func NewBlock() Block {
	var block Block
	switch rand.Intn(7) {
	case 0:
		block.Matrix = []int{0, 1, 0, 1, 1, 1, 0, 0, 0}
		block.Width = 3
		block.Height = 3
		block.X = 0
		block.Y = 0
	case 1:
		block.Matrix = []int{2, 2, 2, 2}
		block.Width = 2
		block.Height = 2
		block.X = 0
		block.Y = 0

	case 2:
		block.Matrix = []int{
			3, 3, 0,
			0, 3, 0,
			0, 3, 0}
		block.Width = 3
		block.Height = 3
		block.X = 0
		block.Y = 0
	case 3:
		block.Matrix = []int{
			0, 4, 4,
			0, 4, 0,
			0, 4, 0}
		block.Width = 3
		block.Height = 3
		block.X = 0
		block.Y = 0
	case 4:
		block.Matrix = []int{
			0, 0, 0, 0,
			5, 5, 5, 5,
			0, 0, 0, 0,
			0, 0, 0, 0}
		block.Width = 4
		block.Height = 4
		block.X = 0
		block.Y = 0
	case 5:
		block.Matrix = []int{
			6, 6, 0,
			0, 6, 6,
			0, 0, 0}
		block.Width = 3
		block.Height = 3
		block.X = 0
		block.Y = 0
	case 6:
		block.Matrix = []int{
			0, 7, 7,
			7, 7, 0,
			0, 0, 0}
		block.Width = 3
		block.Height = 3
		block.X = 0
		block.Y = 0
	}
	return block
}

func (b Block) Copy() Block {
	var tmp Block
	tmp.Matrix = make([]int, len(b.Matrix))
	copy(tmp.Matrix, b.Matrix)
	tmp.X = b.X
	tmp.Y = b.Y
	tmp.Width = b.Width
	tmp.Height = b.Height
	return tmp
}

/*
DrawBlock draws the given Block onto the Field
*/
func (f Field) DrawBlock(b Block) {
	for i := 0; i < len(b.Matrix); i++ {
		x := i%b.Width + b.X
		y := i/b.Width + b.Y

		targeti := y*f.Width + x

		if b.Matrix[i] != 0 {
			f.Matrix[targeti] = b.Matrix[i]
		}
	}
}

/*
Copy creates a faithful copy of a Field
*/
func (f Field) Copy() Field {
	var tf Field
	tf.Width = f.Width
	tf.Height = f.Height
	tf.Matrix = make([]int, tf.Width*tf.Height)
	copy(tf.Matrix, f.Matrix)
	return tf
}

/*
CanMove checks whether the given Block can be moved by dx and dy
without colliding with existing blocks on the Field
*/
func (f Field) CanMove(b Block, dx int, dy int) bool {
	for i := 0; i < len(b.Matrix); i++ {
		x := i%b.Width + b.X + dx
		y := i/b.Width + b.Y + dy

		if b.Matrix[i] == 0 {
			continue
		}

		if x < 0 || y < 0 || x >= f.Width || y >= f.Height {
			return false
		}

		targeti := y*f.Width + x

		if f.Matrix[targeti] != 0 && b.Matrix[i] != 0 {
			return false
		}
	}
	return true
}

func (f Field) CanPlace(b Block) bool {
	for i := 0; i < len(b.Matrix); i++ {
		x := i%b.Width + b.X
		y := i/b.Width + b.Y

		if b.Matrix[i] == 0 {
			continue
		}

		if x < 0 || y < 0 || x >= f.Width || y >= f.Height {
			return false
		}

		targeti := y*f.Width + x

		if f.Matrix[targeti] != 0 && b.Matrix[i] != 0 {
			return false
		}
	}
	return true
}

/*
ClearLines clears complete lines
*/
func (f *Field) ClearLines() int {
	var linesCleared = 0
	for l := f.Height - 1; l >= 0; l-- {
		linestart := l * f.Width
		lineend := (l+1)*f.Width - 1
		line := f.Matrix[linestart:lineend]

		// check if line contains no empty fields (i.e. is complete)
		var lineComplete = true
		for i := range line {
			if line[i] == 0 {
				lineComplete = false
			}
		}
		if !lineComplete {
			continue
		}

		// shift all lines above complete line downward by one line
		for i := linestart - 1; i >= 0; i-- {
			f.Matrix[i+f.Width] = f.Matrix[i]
		}

		// erase topmost line
		for i := 0; i < f.Width; i++ {
			f.Matrix[i] = 0
		}

		linesCleared++
	}
	return linesCleared
}

/*
Rotate rotates the Block clockwise by 90 degrees
*/
func (b *Block) Rotate() {
	matrix := make([]int, len(b.Matrix))
	copy(matrix, b.Matrix)
	for i := 0; i < len(b.Matrix); i++ {
		x := i % b.Width
		y := i / b.Width

		j := x*b.Height + ((b.Height - 1) - y)

		b.Matrix[j] = matrix[i]
	}
}
