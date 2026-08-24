# Scenario

**Feature**: multi iTerm match shows `(+N)` extra count

```
two refs same TTY -> w=<first> t=<first>(+1)
```

## Steps

1. Session TTY matches two iTerm sessions.
2. First listed wins as primary; N = len-1.

```go
import (
	"testing"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Args = []string{}
	req.Sessions = []SessionSeed{
		{
			SessionID: "iterm-s",
			Runner:    "grok",
			Workspace: "/ws/iterm",
			UpdatedAt: "2026-08-05T12:00:00Z",
			Live:      true,
			Sendable:  true,
			State:     "idle",
			TTY:       "/dev/ttys010",
		},
	}
	req.ITermSessions = []run.ITermSession{
		{WindowID: "1", TabIndex: 2, TTY: "/dev/ttys010"},
		{WindowID: "7", TabIndex: 4, TTY: "/dev/ttys010"},
	}
	return nil
}
```
