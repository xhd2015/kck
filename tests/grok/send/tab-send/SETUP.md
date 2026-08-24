# Scenario

```go
import (
	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

const (
	fixtureKckTabSendSessionID = "019f283b-cccc-7ccc-cccc-cccccccccccc"
	pidKckTabSendGrok          = 8400
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.SessionID = fixtureKckTabSendSessionID
	writeKckSendSession(t, req)
	req.CurrentSessionID = "w0t1p0:CURRENT-TAB-UUID"
	req.ControllingTTY = "/dev/ttys101"
	req.ITerm = []iterm2.SessionRef{
		{WindowID: "100", WindowName: "work", TabIndex: 1, SessionID: "w0t1p0:CURRENT-TAB-UUID", TTY: "/dev/ttys101", Name: "current"},
		{WindowID: "100", WindowName: "work", TabIndex: 2, SessionID: "w0t2p0:TAB2-UUID", TTY: "/dev/ttys102", Name: "grok-tab"},
	}
	req.Procs = []sessions.FocusProc{
		{PID: pidKckTabSendGrok, PPID: 1, TTY: "/dev/ttys102", Cmd: "/usr/local/bin/grok"},
	}
	req.OpenFiles[pidKckTabSendGrok] = []string{
		"/Users/fixture/.grok/sessions/%2Ftmp%2Fproj/" + fixtureKckTabSendSessionID + "/events.jsonl",
	}
	req.Args = []string{"grok", "send", "from-tab", "--tab", "2"}
	return nil
}
```
