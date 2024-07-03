// Package app
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-04 03:11
package app

import (
	"go.uber.org/zap"
	"os"
	"time"
)

//var execStopFunc bool

func AppPrepareForceExit(logger *zap.Logger) {
	//if !execStopFunc {
	//	return
	//}
	time.AfterFunc(30*time.Second, func() {
		logger.Info("App server Shutdown timeout, force exit")
		os.Exit(1)
	})
}
