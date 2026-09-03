# Scenario

```go
import (
	"context"

	"github.com/xhd2015/dot-pkgs/go-pkgs/logs"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckCodexWaitSessionID
	writeKckCodexWaitSession(t, req, []string{
		eventMsgLine("task_started"),
		`{"timestamp":"2026-08-01T12:00:01.000Z","type":"event_msg","payload":{"type":"item_completed"}}`,
	})
	markRunning(t, req)
	req.WatchLine = func(ctx context.Context, path string, opts logs.WatchLineOptions, fn func(line string) error) error {
		return fn(eventMsgLine("task_complete"))
	}
	req.Args = []string{"codex", "wait", req.SessionID}
	return nil
}
```
