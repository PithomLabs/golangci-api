package consumers

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golangci/golangci-api/pkg/worker/analytics"
	"github.com/golangci/golangci-api/pkg/worker/analyze/processors"
	"github.com/golangci/golangci-api/pkg/worker/lib/experiments"
	"github.com/golangci/golangci-api/pkg/worker/lib/github"
	"github.com/pkg/errors"
)

type AnalyzeRepo struct {
	baseConsumer

	ec  *experiments.Checker
	rpf *processors.RepoProcessorFactory
}

func NewAnalyzeRepo(ec *experiments.Checker, rpf *processors.RepoProcessorFactory) *AnalyzeRepo {
	return &AnalyzeRepo{
		baseConsumer: baseConsumer{
			eventName: analytics.EventRepoAnalyzed,
		},
		ec:  ec,
		rpf: rpf,
	}
}

func (c AnalyzeRepo) Consume(ctx context.Context, repoName, analysisGUID, branch, privateAccessToken string) error {
	ctx = c.prepareContext(ctx, map[string]interface{}{
		"repoName":     repoName,
		"provider":     "github",
		"analysisGUID": analysisGUID,
		"branch":       branch,
	})

	if os.Getenv("DISABLE_REPO_ANALYSIS") == "1" {
		analytics.Log(ctx).Warnf("Repo analysis is disabled, return error to try it later")
		return errors.New("repo analysis is disabled")
	}

	return c.wrapConsuming(ctx, func() error {
		var cancel context.CancelFunc
		// If you change timeout value don't forget to change it
		// in golangci-api stale analyzes checker
		ctx, cancel = context.WithTimeout(ctx, 10*time.Minute)
		defer cancel()

		return c.analyzeRepo(ctx, repoName, analysisGUID, branch, privateAccessToken)
	})
}

func (c AnalyzeRepo) analyzeRepo(ctx context.Context, repoName, analysisGUID, branch, privateAccessToken string) error {
	parts := strings.Split(repoName, "/")
	repo := &github.Repo{
		Owner: parts[0],
		Name:  parts[1],
	}
	if len(parts) != 2 {
		return fmt.Errorf("invalid repo name %s", repoName)
	}

	repoCtx := &processors.RepoContext{
		Ctx:                ctx,
		AnalysisGUID:       analysisGUID,
		Branch:             branch,
		Repo:               repo,
		PrivateAccessToken: privateAccessToken,
	}
	p, cleanup, err := c.rpf.BuildProcessor(repoCtx)
	if err != nil {
		return errors.Wrap(err, "failed to build repo processor")
	}
	defer cleanup()

	p.Process(repoCtx)
	return nil
}