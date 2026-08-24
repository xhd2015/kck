## Expected

- Exit 0.
- Stdout is valid JSON (object or array containing sessions).
- No ANSI escapes.
- Both `s-new` and `s-old` appear in JSON text.
- Summary / counts reflect 2 sessions (flexible: top-level `sessions` array len 2
  or summary.total == 2).

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
	assertContains(t, resp.Stdout, "s-new", "json")
	assertContains(t, resp.Stdout, "s-old", "json")

	raw := strings.TrimSpace(resp.Stdout)
	var anyJSON any
	if err := json.Unmarshal([]byte(raw), &anyJSON); err != nil {
		t.Fatalf("stdout is not valid JSON: %v\n%s", err, resp.Stdout)
	}
	// Prefer object with sessions array.
	if m, ok := anyJSON.(map[string]any); ok {
		if sess, ok := m["sessions"].([]any); ok {
			if len(sess) != 2 {
				t.Fatalf("sessions len = %d, want 2", len(sess))
			}
		}
	}
}
```
