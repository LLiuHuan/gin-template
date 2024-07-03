// Package configs
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 21:32
package configs

type HashIds struct {
	Alphabet  string   `toml:"alphabet"`
	MinLength uint8    `toml:"min-length"`
	BlockList []string `toml:"block-list"`
}
