package models

import "time"

type SearchResult struct {
	RepoId          string    `json:"repository_id" db:"id"`
	Username        string    `json:"username" db:"username"`
	RepoName        string    `json:"repo_name" db:"name"`
	RepoDescription string    `json:"description,omitempty" db:"description"`
	NumTags         int       `json:"num_tags" db:"num_tags"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
