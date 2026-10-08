package main

import (
	"bytes"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"

	"tetris/internal/game"
)

const (
	cellSize    = 30
	padding     = 20
	sidePanelW  = 160
	statBlockH  = 70
	previewCell = 24
	nextBlockH  = 130
	panelX      = padding + game.BoardWidth*cellSize + padding
)

const (
	fontLabel    = 14
	fontValue    = 28
	fontValueBig = 32
	fontBadge    = 22
	fontHint     = 18
	fontTitle    = 40
)

const (
	screenW = game.BoardWidth*cellSize + 2*padding + sidePanelW
	screenH = game.BoardHeight*cellSize + 2*padding
)

var emptyColor = color.RGBA{R: 40, G: 40, B: 40, A: 255}
var redColor = color.RGBA{R: 255, A: 255}
var labelColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
var goldColor = color.RGBA{R: 255, G: 215, A: 255}
var ghostColor = color.RGBA{R: 110, G: 110, B: 110, A: 255}

var colors = [game.FormCount + 1]color.RGBA{
	0: emptyColor,
	1: {R: 0, G: 255, B: 255, A: 255}, // I
	2: {R: 255, G: 255, B: 0, A: 255}, // O
	3: {R: 128, G: 0, B: 128, A: 255}, // T
	4: {R: 0, G: 255, B: 0, A: 255},   // S
	5: {R: 255, G: 0, B: 0, A: 255},   // Z
	6: {R: 0, G: 0, B: 255, A: 255},   // J
	7: {R: 255, G: 165, B: 0, A: 255}, // L
}

var fontSource *text.GoTextFaceSource

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}

	fontSource = s
}

func newFace(size float64) *text.GoTextFace {
	return &text.GoTextFace{Source: fontSource, Size: size}
}

func drawText(screen *ebiten.Image, s string, size, x, y float64, c color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(c)
	text.Draw(screen, s, newFace(size), op)
}

func drawTextCentered(screen *ebiten.Image, s string, size, y float64, c color.Color) {
	w, _ := text.Measure(s, newFace(size), 0)
	x := (screenW - w) / 2

	drawText(screen, s, size, x, y, c)
}

func drawCellAt(screen *ebiten.Image, px, py, size float32, c color.RGBA) {
	vector.FillRect(screen, px, py, size-1, size-1, c, false)
}

func drawCell(screen *ebiten.Image, x, y int, c color.RGBA) {
	drawCellAt(screen, float32(padding+x*cellSize), float32(padding+y*cellSize), cellSize, c)
}

func drawCellOutline(screen *ebiten.Image, x, y int, c color.RGBA) {
	const stroke = 1
	px := float32(padding + x*cellSize)
	py := float32(padding + y*cellSize)
	size := float32(cellSize - 1 - stroke)

	vector.StrokeRect(screen, px+stroke/2.0, py+stroke/2.0, size, size, stroke, c, false)
}

func drawNext(screen *ebiten.Image, f game.Form, y float64) {
	drawText(screen, "NEXT", fontLabel, panelX, y, labelColor)

	p := game.Piece{Form: f}
	top := float32(y) + 24
	for _, c := range p.Cells() {
		drawCellAt(screen, float32(panelX+c.X*previewCell), top+float32(c.Y*previewCell), previewCell, colors[f.Cell()])
	}
}

func dim(screen *ebiten.Image, alpha uint8) {
	vector.FillRect(screen, 0, 0, screenW, screenH, color.RGBA{A: alpha}, false)
}

func drawStat(screen *ebiten.Image, label string, value int, y float64) {
	drawText(screen, label, fontLabel, panelX, y, labelColor)
	drawText(screen, strconv.Itoa(value), fontValue, panelX, y+20, color.White)
}

func drawStatCentered(screen *ebiten.Image, label string, value int, y float64) {
	drawTextCentered(screen, label, fontLabel, y, labelColor)
	drawTextCentered(screen, strconv.Itoa(value), fontValueBig, y+20, color.White)
}
