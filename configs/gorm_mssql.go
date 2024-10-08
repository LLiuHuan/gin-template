// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 11:08
package configs

type Mssql struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn sqlserver://user:password@localhost:1433?gormDB=dbname
func (m *Mssql) Dsn() string {
	return "sqlserver://" + m.User + ":" + m.Pass + "@" + m.Path + ":" + m.Port + "?gormDB=" + m.DB + "&encrypt=disable"
}
