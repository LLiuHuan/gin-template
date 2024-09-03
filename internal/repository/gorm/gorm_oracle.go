// Package gorm
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-03 22:55
package gorm

import (
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/dzwvip/oracle"
	_ "github.com/godror/godror"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Oracle() *gorm.DB {
	o := configs.Get().Oracle
	if o.Write.DB == "" {
		return nil
	}
	wConf := oracle.Config{
		DSN: o.Dsn(false), // 数据库连接信息
	}
	gormDB, err := gorm.Open(oracle.New(wConf), Config(o.Write.Prefix, o.Write.Singular))
	if err != nil {
		return nil
	}

	if o.IsOpenReadDB {
		rConf := oracle.Config{
			DSN: o.Dsn(true), // 数据库连接信息
		}
		gormDB.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{oracle.New(rConf)},
			Policy:   dbresolver.RandomPolicy{}}).
			SetMaxIdleConns(o.Read.MaxIdleConns).
			SetMaxOpenConns(o.Read.MaxOpenConns))
	}
	gormDB.Use(&TracePlugin{})

	if rawDb, err := gormDB.DB(); err != nil {
		return nil
	} else {
		rawDb.SetMaxIdleConns(o.Read.MaxIdleConns)
		rawDb.SetMaxOpenConns(o.Read.MaxOpenConns)
		return gormDB
	}
}
