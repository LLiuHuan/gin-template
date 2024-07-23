// Package admin
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-09 09:47
package admin

import (
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/internal/services/admin"
	"github.com/LLiuHuan/gin-template/pkg/hashids"
	"github.com/bytedance/godlp/dlpheader"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Login 管理员登录
	// @Tags API.admin
	// @Router /api/v1/login [post]
	Login() core.HandlerFunc

	// Logout 管理员登出
	// @Tags API.admin
	// @Router /api/v1/admin/logout [post]
	Logout() core.HandlerFunc

	// ModifyPassword 修改密码
	// @Tags API.admin
	// @Router /api/v1/admin/modify_password [patch]
	ModifyPassword() core.HandlerFunc

	// Detail 个人信息
	// @Tags API.admin
	// @Router /api/v1/admin/info [get]
	Detail() core.HandlerFunc

	// ModifyPersonalInfo 修改个人信息
	// @Tags API.admin
	// @Router /api/v1/admin/modify_personal_info [patch]
	ModifyPersonalInfo() core.HandlerFunc

	// Create 新增管理员
	// @Tags API.admin
	// @Router /api/v1/admin [post]
	Create() core.HandlerFunc

	// List 管理员列表
	// @Tags API.admin
	// @Router /api/v1/admin [get]
	List() core.HandlerFunc

	// Delete 删除管理员
	// @Tags API.admin
	// @Router /api/v1/admin/{id} [delete]
	Delete() core.HandlerFunc

	// Offline 下线管理员
	// @Tags API.admin
	// @Router /api/v1/admin/offline [patch]
	Offline() core.HandlerFunc

	// UpdateUsed 更新管理员为启用/禁用
	// @Tags API.admin
	// @Router /api/v1/admin/used [patch]
	UpdateUsed() core.HandlerFunc

	// ResetPassword 重置密码
	// @Tags API.admin
	// @Router /api/v1/admin/reset_password/{id} [patch]
	ResetPassword() core.HandlerFunc

	// CreateAdminMenu 提交菜单授权
	// @Tags API.admin
	// @Router /api/v1/admin/menu [post]
	CreateAdminMenu() core.HandlerFunc

	// ListAdminMenu 菜单授权列表
	// @Tags API.admin
	// @Router /api/v1/admin/menu/{id} [get]
	ListAdminMenu() core.HandlerFunc

	// Captcha 获取验证码
	// @Tags API.admin
	// @Router /api/v1/captcha [get]
	Captcha() core.HandlerFunc
}

type handler struct {
	logger       *zap.Logger
	cache        redis.Repo
	hashids      hashids.Hash
	adminService admin.Service
	dlp          dlpheader.EngineAPI
}

func New(logger *zap.Logger, db database.Repo, cache redis.Repo, dlp dlpheader.EngineAPI) Handler {
	return &handler{
		logger: logger,
		cache:  cache,
		hashids: hashids.New(
			hashids.WithAlphabet(configs.Get().HashIds.Alphabet),
			hashids.WithMinLength(configs.Get().HashIds.MinLength),
			hashids.WithBlockList(configs.Get().HashIds.BlockList),
		),
		adminService: admin.New(db, cache, logger),
		dlp:          dlp,
	}
}

func (h *handler) i() {}
