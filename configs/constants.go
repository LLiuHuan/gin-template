// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:32
// @description: 常量配置
package configs

import (
	"strings"
	"sync"
	"time"

	"github.com/LLiuHuan/gin-template/pkg/errors"
)

const (
	// MinGoVersion 支持的最小Go版本
	MinGoVersion = 1.16
	// ProjectVersion 项目版本
	ProjectVersion = "v0.0.1"
	// ProjectName 项目名称
	ProjectName = "gin-template"

	// ProjectLogFile 项目访问日志存放文件
	ProjectLogFile = "./logs/" + ProjectName + ".log"

	// ProjectCronLogFile 项目后台任务日志存放文件
	ProjectCronLogFile = "./logs/" + ProjectName + "-cron.log"

	// ZhCN 简体中文 - 中国
	ZhCN = "zh-cn"

	// EnUS 英文 - 美国
	EnUS = "en-us"

	// HeaderLoginToken 登录验证 Token，Header 中传递的参数
	HeaderLoginToken = "Token"

	// HeaderSignToken 签名验证 Authorization，Header 中传递的参数
	HeaderSignToken = "Authorization"

	// HeaderSignTokenDate 签名验证 Date，Header 中传递的参数
	HeaderSignTokenDate = "Authorization-Date"

	// MaxRequestsPerSecond 每秒最大请求量
	MaxRequestsPerSecond = 10000

	ProjectNotifyMail   = "mail"
	ProjectNotifyWeChat = "wechat"

	ProjectDBMySQL     = "mysql"
	ProjectDBSQLServer = "sqlserver"
	ProjectDBPGSQL     = "postgresql"

	ConfigKeyPrefix = "Config_" //  配置文件键值缓存时，键的前缀

	// RedisKeyPrefixLoginUser Redis Key 前缀 - 登录用户信息
	RedisKeyPrefixLoginUser = ProjectName + ":login-user:"

	// RedisKeyPrefixSignature Redis Key 前缀 - 签名验证信息
	RedisKeyPrefixSignature = ProjectName + ":signature:"

	// LoginSessionTTL 登录有效期为 24 小时
	LoginSessionTTL = time.Hour * 24

	// ProjectInstallMark 项目安装完成标识
	ProjectInstallMark = "INSTALL.lock"

	// HeaderSignTokenTimeout 签名有效期为 2 分钟
	HeaderSignTokenTimeout = time.Minute * 2
)

func init() {
}

// 定义一个全局键值对存储容器
var sMap sync.Map

// CreateContainersFactory 创建一个容器工厂
func CreateContainersFactory() *containers {
	return &containers{}
}

// 定义一个容器结构体
type containers struct {
}

// Set 以键值对的形式将代码注册到容器
func (c *containers) Set(key string, value interface{}) (res bool, err error) {
	if _, exists := c.KeyIsExists(key); exists == false {
		sMap.Store(key, value)
		res = true
	} else {
		return false, errors.New("该键已经注册在容器中了, 相关键：" + key)
	}
	return
}

// Delete 删除
func (c *containers) Delete(key string) {
	sMap.Delete(key)
}

// Get 传递键，从容器获取值
func (c *containers) Get(key string) interface{} {
	if value, exists := c.KeyIsExists(key); exists {
		return value
	}
	return nil
}

// KeyIsExists 判断键是否被注册
func (c *containers) KeyIsExists(key string) (interface{}, bool) {
	return sMap.Load(key)
}

// FuzzyDelete 按照键的前缀模糊删除容器中注册的内容
func (c *containers) FuzzyDelete(keyPre string) {
	sMap.Range(func(key, value interface{}) bool {
		if keyname, ok := key.(string); ok {
			if strings.HasPrefix(keyname, keyPre) {
				sMap.Delete(keyname)
			}
		}
		return true
	})
}
