# Scenario

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.TempDir = t.TempDir()
	req.UserHome = filepath.Join(req.TempDir, "home")
	req.AgentHome = filepath.Join(req.TempDir, ".agent-run")
	req.ProjectDir = filepath.Join(req.TempDir, "proj")
	req.FixedNow = time.Date(2026, 9, 2, 12, 0, 0, 0, time.Local)
	if err := os.MkdirAll(req.ProjectDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(req.UserHome, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(req.AgentHome, 0o755); err != nil {
		return err
	}
	return nil
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

func assertContains(t *testing.T, got, want, label string) {
	t.Helper()
	if !strings.Contains(got, want) {
		t.Fatalf("%s must contain %q; got:\n%s", label, want, got)
	}
}

func assertNotContains(t *testing.T, got, ban, label string) {
	t.Helper()
	if strings.Contains(got, ban) {
		t.Fatalf("%s must not contain %q; got:\n%s", label, ban, got)
	}
}

func expectedCodexSessionID(msg string) string {
	return "brainstorm-" + slugMsg(msg) + "-20260902-120000"
}

func slugMsg(msg string) string {
	s := strings.ToLower(strings.TrimSpace(msg))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if !prevDash {
			b.WriteByte('-')
			prevDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}
```
