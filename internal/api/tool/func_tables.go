package tool

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type tablesRequest struct{}

type tablesResponse struct{}

// Tables 查询 Table
//
//	@Summary		查询 Table
//	@Description	查询 Table
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		tablesRequest	true	"请求信息"
//	@Success		200		{object}	tablesResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/data/tables [post]
func (h *handler) Tables() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
