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
	if strings.Contains(resp.Stdout, "\x1b[") {
		t.Fatalf("JSON must have no ANSI: %q", resp.Stdout)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(resp.Stdout), &doc); err != nil {
		t.Fatalf("json: %v\n%s", err, resp.Stdout)
	}
	path, _ := doc["path"].(string)
	if !strings.Contains(path, "rollout-") {
		t.Fatalf("path = %q", path)
	}
	if strings.HasPrefix(path, "~") {
		t.Fatalf("JSON path must be absolute: %q", path)
	}
	if doc["session_id"] != req.SessionID {
		t.Fatalf("session_id = %v", doc["session_id"])
	}
	if doc["file_active"] != false {
		t.Fatalf("file_active = %v", doc["file_active"])
	}
}
```
