// Package authorized
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 15:03
package authorized

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/internal/repository/database"
	"github.com/LLiuHuan/gin-template/internal/repository/redis"
)

var _ Service = (*service)(nil)

type Service interface {
	i()

	//Create(ctx core.Context, authorizedData *CreateAuthorizedData) (id int32, err error)
	//List(ctx core.Context, searchData *SearchData) (listData []*authorized.Authorized, err error)
	//PageList(ctx core.Context, searchData *SearchData) (listData []*authorized.Authorized, err error)
	//PageListCount(ctx core.Context, searchData *SearchData) (total int64, err error)
	//UpdateUsed(ctx core.Context, id int32, used int32) (err error)
	//Delete(ctx core.Context, id int32) (err error)
	//Detail(ctx core.Context, id int32) (info *authorized.Authorized, err error)
	DetailByKey(ctx core.Context, key string) (data *CacheAuthorizedData, err error)

	//CreateAPI(ctx core.Context, authorizedAPIData *CreateAuthorizedAPIData) (id int32, err error)
	//ListAPI(ctx core.Context, searchAPIData *SearchAPIData) (listData []*authorized_api.AuthorizedApi, err error)
	//DeleteAPI(ctx core.Context, id int32) (err error)
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
