// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-09-04 18:34
package configs

type HashIds struct {
	Alphabet  string   `toml:"alphabet"`
	MinLength int      `toml:"min-length"`
	BlockList []string `toml:"block-list"`
}
