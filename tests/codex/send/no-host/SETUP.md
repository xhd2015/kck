# Scenario

```go
import (
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	req.NoAgentRun = true
	req.Args = []string{"codex", "send", "hi", "--session-id", req.SessionID, "--no-agent-run"}
	return nil
}
```
