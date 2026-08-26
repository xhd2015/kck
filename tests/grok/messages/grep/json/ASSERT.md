## Expected

```go
import (
	"encoding/json"
	"strings"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	assertNoHarnessErr(t, err)
	assertSuccess(t, resp)
	if strings.ContainsRune(resp.Stdout, '\x1b') {
		t.Fatalf("JSON must not contain ANSI:\n%s", resp.Stdout)
	}
	var payload struct {
		SessionID string `json:"session_id"`
		Total     int    `json:"total"`
		Limit     int    `json:"limit"`
		Messages  []struct {
			Kind string `json:"kind"`
			Text string `json:"text"`
		} `json:"messages"`
	}
	if err := json.Unmarshal([]byte(resp.Stdout), &payload); err != nil {
		t.Fatalf("json: %v\n%s", err, resp.Stdout)
	}
	if payload.SessionID != req.SessionID {
		t.Fatalf("session_id=%q", payload.SessionID)
	}
	if payload.Total != 1 {
		t.Fatalf("total=%d want 1 (post-grep)", payload.Total)
	}
	if len(payload.Messages) != 1 || payload.Messages[0].Text != "a1" || payload.Messages[0].Kind != "response" {
		t.Fatalf("messages=%+v", payload.Messages)
	}
}
```
