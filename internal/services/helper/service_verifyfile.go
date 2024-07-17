// Package helper
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-17 23:44
package helper

import (
	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/pkg/file"
	"path"
	"strconv"
	"strings"
)

type VerifyFileData struct {
	Hash string
	Name string
}

type VerifyFileResponse struct {
	ExistFile   bool
	ExistChunks []int
}

// VerifyFile 合并分块文件
func (s *service) VerifyFile(ctx core.Context, data *VerifyFileData) (VerifyFileResponse, error) {
	res := new(VerifyFileResponse)

	// 判断是否存在该hash的文件
	// 最终文件保存路径
	dst := path.Join(DefaultFileSavePath, data.Hash+path.Ext(data.Name))
	// 读取临时文件夹
	dir := path.Join(DefaultBreakpointPath, data.Hash)

	// 检查最终文件是否存在
	if _, ok := file.IsExists(dst); ok {
		res.ExistFile = true
		return *res, nil
	}

	if _, ok := file.IsExists(dir); !ok {
		return *res, nil
	}

	// 检查临时文件夹下的所有文件
	files, err := file.FindAllFile(dir)
	if err != nil {
		res.ExistChunks = []int{}
		return *res, err
	}

	for _, fileName := range files {
		index := strings.Split(fileName, "_")[0]
		i, _ := strconv.Atoi(index)
		res.ExistChunks = append(res.ExistChunks, i)
	}

	return *res, err
}
