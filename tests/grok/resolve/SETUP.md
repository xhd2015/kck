# Scenario

**Feature**: kck grok resolve thin CLI

```
run.MainWith(Options{Args: ["grok","resolve",…], GrokHome, GrokResolveOpts})
  -> sessions.RunResolve
```

## Preconditions

- No live iTerm / ps / real GROK_HOME.
- Session id from open-file paths only (never cmdline `--resume` / `--session-id`).

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

const (
	fixtureSessionID    = "019f283b-aaaa-7aaa-aaaa-aaaaaaaaaaaa"
	fixtureTabSessionID = "019f283b-dddd-7ddd-dddd-dddddddddddd"
	wrongResumeSessionID = "00000000-0000-0000-0000-000000000000"
	wrongFlagSessionID   = "11111111-1111-1111-1111-111111111111"

	pidGrok  = 4242
	pidBash  = 5000
	pidStart = 6000

	pidTabGrok1 = 8100
	pidTabGrok2 = 8200
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.TempDir == "" {
		req.TempDir = t.TempDir()
	}
	if req.GrokHome == "" {
		req.GrokHome = filepath.Join(req.TempDir, ".grok")
	}
	if req.PID == 0 {
		req.PID = pidStart
	}
	if req.OpenFiles == nil {
		req.OpenFiles = map[int][]string{}
	}
	if req.Args == nil {
		req.Args = []string{}
	}
	return nil
}

func grokCmdWithIgnoredFlags() string {
	return "/usr/local/bin/grok --resume " + wrongResumeSessionID + " --session-id " + wrongFlagSessionID
}

func grokSessionPath(sessionID string) string {
	return "/Users/fixture/.grok/sessions/%2Ftmp%2Fproj/" + sessionID + "/events.jsonl"
}

func defaultAncestorChain() []FixtureProc {
	return []FixtureProc{
		{PID: pidGrok, PPID: 1, Cmd: grokCmdWithIgnoredFlags()},
		{PID: pidBash, PPID: pidGrok, Cmd: "/bin/bash"},
		{PID: pidStart, PPID: pidBash, Cmd: "/usr/local/bin/agent-pro"},
	}
}

func seedHit(req *Request, sessionID string, grokPID int) {
	if len(req.Procs) == 0 {
		req.Procs = defaultAncestorChain()
	}
	req.OpenFiles[grokPID] = []string{grokSessionPath(sessionID)}
}

// seedTabWindow installs a 3-tab window; current is tab 1 (/dev/ttys101).
// Tab 2 hosts fixtureTabSessionID on /dev/ttys102; tab 3 is bash-only.
func seedTabWindow(req *Request) {
	req.CurrentSessionID = "w0t1p0:CURRENT-TAB-UUID"
	req.ControllingTTY = "/dev/ttys101"
	req.ITerm = []iterm2.SessionRef{
		{WindowID: "100", WindowName: "work", TabIndex: 1, SessionID: "w0t1p0:CURRENT-TAB-UUID", TTY: "/dev/ttys101", Name: "current"},
		{WindowID: "100", WindowName: "work", TabIndex: 2, SessionID: "w0t2p0:TAB2-UUID", TTY: "/dev/ttys102", Name: "grok-tab"},
		{WindowID: "100", WindowName: "work", TabIndex: 3, SessionID: "w0t3p0:TAB3-UUID", TTY: "/dev/ttys103", Name: "bash-only"},
	}
	req.FocusProcs = []sessions.FocusProc{
		{PID: pidTabGrok1, PPID: 1, TTY: "/dev/ttys101", Cmd: "/usr/local/bin/bash"},
		{PID: pidTabGrok2, PPID: 1, TTY: "/dev/ttys102", Cmd: "/usr/local/bin/grok"},
		{PID: 9100, PPID: 1, TTY: "/dev/ttys103", Cmd: "/bin/bash"},
	}
	req.OpenFiles[pidTabGrok2] = []string{grokSessionPath(fixtureTabSessionID)}
}

func assertNoHarnessErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
}

func assertSuccess(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode != 0 || resp.ErrText != "" {
		t.Fatalf("want success, exit=%d err=%q stderr=%q stdout=%q",
			resp.ExitCode, resp.ErrText, resp.Stderr, resp.Stdout)
	}
}

func assertFailure(t *testing.T, resp *Response) {
	t.Helper()
	if resp.ExitCode == 0 && resp.ErrText == "" {
		t.Fatalf("want failure; stdout=%q stderr=%q", resp.Stdout, resp.Stderr)
	}
}

func assertContains(t *testing.T, got, want, label string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%s must contain %q; got:\n%s", label, want, got)
	}
}

func assertStdoutExact(t *testing.T, got string, lines ...string) {
	t.Helper()
	want := strings.Join(lines, "\n") + "\n"
	if got != want {
		t.Fatalf("stdout mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}
```
