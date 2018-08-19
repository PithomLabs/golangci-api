package repoanalyzes

import (
	"fmt"
	"time"

	"github.com/golangci/golangci-api/app/internal/db"
	"github.com/golangci/golangci-api/app/internal/errors"
	"github.com/golangci/golangci-api/app/models"
	"github.com/golangci/golangci-api/app/utils"
	"github.com/golangci/golib/server/context"
)

func reanalyzeByNewLinters() {
	ctx := utils.NewBackgroundContext()
	analysisStatusesCh := make(chan models.RepoAnalysisStatus, 1024)
	go reanalyzeFromCh(ctx, analysisStatusesCh)

	for range time.NewTicker(time.Minute).C {
		var analysisStatuses []models.RepoAnalysisStatus
		err := models.NewRepoAnalysisStatusQuerySet(db.Get(ctx)).
			LastAnalyzedLintersVersionNe(lintersVersion).
			HasPendingChangesEq(false).
			Limit(100).
			All(&analysisStatuses)
		if err != nil {
			errors.Warnf(ctx, "Can't fetch analysis statuses")
			continue
		}
		if len(analysisStatuses) == 0 {
			ctx.L.Infof("No analysis statuses to reanalyze by new linters")
			break
		}

		ctx.L.Infof("Fetched %d analysis statuses to reanalyze by new linters", len(analysisStatuses))

		for _, as := range analysisStatuses {
			analysisStatusesCh <- as
		}
	}

	close(analysisStatusesCh)
}

func reanalyzeFromCh(ctx *context.C, analysisStatusesCh <-chan models.RepoAnalysisStatus) {
	const avgAnalysisTime = time.Minute
	const maxReanalyzeCapacity = 0.5
	reanalyzeInterval := time.Duration(float64(avgAnalysisTime) / maxReanalyzeCapacity)

	for as := range analysisStatusesCh {
		if err := reanalyzeAnalysisByNewLinters(ctx, &as); err != nil {
			errors.Warnf(ctx, "Can't reanalyze analysis status %#v: %s", as, err)
		}
		time.Sleep(reanalyzeInterval)
	}
}

func reanalyzeAnalysisByNewLinters(ctx *context.C, as *models.RepoAnalysisStatus) error {
	var a models.RepoAnalysis
	err := models.NewRepoAnalysisQuerySet(db.Get(ctx)).
		RepoAnalysisStatusIDEq(as.ID).
		OrderDescByID().
		One(&a)
	if err != nil {
		return fmt.Errorf("can't fetch last repo analysis for %s: %s", as.Name, err)
	}

	if as.LastAnalyzedLintersVersion == "" {
		err = models.NewRepoAnalysisStatusQuerySet(db.Get(ctx)).
			IDEq(as.ID).
			GetUpdater().
			SetLastAnalyzedLintersVersion(a.LintersVersion).
			Update()
		if err != nil {
			return fmt.Errorf("can't set last_analyzed_linters_version to %s", a.LintersVersion)
		}

		// send it to the next iteration to prevent extra checks here
		return nil
	}

	err = models.NewRepoAnalysisStatusQuerySet(db.Get(ctx)).
		IDEq(as.ID).
		GetUpdater().
		SetHasPendingChanges(true).
		SetPendingCommitSHA(a.CommitSHA).
		Update()
	if err != nil {
		return fmt.Errorf("can't update has_pending_changes to true: %s", err)
	}

	ctx.L.Infof("Marked repo %s for reanalysis by new linters: %s -> %s",
		as.Name, as.LastAnalyzedLintersVersion, lintersVersion)
	return nil
}