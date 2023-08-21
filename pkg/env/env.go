// Package env
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-16 15:36
package env

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	active Environment
	dev    Environment = &environment{value: "dev"}
	pro    Environment = &environment{value: "pro"}
)

type Environment interface {
	Value() string
	IsDev() bool
	IsPro() bool
}

type environment struct {
	value string
}

func (e *environment) Value() string {
	return e.value
}

func (e *environment) IsDev() bool {
	return e.value == "dev"
}

func (e *environment) IsPro() bool {
	return e.value == "pro"
}

func init() {
	//env := flag.String("env", "", "请输入运行环境:\n dev:开发环境\n pro:正式环境\n")
	//flag.Parse()
	//
	//switch strings.ToLower(strings.TrimSpace(*env)) {
	//case "dev":
	//	active = dev
	//	gin.SetMode(gin.DebugMode)
	//case "pro":
	//	gin.SetMode(gin.ReleaseMode)
	//	active = pro
	//default:
	active = dev
	gin.SetMode(gin.DebugMode)
	fmt.Println(dev)
	//fmt.Println("Warning: '-env' cannot be found, or it is illegal. The default 'dev' will be used.")
	//}
}

func Active() Environment {
	return active
}
