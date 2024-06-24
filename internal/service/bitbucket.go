package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"time"

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

func (b *BitbucketService) GetSortedPRs() map[string][]model.PullRequest {
	var ret map[string][]model.PullRequest = make(map[string][]model.PullRequest)
	for _, repo := range b.Repositories {
		prs, err := b.getPullRequests(repo.Owner, repo.Name)
		if err != nil {
			fmt.Printf("Error retrieving pull requests for %s/%s: %s\n",
				repo.Owner, repo.Name, err)
			continue
		}

		sort.Slice(prs, func(i, j int) bool {
			return time.Since(prs[i].CreatedOn) > time.Since((prs[j].CreatedOn))
		})

		key := repo.Owner + "/" + repo.Name
		ret[key] = prs
	}
	return ret
}

// DisplaySortedPRs displays pull requests sorted by their open duration in descending order.
func (b *BitbucketService) DisplaySortedPRs() {
	jst, _ := time.LoadLocation(("Asia/Tokyo"))

	reposPRs := b.GetSortedPRs()

	for _, repo := range b.Repositories {
		fmt.Printf("Pull Requests for %s/%s:\n", repo.Owner, repo.Name)
		key := repo.Owner + "/" + repo.Name
		prs := reposPRs[key]
		for _, pr := range prs {
			openDuration := time.Since(pr.CreatedOn)
			daysOpen := openDuration.Hours() / 24
			createdOnFormattedJst := pr.CreatedOn.In(jst).Format("2006-01-02 15:04")
			fmt.Printf("[%d] %s by %s, Open for %.1f days(created on %s)\n",
				pr.ID, pr.Title, pr.Author.Username, daysOpen, createdOnFormattedJst)
		}
	}
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

// getPullRequests retrieves a list of pull requests from a specific repository.
func (b *BitbucketService) getPullRequests(owner, repo string) ([]model.PullRequest, error) {
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

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data.Values, nil
}
