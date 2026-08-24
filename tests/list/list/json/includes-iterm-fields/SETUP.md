# Scenario

**Feature**: JSON rows include iterm field(s)

```
--json + matching iTerm -> session JSON has iterm string (or window/tab fields)
```

## Steps

1. One session with TTY + single iTerm match.
2. `--json`.

```go
import (
	"testing"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Sessions = []SessionSeed{
		{
			SessionID: "json-iterm",
			Runner:    "grok",
			Workspace: "/ws/j",
			UpdatedAt: "2026-08-05T15:00:00Z",
			Live:      true,
			Sendable:  true,
			State:     "idle",
			TTY:       "/dev/ttys042",
		},
	}
	req.ITermSessions = []run.ITermSession{
		{WindowID: "42", TabIndex: 3, TTY: "/dev/ttys042"},
	}
	req.Args = []string{"--json"}
	return nil
}
```
