// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 11:11
package configs

type Oracle struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn oracle://username:password@host:port/dbname?charset=utf8&parseTime=True&loc=Local
func (o *Oracle) Dsn() string {
	return "oracle://" + o.User + ":" + o.Pass + "@" + o.Path + ":" + o.Port + "/" + o.DB + "?" + o.Config
}
