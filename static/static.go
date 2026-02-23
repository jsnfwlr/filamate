// Package static embeds static files into the binary.
package static

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed assets/*
//go:embed *.ico
//go:embed *.html
//go:embed *.png
var Files embed.FS

type FS interface {
	ReadDir(name string) (folderContents []fs.FileInfo, fault error)
	ReadFile(name string) (contents []byte, fault error)
	Glob(pattern string) (matches []string, fault error)
	Open(name string) (file fs.File, fault error)
}

type Embedded struct {
	filesystem embed.FS
}

func NewEmbedded() (fs Embedded, fault error) {
	return Embedded{
		filesystem: Files,
	}, nil
}

func (e Embedded) ReadDir(name string) (folderContents []fs.FileInfo, fault error) {
	files, err := e.filesystem.ReadDir(name)
	if err != nil {
		return nil, fmt.Errorf("could not get the files from the embedded filesystem: %w", err)
	}

	var r []os.FileInfo

	for _, f := range files {
		fi, _ := f.Info()

		r = append(r, fi)
	}

	return r, nil
}

func (e Embedded) ReadFile(name string) (contents []byte, fault error) {
	b, err := e.filesystem.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s' from embedded filesystem: %w", name, err)
	}

	return b, nil
}

func (e Embedded) Glob(pattern string) (matches []string, fault error) {
	matches, err := fs.Glob(e.filesystem, pattern)
	if err != nil {
		return nil, fmt.Errorf("could not get glob matches for pattern '%s': %w", pattern, err)
	}

	return matches, nil
}

func (e Embedded) Open(name string) (file fs.File, fault error) {
	f, err := e.filesystem.Open(name)
	if err != nil {
		return nil, os.ErrNotExist
	}

	return f, nil
}

type Directory struct {
	path string
}

func (c Directory) ReadDir(name string) (folderContents []fs.FileInfo, fault error) {
	files, err := os.ReadDir(filepath.Join(c.path, name))
	if err != nil {
		return nil, fmt.Errorf("could not get the files from the directory filesystem: %w", err)
	}

	var r []os.FileInfo

	for _, f := range files {
		fi, _ := f.Info()

		r = append(r, fi)
	}

	return r, nil
}

func (c Directory) ReadFile(name string) (contents []byte, fault error) {
	b, err := os.ReadFile(filepath.Join(c.path, name))
	if err != nil {
		return nil, fmt.Errorf("could not read file '%s' from directory filesystem: %w", name, err)
	}

	return b, nil
}

func (c Directory) Glob(pattern string) (matches []string, fault error) {
	matches, err := fs.Glob(os.DirFS(c.path), pattern)
	if err != nil {
		return nil, fmt.Errorf("could not get glob matches for pattern '%s': %w", pattern, err)
	}

	return matches, nil
}

func (c Directory) Open(name string) (file fs.File, fault error) {
	f, err := os.Open(filepath.Join(c.path, name))
	if err != nil {
		return nil, os.ErrNotExist
	}

	return f, nil
}

func NewDirectory(path string) (fs Directory, fault error) {
	root, err := os.Getwd()
	if err != nil {
		return Directory{}, fmt.Errorf("could not get current working directory: %w", err)
	}
	d := Directory{
		path: filepath.Join(root, path),
	}
	return d, nil
}
