// Package currency
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-09 11:04
package currency

import (
	"testing"
)

func TestSimplifyNum(t *testing.T) {
	// 生成一个单元测试
	if num := SimplifyNum(1, "元"); num != "1 元" {
		t.Errorf("Expected 1 元, but got %s", num)
	}
	if num := SimplifyNum(100, "元"); num != "100 元" {
		t.Errorf("Expected 100 元, but got %s", num)
	}
	if num := SimplifyNum(1000, "元"); num != "1000 元" {
		t.Errorf("Expected 1000 元, but got %s", num)
	}
	if num := SimplifyNum(10000, "元"); num != "1 万元" {
		t.Errorf("Expected 1 万元, but got %s", num)
	}
	if num := SimplifyNum(100000, "元"); num != "10 万元" {
		t.Errorf("Expected 10 万元, but got %s", num)
	}
	if num := SimplifyNum(140312310, "元"); num != " 1.403 亿元" {
		t.Errorf("Expected  1.403 亿元, but got %s", num)
	}
	if num := SimplifyNum(10000000, "元"); num != "1 千万元" {
		t.Errorf("Expected 1 千万元, but got %s", num)
	}
}

func TestFormatFileSize(t *testing.T) {
	result := FormatFileSize(1)
	if result != "1 B" {
		t.Errorf("Expected 1B, but got %s", result)
	}

	result = FormatFileSize(1 << 10)
	if result != "1 KB" {
		t.Errorf("Expected 1KB, but got %s", result)
	}

	result = FormatFileSize(1 << 20)
	if result != "1 MB" {
		t.Errorf("Expected 1MB, but got %s", result)
	}

	result = FormatFileSize(1 << 30)
	if result != "1 GB" {
		t.Errorf("Expected 1GB, but got %s", result)
	}

	result = FormatFileSize(1 << 40)
	if result != "1 TB" {
		t.Errorf("Expected 1TB, but got %s", result)
	}

	//result = FormatFileSize(48712831672)
	//t.Log(result)
}

func TestReversalFileSize(t *testing.T) {
	result := ReversalFileSize("100 GB")
	t.Log(result)
	//if result != 5000 {
	//	t.Errorf("Expected 5000, but got %f", result)
	//}

}
