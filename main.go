package main

import (
	"fmt"

	"github.com/teihenn/prtime/internal/config"
)

func main() {
	configFilePath := "/Users/y_yoshida/Projects/prtime/internal/config/config.yml"
	config, err := config.Load(configFilePath)
	if err != nil {
		fmt.Printf("Error loading config: %s\n", err)
		return
	}
	fmt.Printf("Loaded config: %+v\n", config)
}
