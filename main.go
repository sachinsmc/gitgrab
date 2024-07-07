package main

import (
	"fmt"
	"log"
	"os"

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

			err := DownloadFolder(URL)
			CheckIfError(err)
			return err
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
