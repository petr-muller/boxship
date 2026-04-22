# Handler Results

## Status
Implemented

## Motivation

During live testing on ota-stage, the ready-for-humans plugin silently did nothing when receiving review events. Sub-plugin handlers returned void, so neither the dispatcher nor an operator could determine whether a plugin considered an event relevant or why it was skipped. The dispatcher itself had no logging for event receipt, routing decisions, or handler outcomes.

This is a gap in the sub-plugin contract: handlers should report their outcome back to the dispatcher.

## Design

### HandlerResult Type

A new `HandlerResult` struct in the `dispatch` package:

```go
type HandlerResult struct {
    Relevant bool
    Reason   string
}
```

Two convenience constructors:
- `Irrelevant(reason)` — the event did not match the plugin's criteria
- `Handled(reason)` — the plugin acted on the event

### Updated SubPlugin Interface

All three handler methods now return `HandlerResult`:

```go
type SubPlugin interface {
    Name() string
    HandlePullRequestEvent(context.Context, *logrus.Entry, github.PullRequestEvent) HandlerResult
    HandleIssueCommentEvent(context.Context, *logrus.Entry, github.IssueCommentEvent) HandlerResult
    HandleReviewEvent(context.Context, *logrus.Entry, github.ReviewEvent) HandlerResult
}
```

It is up to each plugin to decide what "relevant" means. A slash-command plugin receives all comment events but only considers ones containing the command as relevant.

### Dispatcher Logging

The dispatcher now logs at four points:

1. **Event received** (Info): event type, org, repo, PR number, action
2. **Plugin not enabled** (Debug): when `IsEnabled` returns false
3. **Dispatching to plugin** (Debug): before launching goroutine
4. **Plugin completed** (Info if relevant, Debug if irrelevant): `relevant` and `reason` fields

### Sub-Plugin Updates

Each sub-plugin handler returns an appropriate result:
- Early exits return `Irrelevant("reason describing why")`
- Success paths return `Handled("description of action taken")`
- Handlers for unsupported event types return `Irrelevant("only handles X events")`

## Verification

1. `make verify` passes
2. Dev server test: send a review event, confirm dispatcher logs show event receipt, plugin dispatch, and relevance result
