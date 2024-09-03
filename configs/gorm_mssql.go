// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 11:08
package configs

type Mssql struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn sqlserver://user:password@localhost:1433?gorm=dbname
func (m *Mssql) Dsn(isRead bool) string {
	if isRead {
		return "sqlserver://" + m.Read.User + ":" + m.Read.Pass + "@" + m.Read.Path + ":" + m.Read.Port + "?gorm=" + m.Read.DB + "&encrypt=disable"
	}

	return "sqlserver://" + m.Write.User + ":" + m.Write.Pass + "@" + m.Write.Path + ":" + m.Write.Port + "?gorm=" + m.Write.DB + "&encrypt=disable"
}
