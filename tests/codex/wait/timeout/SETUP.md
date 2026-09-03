# Scenario

```go
import (
	"context"
	"time"

	"github.com/xhd2015/dot-pkgs/go-pkgs/logs"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckCodexWaitSessionID
	writeKckCodexWaitSession(t, req, []string{
		eventMsgLine("task_started"),
	})
	markRunning(t, req)
	req.WaitTimeout = 80 * time.Millisecond
	req.WatchLine = func(ctx context.Context, path string, opts logs.WatchLineOptions, fn func(line string) error) error {
		<-ctx.Done()
		return nil
	}
	req.Args = []string{"codex", "wait", req.SessionID, "--timeout", "80ms"}
	return nil
}
```
