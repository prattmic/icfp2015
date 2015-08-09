package main

import (
	//"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"io"
	"math"
)

const ( // Sizes of hexagons to render
	tinySize  = 0
	smallSize = 1
	medSize   = 2
	largeSize = 3
)

type HexSize struct {
	size  int // Hexagon size (width and height) in pixels
	horiz int // Horizontal spacing in pixels (incl. hexagon size)
	vert  int // Vertical spacing in pixels (incl. hexagon size)
	pivot int // Pivot size in pixels
}

var (
	// Colors to draw with
	whiteColor  = color.RGBA{255, 255, 255, 255} // Background
	yellowColor = color.RGBA{255, 255, 0, 255}   // Filled
	redColor    = color.RGBA{255, 0, 0, 255}     // Unit
	greyColor   = color.RGBA{200, 200, 200, 255} // Empty
	blackColor  = color.RGBA{0, 0, 0, 255}       // Pivot
	pal         = []color.Color{whiteColor, yellowColor, redColor, greyColor, blackColor}
	// Prebuilt uniform images of all the colors
	white  = image.Uniform{whiteColor}
	yellow = image.Uniform{yellowColor}
	red    = image.Uniform{redColor}
	grey   = image.Uniform{greyColor}
	black  = image.Uniform{blackColor}
	// Hexagon presize
	tinyHex  = HexSize{size: 3, horiz: 4, vert: 3, pivot: 1}
	smallHex = HexSize{size: 5, horiz: 6, vert: 4, pivot: 2}
	medHex   = HexSize{size: 10, horiz: 12, vert: 9, pivot: 4}
	largeHex = HexSize{size: 20, horiz: 24, vert: 18, pivot: 8}
)

type HexMask struct {
	cell Cell    // Game board coordinate
	hs   HexSize // Size to draw things
}

func (hm *HexMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (hm *HexMask) origin() image.Point {
	px := hm.cell.X*hm.hs.horiz + (hm.cell.Y%2)*hm.hs.horiz/2
	py := hm.cell.Y * hm.hs.vert
	return image.Point{px, py}
}

func (hm *HexMask) Bounds() image.Rectangle {
	return image.Rect(0, 0, hm.hs.size, hm.hs.size)
}

func (hm *HexMask) At(xx, yy int) color.Color {
	x := float64(xx)
	y := float64(yy)
	k := float64(hm.hs.size)
	h := float64(hm.hs.size-1) / 2.0
	if (k - 2*math.Abs(y-h)) <= math.Abs(x-h) {
		return color.Alpha{0}
	}
	return color.Alpha{255}
}

func drawHex(m draw.Image, board image.Rectangle, cell Cell, hs HexSize, i image.Image) {
	hm := HexMask{cell, hs}
	r := hm.Bounds().Add(board.Min).Add(hm.origin())
	draw.DrawMask(m, r, i, image.ZP, &hm, image.ZP, draw.Over)
}

func drawPivot(m draw.Image, board image.Rectangle, cell Cell, hs HexSize) {
	h := (hs.size - hs.pivot) / 2
	x := cell.X*hs.horiz + (cell.Y%2)*hs.horiz/2 + h
	y := cell.Y*hs.vert + h
	min := board.Min.Add(image.Point{x, y})
	max := min.Add(image.Point{hs.pivot, hs.pivot})
	draw.Draw(m, image.Rectangle{min, max}, &black, image.ZP, draw.Over)
}

type GameRenderer struct {
	width, height, border int
	frames                []*image.Paletted
	hs                    HexSize
}

func NewGameRenderer(g *Game, border int, size int) *GameRenderer {
	hs := largeHex // Default
	switch size {
	case tinySize:
		hs = tinyHex
	case smallSize:
		hs = smallHex
	case medSize:
		hs = medHex
	case largeSize:
		hs = largeHex
	}
	return &GameRenderer{width: g.b.width, height: g.b.height, border: border, hs: hs}
}

func gameFillColor(g *Game, x, y int) image.Image {
	c := Cell{x, y}

	if c.EqualsAny(g.currUnit.Members) {
		return &red
	}

	if g.b.IsFilled(c) {
		return &yellow
	}

	return &grey
}

// TODO(myenik) Render Unit
// TODO(myenik)  Merge this with RenderInputProblem?
func (r *GameRenderer) AddFrame(g *Game) {
	height := r.height
	width := r.width
	border := r.border
	hs := r.hs

	board := image.Rect(border, border, hs.horiz/2+hs.horiz*width, hs.vert*height)
	rect := image.Rectangle{image.ZP, board.Max.Add(image.Pt(border, border))}

	// background
	m := image.NewPaletted(rect, pal)
	draw.Draw(m, rect, &white, image.ZP, draw.Src)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			drawHex(m, board, Cell{x, y}, hs, gameFillColor(g, x, y))
		}
	}
	drawPivot(m, board, g.currUnit.Pivot, hs)

	r.frames = append(r.frames, m)
}

func (r *GameRenderer) OutputGIF(w io.Writer, delay int) {
	delays := make([]int, len(r.frames))
	for i := range delays {
		delays[i] = delay
	}

	anim := gif.GIF{Image: r.frames, Delay: delays, LoopCount: 0}
	gif.EncodeAll(w, &anim)
}
