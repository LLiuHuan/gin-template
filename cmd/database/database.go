// Package main
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-23 17:55
package main

import (
	"fmt"
	"github.com/LLiuHuan/gin-template/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"strings"
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
func New(mode, dbUser, dbPass, dbHose, dbPort, dbName string) (Repo, error) {
	db, err := getDBDriver(mode, dbUser, dbPass, dbHose, dbPort, dbName)
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

func getDBDriver(mode, dbUser, dbPass, dbHose, dbPort, dbName string) (*gorm.DB, error) {
	var dialector gorm.Dialector
	if val, err := getDbDialector(mode, dbUser, dbPass, dbHose, dbPort, dbName); err != nil {
		return nil, err
	} else {
		dialector = val
	}

	gormDB, err := gorm.Open(dialector, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// https://github.com/go-gorm/gorm/issues/3789 可通过自定义Logger解决
	})
	if err != nil {
		return nil, err
	}

	return gormDB, nil
}

// 获取一个数据库方言(Dialector),通俗地说就是根据不同的连接参数，获取具体的一类数据库的连接指针
func getDbDialector(mode, dbUser, dbPass, dbHose, dbPort, dbName string) (gorm.Dialector, error) {
	var dbDialector gorm.Dialector
	dsn := getDsn(mode, dbUser, dbPass, dbHose, dbPort, dbName)
	switch strings.ToLower(mode) {
	case "mysql":
		dbDialector = mysql.Open(dsn)
	case "sqlserver", "mssql":
		dbDialector = sqlserver.Open(dsn)
	case "postgres", "postgresql", "postgre":
		dbDialector = postgres.Open(dsn)
	default:
		return nil, errors.New("不支持的数据库类型")
	}
	return dbDialector, nil
}

// 根据配置参数生成数据库驱动 dsn
func getDsn(mode, dbUser, dbPass, dbHose, dbPort, dbName string) string {

	switch strings.ToLower(mode) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", dbUser, dbPass, dbHose, dbPort, dbName, "utf8mb4")
	case "sqlserver", "mssql":
		return fmt.Sprintf("server=%s;port=%s;database=%s;user id=%s;password=%s;encrypt=disable", dbHose, dbPort, dbName, dbUser, dbPass)
	case "postgresql", "postgre", "postgres":
		return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", dbHose, dbPort, dbName, dbUser, dbPass)
	}
	return ""
}
