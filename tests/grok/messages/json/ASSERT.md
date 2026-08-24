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
		SessionID     string `json:"session_id"`
		Total         int    `json:"total"`
		OffsetFromEnd int    `json:"offset_from_end"`
		Limit         int    `json:"limit"`
		Messages      []struct {
			Kind      string `json:"kind"`
			Text      string `json:"text"`
			Timestamp string `json:"timestamp"`
		} `json:"messages"`
	}
	if err := json.Unmarshal([]byte(resp.Stdout), &payload); err != nil {
		t.Fatalf("json: %v\n%s", err, resp.Stdout)
	}
	if payload.SessionID != req.SessionID {
		t.Fatalf("session_id=%q want %q", payload.SessionID, req.SessionID)
	}
	if payload.OffsetFromEnd != 2 || payload.Limit != 2 {
		t.Fatalf("offset/limit = %d/%d", payload.OffsetFromEnd, payload.Limit)
	}
	if payload.Total < 4 {
		t.Fatalf("total=%d", payload.Total)
	}
	if len(payload.Messages) != 2 {
		t.Fatalf("messages len=%d: %+v", len(payload.Messages), payload.Messages)
	}
	if payload.Messages[1].Kind != "response" || payload.Messages[1].Text != "a1" {
		t.Fatalf("want trailing a1 response: %+v", payload.Messages)
	}
	if payload.Messages[1].Timestamp == "" {
		t.Fatalf("want RFC3339 timestamp on response: %+v", payload.Messages[1])
	}
}
```
