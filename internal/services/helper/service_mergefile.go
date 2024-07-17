// Package helper
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-12 20:57
package helper

import (
	"fmt"
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/pkg/file"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type MergeFileData struct {
	Hash  string `form:"hash"`
	Name  string `form:"name"`
	Count int
}

// MergeFile 合并分块文件
func (s *service) MergeFile(ctx core.Context, data *MergeFileData) error {
	start := time.Now().UnixMicro()

	// 最终文件保存路径
	dst := path.Join(DefaultFileSavePath, data.Hash+path.Ext(data.Name))
	// 读取临时文件夹
	dir := path.Join(DefaultBreakpointPath, data.Hash)
	// 读取临时文件夹下的所有文件
	files, err := file.FindAllFile(dir)
	if err != nil {
		return err
	}

	// 检查保存的目录是否存在
	if err = file.MkdirAll(DefaultFileSavePath); err != nil {
		return err
	}

	if _, ok := file.IsExists(dst); ok {
		return fmt.Errorf("文件已存在")
	}

	fmt.Println(data.Count, len(files))
	if data.Count != len(files) {
		return fmt.Errorf("文件不完整")
	}

	// 名称数组可能不是按顺序的，所以需要排序
	sort.Slice(files, func(i, j int) bool {
		iIndex, _ := strconv.Atoi(strings.Split(files[i], "_")[0])
		jIndex, _ := strconv.Atoi(strings.Split(files[j], "_")[0])
		return iIndex < jIndex
	})

	// 创建最终文件
	create, err := os.Create(dst)
	if err != nil {
		return err
	}

	// 循环把分片文件写入最终文件
	for _, fileName := range files {
		if fileName == ".BD_Store" {
			continue
		}
		bytes, err := os.ReadFile(path.Join(dir, fileName))
		if err != nil {
			return err
		}

		_, err = create.Write(bytes)
		if err != nil {
			return err
		}
	}

	// 删除临时文件夹
	err = os.RemoveAll(dir)

	end := time.Now().UnixMicro()
	timeSend := end - start
	// 计算耗时多少秒
	fmt.Printf("耗时：%d秒\n", timeSend/1e6)
	return err
}
