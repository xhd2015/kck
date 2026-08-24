## Expected

- Exit 0 (soft-fail).
- Stderr contains `warning:` (case-insensitive ok for prefix).
- Stdout does not show `w=42` (no successful resolve).
- Session still listed.

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
	assertContains(t, resp.Stdout, "iterm-s", "session")
	assertNotContains(t, resp.Stdout, "w=42", "failed list must not fill ITERM")
	if !strings.Contains(strings.ToLower(resp.Stderr), "warning:") {
		t.Fatalf("want warning: on stderr; got %q", resp.Stderr)
	}
}
```
