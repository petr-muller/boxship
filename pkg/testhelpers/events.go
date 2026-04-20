package testhelpers

import "sigs.k8s.io/prow/pkg/github"

func NewPullRequestEvent(org, repo string, number int, action github.PullRequestEventAction) github.PullRequestEvent {
	return github.PullRequestEvent{
		Action: action,
		Number: number,
		PullRequest: github.PullRequest{
			Number: number,
			Base: github.PullRequestBranch{
				Repo: github.Repo{
					Owner: github.User{Login: org},
					Name:  repo,
				},
			},
		},
		Repo: github.Repo{
			Owner: github.User{Login: org},
			Name:  repo,
		},
	}
}

func NewReviewEvent(org, repo string, number int, action github.ReviewEventAction, reviewer string, state github.ReviewState) github.ReviewEvent {
	return github.ReviewEvent{
		Action: action,
		Review: github.Review{
			User:  github.User{Login: reviewer},
			State: state,
		},
		PullRequest: github.PullRequest{
			Number: number,
		},
		Repo: github.Repo{
			Owner: github.User{Login: org},
			Name:  repo,
		},
	}
}

func NewIssueCommentEvent(org, repo string, number int, body string) github.IssueCommentEvent {
	return github.IssueCommentEvent{
		Action: github.IssueCommentActionCreated,
		Issue: github.Issue{
			Number: number,
		},
		Comment: github.IssueComment{
			Body: body,
			User: github.User{Login: "test-user"},
		},
		Repo: github.Repo{
			Owner: github.User{Login: org},
			Name:  repo,
		},
	}
}
