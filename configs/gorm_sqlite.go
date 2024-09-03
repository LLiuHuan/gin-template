// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-09-02 11:19
package configs

import "path/filepath"

type Sqlite struct {
	GeneralDB `toml:",inline" mapstructure:",squash"`
}

// Dsn /path/to/sqlite.db
func (s *Sqlite) Dsn(isRead bool) string {
	if isRead {
		return filepath.Join(s.Read.Pass, s.Read.DB+".db")
	}
	return filepath.Join(s.Write.Pass, s.Write.DB+".db")
}
