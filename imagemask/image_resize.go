package imagemask

import (
	"github.com/nfnt/resize"
	"image"
)

func ImageResize(img image.Image, w int, h int) image.Image {
	return resize.Resize(uint(w), uint(h), img, resize.Lanczos3)
}
