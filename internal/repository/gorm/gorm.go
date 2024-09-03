// Package gorm
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 10:14
package gorm

import (
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/LLiuHuan/gin-template/configs"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Predicate is a string that acts as a condition in the where clause
type Predicate string

var (
	EqualPredicate              = Predicate("=")
	NotEqualPredicate           = Predicate("<>")
	GreaterThanPredicate        = Predicate(">")
	GreaterThanOrEqualPredicate = Predicate(">=")
	SmallerThanPredicate        = Predicate("<")
	SmallerThanOrEqualPredicate = Predicate("<=")
	LikePredicate               = Predicate("LIKE")
)

var _ Repo = (*dbRepo)(nil)

type Repo interface {
	i()
	GetDB() *gorm.DB
	DBClose() error
}

type dbRepo struct {
	DB *gorm.DB
}

// New 创建一个DB对象
func New() (Repo, error) {
	var db *gorm.DB
	switch configs.Get().Project.DbType {
	case "mysql":
		db = Mysql()
	case "postgres", "postgresql", "postgre", "pgsql":
		db = Pgsql()
	case "sqlserver", "mssql":
		db = Mssql()
	case "oracle":
		db = Oracle()
	case "sqlite":
		db = Sqlite()
	default:
		db = Mysql()
	}
	return &dbRepo{
		DB: db,
	}, nil
}

func (d *dbRepo) i() {}

func (d *dbRepo) GetDB() *gorm.DB {
	return d.DB
}

func (d *dbRepo) DBClose() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
func Config(prefix string, singular bool) *gorm.Config {
	var general configs.GeneralDB
	switch configs.Get().Project.DbType {
	case "mysql":
		general = configs.Get().Mysql.GeneralDB
	case "postgres", "postgresql", "postgre", "pgsql":
		general = configs.Get().Pgsql.GeneralDB
	case "sqlserver", "mssql":
		general = configs.Get().Mssql.GeneralDB
	case "oracle":
		general = configs.Get().Oracle.GeneralDB
	case "sqlite":
		general = configs.Get().Sqlite.GeneralDB
	default:
		general = configs.Get().Mysql.GeneralDB
	}

	return &gorm.Config{
		Logger: logger.New(NewWriter(general, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
			SlowThreshold: 200 * time.Millisecond,
			LogLevel:      general.Write.LogLevel(),
			Colorful:      true,
		}),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   prefix,   // 表前缀
			SingularTable: singular, // 单数表名
		},
		DisableForeignKeyConstraintWhenMigrating: true, // 关闭自动外键约束
	}
}
