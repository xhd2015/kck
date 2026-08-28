# Scenario

**Feature**: kck codex info thin CLI

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

const fixtureKckInfoSessionID = "019f283a-cccc-7ccc-cccc-cccccccccc81"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.TempDir = t.TempDir()
	req.CodexHome = filepath.Join(req.TempDir, ".codex")
	if err := os.MkdirAll(filepath.Join(req.CodexHome, "sessions"), 0o755); err != nil {
		return err
	}
	if req.OpenFiles == nil {
		req.OpenFiles = map[int][]string{}
	}
	return nil
}

func writeKckInfoSession(t *testing.T, req *Request) {
	t.Helper()
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fixtureKckInfoSessionID
		req.SessionID = sessionID
	}
	dir := filepath.Join(req.CodexHome, "sessions", "2026", "08", "01")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, "rollout-2026-08-01T12-00-00-"+sessionID+".jsonl")
	// Include a short user message so Title is populated.
	body := `{"timestamp":"2026-08-01T12:00:00.000Z","type":"session_meta","payload":{"id":"` + sessionID + `","cwd":"/tmp/kck-info-proj"}}` + "\n" +
		`{"timestamp":"2026-08-01T12:00:01.000Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"kck info fixture"}]}}` + "\n"
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}

func codexOpenPath(sessionID string) string {
	return "/Users/fixture/.codex/sessions/2026/08/01/rollout-2026-08-01T12-00-00-" + sessionID + ".jsonl"
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
