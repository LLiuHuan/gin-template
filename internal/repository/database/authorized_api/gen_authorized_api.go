///////////////////////////////////////////////////////////
// THIS FILE IS AUTO GENERATED by gormgen, DON'T EDIT IT //
//        ANY CHANGES DONE HERE WILL BE LOST             //
///////////////////////////////////////////////////////////

package authorized_api

import (
	"fmt"
	"time"

	"github.com/LLiuHuan/gin-template/internal/repository/database"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func NewModel() *AuthorizedApi {
	return new(AuthorizedApi)
}

func NewQueryBuilder() *authorizedApiQueryBuilder {
	return new(authorizedApiQueryBuilder)
}

func (t *AuthorizedApi) Create(db *gorm.DB) (id int, err error) {
	if err = db.Create(t).Error; err != nil {
		return 0, errors.Wrap(err, "create err")
	}
	return t.Id, nil
}

type authorizedApiQueryBuilder struct {
	order []string
	where []struct {
		prefix string
		value  interface{}
	}
	limit  int
	offset int
}

func (qb *authorizedApiQueryBuilder) buildQuery(db *gorm.DB) *gorm.DB {
	ret := db
	for _, where := range qb.where {
		ret = ret.Where(where.prefix, where.value)
	}
	for _, order := range qb.order {
		ret = ret.Order(order)
	}
	ret = ret.Limit(qb.limit).Offset(qb.offset)
	return ret
}

func (qb *authorizedApiQueryBuilder) Updates(db *gorm.DB, m map[string]interface{}) (err error) {
	db = db.Model(&AuthorizedApi{})

	for _, where := range qb.where {
		db.Where(where.prefix, where.value)
	}

	if err = db.Updates(m).Error; err != nil {
		return errors.Wrap(err, "updates err")
	}
	return nil
}

func (qb *authorizedApiQueryBuilder) Delete(db *gorm.DB) (err error) {
	for _, where := range qb.where {
		db = db.Where(where.prefix, where.value)
	}

	if err = db.Delete(&AuthorizedApi{}).Error; err != nil {
		return errors.Wrap(err, "delete err")
	}
	return nil
}

func (qb *authorizedApiQueryBuilder) Count(db *gorm.DB) (int64, error) {
	var c int64
	res := qb.buildQuery(db).Model(&AuthorizedApi{}).Count(&c)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		c = 0
	}
	return c, res.Error
}

func (qb *authorizedApiQueryBuilder) First(db *gorm.DB) (*AuthorizedApi, error) {
	ret := &AuthorizedApi{}
	res := qb.buildQuery(db).First(ret)
	if res.Error != nil && res.Error == gorm.ErrRecordNotFound {
		ret = nil
	}
	return ret, res.Error
}

func (qb *authorizedApiQueryBuilder) QueryOne(db *gorm.DB) (*AuthorizedApi, error) {
	qb.limit = 1
	ret, err := qb.QueryAll(db)
	if len(ret) > 0 {
		return ret[0], err
	}
	return nil, err
}

func (qb *authorizedApiQueryBuilder) QueryAll(db *gorm.DB) ([]*AuthorizedApi, error) {
	var ret []*AuthorizedApi
	err := qb.buildQuery(db).Find(&ret).Error
	return ret, err
}

func (qb *authorizedApiQueryBuilder) Limit(limit int) *authorizedApiQueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *authorizedApiQueryBuilder) Offset(offset int) *authorizedApiQueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereId(p database.Predicate, value int) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereIdIn(value []int) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereIdNotIn(value []int) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "id", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderById(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "id "+order)
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereBusinessKey(p database.Predicate, value string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "business_key", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereBusinessKeyIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "business_key", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereBusinessKeyNotIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "business_key", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderByBusinessKey(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "business_key "+order)
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereMethod(p database.Predicate, value string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "method", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereMethodIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "method", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereMethodNotIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "method", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderByMethod(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "method "+order)
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereApi(p database.Predicate, value string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "api", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereApiIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "api", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereApiNotIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "api", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderByApi(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "api "+order)
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereIsDeleted(p database.Predicate, value int32) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_deleted", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereIsDeletedIn(value []int32) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_deleted", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereIsDeletedNotIn(value []int32) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "is_deleted", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderByIsDeleted(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "is_deleted "+order)
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereCreatedAt(p database.Predicate, value time.Time) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_at", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereCreatedAtIn(value []time.Time) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_at", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereCreatedAtNotIn(value []time.Time) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_at", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderByCreatedAt(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "created_at "+order)
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereCreatedUser(p database.Predicate, value string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_user", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereCreatedUserIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_user", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereCreatedUserNotIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "created_user", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderByCreatedUser(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "created_user "+order)
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereUpdatedAt(p database.Predicate, value time.Time) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_at", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereUpdatedAtIn(value []time.Time) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_at", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereUpdatedAtNotIn(value []time.Time) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_at", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderByUpdatedAt(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "updated_at "+order)
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereUpdatedUser(p database.Predicate, value string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_user", p),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereUpdatedUserIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_user", "IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) WhereUpdatedUserNotIn(value []string) *authorizedApiQueryBuilder {
	qb.where = append(qb.where, struct {
		prefix string
		value  interface{}
	}{
		fmt.Sprintf("%v %v ?", "updated_user", "NOT IN"),
		value,
	})
	return qb
}

func (qb *authorizedApiQueryBuilder) OrderByUpdatedUser(asc bool) *authorizedApiQueryBuilder {
	order := "DESC"
	if asc {
		order = "ASC"
	}

	qb.order = append(qb.order, "updated_user "+order)
	return qb
}