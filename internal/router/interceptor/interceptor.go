// Package interceptor
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 14:59
package interceptor

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/proposal"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/internal/services/authorized"

	"go.uber.org/zap"
)

var _ Interceptor = (*interceptor)(nil)

type Interceptor interface {
	// CheckLogin 验证是否登录
	CheckLogin(ctx core.Context) (info proposal.SessionUserInfo, err core.BusinessError)

	// CheckRBAC 验证 RBAC 权限是否合法
	CheckRBAC() core.HandlerFunc

	// CheckSignature 验证签名是否合法，对用签名算法 pkg/signature
	CheckSignature() core.HandlerFunc

	// i 为了避免被其他包实现
	i()
}

type interceptor struct {
	logger            *zap.Logger
	cache             redis.Repo
	db                database.Repo
	authorizedService authorized.Service
	//adminService      admin.Service
}

func New(logger *zap.Logger, cache redis.Repo, db database.Repo) Interceptor {
	return &interceptor{
		logger:            logger,
		cache:             cache,
		db:                db,
		authorizedService: authorized.New(db, cache),
		//adminService:      admin.New(db, cache),
	}
}

func (i *interceptor) i() {}
