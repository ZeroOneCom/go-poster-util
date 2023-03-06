package core

import (
	"github.com/ZeroOneCom/go-poster-util/imagemask"
	"image"
)

type ImageSourceType int
type ImageHandler = func(image.Image) image.Image

type IImageSource interface {
	ReadImage() (image.Image, error)
}

type ImageSource struct {
	Path    string
	Handler ImageHandler
	Resize  *image.Point
}

func (s *ImageSource) ReadImage() (image.Image, error) {
	srcReader, err := GetResourceReader(s.Path)
	if err != nil {
		return nil, err
	}
	srcImage, _, err := image.Decode(srcReader)
	if err != nil {
		return nil, err
	}

	if s.Handler != nil {
		srcImage = s.Handler(srcImage)
	}

	if s.Resize != nil {
		w := s.Resize.X
		h := s.Resize.Y
		// 如果设置了宽没有设置高 自适应高度
		if w > 0 && h == 0 {
			h = w / srcImage.Bounds().Size().X * srcImage.Bounds().Size().Y
		}
		// 如果设置了高没有设置宽 自适应宽度
		if w == 0 && h > 0 {
			w = h / srcImage.Bounds().Size().Y * srcImage.Bounds().Size().X
		}

		srcImage = imagemask.ImageResize(srcImage, w, h)
	}

	return srcImage, nil
}
