# Scenario

```go
import (
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckOpenSessionID
	writeKckOpenSession(t, req)
	addLiveCodexOpen(req, 5001, "ttys148")
	req.ITerm = oneITermTabOpen()
	req.NoAgentRun = true
	req.Args = []string{"codex", "open", req.SessionID, "--no-agent-run"}
	return nil
}
```
