// Package helper
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-17 23:42
package helper

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/services/helper"
	"net/http"
)

type uploadVerifyRequest struct {
	Hash string `form:"hash" json:"hash"`
	Name string `form:"name" json:"name"`
}

type uploadVerifyResponse struct {
	ExistFile   bool  `json:"exist_file"`
	ExistChunks []int `json:"exist_chunks"`
}

// UploadVerify 大文件分片合并
//
//	@Summary		大文件上传
//	@Description	大文件上传
//	@Tags			Helper
func (h *handler) UploadVerify() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(uploadVerifyRequest)
		res := new(uploadVerifyResponse)

		if err := ctx.ShouldBind(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		data := new(helper.VerifyFileData)
		data.Hash = req.Hash
		data.Name = req.Name

		verifyFile, err := h.helperService.VerifyFile(ctx, data)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.FileIncompleteError,
				code.Text(code.FileIncompleteError)).WithError(err),
			)
			return
		}
		res.ExistChunks = verifyFile.ExistChunks
		res.ExistFile = verifyFile.ExistFile

		ctx.Payload(res)
	}
}
