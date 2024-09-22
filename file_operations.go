package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const rootDir = "./nodian"

var currentDirectory string

func init() {
	// 从配置文件或其他持久化存储中读取上次保存的目录
	// 这里简单起见，我们使用一个环境变量
	currentDirectory = os.Getenv("NODIAN_ROOT_DIR")
	if currentDirectory == "" {
		currentDirectory = "./nodian"
	}
}

func (a *App) CreateMarkdownFile(filename string, content string) error {
	if !strings.HasPrefix(filename, rootDir) {
		filename = filepath.Join(rootDir, filename)
	}
	if filepath.Ext(filename) != ".mk" {
		filename += ".mk"
	}
	return os.WriteFile(filename, []byte(content), 0644)
}

func (a *App) ReadMarkdownFile(filename string) (string, error) {
	if !strings.HasPrefix(filename, rootDir) {
		filename = filepath.Join(rootDir, filename)
	}
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (a *App) SaveMarkdownFile(filename string, content string) error {
	if !strings.HasPrefix(filename, rootDir) {
		filename = filepath.Join(rootDir, filename)
	}
	return os.WriteFile(filename, []byte(content), 0644)
}

func (a *App) ListMarkdownFiles(directory string) ([]string, error) {
	if directory == "" || directory == "." {
		directory = currentDirectory
	}
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil // 跳过根目录
		}
		if info.IsDir() {
			files = append(files, relPath+"/")
		} else if filepath.Ext(path) == ".mk" {
			files = append(files, relPath)
		}
		return nil
	})
	log.Printf("Listed files: %v", files)
	return files, err
}

func (a *App) CreateFolder(path string) error {
	if !strings.HasPrefix(path, rootDir) {
		path = filepath.Join(rootDir, path)
	}
	return os.MkdirAll(path, 0755)
}

func (a *App) RenameItem(oldPath, newPath string) error {
	if !strings.HasPrefix(oldPath, rootDir) {
		oldPath = filepath.Join(rootDir, oldPath)
	}
	if !strings.HasPrefix(newPath, rootDir) {
		newPath = filepath.Join(rootDir, newPath)
	}
	return os.Rename(oldPath, newPath)
}

func (a *App) SelectDirectory() (string, error) {
	selected, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Root Directory",
	})
	if err != nil {
		return "", err
	}
	if selected != "" {
		currentDirectory = selected
		// 保存选择的目录到配置文件或其他持久化存储
		// 这里简单起见，我们使用一个环境变量
		os.Setenv("NODIAN_ROOT_DIR", currentDirectory)
	}
	return currentDirectory, nil
}

func (a *App) GetCurrentDirectory() string {
	return currentDirectory
}

func (a *App) SetCurrentDirectory(dir string) {
	currentDirectory = dir
	// 保存选择的目录到配置文件或其他持久化存储
	// 这里简单起见，我们使用一个环境变量
	os.Setenv("NODIAN_ROOT_DIR", currentDirectory)
}
