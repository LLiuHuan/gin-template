// Package hashids
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 22:00
package hashids

import "github.com/sqids/sqids-go"

func (h *hash) HashidsEncode(params []uint64) (string, error) {
	s, err := sqids.New(sqids.Options{
		Alphabet:  h.alphabet,
		MinLength: h.minLength,
		Blocklist: h.blockList,
	})
	if err != nil {
		return "", err
	}
	encode, err := s.Encode(params)
	if err != nil {
		return "", err
	}
	return encode, nil
}

func (h *hash) HashidsDecode(hash string) ([]uint64, error) {
	s, err := sqids.New(sqids.Options{
		Alphabet:  h.alphabet,
		MinLength: h.minLength,
		Blocklist: h.blockList,
	})
	if err != nil {
		return nil, err
	}
	return s.Decode(hash), nil
}
