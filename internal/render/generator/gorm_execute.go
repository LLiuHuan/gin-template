// Package generator
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-22 03:14
package generator

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/LLiuHuan/gin-template/configs"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
)

type gormExecuteRequest struct {
	Tables string `form:"tables"`
}

func (h *handler) GormExecute() core.HandlerFunc {
	dir, _ := os.Getwd()
	projectPath := strings.Replace(dir, "\\", "/", -1)
	gormgenSh := projectPath + "/scripts/gormgen.sh"
	gormgenBat := projectPath + "/scripts/gormgen.bat"

	return func(c core.Context) {
		req := new(gormExecuteRequest)
		if err := c.ShouldBindPostForm(req); err != nil {
			c.Payload("参数传递有误")
			return
		}

		mysqlConf := configs.Get().DataBase.MySql.Write
		shellPath := fmt.Sprintf("%s %s:%d %s %s %s %s", gormgenSh, mysqlConf.Host, mysqlConf.Port, mysqlConf.User, mysqlConf.Pass, mysqlConf.DataBase, req.Tables)
		batPath := fmt.Sprintf("%s %s:%d %s %s %s %s", gormgenBat, mysqlConf.Host, mysqlConf.Port, mysqlConf.User, mysqlConf.Pass, mysqlConf.DataBase, req.Tables)

		fmt.Println(shellPath)
		fmt.Println(batPath)
		fmt.Println(req.Tables)
		command := new(exec.Cmd)

		if runtime.GOOS == "windows" {
			command = exec.Command("cmd", "/C", batPath)
		} else {
			// runtime.GOOS = linux or darwin
			command = exec.Command("/bin/bash", "-c", shellPath)
		}

		var stderr bytes.Buffer
		command.Stderr = &stderr

		output, err := command.Output()
		if err != nil {
			c.Payload(stderr.String())
			return
		}

		c.Payload(string(output))
	}
}
