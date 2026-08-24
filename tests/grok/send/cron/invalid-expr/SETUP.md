# Scenario

**Feature**: bad --cron fails before any send

```
--cron every-1h-at-60m → Error: invalid --cron:…
```

## Steps

1. Live host present (would send if cron parsed).
2. Args include invalid cron.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	addLiveGrok(req, 5001, "ttys148")
	req.ITerm = oneITermTab()
	req.Args = []string{"grok", "send", "ping", "--session-id", req.SessionID, "--cron", "every-1h-at-60m"}
	return nil
}
```
