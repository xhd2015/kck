# Scenario

**Feature**: --cron --dry-run prints schedule preview and one would-send

```
no loop; no SendText
```

## Steps

1. CronClock set; --dry-run + --cron every-1h.

```go
import (
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	addLiveGrok(req, 5001, "ttys148")
	req.ITerm = oneITermTab()
	req.CronClock = time.Date(2026, 8, 24, 14, 0, 0, 0, time.UTC)
	req.Args = []string{
		"grok", "send", "ping", "--session-id", req.SessionID,
		"--cron", "every-1h", "--dry-run",
	}
	return nil
}
```
