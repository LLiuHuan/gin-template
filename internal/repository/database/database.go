// Package database
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-03 10:14
package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/code"
	"github.com/LLiuHuan/gin-template/pkg/errors"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

// TODO：需要修改，目前只支持MySql

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
	db, err := useDBConn()
	if err != nil {
		return nil, err
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

func useDBConn() (*gorm.DB, error) {
	isOpenReadDB := 0
	switch configs.Get().DataBase.Mode {
	case "mysql":
		isOpenReadDB = configs.Get().DataBase.MySql.IsOpenReadDB
	case "postgres", "postgresql", "postgre":
		isOpenReadDB = configs.Get().DataBase.PostgreSql.IsOpenReadDB
	case "sqlserver", "mssql":
		isOpenReadDB = configs.Get().DataBase.SqlServer.IsOpenReadDB
	default:
		return nil, nil
	}
	return getDBDriver(configs.Get().DataBase.Mode, isOpenReadDB)
}

func GetMySQLClient() (*gorm.DB, error) {
	return getDBDriver("mysql", configs.Get().DataBase.MySql.IsOpenReadDB)
}

func GetSqlServerClient() (*gorm.DB, error) {
	return getDBDriver("sqlserver", configs.Get().DataBase.SqlServer.IsOpenReadDB)
}

func GetPostgreSqlClient() (*gorm.DB, error) {
	return getDBDriver("postgresql", configs.Get().DataBase.PostgreSql.IsOpenReadDB)
}

func getDBDriver(mode string, isOpenReadDB int) (*gorm.DB, error) {
	var dialector gorm.Dialector
	if val, err := getDbDialector(mode, "write"); err != nil {
		return nil, err
	} else {
		dialector = val
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, err
	}

	// 如果开启了读写分离，就再次获取一个读数据库的连接
	if isOpenReadDB == 1 {
		if val, err := getDbDialector(mode, "read"); err != nil {
			return nil, err
		} else {
			dialector = val
		}
		resolverConf := dbresolver.Config{
			Replicas: []gorm.Dialector{dialector}, //  读 操作库，查询类
			Policy:   dbresolver.RandomPolicy{},   // sources/replicas 负载均衡策略适用于
		}

		err = gormDB.Use(dbresolver.Register(resolverConf).SetConnMaxIdleTime(time.Second * 30).
			SetConnMaxLifetime(configs.GetContainer().GetDuration(fmt.Sprintf("database.%s.base.connMaxLifeTime", mode)) * time.Second).
			SetMaxIdleConns(configs.GetContainer().GetInt(fmt.Sprintf("database.%s.base.maxIdleConn", mode))).
			SetMaxOpenConns(configs.GetContainer().GetInt(fmt.Sprintf("database.%s.base.maxOpenConn", mode))))
		if err != nil {
			return nil, err
		}
	}

	gormDB.Use(&TracePlugin{})

	// 为主连接设置连接池(43行返回的数据库驱动指针)
	if rawDb, err := gormDB.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetConnMaxIdleTime(time.Second * 30)
		rawDb.SetConnMaxLifetime(configs.GetContainer().GetDuration(fmt.Sprintf("database.%s.base.connMaxLifeTime", mode)) * time.Second)
		rawDb.SetMaxIdleConns(configs.GetContainer().GetInt(fmt.Sprintf("database.%s.base.maxIdleConn", mode)))
		rawDb.SetMaxOpenConns(configs.GetContainer().GetInt(fmt.Sprintf("database.%s.base.maxOpenConn", mode)))
		return gormDB, nil
	}
}

// 获取一个数据库方言(Dialector),通俗地说就是根据不同的连接参数，获取具体的一类数据库的连接指针
func getDbDialector(sqlType string, readWrite string) (gorm.Dialector, error) {
	var dbDialector gorm.Dialector
	dsn := getDsn(sqlType, readWrite)
	switch strings.ToLower(sqlType) {
	case "mysql":
		dbDialector = mysql.Open(dsn)
	case "sqlserver", "mssql":
		dbDialector = sqlserver.Open(dsn)
	case "postgres", "postgresql", "postgre":
		dbDialector = postgres.Open(dsn)
	default:
		return nil, errors.New(code.Text(code.DBDriverNotExists) + sqlType)
	}
	return dbDialector, nil
}

// 根据配置参数生成数据库驱动 dsn
func getDsn(sqlType string, readWrite string) string {
	var (
		User     = configs.GetContainer().GetString(fmt.Sprintf("database.%s.%s.user", sqlType, readWrite))
		Pass     = configs.GetContainer().GetString(fmt.Sprintf("database.%s.%s.pass", sqlType, readWrite))
		Host     = configs.GetContainer().GetString(fmt.Sprintf("database.%s.%s.host", sqlType, readWrite))
		Port     = configs.GetContainer().GetInt(fmt.Sprintf("database.%s.%s.port", sqlType, readWrite))
		DataBase = configs.GetContainer().GetString(fmt.Sprintf("database.%s.%s.dataBase", sqlType, readWrite))
		Charset  = configs.GetContainer().GetString(fmt.Sprintf("database.%s.%s.charset", sqlType, readWrite))
	)

	switch strings.ToLower(sqlType) {
	case "mysql":
		if Charset == "" {
			Charset = "utf8mb4"
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=false&loc=Local", User, Pass, Host, Port, DataBase, Charset)
	case "sqlserver", "mssql":
		return fmt.Sprintf("server=%s;port=%d;database=%s;user id=%s;password=%s;encrypt=disable", Host, Port, DataBase, User, Pass)
	case "postgresql", "postgre", "postgres":
		return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", Host, Port, DataBase, User, Pass)
	}
	return ""
}
