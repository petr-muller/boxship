package dispatch

import (
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/github"

	"github.com/petr-muller/boxship/pkg/testhelpers"
)

type fakePlugin struct{}

func (f *fakePlugin) Name() string                                                       { return "fake" }
func (f *fakePlugin) HandlePullRequestEvent(_ *logrus.Entry, _ github.PullRequestEvent)   {}
func (f *fakePlugin) HandleIssueCommentEvent(_ *logrus.Entry, _ github.IssueCommentEvent) {}
func (f *fakePlugin) HandleReviewEvent(_ *logrus.Entry, _ github.ReviewEvent)             {}

type recordingPlugin struct {
	name               string
	mu                 sync.Mutex
	prEvents           []github.PullRequestEvent
	issueCommentEvents []github.IssueCommentEvent
	reviewEvents       []github.ReviewEvent
}

func (r *recordingPlugin) Name() string { return r.name }

func (r *recordingPlugin) HandlePullRequestEvent(_ *logrus.Entry, event github.PullRequestEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.prEvents = append(r.prEvents, event)
}

func (r *recordingPlugin) HandleIssueCommentEvent(_ *logrus.Entry, event github.IssueCommentEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.issueCommentEvents = append(r.issueCommentEvents, event)
}

func (r *recordingPlugin) HandleReviewEvent(_ *logrus.Entry, event github.ReviewEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.reviewEvents = append(r.reviewEvents, event)
}

func TestDispatcherRegister(t *testing.T) {
	d := NewDispatcher(logrus.NewEntry(logrus.StandardLogger()))
	d.Register(&fakePlugin{})

	if len(d.plugins) != 1 {
		t.Errorf("expected 1 plugin, got %d", len(d.plugins))
	}
	if d.plugins[0].Name() != "fake" {
		t.Errorf("expected plugin name 'fake', got %q", d.plugins[0].Name())
	}
}

func TestDispatcherHandlePullRequestEvent(t *testing.T) {
	d := NewDispatcher(logrus.NewEntry(logrus.StandardLogger()))
	p1 := &recordingPlugin{name: "plugin-1"}
	p2 := &recordingPlugin{name: "plugin-2"}
	d.Register(p1)
	d.Register(p2)

	event := testhelpers.NewPullRequestEvent("org", "repo", 1, "opened")
	d.HandlePullRequestEvent(logrus.NewEntry(logrus.StandardLogger()), event)

	deadline := time.After(2 * time.Second)
	for {
		p1.mu.Lock()
		p2.mu.Lock()
		p1Count := len(p1.prEvents)
		p2Count := len(p2.prEvents)
		p1.mu.Unlock()
		p2.mu.Unlock()
		if p1Count == 1 && p2Count == 1 {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for events: plugin-1 got %d, plugin-2 got %d", p1Count, p2Count)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	if p1.prEvents[0].Number != 1 {
		t.Errorf("plugin-1: expected PR #1, got #%d", p1.prEvents[0].Number)
	}
	if p2.prEvents[0].Number != 1 {
		t.Errorf("plugin-2: expected PR #1, got #%d", p2.prEvents[0].Number)
	}
}

func TestDispatcherHandleIssueCommentEvent(t *testing.T) {
	d := NewDispatcher(logrus.NewEntry(logrus.StandardLogger()))
	p := &recordingPlugin{name: "recorder"}
	d.Register(p)

	event := testhelpers.NewIssueCommentEvent("org", "repo", 5, "/test")
	d.HandleIssueCommentEvent(logrus.NewEntry(logrus.StandardLogger()), event)

	deadline := time.After(2 * time.Second)
	for {
		p.mu.Lock()
		count := len(p.issueCommentEvents)
		p.mu.Unlock()
		if count == 1 {
			break
		}
		select {
		case <-deadline:
			t.Fatalf("timed out waiting for event")
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}

	if p.issueCommentEvents[0].Comment.Body != "/test" {
		t.Errorf("expected comment body '/test', got %q", p.issueCommentEvents[0].Comment.Body)
	}
}
