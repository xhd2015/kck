# Scenario

**Feature**: sendable idle session shows SENDABLE yes

```
Probe{Live:true, Sendable:true, State:idle} -> row SENDABLE yes; footer sendable >= 1
```

## Steps

1. Single sendable session.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Sessions = []SessionSeed{
		{
			SessionID: "ready-1",
			Runner:    "grok",
			Workspace: "/ws/ready",
			UpdatedAt: "2026-08-05T11:00:00Z",
			Live:      true,
			Sendable:  true,
			State:     "idle",
		},
	}
	req.Args = []string{}
	return nil
}
```
