// Package helper
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-12 11:22
package helper

import (
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/services/helper"
	"mime/multipart"
	"net/http"
)

type uploadFileRequest struct {
	File  *multipart.FileHeader `form:"file"`
	Hash  string                `form:"hash"`
	Index int                   `form:"index"`
	Start int                   `form:"start"`
	End   int                   `form:"end"`
}

type uploadFileResponse struct {
}

// UploadFile 大文件分片上传
//
//	@Summary		大文件上传
//	@Description	大文件上传
//	@Tags			Helper
func (h *handler) UploadFile() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(uploadFileRequest)
		res := new(uploadFileResponse)

		if err := ctx.ShouldBind(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		data := new(helper.UploadFileData)
		data.File = req.File
		data.FileHash = req.Hash
		data.FileIndex = req.Index
		data.FileStart = req.Start
		data.FileEnd = req.End

		err := h.helperService.UploadFile(ctx, data)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.FileIncompleteError,
				code.Text(code.FileIncompleteError)).WithError(err),
			)
			return
		}

		//time.Sleep(1 * time.Second)

		ctx.Payload(res)
	}
}
