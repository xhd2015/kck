# Scenario

**Feature**: kck codex send thin CLI

```
run.MainWith(Options{Args: ["codex","send",…], CodexHome, CodexSendOpts})
  -> sessions.RunSend
```

## Preconditions

- No live iTerm / ps / real CODEX_HOME.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/agent-pro/agent/codex/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

const fixtureKckSendSessionID = "019f283a-ffff-7fff-ffff-ffffffffff61"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.TempDir = t.TempDir()
	req.CodexHome = filepath.Join(req.TempDir, ".codex")
	req.ProjectDir = filepath.Join(req.TempDir, "proj")
	if err := os.MkdirAll(filepath.Join(req.CodexHome, "sessions"), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(req.ProjectDir, 0o755); err != nil {
		return err
	}
	if req.OpenFiles == nil {
		req.OpenFiles = map[int][]string{}
	}
	return nil
}

func writeKckSendSession(t *testing.T, req *Request) {
	t.Helper()
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fixtureKckSendSessionID
		req.SessionID = sessionID
	}
	dir := filepath.Join(req.CodexHome, "sessions", "2026", "08", "01")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, "rollout-2026-08-01T12-00-00-"+sessionID+".jsonl")
	body := `{"timestamp":"2026-08-01T12:00:00.000Z","type":"session_meta","payload":{"id":"` + sessionID + `","cwd":"` + req.ProjectDir + `"}}` + "\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func addLiveCodexSend(req *Request, pid int, tty string) {
	req.Procs = append(req.Procs, sessions.FocusProc{
		PID: pid, PPID: 1, TTY: tty, Cmd: "/usr/local/bin/codex",
	})
	req.OpenFiles[pid] = []string{
		"/Users/fixture/.codex/sessions/2026/08/01/rollout-2026-08-01T12-00-00-" + req.SessionID + ".jsonl",
	}
}

func oneITermTabSend() []iterm2.SessionRef {
	return []iterm2.SessionRef{
		{WindowID: "3", WindowName: "worktrees", TabIndex: 1, SessionID: "w2t1p0", TTY: "/dev/ttys148"},
	}
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
```
