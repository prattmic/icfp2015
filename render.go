package main

import (
	//"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"io"
	"math"
)

var (
	pal    = color.Palette(palette.Plan9) // TODO(aray) smaller palette
	white  = pal.Convert(color.RGBA{255, 255, 255, 255})
	yellow = pal.Convert(color.RGBA{255, 255, 0, 255})   // Filled
	green  = pal.Convert(color.RGBA{0, 255, 0, 255})     // Unit
	grey   = pal.Convert(color.RGBA{200, 200, 200, 255}) // Empty
)

type HexMask struct {
	p    image.Point // top left of bounding rectangle
	size int         // size in pixels (dimension in both width and height)
}

func (hm *HexMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (hm *HexMask) Bounds() image.Rectangle {
	return image.Rect(hm.p.X, hm.p.Y, hm.p.X+hm.size, hm.p.Y+hm.size)
}

func (hm *HexMask) At(xx, yy int) color.Color {
	x := float64(xx - hm.p.X)
	y := float64(yy - hm.p.Y)
	k := float64(hm.size)
	h := float64(hm.size-1) / 2.0
	if (k - 2*math.Abs(y-h)) <= math.Abs(x-h) {
		return color.Alpha{0}
	}
	return color.Alpha{255}
}

func drawHex(m draw.Image, x, y, size int, c color.Color) {
	hm := HexMask{image.Point{x, y}, size}
	u := image.Uniform{c}
	draw.DrawMask(m, m.Bounds(), &u, image.ZP, &hm, image.ZP, draw.Over)
}

func fillColor(ip *InputProblem, x, y int) color.Color {
	for _, cell := range ip.Filled {
		//fmt.Printf("cell.X %d cell.Y %d x %d y %d", cell.X, cell.Y, x, y)
		if cell.X == x && cell.Y == y {
			return yellow
		}
	}
	return grey
}

// This takes an io.Writer, and renders a GIF of the InputProblem to it
func RenderInputProblem(w io.Writer, ip *InputProblem) {
	height := ip.Height
	width := ip.Width
	//fmt.Printf("filling %d cells\n", len(ip.Filled))
	//fmt.Println(ip.Filled)
	rect := image.Rect(0, 0, 28+24*width, 20+18*height)

	// background
	m := image.NewPaletted(rect, pal)
	draw.Draw(m, rect, &image.Uniform{white}, image.ZP, draw.Src)

	// TODO(aray): dont hardcode pixels based on 20px hexagons
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			x := 10 + j*24
			y := 10 + i*18
			if i%2 == 1 {
				x += 12
			}
			drawHex(m, x, y, 20, fillColor(ip, j, i))
		}
	}

	// end of not a loop anymore
	image := []*image.Paletted{m}
	anim := gif.GIF{Image: image, Delay: []int{10}, LoopCount: 1}
	gif.EncodeAll(w, &anim)
}
