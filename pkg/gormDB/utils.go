// Package gormDB
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-10-08 10:55
package gormDB

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/LLiuHuan/gin-template/pkg/errors"

	"github.com/dzwvip/oracle"
	"github.com/glebarez/sqlite"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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

func NewDB(driver, user, pass, host, port, db string) (Repo, error) {
	// 获取DSN
	dialector, err := Dialector(driver, user, pass, host, port, db)
	if err != nil {
		return nil, err
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

	return &dbRepo{
		DB: gormDB,
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

// DSN 根据配置参数生成数据库驱动 dsn
func DSN(driver, user, pass, host, port, db string) string {
	switch strings.ToLower(driver) {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", user, pass, host, port, db, "utf8mb4")
	case "sqlserver", "mssql":
		return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=disable", user, pass, host, port, db)
	case "postgresql", "postgre", "postgres":
		return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable TimeZone=Asia/Shanghai", host, port, db, user, pass)
	case "oracle":
		return fmt.Sprintf("oracle://%s:%s@%s:%s/%s?charset=utf8&parseTime=True&loc=Local", user, pass, host, port, db)
	case "sqlite":
		return filepath.Join(host, db+".db")
	default:
		return ""
	}
}

func Dialector(driver, user, pass, host, port, db string) (gorm.Dialector, error) {
	dsn := DSN(driver, user, pass, host, port, db)
	fmt.Println(dsn)
	switch strings.ToLower(driver) {
	case "mysql":
		return mysql.Open(dsn), nil
	case "sqlserver", "mssql":
		return sqlserver.Open(dsn), nil
	case "postgresql", "postgre", "postgres":
		return postgres.Open(dsn), nil
	case "oracle":
		return oracle.Open(dsn), nil
	case "sqlite":
		return sqlite.Open(dsn), nil
	default:
		return nil, errors.New("不支持的数据库类型")
	}
}
