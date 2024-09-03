// Package gorm
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-03 18:14
package gorm

import (
	"github.com/LLiuHuan/gin-template/configs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func Pgsql() *gorm.DB {
	p := configs.Get().Pgsql
	if p.Write.DB == "" {
		return nil
	}

	wConf := postgres.Config{
		DSN:                  p.Dsn(false),
		PreferSimpleProtocol: true,
	}
	gormDB, err := gorm.Open(postgres.New(wConf), Config(p.Write.Prefix, p.Write.Singular))
	if err != nil {
		return nil
	}

	if p.IsOpenReadDB {
		rConf := postgres.Config{
			DSN:                  p.Dsn(true),
			PreferSimpleProtocol: true,
		}
		gormDB.Use(dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{postgres.New(rConf)},
			Policy:   dbresolver.RandomPolicy{},
		}).SetMaxIdleConns(p.Read.MaxIdleConns).SetMaxOpenConns(p.Read.MaxOpenConns))
	}
	gormDB.Use(&TracePlugin{})

	if rawDB, err := gormDB.DB(); err != nil {
		return nil
	} else {
		rawDB.SetMaxIdleConns(p.Write.MaxIdleConns)
		rawDB.SetMaxOpenConns(p.Write.MaxOpenConns)
		return gormDB
	}
}
