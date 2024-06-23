package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/teihenn/prtime/internal/config"
	"github.com/teihenn/prtime/internal/model"
)

// BitbucketService represents the client for Bitbucket API operations.
type BitbucketService struct {
	Repositories []struct {
		Owner  string
		Name   string
		ApiKey string
	}
}

const API_VERSION = "2.0"

// NewBitbucketService creates a new instance of BitbucketService.
func NewBitbucketService(cfg config.Config) (*BitbucketService, error) {
	if len(cfg.Bitbucket.Repositories) == 0 {
		return nil, fmt.Errorf("no repositories provided")
	}

	repos := make([]struct {
		Owner  string
		Name   string
		ApiKey string
	}, len(cfg.Bitbucket.Repositories))

	for i, repo := range cfg.Bitbucket.Repositories {
		apiKey, err := loadApiKey(cfg.Bitbucket.AllRepoApiKeyPath, repo.Name)
		if err != nil {
			return nil, err
		}
		repos[i] = struct {
			Owner  string
			Name   string
			ApiKey string
		}{
			Owner:  repo.Owner,
			Name:   repo.Name,
			ApiKey: apiKey,
		}
	}

	return &BitbucketService{
		Repositories: repos,
	}, nil
}

func loadApiKey(apiKeyPath, repoName string) (string, error) {
	data, err := os.ReadFile(apiKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed to read API key file: %v", err)
	}

	var keys map[string]string
	if err := json.Unmarshal(data, &keys); err != nil {
		return "", fmt.Errorf("failed to parse API key file: %v", err)
	}

	apiKey, ok := keys[repoName]
	if !ok {
		return "", fmt.Errorf("API key for repository %s not found", repoName)
	}

	return apiKey, nil
}

// GetPullRequests retrieves a list of pull requests from a specific repository.
func (b *BitbucketService) GetPullRequests(owner, repo string) ([]model.PullRequest, error) {
	var apiKey string
	for _, r := range b.Repositories {
		if r.Owner == owner && r.Name == repo {
			apiKey = r.ApiKey
			break
		}
	}
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found for repository %s/%s", owner, repo)
	}

	url := fmt.Sprintf("https://api.bitbucket.org/%s/repositories/%s/%s/pullrequests",
		API_VERSION, owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer"+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data struct {
		Values []model.PullRequest `json:"values"`
	}

	log.Printf("%v", data)

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data.Values, nil
}

// GetPullRequestDetails fetches details for a specific pull request.
func (b *BitbucketService) GetPullRequestDetails(owner, repo string, prID int) (model.PullRequest, error) {
	// Call Bitbucket API to fetch detailed data for a pull request
	pr := model.PullRequest{}
	return pr, nil
}
