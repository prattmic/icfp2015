package main

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"os"
)

var (
	pal = color.Palette(palette.Plan9) // TODO(aray) smaller palette
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

func drawHex(m draw.Image, x, y, size int) {
	red := pal.Convert(color.RGBA{255, 0, 0, 255})
	hm := HexMask{image.Point{x, y}, size}
	u := image.Uniform{red}
	draw.DrawMask(m, m.Bounds(), &u, image.ZP, &hm, image.ZP, draw.Over)
}

func main() {
	var images []*image.Paletted
	var delays []int
	f, _ := os.Create("foo.gif")
	rect := image.Rect(0, 0, 28+24*4, 20+18*10)
	// not a loop anymore
	m := image.NewPaletted(rect, pal)
	black := pal.Convert(color.RGBA{0, 0, 0, 255})
	draw.Draw(m, rect, &image.Uniform{black}, image.ZP, draw.Src)
	for i := 0; i < 10; i++ {
		for j := 0; j < 4; j++ {
			x := 10 + j*24
			y := 10 + i*18
			if i%2 == 1 {
				x += 12
			}
			drawHex(m, x, y, 20)
		}
	}
	images = append(images, m)
	delays = append(delays, 10)
	// end of not a loop anymore
	anim := gif.GIF{Image: images, Delay: delays, LoopCount: 1000}
	gif.EncodeAll(f, &anim)
}
