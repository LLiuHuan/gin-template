// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 11:11
package configs

type Oracle struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn oracle://username:password@host:port/dbname?charset=utf8&parseTime=True&loc=Local
func (o *Oracle) Dsn(isRead bool) string {
	if isRead {
		return "oracle://" + o.Read.User + ":" + o.Read.Pass + "@" + o.Read.Path + ":" + o.Read.Port + "/" + o.Read.DB + "?" + o.Read.Config
	}
	return "oracle://" + o.Write.User + ":" + o.Write.Pass + "@" + o.Write.Path + ":" + o.Write.Port + "/" + o.Write.DB + "?" + o.Write.Config
}
