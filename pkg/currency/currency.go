// Package currency
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-05 17:39
package currency

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatFileSize 字节的单位转换 保留两位小数
func FormatFileSize(fileSize int64) (size string) {
	if fileSize < 1<<10 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2f B", float64(fileSize)/float64(1))
	} else if fileSize < 1<<20 {
		return fmt.Sprintf("%.2f KB", float64(fileSize)/float64(1<<10))
	} else if fileSize < 1<<30 {
		return fmt.Sprintf("%.2f MB", float64(fileSize)/float64(1<<20))
	} else if fileSize < 1<<40 {
		return fmt.Sprintf("%.2f GB", float64(fileSize)/float64(1<<30))
	} else if fileSize < 1<<50 {
		return fmt.Sprintf("%.2f TB", float64(fileSize)/float64(1<<40))
	} else if fileSize < 1<<60 {
		return fmt.Sprintf("%.2f EB", float64(fileSize)/float64(1<<50))
	} else { //if fileSize < (1000 * 1000 * 1000 * 1000 * 1000 * 1000)
		return fmt.Sprintf("%.2f ZB", float64(fileSize)/float64(1<<60))
	}
}

// ReversalFileSize 字节的单位转换 保留两位小数
func ReversalFileSize(fileSize string) (size float64) {
	if strings.HasSuffix(fileSize, "KB") {
		float, err := strconv.ParseFloat(strings.Replace(fileSize, "KB", "", -1), 64)
		if err != nil {
			return 0
		}
		size = float * float64(1<<10)
	} else if strings.HasSuffix(fileSize, "MB") {
		float, err := strconv.ParseFloat(strings.Replace(fileSize, "MB", "", -1), 64)
		if err != nil {
			return 0
		}
		size = float * float64(1<<20)
	} else if strings.HasSuffix(fileSize, "GB") {
		float, err := strconv.ParseFloat(strings.Replace(fileSize, "GB", "", -1), 64)
		if err != nil {
			return 0
		}
		size = float * float64(1<<30)
	} else if strings.HasSuffix(fileSize, "TB") {
		float, err := strconv.ParseFloat(strings.Replace(fileSize, "TB", "", -1), 64)
		if err != nil {
			return 0
		}
		size = float * float64(1<<40)
	} else if strings.HasSuffix(fileSize, "EB") {
		float, err := strconv.ParseFloat(strings.Replace(fileSize, "EB", "", -1), 64)
		if err != nil {
			return 0
		}
		size = float * float64(1<<50)
	} else if strings.HasSuffix(fileSize, "ZB") {
		float, err := strconv.ParseFloat(strings.Replace(fileSize, "ZB", "", -1), 64)
		if err != nil {
			return 0
		}
		size = float * float64(1<<60)
	} else {
		float, err := strconv.ParseFloat(strings.Replace(fileSize, "B", "", -1), 64)
		if err != nil {
			return 0
		}
		size = float * float64(1)
	}
	return size / float64(1<<20)
}
