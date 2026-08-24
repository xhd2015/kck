## Expected

- Exit 0.
- Stdout contains `agent-sess-grok-1`.
- Stdout contains ITERM token `w=42 t=3` (whitespace-flexible around cells ok;
  substring match).
- Stdout contains workspace `/ws/grok`.
- Header mentions SESSION_ID and ITERM and WORKSPACE (case-insensitive ok).

## Exit Code

0

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	_ = d
	_ = req
	assertNoError(t, err)
	assertSuccess(t, resp)
	out := resp.Stdout
	assertContains(t, out, "agent-sess-grok-1", "session id")
	assertContains(t, out, "w=42 t=3", "ITERM")
	assertContains(t, out, "/ws/grok", "workspace")
	assertContains(t, out, "grok", "runner")
	upper := strings.ToUpper(out)
	for _, col := range []string{"SESSION_ID", "ITERM", "WORKSPACE"} {
		if !strings.Contains(upper, col) {
			t.Fatalf("header missing %q; stdout:\n%s", col, out)
		}
	}
}
```
