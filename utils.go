package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
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

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("Please provide a valid folder url [error] %s:%d %v", filename, line, err))
	os.Exit(1)
}

func DownloadFolder(URL string) error {
	repoURL, subdirectoryPath, refName := splitURL(URL)
	savePath := "./" + subdirectoryPath

	storage := memory.NewStorage()

	repo, err := git.Clone(storage, nil, &git.CloneOptions{
		URL:           repoURL,
		SingleBranch:  true,
		ReferenceName: plumbing.ReferenceName(refName),
		Depth:         1, // Shallow clone with depth 1.
	})

	CheckIfError(err)

	ref, err := repo.Head()
	CheckIfError(err)

	commit, err := repo.CommitObject(ref.Hash())
	CheckIfError(err)

	tree, err := commit.Tree()
	CheckIfError(err)

	subTree, err := tree.Tree(subdirectoryPath)
	CheckIfError(err)

	err = os.Mkdir(savePath, os.ModePerm)
	CheckIfError(err)

	err = SaveTree(subTree, savePath)
	CheckIfError(err)

	fmt.Println("Subdirectory saved : ", savePath)
	return nil
}

func splitURL(URL string) (string, string, string) {
	parts := strings.Split(URL, "/")

	repoURL := strings.Join(parts[:5], "/")

	subdirectoryPath := strings.Join(parts[7:], "/")
	ref := parts[6]

	return repoURL, subdirectoryPath, ref
}
