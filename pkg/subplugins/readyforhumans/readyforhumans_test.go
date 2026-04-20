package readyforhumans

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/github"
)

type fakeClient struct {
	labelsAdded []string
}

func (c *fakeClient) AddLabel(org, repo string, number int, label string) error {
	c.labelsAdded = append(c.labelsAdded, fmt.Sprintf("%s/%s#%d:%s", org, repo, number, label))
	return nil
}

func TestHandleReviewEvent(t *testing.T) {
	testCases := []struct {
		name        string
		event       github.ReviewEvent
		expectAdded []string
	}{
		{
			name: "ignore review from different user",
			event: github.ReviewEvent{
				Action: github.ReviewActionSubmitted,
				Review: github.Review{
					User:  github.User{Login: "someone-else"},
					State: github.ReviewStateApproved,
				},
				PullRequest: github.PullRequest{Number: 1},
				Repo:        github.Repo{Owner: github.User{Login: "org"}, Name: "repo"},
			},
		},
		{
			name: "ignore edited action from coderabbit",
			event: github.ReviewEvent{
				Action: github.ReviewActionEdited,
				Review: github.Review{
					User:  github.User{Login: coderabbitBotLogin},
					State: github.ReviewStateApproved,
				},
				PullRequest: github.PullRequest{Number: 1},
				Repo:        github.Repo{Owner: github.User{Login: "org"}, Name: "repo"},
			},
		},
		{
			name: "ignore dismissed action from coderabbit",
			event: github.ReviewEvent{
				Action: github.ReviewActionDismissed,
				Review: github.Review{
					User: github.User{Login: coderabbitBotLogin},
				},
				PullRequest: github.PullRequest{Number: 1},
				Repo:        github.Repo{Owner: github.User{Login: "org"}, Name: "repo"},
			},
		},
		{
			name: "ignore non-approved state",
			event: github.ReviewEvent{
				Action: github.ReviewActionSubmitted,
				Review: github.Review{
					User:  github.User{Login: coderabbitBotLogin},
					State: github.ReviewStateChangesRequested,
				},
				PullRequest: github.PullRequest{Number: 42},
				Repo:        github.Repo{Owner: github.User{Login: "org"}, Name: "repo"},
			},
		},
		{
			name: "ignore comment review state",
			event: github.ReviewEvent{
				Action: github.ReviewActionSubmitted,
				Review: github.Review{
					User:  github.User{Login: coderabbitBotLogin},
					State: github.ReviewStateCommented,
				},
				PullRequest: github.PullRequest{Number: 42},
				Repo:        github.Repo{Owner: github.User{Login: "org"}, Name: "repo"},
			},
		},
		{
			name: "add label on approval when absent",
			event: github.ReviewEvent{
				Action: github.ReviewActionSubmitted,
				Review: github.Review{
					User:  github.User{Login: coderabbitBotLogin},
					State: github.ReviewStateApproved,
				},
				PullRequest: github.PullRequest{Number: 42},
				Repo:        github.Repo{Owner: github.User{Login: "org"}, Name: "repo"},
			},
			expectAdded: []string{"org/repo#42:ready-for-humans"},
		},
		{
			name: "no-op on approval when label already present",
			event: github.ReviewEvent{
				Action: github.ReviewActionSubmitted,
				Review: github.Review{
					User:  github.User{Login: coderabbitBotLogin},
					State: github.ReviewStateApproved,
				},
				PullRequest: github.PullRequest{
					Number: 42,
					Labels: []github.Label{{Name: labelName}},
				},
				Repo: github.Repo{Owner: github.User{Login: "org"}, Name: "repo"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fc := &fakeClient{}
			p := New(fc)
			p.HandleReviewEvent(logrus.NewEntry(logrus.StandardLogger()), tc.event)

			if len(tc.expectAdded) == 0 && len(fc.labelsAdded) == 0 {
				return
			}
			if fmt.Sprintf("%v", tc.expectAdded) != fmt.Sprintf("%v", fc.labelsAdded) {
				t.Errorf("labels added: want %v, got %v", tc.expectAdded, fc.labelsAdded)
			}
		})
	}
}
