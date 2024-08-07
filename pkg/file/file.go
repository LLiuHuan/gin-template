// Package file
//
//	@program:		gin-template
//	@author:		[lliuhuan](https://github.com/lliuhuan)
//	@create:		2024-07-02 21:57
//	@description:	文件操作
package file

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

var (
	buffSize = 1 << 20 // 1MB
)

// ReadLineFromEnd 从文件末尾读取行
type ReadLineFromEnd struct {
	f *os.File

	fileSize int
	bwr      *bytes.Buffer
	lineBuff []byte
	swapBuff []byte

	isFirst bool
}

// SaveFile 保存文件
func SaveFile(path string, file *multipart.FileHeader) (int64, error) {
	f, err := file.Open()
	if err != nil {
		return 0, err
	}
	defer f.Close()

	dst, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, f)
}

// IsExists 文件是否存在
func IsExists(path string) (os.FileInfo, bool) {
	f, err := os.Stat(path)
	return f, err == nil || os.IsExist(err)
}

// MkdirAll 创建文件夹
func MkdirAll(path string) error {
	if _, b := IsExists(path); !b {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// FileSize 文件大小
func FileSize(path string) (int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	stat, err := f.Stat() //获取文件状态
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return stat.Size(), nil
}

// FindAllFile 查找所有文件
func FindAllFile(path string) (s []string, err error) {
	s = make([]string, 0)
	rd, err := os.ReadDir(path)
	if err != nil {
		return s, err
	}
	for _, file := range rd {
		if !file.IsDir() {
			s = append(s, file.Name())
		}
	}
	return s, nil
}

// NewReadLineFromEnd 从末尾读取文件内容
func NewReadLineFromEnd(filename string) (rd *ReadLineFromEnd, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	info, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("not file")
	}

	fileSize := int(info.Size())

	rd = &ReadLineFromEnd{
		f:        f,
		fileSize: fileSize,
		bwr:      bytes.NewBuffer([]byte{}),
		lineBuff: make([]byte, 0),
		swapBuff: make([]byte, buffSize),
		isFirst:  true,
	}
	return rd, nil
}

// ReadLine 结尾包含'\n'
func (c *ReadLineFromEnd) ReadLine() (line []byte, err error) {
	var ok bool
	for {
		ok, err = c.buff()
		if err != nil {
			return nil, err
		}
		if ok {
			break
		}
	}
	line, err = c.bwr.ReadBytes('\n')
	if err == io.EOF && c.fileSize > 0 {
		err = nil
	}
	return line, err
}

// Close 关闭文件
func (c *ReadLineFromEnd) Close() (err error) {
	return c.f.Close()
}

// buff 读取文件内容到缓冲区
func (c *ReadLineFromEnd) buff() (ok bool, err error) {
	if c.fileSize == 0 {
		return true, nil
	}

	if c.bwr.Len() >= buffSize {
		return true, nil
	}

	offset := 0
	if c.fileSize > buffSize {
		offset = c.fileSize - buffSize
	}
	_, err = c.f.Seek(int64(offset), 0)
	if err != nil {
		return false, err
	}

	n, err := c.f.Read(c.swapBuff)
	if err != nil && err != io.EOF {
		return false, err
	}
	if c.fileSize < n {
		n = c.fileSize
	}
	if n == 0 {
		return true, nil
	}

	for {
		m := bytes.LastIndex(c.swapBuff[:n], []byte{'\n'})
		if m == -1 {
			break
		}
		if m < n-1 {
			err = c.writeLine(c.swapBuff[m+1 : n])
			if err != nil {
				return false, err
			}
			ok = true
		} else if m == n-1 && !c.isFirst {
			err = c.writeLine(nil)
			if err != nil {
				return false, err
			}
			ok = true
		}
		n = m
		if n == 0 {
			break
		}
	}
	if n > 0 {
		reverseBytes(c.swapBuff[:n])
		c.lineBuff = append(c.lineBuff, c.swapBuff[:n]...)
	}
	if offset == 0 {
		err = c.writeLine(nil)
		if err != nil {
			return false, err
		}
		ok = true
	}
	c.fileSize = offset
	if c.isFirst {
		c.isFirst = false
	}
	return ok, nil
}

// writeLine 写入行
func (c *ReadLineFromEnd) writeLine(b []byte) (err error) {
	if len(b) > 0 {
		_, err = c.bwr.Write(b)
		if err != nil {
			return err
		}
	}
	if len(c.lineBuff) > 0 {
		reverseBytes(c.lineBuff)
		_, err = c.bwr.Write(c.lineBuff)
		if err != nil {
			return err
		}
		c.lineBuff = c.lineBuff[:0]
	}
	_, err = c.bwr.Write([]byte{'\n'})
	if err != nil {
		return err
	}
	return nil
}

// reverseBytes 反转字节
func reverseBytes(b []byte) {
	n := len(b)
	if n <= 1 {
		return
	}
	for i := 0; i < n; i++ {
		k := n - 1
		if k != i {
			b[i], b[k] = b[k], b[i]
		}
		n--
	}
}
