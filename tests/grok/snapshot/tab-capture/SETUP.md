# Scenario

```go
import (
	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = "019f283b-dddd-7ddd-dddd-dddddddddddd"
	writeKckSnapshotSession(t, req)
	req.CurrentSessionID = "w0t1p0:CURRENT-TAB-UUID"
	req.ControllingTTY = "/dev/ttys101"
	req.ITerm = []iterm2.SessionRef{
		{WindowID: "100", WindowName: "work", TabIndex: 1, SessionID: "w0t1p0:CURRENT-TAB-UUID", TTY: "/dev/ttys101", Name: "current"},
		{WindowID: "100", WindowName: "work", TabIndex: 2, SessionID: "w0t2p0:TAB2-UUID", TTY: "/dev/ttys102", Name: "grok-tab"},
	}
	req.Procs = []sessions.FocusProc{
		{PID: 8200, PPID: 1, TTY: "/dev/ttys102", Cmd: "/usr/local/bin/grok"},
	}
	req.OpenFiles[8200] = []string{
		"/Users/fixture/.grok/sessions/%2Ftmp%2Fproj/" + req.SessionID + "/events.jsonl",
	}
	req.ContentsByID = map[string]iterm2.ContentsResult{
		"TAB2-UUID": {SessionID: "TAB2-UUID", App: "/Applications/iTerm.app", Contents: "kck tab pane"},
	}
	req.Args = []string{"grok", "snapshot", "--tab", "2"}
	return nil
}
```
