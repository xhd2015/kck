# Scenario

```go
import (
	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckListSessionID
	writeKckListSession(t, req)
	addLiveCodexList(req, 5001, "ttys148")
	req.ITerm = oneITermTabList()
	req.Args = []string{"codex", "list"}
	return nil
}
```
