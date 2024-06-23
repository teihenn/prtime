package service

import "github.com/teihenn/prtime/internal/model"

// TODO: How to make model common between GitHub and Bitbucket

// PullRequestService defines the interface for operations related to pull requests.
type PullRequestService interface {
	GetPullRequests(owner, repo string) ([]model.PullRequest, error)
	GetPullRequestDetails(owner, repo string, prID int) (model.PullRequest, error)
}
