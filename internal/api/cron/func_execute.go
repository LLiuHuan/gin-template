package cron

import (
	"net/http"

	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/pkg/validation"

	"github.com/spf13/cast"
)

type executeRequest struct {
	Id string `uri:"id"` // HashID
}

type executeResponse struct {
	Id int `json:"id"` // ID
}

// Execute 手动执行任务
//
//	@Summary		手动执行任务
//	@Description	手动执行任务
//	@Tags			API.cron
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		executeRequest	true	"请求信息"
//	@Success		200		{object}	executeResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/cron/exec/{id} [patch]
func (h *handler) Execute() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(executeRequest)
		res := new(executeResponse)
		if err := ctx.ShouldBindURI(req); err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)).WithError(err),
			)
			return
		}

		ids, err := h.hashids.HashidsDecode(req.Id)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.HashIdsDecodeError,
				code.Text(code.HashIdsDecodeError)).WithError(err),
			)
			return
		}

		err = h.cronService.Execute(ctx, cast.ToInt(ids[0]))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.CronExecuteError,
				code.Text(code.CronExecuteError)).WithError(err),
			)
			return
		}

		res.Id = cast.ToInt(ids[0])
		ctx.Payload(res)
	}
}
