# Scenario

**Feature**: empty live Result prints zero-session summary

```
LiveCapture -> Result{Snapshot: empty windows} -> footer 0 sessions
```

## Steps

1. Inject Result with zero windows (or empty Snapshot).
2. Default list args.

```go
import (
	"testing"

	"github.com/xhd2015/agent-pro/pkgs/itermsnapshot"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2/snapshot"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	req.LiveResult = &itermsnapshot.Result{
		Snapshot: &snapshot.Snapshot{
			CapturedAt: "2026-08-06T12:00:00Z",
			Host:       "testhost",
			Source:     "iterm2",
			Summary:    snapshot.SnapshotSummary{},
			Windows:    nil,
		},
	}
	return nil
}
```
