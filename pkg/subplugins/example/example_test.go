package example

import (
	"testing"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/github/fakegithub"

	"github.com/petr-muller/boxship/pkg/testhelpers"
)

func TestHandlePullRequestEvent(t *testing.T) {
	ghc := fakegithub.NewFakeClient()
	plugin := New(ghc)

	event := testhelpers.NewPullRequestEvent("org", "repo", 42, "opened")
	plugin.HandlePullRequestEvent(logrus.NewEntry(logrus.StandardLogger()), event)

	if len(ghc.IssueCommentsAdded) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(ghc.IssueCommentsAdded))
	}
	expected := "org/repo#42:example plugin noticed PR #42"
	if ghc.IssueCommentsAdded[0] != expected {
		t.Errorf("expected comment %q, got %q", expected, ghc.IssueCommentsAdded[0])
	}
}

func TestHandleIssueCommentEvent(t *testing.T) {
	ghc := fakegithub.NewFakeClient()
	plugin := New(ghc)

	event := testhelpers.NewIssueCommentEvent("org", "repo", 1, "/example")
	plugin.HandleIssueCommentEvent(logrus.NewEntry(logrus.StandardLogger()), event)

	if len(ghc.IssueCommentsAdded) != 0 {
		t.Errorf("expected no comments, got %d", len(ghc.IssueCommentsAdded))
	}
}
