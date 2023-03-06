/**
 * @Author: entere@126.com
 * @Description: 圆型图片遮罩
 * @File:  circle_mask
 * @Version: 1.0.0
 * @Date: 2020/4/22 14:50
 */

package imagemask

import (
	"image"
	"image/color"
	"math"
)

func NewCircleMask(img image.Image) CircleMask {
	width := img.Bounds().Max.X - img.Bounds().Min.X
	height := img.Bounds().Max.Y - img.Bounds().Min.Y

	diameter := width
	if width > height {
		diameter = height
	}

	return CircleMask{img, image.Point{0, 0}, diameter}
}

type CircleMask struct {
	image    image.Image
	point    image.Point // 圆心位置
	diameter int
}

func (ci CircleMask) ColorModel() color.Model {
	return ci.image.ColorModel()
}

func (ci CircleMask) Bounds() image.Rectangle {
	return image.Rect(0, 0, ci.diameter, ci.diameter)
}

func (ci CircleMask) At(x, y int) color.Color {
	d := ci.diameter
	dis := math.Sqrt(math.Pow(float64(x-d/2), 2) + math.Pow(float64(y-d/2), 2))
	if dis > float64(d)/2 {
		return color.Alpha{}
	} else {
		return ci.image.At(ci.point.X+x, ci.point.Y+y)
	}
}
