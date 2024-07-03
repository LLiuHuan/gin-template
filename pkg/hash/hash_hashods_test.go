// Package hash
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-02 22:21
package hash

import "testing"

const secret = "i1ydX9RtHyuJTrw7frcu"
const length = 12

func TestHashidsEncode(t *testing.T) {
	str, _ := New(secret, length).HashidsEncode([]int{99})
	t.Log(str)

	//GyV5pJqXvwAR
}

func TestHashidsDecode(t *testing.T) {
	ids, _ := New(secret, length).HashidsDecode("GyV5pJqXvwAR")
	t.Log(ids)
}
