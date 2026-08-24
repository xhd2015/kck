# Scenario

**Feature**: not-within skips quiet band then resumes

```
every-5m-until-18h45m-not-within-18h30m-to-18h40m @ 18:28
→ send 18:28, skip quiet, send 18:43, done
(suffix order: until before not-within)
```

## Steps

1. CronClock = 18:28 UTC.
2. Compose until + not-within (fixed easycron order).

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
	req.CronClock = time.Date(2026, 8, 24, 18, 28, 0, 0, time.UTC)
	req.Args = []string{
		"grok", "send", "ping", "--session-id", req.SessionID,
		"--cron", "every-5m-until-18h45m-not-within-18h30m-to-18h40m",
	}
	return nil
}
```
