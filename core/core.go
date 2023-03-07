/**
 * @Author: entere@126.com
 * @Description:
 * @File:  core
 * @Version: 1.0.0
 * @Date: 2020/4/22 10:41
 */

package core

import (
	"bytes"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/skip2/go-qrcode"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// 新PNG载体
type Rect struct {
	X0 int
	X1 int
	Y0 int
	Y1 int
}

// 坐标
type Pt struct {
	X int
	Y int
}

// 图片切片
type DImage struct {
	PNG draw.Image //合并到的PNG切片,可用image.NewrRGBA设置
	X   int        //横坐标
	Y   int        //纵坐标
}

// 文字切片
type DText struct {
	PNG   draw.Image //合并到的PNG切片,可用image.NewrRGBA设置
	Title string     //文字
	X     int        //横坐标
	Y     int        //纵坐标
	Size  float64
	R     uint8
	G     uint8
	B     uint8
	A     uint8
}

// 新建文件载体
func NewMerged(path string) (*os.File, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 新建图片载体
func NewPNG(X0 int, Y0 int, X1 int, Y1 int) *image.RGBA {
	return image.NewRGBA(image.Rect(X0, Y0, X1, Y1))
}

// 合并图片到载体
func MergeImage(PNG draw.Image, image image.Image, imageBound image.Point) {
	draw.Draw(PNG, PNG.Bounds(), image, imageBound, draw.Over)
}

// 读取字体类型
func LoadTextType(path string) (*truetype.Font, error) {
	fbyte, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	trueTypeFont, err := freetype.ParseFont(fbyte)
	if err != nil {
		return nil, err
	}
	return trueTypeFont, nil
}

// 创建新字体切片
func NewDrawText(png draw.Image) *DText {
	return &DText{
		PNG:  png,
		Size: 18,
		X:    0,
		Y:    0,
		R:    0,
		G:    0,
		B:    0,
		A:    255,
	}
}

// 设置字体颜色
func (dtext *DText) SetColor(R uint8, G uint8, B uint8) {
	dtext.R = R
	dtext.G = G
	dtext.B = B
}

type TextAlign int

const (
	TEXT_ALIGN_DEFAULT = iota
	TEXT_ALIGN_CENTER
)

// 计算文本实际像素
func (dtext *DText) CalcAdvanceWidth(s string, f *truetype.Font, fontHinting font.Hinting, fontSize float64, dpi float64) fixed.Int26_6 {
	scale := fixed.Int26_6(fontSize * dpi * (64.0 / 72.0))
	p := freetype.Pt(0, 0)
	prev, hasPrev := truetype.Index(0), false

	opts := truetype.Options{
		Size: fontSize,
	}
	face := truetype.NewFace(f, &opts)

	for _, w := range s {
		index := f.Index(w)
		if hasPrev {
			kern := f.Kern(scale, prev, index)
			if fontHinting != font.HintingNone {
				kern = (kern + 32) &^ 63
			}
			p.X += kern
		}

		advanceWidth, ok := face.GlyphAdvance(w)
		if ok {
			p.X += advanceWidth
		}
		prev, hasPrev = index, true
	}

	return p.X
}

// 合并字体到载体
func (dtext *DText) MergeText(title string, size float64, tf *truetype.Font, x int, y int, align TextAlign) error {
	const dpi = 72
	const fontHunting = font.HintingNone

	fc := freetype.NewContext()
	//设置屏幕每英寸的分辨率
	fc.SetDPI(dpi)
	//设置用于绘制文本的字体
	fc.SetFont(tf)
	//以磅为单位设置字体大小
	fc.SetFontSize(size)
	//设置剪裁矩形以进行绘制
	fc.SetClip(dtext.PNG.Bounds())
	//设置目标图像
	fc.SetDst(dtext.PNG)
	//设置绘制操作的源图像，通常为 image.Uniform
	fc.SetSrc(image.NewUniform(color.RGBA{dtext.R, dtext.G, dtext.B, dtext.A}))
	fc.SetHinting(fontHunting)

	pt := freetype.Pt(x, y)

	// 计算文本对齐
	if align == TEXT_ALIGN_CENTER {
		bounds := fixed.I(dtext.PNG.Bounds().Max.X) / 2
		words := dtext.CalcAdvanceWidth(title, tf, fontHunting, size, dpi) / 2
		pt.X = bounds - words
	}

	_, err := fc.DrawString(title, pt)
	if err != nil {
		return err
	}
	return nil
}

// 合并到图片
func Merge(img draw.Image, merged *os.File) error {
	err := png.Encode(merged, img)
	if err != nil {
		return err
	}
	return nil
}

// 获取二维码图像
func DrawQRImage(url string, level qrcode.RecoveryLevel, size int) (image.Image, error) {
	newQr, err := qrcode.New(url, level)
	if err != nil {
		return nil, err
	}
	qrImage := newQr.Image(size)
	return qrImage, nil
}

func GetResourceReader(url string) (r *bytes.Reader, err error) {
	if url[0:4] == "http" {
		resp, err := http.Get(url)
		if err != nil {
			return r, err
		}
		defer resp.Body.Close()
		fileBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return r, err
		}
		r = bytes.NewReader(fileBytes)
	} else {
		fileBytes, err := os.ReadFile(url)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(fileBytes)
	}
	return r, nil
}
