// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 11:15
package configs

type Pgsql struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn host=127.0.0.1 port=5432 user=root password=123456 dbname=gin-template sslmode=disable
func (p *Pgsql) Dsn(isRead bool) string {
	if isRead {
		return "host=" + p.Read.Path + " port=" + p.Read.Port + " user=" + p.Read.User + " password=" + p.Read.Pass + " dbname=" + p.Read.DB + " " + p.Read.Config
	}
	return "host=" + p.Write.Path + " port=" + p.Write.Port + " user=" + p.Write.User + " password=" + p.Write.Pass + " dbname=" + p.Write.DB + " " + p.Write.Config
}
