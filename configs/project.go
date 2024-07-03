// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:33
package configs

type Project struct {
	Domain  string `mapstructure:"domain" json:"domain" toml:"domain"`    // 项目ip/域名
	Port    int    `mapstructure:"port" json:"port" toml:"port"`          // 项目端口
	Local   string `mapstructure:"local" json:"local" toml:"local"`       // 中英文 zh-cn/en-us
	PidFile string `mapstructure:"pidfile" json:"pidfile" toml:"pidfile"` // pid文件
	//Name          string `mapstructure:"name" json:"name" yaml:"name"`                              // 项目名称
	//Version       string `mapstructure:"version" json:"version" yaml:"version"`                     // 项目版本
	//Port          string `mapstructure:"port" json:"port" yaml:"port"`                              // 端口值
	//Login         string `mapstructure:"login" json:"login" yaml:"login"`                           // 登陆接口地址
	//Service       string `mapstructure:"service" json:"service" yaml:"service"`                     // 服务Web地址
	//UseMultipoint bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"` // 多点登录拦截 限制ip等
}
