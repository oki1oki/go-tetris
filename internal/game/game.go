package game

import "math/rand"

const (
	fallInterval     = 30
	minFallInterval  = 3
	levelSpeedup     = 3
	softDropInterval = 5
	linesPerLevel    = 5
)

// Корды новой фигуры
const (
	spawnX = 3
	spawnY = 0
)

type Game struct {
	Board    Board
	Current  Piece
	NextForm Form
	Score    int
	Lines    int
	Level    int
	GameOver bool

	fallTimer int
	bag       []Form
}

func New() *Game {
	g := &Game{}
	g.start()
	return g
}

func (g *Game) Restart() {
	*g = Game{}
	g.start()
}

func (g *Game) Shadow() Piece {
	p := g.Current
	for {
		next := p
		next.Y++
		if !g.Board.canPlace(next) {
			return p
		}
		p = next
	}
}

func (g *Game) nextForm() Form {
	if len(g.bag) == 0 {
		for f := range FormCount {
			g.bag = append(g.bag, f)
		}

		rand.Shuffle(len(g.bag), func(i, j int) { g.bag[i], g.bag[j] = g.bag[j], g.bag[i] })
	}

	f := g.bag[len(g.bag)-1]
	g.bag = g.bag[:len(g.bag)-1]

	return f
}

func (g *Game) spawn() {
	p := Piece{Form: g.NextForm, X: spawnX, Y: spawnY, Rotation: Rot0}

	if !g.Board.canPlace(p) {
		g.GameOver = true
	}

	g.Current = p
	g.NextForm = g.nextForm()
}

func (g *Game) tryMove(p Piece) bool {
	if !g.Board.canPlace(p) {
		return false
	}
	g.Current = p
	return true
}

func (g *Game) start() {
	g.Level = 1
	g.NextForm = g.nextForm()
	g.spawn()
}

type Input struct {
	Left, Right         bool
	RotateCW, RotateCCW bool
	HardDrop            bool
	SoftDrop            bool
}

func (g *Game) Update(in Input) {
	if g.GameOver {
		return
	}

	if in.Left {
		g.shift(-1)
	}
	if in.Right {
		g.shift(1)
	}
	if in.RotateCW {
		g.rotate(g.Current.Rotation.Next())
	}
	if in.RotateCCW {
		g.rotate(g.Current.Rotation.Prev())
	}

	if in.HardDrop {
		g.hardDrop()
		return
	}

	g.tick(in.SoftDrop)
}

func (g *Game) shift(dx int) {
	moved := g.Current
	moved.X += dx
	g.tryMove(moved)
}

// На сколько клеток по x пробуем сдвинуть фигуру, если поворот на месте не проходит.
var kickOffsets = [...]int{0, -1, 1, -2, 2}

func (g *Game) rotate(r Rotation) {
	for _, dx := range kickOffsets {
		moved := g.Current
		moved.Rotation = r
		moved.X += dx
		if g.tryMove(moved) {
			return
		}
	}
}

func (g *Game) hardDrop() {
	landed := g.Shadow()
	g.Score += 2 * (landed.Y - g.Current.Y)
	g.Current = landed

	g.lockPiece()
	g.fallTimer = 0
}

// tick вызывается каждый кадр. Когда накопилось достаточно тиков,
// фигура опускается на клетку или приземляется.
func (g *Game) tick(softDrop bool) {
	g.fallTimer++

	interval := g.fallInterval()
	if softDrop {
		interval = min(softDropInterval, interval)
	}
	if g.fallTimer < interval {
		return
	}
	g.fallTimer = 0

	moved := g.Current
	moved.Y++
	if g.tryMove(moved) {
		if softDrop {
			g.Score++
		}
		return
	}

	g.lockPiece()
}

// Интервал между тиками по уровням
func (g *Game) fallInterval() int {
	return max(minFallInterval, fallInterval-(g.Level-1)*levelSpeedup)
}

func (g *Game) lockPiece() {
	g.Board.lock(g.Current)

	cleared := g.Board.clearLines()
	g.Lines += cleared
	g.Score += lineScores[cleared]
	g.Level = g.Lines/linesPerLevel + 1

	g.spawn()
}
