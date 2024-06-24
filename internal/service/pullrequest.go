package service

import "github.com/teihenn/prtime/internal/model"

// TODO: How to make model common between GitHub and Bitbucket

// PullRequestService defines the interface for operations related to pull requests.
type PullRequestService interface {
	GetSortedPRs() map[string][]model.PullRequest
	DisplaySortedPRs()
}
