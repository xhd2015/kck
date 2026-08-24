# Scenario

**Feature**: unreachable terminal does not crash list

```
Probe{Live:false, State:unknown} -> LIVE no; exit 0; optional warning on stderr
```

## Steps

1. One session with not-live probe.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Sessions = []SessionSeed{
		{
			SessionID: "dead-tty",
			Runner:    "grok",
			Workspace: "/ws/dead",
			UpdatedAt: "2026-08-04T10:00:00Z",
			Live:      false,
			Sendable:  false,
			State:     "unknown",
		},
	}
	req.Args = []string{}
	return nil
}
```
