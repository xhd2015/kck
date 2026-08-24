# Scenario

**Feature**: store list marks AGENT_RUN yes and AGENT_SID from runner_session_id

```
meta.session_id=store-1 runner_session_id=grok-native-99
  -> AGENT_RUN yes, AGENT_SID grok-native-99
```

## Steps

1. Seed one store session with RunnerSessionID.
2. Probe live idle so row appears cleanly.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.Args = nil
	req.Sessions = []SessionSeed{
		{
			SessionID:       "store-sess-ar",
			Runner:          "grok",
			RunnerSessionID: "grok-native-99",
			Workspace:       "/ws/store",
			UpdatedAt:       "2026-08-06T10:00:00Z",
			Live:            true,
			Sendable:        true,
			State:           "idle",
			TTY:             "/dev/ttys099",
		},
	}
	req.SkipITerm = true
	return nil
}
```
