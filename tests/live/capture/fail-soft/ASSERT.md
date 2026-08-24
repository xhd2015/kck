## Expected

- Exit 0 (soft-fail; MainWith returns nil).
- Stderr contains `warning:` (case-insensitive prefix ok).
- Stderr mentions capture failure (contains `live capture` or the error text).
- Stdout footer `0 sessions · 0 needs attention · 0 sendable`.
- No agent session ids from fixtures.

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
	assertWarningPrefix(t, resp.Stderr)
	low := strings.ToLower(resp.Stderr)
	if !strings.Contains(low, "live capture") && !strings.Contains(low, "iterm not running") {
		t.Fatalf("stderr should mention live capture failure; got %q", resp.Stderr)
	}
	assertContains(t, resp.Stdout, "0 sessions · 0 needs attention · 0 sendable", "footer")
	assertNotContains(t, resp.Stdout, "agent-sess-grok-1", "no live rows on fail")
}
```
