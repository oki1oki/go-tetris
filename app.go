package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"tetris/internal/game"
	"tetris/internal/storage"
)

type App struct {
	game      *game.Game
	highScore int
	newRecord bool
	paused    bool
	left      keyRepeat
	right     keyRepeat
}

func newApp() *App {
	return &App{
		game:      game.New(),
		highScore: storage.LoadHighScore(),
	}
}

func (a *App) Update() error {
	g := a.game

	if g.GameOver {
		if g.Score > a.highScore {
			a.highScore = g.Score
			a.newRecord = true
			if err := storage.SaveHighScore(a.highScore); err != nil {
				log.Println(err)
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.Restart()
			a.newRecord = false
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		a.paused = !a.paused
	}

	if a.paused {
		return nil
	}

	g.Update(game.Input{
		Left:      a.left.fire(ebiten.IsKeyPressed(ebiten.KeyLeft)),
		Right:     a.right.fire(ebiten.IsKeyPressed(ebiten.KeyRight)),
		RotateCW:  inpututil.IsKeyJustPressed(ebiten.KeyUp),
		RotateCCW: inpututil.IsKeyJustPressed(ebiten.KeyZ),
		HardDrop:  inpututil.IsKeyJustPressed(ebiten.KeySpace),
		SoftDrop:  ebiten.IsKeyPressed(ebiten.KeyDown),
	})

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	a.drawBoard(screen)
	a.drawPanel(screen)

	if a.game.GameOver {
		a.drawGameOver(screen)
	} else if a.paused {
		a.drawPaused(screen)
	}
}

func (a *App) drawBoard(screen *ebiten.Image) {
	g := a.game

	vector.StrokeRect(screen, float32(padding-1), float32(padding-1), float32(game.BoardWidth*cellSize+2), float32(game.BoardHeight*cellSize+2), 1, color.White, false)
	for y := range game.BoardHeight {
		for x := range game.BoardWidth {
			drawCell(screen, x, y, colors[g.Board.Cells[y][x]])
		}
	}

	for _, c := range g.Shadow().Cells() {
		drawCellOutline(screen, c.X, c.Y, ghostColor)
	}

	for _, c := range g.Current.Cells() {
		drawCell(screen, c.X, c.Y, colors[g.Current.Form.Cell()])
	}
}

func (a *App) drawPanel(screen *ebiten.Image) {
	g := a.game
	y := float64(padding)

	drawStat(screen, "SCORE", g.Score, y)
	y += statBlockH

	drawStat(screen, "BEST", a.highScore, y)
	y += statBlockH

	drawNext(screen, g.NextForm, y)
	y += nextBlockH

	drawStat(screen, "LEVEL", g.Level, y)
	y += statBlockH

	drawStat(screen, "LINES", g.Lines, y)
}

func (a *App) drawPaused(screen *ebiten.Image) {
	dim(screen, 180)

	y := float64(screenH/2 - 40)
	drawTextCentered(screen, "PAUSED", fontTitle, y, color.White)
	drawTextCentered(screen, "Press P to resume", fontHint, y+60, labelColor)
}

func (a *App) drawGameOver(screen *ebiten.Image) {
	dim(screen, 220)

	y := float64(screenH/2 - 120)
	drawTextCentered(screen, "GAME OVER", fontTitle, y, redColor)

	if a.newRecord {
		drawTextCentered(screen, "NEW RECORD!", fontBadge, y+60, goldColor)
	}

	drawStatCentered(screen, "SCORE", a.game.Score, y+105)
	drawStatCentered(screen, "BEST", a.highScore, y+170)

	drawTextCentered(screen, "Press R to restart", fontHint, y+250, labelColor)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return screenW, screenH
}
