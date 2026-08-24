# Scenario

**Feature**: later-tick send failure warns and continues

```
3 legal fires; 2nd SendText fails → warning:; 3rd still sends
```

## Steps

1. every-5m-until-19h05m @ 18:50 → fires 18:50, 18:55, 19:00.
2. FailSendOnTick = 2.

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
	req.FailSendOnTick = 2
	req.Args = []string{
		"grok", "send", "ping", "--session-id", req.SessionID,
		"--cron", "every-5m-until-19h05m",
	}
	return nil
}
```
