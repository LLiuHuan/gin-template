// Package hashids
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 22:00
package hashids

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
	minLength uint8
	blockList []string
}

func New(options ...Option) Hash {
	opt := new(hash)
	for _, f := range options {
		f(opt)
	}
	return opt
}

func (h *hash) i() {}
