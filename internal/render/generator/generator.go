// Package generator
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-22 03:13
package generator

import (
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"

	"go.uber.org/zap"
)

type handler struct {
	db     database.Repo
	logger *zap.Logger
	cache  redis.Repo
}

func New(logger *zap.Logger, db database.Repo, cache redis.Repo) *handler {
	return &handler{
		logger: logger,
		cache:  cache,
		db:     db,
	}
}
