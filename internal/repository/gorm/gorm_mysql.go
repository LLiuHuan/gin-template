// Package gorm
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 18:00
package gorm

import (
	"github.com/LLiuHuan/gin-template/configs"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Mysql() *gorm.DB {
	m := configs.Get().Mysql
	if m.Write.DB == "" {
		return nil
	}
	wConf := mysql.Config{
		DSN:                       m.Dsn(false), // 数据库连接信息
		SkipInitializeWithVersion: false,        // 禁用根据版本自动配置
	}
	gormDB, err := gorm.Open(mysql.New(wConf), Config(m.Write.Prefix, m.Write.Singular))
	if err != nil {
		return nil
	}

	if m.IsOpenReadDB {
		rConf := mysql.Config{
			DSN:                       m.Dsn(true), // 数据库连接信息
			SkipInitializeWithVersion: false,       // 禁用根据版本自动配置
		}
		gormDB.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{mysql.New(rConf)},
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
