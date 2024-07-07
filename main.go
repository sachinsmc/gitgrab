package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gitgrab",
		Usage: "grab folder from git",
		Action: func(ctx *cli.Context) error {
			fmt.Println("Gitgrab : ", ctx.Args().Get(0))

			URL := ctx.Args().Get(0)
			if URL == "" {
				fmt.Println("Error Please provide a folder URL \n example usage : gitgrab https://github.com/go-git/go-git/blob/master/config")
				return nil
			}

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
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
