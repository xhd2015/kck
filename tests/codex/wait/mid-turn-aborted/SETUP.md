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
	})
	markRunning(t, req)
	req.WatchLine = func(ctx context.Context, path string, opts logs.WatchLineOptions, fn func(line string) error) error {
		return fn(eventMsgLine("turn_aborted"))
	}
	req.Args = []string{"codex", "wait", req.SessionID}
	return nil
}
```
