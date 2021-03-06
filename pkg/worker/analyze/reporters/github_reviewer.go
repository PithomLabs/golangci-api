package reporters

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/golangci/golangci-api/pkg/worker/lib/experiments"

	"github.com/golangci/golangci-api/pkg/goenvbuild/config"

	envbuildresult "github.com/golangci/golangci-api/pkg/goenvbuild/result"
	"github.com/golangci/golangci-api/pkg/worker/analyze/linters/result"
	"github.com/golangci/golangci-api/pkg/worker/lib/github"
	gh "github.com/google/go-github/github"
	"github.com/pkg/errors"
)

type GithubReviewer struct {
	*github.Context
	client github.Client
	ec     *experiments.Checker
}

func NewGithubReviewer(c *github.Context, client github.Client, ec *experiments.Checker) *GithubReviewer {
	accessToken := os.Getenv("GITHUB_REVIEWER_ACCESS_TOKEN")
	if accessToken != "" { // review as special user
		cCopy := *c
		cCopy.GithubAccessToken = accessToken
		c = &cCopy
	}
	ret := &GithubReviewer{
		Context: c,
		client:  client,
		ec:      ec,
	}
	return ret
}

type existingComment struct {
	file string
	line int
}

type existingComments []existingComment

func (ecs existingComments) contains(i *result.Issue) bool {
	for _, c := range ecs {
		if c.file == i.File && c.line == i.HunkPos {
			return true
		}
	}

	return false
}

func (gr GithubReviewer) fetchExistingComments(ctx context.Context) (existingComments, error) {
	comments, err := gr.client.GetPullRequestComments(ctx, gr.Context)
	if err != nil {
		return nil, err
	}

	var ret existingComments
	for _, c := range comments {
		if c.Position == nil { // comment on outdated code, skip it
			continue
		}
		ret = append(ret, existingComment{
			file: c.GetPath(),
			line: c.GetPosition(),
		})
	}

	return ret, nil
}

func (gr GithubReviewer) makeSimpleIssueCommentBody(i *result.Issue) string {
	text := i.Text
	if i.FromLinter != "" {
		text += fmt.Sprintf(" (from `%s`)", i.FromLinter)
	}
	return text
}

func (gr GithubReviewer) makeIssueCommentBody(i *result.Issue, buildConfig *config.Service) string {
	if buildConfig.SuggestedChanges.Disabled {
		return gr.makeSimpleIssueCommentBody(i)
	}

	if i.Replacement == nil {
		return gr.makeSimpleIssueCommentBody(i)
	}

	if i.LineRange != nil && i.LineRange.From != i.LineRange.To {
		// github api doesn't support multi-line suggestion
		return gr.makeSimpleIssueCommentBody(i)
	}

	if !gr.ec.IsActiveForRepo("SUGGESTED_CHANGES", gr.Repo.Owner, gr.Repo.Name) {
		return gr.makeSimpleIssueCommentBody(i)
	}

	var suggestionBody string
	if !i.Replacement.NeedOnlyDelete {
		suggestionBody = strings.Join(i.Replacement.NewLines, "\n")
	}
	return fmt.Sprintf("%s\n```suggestion\n%s\n```", gr.makeSimpleIssueCommentBody(i), suggestionBody)
}

func (gr GithubReviewer) makeComments(issues []result.Issue, ec existingComments,
	buildConfig *config.Service) []*gh.DraftReviewComment {

	comments := []*gh.DraftReviewComment{}
	for _, i := range issues {
		if ec.contains(&i) {
			continue // don't be annoying: don't comment on the same line twice
		}

		comment := &gh.DraftReviewComment{
			Path:     gh.String(i.File),
			Position: gh.Int(i.HunkPos),
			Body:     gh.String(gr.makeIssueCommentBody(&i, buildConfig)),
		}
		comments = append(comments, comment)
	}

	return comments
}

func (gr GithubReviewer) Report(ctx context.Context, buildConfig *config.Service, buildLog *envbuildresult.Log,
	ref string, issues []result.Issue) error {

	return buildLog.RunNewGroup("post review", func(sg *envbuildresult.StepGroup) error {
		step := sg.AddStep("check issues")
		if len(issues) == 0 {
			step.AddOutputLine("Nothing to report: no issues found")
			return nil
		}
		step.AddOutputLine("Have %d issues", len(issues))

		sg.AddStep("fetch existing comments")
		existingComments, err := gr.fetchExistingComments(ctx)
		if err != nil {
			return err
		}

		step = sg.AddStep("build new review comments")
		comments := gr.makeComments(issues, existingComments, buildConfig)
		if len(comments) == 0 {
			step.AddOutputLine("No new comments were built")
			return nil // all comments are already exist
		}
		step.AddOutputLine("Send %d comments about new issues", len(comments))

		sg.AddStep("create GitHub review")
		review := &gh.PullRequestReviewRequest{
			CommitID: gh.String(ref),
			Event:    gh.String("COMMENT"),
			Body:     gh.String(""),
			Comments: comments,
		}
		if err := gr.client.CreateReview(ctx, gr.Context, review); err != nil {
			return errors.Wrapf(err, "can't create review %v", review)
		}

		return nil
	})
}
