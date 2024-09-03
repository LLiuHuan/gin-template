// Package gorm
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-03 17:30
package gorm

import (
	"fmt"
	"github.com/LLiuHuan/gin-template/configs"
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type Writer struct {
	config configs.GeneralDB
	writer logger.Writer
}

func NewWriter(config configs.GeneralDB, writer logger.Writer) *Writer {
	return &Writer{config: config, writer: writer}
}

// Printf 格式化打印日志
func (c *Writer) Printf(message string, data ...any) {
	if c.config.Write.LogZap {
		switch c.config.Write.LogLevel() {
		case logger.Silent:
			zap.L().Debug(fmt.Sprintf(message, data...))
		case logger.Error:
			zap.L().Error(fmt.Sprintf(message, data...))
		case logger.Warn:
			zap.L().Warn(fmt.Sprintf(message, data...))
		case logger.Info:
			zap.L().Info(fmt.Sprintf(message, data...))
		default:
			zap.L().Info(fmt.Sprintf(message, data...))
		}
		return
	}
	c.writer.Printf(message, data...)
}
