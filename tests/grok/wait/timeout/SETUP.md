# Scenario

```go
import (
	"context"
	"time"

	"github.com/xhd2015/dot-pkgs/go-pkgs/logs"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckWaitSessionID
	writeKckWaitSession(t, req, []string{
		`{"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"hi"}}`,
	})
	markRunning(t, req)
	req.WaitTimeout = 80 * time.Millisecond
	req.WatchLine = func(ctx context.Context, path string, opts logs.WatchLineOptions, fn func(line string) error) error {
		<-ctx.Done()
		return nil
	}
	req.Args = []string{"grok", "wait", req.SessionID, "--timeout", "80ms"}
	return nil
}
```
