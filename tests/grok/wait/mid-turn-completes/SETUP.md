# Scenario

```go
import (
	"context"

	"github.com/xhd2015/dot-pkgs/go-pkgs/logs"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckWaitSessionID
	writeKckWaitSession(t, req, []string{
		`{"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"hi"}}`,
		`{"sessionUpdate":"tool_call","toolCallId":"t1"}`,
	})
	markRunning(t, req)
	req.WatchLine = func(ctx context.Context, path string, opts logs.WatchLineOptions, fn func(line string) error) error {
		return fn(`{"sessionUpdate":"turn_completed","prompt_id":"p2","stop_reason":"end_turn"}`)
	}
	req.Args = []string{"grok", "wait", req.SessionID}
	return nil
}
```
