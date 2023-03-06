package handler

import (
	"fmt"
	"github.com/ZeroOneCom/go-poster-util/core"
	"image"
)

// ImageHandler 根据URL地址设置图片
type ImageHandler struct {
	// 合成复用Next
	Next
	X           int
	Y           int
	ImageSource core.IImageSource
}

// Do 地址逻辑
func (h *ImageHandler) Do(c *Context) (err error) {
	srcImage, err := h.ImageSource.ReadImage()
	if err != nil {
		fmt.Errorf("图片获取失败：%v", err)
	}
	srcPoint := image.Point{
		X: h.X,
		Y: h.Y,
	}
	core.MergeImage(c.PngCarrier, srcImage, srcImage.Bounds().Min.Sub(srcPoint))
	return
}
