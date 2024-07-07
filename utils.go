package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func SaveTree(tree *object.Tree, savePath string) error {
	for _, entry := range tree.Entries {
		path := filepath.Join(savePath, entry.Name)

		switch entry.Mode {
		case filemode.Dir:
			err := os.MkdirAll(path, os.ModePerm)
			CheckIfError(err)
			subTree, err := tree.Tree(entry.Name)
			CheckIfError(err)
			err = SaveTree(subTree, path)
			CheckIfError(err)
		case filemode.Regular, filemode.Executable:
			fileRaw, err := tree.File(entry.Name)
			CheckIfError(err)

			reader, err := fileRaw.Reader()
			CheckIfError(err)
			defer reader.Close()

			file, err := os.Create(path)
			CheckIfError(err)
			defer file.Close()

			_, err = io.Copy(file, reader)
			CheckIfError(err)
		}
	}
	return nil
}

func CheckIfError(err error) {
	if err == nil {
		return
	}
	_, filename, line, _ := runtime.Caller(1)

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("[error] %s:%d %v", filename, line, err))
	os.Exit(1)
}

func splitURL(URL string) (string, string, string) {
	parts := strings.Split(URL, "/")

	repoURL := strings.Join(parts[:5], "/")

	subdirectoryPath := strings.Join(parts[7:], "/")
	ref := parts[6]

	return repoURL, subdirectoryPath, ref
}
