# Scenario

**Feature**: kck codex resolve thin CLI

```
run.MainWith(Options{Args: ["codex","resolve",…], CodexHome, CodexResolveOpts})
  -> sessions.RunResolve
```

## Preconditions

- No live iTerm / ps / real CODEX_HOME.
- Session id from open-file paths only (never cmdline `resume` flags).

```go
import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

const (
	fixtureSessionID    = "019f283b-aaaa-7aaa-aaaa-aaaaaaaaaaaa"
	fixtureTabSessionID = "019f283b-dddd-7ddd-dddd-dddddddddddd"

	pidCodex = 4242
	pidBash  = 5000
	pidStart = 6000

	pidTabCodex1 = 8100
	pidTabCodex2 = 8200
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	if req.TempDir == "" {
		req.TempDir = t.TempDir()
	}
	if req.CodexHome == "" {
		req.CodexHome = filepath.Join(req.TempDir, ".codex")
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

func codexSessionPath(sessionID string) string {
	return "/Users/fixture/.codex/sessions/2026/08/01/rollout-2026-08-01T12-00-00-" + sessionID + ".jsonl"
}

func defaultAncestorChain() []FixtureProc {
	return []FixtureProc{
		{PID: pidCodex, PPID: 1, Cmd: "/usr/local/bin/codex"},
		{PID: pidBash, PPID: pidCodex, Cmd: "/bin/bash"},
		{PID: pidStart, PPID: pidBash, Cmd: "/usr/local/bin/agent-pro"},
	}
}

func seedHit(req *Request, sessionID string, codexPID int) {
	if len(req.Procs) == 0 {
		req.Procs = defaultAncestorChain()
	}
	req.OpenFiles[codexPID] = []string{codexSessionPath(sessionID)}
}

// seedTabWindow installs a 3-tab window; current is tab 1 (/dev/ttys101).
// Tab 2 hosts fixtureTabSessionID on /dev/ttys102; tab 3 is bash-only.
func seedTabWindow(req *Request) {
	req.CurrentSessionID = "w0t1p0:CURRENT-TAB-UUID"
	req.ControllingTTY = "/dev/ttys101"
	req.ITerm = []iterm2.SessionRef{
		{WindowID: "100", WindowName: "work", TabIndex: 1, SessionID: "w0t1p0:CURRENT-TAB-UUID", TTY: "/dev/ttys101", Name: "current"},
		{WindowID: "100", WindowName: "work", TabIndex: 2, SessionID: "w0t2p0:TAB2-UUID", TTY: "/dev/ttys102", Name: "codex-tab"},
		{WindowID: "100", WindowName: "work", TabIndex: 3, SessionID: "w0t3p0:TAB3-UUID", TTY: "/dev/ttys103", Name: "bash-only"},
	}
	req.FocusProcs = []sessions.FocusProc{
		{PID: pidTabCodex1, PPID: 1, TTY: "/dev/ttys101", Cmd: "/usr/local/bin/bash"},
		{PID: pidTabCodex2, PPID: 1, TTY: "/dev/ttys102", Cmd: "/usr/local/bin/codex"},
		{PID: 9100, PPID: 1, TTY: "/dev/ttys103", Cmd: "/bin/bash"},
	}
	req.OpenFiles[pidTabCodex2] = []string{codexSessionPath(fixtureTabSessionID)}
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
