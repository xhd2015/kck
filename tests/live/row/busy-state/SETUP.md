# Scenario

**Feature**: busy pane maps to LIVE yes, SENDABLE no, STATE busy, needs attention

```
Idle=false + Agent -> LIVE yes SENDABLE no STATE busy; footer 1 needs attention
```

## Steps

1. Inject busyGrokAgentResult().

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveResult = busyGrokAgentResult()
	return nil
}
```
