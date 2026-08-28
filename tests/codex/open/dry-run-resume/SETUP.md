# Scenario

```go
import (
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckOpenSessionID
	writeKckOpenSession(t, req)
	req.NoAgentRun = true
	req.Args = []string{"codex", "open", req.SessionID, "--dry-run", "--no-agent-run"}
	return nil
}
```
