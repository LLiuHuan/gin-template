// Package helper
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 15:33
package helper

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/internal/services/authorized"
	"github.com/LLiuHuan/gin-template/internal/services/helper"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Md5 加密
	// @Tags Helper
	// @Router /helper/md5/{str} [get]
	Md5() core.HandlerFunc

	// Sign 签名
	// @Tags Helper
	// @Router /helper/sign [post]
	Sign() core.HandlerFunc

	// UploadFile 大文件分片上传
	// @Summary 大文件上传
	// @Description 大文件上传
	// @Tags Helper
	// @Router /api/v1/helper/upload [post]
	UploadFile() core.HandlerFunc

	// UploadMerge 大文件分片合并
	// @Summary 大文件上传
	// @Description 大文件上传
	// @Tags Helper
	// @Router /api/v1/helper/merge [post]
	UploadMerge() core.HandlerFunc

	// UploadVerify 大文件分片校验
	// @Summary 大文件上传
	// @Description 大文件上传
	// @Tags Helper
	// @Router /api/v1/helper/verify [post]
	UploadVerify() core.HandlerFunc
}

type handler struct {
	logger            *zap.Logger
	db                database.Repo
	authorizedService authorized.Service
	helperService     helper.Service
}

func New(logger *zap.Logger, db database.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:            logger,
		db:                db,
		authorizedService: authorized.New(db, cache),
		helperService:     helper.New(db, cache),
	}
}

func (h *handler) i() {}
