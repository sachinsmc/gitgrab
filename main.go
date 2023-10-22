package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gitgrab",
		Usage: "grab folder/file from git",
		Action: func(ctx *cli.Context) error {
			fmt.Println("Gitgrab", ctx.Args().Get(0))

			repoURL := "https://github.com/sachinsmc/rest-api-go-sample.git"
			subdirectoryPath := "models"
			downloadDir := "/Users/sachin/Github/test-npm-pkg/temp"

			// Create a new memory storage.
			storage := memory.NewStorage()

			// Clone the entire repository.
			repo, err := git.Clone(storage, nil, &git.CloneOptions{
				URL:           repoURL,
				SingleBranch:  true,
				ReferenceName: plumbing.ReferenceName("refs/heads/master"),
				Depth:         1, // Shallow clone with depth 1.
			})

			if err != nil {
				fmt.Printf("Error cloning the repository: %v\n", err)
				os.Exit(1)
			}

			// Check out only the specific subdirectory.
			tree, err := getSubdirectory(repo, subdirectoryPath)
			if err != nil {
				fmt.Printf("Error retrieving the subdirectory: %v\n", err)
				os.Exit(1)
			}

			// Print the contents of the subdirectory.
			for _, entry := range tree.Entries {
				fmt.Println(entry.Name)
			}
			CheckIfError(err)

			saveFiles(repo, tree, downloadDir)

			// Now, you have the "models" directory cloned in memory. You can access its contents as needed.

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
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

		file, err := repo.BlobObject(entry.Hash)
		if err != nil {
			fmt.Printf("Error retrieving file: %v\n", err)
			continue
		}

		fileContents, err := file.Contents()
		if err != nil {
			fmt.Printf("Error reading file contents: %v\n", err)
			continue
		}

		filePath := fmt.Sprintf("%s/%s", downloadDir, entry.Name)

		err = os.MkdirAll(path.Dir(filePath), os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating directories: %v\n", err)
			continue
		}

		err = os.WriteFile(filePath, []byte(fileContents), 0644)
		if err != nil {
			fmt.Printf("Error saving file: %v\n", err)
		}
	}
}
