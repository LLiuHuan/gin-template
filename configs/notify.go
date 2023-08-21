// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-08-17 10:22
// @description: 报警配置
package configs

type Notify struct {
	Way    string `mapstructure:"way" json:"way" toml:"way"`          // 告警方式
	Mail   Mail   `mapstructure:"mail" json:"mail" toml:"mail"`       // 告警邮箱配置
	WeChat WeChat `mapstructure:"wechat" json:"wechat" toml:"wechat"` // 告警邮箱配置
}

type Mail struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
	User string `toml:"user"`
	Pass string `toml:"pass"`
	To   string `toml:"to"`
}

type WeChat struct {
	Key string `mapstructure:"key" json:"key" toml:"key"` // 告警邮箱配置
}
