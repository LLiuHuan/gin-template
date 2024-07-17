// Package helper
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-17 21:20
package helper

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/services/helper"
	"net/http"
)

type uploadMergeRequest struct {
	Hash  string `form:"hash" json:"hash"`
	Name  string `form:"name" json:"name"`
	Count int    `form:"count" json:"count"`
}

type uploadMergeResponse struct {
}

// UploadMerge 大文件分片合并
//
//	@Summary		大文件上传
//	@Description	大文件上传
//	@Tags			Helper
func (h *handler) UploadMerge() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(uploadMergeRequest)
		res := new(uploadMergeResponse)

		if err := ctx.ShouldBind(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		data := new(helper.MergeFileData)
		data.Hash = req.Hash
		data.Name = req.Name
		data.Count = req.Count

		err := h.helperService.MergeFile(ctx, data)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.FileIncompleteError,
				code.Text(code.FileIncompleteError)).WithError(err),
			)
			return
		}

		ctx.Payload(res)

	}
}
