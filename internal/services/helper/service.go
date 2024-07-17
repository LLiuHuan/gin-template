// Package helper
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-12 14:31
package helper

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
)

const (
	// DefaultBufferSize 默认缓冲区大小
	DefaultBufferSize = 1 << 20 // 1MB
	// DefaultBreakpointPath 默认断点续传文件路径
	DefaultBreakpointPath = "./breakpoint"
	// DefaultFileSavePath 默认文件保存路径
	DefaultFileSavePath = "./upload"
)

var _ Service = (*service)(nil)

type Service interface {
	i()

	UploadFile(ctx core.Context, data *UploadFileData) error
	MergeFile(ctx core.Context, data *MergeFileData) error
	VerifyFile(ctx core.Context, data *VerifyFileData) (VerifyFileResponse, error)
}

type service struct {
	db    database.Repo
	cache redis.Repo
}

func New(db database.Repo, cache redis.Repo) Service {
	return &service{
		db:    db,
		cache: cache,
	}
}

func (s *service) i() {}
