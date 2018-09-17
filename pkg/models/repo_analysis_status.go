package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//go:generate goqueryset -in repo_analysis_status.go

// gen:qs
type RepoAnalysisStatus struct {
	gorm.Model

	Repo   Repo
	RepoID uint

	Name string // TODO: remove

	LastAnalyzedAt             time.Time
	LastAnalyzedLintersVersion string

	HasPendingChanges bool
	PendingCommitSHA  string
	Version           int
	DefaultBranch     string

	Active bool
}
