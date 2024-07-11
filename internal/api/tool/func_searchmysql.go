package tool

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type searchMySQLRequest struct{}

type searchMySQLResponse struct{}

// SearchMySQL 执行 SQL 语句
//
//	@Summary		执行 SQL 语句
//	@Description	执行 SQL 语句
//	@Tags			API.tool
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		searchMySQLRequest	true	"请求信息"
//	@Success		200		{object}	searchMySQLResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/tool/data/mysql [post]
func (h *handler) SearchMySQL() core.HandlerFunc {
	return func(ctx core.Context) {

	}
}
