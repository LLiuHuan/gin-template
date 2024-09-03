// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 10:44
package configs

import (
	"gorm.io/gorm/logger"
	"strings"
)

type DsnProvider interface {
	Dsn(isRead bool) string
}

// GeneralDB 通用数据库配置
type GeneralDB struct {
	IsOpenReadDB bool            `mapstructure:"isOpenReadDB" json:"isOpenReadDB" toml:"isOpenReadDB"` // 是否开启读写分离，如果不开启的话就默认忽略Read库的配置
	Read         GeneralDBConfig `mapstructure:"read" json:"read" toml:"read"`
	Write        GeneralDBConfig `mapstructure:"write" json:"write" toml:"write"`
}

type GeneralDBConfig struct {
	Prefix       string `mapstructure:"prefix" json:"prefix" toml:"prefix"`                   // 数据库前缀
	Path         string `mapstructure:"path" json:"path" toml:"path"`                         // path / host
	Port         string `mapstructure:"port" json:"port" toml:"port"`                         // 数据库端口
	User         string `mapstructure:"user" json:"user" toml:"user"`                         // 数据库账号
	Pass         string `mapstructure:"pass" json:"pass" toml:"pass"`                         // 数据库密码
	DB           string `mapstructure:"db" json:"db" toml:"db"`                               // 数据库名
	Config       string `mapstructure:"config" json:"config" toml:"config"`                   // 高级配置
	Engine       string `mapstructure:"engine" json:"engine" toml:"engine" default:"InnoDB"`  // 数据库引擎，默认InnoDB
	Singular     bool   `mapstructure:"singular" json:"singular" toml:"singular"`             // 是否开启全局禁用复数，true表示开启
	LogMode      string `mapstructure:"logMode" json:"logMode" toml:"logMode"`                // 是否开启Gorm全局日志
	MaxIdleConns int    `mapstructure:"maxIdleConns" json:"maxIdleConns" toml:"maxIdleConns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"maxOpenConns" json:"maxOpenConns" toml:"maxOpenConns"` // 打开到数据库的最大连接数

	LogZap bool `mapstructure:"logZap" json:"logZap" yaml:"logZap"`
}

func (c GeneralDBConfig) LogLevel() logger.LogLevel {
	switch strings.ToLower(c.LogMode) {
	case "silent", "Silent":
		return logger.Silent
	case "error", "Error":
		return logger.Error
	case "warn", "Warn":
		return logger.Warn
	case "info", "Info":
		return logger.Info
	default:
		return logger.Info
	}
}

type SpecializedDB struct {
	Type      string `mapstructure:"type" json:"type" toml:"type"`
	AliasName string `mapstructure:"aliasName" json:"aliasName" toml:"aliasName"`
	GeneralDB `mapstructure:",squash" json:"generalDB" toml:",inline"`
	Disable   bool `mapstructure:"disable" json:"disable" toml:"disable"`
}
