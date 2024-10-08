// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 10:43
package configs

type Mysql struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn user:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local
func (m *Mysql) Dsn() string {
	return m.User + ":" + m.Pass + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.DB + "?" + m.Config
}
