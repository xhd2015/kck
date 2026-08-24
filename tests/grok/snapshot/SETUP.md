# Scenario

**Feature**: kck grok snapshot thin CLI

```
run.MainWith(Options{Args: ["grok","snapshot",…], GrokHome, GrokSnapshotOpts})
  -> sessions.RunSnapshot
```

## Preconditions

- No live iTerm / ps / real GROK_HOME.
- Session fixtures under `req.GrokHome`.

```go
import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/agent-pro/agent/grok/sessions"
	"github.com/xhd2015/doctest/session"
	"github.com/xhd2015/dot-pkgs/go-pkgs/shell/iterm2"
)

const fixtureKckSnapshotSessionID = "019f283a-ffff-7fff-ffff-ffffffffff01"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.TempDir = t.TempDir()
	req.GrokHome = filepath.Join(req.TempDir, ".grok")
	req.ProjectDir = filepath.Join(req.TempDir, "proj")
	if err := os.MkdirAll(filepath.Join(req.GrokHome, "sessions"), 0o755); err != nil {
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

func writeKckSnapshotSession(t *testing.T, req *Request) {
	t.Helper()
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fixtureKckSnapshotSessionID
		req.SessionID = sessionID
	}
	cwd := req.ProjectDir
	absKey, err := filepath.Abs(cwd)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	dir := filepath.Join(req.GrokHome, "sessions", url.PathEscape(absKey), sessionID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	summary := map[string]any{
		"info":            map[string]any{"id": sessionID, "cwd": cwd},
		"generated_title": "kck snapshot fixture",
		"created_at":      "2026-07-01T10:00:00.000Z",
		"updated_at":      "2026-07-01T11:00:00.000Z",
		"last_active_at":  "2026-07-01T11:00:00.000Z",
		"num_messages":    1,
	}
	body, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	body = append(body, '\n')
	if err := os.WriteFile(filepath.Join(dir, "summary.json"), body, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func addLiveGrok(req *Request, pid int, tty string) {
	req.Procs = append(req.Procs, sessions.FocusProc{
		PID: pid, PPID: 1, TTY: tty, Cmd: "/usr/local/bin/grok",
	})
	req.OpenFiles[pid] = []string{
		"/Users/fixture/.grok/sessions/%2Ftmp%2Fproj/" + req.SessionID + "/events.jsonl",
	}
}

func oneITermTab() []iterm2.SessionRef {
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

func assertNoContents(t *testing.T, resp *Response) {
	t.Helper()
	if len(resp.ContentsCalls) != 0 {
		t.Fatalf("ContentsCalls = %v, want none", resp.ContentsCalls)
	}
}
```
