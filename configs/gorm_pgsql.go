// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 11:15
package configs

type Pgsql struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn host=127.0.0.1 port=5432 user=root password=123456 dbname=gin-template sslmode=disable
func (p *Pgsql) Dsn() string {
	return "host=" + p.Path + " port=" + p.Port + " user=" + p.User + " password=" + p.Pass + " dbname=" + p.DB + " " + p.Config
}
