// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 11:19
package configs

import "path/filepath"

type Sqlite struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn /path/to/sqlite.gormDB
func (s *Sqlite) Dsn() string {
	return filepath.Join(s.Pass, s.DB+".db")
}
