// Package currency
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-05 17:39
package currency

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// FormatFileSize 字节的单位转换 保留两位小数, 	基数是1024
func FormatFileSize(fileSize int64) string {
	units := []struct {
		threshold int64
		format    string
	}{
		{1 << 10, "%.2f B"},
		{1 << 20, "%.2f KB"},
		{1 << 30, "%.2f MB"},
		{1 << 40, "%.2f GB"},
		{1 << 50, "%.2f TB"},
		{1 << 60, "%.2f EB"},
		{math.MaxInt64, "%.2f ZB"},
	}

	var unit string
	var value float64

	for _, u := range units {
		if fileSize < u.threshold {
			unit = u.format
			value = float64(fileSize) / float64(u.threshold>>10)
			break
		}
	}

	return strings.ReplaceAll(fmt.Sprintf(unit, value), ".00", "")
}

// ReversalFileSize 字节的单位转换 保留两位小数

func ReversalFileSize(fileSize string) (size float64) {
	units := []struct {
		suffix     string
		multiplier float64
	}{
		{"KB", float64(1 << 10)},
		{"MB", float64(1 << 20)},
		{"GB", float64(1 << 30)},
		{"TB", float64(1 << 40)},
		{"EB", float64(1 << 50)},
		{"ZB", float64(1 << 60)},
		{"B", float64(1)},
	}

	for _, unit := range units {
		if strings.HasSuffix(fileSize, unit.suffix) {
			value, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(fileSize, unit.suffix)), 64)
			if err != nil {
				return 0
			}
			size = value * unit.multiplier
			break
		}
	}
	return size / float64(1<<20)
}

// ReversalFileSize 字节的单位转换 保留两位小数
//func ReversalFileSize(fileSize string) (size float64) {
//	if strings.HasSuffix(fileSize, "KB") {
//		float, err := strconv.ParseFloat(strings.Replace(fileSize, "KB", "", -1), 64)
//		if err != nil {
//			return 0
//		}
//		size = float * float64(1<<10)
//	} else if strings.HasSuffix(fileSize, "MB") {
//		float, err := strconv.ParseFloat(strings.Replace(fileSize, "MB", "", -1), 64)
//		if err != nil {
//			return 0
//		}
//		size = float * float64(1<<20)
//	} else if strings.HasSuffix(fileSize, "GB") {
//		float, err := strconv.ParseFloat(strings.Replace(fileSize, "GB", "", -1), 64)
//		if err != nil {
//			return 0
//		}
//		size = float * float64(1<<30)
//	} else if strings.HasSuffix(fileSize, "TB") {
//		float, err := strconv.ParseFloat(strings.Replace(fileSize, "TB", "", -1), 64)
//		if err != nil {
//			return 0
//		}
//		size = float * float64(1<<40)
//	} else if strings.HasSuffix(fileSize, "EB") {
//		float, err := strconv.ParseFloat(strings.Replace(fileSize, "EB", "", -1), 64)
//		if err != nil {
//			return 0
//		}
//		size = float * float64(1<<50)
//	} else if strings.HasSuffix(fileSize, "ZB") {
//		float, err := strconv.ParseFloat(strings.Replace(fileSize, "ZB", "", -1), 64)
//		if err != nil {
//			return 0
//		}
//		size = float * float64(1<<60)
//	} else {
//		float, err := strconv.ParseFloat(strings.Replace(fileSize, "B", "", -1), 64)
//		if err != nil {
//			return 0
//		}
//		size = float * float64(1)
//	}
//	return size / float64(1<<20)
//}

const (
	UnitWan     = "万"
	UnitQianWan = "千万"
	UnitYi      = "亿"
	UnitBaiYi   = "百亿"
	UnitQianYi  = "千亿"
	UnitWanYi   = "万亿"
)

func SimplifyNum(num float64, suffix string) string {
	units := []struct {
		threshold float64
		unit      string
	}{
		{1e12, UnitWanYi},
		{1e11, UnitQianYi},
		{1e10, UnitBaiYi},
		{1e8, UnitYi},
		{1e7, UnitQianWan},
		{1e4, UnitWan},
	}

	num = math.Round(num*100) / 100
	if num == 0 {
		return "0"
	}

	for _, u := range units {
		if num >= u.threshold {
			return strings.ReplaceAll(fmt.Sprintf("%.3f %s%s", num/u.threshold, u.unit, suffix), ".000", "")
		}
	}

	return strings.ReplaceAll(fmt.Sprintf("%.3f %s", num, suffix), ".000", "")
}

//func SimplifyNum(num float64, suffix string) string {
//	// 保留两位小数
//	num = math.Round(num*100) / 100
//	if num <= 0 {
//		return "0"
//	} else if num >= 1000 && num < 1000000 {
//		return fmt.Sprintf("%.3f 万%s", num/10000, suffix)
//	} else if num >= 10000000 && num < 100000000 {
//		return fmt.Sprintf("%.3f 千万%s", num/10000000, suffix)
//	} else if num >= 100000000 && num < 10000000000 {
//		return fmt.Sprintf("%.3f 亿%s", num/100000000, suffix)
//	} else if num >= 10000000000 && num < 100000000000 {
//		return fmt.Sprintf("%.3f 百亿%s", num/10000000000, suffix)
//	} else if num >= 100000000000 && num < 1000000000000 {
//		return fmt.Sprintf("%.3f 千亿%s", num/100000000000, suffix)
//	} else if num >= 1000000000000 && num < 10000000000000 {
//		return fmt.Sprintf("%.3f 万亿%s", num/1000000000000, suffix)
//	} else {
//		return fmt.Sprintf("%f%s", num, suffix)
//	}
//}
