# Scenario

**Feature**: ITERM column from injectable ListITerm + probe TTY

```
Probe.TTY + ListITerm refs -> w=<id> t=<tab> | (+N) | -
```

## Steps

1. Leaves seed one session with TTY and iTerm fixtures.
2. Default resolve on (no `--no-iterm` unless leaf).

```go
import (
	"testing"

	"kck/run"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	// Common single-session shell; leaves override TTY / ITerm.
	if len(req.Sessions) == 0 {
		req.Sessions = []SessionSeed{
			{
				SessionID: "iterm-s",
				Runner:    "grok",
				Workspace: "/ws/iterm",
				UpdatedAt: "2026-08-05T12:00:00Z",
				Live:      true,
				Sendable:  true,
				State:     "idle",
				TTY:       "/dev/ttys042",
			},
		}
	}
	_ = run.ITermSession{} // keep import for leaf helpers if needed
	return nil
}
```
