package main

import (
	"os"
	"path/filepath"
	"strings"
)

const rootDir = "./assets"

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
		directory = rootDir
	}
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			files = append(files, relPath+"/")
		} else if filepath.Ext(path) == ".mk" {
			files = append(files, relPath)
		}
		return nil
	})
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
