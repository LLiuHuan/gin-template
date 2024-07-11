package tool

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type dbsRequest struct{}

type dbsResponse struct{}

// Dbs 查询 DB
//
//	@Summary		查询 DB
//	@Description	查询 DB
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		dbsRequest	true	"请求信息"
//	@Success		200		{object}	dbsResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/data/dbs [get]
func (h *handler) Dbs() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
