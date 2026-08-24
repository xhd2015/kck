## Expected

- Exit 0 (soft path).
- Session `dead-tty` still listed.
- LIVE shown as not alive (`no` token present in output).
- Optional `warning:` on stderr is allowed; must not hard-fail.

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
	assertContains(t, resp.Stdout, "dead-tty", "session")
	// Row should indicate not live: "no" appears (LIVE column).
	if !strings.Contains(strings.ToLower(resp.Stdout), "no") {
		t.Fatalf("want LIVE no for unreachable; stdout:\n%s", resp.Stdout)
	}
}
```
