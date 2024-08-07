// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 10:39
package cron

import (
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	cronRepo "github.com/LLiuHuan/gin-template/internal/repository/cron"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/internal/services/cron"
	"github.com/LLiuHuan/gin-template/pkg/hashids"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Create 创建任务
	// @Tags API.cron
	// @Router /api/cron [post]
	Create() core.HandlerFunc

	// Modify 编辑任务
	// @Tags API.cron
	// @Router /api/cron/{id} [post]
	Modify() core.HandlerFunc

	// List 任务列表
	// @Tags API.cron
	// @Router /api/cron [get]
	List() core.HandlerFunc

	// UpdateUsed 更新任务为启用/禁用
	// @Tags API.cron
	// @Router /api/cron/used [patch]
	UpdateUsed() core.HandlerFunc

	// Detail 获取单条任务详情
	// @Tags API.cron
	// @Router /api/cron/{id} [get]
	Detail() core.HandlerFunc

	// Execute 手动执行任务
	// @Tags API.cron
	// @Router /api/cron/exec/{id} [patch]
	Execute() core.HandlerFunc
}

type handler struct {
	logger      *zap.Logger
	cache       redis.Repo
	hashids     hashids.Hash
	cronService cron.Service
}

func New(logger *zap.Logger, db database.Repo, cache redis.Repo, cronServer cronRepo.Server) Handler {
	return &handler{
		logger: logger,
		cache:  cache,
		hashids: hashids.New(
			hashids.WithAlphabet(configs.Get().HashIds.Alphabet),
			hashids.WithMinLength(configs.Get().HashIds.MinLength),
			hashids.WithBlockList(configs.Get().HashIds.BlockList),
		),
		cronService: cron.New(db, cache, cronServer),
	}
}

func (h *handler) i() {}
