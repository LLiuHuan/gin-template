// Package hashids
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 22:21
package hashids

import "testing"

func TestHashidsEncode(t *testing.T) {
	str, _ := New(
		WithAlphabet("FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE"),
		WithMinLength(10),
	).HashidsEncode([]uint64{1, 2, 3})
	t.Log(str)

	//B4aajshuHt
}

func TestHashidsDecode(t *testing.T) {
	ids, _ := New(
		WithAlphabet("FxnXM1kBN6cuhsAvjW3Co7l2RePyY8DwaU04Tzt9fHQrqSVKdpimLGIJOgb5ZE"),
		WithMinLength(10),
	).HashidsDecode("B4aajshuHt")
	t.Log(ids)
}
