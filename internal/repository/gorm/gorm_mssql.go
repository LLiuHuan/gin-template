// Package gorm
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-03 22:53
package gorm

import (
	"github.com/LLiuHuan/gin-template/configs"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Mssql() *gorm.DB {
	m := configs.Get().Mssql
	if m.Write.DB == "" {
		return nil
	}
	wConf := sqlserver.Config{
		DSN: m.Dsn(false), // 数据库连接信息
	}
	gormDB, err := gorm.Open(sqlserver.New(wConf), Config(m.Write.Prefix, m.Write.Singular))
	if err != nil {
		return nil
	}

	if m.IsOpenReadDB {
		rConf := sqlserver.Config{
			DSN: m.Dsn(true), // 数据库连接信息
		}
		gormDB.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{sqlserver.New(rConf)},
			Policy:   dbresolver.RandomPolicy{}}).
			SetMaxIdleConns(m.Read.MaxIdleConns).
			SetMaxOpenConns(m.Read.MaxOpenConns))
	}
	gormDB.Use(&TracePlugin{})

	if rawDb, err := gormDB.DB(); err != nil {
		return nil
	} else {
		rawDb.SetMaxIdleConns(m.Read.MaxIdleConns)
		rawDb.SetMaxOpenConns(m.Read.MaxOpenConns)
		return gormDB
	}
}
