# Scenario

**Feature**: Options.Home set uses store list and does not call LiveCapture

```
Home=temp + seed store-sess-1 + LiveCapture spy
  -> MainWith lists store-sess-1
  -> LiveCaptureCalled stays false
```

## Steps

1. Set isolated Home under t.TempDir.
2. Seed one store session `store-sess-1`.
3. Wire LiveCaptureCalled spy; leave LiveResult nil.
4. Args empty (default list).

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.Home = filepath.Join(t.TempDir(), ".agent-run")
	req.Args = []string{}
	req.Sessions = []SessionSeed{
		{
			SessionID: "store-sess-1",
			Runner:    "grok",
			Workspace: "/ws/store",
			UpdatedAt: "2026-08-05T10:00:00Z",
			Live:      true,
			Sendable:  true,
			State:     "idle",
			TTY:       "/dev/ttys099",
		},
	}
	called := false
	req.LiveCaptureCalled = &called
	// Tempt live inject would be wrong path — result present but must not be used.
	req.LiveResult = busyGrokAgentResult()
	return nil
}
```
