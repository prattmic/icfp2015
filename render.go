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

const (
	hexsize   = 20   // Hexagon size in pixels
	hexhoriz  = 24   // Horizontal spacing of hexagons
	hexvert   = 18   // Vertical spacing of hexagons
	hexborder = 20   // Border pixels to edge of rectangle
	hexpivot  = 8    // Pivot size in pixels
	c_empty   = iota // keep this one first, we rely on it being 0
	c_fill
	c_unit
)

var (
	white  = color.RGBA{255, 255, 255, 255} // Background
	yellow = color.RGBA{255, 255, 0, 255}   // Filled
	red    = color.RGBA{255, 0, 0, 255}     // Unit
	grey   = color.RGBA{200, 200, 200, 255} // Empty
	black  = color.RGBA{0, 0, 0, 255}       // Pivot
	pal    = []color.Color{white, yellow, red, grey, black}
)

type HexMask struct {
	cell Cell // Game board coordinate
}

func (hm *HexMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (hm *HexMask) topLeft() image.Point {
	px := hexborder + hm.cell.X*hexhoriz + (hm.cell.Y%2)*hexhoriz/2
	py := hexborder + hm.cell.Y*hexvert
	return image.Point{px, py}
}

func (hm *HexMask) Bounds() image.Rectangle {
	p := hm.topLeft()
	return image.Rect(p.X, p.Y, p.X+hexsize, p.Y+hexsize)
}

func (hm *HexMask) At(xx, yy int) color.Color {
	p := hm.topLeft()
	x := float64(xx - p.X)
	y := float64(yy - p.Y)
	k := float64(hexsize)
	h := float64(hexsize-1) / 2.0
	if (k - 2*math.Abs(y-h)) <= math.Abs(x-h) {
		return color.Alpha{0}
	}
	return color.Alpha{255}
}

func drawHex(m draw.Image, x, y int, c color.Color) {
	hm := HexMask{Cell{x, y}}
	u := image.Uniform{c}
	draw.DrawMask(m, m.Bounds(), &u, image.ZP, &hm, image.ZP, draw.Over)
}

func drawPivot(m draw.Image, cell Cell) {
	h := (hexsize - hexpivot) / 2
	x := hexborder + cell.X*hexhoriz + (cell.Y%2)*hexhoriz/2 + h
	y := hexborder + cell.Y*hexvert + h
	r := image.Rect(x, y, x+hexpivot, y+hexpivot)
	u := image.Uniform{black}
	draw.Draw(m, r, &u, image.ZP, draw.Over)
}

func fillColor(c int) color.Color {
	switch c {
	case c_fill:
		return yellow
	case c_unit:
		return red
	}
	return grey
}

func memoizeCells(fill, unit []Cell, width, height int) []int {
	memoized := make([]int, width*height) // start out all c_empty
	for _, cell := range fill {
		memoized[cell.X*width+cell.Y] = c_fill
	}
	for _, cell := range unit {
		memoized[cell.X*width+cell.Y] = c_unit
	}
	return memoized
}

// This takes an io.Writer, and renders a GIF of the InputProblem to it
func RenderInputProblem(w io.Writer, ip *InputProblem) {
	height := ip.Height
	width := ip.Width
	memo := memoizeCells(ip.Filled, nil, width, height)

	// background
	rect := image.Rect(0, 0, 2*hexborder+hexhoriz/2+hexhoriz*width, 2*hexborder+hexvert*height)
	m := image.NewPaletted(rect, pal)
	draw.Draw(m, rect, &image.Uniform{white}, image.ZP, draw.Src)

	// TODO(aray): dont hardcode pixels based on 20px hexagons
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := memo[x*width+y]
			drawHex(m, x, y, fillColor(c))
		}
	}

	// end of not a loop anymore
	image := []*image.Paletted{m}
	anim := gif.GIF{Image: image, Delay: []int{10}, LoopCount: 1}
	gif.EncodeAll(w, &anim)
}

type GameRenderer struct {
	frames []*image.Paletted
}

func NewGameRenderer() *GameRenderer {
	return &GameRenderer{}
}

func gameFillColor(g *Game, x, y int) color.Color {
	c := Cell{x, y}

	if c.EqualsAny(g.currUnit.Members) {
		return red
	}

	if g.b.IsFilled(c) {
		return yellow
	}

	return grey
}

// TODO(myenik) Render Unit
// TODO(myenik)  Merge this with RenderInputProblem?
func (r *GameRenderer) AddFrame(g *Game) {
	height := g.b.height
	width := g.b.width

	rect := image.Rect(0, 0, 2*hexborder+hexhoriz/2+hexhoriz*width, 2*hexborder+hexvert*height)

	// background
	m := image.NewPaletted(rect, pal)
	draw.Draw(m, rect, &image.Uniform{white}, image.ZP, draw.Src)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			drawHex(m, x, y, gameFillColor(g, x, y))
		}
	}
	drawPivot(m, g.currUnit.Pivot)

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
