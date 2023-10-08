package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "gitgrab",
		Usage: "grab folder/file from git",
		Action: func(ctx *cli.Context) error {
			fmt.Println("Gitgrab", ctx.Args().Get(0))
			repoURL := "https://api.github.com/repos/sachinsmc/rest-api-go-sample/tree/master/config"
			token := "YOUR_GITHUB_ACCESS_TOKEN"
			resp, err := sendRequest(repoURL, token)
			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}
			defer resp.Body.Close()

			files, err := parseResponse(resp.Body)

			if err != nil {
				fmt.Println("Error:", err)
				return nil
			}

			for _, file := range files {
				downloadURL := file["download_url"].(string)
				filename := filepath.Base(downloadURL)

				err := downloadFile(downloadURL, filename)
				if err != nil {
					fmt.Println("Error downloading", filename, ":", err)
				} else {
					fmt.Println("Downloaded", filename)
				}
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func sendRequest(url, token string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set the Authorization header with the GitHub token
	req.Header.Set("Authorization", "token "+token)

	client := &http.Client{}
	return client.Do(req)
}

func parseResponse(body io.Reader) ([]map[string]interface{}, error) {
	var files []map[string]interface{}

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func downloadFile(url, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
