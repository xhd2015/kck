# Scenario

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

const fixtureKckPickupSessionID = "019f283a-pickup-7ccc-cccc-cccccccccc01"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.TempDir = t.TempDir()
	req.UserHome = filepath.Join(req.TempDir, "home")
	req.GrokHome = filepath.Join(req.TempDir, ".grok")
	req.ProjectDir = filepath.Join(req.TempDir, "proj")
	req.CachePath = filepath.Join(req.UserHome, ".cache", "kck-pickup-a-session", "SKILL.md")
	req.SessionID = fixtureKckPickupSessionID
	req.SkillBody = "# kck-pickup-a-session fixture\n"
	if err := os.MkdirAll(req.ProjectDir, 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(req.UserHome, 0o755); err != nil {
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

func assertErrorPrefix(t *testing.T, resp *Response) {
	t.Helper()
	blob := resp.Stderr + "\n" + resp.ErrText
	if !strings.Contains(blob, "Error:") {
		t.Fatalf("want Error: prefix; stderr=%q err=%q", resp.Stderr, resp.ErrText)
	}
}

func draftSkillDisplay(req *Request) string {
	return "~/.cache/kck-pickup-a-session/SKILL.md"
}
```
