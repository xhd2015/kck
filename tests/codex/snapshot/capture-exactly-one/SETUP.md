# Scenario

```go
import (
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckSnapshotSessionID
	writeKckSnapshotSession(t, req)
	addLiveCodex(req, 5001, "ttys148")
	req.ITerm = oneITermTab()
	req.ContentsByID = map[string]iterm2.ContentsResult{
		"w2t1p0": {SessionID: "w2t1p0", App: "/Applications/iTerm.app", Contents: "kck pane text"},
	}
	req.Args = []string{"codex", "snapshot", req.SessionID}
	return nil
}
```
