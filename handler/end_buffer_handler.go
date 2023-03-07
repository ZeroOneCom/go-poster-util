/**
 * @Author: entere@126.com
 * @Description:
 * @File:  background_handler.go
 * @Version: 1.0.0
 * @Date: 2020/5/21 12:31
 */

package handler

import (
	"bytes"
	"fmt"
	"github.com/ZeroOneCom/go-poster-util/core"
)

// EndBufferHandler 结束，写在最后，把图片合并到一张图上
type EndBufferHandler struct {
	// 合成复用Next
	Next
	Buffer *bytes.Buffer
}

// Do 地址逻辑
func (h *EndBufferHandler) Do(c *Context) (err error) {
	// 合并
	err = core.Merge(c.PngCarrier, h.Buffer)
	if err != nil {
		fmt.Errorf("core.Merge err：%v", err)
	}
	return
}
