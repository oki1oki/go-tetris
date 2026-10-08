package game

const (
	BoardWidth  = 10
	BoardHeight = 20
)

var lineScores = [5]int{0, 100, 300, 500, 800}

type Cell uint8

type Board struct {
	Cells [BoardHeight][BoardWidth]Cell
}

func (b *Board) canPlace(p Piece) bool {
	for _, c := range p.Cells() {
		if c.X < 0 || c.X >= BoardWidth || c.Y < 0 || c.Y >= BoardHeight {
			return false
		}
		if b.Cells[c.Y][c.X] != 0 {
			return false
		}
	}
	return true
}

func (b *Board) lock(p Piece) {
	for _, c := range p.Cells() {
		b.Cells[c.Y][c.X] = p.Form.Cell()
	}
}

func isFullRow(row [BoardWidth]Cell) bool {
	for _, c := range row {
		if c == 0 {
			return false
		}
	}
	return true
}

func (b *Board) clearLines() int {
	dst := BoardHeight - 1
	cleared := 0
	for src := BoardHeight - 1; src >= 0; src-- {
		if isFullRow(b.Cells[src]) {
			cleared++
			continue
		}

		b.Cells[dst] = b.Cells[src]
		dst--
	}

	for ; dst >= 0; dst-- {
		b.Cells[dst] = [BoardWidth]Cell{}
	}

	return cleared
}
