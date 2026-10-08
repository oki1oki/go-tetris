package game

type Piece struct {
	Form     Form
	Rotation Rotation
	X, Y     int
}

func (p Piece) Cells() Shape {
	shape := formsTable[p.Form][p.Rotation]
	for i := range shape {
		shape[i].X += p.X
		shape[i].Y += p.Y
	}
	return shape
}
