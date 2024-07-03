// Package cron
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 11:30
package cron

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/database/cron_task"
)

type SearchData struct {
	Page     int    // 第几页
	PageSize int    // 每页显示条数
	Name     string // 任务名称
	Protocol int32  // 执行方式
	IsUsed   int32  // 是否启用
}

func (s *service) PageList(ctx core.Context, searchData *SearchData) (listData []*cron_task.CronTask, err error) {
	page := searchData.Page
	if page == 0 {
		page = 1
	}

	pageSize := searchData.PageSize
	if pageSize == 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	qb := cron_task.NewQueryBuilder()

	if searchData.Name != "" {
		qb.WhereName(database.EqualPredicate, searchData.Name)
	}

	if searchData.Protocol != 0 {
		qb.WhereProtocol(database.EqualPredicate, searchData.Protocol)
	}

	if searchData.IsUsed != 0 {
		qb.WhereIsUsed(database.EqualPredicate, searchData.IsUsed)
	}

	listData, err = qb.
		Limit(pageSize).
		Offset(offset).
		OrderById(false).
		QueryAll(s.db.GetDB().WithContext(ctx.RequestContext()))
	if err != nil {
		return nil, err
	}

	return
}
