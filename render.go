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
	hexsize   = 5    // Hexagon size in pixels
	hexhoriz  = 6    // Horizontal spacing of hexagons
	hexvert   = 4    // Vertical spacing of hexagons
	hexborder = 40   // Border pixels to edge of rectangle
	hexpivot  = 2    // Pivot size in pixels
	c_empty   = iota // keep this one first, we rely on it being 0
	c_fill
	c_unit
)

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
)

type HexMask struct {
	cell Cell // Game board coordinate
}

func (hm *HexMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (hm *HexMask) origin() image.Point {
	px := hm.cell.X*hexhoriz + (hm.cell.Y%2)*hexhoriz/2
	py := hm.cell.Y * hexvert
	return image.Point{px, py}
}

func (hm *HexMask) Bounds() image.Rectangle {
	return image.Rect(0, 0, hexsize, hexsize)
}

func (hm *HexMask) At(xx, yy int) color.Color {
	x := float64(xx)
	y := float64(yy)
	k := float64(hexsize)
	h := float64(hexsize-1) / 2.0
	if (k - 2*math.Abs(y-h)) <= math.Abs(x-h) {
		return color.Alpha{0}
	}
	return color.Alpha{255}
}

func drawHex(m draw.Image, board image.Rectangle, cell Cell, i image.Image) {
	hm := HexMask{cell}
	r := hm.Bounds().Add(board.Min).Add(hm.origin())
	draw.DrawMask(m, r, i, image.ZP, &hm, image.ZP, draw.Over)
}

func drawPivot(m draw.Image, board image.Rectangle, cell Cell) {
	h := (hexsize - hexpivot) / 2
	x := cell.X*hexhoriz + (cell.Y%2)*hexhoriz/2 + h
	y := cell.Y*hexvert + h
	min := board.Min.Add(image.Point{x, y})
	max := min.Add(image.Point{hexpivot, hexpivot})
	draw.Draw(m, image.Rectangle{min, max}, &black, image.ZP, draw.Over)
}

func fillColor(c int) image.Image {
	switch c {
	case c_fill:
		return &yellow
	case c_unit:
		return &red
	}
	return &grey
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
func RenderInputProblem(w io.Writer, ip *InputProblem, border int) {
	height := ip.Height
	width := ip.Width
	memo := memoizeCells(ip.Filled, nil, width, height)

	board := image.Rect(border, border, hexhoriz/2+hexhoriz*width, hexvert*height)
	rect := image.Rectangle{image.ZP, board.Max.Add(image.Pt(border, border))}

	// background
	m := image.NewPaletted(rect, pal)
	draw.Draw(m, rect, &white, image.ZP, draw.Src)

	// TODO(aray): dont hardcode pixels based on 20px hexagons
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := memo[x*width+y]
			drawHex(m, board, Cell{x, y}, fillColor(c))
		}
	}

	// end of not a loop anymore
	image := []*image.Paletted{m}
	anim := gif.GIF{Image: image, Delay: []int{10}, LoopCount: 1}
	gif.EncodeAll(w, &anim)
}

type GameRenderer struct {
	width, height, border int
	frames                []*image.Paletted
}

func NewGameRenderer(g *Game, border int) *GameRenderer {
	return &GameRenderer{width: g.b.width, height: g.b.height, border: border}
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

	board := image.Rect(border, border, hexhoriz/2+hexhoriz*width, hexvert*height)
	rect := image.Rectangle{image.ZP, board.Max.Add(image.Pt(border, border))}

	// background
	m := image.NewPaletted(rect, pal)
	draw.Draw(m, rect, &white, image.ZP, draw.Src)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			drawHex(m, board, Cell{x, y}, gameFillColor(g, x, y))
		}
	}
	drawPivot(m, board, g.currUnit.Pivot)

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
