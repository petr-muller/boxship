package readyforhumans

import (
	"strings"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/github"
)

const (
	pluginName         = "ready-for-humans"
	labelName          = "ready-for-humans"
	coderabbitBotLogin = "coderabbitai[bot]"
)

type githubClient interface {
	AddLabel(org, repo string, number int, label string) error
}

type Plugin struct {
	ghc githubClient
}

func New(ghc githubClient) *Plugin {
	return &Plugin{ghc: ghc}
}

func (p *Plugin) Name() string { return pluginName }

func (p *Plugin) HandlePullRequestEvent(_ *logrus.Entry, _ github.PullRequestEvent) {}

func (p *Plugin) HandleIssueCommentEvent(_ *logrus.Entry, _ github.IssueCommentEvent) {}

func (p *Plugin) HandleReviewEvent(l *logrus.Entry, re github.ReviewEvent) {
	if re.Review.User.Login != coderabbitBotLogin {
		return
	}

	if re.Action != github.ReviewActionSubmitted {
		return
	}

	if re.Review.State != github.ReviewStateApproved {
		return
	}

	if prHasLabel(re.PullRequest, labelName) {
		return
	}

	org := re.Repo.Owner.Login
	repo := re.Repo.Name
	number := re.PullRequest.Number

	if err := p.ghc.AddLabel(org, repo, number, labelName); err != nil {
		l.WithError(err).Error("Failed to add ready-for-humans label")
	}
}

func prHasLabel(pr github.PullRequest, label string) bool {
	for _, l := range pr.Labels {
		if strings.EqualFold(l.Name, label) {
			return true
		}
	}
	return false
}
