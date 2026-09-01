# Scenario

# Shared setup

```go
import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/session"
)

const fixtureKckPromptsSessionID = "019f283a-prmt-7ccc-cccc-cccccccccc01"
const fixtureKckPromptsCWD = "/tmp/kck-prompts-proj"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.TempDir = t.TempDir()
	req.GrokHome = filepath.Join(req.TempDir, ".grok")
	if err := os.MkdirAll(filepath.Join(req.GrokHome, "sessions"), 0o755); err != nil {
		return err
	}
	return nil
}

func userChunkAt(text string, ts time.Time) string {
	line, _ := json.Marshal(map[string]any{
		"sessionUpdate": "user_message_chunk",
		"content":       map[string]any{"type": "text", "text": text},
		"timestamp":     ts.UnixMilli(),
	})
	return string(line)
}

func turnCompleted() string {
	return `{"sessionUpdate":"turn_completed"}`
}

func updatesJSONL(lines ...string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n") + "\n"
}

type promptSessionOpts struct {
	ID           string
	Title        string
	Kind         string // session_kind; empty → omit
	ParentID     string
	LastActiveAt time.Time
	Updates      string
}

func writePromptSession(t *testing.T, req *Request, o promptSessionOpts) string {
	t.Helper()
	if o.ID == "" {
		o.ID = fixtureKckPromptsSessionID
	}
	if o.Title == "" {
		o.Title = "kck prompts fixture"
	}
	if o.LastActiveAt.IsZero() {
		o.LastActiveAt = time.Date(2026, 8, 28, 12, 0, 0, 0, time.UTC)
	}
	absCwd, err := filepath.Abs(fixtureKckPromptsCWD)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	dir := filepath.Join(req.GrokHome, "sessions", url.PathEscape(absCwd), o.ID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	lastActive := o.LastActiveAt.UTC().Format("2006-01-02T15:04:05.000Z")
	summary := map[string]any{
		"info":              map[string]any{"id": o.ID, "cwd": absCwd},
		"generated_title":   o.Title,
		"created_at":        "2026-08-01T10:00:00.000Z",
		"updated_at":        lastActive,
		"last_active_at":    lastActive,
		"num_messages":      5,
		"num_chat_messages": 5,
	}
	if o.Kind != "" {
		summary["session_kind"] = o.Kind
	}
	if o.ParentID != "" {
		summary["parent_session_id"] = o.ParentID
	}
	body, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	body = append(body, '\n')
	if err := os.WriteFile(filepath.Join(dir, "summary.json"), body, 0o644); err != nil {
		t.Fatalf("write summary: %v", err)
	}
	if o.Updates != "" {
		u := o.Updates
		if !strings.HasSuffix(u, "\n") {
			u += "\n"
		}
		if err := os.WriteFile(filepath.Join(dir, "updates.jsonl"), []byte(u), 0o644); err != nil {
			t.Fatalf("write updates: %v", err)
		}
	}
	req.SessionID = o.ID
	return o.ID
}

// fiveUserPrompts builds p1..p5 chronologically with assistant separators.
func fiveUserPrompts() string {
	base := time.Date(2026, 8, 28, 11, 0, 0, 0, time.UTC)
	var lines []string
	for i := 1; i <= 5; i++ {
		lines = append(lines,
			userChunkAt(fmt.Sprintf("p%d", i), base.Add(time.Duration(i)*time.Minute)),
			turnCompleted(),
		)
	}
	return updatesJSONL(lines...)
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

func assertNotContains(t *testing.T, got, drop, label string) {
	t.Helper()
	if strings.Contains(got, drop) {
		t.Fatalf("%s must not contain %q; got:\n%s", label, drop, got)
	}
}
```
