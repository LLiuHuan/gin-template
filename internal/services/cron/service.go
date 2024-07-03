// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 10:41
package cron

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/cron"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/database/cron_task"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
)

var _ Service = (*service)(nil)

type Service interface {
	i()

	Create(ctx core.Context, createData *CreateCronTaskData) (id int, err error)
	Modify(ctx core.Context, id int, modifyData *ModifyCronTaskData) (err error)
	PageList(ctx core.Context, searchData *SearchData) (listData []*cron_task.CronTask, err error)
	PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error)
	UpdateUsed(ctx core.Context, id int, used int32) (err error)
	Execute(ctx core.Context, id int) (err error)
	Detail(ctx core.Context, searchOneData *SearchOneData) (info *cron_task.CronTask, err error)
}

type service struct {
	db         database.Repo
	cache      redis.Repo
	cronServer cron.Server
}

func New(db database.Repo, cache redis.Repo, cron cron.Server) Service {
	return &service{
		db:         db,
		cache:      cache,
		cronServer: cron,
	}
}

func (s *service) i() {}
