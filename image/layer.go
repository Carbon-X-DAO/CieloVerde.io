package image

import (
	"image"
	"image/draw"
)

func Layer(bottom image.Image, top image.Image, x, y int) image.Image {
	dims := bottom.Bounds()
	offset := image.Pt(x, y+dims.Dy()-top.Bounds().Dy())
	res := image.NewRGBA(dims)
	draw.Draw(res, dims, bottom, image.Point{}, draw.Src)
	draw.Draw(res, top.Bounds().Add(offset), top, image.Point{}, draw.Over)

	return res
}
