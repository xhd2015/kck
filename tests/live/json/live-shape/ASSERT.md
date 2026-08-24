## Expected

- Exit 0.
- Stdout is valid JSON object with `sessions` array and `summary` object.
- `sessions` length 1; includes `agent-sess-grok-1`.
- No ANSI escapes.
- Trailing newline.
- Summary total == 1; needs_attention == 1 for busy row; sendable == 0.

## Exit Code

0

```go
import (
	"encoding/json"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	assertNoANSI(t, resp.Stdout, "json stdout")
	assertTrailingNewline(t, resp.Stdout, "json stdout")
	assertContains(t, resp.Stdout, "agent-sess-grok-1", "session id in json")
	assertContains(t, resp.Stdout, "agent_session_id", "json field")
	assertContains(t, resp.Stdout, "agent_run", "json field")

	raw := strings.TrimSpace(resp.Stdout)
	var root map[string]any
	if err := json.Unmarshal([]byte(raw), &root); err != nil {
		t.Fatalf("stdout is not a JSON object: %v\n%s", err, resp.Stdout)
	}
	sess, ok := root["sessions"].([]any)
	if !ok {
		t.Fatalf("missing sessions array; got %T", root["sessions"])
	}
	if len(sess) != 1 {
		t.Fatalf("sessions len = %d, want 1", len(sess))
	}
	row, ok := sess[0].(map[string]any)
	if !ok {
		t.Fatalf("session[0] not object: %T", sess[0])
	}
	if sid, _ := row["agent_session_id"].(string); sid != "agent-sess-grok-1" {
		t.Fatalf("agent_session_id = %v, want agent-sess-grok-1", row["agent_session_id"])
	}
	// bare grok tree → agent_run false
	if ar, ok := row["agent_run"].(bool); !ok || ar {
		t.Fatalf("agent_run = %v, want false", row["agent_run"])
	}
	sum, ok := root["summary"].(map[string]any)
	if !ok {
		t.Fatalf("missing summary object; got %T", root["summary"])
	}
	// JSON numbers decode as float64
	if tot, _ := sum["total"].(float64); tot != 1 {
		t.Fatalf("summary.total = %v, want 1", sum["total"])
	}
	if na, _ := sum["needs_attention"].(float64); na != 1 {
		t.Fatalf("summary.needs_attention = %v, want 1", sum["needs_attention"])
	}
	if sn, _ := sum["sendable"].(float64); sn != 0 {
		t.Fatalf("summary.sendable = %v, want 0", sum["sendable"])
	}
}
```
