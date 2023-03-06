/**
 * @Author: entere@126.com
 * @Description:
 * @File:  演示海报生成过程
 * @Version: 1.0.0
 * @Date: 2020/5/21 21:54
 */

package main

import (
	"fmt"
	"github.com/ZeroOneCom/go-poster-util/core"
	"github.com/ZeroOneCom/go-poster-util/handler"
	"github.com/ZeroOneCom/go-poster-util/imagemask"
	"github.com/rs/xid"
	"image"
)

// 海报组件使用示例
func main() {

	nullHandler := &handler.NullHandler{}
	ctx := &handler.Context{
		//图片都绘在这个PNG载体上
		PngCarrier: core.NewPNG(0, 0, 750, 1334),
	}

	//绘制文字
	textHandler1 := &handler.TextHandler{
		Next:     handler.Next{},
		X:        0,
		Y:        105,
		Size:     20,
		R:        255,
		G:        241,
		B:        250,
		Align:    core.TEXT_ALIGN_CENTER,
		Text:     "StrAlign takes two你好大家好",
		FontPath: "./assets/msyh.ttf",
	}

	//结束绘制，把前面的内容合并成一张图片
	endHandler := &handler.EndHandler{
		Output: "./build/poster_" + xid.New().String() + ".png",
	}

	// 链式调用绘制过程
	nullHandler.
		SetNext(textHandler1).
		SetNext(&handler.ImageHandler{
			ImageSource: &core.ImageSource{
				Path:   "https://image.lingyi360.com/2023/3/5/620fccb0efaca36ae7cb9c7158b6e9d1.jpg",
				Resize: &image.Point{X: 650},
				Handler: func(img image.Image) image.Image {
					//return imagemask.NewRadiusMask(img, 40)
					return imagemask.NewCircleMask(img)
				},
			},
			X: 30,
			Y: 50,
		}).
		SetNext(endHandler)

	// 开始执行业务
	if err := nullHandler.Run(ctx); err != nil {
		// 异常
		fmt.Println("Fail | Error:" + err.Error())
		return
	}
	// 成功
	fmt.Println("Success")
	return

}
