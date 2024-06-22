package main

import (
	"fmt"

	"github.com/teihenn/prtime/internal/config"
	"github.com/teihenn/prtime/internal/service"
)

func main() {
	// TODO: make configFilePath to relative path from project root
	configFilePath := "/Users/y_yoshida/Projects/prtime/internal/config/config.yml"
	config, err := config.Load(configFilePath)
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		return
	}

	bitbucketService, err := service.NewBitbucketService(*config)
	if err != nil {
		fmt.Printf("Error creating bitbucket service: %s\n", err)
		return
	}

	for _, repo := range bitbucketService.Repositories {
		prs, err := bitbucketService.ListPullRequests(repo.Owner, repo.Name)
		if err != nil {
			fmt.Printf("Error retrieving pull requests for %s/%s: %s\n", repo.Owner, repo.Name, err)
			continue
		}

		fmt.Printf("Pull Requests for %s/%s:\n", repo.Owner, repo.Name)
		for _, pr := range prs {
			fmt.Printf("- %d: %s by %s\n, created at %v", pr.ID, pr.Title, pr.Author.Username, pr.CreatedOn)
		}
	}
}
