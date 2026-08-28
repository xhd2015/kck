# Scenario

```go
import (
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSendSessionID
	writeKckSendSession(t, req)
	addLiveCodexSend(req, 5001, "ttys148")
	req.ITerm = oneITermTabSend()
	req.NoAgentRun = true
	req.Args = []string{"codex", "send", "hello", "--session-id", req.SessionID, "--no-agent-run"}
	return nil
}
```
