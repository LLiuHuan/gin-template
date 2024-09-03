// Package gorm
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-03 22:57
package gorm

import (
	"github.com/LLiuHuan/gin-template/configs"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Sqlite() *gorm.DB {
	s := configs.Get().Sqlite
	if s.Write.DB == "" {
		return nil
	}
	gormDB, err := gorm.Open(sqlite.Open(s.Dsn(false)), Config(s.Write.Prefix, s.Write.Singular))
	if err != nil {
		return nil
	}

	gormDB.Use(&TracePlugin{})

	if rawDb, err := gormDB.DB(); err != nil {
		return nil
	} else {
		rawDb.SetMaxIdleConns(s.Read.MaxIdleConns)
		rawDb.SetMaxOpenConns(s.Read.MaxOpenConns)
		return gormDB
	}
}
