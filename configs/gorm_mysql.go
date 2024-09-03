// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 10:43
package configs

type Mysql struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func (m *Mysql) Dsn(isRead bool) string {
	if isRead {
		return m.Read.User + ":" + m.Read.Pass + "@tcp(" + m.Read.Path + ":" + m.Read.Port + ")/" + m.Read.DB + "?" + m.Read.Config
	}
	return m.Write.User + ":" + m.Write.Pass + "@tcp(" + m.Write.Path + ":" + m.Write.Port + ")/" + m.Write.DB + "?" + m.Write.Config
}
