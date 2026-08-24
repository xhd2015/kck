# Scenario

**Feature**: kck grok messages thin CLI

```go
import (
	"encoding/json"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

const fixtureKckMessagesSessionID = "019f283a-msgs-7ccc-cccc-cccccccccc01"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	req.TempDir = t.TempDir()
	req.GrokHome = filepath.Join(req.TempDir, ".grok")
	if err := os.MkdirAll(filepath.Join(req.GrokHome, "sessions"), 0o755); err != nil {
		return err
	}
	return nil
}

func writeKckMessagesSession(t *testing.T, req *Request, updates string) {
	t.Helper()
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fixtureKckMessagesSessionID
		req.SessionID = sessionID
	}
	cwd := "/tmp/kck-messages-proj"
	absCwd, err := filepath.Abs(cwd)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}
	dir := filepath.Join(req.GrokHome, "sessions", url.PathEscape(absCwd), sessionID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	summary := map[string]any{
		"info":              map[string]any{"id": sessionID, "cwd": absCwd},
		"generated_title":   "kck messages fixture",
		"created_at":        "2026-07-01T10:00:00.000Z",
		"updated_at":        "2026-07-01T11:00:00.000Z",
		"last_active_at":    "2026-07-01T11:00:00.000Z",
		"num_messages":      3,
		"num_chat_messages": 3,
		"current_model_id":  "grok-4",
	}
	body, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	body = append(body, '\n')
	if err := os.WriteFile(filepath.Join(dir, "summary.json"), body, 0o644); err != nil {
		t.Fatalf("write summary: %v", err)
	}
	if updates != "" {
		if !strings.HasSuffix(updates, "\n") {
			updates += "\n"
		}
		if err := os.WriteFile(filepath.Join(dir, "updates.jsonl"), []byte(updates), 0o644); err != nil {
			t.Fatalf("write updates: %v", err)
		}
	}
}

func fixtureMultiTurnUpdates() string {
	// Outer timestamp is unix seconds (converter normalizes to ms).
	return strings.Join([]string{
		`{"timestamp":1785493071,"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"u0"}}`,
		`{"timestamp":1785493072,"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"a0"}}`,
		`{"timestamp":1785493072,"sessionUpdate":"turn_completed"}`,
		`{"timestamp":1785493073,"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"u1"}}`,
		`{"timestamp":1785493074,"sessionUpdate":"agent_thought_chunk","content":{"type":"text","text":"th1"}}`,
		`{"timestamp":1785493075,"sessionUpdate":"tool_call","toolCallId":"t1","title":"run_terminal_command","rawInput":{"command":"echo hi","description":"say hi"},"_meta":{"x.ai/tool":{"name":"run_terminal_command","kind":"execute"}}}`,
		`{"timestamp":1785493076,"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"a1"}}`,
		`{"timestamp":1785493076,"sessionUpdate":"turn_completed"}`,
		`{"timestamp":1785493077,"sessionUpdate":"user_message_chunk","content":{"type":"text","text":"u2"}}`,
		`{"timestamp":1785493078,"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"a2"}}`,
		`{"timestamp":1785493078,"sessionUpdate":"turn_completed"}`,
	}, "\n")
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
