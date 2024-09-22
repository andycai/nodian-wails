package main

import (
	"os"
	"path/filepath"
)

func (a *App) CreateMarkdownFile(filename string, content string) error {
	if filepath.Ext(filename) != ".mk" {
		filename += ".mk"
	}
	return os.WriteFile(filename, []byte(content), 0644)
}

func (a *App) ReadMarkdownFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (a *App) SaveMarkdownFile(filename string, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func (a *App) ListMarkdownFiles(directory string) ([]string, error) {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".mk" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
