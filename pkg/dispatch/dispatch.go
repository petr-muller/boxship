package dispatch

import (
	"github.com/sirupsen/logrus"
	"sigs.k8s.io/prow/pkg/github"
)

// SubPlugin defines the interface that all boxship sub-plugins must implement.
type SubPlugin interface {
	Name() string
	HandlePullRequestEvent(*logrus.Entry, github.PullRequestEvent)
	HandleIssueCommentEvent(*logrus.Entry, github.IssueCommentEvent)
	HandleReviewEvent(*logrus.Entry, github.ReviewEvent)
}

type Dispatcher struct {
	plugins []SubPlugin
	logger  *logrus.Entry
}

func NewDispatcher(logger *logrus.Entry) *Dispatcher {
	return &Dispatcher{
		logger: logger,
	}
}

func (d *Dispatcher) Register(p SubPlugin) {
	d.logger.WithField("plugin", p.Name()).Info("Registering sub-plugin")
	d.plugins = append(d.plugins, p)
}

func (d *Dispatcher) HandlePullRequestEvent(l *logrus.Entry, event github.PullRequestEvent) {
	for _, p := range d.plugins {
		plugin := p
		go plugin.HandlePullRequestEvent(l.WithField("plugin", plugin.Name()), event)
	}
}

func (d *Dispatcher) HandleIssueCommentEvent(l *logrus.Entry, event github.IssueCommentEvent) {
	for _, p := range d.plugins {
		plugin := p
		go plugin.HandleIssueCommentEvent(l.WithField("plugin", plugin.Name()), event)
	}
}

func (d *Dispatcher) HandleReviewEvent(l *logrus.Entry, event github.ReviewEvent) {
	for _, p := range d.plugins {
		plugin := p
		go plugin.HandleReviewEvent(l.WithField("plugin", plugin.Name()), event)
	}
}
