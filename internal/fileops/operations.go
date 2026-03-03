package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFileProgress 带进度回调的文件复制
// src: 源文件路径, dst: 目标文件路径, progress: 进度回调（已复制字节数, 总字节数）
func CopyFileProgress(src, dst string, progress func(done, total int64)) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("打开源文件失败: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("获取源文件信息失败: %w", err)
	}
	total := srcInfo.Size()

	// 确保目标目录存在
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %w", err)
	}
	defer dstFile.Close()

	buf := make([]byte, 32*1024)
	var done int64

	for {
		nr, readErr := srcFile.Read(buf)
		if nr > 0 {
			nw, writeErr := dstFile.Write(buf[:nr])
			if writeErr != nil {
				return fmt.Errorf("写入目标文件失败: %w", writeErr)
			}
			done += int64(nw)
			if progress != nil {
				progress(done, total)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return fmt.Errorf("读取源文件失败: %w", readErr)
		}
	}

	// 复制文件权限
	if err := os.Chmod(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("设置文件权限失败: %w", err)
	}

	return nil
}

// CopyDir 递归复制目录
// src: 源目录路径, dst: 目标目录路径, progress: 进度回调
func CopyDir(src, dst string, progress func(done, total int64)) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源目录信息失败: %w", err)
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("读取源目录失败: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath, progress); err != nil {
				return err
			}
		} else {
			if err := CopyFileProgress(srcPath, dstPath, progress); err != nil {
				return err
			}
		}
	}

	return nil
}

// MoveEntry 移动文件或目录（先尝试 rename，跨设备则复制后删除）
// src: 源路径, dst: 目标路径
func MoveEntry(src, dst string) error {
	// 先尝试直接 rename（同设备最快）
	if err := os.Rename(src, dst); err == nil {
		return nil
	}

	// rename 失败则先复制再删除
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("获取源路径信息失败: %w", err)
	}

	if srcInfo.IsDir() {
		if err := CopyDir(src, dst, nil); err != nil {
			return fmt.Errorf("复制目录失败: %w", err)
		}
	} else {
		if err := CopyFileProgress(src, dst, nil); err != nil {
			return fmt.Errorf("复制文件失败: %w", err)
		}
	}

	return os.RemoveAll(src)
}

// DeleteEntry 删除文件或目录（递归）
// path: 要删除的路径
func DeleteEntry(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("删除失败: %w", err)
	}
	return nil
}

// RenameEntry 重命名文件或目录
// oldPath: 原路径, newName: 新文件名（不含路径）
func RenameEntry(oldPath, newName string) error {
	dir := filepath.Dir(oldPath)
	newPath := filepath.Join(dir, newName)

	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("重命名失败: %w", err)
	}
	return nil
}

// CreateDir 创建目录
// parentPath: 父目录路径, name: 新目录名
func CreateDir(parentPath, name string) error {
	newPath := filepath.Join(parentPath, name)
	if err := os.MkdirAll(newPath, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	return nil
}

// GetDirSize 计算目录总大小（递归）
func GetDirSize(path string) (int64, error) {
	var total int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // 忽略无法访问的文件
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})
	return total, err
}
