// Package helper
// @program: gin-template
// @author: [lliuhuan](https://github.com/lliuhuan)
// @create: 2024-07-12 14:59
package helper

import (
	"mime/multipart"
	"path/filepath"
	"strconv"

	"github.com/LLiuHuan/gin-template/internal/pkg/core"
	"github.com/LLiuHuan/gin-template/pkg/file"
)

type UploadFileData struct {
	File      *multipart.FileHeader
	FileHash  string
	FileIndex int
	FileStart int
	FileEnd   int
}

// UploadFile 保存分块文件
func (s *service) UploadFile(ctx core.Context, data *UploadFileData) error {
	// 检查临时文件存不存在，不存在则创建
	dir := filepath.Join(DefaultBreakpointPath, data.FileHash)
	if err := file.MkdirAll(dir); err != nil {
		return err
	}

	// 检查分块文件是否存在，不存在则创建
	dst := filepath.Join(dir, strconv.Itoa(data.FileIndex)+"_"+data.FileHash)
	// 检查一下文件是否存在
	if _, b := file.IsExists(dst); b {
		return nil
	}
	_, err := file.SaveFile(dst, data.File)
	if err != nil {
		return err
	}

	return nil
}
