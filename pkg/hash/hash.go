// Package hash
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2023-09-04 18:17
package hash

var _ Hash = (*hash)(nil)

type Hash interface {
	i()

	// HashidsEncode 加密
	HashidsEncode(params []uint64) (string, error)

	// HashidsDecode 解密
	HashidsDecode(hash string) ([]uint64, error)
}

type hash struct {
	alphabet  string
	minLength int
	blockList []string
}

func New(alphabet string, minLength int, blockList []string) Hash {
	return &hash{
		alphabet:  alphabet,
		minLength: minLength,
		blockList: blockList,
	}
}

func (h *hash) i() {}
