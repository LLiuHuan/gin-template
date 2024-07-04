// Package install
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-04 18:06
package install

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Install 安装
	// @Tags API.install
	// @Router /v1/api/install [post]
	Install() core.HandlerFunc
}

type handler struct {
	logger *zap.Logger
	cache  redis.Repo
	db     database.Repo
}

func New(logger *zap.Logger, db database.Repo, cache redis.Repo) Handler {
	return &handler{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

func (h *handler) i() {}
