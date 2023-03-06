package imagemask

import (
	"image"
	"image/color"
)

func NewRadiusMask(img image.Image, r int) RadiusMask {
	return RadiusMask{img, image.Point{img.Bounds().Dx(), img.Bounds().Dy()}, r}
}

type RadiusMask struct {
	image image.Image
	point image.Point // 矩形右下角位置
	r     int         // 圆角值
}

func (r RadiusMask) ColorModel() color.Model {
	return r.image.ColorModel()
}

func (r RadiusMask) Bounds() image.Rectangle {
	return image.Rect(0, 0, r.point.X, r.point.Y)
}

func (r RadiusMask) At(x, y int) color.Color {
	var xx, yy, rr float64
	var inArea bool
	// left up
	if x <= r.r && y <= r.r {
		xx, yy, rr = float64(r.r-x)+0.5, float64(y-r.r)+0.5, float64(r.r)
		inArea = true
	}
	// right up
	if x >= (r.point.X-r.r) && y <= r.r {
		xx, yy, rr = float64(x-(r.point.X-r.r))+0.5, float64(y-r.r)+0.5, float64(r.r)
		inArea = true
	}
	// left bottom
	if x <= r.r && y >= (r.point.Y-r.r) {
		xx, yy, rr = float64(r.r-x)+0.5, float64(y-(r.point.Y-r.r))+0.5, float64(r.r)
		inArea = true
	}
	// right bottom
	if x >= (r.point.X-r.r) && y >= (r.point.Y-r.r) {
		xx, yy, rr = float64(x-(r.point.X-r.r))+0.5, float64(y-(r.point.Y-r.r))+0.5, float64(r.r)
		inArea = true
	}

	if inArea && xx*xx+yy*yy >= rr*rr {
		return color.Alpha{}
	}
	return r.image.At(x, y)
}
