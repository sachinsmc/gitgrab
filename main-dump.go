package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/urfave/cli/v2"
)

func notmain() {
	app := &cli.App{
		Name:  "gitgrab",
		Usage: "grab folder/file from git",
		Action: func(ctx *cli.Context) error {
			fmt.Println("Gitgrab", ctx.Args().Get(0))

			repoURL := "https://github.com/sachinsmc/rest-api-go-sample/tree/master/models"
			subdirectoryPath := "models"
			savePath := "/Users/sachin/go/src/github.com/sachinsmc/gitgrab/models"
			// downloadDir := "/Users/sachin/Github/test-npm-pkg/temp"

			// Create a new memory storage.
			storage := memory.NewStorage()

			// Clone the entire repository.
			repo, err := git.Clone(storage, nil, &git.CloneOptions{
				URL:          repoURL,
				SingleBranch: true,
				// ReferenceName: plumbing.ReferenceName("refs/heads/master"),
				Depth: 1, // Shallow clone with depth 1.
			})
			fmt.Println("cloning")

			CheckIfError(err)

			// ... retrieves the branch pointed by HEAD
			ref, err := repo.Head()
			CheckIfError(err)
			// fmt.Println(ref)
			// fmt.Println("storage", storage.Trees)

			// Get the commit from the reference.
			commit, err := repo.CommitObject(ref.Hash())
			CheckIfError(err)

			// Get the tree object for the commit.
			tree, err := commit.Tree()
			CheckIfError(err)
			// fmt.Println(tree)

			// Get the subdirectory tree.
			subTree, err := tree.Tree(subdirectoryPath)
			CheckIfError(err)

			os.Mkdir(savePath, os.ModePerm)

			// Recursively save the subdirectory.
			err = SaveTree(subTree, savePath)
			CheckIfError(err)

			fmt.Println("Subdirectory saved")

			// Check out only the specific subdirectory.
			// tree, err := getSubdirectory(repo, subdirectoryPath)
			// if err != nil {
			// 	fmt.Printf("Error retrieving the subdirectory: %v\n", err)
			// 	os.Exit(1)
			// }

			// // Print the contents of the subdirectory.
			// for _, entry := range tree.Entries {
			// 	fmt.Println(entry.Name)
			// }
			// CheckIfError(err)

			// saveFiles(repo, tree, downloadDir)

			// Now, you have the "models" directory cloned in memory. You can access its contents as needed.

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

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

func getSubdirectory(repo *git.Repository, subdirectoryPath string) (*object.Tree, error) {
	// Get the HEAD reference.
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	// Get the commit from the reference.
	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	// Get the tree object for the commit.
	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	// Traverse the tree to the specified subdirectory.
	for _, segment := range splitPath(subdirectoryPath) {
		entry, err := tree.FindEntry(segment)
		if err != nil {
			return nil, err
		}

		if !entry.Mode.IsFile() {
			tree, err = repo.TreeObject(entry.Hash)
			if err != nil {
				return nil, err
			}
		}
	}

	return tree, nil
}

func splitPath(path string) []string {
	return strings.Split(path, "/")
}

func saveFiles(repo *git.Repository, tree *object.Tree, downloadDir string) {
	for _, entry := range tree.Entries {
		if !entry.Mode.IsFile() {
			continue
		}

		_, err := repo.BlobObject(entry.Hash)
		if err != nil {
			fmt.Printf("Error retrieving file: %v\n", err)
			continue
		}

		// fileContents, err := file.Contents()
		// if err != nil {
		// 	fmt.Printf("Error reading file contents: %v\n", err)
		// 	continue
		// }

		// filePath := fmt.Sprintf("%s/%s", downloadDir, entry.Name)

		// err = os.MkdirAll(path.Dir(filePath), os.ModePerm)
		// if err != nil {
		// 	fmt.Printf("Error creating directories: %v\n", err)
		// 	continue
		// }

		// err = os.WriteFile(filePath, []byte(fileContents), 0644)
		// if err != nil {
		// 	fmt.Printf("Error saving file: %v\n", err)
		// }
	}
}
