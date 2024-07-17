package admin

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/pkg/password"
	"github.com/LLiuHuan/gin-template/internal/pkg/validation"
	"github.com/LLiuHuan/gin-template/internal/proposal"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/internal/services/admin"
	"github.com/LLiuHuan/gin-template/pkg/errors"
)

type loginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`              // 用户名
	Password string `form:"password" json:"password" binding:"min=4,max=16,required"` // 密码
}

type loginResponse struct {
	Token        string `json:"token"`        // 用户身份标识
	RefreshToken string `json:"refreshToken"` // 刷新token
}

// Login 管理员登录
//
//	@Summary		管理员登录
//	@Description	管理员登录
//	@Tags			API.admin
//	@Accept			application/x-www-form-urlencoded
//	@Produce		json
//	@Param			Request	body		loginRequest	true	"请求信息"
//	@Success		200		{object}	loginResponse
//	@Failure		400		{object}	code.Failure
//	@Router			/api/v1/login [post]
//
// @Security LoginToken
func (h *handler) Login() core.HandlerFunc {
	return func(ctx core.Context) {
		req := new(loginRequest)
		res := new(loginResponse)
		if err := ctx.ShouldBindJSON(req); err != nil {
			fmt.Println(err)
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).
				WithError(validation.ErrorE(err)),
			)
			return
		}

		searchOneData := new(admin.SearchOneData)
		searchOneData.Username = req.Username
		searchOneData.Password = password.GeneratePassword(req.Password)
		searchOneData.IsUsed = 1

		info, err := h.adminService.Detail(ctx, searchOneData)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		if info == nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(errors.New("未查询出符合条件的用户")),
			)
			return
		}

		token := password.GenerateLoginToken(info.Id)

		// 用户信息
		sessionUserInfo := &proposal.SessionUserInfo{
			UserID:   info.Id,
			UserName: info.Username,
		}

		// 将用户信息记录到 Redis 中
		err = h.cache.Set(configs.RedisKeyPrefixLoginUser+token, string(sessionUserInfo.Marshal()), configs.LoginSessionTTL, redis.WithTrace(ctx.Trace()))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		searchMenuData := new(admin.SearchMyMenuData)
		searchMenuData.AdminId = info.Id
		menu, err := h.adminService.MyMenu(ctx, searchMenuData)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		// 菜单栏信息
		menuJsonInfo, _ := json.Marshal(menu)

		// 将菜单栏信息记录到 Redis 中
		err = h.cache.Set(configs.RedisKeyPrefixLoginUser+token+":menu", string(menuJsonInfo), configs.LoginSessionTTL, redis.WithTrace(ctx.Trace()))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		searchActionData := new(admin.SearchMyActionData)
		searchActionData.AdminId = info.Id
		action, err := h.adminService.MyAction(ctx, searchActionData)
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		// 可访问接口信息
		actionJsonInfo, _ := json.Marshal(action)

		// 将可访问接口信息记录到 Redis 中
		err = h.cache.Set(configs.RedisKeyPrefixLoginUser+token+":action", string(actionJsonInfo), configs.LoginSessionTTL, redis.WithTrace(ctx.Trace()))
		if err != nil {
			ctx.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		res.Token = token
		ctx.Payload(res)
	}
}
