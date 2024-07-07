package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
			if err != nil {
				return err
			}
			subTree, err := tree.Tree(entry.Name)
			if err != nil {
				return err
			}
			err = SaveTree(subTree, path)
			if err != nil {
				return err
			}
		case filemode.Regular, filemode.Executable:
			fileRaw, err := tree.File(entry.Name)
			if err != nil {
				return err
			}

			reader, err := fileRaw.Reader()
			if err != nil {
				return err
			}
			defer reader.Close()

			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(file, reader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func splitURL(URL string) (string, string, string) {
	parts := strings.Split(URL, "/")

	repoURL := strings.Join(parts[:5], "/")

	subdirectoryPath := strings.Join(parts[7:], "/")
	ref := parts[6]

	return repoURL, subdirectoryPath, ref
}
