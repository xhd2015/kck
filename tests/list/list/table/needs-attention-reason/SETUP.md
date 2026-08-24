# Scenario

**Feature**: D1 needs_attention row shows REASON

```
live + !sendable + !exited + Reason -> row shows REASON text; footer counts attention
```

## Steps

1. Single session with attention probe.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Sessions = []SessionSeed{
		{
			SessionID: "attn-1",
			Runner:    "grok",
			Workspace: "/ws/attn",
			UpdatedAt: "2026-08-05T10:00:00Z",
			Live:      true,
			Sendable:  false,
			State:     "running",
			Reason:    "awaiting confirmation",
			Exited:    false,
		},
	}
	req.Args = []string{}
	return nil
}
```
