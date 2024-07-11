// Package configs
//
//	@program:		gin-template
//	@author:		[lliuhuan](https://github.com/lliuhuan)
//	@create:		2024-07-02 21:33
//	@description:	报警配置
package configs

type Notify struct {
	Way    string `mapstructure:"way" json:"way" toml:"way"`          // 告警方式
	Mail   Mail   `mapstructure:"mail" json:"mail" toml:"mail"`       // 告警邮箱配置
	WeChat WeChat `mapstructure:"wechat" json:"wechat" toml:"wechat"` // 告警微信配置
}

type Mail struct {
	Host string `toml:"host"` // 邮箱服务器地址
	Port int    `toml:"port"` // 邮箱服务器端口
	User string `toml:"user"` // 邮箱账号
	Pass string `toml:"pass"` // 邮箱密码
	To   string `toml:"to"`   // 接收人
}

type WeChat struct {
	Key string `mapstructure:"key" json:"key" toml:"key"` // 企业微信机器人key
}
