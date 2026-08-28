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
		t.Fatalf("ANSI: %q", resp.Stdout)
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(resp.Stdout), &doc); err != nil {
		t.Fatalf("%v\n%s", err, resp.Stdout)
	}
	if doc["session_id"] != req.SessionID {
		t.Fatalf("%v", doc)
	}
}
```
