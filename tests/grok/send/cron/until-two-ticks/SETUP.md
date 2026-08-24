# Scenario

**Feature**: until schedule sends twice then exits cleanly

```
every-5m-until-19h00m @ 18:50 → send 18:50, 18:55 → cron done
```

## Steps

1. CronClock = 2026-08-24T18:50:00Z.
2. Expr every-5m-until-19h00m.

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
	req.CronClock = time.Date(2026, 8, 24, 18, 50, 0, 0, time.UTC)
	req.Args = []string{
		"grok", "send", "ping", "--session-id", req.SessionID,
		"--cron", "every-5m-until-19h00m",
	}
	return nil
}
```
