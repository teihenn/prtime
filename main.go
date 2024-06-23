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

	bitbucketService.DisplaySortedPRs()
}
