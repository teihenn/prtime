package model

import "time"

type AuthorInfo struct {
	Username string `json:"display_name"`
}

// PullRequest represents the information of a pull request from GitHub or Bitbucket.
type PullRequest struct {
	ID           int        `json:"id"`
	Title        string     `json:"title"`
	Author       AuthorInfo `json:"author"`
	CreatedOn    time.Time  `json:"created_on"`
	State        string     `json:"state"`
	CommentCount int        `json:"comment_count"`
}
