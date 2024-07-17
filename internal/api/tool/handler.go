// Package tool
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-05 17:22
package tool

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
	"github.com/LLiuHuan/gin-template/pkg/hashids"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// HashIdsEncode HashIds 加密
	// @Tags API.tool
	// @Router /api/v1/tool/hashids/encode/{id} [get]
	HashIdsEncode() core.HandlerFunc

	// HashIdsDecode HashIds 解密
	// @Tags API.tool
	// @Router /api/v1/tool/hashids/decode/{id} [get]
	HashIdsDecode() core.HandlerFunc

	// SearchCache 查询缓存
	// @Tags API.tool
	// @Router /api/v1/tool/cache/search [post]
	SearchCache() core.HandlerFunc

	// ClearCache 清空缓存
	// @Tags API.tool
	// @Router /api/v1/tool/cache/clear [patch]
	ClearCache() core.HandlerFunc

	// Dbs 查询 DB
	// @Tags API.tool
	// @Router /api/v1/tool/data/dbs [get]
	Dbs() core.HandlerFunc

	// Tables 查询 Table
	// @Tags API.tool
	// @Router /api/v1/tool/data/tables [post]
	Tables() core.HandlerFunc

	// SearchMySQL 执行 SQL 语句
	// @Tags API.tool
	// @Router /api/v1/tool/data/mysql [post]
	SearchMySQL() core.HandlerFunc

	// ProjectInfo 项目基础信息
	// @Tags API.tool
	// @Router /api/v1/tool/project/info [get]
	ProjectInfo() core.HandlerFunc

	//// SendMessage 发送消息
	//// @Tags API.tool
	//// @Router /api/v1/tool/send_message [post]
	//SendMessage() core.HandlerFunc
}

type handler struct {
	logger  *zap.Logger
	db      database.Repo
	cache   redis.Repo
	hashids hashids.Hash
}

func New(logger *zap.Logger, db database.Repo, cache redis.Repo) Handler {
	return &handler{
		logger: logger,
		db:     db,
		cache:  cache,
	}
}

func (h *handler) i() {}
