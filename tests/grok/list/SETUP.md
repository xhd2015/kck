# Scenario

**Feature**: `kck grok list` thin wrapper over list-live

```
run.MainWith(["grok","list",…], GrokListLiveOpts) -> table/JSON
```

## Preconditions

- Injectable `GrokListLiveOpts` / `ListLiveFake`; no live iTerm.

```go
import (
	"strings"
	"testing"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

const fixtureKckListSID = "019f283a-aaaa-7aaa-aaaa-aaaaaaaaaa01"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = t
	_ = d
	if req.OpenFiles == nil {
		req.OpenFiles = map[int][]string{}
	}
	if req.PaneByTTY == nil {
		req.PaneByTTY = map[string]sessions.LivePaneInfo{}
	}
	return nil
}

func grokKckListPath(sessionID string) string {
	return "/Users/fixture/.grok/sessions/%2Ftmp%2Fproj/" + sessionID + "/events.jsonl"
}

func addKckListHost(req *Request, pid int, ttyBare, sessionID, windowID string, tabIndex int) {
	req.Procs = append(req.Procs, sessions.FocusProc{
		PID: pid, PPID: 1, TTY: ttyBare, Cmd: "/usr/local/bin/grok",
	})
	req.OpenFiles[pid] = []string{grokKckListPath(sessionID)}
	req.ITerm = append(req.ITerm, iterm2.SessionRef{
		WindowID: windowID, TabIndex: tabIndex,
		SessionID: "iterm-" + ttyBare, TTY: "/dev/" + ttyBare,
	})
	idle := true
	req.PaneByTTY["/dev/"+ttyBare] = sessions.LivePaneInfo{Idle: &idle, Cwd: "/tmp/proj"}
}

func assertNoHarnessErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected harness error: %v", err)
	}
}

func assertSuccess(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode != 0 || resp.ErrText != "" {
		t.Fatalf("want success, exit=%d err=%q stderr=%q stdout=%q",
			resp.ExitCode, resp.ErrText, resp.Stderr, resp.Stdout)
	}
}

func assertContains(t *testing.T, got, want, label string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%s must contain %q; got:\n%s", label, want, got)
	}
}
```
