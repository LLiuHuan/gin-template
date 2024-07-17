// Package hashids
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-12 11:08
package hashids

type Option func(*hash)

func WithAlphabet(alphabet string) Option {
	return func(h *hash) {
		h.alphabet = alphabet
	}
}

func WithMinLength(minLength uint8) Option {
	return func(h *hash) {
		h.minLength = minLength
	}
}

func WithBlockList(blockList []string) Option {
	return func(h *hash) {
		h.blockList = blockList
	}
}
